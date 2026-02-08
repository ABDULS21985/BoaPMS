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
// projectService handles full project lifecycle management.
//
// Mirrors .NET methods:
//   - ProjectSetup (full workflow: Draft/Add/Update/Approve/Reject/Return/
//     ReSubmit/Close/Pause/Cancel)
//   - GetProject / GetProjects / GetProjectsByManager
//   - ProjectObjectiveSetup (full workflow)
//   - GetProjectObjectives
//   - ProjectMembersSetup (full workflow)
//   - GetProjectMembers
//   - GetProjectWorkProductStaffList
//   - GetProjectsAssigned / GetStaffProjects
//   - ChangeAdhocAssignmentLead
//   - ValidateStaffEligibilityForAdhoc
// ---------------------------------------------------------------------------

type projectService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	parent *performanceManagementService

	projectRepo         *repository.PMSRepository[performance.Project]
	projectObjRepo      *repository.PMSRepository[performance.ProjectObjective]
	projectMemberRepo   *repository.PMSRepository[performance.ProjectMember]
	projectWPRepo       *repository.PMSRepository[performance.ProjectWorkProduct]
	projAssignedWPRepo  *repository.PMSRepository[performance.ProjectAssignedWorkProduct]
	reviewPeriodRepo    *repository.PMSRepository[performance.PerformanceReviewPeriod]
	plannedObjRepo      *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
}

func newProjectService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *projectService {
	return &projectService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "project").Logger(),
		parent: parent,

		projectRepo:        repository.NewPMSRepository[performance.Project](db),
		projectObjRepo:     repository.NewPMSRepository[performance.ProjectObjective](db),
		projectMemberRepo:  repository.NewPMSRepository[performance.ProjectMember](db),
		projectWPRepo:      repository.NewPMSRepository[performance.ProjectWorkProduct](db),
		projAssignedWPRepo: repository.NewPMSRepository[performance.ProjectAssignedWorkProduct](db),
		reviewPeriodRepo:   repository.NewPMSRepository[performance.PerformanceReviewPeriod](db),
		plannedObjRepo:     repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
	}
}

// =========================================================================
// ProjectSetup – full project lifecycle workflow.
// Mirrors .NET ProjectSetup with OperationType switch.
// =========================================================================

func (ps *projectService) ProjectSetup(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}
	resp.Message = "an error occurred"

	switch req.Status {
	case enums.OperationDraft.String():
		return ps.saveDraftProject(ctx, req)
	case enums.OperationAdd.String():
		return ps.addProject(ctx, req)
	case enums.OperationUpdate.String():
		return ps.updateProject(ctx, req)
	case enums.OperationApprove.String():
		return ps.approveProject(ctx, req)
	case enums.OperationReject.String():
		return ps.rejectProject(ctx, req)
	case enums.OperationReturn.String():
		return ps.returnProject(ctx, req)
	case enums.OperationReSubmit.String():
		return ps.reSubmitProject(ctx, req)
	case enums.OperationClose.String():
		return ps.closeProject(ctx, req)
	case enums.OperationPause.String():
		return ps.pauseProject(ctx, req)
	case enums.OperationCancel.String():
		return ps.cancelProject(ctx, req)
	default:
		return ps.addProject(ctx, req)
	}
}

func (ps *projectService) saveDraftProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	project := performance.Project{
		ProjectManager: req.ProjectManager,
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
	project.RecordStatus = enums.StatusDraft.String()
	project.IsActive = true
	project.CreatedBy = req.CreatedBy

	if err := ps.db.WithContext(ctx).Create(&project).Error; err != nil {
		return resp, fmt.Errorf("saving draft project: %w", err)
	}

	resp.ID = project.ProjectID
	resp.Message = "project draft saved successfully"
	return resp, nil
}

func (ps *projectService) addProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	// Validate review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := ps.db.WithContext(ctx).
		Where("period_id = ?", req.ReviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		return resp, fmt.Errorf("review period not found: %w", err)
	}

	// Check duplicate name in department
	var existing performance.Project
	err := ps.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND department_id = ? AND record_status NOT IN ?",
			req.Name, req.DepartmentID, []string{enums.StatusCancelled.String()}).
		First(&existing).Error
	if err == nil {
		return resp, fmt.Errorf("project name already exists in this department")
	}

	project := performance.Project{
		ProjectManager: req.ProjectManager,
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
	project.RecordStatus = enums.StatusPendingApproval.String()
	project.IsActive = true
	project.CreatedBy = req.CreatedBy

	if err := ps.db.WithContext(ctx).Create(&project).Error; err != nil {
		return resp, fmt.Errorf("creating project: %w", err)
	}

	resp.ID = project.ProjectID
	resp.Message = "project created successfully"
	return resp, nil
}

