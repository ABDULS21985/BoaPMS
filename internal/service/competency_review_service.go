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
// competencyReviewService handles 360-degree review and competency feedback.
//
// Mirrors .NET methods:
//   - ReviewPeriod360ReviewSetup / Initiate360Review / Complete360Review
//   - CompetencyReviewFeedbackSetup / CompetencyReviewerSetup / CompetencyRatingSetup
//   - GetCompetencyReviewFeedback / GetCompetencyReviewFeedbackDetails
//   - GetAllCompetencyReviewFeedbacks / GetCompetencyReviews
//   - GetReviewerFeedbackDetails / GetQuestionnaire
//   - CompetencyGapClosureSetup
// ---------------------------------------------------------------------------

type competencyReviewService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	parent *performanceManagementService

	competencyFeedbackRepo *repository.PMSRepository[performance.CompetencyReviewFeedback]
	competencyReviewerRepo *repository.PMSRepository[performance.CompetencyReviewer]
	competencyRatingRepo   *repository.PMSRepository[performance.CompetencyReviewerRating]
	competencyGapRepo      *repository.PMSRepository[performance.CompetencyGapClosure]
	pmsCompetencyRepo      *repository.PMSRepository[performance.PmsCompetency]
	questionaireRepo       *repository.PMSRepository[performance.FeedbackQuestionaire]
	questionaireOptRepo    *repository.PMSRepository[performance.FeedbackQuestionaireOption]
	feedbackLogRepo        *repository.PMSRepository[performance.FeedbackRequestLog]
}

func newCompetencyReviewService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *competencyReviewService {
	return &competencyReviewService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "competency_review").Logger(),
		parent: parent,

		competencyFeedbackRepo: repository.NewPMSRepository[performance.CompetencyReviewFeedback](db),
		competencyReviewerRepo: repository.NewPMSRepository[performance.CompetencyReviewer](db),
		competencyRatingRepo:   repository.NewPMSRepository[performance.CompetencyReviewerRating](db),
		competencyGapRepo:      repository.NewPMSRepository[performance.CompetencyGapClosure](db),
		pmsCompetencyRepo:      repository.NewPMSRepository[performance.PmsCompetency](db),
		questionaireRepo:       repository.NewPMSRepository[performance.FeedbackQuestionaire](db),
		questionaireOptRepo:    repository.NewPMSRepository[performance.FeedbackQuestionaireOption](db),
		feedbackLogRepo:        repository.NewPMSRepository[performance.FeedbackRequestLog](db),
	}
}

// =========================================================================
// CompetencyReviewFeedbackSetup -- manages competency review feedback records.
// Mirrors .NET CompetencyReviewFeedbackSetup.
// =========================================================================

func (cr *competencyReviewService) CompetencyReviewFeedbackSetup(ctx context.Context, req *performance.CompetencyReviewFeedbackRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.Status {
	case enums.OperationAdd.String():
		// Check if feedback already exists for this staff in this review period
		var existing performance.CompetencyReviewFeedback
		if err := cr.db.WithContext(ctx).
			Where("staff_id = ? AND review_period_id = ? AND record_status != ?",
				req.StaffID, req.ReviewPeriodID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("competency review feedback already exists for this staff in this review period")
		}

		feedback := performance.CompetencyReviewFeedback{
			StaffID:        req.StaffID,
			MaxPoints:      0, // Will be calculated based on competencies
			FinalScore:     req.FinalScore,
			ReviewPeriodID: req.ReviewPeriodID,
		}
		feedback.RecordStatus = enums.StatusActive.String()
		feedback.IsActive = true

		if err := cr.db.WithContext(ctx).Create(&feedback).Error; err != nil {
			return resp, fmt.Errorf("creating competency review feedback: %w", err)
		}

		resp.ID = feedback.CompetencyReviewFeedbackID
		resp.Message = "competency review feedback created successfully"

	case enums.OperationUpdate.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyReviewFeedback{}).
			Where("competency_review_feedback_id = ?", req.CompetencyReviewFeedbackID).
			Updates(map[string]interface{}{
				"final_score": req.FinalScore,
			})
		resp.ID = req.CompetencyReviewFeedbackID
		resp.Message = "competency review feedback updated"

	case enums.OperationClose.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyReviewFeedback{}).
			Where("competency_review_feedback_id = ?", req.CompetencyReviewFeedbackID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusClosed.String(),
			})
		resp.ID = req.CompetencyReviewFeedbackID
		resp.Message = "competency review feedback closed"

	case enums.OperationCancel.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyReviewFeedback{}).
			Where("competency_review_feedback_id = ?", req.CompetencyReviewFeedbackID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CompetencyReviewFeedbackID
		resp.Message = "competency review feedback cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for competency review feedback")
	}

	return resp, nil
}

