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
// committeeService handles full committee lifecycle management.
//
// Mirrors .NET methods:
//   - CommitteeSetup (full workflow: Draft/Add/Update/Approve/Reject/Return/
//     ReSubmit/Close/Pause/Cancel)
//   - GetCommittee / GetCommittees / GetCommitteesByChairperson
//   - CommitteeObjectiveSetup (full workflow)
//   - GetCommitteeObjectives
//   - CommitteeMembersSetup (full workflow)
//   - GetCommitteeMembers
//   - GetCommitteeWorkProductStaffList
//   - GetCommitteesAssigned / GetStaffCommittees
//   - ChangeCommitteeChairperson
// ---------------------------------------------------------------------------

type committeeService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	parent *performanceManagementService

	committeeRepo       *repository.PMSRepository[performance.Committee]
	committeeObjRepo    *repository.PMSRepository[performance.CommitteeObjective]
	committeeMemberRepo *repository.PMSRepository[performance.CommitteeMember]
	committeeWPRepo     *repository.PMSRepository[performance.CommitteeWorkProduct]
	comAssignedWPRepo   *repository.PMSRepository[performance.CommitteeAssignedWorkProduct]
	reviewPeriodRepo    *repository.PMSRepository[performance.PerformanceReviewPeriod]
	plannedObjRepo      *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
}

func newCommitteeService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *committeeService {
	return &committeeService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "committee").Logger(),
		parent: parent,

		committeeRepo:       repository.NewPMSRepository[performance.Committee](db),
		committeeObjRepo:    repository.NewPMSRepository[performance.CommitteeObjective](db),
		committeeMemberRepo: repository.NewPMSRepository[performance.CommitteeMember](db),
		committeeWPRepo:     repository.NewPMSRepository[performance.CommitteeWorkProduct](db),
		comAssignedWPRepo:   repository.NewPMSRepository[performance.CommitteeAssignedWorkProduct](db),
		reviewPeriodRepo:    repository.NewPMSRepository[performance.PerformanceReviewPeriod](db),
		plannedObjRepo:      repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
	}
}

// =========================================================================
// CommitteeSetup -- full committee lifecycle workflow.
// Mirrors .NET CommitteeSetup with OperationType switch.
// =========================================================================

func (cs *committeeService) CommitteeSetup(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}
	resp.Message = "an error occurred"

	switch req.Status {
	case enums.OperationDraft.String():
		return cs.saveDraftCommittee(ctx, req)
	case enums.OperationAdd.String():
		return cs.addCommittee(ctx, req)
	case enums.OperationUpdate.String():
		return cs.updateCommittee(ctx, req)
	case enums.OperationApprove.String():
		return cs.approveCommittee(ctx, req)
	case enums.OperationReject.String():
		return cs.rejectCommittee(ctx, req)
	case enums.OperationReturn.String():
		return cs.returnCommittee(ctx, req)
	case enums.OperationReSubmit.String():
		return cs.reSubmitCommittee(ctx, req)
	case enums.OperationClose.String():
		return cs.closeCommittee(ctx, req)
	case enums.OperationPause.String():
		return cs.pauseCommittee(ctx, req)
	case enums.OperationCancel.String():
		return cs.cancelCommittee(ctx, req)
	default:
		return cs.addCommittee(ctx, req)
	}
}

func (cs *committeeService) saveDraftCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	committee := performance.Committee{
		Chairperson: req.Chairperson,
		BaseProject: performance.BaseProject{
			Name:           req.Name,
			Description:    req.Description,
			StartDate:      req.StartDate,
			EndDate:        req.EndDate,
			Deliverables:   req.Deliverables,
			ReviewPeriodID: req.ReviewPeriodID,
			DepartmentID:   req.DepartmentID,
		},
	}
	committee.RecordStatus = enums.StatusDraft.String()
	committee.IsActive = true
	committee.CreatedBy = req.CreatedBy

	if err := cs.db.WithContext(ctx).Create(&committee).Error; err != nil {
		return resp, fmt.Errorf("saving draft committee: %w", err)
	}

	resp.ID = committee.CommitteeID
	resp.Message = "committee draft saved successfully"
	return resp, nil
}

