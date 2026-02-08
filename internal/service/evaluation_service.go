package service

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// evaluationService handles period objective evaluations and department
// objective evaluations for review periods. Also provides retrieval of
// work products per enterprise objective and departments per objective.
//
// Mirrors .NET methods:
//   - ReviewPeriodObjectiveEvaluation (full workflow: Draft/Add/Update/
//     Approve/Reject/Return/ReSubmit/Cancel)
//   - GetReviewPeriodObjectiveEvaluation
//   - GetReviewPeriodObjectiveEvaluations
//   - ReviewPeriodDepartmentObjectiveEvaluation (full workflow)
//   - GetReviewPeriodDepartmentObjectiveEvaluation
//   - GetReviewPeriodDepartmentObjectiveEvaluations
//   - GetWorkProductsPerEnterpriseObjective
//   - GetDepartmentsbyObjective
// ---------------------------------------------------------------------------

type evaluationService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	parent *performanceManagementService

	periodObjEvalRepo     *repository.PMSRepository[performance.PeriodObjectiveEvaluation]
	periodObjDeptEvalRepo *repository.PMSRepository[performance.PeriodObjectiveDepartmentEvaluation]
	periodObjRepo         *repository.PMSRepository[performance.PeriodObjective]
	enterpriseObjRepo     *repository.PMSRepository[performance.EnterpriseObjective]
	deptObjRepo           *repository.PMSRepository[performance.DepartmentObjective]
	workProductRepo       *repository.PMSRepository[performance.WorkProduct]
	opObjWPRepo           *repository.PMSRepository[performance.OperationalObjectiveWorkProduct]
	plannedObjRepo        *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	reviewPeriodRepo      *repository.PMSRepository[performance.PerformanceReviewPeriod]
}

func newEvaluationService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *evaluationService {
	return &evaluationService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "evaluation").Logger(),
		parent: parent,

		periodObjEvalRepo:     repository.NewPMSRepository[performance.PeriodObjectiveEvaluation](db),
		periodObjDeptEvalRepo: repository.NewPMSRepository[performance.PeriodObjectiveDepartmentEvaluation](db),
		periodObjRepo:         repository.NewPMSRepository[performance.PeriodObjective](db),
		enterpriseObjRepo:     repository.NewPMSRepository[performance.EnterpriseObjective](db),
		deptObjRepo:           repository.NewPMSRepository[performance.DepartmentObjective](db),
		workProductRepo:       repository.NewPMSRepository[performance.WorkProduct](db),
		opObjWPRepo:           repository.NewPMSRepository[performance.OperationalObjectiveWorkProduct](db),
		plannedObjRepo:        repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
		reviewPeriodRepo:      repository.NewPMSRepository[performance.PerformanceReviewPeriod](db),
	}
}

// =========================================================================
// ReviewPeriodObjectiveEvaluation -- full workflow.
// Mirrors .NET ReviewPeriodObjectiveEvaluation(model, operationType).
// =========================================================================

func (e *evaluationService) ReviewPeriodObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}
	resp.Message = "an error occurred"

	switch req.RecordStatus {
	case enums.OperationDraft.String():
		return e.draftObjectiveEvaluation(ctx, req)
	case enums.OperationAdd.String():
		return e.addObjectiveEvaluation(ctx, req)
	case enums.OperationUpdate.String():
		return e.updateObjectiveEvaluation(ctx, req)
	case enums.OperationCommitDraft.String():
		return e.commitDraftObjectiveEvaluation(ctx, req)
	case enums.OperationApprove.String():
		return e.approveObjectiveEvaluation(ctx, req)
	case enums.OperationReject.String():
		return e.rejectObjectiveEvaluation(ctx, req)
	case enums.OperationReturn.String():
		return e.returnObjectiveEvaluation(ctx, req)
	case enums.OperationReSubmit.String():
		return e.resubmitObjectiveEvaluation(ctx, req)
	case enums.OperationCancel.String():
		return e.cancelObjectiveEvaluation(ctx, req)
	default:
		resp.HasError = true
		resp.Message = fmt.Sprintf("unsupported operation: %s", req.RecordStatus)
		return resp, fmt.Errorf("unsupported operation: %s", req.RecordStatus)
	}
}

