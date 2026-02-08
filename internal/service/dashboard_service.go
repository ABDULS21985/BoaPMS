package service

import (
	"context"
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// dashboardService aggregates dashboard statistics for a staff member's
// active review period: feedback request stats, performance points,
// work product counts, detailed work product breakdowns, and the
// full score-card computation.
//
// Mirrors .NET methods:
//   - GetRequestStatistics
//   - GetStaffPerformanceStatistics
//   - GetStaffWorkProductsStatistics
//   - GetStaffWorkProductsDetailsStatistics
//   - GetStaffPerformanceScoreCardStatistics
// ---------------------------------------------------------------------------

type dashboardService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	// Back-reference to parent for shared helpers and peer services.
	parent *performanceManagementService

	// Repositories
	workProductRepo         *repository.PMSRepository[performance.WorkProduct]
	feedbackLogRepo         *repository.PMSRepository[performance.FeedbackRequestLog]
	competencyFeedbackRepo  *repository.PMSRepository[performance.CompetencyReviewFeedback]
	competencyReviewerRepo  *repository.PMSRepository[performance.CompetencyReviewer]
	categoryDefRepo         *repository.PMSRepository[performance.CategoryDefinition]
	competencyGapClosureRepo *repository.PMSRepository[performance.CompetencyGapClosure]
	pmsCompetencyRepo       *repository.PMSRepository[performance.PmsCompetency]
	periodScoreRepo         *repository.PMSRepository[performance.PeriodScore]
	reviewPeriodRepo        *repository.PMSRepository[performance.PerformanceReviewPeriod]
}

func newDashboardService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *dashboardService {
	return &dashboardService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "dashboard").Logger(),
		parent: parent,

		workProductRepo:          repository.NewPMSRepository[performance.WorkProduct](db),
		feedbackLogRepo:          repository.NewPMSRepository[performance.FeedbackRequestLog](db),
		competencyFeedbackRepo:   repository.NewPMSRepository[performance.CompetencyReviewFeedback](db),
		competencyReviewerRepo:   repository.NewPMSRepository[performance.CompetencyReviewer](db),
		categoryDefRepo:          repository.NewPMSRepository[performance.CategoryDefinition](db),
		competencyGapClosureRepo: repository.NewPMSRepository[performance.CompetencyGapClosure](db),
		pmsCompetencyRepo:        repository.NewPMSRepository[performance.PmsCompetency](db),
		periodScoreRepo:          repository.NewPMSRepository[performance.PeriodScore](db),
		reviewPeriodRepo:         repository.NewPMSRepository[performance.PerformanceReviewPeriod](db),
	}
}

// =========================================================================
// GetDashboardStats is the top-level method delegated from the main service.
// It computes the score-card for the staff member's active review period.
// =========================================================================

func (d *dashboardService) GetDashboardStats(ctx context.Context, staffID string) (interface{}, error) {
	// Delegate to the score-card method using the active review period.
	reviewPeriod, err := d.parent.getStaffActiveReviewPeriod(ctx, staffID)
	if err != nil {
		d.log.Error().Err(err).Str("staffID", staffID).Msg("no active review period for dashboard")
		return performance.StaffScoreCardResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{
				HasError: true,
				Message:  "no active review period found",
			},
		}, err
	}
	return d.GetStaffPerformanceScoreCardStatistics(ctx, staffID, reviewPeriod.PeriodID)
}

// =========================================================================
// GetRequestStatistics – Feedback request SLA dashboard.
// Mirrors .NET GetRequestStatistics.
// =========================================================================