func (cs *committeeService) addCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	// Validate review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := cs.db.WithContext(ctx).
		Where("period_id = ?", req.ReviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		return resp, fmt.Errorf("review period not found: %w", err)
	}

	// Check duplicate name in department
	var existing performance.Committee
	err := cs.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND department_id = ? AND record_status NOT IN ?",
			req.Name, req.DepartmentID, []string{enums.StatusCancelled.String()}).
		First(&existing).Error
	if err == nil {
		return resp, fmt.Errorf("committee name already exists in this department")
	}

	committee := performance.Committee{
		Chairperson: req.Chairperson,
		BaseProject: performance.BaseProject{
			Name:           req.Name,
			Description:    req.Description,
			StartDate:      req.StartDate,
			EndDate:        req.EndDate,
			Deliverables:   req.Deliverables,
			ReviewPeriodID: req.ReviewPeriodID,
			DepartmentID:   req.DepartmentID,
		},
	}
	committee.RecordStatus = enums.StatusPendingApproval.String()
	committee.IsActive = true
	committee.CreatedBy = req.CreatedBy

	if err := cs.db.WithContext(ctx).Create(&committee).Error; err != nil {
		return resp, fmt.Errorf("creating committee: %w", err)
	}

	resp.ID = committee.CommitteeID
	resp.Message = "committee created successfully"
	return resp, nil
}

func (cs *committeeService) updateCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	var committee performance.Committee
	if err := cs.db.WithContext(ctx).
		Where("committee_id = ?", req.CommitteeID).
		First(&committee).Error; err != nil {
		return resp, fmt.Errorf("committee not found: %w", err)
	}

	committee.Name = req.Name
	committee.Description = req.Description
	committee.StartDate = req.StartDate
	committee.EndDate = req.EndDate
	committee.Deliverables = req.Deliverables
	committee.Chairperson = req.Chairperson

	if err := cs.db.WithContext(ctx).Save(&committee).Error; err != nil {
		return resp, fmt.Errorf("updating committee: %w", err)
	}

	resp.ID = committee.CommitteeID
	resp.Message = "committee updated successfully"
	return resp, nil
}

func (cs *committeeService) approveCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusActive.String(),
			"is_active":     true,
			"is_approved":   true,
			"approved_by":   req.UpdatedBy,
			"date_approved": now,
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee approved successfully"
	return resp, nil
}

func (cs *committeeService) rejectCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status":    enums.StatusRejected.String(),
			"is_rejected":      true,
			"rejected_by":      req.UpdatedBy,
			"rejection_reason": req.RejectionReason,
			"date_rejected":    now,
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee rejected successfully"
	return resp, nil
}

func (cs *committeeService) returnCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusReturned.String(),
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee returned successfully"
	return resp, nil
}

func (cs *committeeService) reSubmitCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPendingApproval.String(),
			"is_rejected":   false,
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee re-submitted successfully"
	return resp, nil
}

func (cs *committeeService) closeCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusClosed.String(),
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee closed successfully"
	return resp, nil
}

func (cs *committeeService) pauseCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPaused.String(),
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee paused successfully"
	return resp, nil
}

func (cs *committeeService) cancelCommittee(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.CommitteeID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusCancelled.String(),
			"is_active":     false,
		})

	resp.ID = req.CommitteeID
	resp.Message = "committee cancelled successfully"
	return resp, nil
}

// =========================================================================
// GetCommittee -- retrieves a single committee with all associations.
// Mirrors .NET GetCommittee.
// =========================================================================