func (e *evaluationService) draftObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	// Validate review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := e.db.WithContext(ctx).
		Where("period_id = ?", req.ReviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, fmt.Errorf("review period not found: %w", err)
	}

	// Validate period objective exists
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", req.EnterpriseObjectiveID, req.ReviewPeriodID).
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found for the given review period"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	eval := performance.PeriodObjectiveEvaluation{
		TotalOutcomeScore: req.TotalOutcomeScore,
		OutcomeScore:      req.OutcomeScore,
		PeriodObjectiveID: periodObj.PeriodObjectiveID,
	}
	eval.RecordStatus = enums.StatusDraft.String()
	eval.IsActive = true
	eval.CreatedBy = req.CreatedBy

	if err := e.db.WithContext(ctx).Create(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to save draft evaluation"
		return resp, fmt.Errorf("creating draft evaluation: %w", err)
	}

	data := e.mapEvaluationToData(eval, &periodObj, req.ReviewPeriodID, reviewPeriod.Name)
	resp.ObjectiveEvaluation = &data
	resp.Message = "draft evaluation saved successfully"
	return resp, nil
}

func (e *evaluationService) addObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	// Validate review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := e.db.WithContext(ctx).
		Where("period_id = ?", req.ReviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, fmt.Errorf("review period not found: %w", err)
	}

	// Validate period objective
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", req.EnterpriseObjectiveID, req.ReviewPeriodID).
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found for the given review period"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	// Check for existing non-cancelled evaluation
	var existing performance.PeriodObjectiveEvaluation
	err := e.db.WithContext(ctx).
		Where("period_objective_id = ? AND record_status != ?",
			periodObj.PeriodObjectiveID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		resp.HasError = true
		resp.Message = "an evaluation already exists for this period objective"
		return resp, fmt.Errorf("evaluation already exists for period objective %s", periodObj.PeriodObjectiveID)
	}

	eval := performance.PeriodObjectiveEvaluation{
		TotalOutcomeScore: req.TotalOutcomeScore,
		OutcomeScore:      req.OutcomeScore,
		PeriodObjectiveID: periodObj.PeriodObjectiveID,
	}
	eval.RecordStatus = enums.StatusPendingApproval.String()
	eval.IsActive = true
	eval.CreatedBy = req.CreatedBy

	if err := e.db.WithContext(ctx).Create(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to create evaluation"
		return resp, fmt.Errorf("creating evaluation: %w", err)
	}

	e.log.Info().
		Str("evaluationID", eval.PeriodObjectiveEvaluationID).
		Str("periodObjectiveID", periodObj.PeriodObjectiveID).
		Msg("period objective evaluation created")

	data := e.mapEvaluationToData(eval, &periodObj, req.ReviewPeriodID, reviewPeriod.Name)
	resp.ObjectiveEvaluation = &data
	resp.Message = "evaluation added successfully"
	return resp, nil
}

func (e *evaluationService) updateObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		Preload("PeriodObjective").
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusDraft.String() &&
		eval.RecordStatus != enums.StatusReturned.String() {
		resp.HasError = true
		resp.Message = "evaluation can only be updated when in Draft or Returned status"
		return resp, fmt.Errorf("cannot update evaluation in status %s", eval.RecordStatus)
	}

	eval.TotalOutcomeScore = req.TotalOutcomeScore
	eval.OutcomeScore = req.OutcomeScore
	eval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to update evaluation"
		return resp, fmt.Errorf("updating evaluation: %w", err)
	}

	data := e.mapEvaluationToData(eval, eval.PeriodObjective, req.ReviewPeriodID, "")
	resp.ObjectiveEvaluation = &data
	resp.Message = "evaluation updated successfully"
	return resp, nil
}

func (e *evaluationService) commitDraftObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusDraft.String() {
		resp.HasError = true
		resp.Message = "only draft evaluations can be submitted"
		return resp, fmt.Errorf("cannot commit evaluation in status %s", eval.RecordStatus)
	}

	eval.RecordStatus = enums.StatusPendingApproval.String()
	eval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to submit draft evaluation"
		return resp, fmt.Errorf("committing draft evaluation: %w", err)
	}

	resp.Message = "draft evaluation submitted for approval"
	return resp, nil
}

func (e *evaluationService) approveObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		Preload("PeriodObjective").
		Preload("PeriodObjective.Objective").
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be approved"
		return resp, fmt.Errorf("cannot approve evaluation in status %s", eval.RecordStatus)
	}

	eval.RecordStatus = enums.StatusActive.String()
	eval.IsApproved = true
	eval.ApprovedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to approve evaluation"
		return resp, fmt.Errorf("approving evaluation: %w", err)
	}

	// Update the period objective with the evaluation score
	if eval.PeriodObjective != nil {
		if eval.PeriodObjective.Objective != nil {
			// Propagate the outcome score to the enterprise objective's
			// evaluation tracking fields if needed.
			e.log.Info().
				Str("evaluationID", eval.PeriodObjectiveEvaluationID).
				Str("objectiveID", eval.PeriodObjective.ObjectiveID).
				Float64("outcomeScore", eval.OutcomeScore).
				Msg("period objective evaluation approved")
		}
	}

	resp.Message = "evaluation approved successfully"
	return resp, nil
}