func (ps *projectService) updateProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	var project performance.Project
	if err := ps.db.WithContext(ctx).
		Where("project_id = ?", req.ProjectID).
		First(&project).Error; err != nil {
		return resp, fmt.Errorf("project not found: %w", err)
	}

	project.Name = req.Name
	project.Description = req.Description
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate
	project.Deliverables = req.Deliverables
	project.ProjectManager = req.ProjectManager

	if err := ps.db.WithContext(ctx).Save(&project).Error; err != nil {
		return resp, fmt.Errorf("updating project: %w", err)
	}

	resp.ID = project.ProjectID
	resp.Message = "project updated successfully"
	return resp, nil
}

func (ps *projectService) approveProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusActive.String(),
			"is_active":     true,
			"is_approved":   true,
			"approved_by":   req.UpdatedBy,
			"date_approved":  now,
		})

	resp.ID = req.ProjectID
	resp.Message = "project approved successfully"
	return resp, nil
}

func (ps *projectService) rejectProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	now := time.Now().UTC()
	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status":    enums.StatusRejected.String(),
			"is_rejected":      true,
			"rejected_by":      req.UpdatedBy,
			"rejection_reason": req.RejectionReason,
			"date_rejected":    now,
		})

	resp.ID = req.ProjectID
	resp.Message = "project rejected successfully"
	return resp, nil
}

func (ps *projectService) returnProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusReturned.String(),
		})

	resp.ID = req.ProjectID
	resp.Message = "project returned successfully"
	return resp, nil
}

func (ps *projectService) reSubmitProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPendingApproval.String(),
			"is_rejected":   false,
		})

	resp.ID = req.ProjectID
	resp.Message = "project re-submitted successfully"
	return resp, nil
}

func (ps *projectService) closeProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusClosed.String(),
		})

	resp.ID = req.ProjectID
	resp.Message = "project closed successfully"
	return resp, nil
}

func (ps *projectService) pauseProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusPaused.String(),
		})

	resp.ID = req.ProjectID
	resp.Message = "project paused successfully"
	return resp, nil
}

func (ps *projectService) cancelProject(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ProjectID).
		Updates(map[string]interface{}{
			"record_status": enums.StatusCancelled.String(),
			"is_active":     false,
		})

	resp.ID = req.ProjectID
	resp.Message = "project cancelled successfully"
	return resp, nil
}

// =========================================================================
// GetProject – retrieves a single project with all associations.
// Mirrors .NET GetProject.
// =========================================================================