func (cs *committeeService) GetCommittee(ctx context.Context, committeeID string) (performance.CommitteeResponseVm, error) {
	resp := performance.CommitteeResponseVm{}
	resp.Message = "an error occurred"

	var committee performance.Committee
	err := cs.db.WithContext(ctx).
		Preload("CommitteeMembers").
		Preload("CommitteeMembers.PlannedObjective").
		Preload("CommitteeObjectives").
		Preload("CommitteeObjectives.Objective").
		Preload("CommitteeWorkProducts").
		Preload("CommitteeWorkProducts.WorkProduct").
		Preload("CommitteeAssignedWorkProducts").
		Preload("ReviewPeriod").
		Preload("Department").
		Where("committee_id = ?", committeeID).
		First(&committee).Error
	if err != nil {
		cs.log.Error().Err(err).Str("committeeID", committeeID).Msg("committee not found")
		resp.HasError = true
		resp.Message = "committee not found"
		return resp, fmt.Errorf("committee not found: %w", err)
	}

	data := cs.mapCommitteeToData(ctx, committee)
	resp.Committee = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCommittees -- retrieves all committees.
// Mirrors .NET GetCommittees().
func (cs *committeeService) GetCommittees(ctx context.Context) (performance.CommitteeListResponseVm, error) {
	resp := performance.CommitteeListResponseVm{}
	resp.Message = "an error occurred"

	var committees []performance.Committee
	err := cs.db.WithContext(ctx).
		Where("record_status NOT IN ?", []string{enums.StatusCancelled.String()}).
		Preload("CommitteeMembers").
		Preload("CommitteeObjectives").
		Preload("CommitteeObjectives.Objective").
		Preload("ReviewPeriod").
		Preload("Department").
		Find(&committees).Error
	if err != nil {
		cs.log.Error().Err(err).Msg("failed to get all committees")
		resp.HasError = true
		return resp, err
	}

	var vms []performance.CommitteeViewModel
	for _, c := range committees {
		vms = append(vms, cs.mapCommitteeToViewModel(ctx, c))
	}

	resp.Committees = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetCommitteesByChairperson -- retrieves committees chaired by a specific staff member.
// Mirrors .NET GetCommittees(chairpersonId).
func (cs *committeeService) GetCommitteesByChairperson(ctx context.Context, chairpersonID string) (performance.CommitteeListResponseVm, error) {
	resp := performance.CommitteeListResponseVm{}
	resp.Message = "an error occurred"

	var committees []performance.Committee
	err := cs.db.WithContext(ctx).
		Where("chairperson = ? AND record_status NOT IN ?", chairpersonID, []string{enums.StatusCancelled.String()}).
		Preload("CommitteeMembers").
		Preload("CommitteeObjectives").
		Preload("CommitteeObjectives.Objective").
		Preload("ReviewPeriod").
		Find(&committees).Error
	if err != nil {
		cs.log.Error().Err(err).Str("chairpersonID", chairpersonID).Msg("failed to get committees by chairperson")
		resp.HasError = true
		return resp, err
	}

	var vms []performance.CommitteeViewModel
	for _, c := range committees {
		vms = append(vms, cs.mapCommitteeToViewModel(ctx, c))
	}

	resp.Committees = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// CommitteeObjectiveSetup -- full committee objective lifecycle workflow.
// Mirrors .NET CommitteeObjectiveSetup.
// =========================================================================

func (cs *committeeService) CommitteeObjectiveSetup(ctx context.Context, req *performance.CommitteeObjectiveRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.RecordStatus {
	case enums.OperationAdd.String():
		// Validate committee exists
		var committee performance.Committee
		if err := cs.db.WithContext(ctx).
			Where("committee_id = ?", req.CommitteeID).First(&committee).Error; err != nil {
			return resp, fmt.Errorf("committee not found: %w", err)
		}

		// Check duplicate
		var existing performance.CommitteeObjective
		if err := cs.db.WithContext(ctx).
			Where("committee_id = ? AND objective_id = ? AND record_status != ?",
				req.CommitteeID, req.ObjectiveID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("objective already linked to this committee")
		}

		obj := performance.CommitteeObjective{
			ObjectiveID: req.ObjectiveID,
			CommitteeID: req.CommitteeID,
		}
		obj.RecordStatus = enums.StatusPendingApproval.String()
		obj.IsActive = true

		if err := cs.db.WithContext(ctx).Create(&obj).Error; err != nil {
			return resp, fmt.Errorf("adding committee objective: %w", err)
		}

		resp.ID = obj.CommitteeObjectiveID
		resp.Message = "committee objective added successfully"

	case enums.OperationApprove.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeObjective{}).
			Where("committee_objective_id = ?", req.CommitteeObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_active":     true,
			})
		resp.ID = req.CommitteeObjectiveID
		resp.Message = "committee objective approved"

	case enums.OperationReject.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeObjective{}).
			Where("committee_objective_id = ?", req.CommitteeObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusRejected.String(),
			})
		resp.ID = req.CommitteeObjectiveID
		resp.Message = "committee objective rejected"

	case enums.OperationCancel.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeObjective{}).
			Where("committee_objective_id = ?", req.CommitteeObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CommitteeObjectiveID
		resp.Message = "committee objective cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for committee objective")
	}

	return resp, nil
}