func (e *evaluationService) rejectObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be rejected"
		return resp, fmt.Errorf("cannot reject evaluation in status %s", eval.RecordStatus)
	}

	eval.RecordStatus = enums.StatusRejected.String()
	eval.IsRejected = true
	eval.RejectedBy = req.UpdatedBy
	eval.RejectionReason = req.RejectionReason

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to reject evaluation"
		return resp, fmt.Errorf("rejecting evaluation: %w", err)
	}

	resp.Message = "evaluation rejected"
	return resp, nil
}

func (e *evaluationService) returnObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be returned"
		return resp, fmt.Errorf("cannot return evaluation in status %s", eval.RecordStatus)
	}

	eval.RecordStatus = enums.StatusReturned.String()
	eval.RejectionReason = req.RejectionReason
	eval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to return evaluation"
		return resp, fmt.Errorf("returning evaluation: %w", err)
	}

	resp.Message = "evaluation returned for revision"
	return resp, nil
}

func (e *evaluationService) resubmitObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	if eval.RecordStatus != enums.StatusReturned.String() &&
		eval.RecordStatus != enums.StatusRejected.String() {
		resp.HasError = true
		resp.Message = "only returned or rejected evaluations can be re-submitted"
		return resp, fmt.Errorf("cannot resubmit evaluation in status %s", eval.RecordStatus)
	}

	eval.TotalOutcomeScore = req.TotalOutcomeScore
	eval.OutcomeScore = req.OutcomeScore
	eval.RecordStatus = enums.StatusPendingApproval.String()
	eval.IsRejected = false
	eval.RejectionReason = ""
	eval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to re-submit evaluation"
		return resp, fmt.Errorf("re-submitting evaluation: %w", err)
	}

	resp.Message = "evaluation re-submitted for approval"
	return resp, nil
}

func (e *evaluationService) cancelObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveEvaluationRequestModel,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}

	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_evaluation_id = ?", req.PeriodObjectiveEvaluationID).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	eval.RecordStatus = enums.StatusCancelled.String()
	eval.IsActive = false
	eval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to cancel evaluation"
		return resp, fmt.Errorf("cancelling evaluation: %w", err)
	}

	resp.Message = "evaluation cancelled"
	return resp, nil
}

// =========================================================================
// GetReviewPeriodObjectiveEvaluation -- single evaluation by objective.
// Mirrors .NET GetReviewPeriodObjectiveEvaluation.
// =========================================================================