func (ps *projectService) GetProject(ctx context.Context, projectID string) (performance.ProjectResponseVm, error) {
	resp := performance.ProjectResponseVm{}
	resp.Message = "an error occurred"

	var project performance.Project
	err := ps.db.WithContext(ctx).
		Preload("ProjectMembers").
		Preload("ProjectMembers.PlannedObjective").
		Preload("ProjectObjectives").
		Preload("ProjectObjectives.Objective").
		Preload("ProjectWorkProducts").
		Preload("ProjectWorkProducts.WorkProduct").
		Preload("ProjectAssignedWorkProducts").
		Preload("ReviewPeriod").
		Preload("Department").
		Where("project_id = ?", projectID).
		First(&project).Error
	if err != nil {
		ps.log.Error().Err(err).Str("projectID", projectID).Msg("project not found")
		resp.HasError = true
		resp.Message = "project not found"
		return resp, fmt.Errorf("project not found: %w", err)
	}

	data := ps.mapProjectToData(ctx, project)
	resp.Project = &data
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetProjects – retrieves all projects.
// Mirrors .NET GetProjects().
func (ps *projectService) GetProjects(ctx context.Context) (performance.ProjectListResponseVm, error) {
	resp := performance.ProjectListResponseVm{}
	resp.Message = "an error occurred"

	var projects []performance.Project
	err := ps.db.WithContext(ctx).
		Where("record_status NOT IN ?", []string{enums.StatusCancelled.String()}).
		Preload("ProjectMembers").
		Preload("ProjectObjectives").
		Preload("ProjectObjectives.Objective").
		Preload("ReviewPeriod").
		Preload("Department").
		Find(&projects).Error
	if err != nil {
		ps.log.Error().Err(err).Msg("failed to get all projects")
		resp.HasError = true
		return resp, err
	}

	var vms []performance.ProjectViewModel
	for _, p := range projects {
		vms = append(vms, ps.mapProjectToViewModel(ctx, p))
	}

	resp.Projects = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetProjectsByManager – retrieves projects managed by a specific staff member.
// Mirrors .NET GetProjects(managerId).
func (ps *projectService) GetProjectsByManager(ctx context.Context, managerID string) (performance.ProjectListResponseVm, error) {
	resp := performance.ProjectListResponseVm{}
	resp.Message = "an error occurred"

	var projects []performance.Project
	err := ps.db.WithContext(ctx).
		Where("project_manager = ? AND record_status NOT IN ?", managerID, []string{enums.StatusCancelled.String()}).
		Preload("ProjectMembers").
		Preload("ProjectObjectives").
		Preload("ProjectObjectives.Objective").
		Preload("ReviewPeriod").
		Find(&projects).Error
	if err != nil {
		ps.log.Error().Err(err).Str("managerID", managerID).Msg("failed to get projects by manager")
		resp.HasError = true
		return resp, err
	}

	var vms []performance.ProjectViewModel
	for _, p := range projects {
		vms = append(vms, ps.mapProjectToViewModel(ctx, p))
	}

	resp.Projects = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// ProjectObjectiveSetup – full project objective lifecycle workflow.
// Mirrors .NET ProjectObjectiveSetup.
// =========================================================================

func (ps *projectService) ProjectObjectiveSetup(ctx context.Context, req *performance.ProjectObjectiveRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.RecordStatus {
	case enums.OperationAdd.String():
		// Validate project and objective exist
		var project performance.Project
		if err := ps.db.WithContext(ctx).
			Where("project_id = ?", req.ProjectID).First(&project).Error; err != nil {
			return resp, fmt.Errorf("project not found: %w", err)
		}

		// Check duplicate
		var existing performance.ProjectObjective
		if err := ps.db.WithContext(ctx).
			Where("project_id = ? AND objective_id = ? AND record_status != ?",
				req.ProjectID, req.ObjectiveID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("objective already linked to this project")
		}

		obj := performance.ProjectObjective{
			ObjectiveID: req.ObjectiveID,
			ProjectID:   req.ProjectID,
		}
		obj.RecordStatus = enums.StatusPendingApproval.String()
		obj.IsActive = true

		if err := ps.db.WithContext(ctx).Create(&obj).Error; err != nil {
			return resp, fmt.Errorf("adding project objective: %w", err)
		}

		resp.ID = obj.ProjectObjectiveID
		resp.Message = "project objective added successfully"

	case enums.OperationApprove.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectObjective{}).
			Where("project_objective_id = ?", req.ProjectObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_active":     true,
			})
		resp.ID = req.ProjectObjectiveID
		resp.Message = "project objective approved"

	case enums.OperationReject.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectObjective{}).
			Where("project_objective_id = ?", req.ProjectObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusRejected.String(),
			})
		resp.ID = req.ProjectObjectiveID
		resp.Message = "project objective rejected"

	case enums.OperationCancel.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectObjective{}).
			Where("project_objective_id = ?", req.ProjectObjectiveID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.ProjectObjectiveID
		resp.Message = "project objective cancelled"

	default:
		return resp, fmt.Errorf("unsupported operation for project objective")
	}

	return resp, nil
}