func (d *dashboardService) GetRequestStatistics(ctx context.Context, staffID string) (performance.FeedbackRequestDashboardResponseVm, error) {
	resp := performance.FeedbackRequestDashboardResponseVm{
		StaffID: staffID,
	}
	resp.Message = "an error occurred"

	reviewPeriod, err := d.parent.getStaffActiveReviewPeriod(ctx, staffID)
	if err != nil {
		d.log.Error().Err(err).Msg("GetRequestStatistics: no active review period")
		resp.HasError = true
		return resp, err
	}

	requestSLA, pms360SLA := d.parent.getSLAConfig(ctx)
	resp.ReviewPeriodID = reviewPeriod.PeriodID
	reviewPeriodMaxPoints := reviewPeriod.MaxPoints

	// Get all feedback requests for this staff in this review period
	var requests []performance.FeedbackRequestLog
	d.db.WithContext(ctx).
		Where("assigned_staff_id = ? AND review_period_id = ?", staffID, reviewPeriod.PeriodID).
		Find(&requests)

	now := time.Now()

	resp.CompletedRequests = countWhere(requests, func(r performance.FeedbackRequestLog) bool {
		return r.TimeCompleted != nil
	})

	resp.PendingRequests = countWhere(requests, func(r performance.FeedbackRequestLog) bool {
		return r.RecordStatus == enums.StatusActive.String()
	})

	// Completed overdue: non-360 requests
	completedOverdueNon360 := filterRequests(requests, func(r performance.FeedbackRequestLog) bool {
		return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
			isClosedOrCompletedOrBreached(r.RecordStatus) &&
			r.HasSLA &&
			r.TimeCompleted != nil &&
			r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(requestSLA)
	})
	completedOverdueNon360Count := d.parent.countSLABreachedRequests(completedOverdueNon360, staffID, requestSLA)

	// Completed overdue: 360 requests
	completedOverdue360 := filterRequests(requests, func(r performance.FeedbackRequestLog) bool {
		return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
			isClosedOrCompletedOrBreached(r.RecordStatus) &&
			r.HasSLA &&
			r.TimeCompleted != nil &&
			r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(pms360SLA)
	})
	completedOverdue360Count := d.parent.countSLABreachedRequests(completedOverdue360, staffID, pms360SLA)
	resp.CompletedOverdueRequests = completedOverdueNon360Count + completedOverdue360Count

	// Pending overdue: non-360 requests
	pendingOverdueNon360Count := countWhere(requests, func(r performance.FeedbackRequestLog) bool {
		return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
			r.RecordStatus == enums.StatusActive.String() &&
			r.HasSLA &&
			r.TimeInitiated.Add(time.Duration(requestSLA)*time.Hour).Before(now)
	})

	// Pending overdue: 360 requests
	pendingOverdue360Count := countWhere(requests, func(r performance.FeedbackRequestLog) bool {
		return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
			r.RecordStatus == enums.StatusActive.String() &&
			r.HasSLA &&
			r.TimeInitiated.Add(time.Duration(pms360SLA)*time.Hour).Before(now)
	})
	resp.PendingOverdueRequests = pendingOverdueNon360Count + pendingOverdue360Count

	resp.BreachedRequests = resp.CompletedOverdueRequests + resp.PendingOverdueRequests

	deducted := float64(resp.CompletedOverdueRequests + resp.PendingOverdueRequests)
	if deducted > reviewPeriodMaxPoints {
		deducted = reviewPeriodMaxPoints
	}
	resp.DeductedPoints = deducted

	// Pending 360 feedbacks to treat
	var feedbackReviewers []performance.CompetencyReviewer
	d.db.WithContext(ctx).
		Where("review_staff_id = ? AND created_at >= ? AND created_at <= ?",
			staffID, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Find(&feedbackReviewers)

	resp.Pending360FeedbacksToTreat = countWhere(feedbackReviewers, func(r performance.CompetencyReviewer) bool {
		return r.RecordStatus == enums.StatusActive.String()
	})

	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetStaffPerformanceStatistics – Points dashboard.
// Mirrors .NET GetStaffPerformanceStatistics.
// =========================================================================

func (d *dashboardService) GetStaffPerformanceStatistics(ctx context.Context, staffID string) (performance.ReviewPeriodPointsDashboardResponseVm, error) {
	resp := performance.ReviewPeriodPointsDashboardResponseVm{
		StaffID: staffID,
	}
	resp.Message = "an error occurred"

	requestSLA, pms360SLA := d.parent.getSLAConfig(ctx)

	reviewPeriod, err := d.parent.getStaffActiveReviewPeriod(ctx, staffID)
	if err != nil {
		d.log.Error().Err(err).Msg("GetStaffPerformanceStatistics: no active review period")
		resp.HasError = true
		return resp, err
	}

	resp.ReviewPeriodID = reviewPeriod.PeriodID
	resp.MaxPoints = reviewPeriod.MaxPoints
	reviewPeriodMaxPoints := reviewPeriod.MaxPoints

	excluded := excludedStatuses()

	// Get work products within the review period date range
	var workProducts []performance.WorkProduct
	d.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
			staffID, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Find(&workProducts)

	var livingTheValuesPoints float64
	var totalGapClosureScore float64
	var feedbackRatingScore float64
	var workProductPoints float64
	var deductedPoints float64

	// Get staff info for grade lookup
	// NOTE: ERP employee lookup for grade assignment uses peer service.
	// For now we compute based on available data, skipping grade-group
	// resolution that requires ERP integration.
	staffJobGradeGroupID := 0

	// Get competency review feedback (360 ratings)
	var ratings []performance.CompetencyReviewFeedback
	d.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ? AND created_at >= ? AND created_at <= ?",
			staffID, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Preload("CompetencyReviewers").
		Find(&ratings)

	competencyCategoryDefined := false
	if len(ratings) > 0 {
		rating := ratings[0]

		// Find category definition with PMS competency mapping
		var categoryDef performance.CategoryDefinition
		err := d.db.WithContext(ctx).
			Joins("JOIN pms.objective_categories oc ON oc.objective_category_id = pms.category_definitions.objective_category_id").
			Joins("JOIN pms.pms_competencies pc ON pc.object_category_id = oc.objective_category_id").
			Where("pms.category_definitions.review_period_id = ? AND pms.category_definitions.grade_group_id = ?",
				resp.ReviewPeriodID, staffJobGradeGroupID).
			First(&categoryDef).Error

		if err == nil {
			competencyCategoryDefined = true
			maxPts := categoryDef.MaxPoints
			countAll := len(rating.CompetencyReviewers)
			if countAll > 0 {
				var totalRating float64
				for _, r := range rating.CompetencyReviewers {
					totalRating += r.FinalRating
				}
				scoreAverage := totalRating / float64(countAll)
				feedbackRatingScore = maxPts * (scoreAverage / 100)
			}
		}
	}

	workProductPoints = sumFloat64(workProducts, func(wp performance.WorkProduct) float64 {
		return wp.FinalScore
	})

	if !competencyCategoryDefined {
		// Fallback: sum FinalScore from CompetencyReviewFeedback
		var ltvFeedbacks []performance.CompetencyReviewFeedback
		d.db.WithContext(ctx).
			Where("staff_id = ? AND record_status NOT IN ? AND created_at >= ? AND created_at <= ?",
				staffID, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
			Find(&ltvFeedbacks)
		for _, ltv := range ltvFeedbacks {
			livingTheValuesPoints += ltv.FinalScore
		}
	} else {
		livingTheValuesPoints = feedbackRatingScore
	}

	// Competency gap closure score
	var gapClosure performance.CompetencyGapClosure
	err = d.db.WithContext(ctx).
		Where("review_period_id = ? AND staff_id = ?", resp.ReviewPeriodID, staffID).
		Preload("ObjectiveCategory").
		First(&gapClosure).Error
	if err == nil {
		var catDef performance.CategoryDefinition
		catErr := d.db.WithContext(ctx).
			Where("review_period_id = ? AND grade_group_id = ? AND objective_category_id = ?",
				resp.ReviewPeriodID, staffJobGradeGroupID, gapClosure.ObjectiveCategoryID).
			First(&catDef).Error
		if catErr == nil {
			totalGapClosureScore = catDef.MaxPoints * (gapClosure.FinalScore / 100)
		}
	}

	// SLA deduction calculation
	var allRequests []performance.FeedbackRequestLog
	d.db.WithContext(ctx).
		Where("assigned_staff_id = ? AND time_initiated >= ? AND time_initiated <= ?",
			staffID, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Find(&allRequests)

	if len(allRequests) > 0 {
		now := time.Now()

		// Completed overdue: non-360
		completedOverdueList := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
				isClosedOrCompletedOrBreached(r.RecordStatus) &&
				r.HasSLA && r.TimeCompleted != nil &&
				r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(requestSLA)
		})
		// Completed overdue: 360
		completedOverdue360List := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
				isClosedOrCompletedOrBreached(r.RecordStatus) &&
				r.HasSLA && r.TimeCompleted != nil &&
				r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(pms360SLA)
		})
		completedOverdue := d.parent.countSLABreachedRequests(completedOverdueList, staffID, requestSLA) +
			d.parent.countSLABreachedRequests(completedOverdue360List, staffID, pms360SLA)

		// Pending overdue: non-360
		pendingOverdueList := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
				r.RecordStatus == enums.StatusActive.String() &&
				r.HasSLA && r.TimeInitiated.Add(time.Duration(requestSLA)*time.Hour).Before(now)
		})
		// Pending overdue: 360
		pendingOverdue360List := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
				r.RecordStatus == enums.StatusActive.String() &&
				r.HasSLA && r.TimeInitiated.Add(time.Duration(pms360SLA)*time.Hour).Before(now)
		})
		pendingOverdue := d.parent.countSLABreachedRequests(pendingOverdueList, staffID, requestSLA) +
			d.parent.countSLABreachedRequests(pendingOverdue360List, staffID, pms360SLA)

		deductedPoints = float64(completedOverdue + pendingOverdue)
		if deductedPoints > reviewPeriodMaxPoints {
			deductedPoints = reviewPeriodMaxPoints
		}
	}

	resp.DeductedPoints = deductedPoints
	accumulatedPoints := workProductPoints + livingTheValuesPoints + totalGapClosureScore
	actualPoints := accumulatedPoints - deductedPoints
	if actualPoints < 0 {
		actualPoints = 0
	}

	resp.AccumulatedPoints = accumulatedPoints
	resp.ActualPoints = actualPoints
	resp.Message = "operation completed successfully"

	return resp, nil
}