func (e *evaluationService) GetReviewPeriodObjectiveEvaluation(
	ctx context.Context,
	reviewPeriodID, enterpriseObjectiveID string,
) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationResponseVm{}
	resp.Message = "an error occurred"

	// Find the period objective
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", enterpriseObjectiveID, reviewPeriodID).
		Preload("Objective").
		Preload("ReviewPeriod").
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	// Find evaluation
	var eval performance.PeriodObjectiveEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_id = ? AND record_status != ?",
			periodObj.PeriodObjectiveID, enums.StatusCancelled.String()).
		First(&eval).Error; err != nil {
		resp.HasError = true
		resp.Message = "evaluation not found for this period objective"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	reviewPeriodName := ""
	if periodObj.ReviewPeriod != nil {
		reviewPeriodName = periodObj.ReviewPeriod.Name
	}

	data := e.mapEvaluationToData(eval, &periodObj, reviewPeriodID, reviewPeriodName)
	resp.ObjectiveEvaluation = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetReviewPeriodObjectiveEvaluations -- all evaluations for a review period.
// Mirrors .NET GetReviewPeriodObjectiveEvaluations.
// =========================================================================

func (e *evaluationService) GetReviewPeriodObjectiveEvaluations(
	ctx context.Context,
	reviewPeriodID string,
) (performance.PeriodObjectiveEvaluationListResponseVm, error) {
	resp := performance.PeriodObjectiveEvaluationListResponseVm{}
	resp.Message = "an error occurred"

	// Get all period objectives for this review period
	var periodObjs []performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("Objective").
		Preload("ReviewPeriod").
		Find(&periodObjs).Error; err != nil {
		resp.HasError = true
		return resp, err
	}

	// Build a map of period objective ID -> period objective
	periodObjMap := make(map[string]*performance.PeriodObjective)
	var periodObjIDs []string
	for i := range periodObjs {
		periodObjMap[periodObjs[i].PeriodObjectiveID] = &periodObjs[i]
		periodObjIDs = append(periodObjIDs, periodObjs[i].PeriodObjectiveID)
	}

	// Get all evaluations for these period objectives
	var evals []performance.PeriodObjectiveEvaluation
	if len(periodObjIDs) > 0 {
		e.db.WithContext(ctx).
			Where("period_objective_id IN ? AND record_status != ?",
				periodObjIDs, enums.StatusCancelled.String()).
			Find(&evals)
	}

	var dataList []performance.PeriodObjectiveEvaluationData
	for _, eval := range evals {
		po := periodObjMap[eval.PeriodObjectiveID]
		reviewPeriodName := ""
		if po != nil && po.ReviewPeriod != nil {
			reviewPeriodName = po.ReviewPeriod.Name
		}
		data := e.mapEvaluationToData(eval, po, reviewPeriodID, reviewPeriodName)
		dataList = append(dataList, data)
	}

	resp.ObjectiveEvaluations = dataList
	resp.TotalRecords = len(dataList)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// ReviewPeriodDepartmentObjectiveEvaluation -- full workflow.
// Mirrors .NET ReviewPeriodDepartmentObjectiveEvaluation(model, operationType).
// =========================================================================

func (e *evaluationService) ReviewPeriodDepartmentObjectiveEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}
	resp.Message = "an error occurred"

	switch req.RecordStatus {
	case enums.OperationDraft.String():
		return e.draftDeptEvaluation(ctx, req)
	case enums.OperationAdd.String():
		return e.addDeptEvaluation(ctx, req)
	case enums.OperationUpdate.String():
		return e.updateDeptEvaluation(ctx, req)
	case enums.OperationCommitDraft.String():
		return e.commitDraftDeptEvaluation(ctx, req)
	case enums.OperationApprove.String():
		return e.approveDeptEvaluation(ctx, req)
	case enums.OperationReject.String():
		return e.rejectDeptEvaluation(ctx, req)
	case enums.OperationReturn.String():
		return e.returnDeptEvaluation(ctx, req)
	case enums.OperationReSubmit.String():
		return e.resubmitDeptEvaluation(ctx, req)
	case enums.OperationCancel.String():
		return e.cancelDeptEvaluation(ctx, req)
	default:
		resp.HasError = true
		resp.Message = fmt.Sprintf("unsupported operation: %s", req.RecordStatus)
		return resp, fmt.Errorf("unsupported operation: %s", req.RecordStatus)
	}
}

func (e *evaluationService) draftDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	// Find the period objective
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", req.EnterpriseObjectiveID, req.ReviewPeriodID).
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	deptEval := performance.PeriodObjectiveDepartmentEvaluation{
		OverallOutcomeScored: req.OverallOutcomeScored,
		AllocatedOutcome:     req.AllocatedOutcome,
		OutcomeScore:         req.OutcomeScore,
		DepartmentID:         req.DepartmentID,
		PeriodObjectiveID:    periodObj.PeriodObjectiveID,
	}
	deptEval.RecordStatus = enums.StatusDraft.String()
	deptEval.IsActive = true
	deptEval.CreatedBy = req.CreatedBy

	if err := e.db.WithContext(ctx).Create(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to save draft department evaluation"
		return resp, fmt.Errorf("creating draft department evaluation: %w", err)
	}

	data := e.mapDeptEvaluationToData(deptEval, &periodObj, req.ReviewPeriodID, "")
	resp.DepartmentObjectiveEvaluation = &data
	resp.Message = "draft department evaluation saved successfully"
	return resp, nil
}

func (e *evaluationService) addDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	// Find the period objective
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", req.EnterpriseObjectiveID, req.ReviewPeriodID).
		Preload("Objective").
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	// Check for existing non-cancelled department evaluation
	var existing performance.PeriodObjectiveDepartmentEvaluation
	err := e.db.WithContext(ctx).
		Where("period_objective_id = ? AND department_id = ? AND record_status != ?",
			periodObj.PeriodObjectiveID, req.DepartmentID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		resp.HasError = true
		resp.Message = "a department evaluation already exists for this objective and department"
		return resp, fmt.Errorf("department evaluation already exists")
	}

	deptEval := performance.PeriodObjectiveDepartmentEvaluation{
		OverallOutcomeScored: req.OverallOutcomeScored,
		AllocatedOutcome:     req.AllocatedOutcome,
		OutcomeScore:         req.OutcomeScore,
		DepartmentID:         req.DepartmentID,
		PeriodObjectiveID:    periodObj.PeriodObjectiveID,
	}
	deptEval.RecordStatus = enums.StatusPendingApproval.String()
	deptEval.IsActive = true
	deptEval.CreatedBy = req.CreatedBy

	if err := e.db.WithContext(ctx).Create(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to create department evaluation"
		return resp, fmt.Errorf("creating department evaluation: %w", err)
	}

	e.log.Info().
		Str("deptEvalID", deptEval.PeriodObjectiveDepartmentEvaluationID).
		Int("departmentID", req.DepartmentID).
		Msg("department objective evaluation created")

	data := e.mapDeptEvaluationToData(deptEval, &periodObj, req.ReviewPeriodID, "")
	resp.DepartmentObjectiveEvaluation = &data
	resp.Message = "department evaluation added successfully"
	return resp, nil
}