// GetProjectObjectives – retrieves objectives for a project.
// Mirrors .NET GetProjectObjectives.
func (ps *projectService) GetProjectObjectives(ctx context.Context, projectID string) (performance.ProjectObjectiveListResponseVm, error) {
	resp := performance.ProjectObjectiveListResponseVm{}
	resp.Message = "an error occurred"

	var objectives []performance.ProjectObjective
	err := ps.db.WithContext(ctx).
		Where("project_id = ? AND record_status != ?", projectID, enums.StatusCancelled.String()).
		Preload("Objective").
		Find(&objectives).Error
	if err != nil {
		ps.log.Error().Err(err).Str("projectID", projectID).Msg("failed to get project objectives")
		resp.HasError = true
		return resp, err
	}

	var data []performance.ProjectObjectiveData
	for _, obj := range objectives {
		d := performance.ProjectObjectiveData{
			ProjectObjectiveID: obj.ProjectObjectiveID,
			ObjectiveID:        obj.ObjectiveID,
			ProjectID:          obj.ProjectID,
			RecordStatusName:   obj.RecordStatus,
		}
		if obj.Objective != nil {
			d.Objective = obj.Objective.Name
			d.Kpi = obj.Objective.Kpi
		}
		data = append(data, d)
	}

	resp.ProjectObjectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// ProjectMembersSetup – full project member lifecycle workflow.
// Mirrors .NET ProjectMembersSetup.
// =========================================================================

func (ps *projectService) ProjectMembersSetup(ctx context.Context, req *performance.ProjectMemberRequestModel) (performance.ResponseVm, error) {
	resp := performance.ResponseVm{}

	switch req.Status {
	case enums.OperationAdd.String():
		// Validate project
		var project performance.Project
		if err := ps.db.WithContext(ctx).
			Where("project_id = ? AND record_status = ?", req.ProjectID, enums.StatusActive.String()).
			First(&project).Error; err != nil {
			return resp, fmt.Errorf("active project not found: %w", err)
		}

		// Check duplicate
		var existing performance.ProjectMember
		if err := ps.db.WithContext(ctx).
			Where("project_id = ? AND staff_id = ? AND record_status != ?",
				req.ProjectID, req.StaffID, enums.StatusCancelled.String()).
			First(&existing).Error; err == nil {
			return resp, fmt.Errorf("staff is already a member of this project")
		}

		member := performance.ProjectMember{
			StaffID:            req.StaffID,
			ProjectID:          req.ProjectID,
			PlannedObjectiveID: req.PlannedObjectiveID,
		}
		member.RecordStatus = enums.StatusActive.String()
		member.IsActive = true
		member.IsApproved = true

		if err := ps.db.WithContext(ctx).Create(&member).Error; err != nil {
			return resp, fmt.Errorf("adding project member: %w", err)
		}

		resp.ID = member.ProjectMemberID
		resp.Message = "project member added successfully"

	case enums.OperationApprove.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectMember{}).
			Where("project_member_id = ?", req.ProjectMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusActive.String(),
				"is_active":     true,
				"is_approved":   true,
			})
		resp.ID = req.ProjectMemberID
		resp.Message = "project member approved"

	case enums.OperationReject.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectMember{}).
			Where("project_member_id = ?", req.ProjectMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusRejected.String(),
				"is_rejected":   true,
			})
		resp.ID = req.ProjectMemberID
		resp.Message = "project member rejected"

	case enums.OperationCancel.String():
		ps.db.WithContext(ctx).Model(&performance.ProjectMember{}).
			Where("project_member_id = ?", req.ProjectMemberID).
			Updates(map[string]interface{}{
				"record_status": enums.StatusCancelled.String(),
				"is_active":     false,
			})
		resp.ID = req.ProjectMemberID
		resp.Message = "project member removed"

	default:
		return resp, fmt.Errorf("unsupported operation for project member")
	}

	return resp, nil
}

