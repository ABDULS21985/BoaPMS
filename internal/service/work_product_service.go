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
// workProductService handles full work product lifecycle management.
//
// Mirrors .NET methods:
//   - WorkProductSetup (full workflow)
//   - ProjectAssignedWorkProductSetup / CommitteeAssignedWorkProductSetup
//   - GetWorkProduct / GetProjectWorkProduct / GetCommitteeWorkProduct
//   - GetProjectAssignedWorkProducts / GetCommitteeAssignedWorkProducts
//   - GetOperationalWorkProducts / GetStaffWorkProducts / GetAllStaffWorkProducts
//   - GetObjectiveWorkProducts
//   - WorkProductTaskSetup / GetWorkProductTasks
//   - ReCalculateWorkProductPoints
//   - WorkProductEvaluation / GetWorkProductEvaluation / InitiateWorkProductReEvaluation
// ---------------------------------------------------------------------------

type workProductService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	parent *performanceManagementService

	workProductRepo    *repository.PMSRepository[performance.WorkProduct]
	wpTaskRepo         *repository.PMSRepository[performance.WorkProductTask]
	wpEvalRepo         *repository.PMSRepository[performance.WorkProductEvaluation]
	evalOptionRepo     *repository.PMSRepository[performance.EvaluationOption]
	opObjWPRepo        *repository.PMSRepository[performance.OperationalObjectiveWorkProduct]
	projectWPRepo      *repository.PMSRepository[performance.ProjectWorkProduct]
	committeeWPRepo    *repository.PMSRepository[performance.CommitteeWorkProduct]
	projAssignedWPRepo *repository.PMSRepository[performance.ProjectAssignedWorkProduct]
	comAssignedWPRepo  *repository.PMSRepository[performance.CommitteeAssignedWorkProduct]
	plannedObjRepo     *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	periodScoreRepo    *repository.PMSRepository[performance.PeriodScore]
}

func newWorkProductService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *workProductService {
	return &workProductService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "work_product").Logger(),
		parent: parent,

		workProductRepo:    repository.NewPMSRepository[performance.WorkProduct](db),
		wpTaskRepo:         repository.NewPMSRepository[performance.WorkProductTask](db),
		wpEvalRepo:         repository.NewPMSRepository[performance.WorkProductEvaluation](db),
		evalOptionRepo:     repository.NewPMSRepository[performance.EvaluationOption](db),
		opObjWPRepo:        repository.NewPMSRepository[performance.OperationalObjectiveWorkProduct](db),
		projectWPRepo:      repository.NewPMSRepository[performance.ProjectWorkProduct](db),
		committeeWPRepo:    repository.NewPMSRepository[performance.CommitteeWorkProduct](db),
		projAssignedWPRepo: repository.NewPMSRepository[performance.ProjectAssignedWorkProduct](db),
		comAssignedWPRepo:  repository.NewPMSRepository[performance.CommitteeAssignedWorkProduct](db),
		plannedObjRepo:     repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
		periodScoreRepo:    repository.NewPMSRepository[performance.PeriodScore](db),
	}
}

// =========================================================================
// WorkProductSetup -- full work product lifecycle workflow.
// Mirrors .NET WorkProductSetup with OperationType switch.
// =========================================================================

func (ws *workProductService) WorkProductSetup(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}
	resp.Message = "an error occurred"

	switch req.Status {
	case enums.OperationDraft.String():
		return ws.saveDraftWorkProduct(ctx, req)
	case enums.OperationAdd.String():
		return ws.addWorkProduct(ctx, req)
	case enums.OperationUpdate.String():
		return ws.updateWorkProduct(ctx, req)
	case enums.OperationApprove.String():
		return ws.approveWorkProduct(ctx, req)
	case enums.OperationReject.String():
		return ws.rejectWorkProduct(ctx, req)
	case enums.OperationReturn.String():
		return ws.returnWorkProduct(ctx, req)
	case enums.OperationReSubmit.String():
		return ws.reSubmitWorkProduct(ctx, req)
	case enums.OperationClose.String():
		return ws.closeWorkProduct(ctx, req)
	case enums.OperationPause.String():
		return ws.pauseWorkProduct(ctx, req)
	case enums.OperationCancel.String():
		return ws.cancelWorkProduct(ctx, req)
	case enums.OperationAccept.String():
		return ws.acceptWorkProduct(ctx, req)
	case enums.OperationComplete.String():
		return ws.completeWorkProduct(ctx, req)
	case enums.OperationSuspend.String():
		return ws.suspendWorkProduct(ctx, req)
	case enums.OperationResume.String():
		return ws.resumeWorkProduct(ctx, req)
	default:
		return ws.addWorkProduct(ctx, req)
	}
}

func (ws *workProductService) saveDraftWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	wp := performance.WorkProduct{
		Name:            req.Name,
		Description:     req.Description,
		MaxPoint:        req.MaxPoint,
		WorkProductType: enums.WorkProductType(req.WorkProductType),
		IsSelfCreated:   req.IsSelfCreated,
		StaffID:         req.StaffID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Deliverables:    req.Deliverables,
	}
	wp.RecordStatus = enums.StatusDraft.String()
	wp.IsActive = true
	wp.CreatedBy = req.CreatedBy

	if err := ws.db.WithContext(ctx).Create(&wp).Error; err != nil {
		return resp, fmt.Errorf("saving draft work product: %w", err)
	}

	// Link to planned objective
	if req.PlannedObjectiveID != "" {
		link := performance.OperationalObjectiveWorkProduct{
			WorkProductID:      wp.WorkProductID,
			PlannedObjectiveID: req.PlannedObjectiveID,
		}
		ws.db.WithContext(ctx).Create(&link)
	}

	resp.ID = wp.WorkProductID
	resp.Message = "work product draft saved successfully"
	return resp, nil
}

