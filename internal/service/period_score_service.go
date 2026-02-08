package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// periodScoreService handles period-score retrieval and reporting.
//
// Mirrors .NET methods in the #region reporting block:
//   - GetPerformanceScore       – composite: returns the score-card via dashboard
//   - GetPeriodScoreDetails     – single staff + review period score
//   - GetPeriodScores           – all scores for a review period
//   - GetStaffReviewPeriods     – review periods that a staff has participated in
// ---------------------------------------------------------------------------

type periodScoreService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	// Back-reference to parent for shared helpers and peer services.
	parent *performanceManagementService

	// Repositories
	periodScoreRepo  *repository.PMSRepository[performance.PeriodScore]
	reviewPeriodRepo *repository.PMSRepository[performance.PerformanceReviewPeriod]
	plannedObjRepo   *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	competencyFBRepo *repository.PMSRepository[performance.CompetencyReviewFeedback]
}

func newPeriodScoreService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *periodScoreService {
	return &periodScoreService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "period_score").Logger(),
		parent: parent,

		periodScoreRepo:  repository.NewPMSRepository[performance.PeriodScore](db),
		reviewPeriodRepo: repository.NewPMSRepository[performance.PerformanceReviewPeriod](db),
		plannedObjRepo:   repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
		competencyFBRepo: repository.NewPMSRepository[performance.CompetencyReviewFeedback](db),
	}
}

// =========================================================================
// GetPerformanceScore – top-level method delegated from the main service.
// Returns the score-card for the staff member's active review period.
// =========================================================================

func (ps *periodScoreService) GetPerformanceScore(ctx context.Context, staffID string) (interface{}, error) {
	// Use the dashboard sub-service to compute the full score card.
	return ps.parent.dashboard.GetDashboardStats(ctx, staffID)
}

// =========================================================================
// GetPeriodScoreDetails – retrieves a single PeriodScore for a staff member
// in a specific review period, enriched with organogram and strategy data.
// Mirrors .NET GetPeriodScoreDetails.
// =========================================================================

func (ps *periodScoreService) GetPeriodScoreDetails(ctx context.Context, reviewPeriodID, staffID string) (performance.PeriodScoreResponseVm, error) {
	resp := performance.PeriodScoreResponseVm{}
	resp.Message = "an error occurred"

	var score performance.PeriodScore
	err := ps.db.WithContext(ctx).
		Preload("ReviewPeriod").
		Preload("Strategy").
		Where("review_period_id = ? AND staff_id = ?", reviewPeriodID, staffID).
		First(&score).Error

	if err != nil {
		ps.log.Error().Err(err).
			Str("reviewPeriodID", reviewPeriodID).
			Str("staffID", staffID).
			Msg("period score not found")
		resp.HasError = true
		resp.Message = "period score record not found"
		return resp, fmt.Errorf("period score not found: %w", err)
	}

	data := ps.mapPeriodScoreToData(score)

	// Enrich with staff name from ERP service (if available)
	if ps.parent.erpEmployeeSvc != nil {
		empDetail, empErr := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
		if empErr == nil && empDetail != nil {
			if nameHolder, ok := empDetail.(interface{ GetFullName() string }); ok {
				data.StaffFullName = nameHolder.GetFullName()
			}
		}
	}

	resp.PeriodScore = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetPeriodScores – retrieves all PeriodScores for a review period,
// ordered by staff name.
// Mirrors .NET GetPeriodScores.
// =========================================================================

func (ps *periodScoreService) GetPeriodScores(ctx context.Context, reviewPeriodID string) (performance.PeriodScoreListResponseVm, error) {
	resp := performance.PeriodScoreListResponseVm{}
	resp.Message = "an error occurred"

	var scores []performance.PeriodScore
	err := ps.db.WithContext(ctx).
		Preload("ReviewPeriod").
		Preload("Strategy").
		Where("review_period_id = ?", reviewPeriodID).
		Find(&scores).Error

	if err != nil {
		ps.log.Error().Err(err).Str("reviewPeriodID", reviewPeriodID).Msg("failed to get period scores")
		resp.HasError = true
		return resp, err
	}

	var dataList []performance.PeriodScoreData
	for _, score := range scores {
		data := ps.mapPeriodScoreToData(score)

		// Enrich with staff name
		if ps.parent.erpEmployeeSvc != nil {
			empDetail, empErr := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, score.StaffID)
			if empErr == nil && empDetail != nil {
				if nameHolder, ok := empDetail.(interface{ GetFullName() string }); ok {
					data.StaffFullName = nameHolder.GetFullName()
				}
			}
		}

		dataList = append(dataList, data)
	}

	// Sort by staff full name ascending
	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].StaffFullName < dataList[j].StaffFullName
	})

	resp.PeriodScores = dataList
	resp.TotalRecords = len(dataList)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetStaffReviewPeriods – returns the review periods in which a staff