// GetProjectMembers – retrieves members of a project.
// Mirrors .NET GetProjectMembers.
func (ps *projectService) GetProjectMembers(ctx context.Context, projectID string) (performance.ProjectMemberListResponseVm, error) {
	resp := performance.ProjectMemberListResponseVm{}
	resp.Message = "an error occurred"

	var members []performance.ProjectMember
	err := ps.db.WithContext(ctx).
		Where("project_id = ? AND record_status != ?", projectID, enums.StatusCancelled.String()).
		Preload("PlannedObjective").
		Find(&members).Error
	if err != nil {
		ps.log.Error().Err(err).Str("projectID", projectID).Msg("failed to get project members")
		resp.HasError = true
		return resp, err
	}

	var data []performance.ProjectMemberData
	for _, m := range members {
		d := performance.ProjectMemberData{
			ProjectMemberID:    m.ProjectMemberID,
			StaffID:            m.StaffID,
			ProjectID:          m.ProjectID,
			PlannedObjectiveID: m.PlannedObjectiveID,
		}

		// Enrich staff name
		if ps.parent.erpEmployeeSvc != nil {
			if detail, empErr := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, m.StaffID); empErr == nil && detail != nil {
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

	resp.ProjectMembers = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetProjectsAssigned – retrieves projects assigned to a staff member.
// Mirrors .NET GetProjectsAssigned.
// =========================================================================

func (ps *projectService) GetProjectsAssigned(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error) {
	resp := performance.ProjectAssignedListResponseVm{}
	resp.Message = "an error occurred"

	// Get projects where staff is a member
	var members []performance.ProjectMember
	ps.db.WithContext(ctx).
		Where("staff_id = ? AND record_status = ?", staffID, enums.StatusActive.String()).
		Preload("Project").
		Preload("Project.ProjectObjectives").
		Preload("Project.ProjectObjectives.Objective").
		Preload("Project.ReviewPeriod").
		Preload("Project.Department").
		Find(&members)

	var projects []performance.ProjectData
	for _, m := range members {
		if m.Project != nil {
			projects = append(projects, ps.mapProjectToData(ctx, *m.Project))
		}
	}

	resp.Projects = projects
	resp.TotalRecords = len(projects)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// GetStaffProjects – retrieves projects for a staff member (as manager or member).
// Mirrors .NET GetStaffProjects.
func (ps *projectService) GetStaffProjects(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error) {
	resp := performance.ProjectAssignedListResponseVm{}
	resp.Message = "an error occurred"

	// Projects as manager
	var managedProjects []performance.Project
	ps.db.WithContext(ctx).
		Where("project_manager = ? AND record_status NOT IN ?", staffID, []string{enums.StatusCancelled.String()}).
		Preload("ProjectObjectives").
		Preload("ProjectObjectives.Objective").
		Preload("ReviewPeriod").
		Preload("Department").
		Find(&managedProjects)

	// Projects as member
	var memberProjects []performance.ProjectMember
	ps.db.WithContext(ctx).
		Where("staff_id = ? AND record_status = ?", staffID, enums.StatusActive.String()).
		Preload("Project").
		Preload("Project.ProjectObjectives").
		Preload("Project.ProjectObjectives.Objective").
		Preload("Project.ReviewPeriod").
		Preload("Project.Department").
		Find(&memberProjects)

	seen := make(map[string]bool)
	var projects []performance.ProjectData

	for _, p := range managedProjects {
		if !seen[p.ProjectID] {
			seen[p.ProjectID] = true
			projects = append(projects, ps.mapProjectToData(ctx, p))
		}
	}

	for _, m := range memberProjects {
		if m.Project != nil && !seen[m.Project.ProjectID] {
			seen[m.Project.ProjectID] = true
			projects = append(projects, ps.mapProjectToData(ctx, *m.Project))
		}
	}

	resp.Projects = projects
	resp.TotalRecords = len(projects)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// GetProjectWorkProductStaffList – retrieves staff who have work products
// for a project. Mirrors .NET GetProjectWorkProductStaffList.
// =========================================================================

func (ps *projectService) GetProjectWorkProductStaffList(ctx context.Context, projectID string) ([]string, error) {
	var projectWPs []performance.ProjectWorkProduct
	ps.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("WorkProduct").
		Find(&projectWPs)

	seen := make(map[string]bool)
	var staffIDs []string
	for _, pwp := range projectWPs {
		if pwp.WorkProduct != nil && !seen[pwp.WorkProduct.StaffID] {
			seen[pwp.WorkProduct.StaffID] = true
			staffIDs = append(staffIDs, pwp.WorkProduct.StaffID)
		}
	}

	return staffIDs, nil
}

// =========================================================================
// ChangeAdhocAssignmentLead – changes the project manager.
// Mirrors .NET ChangeAdhocAssignmentLead.
// =========================================================================

func (ps *projectService) ChangeProjectLead(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error {
	result := ps.db.WithContext(ctx).Model(&performance.Project{}).
		Where("project_id = ?", req.ReferenceID).
		Update("project_manager", req.StaffID)

	if result.Error != nil {
		return fmt.Errorf("changing project lead: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	ps.log.Info().
		Str("projectID", req.ReferenceID).
		Str("newLead", req.StaffID).
		Msg("project lead changed")

	return nil
}

// =========================================================================
// ValidateStaffEligibilityForAdhoc – validates that a staff member is
// eligible to be assigned to a project/committee.
// Mirrors .NET ValidateStaffEligibilityForAdhoc.
// =========================================================================

func (ps *projectService) ValidateStaffEligibilityForAdhoc(ctx context.Context, staffID, reviewPeriodID string) (performance.AdhocStaffResponseVm, error) {
	resp := performance.AdhocStaffResponseVm{}
	resp.Message = "an error occurred"

	// Validate staff exists
	if ps.parent.erpEmployeeSvc != nil {
		empDetail, empErr := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
		if empErr != nil || empDetail == nil {
			resp.HasError = true
			resp.Message = "staff record not found"
			return resp, fmt.Errorf("staff not found: %s", staffID)
		}
		resp.Employee = empDetail
	}

	// Check if the staff has an active planned objective in this review period
	var plannedObj performance.ReviewPeriodIndividualPlannedObjective
	err := ps.db.WithContext(ctx).
		Where("staff_id = ? AND review_period_id = ? AND record_status = ?",
			staffID, reviewPeriodID, enums.StatusActive.String()).
		First(&plannedObj).Error
	if err != nil {
		resp.HasError = true
		resp.Message = "staff does not have an active planned objective in this review period"
		return resp, fmt.Errorf("no active planned objective found for staff")
	}

	resp.PlannedObjectiveID = plannedObj.PlannedObjectiveID
	resp.Message = "staff is eligible"
	return resp, nil
}

// =========================================================================
// Internal helpers
// =========================================================================

func (ps *projectService) mapProjectToData(ctx context.Context, p performance.Project) performance.ProjectData {
	data := performance.ProjectData{
		BaseProjectData: performance.BaseProjectData{
			Name:           p.Name,
			Description:    p.Description,
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			Deliverables:   p.Deliverables,
			ReviewPeriodID: p.ReviewPeriodID,
			DepartmentID:   p.DepartmentID,
		},
		ProjectID:      p.ProjectID,
		ProjectManager: p.ProjectManager,
	}
	data.RecordStatus = p.RecordStatus
	data.IsActive = p.IsActive

	// Enrich manager name
	if ps.parent.erpEmployeeSvc != nil {
		if detail, err := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, p.ProjectManager); err == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				data.ProjectManagerName = nameHolder.GetFullName()
			}
		}
	}

	if p.Department != nil {
		data.DepartmentName = p.Department.DepartmentName
	}

	// Map objectives
	for _, obj := range p.ProjectObjectives {
		objData := performance.ProjectObjectiveData{
			ProjectObjectiveID: obj.ProjectObjectiveID,
			ObjectiveID:        obj.ObjectiveID,
			ProjectID:          obj.ProjectID,
			RecordStatusName:   obj.RecordStatus,
		}
		if obj.Objective != nil {
			objData.Objective = obj.Objective.Name
			objData.Kpi = obj.Objective.Kpi
		}
		data.ProjectObjectives = append(data.ProjectObjectives, objData)
	}

	// Map members
	for _, m := range p.ProjectMembers {
		memberData := performance.ProjectMemberData{
			ProjectMemberID:    m.ProjectMemberID,
			StaffID:            m.StaffID,
			ProjectID:          m.ProjectID,
			PlannedObjectiveID: m.PlannedObjectiveID,
		}
		data.ProjectMembers = append(data.ProjectMembers, memberData)
	}

	return data
}

func (ps *projectService) mapProjectToViewModel(ctx context.Context, p performance.Project) performance.ProjectViewModel {
	vm := performance.ProjectViewModel{
		ProjectID:      p.ProjectID,
		ProjectManager: p.ProjectManager,
		Name:           p.Name,
		Description:    p.Description,
		StartDate:      p.StartDate,
		EndDate:        p.EndDate,
		Deliverables:   p.Deliverables,
		ReviewPeriodID: p.ReviewPeriodID,
		DepartmentID:   p.DepartmentID,
	}
	vm.RecordStatus = p.RecordStatus
	vm.IsActive = p.IsActive

	// Enrich manager name
	if ps.parent.erpEmployeeSvc != nil {
		if detail, err := ps.parent.erpEmployeeSvc.GetEmployeeDetail(ctx, p.ProjectManager); err == nil && detail != nil {
			if nameHolder, ok := detail.(interface{ GetFullName() string }); ok {
				vm.ProjectManagerName = nameHolder.GetFullName()
			}
		}
	}

	for _, obj := range p.ProjectObjectives {
		objData := performance.ProjectObjectiveData{
			ProjectObjectiveID: obj.ProjectObjectiveID,
			ObjectiveID:        obj.ObjectiveID,
			ProjectID:          obj.ProjectID,
		}
		if obj.Objective != nil {
			objData.Objective = obj.Objective.Name
			objData.Kpi = obj.Objective.Kpi
		}
		vm.ProjectObjectives = append(vm.ProjectObjectives, objData)
	}

	return vm
}