// =========================================================================
// CompetencyReviewerSetup -- manages competency reviewers.
// Mirrors .NET CompetencyReviewerSetup.
// =========================================================================

func (cr *competencyReviewService) CompetencyReviewerSetup(ctx context.Context, req *performance.CompetencyReviewerRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.RecordStatus {
	case enums.OperationAdd.String():
		// Check duplicate
		var existing performance.CompetencyReviewer
		if err := cr.db.WithContext(ctx).
			Where("review_staff_id = ? AND competency_review_feedback_id = ? AND record_status != ?",
				req.ReviewStaffID, req.CompetencyReviewFeedbackID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("reviewer already assigned to this feedback")
		}

		reviewer := performance.CompetencyReviewer{
			ReviewStaffID:              req.ReviewStaffID,
			CompetencyReviewFeedbackID: req.CompetencyReviewFeedbackID,
		}
		reviewer.RecordStatus = enums.StatusActive.String()
		reviewer.IsActive = true

		if err := cr.db.WithContext(ctx).Create(&reviewer).Error; err != nil {
			return resp, fmt.Errorf("creating competency reviewer: %w", err)
		}

		resp.ID = reviewer.CompetencyReviewerID
		resp.Message = "competency reviewer added successfully"

	case enums.OperationCancel.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyReviewer{}).
			Where("competency_reviewer_id = ?", req.CompetencyReviewerID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CompetencyReviewerID
		resp.Message = "competency reviewer removed"

	case enums.OperationComplete.String():
		// Complete the reviewer's feedback
		cr.completeReviewerFeedback(ctx, req.CompetencyReviewerID)
		resp.ID = req.CompetencyReviewerID
		resp.Message = "competency reviewer feedback completed"

	default:
		return resp, fmt.Errorf("unsupported operation for competency reviewer")
	}

	return resp, nil
}

// =========================================================================
// CompetencyRatingSetup -- manages individual competency ratings.
// Mirrors .NET CompetencyRatingSetup / SavePmsCompetencyRating.
// =========================================================================

func (cr *competencyReviewService) CompetencyRatingSetup(ctx context.Context, req *performance.SavePmsCompetencyRequestVm) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	// Look up the option score
	var option performance.FeedbackQuestionaireOption
	if err := cr.db.WithContext(ctx).
		Where("feedback_questionaire_option_id = ?", req.FeedbackQuestionaireOptionID).
		First(&option).Error; err != nil {
		return resp, fmt.Errorf("feedback questionnaire option not found: %w", err)
	}

	// Check if rating already exists
	var existing performance.CompetencyReviewerRating
	err := cr.db.WithContext(ctx).
		Where("competency_reviewer_id = ? AND pms_competency_id = ?",
			req.CompetencyReviewerID, req.PmsCompetencyID).
		First(&existing).Error

	if err == nil {
		// Update existing
		existing.FeedbackQuestionaireOptionID = req.FeedbackQuestionaireOptionID
		existing.Rating = option.Score
		if err := cr.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return resp, fmt.Errorf("updating competency rating: %w", err)
		}
		resp.ID = existing.CompetencyReviewerRatingID
		resp.Message = "competency rating updated"
	} else {
		// Create new
		rating := performance.CompetencyReviewerRating{
			PmsCompetencyID:              req.PmsCompetencyID,
			FeedbackQuestionaireOptionID: req.FeedbackQuestionaireOptionID,
			Rating:                       option.Score,
			CompetencyReviewerID:         req.CompetencyReviewerID,
		}
		rating.RecordStatus = enums.StatusActive.String()
		rating.IsActive = true

		if err := cr.db.WithContext(ctx).Create(&rating).Error; err != nil {
			return resp, fmt.Errorf("creating competency rating: %w", err)
		}
		resp.ID = rating.CompetencyReviewerRatingID
		resp.Message = "competency rating created"
	}

	return resp, nil
}