// =========================================================================
// GetStaffWorkProductsStatistics – Work product count summary.
// Mirrors .NET GetStaffWorkProductsStatistics.
// =========================================================================

func (d *dashboardService) GetStaffWorkProductsStatistics(ctx context.Context, staffID string) (performance.ReviewPeriodWorkProductDashboardResponseVm, error) {
	resp := performance.ReviewPeriodWorkProductDashboardResponseVm{
		StaffID: staffID,
	}
	resp.Message = "an error occurred"

	reviewPeriod, err := d.parent.getStaffActiveReviewPeriod(ctx, staffID)
	if err != nil {
		d.log.Error().Err(err).Msg("GetStaffWorkProductsStatistics: no active review period")
		resp.HasError = true
		return resp, err
	}

	resp.ReviewPeriodID = reviewPeriod.PeriodID
	excluded := excludedStatuses()

	var workProducts []performance.WorkProduct
	d.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
			staffID, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Preload("WorkProductTasks").
		Find(&workProducts)

	// Count tasks across all work products (excluding excluded statuses)
	for _, wp := range workProducts {
		for _, t := range wp.WorkProductTasks {
			if !isInStatuses(t.RecordStatus, excluded) {
				resp.TotalWorkProductTasks++
			}
		}
	}

	resp.NoAllWorkProducts = len(workProducts)
	resp.NoActiveWorkProducts = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusActive.String()
	})
	resp.NoWorkProductsAwaitingEvaluation = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusAwaitingEvaluation.String()
	})
	resp.NoWorkProductsClosed = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusClosed.String()
	})
	resp.NoWorkProductsPendingApproval = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusPendingApproval.String()
	})

	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetStaffWorkProductsDetailsStatistics – Detailed work product breakdown.
// Mirrors .NET GetStaffWorkProductsDetailsStatistics.
// =========================================================================

func (d *dashboardService) GetStaffWorkProductsDetailsStatistics(ctx context.Context, staffID string) (performance.ReviewPeriodWorkProductDetailsDashboardResponseVm, error) {
	resp := performance.ReviewPeriodWorkProductDetailsDashboardResponseVm{
		StaffID: staffID,
	}
	resp.Message = "an error occurred"

	reviewPeriod, err := d.parent.getStaffActiveReviewPeriod(ctx, staffID)
	if err != nil {
		d.log.Error().Err(err).Msg("GetStaffWorkProductsDetailsStatistics: no active review period")
		resp.HasError = true
		return resp, err
	}

	resp.ReviewPeriodID = reviewPeriod.PeriodID
	excluded := excludedStatuses()

	var workProducts []performance.WorkProduct
	d.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
			staffID, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
		Preload("WorkProductTasks").
		Order("end_date ASC").
		Find(&workProducts)

	var dashDetails []performance.WorkProductDashDetails

	for _, wp := range workProducts {
		detail := performance.WorkProductDashDetails{
			BaseWorkFlowVm: performance.BaseWorkFlowVm{
				ID:              wp.ID,
				RecordStatus:    wp.RecordStatus,
				CreatedAt:       wp.CreatedAt,
				Status:          wp.Status,
				UpdatedAt:       wp.UpdatedAt,
				CreatedBy:       wp.CreatedBy,
				UpdatedBy:       wp.UpdatedBy,
				IsActive:        wp.IsActive,
				ApprovedBy:      wp.ApprovedBy,
				DateApproved:    wp.DateApproved,
				IsApproved:      wp.IsApproved,
				IsRejected:      wp.IsRejected,
				RejectedBy:      wp.RejectedBy,
				RejectionReason: wp.RejectionReason,
				DateRejected:    wp.DateRejected,
			},
			WorkProductID:     wp.WorkProductID,
			Name:              wp.Name,
			Description:       wp.Description,
			MaxPoint:          wp.MaxPoint,
			WorkProductType:   int(wp.WorkProductType),
			IsSelfCreated:     wp.IsSelfCreated,
			StaffID:           wp.StaffID,
			AcceptanceComment: wp.AcceptanceComment,
			StartDate:         wp.StartDate,
			EndDate:           wp.EndDate,
			Deliverables:      wp.Deliverables,
			FinalScore:        wp.FinalScore,
			NoReturned:        wp.NoReturned,
			ApproverComment:   wp.ApproverComment,
		}

		if wp.CompletionDate != nil {
			detail.CompletionDate = *wp.CompletionDate
		}

		// Map tasks
		var taskList []performance.WorkProductTaskDetail
		for _, t := range wp.WorkProductTasks {
			td := performance.WorkProductTaskDetail{
				BaseEntityVm: performance.BaseEntityVm{
					ID:           t.ID,
					RecordStatus: t.RecordStatus,
					CreatedAt:    t.CreatedAt,
					UpdatedAt:    t.UpdatedAt,
					CreatedBy:    t.CreatedBy,
					UpdatedBy:    t.UpdatedBy,
					IsActive:     t.IsActive,
				},
				WorkProductTaskID: t.WorkProductTaskID,
				Name:              t.Name,
				Description:       t.Description,
				StartDate:         t.StartDate,
				EndDate:           t.EndDate,
				WorkProductID:     t.WorkProductID,
			}
			if t.CompletionDate != nil {
				td.CompletionDate = *t.CompletionDate
			}
			taskList = append(taskList, td)
		}

		detail.WorkProductTasks = taskList
		detail.TotalTasks = len(taskList)
		detail.TasksCompleted = countWhere(taskList, func(t performance.WorkProductTaskDetail) bool {
			return t.RecordStatus == enums.StatusCompleted.String()
		})

		if detail.TotalTasks > 0 {
			detail.PercentageTaskCompletion = 100.0 * float64(detail.TasksCompleted) / float64(detail.TotalTasks)
		}

		dashDetails = append(dashDetails, detail)
	}

	// Count tasks (excluding excluded statuses)
	for _, wp := range workProducts {
		for _, t := range wp.WorkProductTasks {
			if !isInStatuses(t.RecordStatus, excluded) {
				resp.TotalWorkProductTasks++
			}
		}
	}

	resp.WorkProducts = dashDetails
	resp.NoAllWorkProducts = len(workProducts)
	resp.NoActiveWorkProducts = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusActive.String()
	})
	resp.NoWorkProductsAwaitingEvaluation = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusAwaitingEvaluation.String()
	})
	resp.NoWorkProductsClosed = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusClosed.String()
	})
	resp.NoWorkProductsPendingApproval = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return wp.RecordStatus == enums.StatusPendingApproval.String()
	})

	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetStaffPerformanceScoreCardStatistics – Full score card.