func (e *evaluationService) updateDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusDraft.String() &&
		deptEval.RecordStatus != enums.StatusReturned.String() {
		resp.HasError = true
		resp.Message = "department evaluation can only be updated in Draft or Returned status"
		return resp, fmt.Errorf("cannot update in status %s", deptEval.RecordStatus)
	}

	deptEval.OverallOutcomeScored = req.OverallOutcomeScored
	deptEval.AllocatedOutcome = req.AllocatedOutcome
	deptEval.OutcomeScore = req.OutcomeScore
	deptEval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to update department evaluation"
		return resp, fmt.Errorf("updating department evaluation: %w", err)
	}

	resp.Message = "department evaluation updated successfully"
	return resp, nil
}

func (e *evaluationService) commitDraftDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusDraft.String() {
		resp.HasError = true
		resp.Message = "only draft department evaluations can be submitted"
		return resp, fmt.Errorf("cannot commit in status %s", deptEval.RecordStatus)
	}

	deptEval.RecordStatus = enums.StatusPendingApproval.String()
	deptEval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to submit draft department evaluation"
		return resp, fmt.Errorf("committing draft: %w", err)
	}

	resp.Message = "draft department evaluation submitted for approval"
	return resp, nil
}

func (e *evaluationService) approveDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be approved"
		return resp, fmt.Errorf("cannot approve in status %s", deptEval.RecordStatus)
	}

	deptEval.RecordStatus = enums.StatusActive.String()
	deptEval.IsApproved = true
	deptEval.ApprovedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to approve department evaluation"
		return resp, fmt.Errorf("approving: %w", err)
	}

	e.log.Info().
		Str("deptEvalID", deptEval.PeriodObjectiveDepartmentEvaluationID).
		Msg("department objective evaluation approved")

	resp.Message = "department evaluation approved successfully"
	return resp, nil
}

func (e *evaluationService) rejectDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be rejected"
		return resp, fmt.Errorf("cannot reject in status %s", deptEval.RecordStatus)
	}

	deptEval.RecordStatus = enums.StatusRejected.String()
	deptEval.IsRejected = true
	deptEval.RejectedBy = req.UpdatedBy
	deptEval.RejectionReason = req.RejectionReason

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to reject department evaluation"
		return resp, fmt.Errorf("rejecting: %w", err)
	}

	resp.Message = "department evaluation rejected"
	return resp, nil
}

func (e *evaluationService) returnDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusPendingApproval.String() {
		resp.HasError = true
		resp.Message = "only pending-approval evaluations can be returned"
		return resp, fmt.Errorf("cannot return in status %s", deptEval.RecordStatus)
	}

	deptEval.RecordStatus = enums.StatusReturned.String()
	deptEval.RejectionReason = req.RejectionReason
	deptEval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to return department evaluation"
		return resp, fmt.Errorf("returning: %w", err)
	}

	resp.Message = "department evaluation returned for revision"
	return resp, nil
}

func (e *evaluationService) resubmitDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	if deptEval.RecordStatus != enums.StatusReturned.String() &&
		deptEval.RecordStatus != enums.StatusRejected.String() {
		resp.HasError = true
		resp.Message = "only returned or rejected evaluations can be re-submitted"
		return resp, fmt.Errorf("cannot resubmit in status %s", deptEval.RecordStatus)
	}

	deptEval.OverallOutcomeScored = req.OverallOutcomeScored
	deptEval.AllocatedOutcome = req.AllocatedOutcome
	deptEval.OutcomeScore = req.OutcomeScore
	deptEval.RecordStatus = enums.StatusPendingApproval.String()
	deptEval.IsRejected = false
	deptEval.RejectionReason = ""
	deptEval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to re-submit department evaluation"
		return resp, fmt.Errorf("re-submitting: %w", err)
	}

	resp.Message = "department evaluation re-submitted for approval"
	return resp, nil
}