func (ws *workProductService) addWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	// Validate planned objective
	if req.PlannedObjectiveID != "" {
		var plannedObj performance.ReviewPeriodIndividualPlannedObjective
		if err := ws.db.WithContext(ctx).
			Where("planned_objective_id = ? AND staff_id = ? AND record_status = ?",
				req.PlannedObjectiveID, req.StaffID, enums.StatusActive.String()).
			First(&plannedObj).Error; err != nil {
			return resp, fmt.Errorf("planned objective not found or not active: %w", err)
		}
	}

	wp := performance.WorkProduct{
		Name:            req.Name,
		Description:     req.Description,
		MaxPoint:        req.MaxPoint,
		WorkProductType: enums.WorkProductType(req.WorkProductType),
		IsSelfCreated:   req.IsSelfCreated,
		StaffID:         req.StaffID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Deliverables:    req.Deliverables,
		Remark:          req.Remark,
	}
	wp.RecordStatus = enums.StatusPendingApproval.String()
	wp.IsActive = true
	wp.CreatedBy = req.CreatedBy

	if err := ws.db.WithContext(ctx).Create(&wp).Error; err != nil {
		return resp, fmt.Errorf("creating work product: %w", err)
	}

	// Link to planned objective
	if req.PlannedObjectiveID != "" {
		wpDefID := ""
		if req.WorkProductDefinitionID != nil {
			wpDefID = *req.WorkProductDefinitionID
		}
		link := performance.OperationalObjectiveWorkProduct{
			WorkProductID:           wp.WorkProductID,
			PlannedObjectiveID:      req.PlannedObjectiveID,
			WorkProductDefinitionID: wpDefID,
		}
		ws.db.WithContext(ctx).Create(&link)
	}

	// Link to project
	if req.ProjectID != "" {
		projWP := performance.ProjectWorkProduct{
			WorkProductID: wp.WorkProductID,
			ProjectID:     req.ProjectID,
		}
		if req.ProjectAssignedWorkProductID != nil && *req.ProjectAssignedWorkProductID != "" {
			projWP.ProjectAssignedWorkProductID = *req.ProjectAssignedWorkProductID
		}
		ws.db.WithContext(ctx).Create(&projWP)
	}

	// Link to committee
	if req.CommitteeID != "" {
		comWP := performance.CommitteeWorkProduct{
			WorkProductID: wp.WorkProductID,
			CommitteeID:   req.CommitteeID,
		}
		if req.CommitteeAssignedWorkProductID != nil && *req.CommitteeAssignedWorkProductID != "" {
			comWP.CommitteeAssignedWorkProductID = *req.CommitteeAssignedWorkProductID
		}
		ws.db.WithContext(ctx).Create(&comWP)
	}

	resp.ID = wp.WorkProductID
	resp.Message = "work product created successfully"
	return resp, nil
}

func (ws *workProductService) updateWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	var wp performance.WorkProduct
	if err := ws.db.WithContext(ctx).
		Where("work_product_id = ?", req.WorkProductID).
		First(&wp).Error; err != nil {
		return resp, fmt.Errorf("work product not found: %w", err)
	}

	wp.Name = req.Name
	wp.Description = req.Description
	wp.MaxPoint = req.MaxPoint
	wp.StartDate = req.StartDate
	wp.EndDate = req.EndDate
	wp.Deliverables = req.Deliverables
	wp.Remark = req.Remark

	if err := ws.db.WithContext(ctx).Save(&wp).Error; err != nil {
		return resp, fmt.Errorf("updating work product: %w", err)
	}

	resp.ID = wp.WorkProductID
	resp.Message = "work product updated successfully"
	return resp, nil
}

func (ws *workProductService) approveWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status":   enums.StatusActive.String(),
			"is_active":       true,
			"is_approved":     true,
			"approved_by":     req.UpdatedBy,
			"date_approved":   now,
			"approver_comment": req.ApproverComment,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product approved successfully"
	return resp, nil
}

func (ws *workProductService) rejectWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status":    enums.StatusRejected.String(),
			"is_rejected":      true,
			"rejected_by":      req.UpdatedBy,
			"rejection_reason": req.RejectionReason,
			"date_rejected":    now,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product rejected successfully"
	return resp, nil
}

func (ws *workProductService) returnWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	var wp performance.WorkProduct
	if err := ws.db.WithContext(ctx).
		Where("work_product_id = ?", req.WorkProductID).First(&wp).Error; err != nil {
		return resp, fmt.Errorf("work product not found: %w", err)
	}

	wp.RecordStatus = enums.StatusReturned.String()
	wp.NoReturned++
	wp.ApproverComment = req.ApproverComment

	if err := ws.db.WithContext(ctx).Save(&wp).Error; err != nil {
		return resp, fmt.Errorf("returning work product: %w", err)
	}

	resp.ID = wp.WorkProductID
	resp.Message = "work product returned successfully"
	return resp, nil
}

func (ws *workProductService) reSubmitWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPendingApproval.String(),
			"is_rejected":   false,
			"remark":         req.Remark,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product re-submitted successfully"
	return resp, nil
}

func (ws *workProductService) closeWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusClosed.String(),
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product closed successfully"
	return resp, nil
}

func (ws *workProductService) pauseWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPaused.String(),
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product paused successfully"
	return resp, nil
}