// Mirrors .NET GetStaffPerformanceScoreCardStatistics.
// =========================================================================

func (d *dashboardService) GetStaffPerformanceScoreCardStatistics(ctx context.Context, staffID, reviewPeriodID string) (performance.StaffScoreCardResponseVm, error) {
	resp := performance.StaffScoreCardResponseVm{}
	resp.Message = "an error occurred"

	requestSLA, pms360SLA := d.parent.getSLAConfig(ctx)

	// Look up the review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := d.db.WithContext(ctx).
		Where("period_id = ?", reviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		d.log.Error().Err(err).Str("reviewPeriodID", reviewPeriodID).Msg("review period not found")
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, err
	}

	scoreCard := performance.StaffScoreCardDetails{
		ReviewPeriodID:        reviewPeriod.PeriodID,
		ReviewPeriod:          reviewPeriod.Name,
		ReviewPeriodShortName: reviewPeriod.ShortName,
		MaxPoints:             reviewPeriod.MaxPoints,
		Year:                  reviewPeriod.Year,
	}
	if scoreCard.ReviewPeriodShortName == "" {
		scoreCard.ReviewPeriodShortName = reviewPeriod.Name
	}

	startDate := reviewPeriod.StartDate
	endDate := reviewPeriod.EndDate
	reviewPeriodMaxPoints := reviewPeriod.MaxPoints

	excluded := excludedStatuses()

	// Fetch work products
	var workProducts []performance.WorkProduct
	d.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
			staffID, excluded, startDate, endDate).
		Find(&workProducts)

	// NOTE: Staff grade lookup (ERP integration) would normally resolve the
	// staff grade group ID. We default to 0 pending ERP service integration.
	staffJobGradeGroupID := 0

	// Get 360 competency review ratings
	var allRatings []performance.CompetencyReviewFeedback
	d.db.WithContext(ctx).
		Where("staff_id = ? AND review_period_id = ? AND record_status NOT IN ?",
			staffID, reviewPeriodID, excluded).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Preload("CompetencyReviewers.CompetencyReviewerRatings.PmsCompetency").
		Find(&allRatings)

	var feedbackRatingScore float64
	var workProductPoints float64
	var livingTheValuesPoints float64
	var totalGapClosureScore float64

	competencyCategoryDefined := false

	if len(allRatings) > 0 {
		rating := allRatings[0]

		// Find category definitions with PMS competency mapping
		var categoryDefs []performance.CategoryDefinition
		d.db.WithContext(ctx).
			Preload("Category").
			Joins("JOIN pms.objective_categories oc ON oc.objective_category_id = pms.category_definitions.objective_category_id").
			Joins("JOIN pms.pms_competencies pc ON pc.object_category_id = oc.objective_category_id").
			Where("pms.category_definitions.review_period_id = ? AND pms.category_definitions.grade_group_id = ?",
				reviewPeriodID, staffJobGradeGroupID).
			Find(&categoryDefs)

		if len(categoryDefs) > 0 {
			competencyCategoryDefined = true
		}

		// Calculate feedback rating score per category
		for _, cat := range categoryDefs {
			maxPts := cat.MaxPoints

			// Get reviewers that have ratings matching this category
			var catRatings []float64
			for _, reviewer := range rating.CompetencyReviewers {
				for _, rr := range reviewer.CompetencyReviewerRatings {
					if rr.PmsCompetency != nil && rr.PmsCompetency.ObjectCategoryID == cat.ObjectiveCategoryID {
						catRatings = append(catRatings, rr.Rating)
					}
				}
			}
			if len(catRatings) > 0 {
				scoreAverage := average(catRatings)
				finalRating := maxPts * (scoreAverage / 100)
				feedbackRatingScore += finalRating
			}
		}

		// Build LTV ratings per competency (average across all reviewers)
		type ratingKey struct {
			StaffID         string
			ReviewPeriodID  string
			PmsCompetencyID string
		}
		ratingMap := make(map[ratingKey][]float64)
		competencyNames := make(map[string]string)

		for _, reviewer := range rating.CompetencyReviewers {
			for _, rr := range reviewer.CompetencyReviewerRatings {
				key := ratingKey{
					StaffID:         staffID,
					ReviewPeriodID:  rating.ReviewPeriodID,
					PmsCompetencyID: rr.PmsCompetencyID,
				}
				ratingMap[key] = append(ratingMap[key], rr.Rating)
				if rr.PmsCompetency != nil {
					competencyNames[rr.PmsCompetencyID] = rr.PmsCompetency.Name
				}
			}
		}

		var ltvDetails []performance.StaffLivingTheValueRatingsDetails
		for key, scores := range ratingMap {
			ltvDetails = append(ltvDetails, performance.StaffLivingTheValueRatingsDetails{
				StaffID:         key.StaffID,
				ReviewPeriodID:  key.ReviewPeriodID,
				PmsCompetencyID: key.PmsCompetencyID,
				PmsCompetency:   competencyNames[key.PmsCompetencyID],
				RatingScore:     average(scores),
			})
		}
		scoreCard.PmsCompetencies = ltvDetails

		// Calculate per-category competency scores
		pmsCompetencyCat := make(map[string]float64)
		for _, cat := range categoryDefs {
			maxCategoryPoints := cat.MaxPoints
			var catScores []float64
			for _, reviewer := range rating.CompetencyReviewers {
				for _, rr := range reviewer.CompetencyReviewerRatings {
					if rr.PmsCompetency != nil && rr.PmsCompetency.ObjectCategoryID == cat.ObjectiveCategoryID {
						catScores = append(catScores, rr.Rating)
					}
				}
			}
			var catFeedbackRating float64
			if len(catScores) > 0 {
				scoreAvg := average(catScores)
				catFeedbackRating = maxCategoryPoints * (scoreAvg / 100)
			}
			categoryName := ""
			if cat.Category != nil {
				categoryName = cat.Category.Name
			}
			pmsCompetencyCat[categoryName] = catFeedbackRating
		}
		scoreCard.PmsCompetencyCategory = pmsCompetencyCat
	}

	// Work products stats
	if len(workProducts) > 0 {
		workProductPoints = sumFloat64(workProducts, func(wp performance.WorkProduct) float64 {
			return wp.FinalScore
		})

		now := time.Now()
		completedStatuses := []string{
			enums.StatusCompleted.String(),
			enums.StatusAwaitingEvaluation.String(),
			enums.StatusClosed.String(),
		}

		onTimeCompleted := countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return isInStatuses(wp.RecordStatus, completedStatuses) &&
				wp.CompletionDate != nil && !wp.CompletionDate.After(wp.EndDate)
		})
		completedCount := countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return isInStatuses(wp.RecordStatus, completedStatuses)
		})
		overdueCount := countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return (wp.RecordStatus == enums.StatusActive.String() && wp.EndDate.Before(now)) ||
				(isInStatuses(wp.RecordStatus, completedStatuses) && wp.CompletionDate != nil && wp.CompletionDate.After(wp.EndDate))
		})

		scoreCard.TotalWorkProducts = len(workProducts)
		scoreCard.TotalWorkProductsCompletedOnSchedule = onTimeCompleted
		scoreCard.TotalWorkProductsBehindSchedule = overdueCount

		if completedCount > 0 {
			scoreCard.PercentageWorkProductsCompletion = 100.0 * float64(completedCount) / float64(len(workProducts))
		}
	}

	// Living the values points
	if !competencyCategoryDefined {
		var ltvFeedbacks []performance.CompetencyReviewFeedback
		d.db.WithContext(ctx).
			Where("staff_id = ? AND record_status NOT IN ? AND review_period_id = ?",
				staffID, excluded, reviewPeriodID).
			Find(&ltvFeedbacks)
		for _, ltv := range ltvFeedbacks {
			livingTheValuesPoints += ltv.FinalScore
		}
	} else {
		livingTheValuesPoints = feedbackRatingScore
	}

	// Competency gap closure score
	var gapClosure performance.CompetencyGapClosure
	err := d.db.WithContext(ctx).
		Where("review_period_id = ? AND staff_id = ?", reviewPeriodID, staffID).
		Preload("ObjectiveCategory").
		First(&gapClosure).Error
	if err == nil {
		var catDef performance.CategoryDefinition
		catErr := d.db.WithContext(ctx).
			Where("review_period_id = ? AND grade_group_id = ? AND objective_category_id = ?",
				reviewPeriodID, staffJobGradeGroupID, gapClosure.ObjectiveCategoryID).
			First(&catDef).Error
		if catErr == nil {
			totalGapClosureScore = catDef.MaxPoints * (gapClosure.FinalScore / 100)
			scoreCard.PercentageGapsClosureScore = totalGapClosureScore
		}
	}

	// SLA deduction
	var allRequests []performance.FeedbackRequestLog
	d.db.WithContext(ctx).
		Where("assigned_staff_id = ? AND time_initiated >= ? AND time_initiated <= ?",
			staffID, startDate, endDate).
		Find(&allRequests)

	var deductedPoints float64
	if len(allRequests) > 0 {
		now := time.Now()

		completedOverdueList := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
				isClosedOrCompletedOrBreached(r.RecordStatus) &&
				r.HasSLA && r.TimeCompleted != nil &&
				r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(requestSLA)
		})
		completedOverdue360List := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
				isClosedOrCompletedOrBreached(r.RecordStatus) &&
				r.HasSLA && r.TimeCompleted != nil &&
				r.TimeCompleted.Sub(r.TimeInitiated).Hours() > float64(pms360SLA)
		})
		completedOverdue := d.parent.countSLABreachedRequests(completedOverdueList, staffID, requestSLA) +
			d.parent.countSLABreachedRequests(completedOverdue360List, staffID, pms360SLA)

		pendingOverdueList := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
				r.RecordStatus == enums.StatusActive.String() &&
				r.HasSLA && r.TimeInitiated.Add(time.Duration(requestSLA)*time.Hour).Before(now)
		})
		pendingOverdue360List := filterRequests(allRequests, func(r performance.FeedbackRequestLog) bool {
			return r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
				r.RecordStatus == enums.StatusActive.String() &&
				r.HasSLA && r.TimeInitiated.Add(time.Duration(pms360SLA)*time.Hour).Before(now)
		})
		pendingOverdue := d.parent.countSLABreachedRequests(pendingOverdueList, staffID, requestSLA) +
			d.parent.countSLABreachedRequests(pendingOverdue360List, staffID, pms360SLA)

		deductedPoints = float64(completedOverdue + pendingOverdue)
		if deductedPoints > reviewPeriodMaxPoints {
			deductedPoints = reviewPeriodMaxPoints
		}
	}

	scoreCard.DeductedPoints = deductedPoints
	accumulatedPoints := workProductPoints + livingTheValuesPoints + totalGapClosureScore
	actualPoints := accumulatedPoints - deductedPoints
	if actualPoints < 0 {
		actualPoints = 0
	}

	scoreCard.StaffID = staffID
	scoreCard.AccumulatedPoints = accumulatedPoints
	scoreCard.ActualPoints = actualPoints

	// Performance grade
	if reviewPeriod.MaxPoints > 0 {
		performancePercentage := 100 * (actualPoints / reviewPeriod.MaxPoints)
		grade := getGrade(performancePercentage)
		scoreCard.PercentageScore = performancePercentage
		scoreCard.StaffPerformanceGrade = grade.String()
	}

	resp.ScoreCard = &scoreCard
	resp.Message = "operation completed successfully"

	return resp, nil
}