func (e *evaluationService) cancelDeptEvaluation(
	ctx context.Context,
	req *performance.PeriodObjectiveDepartmentEvaluationRequestModel,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_department_evaluation_id = ?",
			req.PeriodObjectiveDepartmentEvaluationID).
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	deptEval.RecordStatus = enums.StatusCancelled.String()
	deptEval.IsActive = false
	deptEval.UpdatedBy = req.UpdatedBy

	if err := e.db.WithContext(ctx).Save(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "failed to cancel department evaluation"
		return resp, fmt.Errorf("cancelling: %w", err)
	}

	resp.Message = "department evaluation cancelled"
	return resp, nil
}

// =========================================================================
// GetReviewPeriodDepartmentObjectiveEvaluation -- single department evaluation.
// Mirrors .NET GetReviewPeriodDepartmentObjectiveEvaluation.
// =========================================================================

func (e *evaluationService) GetReviewPeriodDepartmentObjectiveEvaluation(
	ctx context.Context,
	reviewPeriodID, enterpriseObjectiveID string,
	departmentID int,
) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationResponseVm{}
	resp.Message = "an error occurred"

	// Find the period objective
	var periodObj performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("objective_id = ? AND review_period_id = ?", enterpriseObjectiveID, reviewPeriodID).
		Preload("Objective").
		Preload("ReviewPeriod").
		First(&periodObj).Error; err != nil {
		resp.HasError = true
		resp.Message = "period objective not found"
		return resp, fmt.Errorf("period objective not found: %w", err)
	}

	var deptEval performance.PeriodObjectiveDepartmentEvaluation
	if err := e.db.WithContext(ctx).
		Where("period_objective_id = ? AND department_id = ? AND record_status != ?",
			periodObj.PeriodObjectiveID, departmentID, enums.StatusCancelled.String()).
		Preload("Department").
		First(&deptEval).Error; err != nil {
		resp.HasError = true
		resp.Message = "department evaluation not found"
		return resp, fmt.Errorf("department evaluation not found: %w", err)
	}

	reviewPeriodName := ""
	if periodObj.ReviewPeriod != nil {
		reviewPeriodName = periodObj.ReviewPeriod.Name
	}

	data := e.mapDeptEvaluationToData(deptEval, &periodObj, reviewPeriodID, reviewPeriodName)
	resp.DepartmentObjectiveEvaluation = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetReviewPeriodDepartmentObjectiveEvaluations -- all department evaluations.
// Mirrors .NET GetReviewPeriodDepartmentObjectiveEvaluations.
// =========================================================================