func (ws *workProductService) cancelWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusCancelled.String(),
			"is_active":     false,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product cancelled successfully"
	return resp, nil
}

func (ws *workProductService) acceptWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status":     enums.StatusAwaitingEvaluation.String(),
			"acceptance_comment": req.AcceptanceComment,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product accepted and awaiting evaluation"
	return resp, nil
}

func (ws *workProductService) completeWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status":     enums.StatusPendingAcceptance.String(),
			"completion_date":   now,
			"remark":            req.Remark,
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product completed and pending acceptance"
	return resp, nil
}

func (ws *workProductService) suspendWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusSuspended.String(),
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product suspended successfully"
	return resp, nil
}

func (ws *workProductService) resumeWorkProduct(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ws.db.WithContext(ctx).Model(&performance.WorkProduct{}).
		Where("work_product_id = ?", req.WorkProductID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusActive.String(),
		})

	resp.ID = req.WorkProductID
	resp.Message = "work product resumed successfully"
	return resp, nil
}

// =========================================================================
// ProjectAssignedWorkProductSetup -- manages project-level work product
// assignments. Mirrors .NET ProjectAssignedWorkProductSetup.
// =========================================================================

func (ws *workProductService) ProjectAssignedWorkProductSetup(ctx context.Context, req *performance.ProjectAssignedWorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.Status {
	case enums.OperationAdd.String():
		awp := performance.ProjectAssignedWorkProduct{
			WorkProductDefinitionID: req.WorkProductDefinitionID,
			Name:                    req.Name,
			Description:             req.Description,
			ProjectID:               req.ProjectID,
			StartDate:               req.StartDate,
			EndDate:                 req.EndDate,
			Deliverables:            req.Deliverables,
		}
		awp.RecordStatus = enums.StatusActive.String()
		awp.IsActive = true

		if err := ws.db.WithContext(ctx).Create(&awp).Error; err != nil {
			return resp, fmt.Errorf("creating project assigned work product: %w", err)
		}

		resp.ID = awp.ProjectAssignedWorkProductID
		resp.Message = "project assigned work product created successfully"

	case enums.OperationUpdate.String():
		ws.db.WithContext(ctx).Model(&performance.ProjectAssignedWorkProduct{}).
			Where("project_assigned_work_product_id = ?", req.ProjectAssignedWorkProductID).
			Updates(map[string]interface{}{
				"name":        req.Name,
				"description": req.Description,
				"start_date":  req.StartDate,
				"end_date":    req.EndDate,
				"deliverables": req.Deliverables,
			})
		resp.ID = req.ProjectAssignedWorkProductID
		resp.Message = "project assigned work product updated"

	case enums.OperationApprove.String():
		ws.db.WithContext(ctx).Model(&performance.ProjectAssignedWorkProduct{}).
			Where("project_assigned_work_product_id = ?", req.ProjectAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_approved":   true,
			})
		resp.ID = req.ProjectAssignedWorkProductID
		resp.Message = "project assigned work product approved"

	case enums.OperationReject.String():
		ws.db.WithContext(ctx).Model(&performance.ProjectAssignedWorkProduct{}).
			Where("project_assigned_work_product_id = ?", req.ProjectAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status":    enums.StatusRejected.String(),
				"is_rejected":      true,
				"rejection_reason": req.RejectionReason,
			})
		resp.ID = req.ProjectAssignedWorkProductID
		resp.Message = "project assigned work product rejected"

	case enums.OperationCancel.String():
		ws.db.WithContext(ctx).Model(&performance.ProjectAssignedWorkProduct{}).
			Where("project_assigned_work_product_id = ?", req.ProjectAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.ProjectAssignedWorkProductID
		resp.Message = "project assigned work product cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for project assigned work product")
	}

	return resp, nil
}

// CommitteeAssignedWorkProductSetup -- manages committee-level work product
// assignments. Mirrors .NET CommitteeAssignedWorkProductSetup.
func (ws *workProductService) CommitteeAssignedWorkProductSetup(ctx context.Context, req *performance.CommitteeAssignedWorkProductRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.Status {
	case enums.OperationAdd.String():
		awp := performance.CommitteeAssignedWorkProduct{
			WorkProductDefinitionID: req.WorkProductDefinitionID,
			Name:                    req.Name,
			Description:             req.Description,
			CommitteeID:             req.CommitteeID,
			StartDate:               req.StartDate,
			EndDate:                 req.EndDate,
			Deliverables:            req.Deliverables,
		}
		awp.RecordStatus = enums.StatusActive.String()
		awp.IsActive = true

		if err := ws.db.WithContext(ctx).Create(&awp).Error; err != nil {
			return resp, fmt.Errorf("creating committee assigned work product: %w", err)
		}

		resp.ID = awp.CommitteeAssignedWorkProductID
		resp.Message = "committee assigned work product created successfully"

	case enums.OperationUpdate.String():
		ws.db.WithContext(ctx).Model(&performance.CommitteeAssignedWorkProduct{}).
			Where("committee_assigned_work_product_id = ?", req.CommitteeAssignedWorkProductID).
			Updates(map[string]interface{}{
				"name":        req.Name,
				"description": req.Description,
				"start_date":  req.StartDate,
				"end_date":    req.EndDate,
				"deliverables": req.Deliverables,
			})
		resp.ID = req.CommitteeAssignedWorkProductID
		resp.Message = "committee assigned work product updated"

	case enums.OperationApprove.String():
		ws.db.WithContext(ctx).Model(&performance.CommitteeAssignedWorkProduct{}).
			Where("committee_assigned_work_product_id = ?", req.CommitteeAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_approved":   true,
			})
		resp.ID = req.CommitteeAssignedWorkProductID
		resp.Message = "committee assigned work product approved"

	case enums.OperationReject.String():
		ws.db.WithContext(ctx).Model(&performance.CommitteeAssignedWorkProduct{}).
			Where("committee_assigned_work_product_id = ?", req.CommitteeAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status":    enums.StatusRejected.String(),
				"is_rejected":      true,
				"rejection_reason": req.RejectionReason,
			})
		resp.ID = req.CommitteeAssignedWorkProductID
		resp.Message = "committee assigned work product rejected"

	case enums.OperationCancel.String():
		ws.db.WithContext(ctx).Model(&performance.CommitteeAssignedWorkProduct{}).
			Where("committee_assigned_work_product_id = ?", req.CommitteeAssignedWorkProductID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CommitteeAssignedWorkProductID
		resp.Message = "committee assigned work product cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for committee assigned work product")
	}

	return resp, nil
}