// =========================================================================
// GetStaffAnnualPerformanceScoreCardStatistics -- annual score card.
// Mirrors .NET GetStaffAnnualPerformanceScoreCardStatistics.
// Returns score cards across all review periods for a staff member in a year.
// =========================================================================

func (d *dashboardService) GetStaffAnnualPerformanceScoreCardStatistics(
	ctx context.Context,
	staffID string,
	year int,
) (performance.StaffAnnualScoreCardResponseVm, error) {
	resp := performance.StaffAnnualScoreCardResponseVm{
		StaffID: staffID,
		Year:    year,
	}
	resp.Message = "an error occurred"

	// Get all review periods for the given year
	var reviewPeriods []performance.PerformanceReviewPeriod
	d.db.WithContext(ctx).
		Where("year = ? AND record_status IN ?", year, []string{
			enums.StatusActive.String(),
			enums.StatusClosed.String(),
			enums.StatusCompleted.String(),
		}).
		Order("start_date ASC").
		Find(&reviewPeriods)

	if len(reviewPeriods) == 0 {
		resp.Message = "no review periods found for the given year"
		return resp, nil
	}

	var scoreCards []performance.StaffScoreCardDetails
	for _, rp := range reviewPeriods {
		scResp, err := d.GetStaffPerformanceScoreCardStatistics(ctx, staffID, rp.PeriodID)
		if err != nil {
			d.log.Warn().Err(err).
				Str("staffID", staffID).
				Str("periodID", rp.PeriodID).
				Msg("failed to compute score card for period, skipping")
			continue
		}
		if scResp.ScoreCard != nil {
			scoreCards = append(scoreCards, *scResp.ScoreCard)
		}
	}

	resp.ScoreCards = scoreCards
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetSubordinatesStaffPerformanceScoreCardStatistics -- subordinate scores.
// Mirrors .NET GetSubordinatesStaffPerformanceScoreCardStatistics.
// Returns score cards for all subordinates of a given manager.
// =========================================================================

func (d *dashboardService) GetSubordinatesStaffPerformanceScoreCardStatistics(
	ctx context.Context,
	managerStaffID, reviewPeriodID string,
) (performance.AllStaffScoreCardResponseVm, error) {
	resp := performance.AllStaffScoreCardResponseVm{}
	resp.Message = "an error occurred"

	// Get subordinates from ERP service
	var subordinateIDs []string
	if d.parent.erpEmployeeSvc != nil {
		subsResult, err := d.parent.erpEmployeeSvc.GetEmployeeSubordinates(ctx, managerStaffID)
		if err == nil && subsResult != nil {
			// Extract staff IDs from subordinates result.
			// The ERP service returns an interface{} -- we try to cast to a
			// slice of maps or structs with a StaffID/EmployeeNumber field.
			if subs, ok := subsResult.([]interface{}); ok {
				for _, sub := range subs {
					if empMap, ok := sub.(map[string]interface{}); ok {
						if empNum, ok := empMap["employeeNumber"].(string); ok {
							subordinateIDs = append(subordinateIDs, empNum)
						}
					}
				}
			}
		}
	}

	// If ERP lookup fails or returns nothing, fall back to planned objectives
	if len(subordinateIDs) == 0 {
		// Find staff who have planned objectives in this review period
		// and whose approver (line manager) is the given manager.
		var plannedObjs []performance.ReviewPeriodIndividualPlannedObjective
		d.db.WithContext(ctx).
			Where("review_period_id = ? AND approved_by = ? AND record_status IN ?",
				reviewPeriodID, managerStaffID, []string{
					enums.StatusActive.String(),
					enums.StatusCompleted.String(),
					enums.StatusClosed.String(),
				}).
			Find(&plannedObjs)

		seen := make(map[string]bool)
		for _, po := range plannedObjs {
			if !seen[po.StaffID] {
				seen[po.StaffID] = true
				subordinateIDs = append(subordinateIDs, po.StaffID)
			}
		}
	}

	var scoreCards []performance.StaffScoreCardDetails
	for _, subID := range subordinateIDs {
		scResp, err := d.GetStaffPerformanceScoreCardStatistics(ctx, subID, reviewPeriodID)
		if err != nil {
			d.log.Warn().Err(err).
				Str("subordinateID", subID).
				Msg("failed to compute subordinate score card, skipping")
			continue
		}
		if scResp.ScoreCard != nil {
			// Enrich with staff name from ERP if available
			sc := *scResp.ScoreCard
			if d.parent.erpEmployeeSvc != nil {
				empDetail, empErr := d.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, subID)
				if empErr == nil && empDetail != nil {
					if nameHolder, ok := empDetail.(interface{ GetFullName() string }); ok {
						sc.StaffName = nameHolder.GetFullName()
					}
				}
			}
			scoreCards = append(scoreCards, sc)
		}
	}

	resp.StaffScoreCards = scoreCards
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetOrganogramPerformanceSummaryStatistics -- single organogram unit summary.
// Mirrors .NET GetOrganogramPerformanceSummaryStatistics.
// Returns aggregated performance summary for a department/division/office.
// =========================================================================

func (d *dashboardService) GetOrganogramPerformanceSummaryStatistics(
	ctx context.Context,
	referenceID, reviewPeriodID string,
	organogramLevel enums.OrganogramLevel,
) (performance.OrganogramPerformanceSummaryResponseVm, error) {
	resp := performance.OrganogramPerformanceSummaryResponseVm{
		ReferenceID:     referenceID,
		ReviewPeriodID:  reviewPeriodID,
		OrganogramLevel: organogramLevel,
	}
	resp.Message = "an error occurred"

	// Get review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := d.db.WithContext(ctx).
		Where("period_id = ?", reviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, err
	}

	resp.ReviewPeriod = reviewPeriod.Name
	resp.ReviewPeriodShortName = reviewPeriod.ShortName
	if resp.ReviewPeriodShortName == "" {
		resp.ReviewPeriodShortName = reviewPeriod.Name
	}
	resp.MaxPoint = reviewPeriod.MaxPoints
	resp.Year = reviewPeriod.Year

	// Get all period scores for this review period
	var periodScores []performance.PeriodScore
	d.db.WithContext(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&periodScores)

	// Filter scores by the organogram reference.
	// This is a simplified implementation -- in production the organogram
	// hierarchy (office -> division -> department) would be resolved via
	// the ERP service to collect all staff IDs under the reference.
	// For now we use the OfficeID field on PeriodScore as a proxy.
	var filteredScores []performance.PeriodScore
	for _, ps := range periodScores {
		// Match based on organogram level
		switch organogramLevel {
		case enums.OrganogramLevelOffice:
			if fmt.Sprintf("%d", ps.OfficeID) == referenceID {
				filteredScores = append(filteredScores, ps)
			}
		default:
			// For department/division/bankwide, include all scores.
			// A full implementation would resolve the hierarchy.
			filteredScores = append(filteredScores, ps)
		}
	}

	totalStaff := len(filteredScores)
	resp.TotalStaff = totalStaff

	if totalStaff == 0 {
		resp.Message = "operation completed successfully"
		return resp, nil
	}

	// Aggregate scores
	var totalActual float64
	for _, ps := range filteredScores {
		totalActual += ps.FinalScore
	}
	resp.ActualScore = totalActual
	if totalStaff > 0 {
		avgPercentage := (totalActual / (float64(totalStaff) * reviewPeriod.MaxPoints)) * 100
		resp.PerformanceScore = avgPercentage
		grade := getGrade(avgPercentage)
		resp.EarnedPerformanceGrade = grade.String()
	}

	// Work product stats across all staff in scope
	excluded := excludedStatuses()
	staffIDs := make([]string, 0, totalStaff)
	for _, ps := range filteredScores {
		staffIDs = append(staffIDs, ps.StaffID)
	}

	var workProducts []performance.WorkProduct
	if len(staffIDs) > 0 {
		d.db.WithContext(ctx).
			Where("staff_id IN ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
				staffIDs, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
			Find(&workProducts)
	}

	completedStatuses := []string{
		enums.StatusCompleted.String(),
		enums.StatusAwaitingEvaluation.String(),
		enums.StatusClosed.String(),
	}
	resp.TotalWorkProducts = len(workProducts)
	resp.TotalWorkProductsCompletedOnSchedule = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return isInStatuses(wp.RecordStatus, completedStatuses) &&
			wp.CompletionDate != nil && !wp.CompletionDate.After(wp.EndDate)
	})
	resp.TotalWorkProductsBehindSchedule = countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return (wp.RecordStatus == enums.StatusActive.String() && wp.EndDate.Before(time.Now())) ||
			(isInStatuses(wp.RecordStatus, completedStatuses) && wp.CompletionDate != nil && wp.CompletionDate.After(wp.EndDate))
	})

	closedCount := countWhere(workProducts, func(wp performance.WorkProduct) bool {
		return isInStatuses(wp.RecordStatus, completedStatuses)
	})
	if len(workProducts) > 0 {
		resp.PercentageWorkProductsClosed = 100.0 * float64(closedCount) / float64(len(workProducts))
		resp.PercentageWorkProductsPending = 100.0 * float64(len(workProducts)-closedCount) / float64(len(workProducts))
	}

	// 360 feedback stats
	var feedbacks []performance.CompetencyReviewFeedback
	if len(staffIDs) > 0 {
		d.db.WithContext(ctx).
			Where("staff_id IN ? AND review_period_id = ? AND record_status NOT IN ?",
				staffIDs, reviewPeriodID, excluded).
			Preload("CompetencyReviewers").
			Find(&feedbacks)
	}

	resp.Total360Feedbacks = len(feedbacks)
	for _, fb := range feedbacks {
		for _, cr := range fb.CompetencyReviewers {
			if cr.RecordStatus == enums.StatusCompleted.String() || cr.RecordStatus == enums.StatusClosed.String() {
				resp.Completed360FeedbacksToTreat++
			} else if cr.RecordStatus == enums.StatusActive.String() {
				resp.Pending360FeedbacksToTreat++
			}
		}
	}

	// Competency gap closure stats
	var gapClosures []performance.CompetencyGapClosure
	if len(staffIDs) > 0 {
		d.db.WithContext(ctx).
			Where("staff_id IN ? AND review_period_id = ?", staffIDs, reviewPeriodID).
			Find(&gapClosures)
	}

	resp.TotalCompetencyGaps = len(gapClosures)
	if len(gapClosures) > 0 {
		closedGaps := countWhere(gapClosures, func(gc performance.CompetencyGapClosure) bool {
			return gc.FinalScore > 0
		})
		resp.PercentageGapsClosure = 100.0 * float64(closedGaps) / float64(len(gapClosures))
	}

	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetOrganogramPerformanceSummaryListStatistics -- org performance summary list.