func (e *evaluationService) GetReviewPeriodDepartmentObjectiveEvaluations(
	ctx context.Context,
	reviewPeriodID string,
) (performance.PeriodObjectiveDepartmentEvaluationListResponseVm, error) {
	resp := performance.PeriodObjectiveDepartmentEvaluationListResponseVm{}
	resp.Message = "an error occurred"

	// Get all period objectives for this review period
	var periodObjs []performance.PeriodObjective
	if err := e.db.WithContext(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("Objective").
		Preload("ReviewPeriod").
		Find(&periodObjs).Error; err != nil {
		resp.HasError = true
		return resp, err
	}

	periodObjMap := make(map[string]*performance.PeriodObjective)
	var periodObjIDs []string
	for i := range periodObjs {
		periodObjMap[periodObjs[i].PeriodObjectiveID] = &periodObjs[i]
		periodObjIDs = append(periodObjIDs, periodObjs[i].PeriodObjectiveID)
	}

	// Get all department evaluations
	var deptEvals []performance.PeriodObjectiveDepartmentEvaluation
	if len(periodObjIDs) > 0 {
		e.db.WithContext(ctx).
			Where("period_objective_id IN ? AND record_status != ?",
				periodObjIDs, enums.StatusCancelled.String()).
			Preload("Department").
			Find(&deptEvals)
	}

	var dataList []performance.PeriodObjectiveDepartmentEvaluationData
	for _, de := range deptEvals {
		po := periodObjMap[de.PeriodObjectiveID]
		reviewPeriodName := ""
		if po != nil && po.ReviewPeriod != nil {
			reviewPeriodName = po.ReviewPeriod.Name
		}
		data := e.mapDeptEvaluationToData(de, po, reviewPeriodID, reviewPeriodName)
		dataList = append(dataList, data)
	}

	resp.DepartmentObjectiveEvaluations = dataList
	resp.TotalRecords = len(dataList)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetWorkProductsPerEnterpriseObjective -- work products linked to an
// enterprise objective via operational objective work products.
// Mirrors .NET GetWorkProductsPerEnterpriseObjective.
// =========================================================================

func (e *evaluationService) GetWorkProductsPerEnterpriseObjective(
	ctx context.Context,
	enterpriseObjectiveID, reviewPeriodID string,
) (performance.ObjectiveWorkProductListResponseVm, error) {
	resp := performance.ObjectiveWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	// Get review period for date range
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := e.db.WithContext(ctx).
		Where("period_id = ?", reviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		resp.HasError = true
		resp.Message = "review period not found"
		return resp, fmt.Errorf("review period not found: %w", err)
	}

	// Find planned objectives that reference this enterprise objective
	// at any level (department, division, office, enterprise)
	var plannedObjs []performance.ReviewPeriodIndividualPlannedObjective
	e.db.WithContext(ctx).
		Where("review_period_id = ? AND record_status = ?",
			reviewPeriodID, enums.StatusActive.String()).
		Find(&plannedObjs)

	// Collect planned objective IDs that resolve to the enterprise objective
	var matchingPlannedObjIDs []string
	for _, po := range plannedObjs {
		objID := po.ObjectiveID
		matchesEnterprise := false

		switch po.ObjectiveLevel {
		case enums.ObjectiveLevelEnterprise:
			matchesEnterprise = (objID == enterpriseObjectiveID)
		case enums.ObjectiveLevelDepartment:
			var deptObj performance.DepartmentObjective
			if err := e.db.WithContext(ctx).
				Where("department_objective_id = ?", objID).
				First(&deptObj).Error; err == nil {
				matchesEnterprise = (deptObj.EnterpriseObjectiveID == enterpriseObjectiveID)
			}
		case enums.ObjectiveLevelDivision:
			var divObj performance.DivisionObjective
			if err := e.db.WithContext(ctx).
				Where("division_objective_id = ?", objID).
				Preload("DepartmentObjective").
				First(&divObj).Error; err == nil && divObj.DepartmentObjective != nil {
				matchesEnterprise = (divObj.DepartmentObjective.EnterpriseObjectiveID == enterpriseObjectiveID)
			}
		case enums.ObjectiveLevelOffice:
			var offObj performance.OfficeObjective
			if err := e.db.WithContext(ctx).
				Where("office_objective_id = ?", objID).
				Preload("DivisionObjective").
				Preload("DivisionObjective.DepartmentObjective").
				First(&offObj).Error; err == nil &&
				offObj.DivisionObjective != nil &&
				offObj.DivisionObjective.DepartmentObjective != nil {
				matchesEnterprise = (offObj.DivisionObjective.DepartmentObjective.EnterpriseObjectiveID == enterpriseObjectiveID)
			}
		}

		if matchesEnterprise {
			matchingPlannedObjIDs = append(matchingPlannedObjIDs, po.PlannedObjectiveID)
		}
	}

	if len(matchingPlannedObjIDs) == 0 {
		resp.ObjectiveWorkProducts = []performance.ObjectiveWorkProductData{}
		resp.TotalRecords = 0
		resp.Message = "operation completed successfully"
		return resp, nil
	}

	// Get operational objective work products linked to these planned objectives
	var opObjWPs []performance.OperationalObjectiveWorkProduct
	e.db.WithContext(ctx).
		Where("planned_objective_id IN ?", matchingPlannedObjIDs).
		Preload("WorkProduct").
		Find(&opObjWPs)

	excluded := excludedStatuses()
	var wpDataList []performance.ObjectiveWorkProductData
	for _, opwp := range opObjWPs {
		if opwp.WorkProduct == nil {
			continue
		}
		wp := opwp.WorkProduct
		if isInStatuses(wp.RecordStatus, excluded) {
			continue
		}

		data := performance.ObjectiveWorkProductData{
			ReviewPeriodID:          reviewPeriodID,
			ReviewPeriod:            reviewPeriod.Name,
			ObjectiveID:             enterpriseObjectiveID,
			WorkProductDefinitionID: opwp.WorkProductDefinitionID,
			WorkProductName:         wp.Name,
			Description:             wp.Description,
			Deliverables:            wp.Deliverables,
			StaffID:                 wp.StaffID,
		}
		wpDataList = append(wpDataList, data)
	}

	resp.ObjectiveWorkProducts = wpDataList
	resp.TotalRecords = len(wpDataList)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetDepartmentsbyObjective -- returns departments linked to an enterprise
// objective via department objectives.
// Mirrors .NET GetDepartmentsbyObjective.
// =========================================================================

func (e *evaluationService) GetDepartmentsbyObjective(
	ctx context.Context,
	enterpriseObjectiveID string,
) (performance.DepartmentListResponseVm, error) {
	resp := performance.DepartmentListResponseVm{}
	resp.Message = "an error occurred"
	resp.EnterpriseObjectiveID = enterpriseObjectiveID

	var deptObjs []performance.DepartmentObjective
	if err := e.db.WithContext(ctx).
		Where("enterprise_objective_id = ? AND record_status = ?",
			enterpriseObjectiveID, enums.StatusActive.String()).
		Preload("Department").
		Find(&deptObjs).Error; err != nil {
		resp.HasError = true
		return resp, err
	}

	var departments []interface{}
	seen := make(map[int]bool)
	for _, do := range deptObjs {
		if do.Department == nil {
			continue
		}
		if seen[do.DepartmentID] {
			continue
		}
		seen[do.DepartmentID] = true
		departments = append(departments, do.Department)
	}

	resp.Departments = departments
	resp.TotalRecords = len(departments)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// Internal helpers
// =========================================================================

func (e *evaluationService) mapEvaluationToData(
	eval performance.PeriodObjectiveEvaluation,
	periodObj *performance.PeriodObjective,
	reviewPeriodID, reviewPeriodName string,
) performance.PeriodObjectiveEvaluationData {
	data := performance.PeriodObjectiveEvaluationData{
		BaseWorkFlowVm: performance.BaseWorkFlowVm{
			ID:              eval.ID,
			RecordStatus:    eval.RecordStatus,
			CreatedAt:       eval.CreatedAt,
			UpdatedAt:       eval.UpdatedAt,
			CreatedBy:       eval.CreatedBy,
			UpdatedBy:       eval.UpdatedBy,
			IsActive:        eval.IsActive,
			ApprovedBy:      eval.ApprovedBy,
			DateApproved:    eval.DateApproved,
			IsApproved:      eval.IsApproved,
			IsRejected:      eval.IsRejected,
			RejectedBy:      eval.RejectedBy,
			RejectionReason: eval.RejectionReason,
			DateRejected:    eval.DateRejected,
		},
		PeriodObjectiveEvaluationID: eval.PeriodObjectiveEvaluationID,
		TotalOutcomeScore:           eval.TotalOutcomeScore,
		OutcomeScore:                eval.OutcomeScore,
		ReviewPeriodID:              reviewPeriodID,
		ReviewPeriod:                reviewPeriodName,
	}

	if periodObj != nil {
		data.PeriodObjectiveID = periodObj.PeriodObjectiveID
		data.EnterpriseObjectiveID = periodObj.ObjectiveID
		if periodObj.Objective != nil {
			data.EnterpriseObjective = periodObj.Objective.Name
		}
	}

	return data
}

func (e *evaluationService) mapDeptEvaluationToData(
	deptEval performance.PeriodObjectiveDepartmentEvaluation,
	periodObj *performance.PeriodObjective,
	reviewPeriodID, reviewPeriodName string,
) performance.PeriodObjectiveDepartmentEvaluationData {
	data := performance.PeriodObjectiveDepartmentEvaluationData{
		BaseWorkFlowVm: performance.BaseWorkFlowVm{
			ID:              deptEval.ID,
			RecordStatus:    deptEval.RecordStatus,
			CreatedAt:       deptEval.CreatedAt,
			UpdatedAt:       deptEval.UpdatedAt,
			CreatedBy:       deptEval.CreatedBy,
			UpdatedBy:       deptEval.UpdatedBy,
			IsActive:        deptEval.IsActive,
			ApprovedBy:      deptEval.ApprovedBy,
			DateApproved:    deptEval.DateApproved,
			IsApproved:      deptEval.IsApproved,
			IsRejected:      deptEval.IsRejected,
			RejectedBy:      deptEval.RejectedBy,
			RejectionReason: deptEval.RejectionReason,
			DateRejected:    deptEval.DateRejected,
		},
		PeriodObjectiveDepartmentEvaluationID: deptEval.PeriodObjectiveDepartmentEvaluationID,
		OverallOutcomeScored:                  deptEval.OverallOutcomeScored,
		AllocatedOutcome:                      deptEval.AllocatedOutcome,
		OutcomeScore:                          deptEval.OutcomeScore,
		DepartmentID:                          deptEval.DepartmentID,
		ReviewPeriodID:                        reviewPeriodID,
		ReviewPeriod:                          reviewPeriodName,
	}

	if periodObj != nil {
		data.PeriodObjectiveID = periodObj.PeriodObjectiveID
		data.EnterpriseObjectiveID = periodObj.ObjectiveID
	}

	if deptEval.Department != nil {
		data.DepartmentName = deptEval.Department.DepartmentName
	}

	// Resolve department objective name if available
	if periodObj != nil {
		var deptObj performance.DepartmentObjective
		if err := e.db.
			Where("enterprise_objective_id = ? AND department_id = ?",
				periodObj.ObjectiveID, deptEval.DepartmentID).
			First(&deptObj).Error; err == nil {
			data.DepartmentObjectiveID = deptObj.DepartmentObjectiveID
			data.DepartmentObjective = deptObj.Name
		}
	}

	return data
}