// GetCommitteeObjectives -- retrieves objectives for a committee.
// Mirrors .NET GetCommitteeObjectives.
func (cs *committeeService) GetCommitteeObjectives(ctx context.Context, committeeID string) (performance.CommitteeObjectiveListResponseVm, error) {
	resp := performance.CommitteeObjectiveListResponseVm{}
	resp.Message = "an error occurred"

	var objectives []performance.CommitteeObjective
	err := cs.db.WithContext(ctx).
		Where("committee_id = ? AND record_status != ?", committeeID, enums.StatusCancelled.String()).
		Preload("Objective").
		Find(&objectives).Error
	if err != nil {
		cs.log.Error().Err(err).Str("committeeID", committeeID).Msg("failed to get committee objectives")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CommitteeObjectiveData
	for _, obj := range objectives {
		d := performance.CommitteeObjectiveData{
			CommitteeObjectiveID: obj.CommitteeObjectiveID,
			ObjectiveID:          obj.ObjectiveID,
			CommitteeID:          obj.CommitteeID,
			RecordStatusName:     obj.RecordStatus,
		}
		if obj.Objective != nil {
			d.Objective = obj.Objective.Name
			d.Kpi = obj.Objective.Kpi
		}
		data = append(data, d)
	}

	resp.CommitteeObjectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// CommitteeMembersSetup -- full committee member lifecycle workflow.
// Mirrors .NET CommitteeMembersSetup.
// =========================================================================

func (cs *committeeService) CommitteeMembersSetup(ctx context.Context, req *performance.CommitteeMemberRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.Status {
	case enums.OperationAdd.String():
		// Validate committee
		var committee performance.Committee
		if err := cs.db.WithContext(ctx).
			Where("committee_id = ? AND record_status = ?", req.CommitteeID, enums.StatusActive.String()).
			First(&committee).Error; err != nil {
			return resp, fmt.Errorf("active committee not found: %w", err)
		}

		// Check duplicate
		var existing performance.CommitteeMember
		if err := cs.db.WithContext(ctx).
			Where("committee_id = ? AND staff_id = ? AND record_status != ?",
				req.CommitteeID, req.StaffID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("staff is already a member of this committee")
		}

		member := performance.CommitteeMember{
			StaffID:            req.StaffID,
			CommitteeID:        req.CommitteeID,
			PlannedObjectiveID: req.PlannedObjectiveID,
		}
		member.RecordStatus = enums.StatusActive.String()
		member.IsActive = true
		member.IsApproved = true

		if err := cs.db.WithContext(ctx).Create(&member).Error; err != nil {
			return resp, fmt.Errorf("adding committee member: %w", err)
		}

		resp.ID = member.CommitteeMemberID
		resp.Message = "committee member added successfully"

	case enums.OperationApprove.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeMember{}).
			Where("committee_member_id = ?", req.CommitteeMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_active":     true,
				"is_approved":   true,
			})
		resp.ID = req.CommitteeMemberID
		resp.Message = "committee member approved"

	case enums.OperationReject.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeMember{}).
			Where("committee_member_id = ?", req.CommitteeMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusRejected.String(),
				"is_rejected":   true,
			})
		resp.ID = req.CommitteeMemberID
		resp.Message = "committee member rejected"

	case enums.OperationCancel.String():
		cs.db.WithContext(ctx).Model(&performance.CommitteeMember{}).
			Where("committee_member_id = ?", req.CommitteeMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.CommitteeMemberID
		resp.Message = "committee member removed"

	default:
		return resp, fmt.Errorf("unsupported operation for committee member")
	}

	return resp, nil
}