// Mirrors .NET GetOrganogramPerformanceSummaryListStatistics.
// Returns a list of performance summaries for all units at a given organogram level.
// =========================================================================

func (d *dashboardService) GetOrganogramPerformanceSummaryListStatistics(
	ctx context.Context,
	headOfUnitID, reviewPeriodID string,
	organogramLevel enums.OrganogramLevel,
) (performance.OrganogramPerformanceSummaryListResponseVm, error) {
	resp := performance.OrganogramPerformanceSummaryListResponseVm{
		HeadOfUnitID:    headOfUnitID,
		ReviewPeriodID:  reviewPeriodID,
		OrganogramLevel: organogramLevel,
	}
	resp.Message = "an error occurred"

	// Get review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := d.db.WithContext(ctx).
		Where("period_id = ?", reviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, err
	}

	resp.ReviewPeriod = reviewPeriod.Name
	resp.ReviewPeriodShortName = reviewPeriod.ShortName
	if resp.ReviewPeriodShortName == "" {
		resp.ReviewPeriodShortName = reviewPeriod.Name
	}
	resp.MaxPoint = reviewPeriod.MaxPoints
	resp.Year = reviewPeriod.Year

	// Get all period scores for this review period
	var periodScores []performance.PeriodScore
	d.db.WithContext(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&periodScores)

	if len(periodScores) == 0 {
		resp.OrganogramPerformances = []performance.OrganogramPerformanceSummaryDetails{}
		resp.TotalRecords = 0
		resp.Message = "operation completed successfully"
		return resp, nil
	}

	// Group period scores by organogram reference.
	// This implementation uses OfficeID as a basic grouping key.
	// A full implementation would resolve the organogram hierarchy from ERP
	// to build department/division/office groupings.
	type groupData struct {
		referenceID   string
		referenceName string
		managerID     string
		managerName   string
		scores        []performance.PeriodScore
	}

	groupMap := make(map[string]*groupData)

	for _, ps := range periodScores {
		key := fmt.Sprintf("%d", ps.OfficeID)
		if _, exists := groupMap[key]; !exists {
			groupMap[key] = &groupData{
				referenceID: key,
			}
		}
		groupMap[key].scores = append(groupMap[key].scores, ps)
	}

	excluded := excludedStatuses()
	completedStatuses := []string{
		enums.StatusCompleted.String(),
		enums.StatusAwaitingEvaluation.String(),
		enums.StatusClosed.String(),
	}

	var summaries []performance.OrganogramPerformanceSummaryDetails
	for _, grp := range groupMap {
		detail := performance.OrganogramPerformanceSummaryDetails{
			ReferenceID: grp.referenceID,
			TotalStaff:  len(grp.scores),
		}

		// Aggregate scores
		var totalActual float64
		staffIDs := make([]string, 0, len(grp.scores))
		for _, ps := range grp.scores {
			totalActual += ps.FinalScore
			staffIDs = append(staffIDs, ps.StaffID)
		}
		detail.ActualScore = totalActual
		if detail.TotalStaff > 0 && reviewPeriod.MaxPoints > 0 {
			avgPct := (totalActual / (float64(detail.TotalStaff) * reviewPeriod.MaxPoints)) * 100
			detail.PerformanceScore = avgPct
			detail.EarnedPerformanceGrade = getGrade(avgPct).String()
		}

		// Work product stats for this group
		var workProducts []performance.WorkProduct
		if len(staffIDs) > 0 {
			d.db.WithContext(ctx).
				Where("staff_id IN ? AND record_status NOT IN ? AND start_date >= ? AND end_date <= ?",
					staffIDs, excluded, reviewPeriod.StartDate, reviewPeriod.EndDate).
				Find(&workProducts)
		}

		detail.TotalWorkProducts = len(workProducts)
		detail.TotalWorkProductsCompletedOnSchedule = countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return isInStatuses(wp.RecordStatus, completedStatuses) &&
				wp.CompletionDate != nil && !wp.CompletionDate.After(wp.EndDate)
		})
		detail.TotalWorkProductsBehindSchedule = countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return (wp.RecordStatus == enums.StatusActive.String() && wp.EndDate.Before(time.Now())) ||
				(isInStatuses(wp.RecordStatus, completedStatuses) && wp.CompletionDate != nil && wp.CompletionDate.After(wp.EndDate))
		})

		closedCount := countWhere(workProducts, func(wp performance.WorkProduct) bool {
			return isInStatuses(wp.RecordStatus, completedStatuses)
		})
		if len(workProducts) > 0 {
			detail.PercentageWorkProductsClosed = 100.0 * float64(closedCount) / float64(len(workProducts))
			detail.PercentageWorkProductsPending = 100.0 * float64(len(workProducts)-closedCount) / float64(len(workProducts))
		}

		// 360 feedback stats
		var feedbacks []performance.CompetencyReviewFeedback
		if len(staffIDs) > 0 {
			d.db.WithContext(ctx).
				Where("staff_id IN ? AND review_period_id = ? AND record_status NOT IN ?",
					staffIDs, reviewPeriodID, excluded).
				Preload("CompetencyReviewers").
				Find(&feedbacks)
		}

		detail.Total360Feedbacks = len(feedbacks)
		for _, fb := range feedbacks {
			for _, cr := range fb.CompetencyReviewers {
				if cr.RecordStatus == enums.StatusCompleted.String() || cr.RecordStatus == enums.StatusClosed.String() {
					detail.Completed360FeedbacksToTreat++
				} else if cr.RecordStatus == enums.StatusActive.String() {
					detail.Pending360FeedbacksToTreat++
				}
			}
		}

		// Gap closure stats
		var gapClosures []performance.CompetencyGapClosure
		if len(staffIDs) > 0 {
			d.db.WithContext(ctx).
				Where("staff_id IN ? AND review_period_id = ?", staffIDs, reviewPeriodID).
				Find(&gapClosures)
		}

		detail.TotalCompetencyGaps = len(gapClosures)
		if len(gapClosures) > 0 {
			closedGaps := countWhere(gapClosures, func(gc performance.CompetencyGapClosure) bool {
				return gc.FinalScore > 0
			})
			detail.PercentageGapsClosure = 100.0 * float64(closedGaps) / float64(len(gapClosures))
		}

		summaries = append(summaries, detail)
	}

	resp.OrganogramPerformances = summaries
	resp.TotalRecords = len(summaries)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// Helper functions for filtering, counting, and aggregating.