// =========================================================================
// GetCompetencyReviewFeedback -- retrieves a single competency review feedback.
// Mirrors .NET GetCompetencyReviewFeedback.
// =========================================================================

func (cr *competencyReviewService) GetCompetencyReviewFeedback(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackResponseVm, error) {
	resp := performance.CompetencyReviewFeedbackResponseVm{}
	resp.Message = "an error occurred"

	var feedback performance.CompetencyReviewFeedback
	err := cr.db.WithContext(ctx).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Where("competency_review_feedback_id = ?", feedbackID).
		First(&feedback).Error
	if err != nil {
		cr.log.Error().Err(err).Str("feedbackID", feedbackID).Msg("competency review feedback not found")
		resp.HasError = true
		resp.Message = "competency review feedback not found"
		return resp, fmt.Errorf("competency review feedback not found: %w", err)
	}

	data := cr.mapFeedbackToData(ctx, feedback)
	resp.CompetencyReviewFeedback = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCompetencyReviewFeedbackDetails -- retrieves detailed competency review feedback.
// Mirrors .NET GetCompetencyReviewFeedbackDetails.
func (cr *competencyReviewService) GetCompetencyReviewFeedbackDetails(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackDetailsResponseVm, error) {
	resp := performance.CompetencyReviewFeedbackDetailsResponseVm{}
	resp.Message = "an error occurred"

	var feedback performance.CompetencyReviewFeedback
	err := cr.db.WithContext(ctx).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Preload("CompetencyReviewers.CompetencyReviewerRatings.PmsCompetency").
		Where("competency_review_feedback_id = ?", feedbackID).
		First(&feedback).Error
	if err != nil {
		resp.HasError = true
		resp.Message = "competency review feedback not found"
		return resp, fmt.Errorf("not found: %w", err)
	}

	// Build summary ratings by competency
	competencyRatings := make(map[string][]float64)
	competencyNames := make(map[string]string)

	for _, reviewer := range feedback.CompetencyReviewers {
		for _, rating := range reviewer.CompetencyReviewerRatings {
			competencyRatings[rating.PmsCompetencyID] = append(competencyRatings[rating.PmsCompetencyID], rating.Rating)
			if rating.PmsCompetency != nil {
				competencyNames[rating.PmsCompetencyID] = rating.PmsCompetency.Name
			}
		}
	}

	var ratings []performance.CompetencyReviewerRatingSummaryData
	for compID, scores := range competencyRatings {
		avg := 0.0
		if len(scores) > 0 {
			sum := 0.0
			for _, s := range scores {
				sum += s
			}
			avg = sum / float64(len(scores))
		}
		ratings = append(ratings, performance.CompetencyReviewerRatingSummaryData{
			PmsCompetencyID: compID,
			PmsCompetency:   competencyNames[compID],
			AverageRating:   avg,
		})
	}

	details := performance.CompetencyReviewFeedbackDetails{
		CompetencyReviewFeedbackID: feedback.CompetencyReviewFeedbackID,
		StaffID:                    feedback.StaffID,
		MaxPoints:                  feedback.MaxPoints,
		FinalScore:                 feedback.FinalScore,
		ReviewPeriodID:             feedback.ReviewPeriodID,
		RecordStatusName:           feedback.RecordStatus,
		Ratings:                    ratings,
	}

	if feedback.MaxPoints > 0 {
		details.FinalScorePercentage = (feedback.FinalScore / feedback.MaxPoints) * 100
	}

	// Enrich staff name
	if cr.parent.erpEmployeeSvc != nil {
		if detail, empErr := cr.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, feedback.StaffID); empErr == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				details.StaffName = nameHolder.GetFullName()
			}
		}
	}

	resp.CompetencyReviewFeedback = &details
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetAllCompetencyReviewFeedbacks -- retrieves all competency review feedbacks for a staff member.
// Mirrors .NET GetAllCompetencyReviewFeedbacks.
func (cr *competencyReviewService) GetAllCompetencyReviewFeedbacks(ctx context.Context, staffID string) (performance.CompetencyReviewFeedbackListResponseVm, error) {
	resp := performance.CompetencyReviewFeedbackListResponseVm{}
	resp.Message = "an error occurred"

	var feedbacks []performance.CompetencyReviewFeedback
	err := cr.db.WithContext(ctx).
		Where("staff_id = ?", staffID).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Find(&feedbacks).Error
	if err != nil {
		cr.log.Error().Err(err).Str("staffID", staffID).Msg("failed to get competency review feedbacks")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CompetencyReviewFeedbackData
	for _, fb := range feedbacks {
		data = append(data, cr.mapFeedbackToData(ctx, fb))
	}

	resp.CompetencyReviewFeedbacks = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCompetencyReviews -- retrieves competency reviewers for a staff member.
// Mirrors .NET GetCompetencyReviews.
func (cr *competencyReviewService) GetCompetencyReviews(ctx context.Context, reviewerStaffID string) (performance.CompetencyReviewersListResponseVm, error) {
	resp := performance.CompetencyReviewersListResponseVm{}
	resp.Message = "an error occurred"

	var reviewers []performance.CompetencyReviewer
	err := cr.db.WithContext(ctx).
		Where("review_staff_id = ? AND record_status = ?", reviewerStaffID, enums.StatusActive.String()).
		Preload("CompetencyReviewFeedback").
		Preload("CompetencyReviewerRatings").
		Preload("CompetencyReviewerRatings.PmsCompetency").
		Preload("CompetencyReviewerRatings.FeedbackQuestionaireOption").
		Find(&reviewers).Error
	if err != nil {
		cr.log.Error().Err(err).Str("reviewerStaffID", reviewerStaffID).Msg("failed to get competency reviews")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CompetencyReviewerData
	for _, r := range reviewers {
		d := cr.mapReviewerToData(ctx, r)
		data = append(data, d)
	}

	resp.CompetencyReviewers = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetReviewerFeedbackDetails -- retrieves detailed feedback from a specific reviewer.
// Mirrors .NET GetReviewerFeedbackDetails.
func (cr *competencyReviewService) GetReviewerFeedbackDetails(ctx context.Context, reviewerID string) (performance.CompetencyReviewersResponseVm, error) {
	resp := performance.CompetencyReviewersResponseVm{}
	resp.Message = "an error occurred"

	var reviewer performance.CompetencyReviewer
	err := cr.db.WithContext(ctx).
		Preload("CompetencyReviewFeedback").
		Preload("CompetencyReviewerRatings").
		Preload("CompetencyReviewerRatings.PmsCompetency").
		Preload("CompetencyReviewerRatings.FeedbackQuestionaireOption").
		Where("competency_reviewer_id = ?", reviewerID).
		First(&reviewer).Error
	if err != nil {
		cr.log.Error().Err(err).Str("reviewerID", reviewerID).Msg("competency reviewer not found")
		resp.HasError = true
		resp.Message = "reviewer not found"
		return resp, fmt.Errorf("reviewer not found: %w", err)
	}

	data := cr.mapReviewerToData(ctx, reviewer)
	resp.CompetencyReview = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetQuestionnaire -- retrieves the PMS competency questionnaire with options.
// Mirrors .NET GetQuestionnaire.
func (cr *competencyReviewService) GetQuestionnaire(ctx context.Context, staffID string) (performance.QuestionnaireListResponseVm, error) {
	resp := performance.QuestionnaireListResponseVm{}
	resp.Message = "an error occurred"

	var competencies []performance.PmsCompetency
	err := cr.db.WithContext(ctx).
		Where("record_status = ?", enums.StatusActive.String()).
		Preload("FeedbackQuestionaires").
		Preload("FeedbackQuestionaires.Options").
		Find(&competencies).Error
	if err != nil {
		cr.log.Error().Err(err).Msg("failed to get questionnaire")
		resp.HasError = true
		return resp, err
	}

	var data []performance.PmsCompetencyData
	for _, comp := range competencies {
		cd := performance.PmsCompetencyData{
			PmsCompetencyID:  comp.PmsCompetencyID,
			Name:             comp.Name,
			Description:      comp.Description,
			ObjectCategoryID: comp.ObjectCategoryID,
		}

		for _, q := range comp.FeedbackQuestionaires {
			qd := performance.FeedbackQuestionaireData{
				FeedbackQuestionaireID: q.FeedbackQuestionaireID,
				Question:               q.Question,
				Description:            q.Description,
				PmsCompetencyID:        q.PmsCompetencyID,
			}

			for _, opt := range q.Options {
				qd.Options = append(qd.Options, performance.FeedbackQuestionaireOptionData{
					FeedbackQuestionaireOptionID: opt.FeedbackQuestionaireOptionID,
					OptionStatement:              opt.OptionStatement,
					Description:                  opt.Description,
					Score:                         opt.Score,
					QuestionID:                    opt.QuestionID,
				})
			}

			cd.FeedbackQuestionaires = append(cd.FeedbackQuestionaires, qd)
		}

		data = append(data, cd)
	}

	resp.StaffID = staffID
	resp.PmsCompetencyData = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// CompetencyGapClosureSetup -- manages competency gap closure records.
// Mirrors .NET CompetencyGapClosureSetup.
// =========================================================================

func (cr *competencyReviewService) CompetencyGapClosureSetup(ctx context.Context, req *performance.CompetencyGapClosureRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.RecordStatus {
	case enums.OperationAdd.String():
		gap := performance.CompetencyGapClosure{
			StaffID:             req.StaffID,
			MaxPoints:           req.MaxPoints,
			FinalScore:          req.FinalScore,
			ReviewPeriodID:      req.ReviewPeriodID,
			ObjectiveCategoryID: req.ObjectiveCategoryID,
		}
		gap.RecordStatus = enums.StatusActive.String()
		gap.IsActive = true

		if err := cr.db.WithContext(ctx).Create(&gap).Error; err != nil {
			return resp, fmt.Errorf("creating competency gap closure: %w", err)
		}

		resp.ID = gap.CompetencyGapClosureID
		resp.Message = "competency gap closure created"

	case enums.OperationUpdate.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyGapClosure{}).
			Where("competency_gap_closure_id = ?", req.CompetencyGapClosureID).
			Updates(map[string]interface{}{
				"max_points":  req.MaxPoints,
				"final_score": req.FinalScore,
			})
		resp.ID = req.CompetencyGapClosureID
		resp.Message = "competency gap closure updated"

	case enums.OperationCancel.String():
		cr.db.WithContext(ctx).Model(&performance.CompetencyGapClosure{}).
			Where("competency_gap_closure_id = ?", req.CompetencyGapClosureID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CompetencyGapClosureID
		resp.Message = "competency gap closure cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for competency gap closure")
	}

	return resp, nil
}

// =========================================================================
// Initiate360Review -- initiates 360-degree review for staff members.
// Mirrors .NET Initiate360Review.
// =========================================================================

func (cr *competencyReviewService) Initiate360Review(ctx context.Context, req *performance.Initiate360ReviewRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	for _, staffID := range req.StaffID {
		// Check if feedback already exists
		var existing performance.CompetencyReviewFeedback
		if err := cr.db.WithContext(ctx).
			Where("staff_id = ? AND review_period_id = ? AND record_status != ?",
				staffID, req.ReviewPeriodID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			// Already exists, skip
			cr.log.Info().Str("staffID", staffID).Msg("360 review already initiated, skipping")
			continue
		}

		feedback := performance.CompetencyReviewFeedback{
			StaffID:        staffID,
			ReviewPeriodID: req.ReviewPeriodID,
		}
		feedback.RecordStatus = enums.StatusActive.String()
		feedback.IsActive = true

		if err := cr.db.WithContext(ctx).Create(&feedback).Error; err != nil {
			cr.log.Error().Err(err).Str("staffID", staffID).Msg("failed to initiate 360 review")
			continue
		}

		// Get the staff's line manager and peers to add as reviewers
		if cr.parent.erpEmployeeSvc != nil {
			cr.addDefaultReviewers(ctx, feedback.CompetencyReviewFeedbackID, staffID)
		}

		cr.log.Info().
			Str("staffID", staffID).
			Str("feedbackID", feedback.CompetencyReviewFeedbackID).
			Msg("360 review initiated")
	}

	resp.Message = "360 review initiated successfully"
	return resp, nil
}

// Complete360Review -- completes all 360 reviews for a review period.
// Mirrors .NET Complete360Review.
func (cr *competencyReviewService) Complete360Review(ctx context.Context, req *performance.Complete360ReviewRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	// Get all active feedbacks for the review period
	var feedbacks []performance.CompetencyReviewFeedback
	cr.db.WithContext(ctx).
		Where("review_period_id = ? AND record_status = ?", req.ReviewPeriodID, enums.StatusActive.String()).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Find(&feedbacks)

	for _, fb := range feedbacks {
		// Calculate final score from all reviewer ratings
		totalScore := 0.0
		totalRaters := 0

		for _, reviewer := range fb.CompetencyReviewers {
			if reviewer.RecordStatus == enums.StatusCompleted.String() {
				if reviewer.FinalRating > 0 {
					totalScore += reviewer.FinalRating
					totalRaters++
				}
			}
		}

		if totalRaters > 0 {
			fb.FinalScore = totalScore / float64(totalRaters)
		}

		fb.RecordStatus = enums.StatusCompleted.String()

		cr.db.WithContext(ctx).Save(&fb)
	}

	resp.Message = "360 reviews completed successfully"
	return resp, nil
}

// =========================================================================
// Internal helpers
// =========================================================================

// completeReviewerFeedback calculates the reviewer's final rating
// from their individual competency ratings.
func (cr *competencyReviewService) completeReviewerFeedback(ctx context.Context, reviewerID string) {
	var reviewer performance.CompetencyReviewer
	if err := cr.db.WithContext(ctx).
		Preload("CompetencyReviewerRatings").
		Where("competency_reviewer_id = ?", reviewerID).
		First(&reviewer).Error; err != nil {
		cr.log.Error().Err(err).Str("reviewerID", reviewerID).Msg("reviewer not found for completion")
		return
	}

	// Calculate average rating
	if len(reviewer.CompetencyReviewerRatings) > 0 {
		sum := 0.0
		for _, r := range reviewer.CompetencyReviewerRatings {
			sum += r.Rating
		}
		reviewer.FinalRating = sum / float64(len(reviewer.CompetencyReviewerRatings))
	}

	now := time.Now().UTC()
	_ = now
	reviewer.RecordStatus = enums.StatusCompleted.String()

	cr.db.WithContext(ctx).Save(&reviewer)

	// Close the associated feedback request
	var request performance.FeedbackRequestLog
	if err := cr.db.WithContext(ctx).
		Where("reference_id = ? AND feedback_request_type = ? AND record_status = ?",
			reviewer.CompetencyReviewFeedbackID, enums.FeedbackRequest360ReviewFeedback, enums.StatusActive.String()).
		First(&request).Error; err == nil {
		completedAt := time.Now().UTC()
		request.RecordStatus = enums.StatusCompleted.String()
		request.TimeCompleted = &completedAt
		cr.db.WithContext(ctx).Save(&request)
	}
}

// addDefaultReviewers attempts to auto-add reviewers (supervisor, peers, subordinates).
func (cr *competencyReviewService) addDefaultReviewers(ctx context.Context, feedbackID, staffID string) {
	// Get subordinates (this gives us the supervisor info)
	if subordinateInfo, err := cr.parent.erpEmployeeSvc.GetHeadSubordinates(ctx, staffID); err == nil && subordinateInfo != nil {
		cr.log.Debug().
			Str("staffID", staffID).
			Msg("auto-adding reviewers: ERP integration needed for specific reviewer assignment")
	}

	// Get peers
	if peers, err := cr.parent.erpEmployeeSvc.GetEmployeePeers(ctx, staffID); err == nil && peers != nil {
		cr.log.Debug().
			Str("staffID", staffID).
			Msg("auto-adding peer reviewers: ERP integration needed for specific peer IDs")
	}
}

func (cr *competencyReviewService) mapFeedbackToData(ctx context.Context, fb performance.CompetencyReviewFeedback) performance.CompetencyReviewFeedbackData {
	data := performance.CompetencyReviewFeedbackData{
		CompetencyReviewFeedbackID: fb.CompetencyReviewFeedbackID,
		StaffID:                    fb.StaffID,
		MaxPoints:                  fb.MaxPoints,
		FinalScore:                 fb.FinalScore,
		ReviewPeriodID:             fb.ReviewPeriodID,
		RecordStatusName:           fb.RecordStatus,
	}

	if fb.MaxPoints > 0 {
		data.FinalScorePercentage = (fb.FinalScore / fb.MaxPoints) * 100
	}

	// Enrich staff name
	if cr.parent.erpEmployeeSvc != nil {
		if detail, err := cr.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, fb.StaffID); err == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				data.StaffName = nameHolder.GetFullName()
			}
		}
	}

	// Map reviewers
	for _, r := range fb.CompetencyReviewers {
		data.CompetencyReviewers = append(data.CompetencyReviewers, cr.mapReviewerToData(ctx, r))
	}

	return data
}

func (cr *competencyReviewService) mapReviewerToData(ctx context.Context, r performance.CompetencyReviewer) performance.CompetencyReviewerData {
	data := performance.CompetencyReviewerData{
		CompetencyReviewerID:       r.CompetencyReviewerID,
		ReviewStaffID:              r.ReviewStaffID,
		FinalRating:                r.FinalRating,
		CompetencyReviewFeedbackID: r.CompetencyReviewFeedbackID,
		RecordStatusName:           r.RecordStatus,
	}

	// Map ratings
	for _, rating := range r.CompetencyReviewerRatings {
		rd := performance.CompetencyReviewerRatingData{
			CompetencyReviewerRatingID:   rating.CompetencyReviewerRatingID,
			PmsCompetencyID:              rating.PmsCompetencyID,
			FeedbackQuestionaireOptionID: rating.FeedbackQuestionaireOptionID,
			Rating:                       rating.Rating,
			CompetencyReviewerID:         rating.CompetencyReviewerID,
		}

		if rating.PmsCompetency != nil {
			rd.PmsCompetency = &performance.PmsCompetencyData{
				PmsCompetencyID:  rating.PmsCompetency.PmsCompetencyID,
				Name:             rating.PmsCompetency.Name,
				Description:      rating.PmsCompetency.Description,
				ObjectCategoryID: rating.PmsCompetency.ObjectCategoryID,
			}
		}

		if rating.FeedbackQuestionaireOption != nil {
			rd.FeedbackQuestionaireOption = &performance.FeedbackQuestionaireOptionData{
				FeedbackQuestionaireOptionID: rating.FeedbackQuestionaireOption.FeedbackQuestionaireOptionID,
				OptionStatement:              rating.FeedbackQuestionaireOption.OptionStatement,
				Description:                  rating.FeedbackQuestionaireOption.Description,
				Score:                         rating.FeedbackQuestionaireOption.Score,
				QuestionID:                    rating.FeedbackQuestionaireOption.QuestionID,
			}
		}

		data.CompetencyReviewerRatings = append(data.CompetencyReviewerRatings, rd)
	}

	return data
}