// GetCommitteeMembers -- retrieves members of a committee.
// Mirrors .NET GetCommitteeMembers.
func (cs *committeeService) GetCommitteeMembers(ctx context.Context, committeeID string) (performance.CommitteeMemberListResponseVm, error) {
	resp := performance.CommitteeMemberListResponseVm{}
	resp.Message = "an error occurred"

	var members []performance.CommitteeMember
	err := cs.db.WithContext(ctx).
		Where("committee_id = ? AND record_status != ?", committeeID, enums.StatusCancelled.String()).
		Preload("PlannedObjective").
		Find(&members).Error
	if err != nil {
		cs.log.Error().Err(err).Str("committeeID", committeeID).Msg("failed to get committee members")
		resp.HasError = true
		return resp, err
	}

	var data []performance.CommitteeMemberData
	for _, m := range members {
		d := performance.CommitteeMemberData{
			CommitteeMemberID:  m.CommitteeMemberID,
			StaffID:            m.StaffID,
			CommitteeID:        m.CommitteeID,
			PlannedObjectiveID: m.PlannedObjectiveID,
		}

		// Enrich staff name
		if cs.parent.erpEmployeeSvc != nil {
			if detail, empErr := cs.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, m.StaffID); empErr == nil && detail != nil {
				if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
					d.StaffName = nameHolder.GetFullName()
				}
			}
		}

		if m.PlannedObjective != nil {
			d.ObjectiveName = m.PlannedObjective.ObjectiveID
		}

		data = append(data, d)
	}

	resp.CommitteeMembers = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetCommitteesAssigned -- retrieves committees assigned to a staff member.
// Mirrors .NET GetCommitteesAssigned.
// =========================================================================