// =========================================================================
// GetWorkProduct -- retrieves a single work product with all associations.
// Mirrors .NET GetWorkProduct.
// =========================================================================

func (ws *workProductService) GetWorkProduct(ctx context.Context, workProductID string) (performance.WorkProductResponseVm, error) {
	resp := performance.WorkProductResponseVm{}
	resp.Message = "an error occurred"

	var wp performance.WorkProduct
	err := ws.db.WithContext(ctx).
		Preload("WorkProductTasks").
		Preload("OperationalObjectiveWorkProducts").
		Preload("OperationalObjectiveWorkProducts.PlannedObjective").
		Preload("ProjectWorkProducts").
		Preload("ProjectWorkProducts.Project").
		Preload("CommitteeWorkProducts").
		Preload("CommitteeWorkProducts.Committee").
		Where("work_product_id = ?", workProductID).
		First(&wp).Error
	if err != nil {
		ws.log.Error().Err(err).Str("workProductID", workProductID).Msg("work product not found")
		resp.HasError = true
		resp.Message = "work product not found"
		return resp, fmt.Errorf("work product not found: %w", err)
	}

	data := ws.mapWorkProductToData(ctx, wp)
	resp.WorkProduct = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetProjectWorkProducts -- retrieves work products for a project.
// Mirrors .NET GetProjectWorkProduct.
func (ws *workProductService) GetProjectWorkProducts(ctx context.Context, projectID string) (performance.ProjectWorkProductListResponseVm, error) {
	resp := performance.ProjectWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var pwps []performance.ProjectWorkProduct
	err := ws.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("WorkProduct").
		Preload("WorkProduct.WorkProductTasks").
		Find(&pwps).Error
	if err != nil {
		ws.log.Error().Err(err).Str("projectID", projectID).Msg("failed to get project work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.ProjectWorkProductDataResponse
	for _, pwp := range pwps {
		d := performance.ProjectWorkProductDataResponse{
			ProjectWorkProductID: pwp.ProjectWorkProductID,
			WorkProductID:        pwp.WorkProductID,
			ProjectID:            pwp.ProjectID,
		}
		if pwp.WorkProduct != nil {
			wpData := ws.mapWorkProductToData(ctx, *pwp.WorkProduct)
			d.WorkProduct = &wpData
		}
		data = append(data, d)
	}

	resp.ProjectWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetProjectAssignedWorkProducts -- retrieves assigned work products for a project.
// Mirrors .NET GetProjectAssignedWorkProducts.
func (ws *workProductService) GetProjectAssignedWorkProducts(ctx context.Context, projectID string) (performance.ProjectAssignedWorkProductListResponseVm, error) {
	resp := performance.ProjectAssignedWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var awps []performance.ProjectAssignedWorkProduct
	err := ws.db.WithContext(ctx).
		Where("project_id = ? AND record_status != ?", projectID, enums.StatusCancelled.String()).
		Find(&awps).Error
	if err != nil {
		ws.log.Error().Err(err).Str("projectID", projectID).Msg("failed to get project assigned work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.ProjectAssignedWorkProductData
	for _, awp := range awps {
		d := performance.ProjectAssignedWorkProductData{
			ProjectAssignedWorkProductID: awp.ProjectAssignedWorkProductID,
			WorkProductDefinitionID:      awp.WorkProductDefinitionID,
			Name:                         awp.Name,
			Description:                  awp.Description,
			ProjectID:                    awp.ProjectID,
			ReviewPeriodID:               awp.ReviewPeriodID,
			StartDate:                    awp.StartDate,
			EndDate:                      awp.EndDate,
			Deliverables:                 awp.Deliverables,
		}
		data = append(data, d)
	}

	resp.ProjectWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCommitteeWorkProducts -- retrieves work products for a committee.
// Mirrors .NET GetCommitteeWorkProduct.
func (ws *workProductService) GetCommitteeWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeWorkProductListResponseVm, error) {
	resp := performance.CommitteeWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var cwps []performance.CommitteeWorkProduct
	err := ws.db.WithContext(ctx).
		Where("committee_id = ?", committeeID).
		Preload("WorkProduct").
		Preload("WorkProduct.WorkProductTasks").
		Find(&cwps).Error
	if err != nil {
		ws.log.Error().Err(err).Str("committeeID", committeeID).Msg("failed to get committee work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CommitteeWorkProductDataResponse
	for _, cwp := range cwps {
		d := performance.CommitteeWorkProductDataResponse{
			CommitteeWorkProductID: cwp.CommitteeWorkProductID,
			WorkProductID:          cwp.WorkProductID,
			CommitteeID:            cwp.CommitteeID,
		}
		if cwp.WorkProduct != nil {
			wpData := ws.mapWorkProductToData(ctx, *cwp.WorkProduct)
			d.WorkProduct = &wpData
		}
		data = append(data, d)
	}

	resp.CommitteeWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCommitteeAssignedWorkProducts -- retrieves assigned work products for a committee.
// Mirrors .NET GetCommitteeAssignedWorkProducts.
func (ws *workProductService) GetCommitteeAssignedWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeAssignedWorkProductListResponseVm, error) {
	resp := performance.CommitteeAssignedWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var awps []performance.CommitteeAssignedWorkProduct
	err := ws.db.WithContext(ctx).
		Where("committee_id = ? AND record_status != ?", committeeID, enums.StatusCancelled.String()).
		Find(&awps).Error
	if err != nil {
		ws.log.Error().Err(err).Str("committeeID", committeeID).Msg("failed to get committee assigned work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CommitteeAssignedWorkProductData
	for _, awp := range awps {
		d := performance.CommitteeAssignedWorkProductData{
			CommitteeAssignedWorkProductID: awp.CommitteeAssignedWorkProductID,
			WorkProductDefinitionID:        awp.WorkProductDefinitionID,
			Name:                           awp.Name,
			Description:                    awp.Description,
			CommitteeID:                    awp.CommitteeID,
			ReviewPeriodID:                 awp.ReviewPeriodID,
			StartDate:                      awp.StartDate,
			EndDate:                        awp.EndDate,
			Deliverables:                   awp.Deliverables,
		}
		data = append(data, d)
	}

	resp.CommitteeWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetOperationalWorkProducts -- retrieves operational work products for a planned objective.
// Mirrors .NET GetOperationalWorkProducts.
func (ws *workProductService) GetOperationalWorkProducts(ctx context.Context, plannedObjectiveID string) (performance.OperationalObjectiveWorkProductListResponseVm, error) {
	resp := performance.OperationalObjectiveWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var links []performance.OperationalObjectiveWorkProduct
	err := ws.db.WithContext(ctx).
		Where("planned_objective_id = ?", plannedObjectiveID).
		Preload("WorkProduct").
		Find(&links).Error
	if err != nil {
		ws.log.Error().Err(err).Str("plannedObjectiveID", plannedObjectiveID).Msg("failed to get operational work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.OperationalObjectiveWorkProductData
	for _, link := range links {
		d := performance.OperationalObjectiveWorkProductData{
			OperationalObjectiveWorkProductID: link.OperationalObjectiveWorkProductID,
			WorkProductID:                     link.WorkProductID,
			WorkProductDefinitionID:           link.WorkProductDefinitionID,
			PlannedObjectiveID:                link.PlannedObjectiveID,
		}
		data = append(data, d)
	}

	resp.OperationalObjectiveWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetStaffWorkProducts -- retrieves all work products for a staff member in a review period.
// Mirrors .NET GetStaffWorkProducts.
func (ws *workProductService) GetStaffWorkProducts(ctx context.Context, staffID, reviewPeriodID string) (performance.StaffWorkProductListResponseVm, error) {
	resp := performance.StaffWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var wps []performance.WorkProduct
	query := ws.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ?", staffID, excludedStatuses()).
		Preload("WorkProductTasks").
		Preload("OperationalObjectiveWorkProducts").
		Preload("OperationalObjectiveWorkProducts.PlannedObjective")

	if reviewPeriodID != "" {
		// Filter by review period through planned objective
		query = query.
			Joins("JOIN pms.operational_objective_work_products ON pms.operational_objective_work_products.work_product_id = pms.work_products.work_product_id").
			Joins("JOIN pms.review_period_individual_planned_objectives ON pms.review_period_individual_planned_objectives.planned_objective_id = pms.operational_objective_work_products.planned_objective_id").
			Where("pms.review_period_individual_planned_objectives.review_period_id = ?", reviewPeriodID)
	}

	err := query.Find(&wps).Error
	if err != nil {
		ws.log.Error().Err(err).Str("staffID", staffID).Msg("failed to get staff work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.WorkProductData
	for _, wp := range wps {
		data = append(data, ws.mapWorkProductToData(ctx, wp))
	}

	resp.WorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetAllStaffWorkProducts -- retrieves all work products for a staff member.
// Mirrors .NET GetAllStaffWorkProducts.
func (ws *workProductService) GetAllStaffWorkProducts(ctx context.Context, staffID string) (performance.StaffWorkProductListResponseVm, error) {
	return ws.GetStaffWorkProducts(ctx, staffID, "")
}

// GetObjectiveWorkProducts -- retrieves work products linked to an objective.
// Mirrors .NET GetObjectiveWorkProducts.
func (ws *workProductService) GetObjectiveWorkProducts(ctx context.Context, objectiveID string) (performance.ObjectiveWorkProductListResponseVm, error) {
	resp := performance.ObjectiveWorkProductListResponseVm{}
	resp.Message = "an error occurred"

	var links []performance.OperationalObjectiveWorkProduct
	err := ws.db.WithContext(ctx).
		Preload("WorkProduct").
		Preload("PlannedObjective").
		Where("planned_objective_id IN (?)",
			ws.db.Model(&performance.ReviewPeriodIndividualPlannedObjective{}).
				Select("planned_objective_id").
				Where("objective_id = ?", objectiveID)).
		Find(&links).Error
	if err != nil {
		ws.log.Error().Err(err).Str("objectiveID", objectiveID).Msg("failed to get objective work products")
		resp.HasError = true
		return resp, err
	}

	var data []performance.ObjectiveWorkProductData
	for _, link := range links {
		d := performance.ObjectiveWorkProductData{
			ObjectiveID: objectiveID,
		}
		if link.WorkProduct != nil {
			d.WorkProductName = link.WorkProduct.Name
			d.Description = link.WorkProduct.Description
			d.Deliverables = link.WorkProduct.Deliverables
			d.StaffID = link.WorkProduct.StaffID
		}
		if link.PlannedObjective != nil {
			d.ReviewPeriodID = link.PlannedObjective.ReviewPeriodID
		}
		d.WorkProductDefinitionID = link.WorkProductDefinitionID
		data = append(data, d)
	}

	resp.ObjectiveWorkProducts = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// WorkProductTaskSetup -- full work product task lifecycle workflow.
// Mirrors .NET WorkProductTaskSetup.
// =========================================================================

func (ws *workProductService) WorkProductTaskSetup(ctx context.Context, req *performance.WorkProductTaskRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.RecordStatus {
	case enums.OperationAdd.String():
		task := performance.WorkProductTask{
			Name:          req.Name,
			Description:   req.Description,
			StartDate:     req.StartDate,
			EndDate:       req.EndDate,
			WorkProductID: req.WorkProductID,
		}
		task.RecordStatus = enums.StatusActive.String()
		task.IsActive = true

		if err := ws.db.WithContext(ctx).Create(&task).Error; err != nil {
			return resp, fmt.Errorf("creating work product task: %w", err)
		}

		resp.ID = task.WorkProductTaskID
		resp.Message = "work product task created successfully"

	case enums.OperationUpdate.String():
		ws.db.WithContext(ctx).Model(&performance.WorkProductTask{}).
			Where("work_product_task_id = ?", req.WorkProductTaskID).
			Updates(map[string]interface{}{
				"name":        req.Name,
				"description": req.Description,
				"start_date":  req.StartDate,
				"end_date":    req.EndDate,
			})
		resp.ID = req.WorkProductTaskID
		resp.Message = "work product task updated"

	case enums.OperationComplete.String():
		now := time.Now().UTC()
		ws.db.WithContext(ctx).Model(&performance.WorkProductTask{}).
			Where("work_product_task_id = ?", req.WorkProductTaskID).
			Updates(map[string]interface{}{
				"record_status":   enums.StatusCompleted.String(),
				"completion_date": now,
			})
		resp.ID = req.WorkProductTaskID
		resp.Message = "work product task completed"

	case enums.OperationCancel.String():
		ws.db.WithContext(ctx).Model(&performance.WorkProductTask{}).
			Where("work_product_task_id = ?", req.WorkProductTaskID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.WorkProductTaskID
		resp.Message = "work product task cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for work product task")
	}

	return resp, nil
}

// GetWorkProductTasks -- retrieves tasks for a work product.
// Mirrors .NET GetWorkProductTasks.
func (ws *workProductService) GetWorkProductTasks(ctx context.Context, workProductID string) (performance.WorkProductTaskListResponseVm, error) {
	resp := performance.WorkProductTaskListResponseVm{}
	resp.Message = "an error occurred"

	var tasks []performance.WorkProductTask
	err := ws.db.WithContext(ctx).
		Where("work_product_id = ? AND record_status != ?", workProductID, enums.StatusCancelled.String()).
		Find(&tasks).Error
	if err != nil {
		ws.log.Error().Err(err).Str("workProductID", workProductID).Msg("failed to get work product tasks")
		resp.HasError = true
		return resp, err
	}

	var data []performance.WorkProductTaskData
	for _, t := range tasks {
		d := performance.WorkProductTaskData{
			WorkProductTaskID: t.WorkProductTaskID,
			Name:              t.Name,
			Description:       t.Description,
			WorkProductID:     t.WorkProductID,
		}
		if t.CompletionDate != nil {
			d.CompletionDate = *t.CompletionDate
		}
		data = append(data, d)
	}

	resp.WorkProductTasks = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetWorkProductTaskDetail -- retrieves a single work product task by ID.
// Mirrors .NET GetWorkProductTaskDetail.
func (ws *workProductService) GetWorkProductTaskDetail(ctx context.Context, taskID string) (performance.WorkProductTaskResponseVm, error) {
	resp := performance.WorkProductTaskResponseVm{}
	resp.Message = "an error occurred"

	var task performance.WorkProductTask
	err := ws.db.WithContext(ctx).
		Where("work_product_task_id = ?", taskID).
		First(&task).Error
	if err != nil {
		ws.log.Error().Err(err).Str("taskID", taskID).Msg("failed to get work product task detail")
		resp.HasError = true
		return resp, err
	}

	detail := performance.WorkProductTaskDetail{
		WorkProductTaskID: task.WorkProductTaskID,
		Name:              task.Name,
		Description:       task.Description,
		StartDate:         task.StartDate,
		EndDate:           task.EndDate,
		WorkProductID:     task.WorkProductID,
	}
	if task.CompletionDate != nil {
		detail.CompletionDate = *task.CompletionDate
	}
	resp.WorkProductTask = &detail
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// ReCalculateWorkProductPoints -- recalculates points for all work products
// of a staff member in a review period.
// Mirrors .NET ReCalculateWorkProductPoints.
// =========================================================================

func (ws *workProductService) ReCalculateWorkProductPoints(ctx context.Context, staffID, reviewPeriodID string) (performance.RecalculateWorkProductResponseVm, error) {
	resp := performance.RecalculateWorkProductResponseVm{}
	resp.Message = "an error occurred"

	// Get all active/closed work products for the staff
	var wps []performance.WorkProduct
	ws.db.WithContext(ctx).
		Where("staff_id = ? AND record_status NOT IN ?", staffID, excludedStatuses()).
		Find(&wps)

	// Recalculate period score based on work product evaluations
	var totalPoints float64
	var maxPoints float64
	for _, wp := range wps {
		maxPoints += wp.MaxPoint
		if wp.RecordStatus == enums.StatusClosed.String() {
			totalPoints += wp.FinalScore
		}
	}

	// Update period score
	var periodScore performance.PeriodScore
	err := ws.db.WithContext(ctx).
		Where("review_period_id = ? AND staff_id = ?", reviewPeriodID, staffID).
		First(&periodScore).Error

	if err == nil {
		periodScore.FinalScore = totalPoints
		if maxPoints > 0 {
			periodScore.ScorePercentage = (totalPoints / maxPoints) * 100
		}
		ws.db.WithContext(ctx).Save(&periodScore)
	}

	resp.StaffID = staffID
	resp.ReviewPeriodID = reviewPeriodID
	resp.Message = "work product points recalculated successfully"
	return resp, nil
}

// =========================================================================
// WorkProductEvaluation -- evaluates a work product with T/Q/O scores.
// Mirrors .NET WorkProductEvaluation (full workflow).
// =========================================================================

func (ws *workProductService) WorkProductEvaluation(ctx context.Context, req *performance.WorkProductEvaluationRequestModel) (performance.EvaluationResponseVm, error) {
	resp := performance.EvaluationResponseVm{}
	resp.Message = "an error occurred"

	// Validate work product
	var wp performance.WorkProduct
	if err := ws.db.WithContext(ctx).
		Where("work_product_id = ?", req.WorkProductID).
		First(&wp).Error; err != nil {
		return resp, fmt.Errorf("work product not found: %w", err)
	}

	if wp.RecordStatus != enums.StatusAwaitingEvaluation.String() &&
		wp.RecordStatus != enums.StatusReEvaluate.String() {
		return resp, fmt.Errorf("work product is not in a state that allows evaluation (current: %s)", wp.RecordStatus)
	}

	// Get evaluation option scores
	var timelinessOpt, qualityOpt, outputOpt performance.EvaluationOption
	if req.TimelinessEvaluationOptionID != "" {
		ws.db.WithContext(ctx).Where("evaluation_option_id = ?", req.TimelinessEvaluationOptionID).First(&timelinessOpt)
	}
	if req.QualityEvaluationOptionID != "" {
		ws.db.WithContext(ctx).Where("evaluation_option_id = ?", req.QualityEvaluationOptionID).First(&qualityOpt)
	}
	if req.OutputEvaluationOptionID != "" {
		ws.db.WithContext(ctx).Where("evaluation_option_id = ?", req.OutputEvaluationOptionID).First(&outputOpt)
	}

	timeliness := timelinessOpt.Score
	quality := qualityOpt.Score
	output := outputOpt.Score

	// Calculate: (T + Q + O) / 3 * maxPoint / 100
	avgScore := (timeliness + quality + output) / 3.0
	finalScore := wp.MaxPoint * avgScore / 100.0

	// Check if re-evaluation
	var existingEval performance.WorkProductEvaluation
	isReEval := false
	if err := ws.db.WithContext(ctx).
		Where("work_product_id = ?", req.WorkProductID).
		First(&existingEval).Error; err == nil {
		isReEval = true
	}

	if isReEval {
		existingEval.Timeliness = timeliness
		existingEval.TimelinessEvaluationOptionID = req.TimelinessEvaluationOptionID
		existingEval.Quality = quality
		existingEval.QualityEvaluationOptionID = req.QualityEvaluationOptionID
		existingEval.Output = output
		existingEval.OutputEvaluationOptionID = req.OutputEvaluationOptionID
		existingEval.Outcome = finalScore
		existingEval.EvaluatorStaffID = req.EvaluatorStaffID
		existingEval.IsReEvaluated = true

		ws.db.WithContext(ctx).Save(&existingEval)
	} else {
		eval := performance.WorkProductEvaluation{
			WorkProductID:                req.WorkProductID,
			Timeliness:                   timeliness,
			TimelinessEvaluationOptionID: req.TimelinessEvaluationOptionID,
			Quality:                      quality,
			QualityEvaluationOptionID:    req.QualityEvaluationOptionID,
			Output:                       output,
			OutputEvaluationOptionID:     req.OutputEvaluationOptionID,
			Outcome:                      finalScore,
			EvaluatorStaffID:             req.EvaluatorStaffID,
		}

		ws.db.WithContext(ctx).Create(&eval)
	}

	// Update work product
	now := time.Now().UTC()
	wp.FinalScore = finalScore
	wp.RecordStatus = enums.StatusClosed.String()
	wp.CompletionDate = &now
	ws.db.WithContext(ctx).Save(&wp)

	resp.IsSuccess = true
	resp.Score = finalScore
	resp.Message = "work product evaluated successfully"
	return resp, nil
}

// GetWorkProductEvaluation -- retrieves the evaluation for a work product.
// Mirrors .NET GetWorkProductEvaluation.
func (ws *workProductService) GetWorkProductEvaluation(ctx context.Context, workProductID string) (performance.WorkProductEvaluationResponseVm, error) {
	resp := performance.WorkProductEvaluationResponseVm{}
	resp.Message = "an error occurred"

	var eval performance.WorkProductEvaluation
	err := ws.db.WithContext(ctx).
		Preload("WorkProduct").
		Preload("TimelinessEvaluationOption").
		Preload("QualityEvaluationOption").
		Preload("OutputEvaluationOption").
		Where("work_product_id = ?", workProductID).
		First(&eval).Error
	if err != nil {
		ws.log.Error().Err(err).Str("workProductID", workProductID).Msg("work product evaluation not found")
		resp.HasError = true
		resp.Message = "evaluation not found"
		return resp, fmt.Errorf("evaluation not found: %w", err)
	}

	data := performance.WorkProductEvaluationDataResponse{
		WorkProductEvaluationID:      eval.WorkProductEvaluationID,
		WorkProductID:                eval.WorkProductID,
		Timeliness:                   eval.Timeliness,
		TimelinessEvaluationOptionID: eval.TimelinessEvaluationOptionID,
		Quality:                      eval.Quality,
		QualityEvaluationOptionID:    eval.QualityEvaluationOptionID,
		Output:                       eval.Output,
		OutputEvaluationOptionID:     eval.OutputEvaluationOptionID,
		Outcome:                      eval.Outcome,
		EvaluatorStaffID:             eval.EvaluatorStaffID,
		IsReEvaluated:                eval.IsReEvaluated,
	}

	// Map evaluation options
	if eval.TimelinessEvaluationOption != nil {
		data.TimelinessEvaluationOption = &performance.EvaluationOptionData{
			EvaluationOptionID: eval.TimelinessEvaluationOption.EvaluationOptionID,
			Name:               eval.TimelinessEvaluationOption.Name,
			Description:        eval.TimelinessEvaluationOption.Description,
			Score:              eval.TimelinessEvaluationOption.Score,
		}
	}
	if eval.QualityEvaluationOption != nil {
		data.QualityEvaluationOption = &performance.EvaluationOptionData{
			EvaluationOptionID: eval.QualityEvaluationOption.EvaluationOptionID,
			Name:               eval.QualityEvaluationOption.Name,
			Description:        eval.QualityEvaluationOption.Description,
			Score:              eval.QualityEvaluationOption.Score,
		}
	}
	if eval.OutputEvaluationOption != nil {
		data.OutputEvaluationOption = &performance.EvaluationOptionData{
			EvaluationOptionID: eval.OutputEvaluationOption.EvaluationOptionID,
			Name:               eval.OutputEvaluationOption.Name,
			Description:        eval.OutputEvaluationOption.Description,
			Score:              eval.OutputEvaluationOption.Score,
		}
	}

	resp.WorkProductEvaluation = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// InitiateWorkProductReEvaluation -- marks a work product for re-evaluation.
// Mirrors .NET InitiateWorkProductReEvaluation.
func (ws *workProductService) InitiateWorkProductReEvaluation(ctx context.Context, workProductID string) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	var wp performance.WorkProduct
	if err := ws.db.WithContext(ctx).
		Where("work_product_id = ?", workProductID).
		First(&wp).Error; err != nil {
		return resp, fmt.Errorf("work product not found: %w", err)
	}

	if wp.RecordStatus != enums.StatusClosed.String() {
		return resp, fmt.Errorf("work product must be closed to initiate re-evaluation")
	}

	wp.RecordStatus = enums.StatusReEvaluate.String()
	wp.ReEvaluationReInitiated = true

	if err := ws.db.WithContext(ctx).Save(&wp).Error; err != nil {
		return resp, fmt.Errorf("initiating re-evaluation: %w", err)
	}

	resp.ID = wp.WorkProductID
	resp.Message = "work product re-evaluation initiated"
	return resp, nil
}

// =========================================================================
// Internal helpers
// =========================================================================

func (ws *workProductService) mapWorkProductToData(ctx context.Context, wp performance.WorkProduct) performance.WorkProductData {
	data := performance.WorkProductData{
		WorkProductID:   wp.WorkProductID,
		Name:            wp.Name,
		Description:     wp.Description,
		MaxPoint:        wp.MaxPoint,
		WorkProductType: int(wp.WorkProductType),
		IsSelfCreated:   wp.IsSelfCreated,
		StaffID:         wp.StaffID,
		AcceptanceComment: wp.AcceptanceComment,
		StartDate:       wp.StartDate,
		EndDate:         wp.EndDate,
		Deliverables:    wp.Deliverables,
		FinalScore:      wp.FinalScore,
		NoReturned:      wp.NoReturned,
		ApproverComment: wp.ApproverComment,
	}
	data.RecordStatus = wp.RecordStatus
	data.IsActive = wp.IsActive

	if wp.CompletionDate != nil {
		data.CompletionDate = *wp.CompletionDate
	}

	// Map work product type name
	switch wp.WorkProductType {
	case enums.WorkProductTypeOperational:
		data.WorkProductTypeName = "Operational"
	case enums.WorkProductTypeProject:
		data.WorkProductTypeName = "Project"
	case enums.WorkProductTypeCommittee:
		data.WorkProductTypeName = "Committee"
	}

	// Map tasks
	for _, t := range wp.WorkProductTasks {
		td := performance.WorkProductTaskDetail{
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
		data.WorkProductTasks = append(data.WorkProductTasks, td)
	}

	// Extract planned objective link
	if len(wp.OperationalObjectiveWorkProducts) > 0 {
		link := wp.OperationalObjectiveWorkProducts[0]
		data.PlannedObjectiveID = link.PlannedObjectiveID
		data.WorkProductDefinitionID = link.WorkProductDefinitionID
		if link.PlannedObjective != nil {
			data.ObjectiveID = link.PlannedObjective.ObjectiveID
			data.ReviewPeriodID = link.PlannedObjective.ReviewPeriodID
		}
	}

	return data
}