// member has participated (via planned objectives or competency feedbacks).
// Mirrors .NET GetStaffReviewPeriods.
// =========================================================================

func (ps *periodScoreService) GetStaffReviewPeriods(ctx context.Context, staffID string) (performance.GetStaffReviewPeriodResponseVm, error) {
	resp := performance.GetStaffReviewPeriodResponseVm{}
	resp.Message = "an error occurred"
	resp.StaffID = staffID

	// Validate staff exists via ERP service
	if ps.parent.erpEmployeeSvc != nil {
		empDetail, empErr := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
		if empErr != nil || empDetail == nil {
			resp.HasError = true
			resp.Message = "staff record not found"
			return resp, fmt.Errorf("staff not found: %s", staffID)
		}
	}

	participationStatuses := []string{
		enums.StatusActive.String(),
		enums.StatusCompleted.String(),
		enums.StatusClosed.String(),
		enums.StatusReturned.String(),
		enums.StatusPendingApproval.String(),
		enums.StatusPaused.String(),
	}

	// Try planned objectives first
	var plannedObjs []performance.ReviewPeriodIndividualPlannedObjective
	ps.db.WithContext(ctx).
		Where("staff_id = ? AND record_status IN ?", staffID, participationStatuses).
		Preload("ReviewPeriod").
		Find(&plannedObjs)

	// Collect unique review periods
	seenPeriods := make(map[string]bool)
	var reviewPeriods []performance.PerformanceReviewPeriod

	if len(plannedObjs) > 0 {
		for _, po := range plannedObjs {
			if po.ReviewPeriod != nil && !seenPeriods[po.ReviewPeriod.PeriodID] {
				seenPeriods[po.ReviewPeriod.PeriodID] = true
				reviewPeriods = append(reviewPeriods, *po.ReviewPeriod)
			}
		}
	} else {
		// Fallback: use competency review feedbacks
		var feedbacks []performance.CompetencyReviewFeedback
		ps.db.WithContext(ctx).
			Where("staff_id = ? AND record_status IN ?", staffID, participationStatuses).
			Preload("ReviewPeriod").
			Find(&feedbacks)

		for _, fb := range feedbacks {
			if fb.ReviewPeriod != nil && !seenPeriods[fb.ReviewPeriod.PeriodID] {
				seenPeriods[fb.ReviewPeriod.PeriodID] = true
				reviewPeriods = append(reviewPeriods, *fb.ReviewPeriod)
			}
		}
	}

	// Sort by year DESC, then start date ASC
	sort.Slice(reviewPeriods, func(i, j int) bool {
		if reviewPeriods[i].Year != reviewPeriods[j].Year {
			return reviewPeriods[i].Year > reviewPeriods[j].Year
		}
		return reviewPeriods[i].StartDate.Before(reviewPeriods[j].StartDate)
	})

	// Map to VMs
	var vms []performance.PerformanceReviewPeriodVm
	for _, rp := range reviewPeriods {
		vm := performance.PerformanceReviewPeriodVm{
			PeriodID:                   rp.PeriodID,
			Year:                       rp.Year,
			Range:                      int(rp.Range),
			RangeValue:                 rp.RangeValue,
			Name:                       rp.Name,
			ShortName:                  rp.ShortName,
			Description:                rp.Description,
			StartDate:                  rp.StartDate,
			EndDate:                    rp.EndDate,
			MaxPoints:                  rp.MaxPoints,
			MinNoOfObjectives:          rp.MinNoOfObjectives,
			MaxNoOfObjectives:          rp.MaxNoOfObjectives,
			StrategyID:                 rp.StrategyID,
			AllowObjectivePlanning:     rp.AllowObjectivePlanning,
			AllowWorkProductPlanning:   rp.AllowWorkProductPlanning,
			AllowWorkProductEvaluation: rp.AllowWorkProductEvaluation,
		}
		vm.RecordStatus = rp.RecordStatus
		vm.IsActive = rp.IsActive
		vm.CreatedAt = rp.CreatedAt
		vm.CreatedBy = rp.CreatedBy
		vm.UpdatedAt = rp.UpdatedAt
		vm.UpdatedBy = rp.UpdatedBy
		vm.IsApproved = rp.IsApproved
		vm.IsRejected = rp.IsRejected
		vm.ApprovedBy = rp.ApprovedBy
		vm.DateApproved = rp.DateApproved
		vm.RejectedBy = rp.RejectedBy
		vm.RejectionReason = rp.RejectionReason
		vm.DateRejected = rp.DateRejected

		vms = append(vms, vm)
	}

	resp.PerformanceReviewPeriods = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// Internal helpers
// =========================================================================

// mapPeriodScoreToData converts a PeriodScore domain entity to the
// PeriodScoreData response DTO, enriching with review period and strategy.
func (ps *periodScoreService) mapPeriodScoreToData(score performance.PeriodScore) performance.PeriodScoreData {
	data := performance.PeriodScoreData{
		BaseEntityVm: performance.BaseEntityVm{
			ID:           score.ID,
			RecordStatus: score.RecordStatus,
			CreatedAt:    score.CreatedAt,
			UpdatedAt:    score.UpdatedAt,
			CreatedBy:    score.CreatedBy,
			UpdatedBy:    score.UpdatedBy,
			IsActive:     score.IsActive,
		},
		PeriodScoreID:     score.PeriodScoreID,
		ReviewPeriodID:    score.ReviewPeriodID,
		StaffID:           score.StaffID,
		FinalScore:        score.FinalScore,
		ScorePercentage:   score.ScorePercentage,
		FinalGrade:        int(score.FinalGrade),
		FinalGradeName:    score.FinalGrade.String(),
		EndDate:           score.EndDate,
		OfficeID:          score.OfficeID,
		MinNoOfObjectives: score.MinNoOfObjectives,
		MaxNoOfObjectives: score.MaxNoOfObjectives,
		StrategyID:        score.StrategyID,
		StaffGrade:        score.StaffGrade,
		LocationID:        score.LocationID,
		HRDDeductedPoints: score.HRDDeductedPoints,
		IsUnderPerforming: score.IsUnderPerforming,
	}

	// Enrich from ReviewPeriod navigation
	if score.ReviewPeriod != nil {
		data.ReviewPeriod = score.ReviewPeriod.Name
		data.Year = score.ReviewPeriod.Year
		data.MaxPoint = score.ReviewPeriod.MaxPoints
		data.StartDate = score.ReviewPeriod.StartDate
		data.EndDate = score.ReviewPeriod.EndDate
	}

	// Enrich from Strategy navigation
	if score.Strategy != nil {
		data.StrategyName = score.Strategy.Name
	}

	// NOTE: StaffOffice/Division/Department organogram enrichment requires
	// the Office model to be loaded. This is a placeholder until the
	// organogram relationships are fully wired on PeriodScore.
	// In the .NET source: obj.StaffOffice.OfficeCode, Division.*, Department.*

	return data
}