func (cs *committeeService) GetCommitteesAssigned(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error) {
	resp := performance.CommitteeAssignedListResponseVm{}
	resp.Message = "an error occurred"

	// Get committees where staff is a member
	var members []performance.CommitteeMember
	cs.db.WithContext(ctx).
		Where("staff_id = ? AND record_status = ?", staffID, enums.StatusActive.String()).
		Preload("Committee").
		Preload("Committee.CommitteeObjectives").
		Preload("Committee.CommitteeObjectives.Objective").
		Preload("Committee.ReviewPeriod").
		Preload("Committee.Department").
		Find(&members)

	var committees []performance.CommitteeData
	for _, m := range members {
		if m.Committee != nil {
			committees = append(committees, cs.mapCommitteeToData(ctx, *m.Committee))
		}
	}

	resp.Committees = committees
	resp.TotalRecords = len(committees)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetStaffCommittees -- retrieves committees for a staff member (as chairperson or member).
// Mirrors .NET GetStaffCommittees.
func (cs *committeeService) GetStaffCommittees(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error) {
	resp := performance.CommitteeAssignedListResponseVm{}
	resp.Message = "an error occurred"

	// Committees as chairperson
	var chairedCommittees []performance.Committee
	cs.db.WithContext(ctx).
		Where("chairperson = ? AND record_status NOT IN ?", staffID, []string{enums.StatusCancelled.String()}).
		Preload("CommitteeObjectives").
		Preload("CommitteeObjectives.Objective").
		Preload("ReviewPeriod").
		Preload("Department").
		Find(&chairedCommittees)

	// Committees as member
	var memberCommittees []performance.CommitteeMember
	cs.db.WithContext(ctx).
		Where("staff_id = ? AND record_status = ?", staffID, enums.StatusActive.String()).
		Preload("Committee").
		Preload("Committee.CommitteeObjectives").
		Preload("Committee.CommitteeObjectives.Objective").
		Preload("Committee.ReviewPeriod").
		Preload("Committee.Department").
		Find(&memberCommittees)

	seen := make(map[string]bool)
	var committees []performance.CommitteeData

	for _, c := range chairedCommittees {
		if !seen[c.CommitteeID] {
			seen[c.CommitteeID] = true
			committees = append(committees, cs.mapCommitteeToData(ctx, c))
		}
	}

	for _, m := range memberCommittees {
		if m.Committee != nil && !seen[m.Committee.CommitteeID] {
			seen[m.Committee.CommitteeID] = true
			committees = append(committees, cs.mapCommitteeToData(ctx, *m.Committee))
		}
	}

	resp.Committees = committees
	resp.TotalRecords = len(committees)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetCommitteeWorkProductStaffList -- retrieves staff who have work products
// for a committee. Mirrors .NET GetCommitteeWorkProductStaffList.
// =========================================================================

func (cs *committeeService) GetCommitteeWorkProductStaffList(ctx context.Context, committeeID string) ([]string, error) {
	var committeeWPs []performance.CommitteeWorkProduct
	cs.db.WithContext(ctx).
		Where("committee_id = ?", committeeID).
		Preload("WorkProduct").
		Find(&committeeWPs)

	seen := make(map[string]bool)
	var staffIDs []string
	for _, cwp := range committeeWPs {
		if cwp.WorkProduct != nil && !seen[cwp.WorkProduct.StaffID] {
			seen[cwp.WorkProduct.StaffID] = true
			staffIDs = append(staffIDs, cwp.WorkProduct.StaffID)
		}
	}

	return staffIDs, nil
}

// =========================================================================
// ChangeCommitteeChairperson -- changes the committee chairperson.
// Mirrors .NET ChangeAdhocAssignmentLead for committee type.
// =========================================================================

func (cs *committeeService) ChangeCommitteeChairperson(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error {
	result := cs.db.WithContext(ctx).Model(&performance.Committee{}).
		Where("committee_id = ?", req.ReferenceID).
		Update("chairperson", req.StaffID)

	if result.Error != nil {
		return fmt.Errorf("changing committee chairperson: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("committee not found")
	}

	cs.log.Info().
		Str("committeeID", req.ReferenceID).
		Str("newChairperson", req.StaffID).
		Msg("committee chairperson changed")

	return nil
}

// =========================================================================
// Internal helpers
// =========================================================================

func (cs *committeeService) mapCommitteeToData(ctx context.Context, c performance.Committee) performance.CommitteeData {
	data := performance.CommitteeData{
		BaseProjectData: performance.BaseProjectData{
			Name:           c.Name,
			Description:    c.Description,
			StartDate:      c.StartDate,
			EndDate:        c.EndDate,
			Deliverables:   c.Deliverables,
			ReviewPeriodID: c.ReviewPeriodID,
			DepartmentID:   c.DepartmentID,
		},
		CommitteeID: c.CommitteeID,
		Chairperson: c.Chairperson,
	}
	data.RecordStatus = c.RecordStatus
	data.IsActive = c.IsActive

	// Enrich chairperson name
	if cs.parent.erpEmployeeSvc != nil {
		if detail, err := cs.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, c.Chairperson); err == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				data.ChairpersonName = nameHolder.GetFullName()
			}
		}
	}

	if c.Department != nil {
		data.DepartmentName = c.Department.DepartmentName
	}

	// Map objectives
	for _, obj := range c.CommitteeObjectives {
		objData := performance.CommitteeObjectiveData{
			CommitteeObjectiveID: obj.CommitteeObjectiveID,
			ObjectiveID:          obj.ObjectiveID,
			CommitteeID:          obj.CommitteeID,
			RecordStatusName:     obj.RecordStatus,
		}
		if obj.Objective != nil {
			objData.Objective = obj.Objective.Name
			objData.Kpi = obj.Objective.Kpi
		}
		data.CommitteeObjectives = append(data.CommitteeObjectives, objData)
	}

	// Map members
	for _, m := range c.CommitteeMembers {
		memberData := performance.CommitteeMemberData{
			CommitteeMemberID:  m.CommitteeMemberID,
			StaffID:            m.StaffID,
			CommitteeID:        m.CommitteeID,
			PlannedObjectiveID: m.PlannedObjectiveID,
		}
		data.CommitteeMembers = append(data.CommitteeMembers, memberData)
	}

	return data
}

func (cs *committeeService) mapCommitteeToViewModel(ctx context.Context, c performance.Committee) performance.CommitteeViewModel {
	vm := performance.CommitteeViewModel{
		CommitteeID:    c.CommitteeID,
		Chairperson:    c.Chairperson,
		Name:           c.Name,
		Description:    c.Description,
		StartDate:      c.StartDate,
		EndDate:        c.EndDate,
		Deliverables:   c.Deliverables,
		ReviewPeriodID: c.ReviewPeriodID,
		DepartmentID:   c.DepartmentID,
	}
	vm.RecordStatus = c.RecordStatus
	vm.IsActive = c.IsActive

	// Enrich chairperson name
	if cs.parent.erpEmployeeSvc != nil {
		if detail, err := cs.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, c.Chairperson); err == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				vm.ChairpersonName = nameHolder.GetFullName()
			}
		}
	}

	for _, obj := range c.CommitteeObjectives {
		objData := performance.CommitteeObjectiveData{
			CommitteeObjectiveID: obj.CommitteeObjectiveID,
			ObjectiveID:          obj.ObjectiveID,
			CommitteeID:          obj.CommitteeID,
		}
		if obj.Objective != nil {
			objData.Objective = obj.Objective.Name
			objData.Kpi = obj.Objective.Kpi
		}
		vm.CommitteeObjectives = append(vm.CommitteeObjectives, objData)
	}

	return vm
}