// =========================================================================

// countWhere counts elements in a slice matching a predicate.
func countWhere[T any](items []T, pred func(T) bool) int {
	count := 0
	for _, item := range items {
		if pred(item) {
			count++
		}
	}
	return count
}

// filterRequests filters feedback request logs by a predicate.
func filterRequests(requests []performance.FeedbackRequestLog, pred func(performance.FeedbackRequestLog) bool) []performance.FeedbackRequestLog {
	var result []performance.FeedbackRequestLog
	for _, r := range requests {
		if pred(r) {
			result = append(result, r)
		}
	}
	return result
}

// sumFloat64 sums a float64 field across a slice.
func sumFloat64[T any](items []T, fn func(T) float64) float64 {
	var total float64
	for _, item := range items {
		total += fn(item)
	}
	return total
}

// average computes the arithmetic mean of a float64 slice.
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// isClosedOrCompletedOrBreached checks if a status string is one of the
// "finished" statuses used in SLA breach calculations.
func isClosedOrCompletedOrBreached(status string) bool {
	return status == enums.StatusClosed.String() ||
		status == enums.StatusCompleted.String() ||
		status == enums.StatusBreached.String()
}

// isInStatuses checks if a status string is in a list of statuses.
func isInStatuses(status string, statuses []string) bool {
	for _, s := range statuses {
		if status == s {
			return true
		}
	}
	return false
}
