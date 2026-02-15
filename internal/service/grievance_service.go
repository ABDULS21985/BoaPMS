package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// grievanceManagementService implements GrievanceManagementService.
// Mirrors the .NET GrievanceManagementService with IPMSRepo<Grievance> and
// IPMSRepo<GrievanceResolution>.
type grievanceManagementService struct {
	grievanceRepo    *repository.PMSRepository[performance.Grievance]
	resolutionRepo   *repository.PMSRepository[performance.GrievanceResolution]
	db               *gorm.DB
	erpDataDB        *sqlx.DB
	erpEmployeeSvc   ErpEmployeeService
	emailSvc         EmailService
	globalSettingSvc GlobalSettingService
	userContextSvc   UserContextService
	reviewPeriodSvc  ReviewPeriodService
	cfg              *config.Config
	log              zerolog.Logger
}

func newGrievanceManagementService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	erpSvc ErpEmployeeService,
	gsSvc GlobalSettingService,
	ucSvc UserContextService,
	emailSvc EmailService,
	rpSvc ReviewPeriodService,
) GrievanceManagementService {
	return &grievanceManagementService{
		grievanceRepo:    repository.NewPMSRepository[performance.Grievance](repos.GormDB),
		resolutionRepo:   repository.NewPMSRepository[performance.GrievanceResolution](repos.GormDB),
		db:               repos.GormDB,
		erpDataDB:        repos.ErpSQL,
		erpEmployeeSvc:   erpSvc,
		emailSvc:         emailSvc,
		globalSettingSvc: gsSvc,
		userContextSvc:   ucSvc,
		reviewPeriodSvc:  rpSvc,
		cfg:              cfg,
		log:              log.With().Str("service", "grievance").Logger(),
	}
}

// ---------------------------------------------------------------------------
// RaiseNewGrievance creates a new grievance record.
// Maps to .NET GrievanceSetup with OperationTypes.Add.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) RaiseNewGrievance(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.GrievanceVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *GrievanceVm")
	}

	s.log.Info().Str("complainant", vm.ComplainantStaffID).Msg("raising new grievance")

	// Generate a unique grievance ID from the sequence table.
	grievanceID, err := s.generateCode(ctx, enums.SeqGrievance, 8)
	if err != nil {
		return nil, fmt.Errorf("generating grievance ID: %w", err)
	}

	// Determine respondent (complainant's supervisor) and mediator (respondent's supervisor).
	complainantData, err := s.getEmployeeData(ctx, vm.ComplainantStaffID)
	if err != nil {
		return nil, fmt.Errorf("getting complainant employee data: %w", err)
	}
	respondentStaffID := complainantData.SupervisorID

	respondentData, err := s.getEmployeeData(ctx, respondentStaffID)
	if err != nil {
		return nil, fmt.Errorf("getting respondent employee data: %w", err)
	}
	mediatorStaffID := respondentData.SupervisorID

	grievance := performance.Grievance{
		GrievanceID:               grievanceID,
		GrievanceType:             enums.GrievanceType(vm.GrievanceType),
		ReviewPeriodID:            vm.ReviewPeriodID,
		SubjectID:                 vm.SubjectID,
		Subject:                   vm.Subject,
		Description:               vm.Description,
		ComplainantStaffID:        vm.ComplainantStaffID,
		ComplainantEvidenceUpload: vm.ComplainantEvidenceUpload,
		RespondentStaffID:         respondentStaffID,
		CurrentMediatorStaffID:    mediatorStaffID,
		CurrentResolutionLevel:    enums.ResolutionLevelSBU,
	}
	grievance.RecordStatus = enums.StatusAwaitingRespondentComment.String()

	if err := s.grievanceRepo.InsertAndSave(ctx, &grievance); err != nil {
		return nil, fmt.Errorf("saving new grievance: %w", err)
	}

	s.log.Info().Str("grievanceId", grievanceID).Msg("grievance raised successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Operation completed successfully",
		Data:      grievanceID,
	}, nil
}

// ---------------------------------------------------------------------------
// UpdateGrievance updates an existing grievance (respondent comment/evidence).
// Maps to .NET GrievanceSetup with OperationTypes.Update.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) UpdateGrievance(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.GrievanceVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *GrievanceVm")
	}

	s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("updating grievance")

	grievance, err := s.grievanceRepo.FirstOrDefault(ctx, "grievance_id = ?", vm.GrievanceID)
	if err != nil {
		return nil, fmt.Errorf("fetching grievance: %w", err)
	}
	if grievance == nil {
		return nil, fmt.Errorf("no record found for grievance(%s)", vm.GrievanceID)
	}

	grievance.RespondentComment = vm.RespondentComment
	grievance.RespondentEvidenceUpload = vm.RespondentEvidenceUpload
	grievance.RecordStatus = enums.StatusPendingResolution.String()

	if err := s.grievanceRepo.UpdateAndSave(ctx, grievance); err != nil {
		return nil, fmt.Errorf("updating grievance: %w", err)
	}

	s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("grievance updated successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Operation completed successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// CreateGrievanceResolution creates a new resolution for a grievance.
// Maps to .NET LogGrievanceResolution with OperationTypes.Add.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) CreateGrievanceResolution(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.GrievanceResolutionVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *GrievanceResolutionVm")
	}

	s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("creating grievance resolution")

	// Start a transaction since we update both resolution and grievance.
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	grievance, err := s.findGrievanceInTx(tx, vm.GrievanceID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Generate a unique resolution ID.
	resolutionID, err := s.generateCode(ctx, enums.SeqGrievanceComment, 8)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("generating resolution ID: %w", err)
	}

	resolution := performance.GrievanceResolution{
		GrievanceResolutionID: resolutionID,
		ResolutionComment:     vm.ResolutionComment,
		Level:                 vm.Level,
		ResolutionLevel:       resolutionLevelName(vm.Level),
		MediatorStaffID:       vm.MediatorStaffID,
		EvidenceUpload:        vm.EvidenceUpload,
		GrievanceID:           vm.GrievanceID,
	}

	// Read global settings for email notifications (best-effort, mirrors .NET try/catch).
	useActualUserMail, _ := s.getGlobalBool(ctx, "USE_ACTUAL_USER_MAIL")
	hrdGrievanceNotificationMail, _ := s.getGlobalString(ctx, "HRD_GRIEVANCE_NOTIFICATION_MAIL")

	// Determine remark statuses based on role and level.
	isHrApprover := s.userContextSvc.IsInRole(ctx, auth.RoleHrApprover)
	if isHrApprover && vm.Level == enums.ResolutionLevelHRD {
		if vm.RespondentRemark == enums.ResolutionRemarkClosed {
			// HR closes the grievance.
			resolution.ComplainantRemark = enums.ResolutionRemarkClosed
			resolution.RespondentRemark = enums.ResolutionRemarkClosed
			grievance.RecordStatus = enums.StatusClosed.String()

			// Send email notification for closure.
			s.sendGrievanceResolutionEmail(ctx, grievance, useActualUserMail, hrdGrievanceNotificationMail,
				"has been resolved and closed")

			s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("grievance closed by HRD")
		} else if vm.RespondentRemark == enums.ResolutionRemarkReEvaluate {
			// HR sends for re-evaluation.
			resolution.ComplainantRemark = enums.ResolutionRemarkReEvaluate
			resolution.RespondentRemark = enums.ResolutionRemarkReEvaluate
			grievance.RecordStatus = enums.StatusReEvaluate.String()

			// Send email notification for re-evaluation.
			s.sendGrievanceResolutionEmail(ctx, grievance, useActualUserMail, hrdGrievanceNotificationMail,
				"has been sent to your line manager for re-evaluation by HR")

			s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("grievance sent for re-evaluation by HRD")
		}
	} else {
		// Standard resolution: awaiting feedback from both parties.
		resolution.ComplainantRemark = enums.ResolutionRemarkPending
		resolution.RespondentRemark = enums.ResolutionRemarkPending
		grievance.RecordStatus = enums.StatusResolvedAwaitingFeedback.String()
	}

	if err := tx.Create(&resolution).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("saving grievance resolution: %w", err)
	}

	if err := tx.Save(grievance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("updating grievance status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	s.log.Info().Str("resolutionId", resolutionID).Msg("grievance resolution created successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Operation completed successfully",
		Data:      resolutionID,
	}, nil
}

// ---------------------------------------------------------------------------
// UpdateGrievanceResolution updates an existing resolution with feedback.
// Maps to .NET LogGrievanceResolution with OperationTypes.Update.
// Handles escalation logic when either party escalates.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) UpdateGrievanceResolution(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.GrievanceResolutionVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *GrievanceResolutionVm")
	}

	s.log.Info().
		Str("resolutionId", vm.GrievanceResolutionID).
		Str("grievanceId", vm.GrievanceID).
		Msg("updating grievance resolution")

	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	grievance, err := s.findGrievanceInTx(tx, vm.GrievanceID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var resolution performance.GrievanceResolution
	if err := tx.Where("grievance_resolution_id = ? AND soft_deleted = false", vm.GrievanceResolutionID).
		First(&resolution).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no record found for resolution(%s)", vm.GrievanceResolutionID)
		}
		return nil, fmt.Errorf("fetching resolution: %w", err)
	}

	// Read global settings for email notifications (best-effort, mirrors .NET try/catch).
	useActualUserMail, _ := s.getGlobalBool(ctx, "USE_ACTUAL_USER_MAIL")
	hrdGrievanceNotificationMail, _ := s.getGlobalString(ctx, "HRD_GRIEVANCE_NOTIFICATION_MAIL")

	// Update resolution fields.
	resolution.EvidenceUpload = vm.EvidenceUpload
	resolution.Level = vm.Level
	resolution.ResolutionLevel = resolutionLevelName(vm.Level)
	resolution.RespondentRemark = vm.RespondentRemark
	resolution.ComplainantRemark = vm.ComplainantRemark
	resolution.RespondentFeedback = vm.RespondentFeedback
	resolution.ComplainantFeedback = vm.ComplainantFeedback
	resolution.ResolutionComment = vm.ResolutionComment

	// Escalation logic: when both parties have responded and at least one escalated.
	if vm.RespondentRemark != enums.ResolutionRemarkPending && vm.ComplainantRemark != enums.ResolutionRemarkPending {
		if vm.RespondentRemark == enums.ResolutionRemarkEscalated || vm.ComplainantRemark == enums.ResolutionRemarkEscalated {
			grievance.RecordStatus = enums.StatusEscalated.String()

			if err := s.escalateGrievance(ctx, grievance); err != nil {
				s.log.Warn().Err(err).Msg("escalation hierarchy lookup failed, defaulting to HRD")
				grievance.CurrentMediatorStaffID = "HRD"
				grievance.CurrentResolutionLevel = enums.ResolutionLevelHRD
			}

			// Send escalation email notification.
			s.sendGrievanceEscalationEmail(ctx, grievance, useActualUserMail, hrdGrievanceNotificationMail)

			s.log.Info().
				Str("grievanceId", vm.GrievanceID).
				Int("newLevel", int(grievance.CurrentResolutionLevel)).
				Msg("grievance escalated")
		}
	}

	// Both parties accepted: close the grievance.
	if vm.RespondentRemark == enums.ResolutionRemarkAccepted && vm.ComplainantRemark == enums.ResolutionRemarkAccepted {
		grievance.RecordStatus = enums.StatusClosed.String()

		// Send resolution/closure email notification.
		s.sendGrievanceResolutionEmail(ctx, grievance, useActualUserMail, hrdGrievanceNotificationMail,
			"has been resolved and closed")

		s.log.Info().Str("grievanceId", vm.GrievanceID).Msg("grievance closed by mutual acceptance")
	}

	if err := tx.Save(&resolution).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("updating resolution: %w", err)
	}
	if err := tx.Save(grievance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("updating grievance: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	s.log.Info().Str("resolutionId", vm.GrievanceResolutionID).Msg("grievance resolution updated successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Operation completed successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// GetStaffGrievances retrieves all grievances a staff member is involved in
// (as complainant, respondent, or mediator).
// Maps to .NET GetStaffGrievances.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) GetStaffGrievances(ctx context.Context, staffID string) (interface{}, error) {
	s.log.Info().Str("staffId", staffID).Msg("fetching staff grievances")

	// Fetch grievances where the staff is complainant, respondent, or current mediator.
	grievances, err := s.grievanceRepo.WhereWithPreload(
		ctx,
		[]string{"GrievanceResolutions"},
		"complainant_staff_id = ? OR respondent_staff_id = ? OR current_mediator_staff_id = ?",
		staffID, staffID, staffID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying grievances for staff: %w", err)
	}

	// Also fetch resolutions where staff is the mediator (to include grievances
	// they mediated but are not directly complainant/respondent/current mediator).
	var mediatedResolutions []performance.GrievanceResolution
	err = s.db.WithContext(ctx).
		Where("mediator_staff_id = ? AND soft_deleted = false", staffID).
		Preload("Grievance").
		Preload("Grievance.GrievanceResolutions").
		Find(&mediatedResolutions).Error
	if err != nil {
		return nil, fmt.Errorf("querying mediated resolutions: %w", err)
	}

	// Merge mediated grievances into the main list (dedup by GrievanceID).
	existingIDs := make(map[string]bool, len(grievances))
	for _, g := range grievances {
		existingIDs[g.GrievanceID] = true
	}
	for _, r := range mediatedResolutions {
		if r.Grievance != nil && !existingIDs[r.GrievanceID] {
			grievances = append(grievances, *r.Grievance)
			existingIDs[r.GrievanceID] = true
		}
	}

	// Map to view models with enriched employee names.
	vmList, err := s.mapGrievancesToVm(ctx, grievances)
	if err != nil {
		return nil, err
	}

	// Sort by creation date descending.
	sort.Slice(vmList, func(i, j int) bool {
		if vmList[i].DateCreated == nil || vmList[j].DateCreated == nil {
			return vmList[i].DateCreated != nil
		}
		return vmList[i].DateCreated.After(*vmList[j].DateCreated)
	})

	s.log.Info().Int("count", len(vmList)).Str("staffId", staffID).Msg("staff grievances fetched")

	return &performance.GenericListVm{
		BaseAPIResponse: performance.BaseAPIResponse{
			HasError: false,
			Message:  "Operation completed successfully",
		},
		ListData:    vmList,
		TotalRecord: len(vmList),
	}, nil
}

// ---------------------------------------------------------------------------
// GetGrievancesReport retrieves all grievances (admin report view).
// Maps to .NET GetGrievancesReport.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) GetGrievancesReport(ctx context.Context) (interface{}, error) {
	s.log.Info().Msg("fetching grievances report")

	grievances, err := s.grievanceRepo.GetAllIncluding(ctx, "GrievanceResolutions")
	if err != nil {
		return nil, fmt.Errorf("querying all grievances: %w", err)
	}

	vmList, err := s.mapGrievancesToVm(ctx, grievances)
	if err != nil {
		return nil, err
	}

	// Sort by creation date descending.
	sort.Slice(vmList, func(i, j int) bool {
		if vmList[i].DateCreated == nil || vmList[j].DateCreated == nil {
			return vmList[i].DateCreated != nil
		}
		return vmList[i].DateCreated.After(*vmList[j].DateCreated)
	})

	s.log.Info().Int("count", len(vmList)).Msg("grievances report fetched")

	return &performance.GenericListVm{
		BaseAPIResponse: performance.BaseAPIResponse{
			HasError: false,
			Message:  "Operation completed successfully",
		},
		ListData:    vmList,
		TotalRecord: len(vmList),
	}, nil
}

// ===========================================================================
// Private helpers
// ===========================================================================

// resolutionLevelName returns the human-readable name for a ResolutionLevel.
// Mirrors .NET's Humanize() call on the ResolutionLevl enum.
func resolutionLevelName(level enums.ResolutionLevel) string {
	switch level {
	case enums.ResolutionLevelSBU:
		return "SBU"
	case enums.ResolutionLevelDepartment:
		return "Department"
	case enums.ResolutionLevelHRD:
		return "HRD"
	default:
		return fmt.Sprintf("Level(%d)", level)
	}
}

// findGrievanceInTx locates a grievance within an existing transaction.
func (s *grievanceManagementService) findGrievanceInTx(tx *gorm.DB, grievanceID string) (*performance.Grievance, error) {
	var grievance performance.Grievance
	err := tx.Where("grievance_id = ? AND soft_deleted = false", grievanceID).
		First(&grievance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no record found for grievance(%s)", grievanceID)
		}
		return nil, fmt.Errorf("fetching grievance: %w", err)
	}
	return &grievance, nil
}

// escalateGrievance advances the resolution level based on the current level.
// SBU -> Department (head of dept) or HRD if already at head of dept.
// Department -> HRD.
// HRD stays at HRD.
func (s *grievanceManagementService) escalateGrievance(ctx context.Context, grievance *performance.Grievance) error {
	switch grievance.CurrentResolutionLevel {
	case enums.ResolutionLevelSBU:
		empData, err := s.getEmployeeData(ctx, grievance.CurrentMediatorStaffID)
		if err != nil {
			return err
		}
		if grievance.CurrentMediatorStaffID == empData.HeadOfDeptID {
			// Already at head of department; escalate to HRD.
			grievance.CurrentMediatorStaffID = "HRD"
			grievance.CurrentResolutionLevel = enums.ResolutionLevelHRD
		} else {
			// Escalate to head of department.
			grievance.CurrentMediatorStaffID = empData.HeadOfDeptID
			grievance.CurrentResolutionLevel = enums.ResolutionLevelDepartment
		}

	case enums.ResolutionLevelDepartment:
		grievance.CurrentMediatorStaffID = "HRD"
		grievance.CurrentResolutionLevel = enums.ResolutionLevelHRD

	case enums.ResolutionLevelHRD:
		// Already at the highest level; stays at HRD.
		grievance.CurrentMediatorStaffID = "HRD"
		grievance.CurrentResolutionLevel = enums.ResolutionLevelHRD
	}

	return nil
}

// mapGrievancesToVm converts domain Grievance entities to GrievanceVm DTOs
// with employee name enrichment.
func (s *grievanceManagementService) mapGrievancesToVm(ctx context.Context, grievances []performance.Grievance) ([]performance.GrievanceVm, error) {
	vmList := make([]performance.GrievanceVm, 0, len(grievances))

	for i := range grievances {
		g := &grievances[i]

		vm := performance.GrievanceVm{
			GrievanceID:               g.GrievanceID,
			GrievanceType:             int(g.GrievanceType),
			Description:               g.Description,
			RespondentComment:         g.RespondentComment,
			SubjectID:                 g.SubjectID,
			Subject:                   g.Subject,
			ReviewPeriodID:            g.ReviewPeriodID,
			ComplainantStaffID:        g.ComplainantStaffID,
			ComplainantEvidenceUpload: g.ComplainantEvidenceUpload,
			CurrentResolutionLevel:    g.CurrentResolutionLevel,
			CurrentMediatorStaffID:    g.CurrentMediatorStaffID,
			RespondentStaffID:         g.RespondentStaffID,
			RespondentEvidenceUpload:  g.RespondentEvidenceUpload,
		}

		// Carry base entity fields.
		vm.BaseAuditVm = performance.BaseAuditVm{
			CreatedBy:   g.CreatedBy,
			DateCreated: g.CreatedAt,
			IsActive:    g.IsActive,
			Status:      g.RecordStatus,
			DateUpdated: g.UpdatedAt,
			UpdatedBy:   g.UpdatedBy,
		}

		// Enrich with employee names (best-effort; log warning on failures).
		vm.ComplainantStaff = s.safeGetEmployeeName(ctx, g.ComplainantStaffID)
		vm.RespondentStaff = s.safeGetEmployeeName(ctx, g.RespondentStaffID)

		if g.CurrentResolutionLevel == enums.ResolutionLevelHRD {
			vm.CurrentMediatorStaff = "HRD"
		} else {
			vm.CurrentMediatorStaff = s.safeGetEmployeeName(ctx, g.CurrentMediatorStaffID)
		}

		// Map resolutions.
		vm.GrievanceResolutions = make([]performance.GrievanceResolutionVm, 0, len(g.GrievanceResolutions))
		for j := range g.GrievanceResolutions {
			r := &g.GrievanceResolutions[j]
			rvm := performance.GrievanceResolutionVm{
				GrievanceResolutionID: r.GrievanceResolutionID,
				ResolutionComment:     r.ResolutionComment,
				ResolutionLevel:       r.ResolutionLevel,
				Level:                 r.Level,
				MediatorStaffID:       r.MediatorStaffID,
				MediatorStaff:         s.safeGetEmployeeName(ctx, r.MediatorStaffID),
				EvidenceUpload:        r.EvidenceUpload,
				RespondentFeedback:    r.RespondentFeedback,
				ComplainantFeedback:   r.ComplainantFeedback,
				ComplainantRemark:     r.ComplainantRemark,
				RespondentRemark:      r.RespondentRemark,
				GrievanceID:           r.GrievanceID,
			}
			rvm.BaseAuditVm = performance.BaseAuditVm{
				CreatedBy:   r.CreatedBy,
				DateCreated: r.CreatedAt,
				IsActive:    r.IsActive,
				Status:      r.RecordStatus,
				DateUpdated: r.UpdatedAt,
				UpdatedBy:   r.UpdatedBy,
			}
			vm.GrievanceResolutions = append(vm.GrievanceResolutions, rvm)
		}

		// IsResolved if the status is Closed.
		vm.IsResolved = strings.EqualFold(g.RecordStatus, enums.StatusClosed.String())

		vmList = append(vmList, vm)
	}

	return vmList, nil
}

// employeeDataResult is a service-local struct for extracting employee data
// from the untyped ErpEmployeeService response.
type employeeDataResult struct {
	SupervisorID string
	HeadOfDeptID string
	FullName     string
}

// getEmployeeData calls the ERP employee service and extracts the fields
// needed for grievance processing.
func (s *grievanceManagementService) getEmployeeData(ctx context.Context, staffID string) (*employeeDataResult, error) {
	result, err := s.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("getting employee data for %s: %w", staffID, err)
	}
	if result == nil {
		return nil, fmt.Errorf("employee not found: %s", staffID)
	}

	// The ErpEmployeeService returns interface{} -- try known concrete types first.
	switch v := result.(type) {
	case *erp.EmployeeData:
		return &employeeDataResult{
			SupervisorID: v.SupervisorID,
			HeadOfDeptID: v.HeadOfDeptID,
			FullName:     strings.TrimSpace(v.FullName()),
		}, nil
	case erp.EmployeeData:
		return &employeeDataResult{
			SupervisorID: v.SupervisorID,
			HeadOfDeptID: v.HeadOfDeptID,
			FullName:     strings.TrimSpace(v.FullName()),
		}, nil
	case *erp.EmployeeErpDetailsDTO:
		return &employeeDataResult{
			SupervisorID: v.SupervisorID,
			HeadOfDeptID: v.HeadOfDeptID,
			FullName:     strings.TrimSpace(v.FullName()),
		}, nil
	case erp.EmployeeErpDetailsDTO:
		return &employeeDataResult{
			SupervisorID: v.SupervisorID,
			HeadOfDeptID: v.HeadOfDeptID,
			FullName:     strings.TrimSpace(v.FullName()),
		}, nil
	case map[string]interface{}:
		return &employeeDataResult{
			SupervisorID: fmt.Sprintf("%v", v["supervisorId"]),
			HeadOfDeptID: fmt.Sprintf("%v", v["headOfDeptId"]),
			FullName:     fmt.Sprintf("%v", v["fullName"]),
		}, nil
	default:
		s.log.Warn().Str("staffId", staffID).Type("type", result).
			Msg("unable to extract employee details from unknown result type; using fallback")
		return &employeeDataResult{
			SupervisorID: "",
			HeadOfDeptID: "",
			FullName:     staffID,
		}, nil
	}
}

// safeGetEmployeeName fetches an employee's full name, returning the staffID
// on failure rather than propagating the error (best-effort enrichment).
func (s *grievanceManagementService) safeGetEmployeeName(ctx context.Context, staffID string) string {
	if staffID == "" || staffID == "HRD" {
		return staffID
	}
	data, err := s.getEmployeeData(ctx, staffID)
	if err != nil {
		s.log.Warn().Err(err).Str("staffId", staffID).Msg("failed to fetch employee name")
		return staffID
	}
	if data.FullName != "" {
		return strings.TrimSpace(data.FullName)
	}
	return staffID
}

// generateCode generates a unique code for the given sequence type.
// Mirrors .NET BaseService.GenerateCode using the pms.sequence_numbers table.
func (s *grievanceManagementService) generateCode(ctx context.Context, seqType enums.SequenceNumberTypes, length int) (string, error) {
	var seq struct {
		NextNumber int64  `gorm:"column:next_number"`
		Prefix     string `gorm:"column:prefix"`
		UsePrefix  bool   `gorm:"column:use_prefix"`
	}

	// Atomically read and increment the next number.
	err := s.db.WithContext(ctx).
		Raw(`UPDATE pms.sequence_numbers SET next_number = next_number + 1
		     WHERE sequence_number_type = ? RETURNING next_number - 1 AS next_number, prefix, use_prefix`, int(seqType)).
		Scan(&seq).Error
	if err != nil {
		return "", fmt.Errorf("incrementing sequence %d: %w", seqType, err)
	}

	code := fmt.Sprintf("%0*d", length, seq.NextNumber)
	if seq.UsePrefix && seq.Prefix != "" {
		code = seq.Prefix + code
	}
	return code, nil
}

// ---------------------------------------------------------------------------
// Global setting helpers (best-effort, mirrors .NET try/catch pattern)
// ---------------------------------------------------------------------------

// getGlobalBool reads a boolean global setting, returning false on any error.
func (s *grievanceManagementService) getGlobalBool(ctx context.Context, key string) (bool, error) {
	if s.globalSettingSvc == nil {
		return false, nil
	}
	return s.globalSettingSvc.GetBoolValue(ctx, key)
}

// getGlobalString reads a string global setting, returning "" on any error.
func (s *grievanceManagementService) getGlobalString(ctx context.Context, key string) (string, error) {
	if s.globalSettingSvc == nil {
		return "", nil
	}
	return s.globalSettingSvc.GetStringValue(ctx, key)
}

// ---------------------------------------------------------------------------
// Email notification helpers
// Mirrors the .NET email sending logic in LogGrievanceResolution.
// ---------------------------------------------------------------------------

// sendGrievanceResolutionEmail sends notification emails when a grievance
// resolution is closed or re-evaluated. Mirrors .NET's inline email sending
// inside LogGrievanceResolution for both the complainant and HRD.
func (s *grievanceManagementService) sendGrievanceResolutionEmail(
	ctx context.Context,
	grievance *performance.Grievance,
	useActualUserMail bool,
	hrdGrievanceNotificationMail string,
	actionDescription string,
) {
	if s.emailSvc == nil {
		s.log.Warn().Msg("email service not configured, skipping grievance resolution notification")
		return
	}

	complainant, err := s.getEmployeeData(ctx, grievance.ComplainantStaffID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to fetch complainant for resolution email")
		return
	}

	respondent, err := s.getEmployeeData(ctx, grievance.RespondentStaffID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to fetch respondent for resolution email")
		return
	}

	action := "PMS_GRIEVANCE_RESOLUTION_TRIGGER"
	subject := fmt.Sprintf("GRIEVANCE RESOLUTION ON: %s (ID: %s)", grievance.Subject, grievance.SubjectID)
	now := time.Now().Format("02 Jan 2006, 03:04:05 PM")

	// Email to complainant.
	body := fmt.Sprintf(
		`<p>Dear %s,</p><br/><p>This is to notify you that the grievance raised on the above subject %s`+
			`<br/> Date: %s </p>`+
			`<br/><p>Regards, <br/>Performance Management System (PMS)</p>`,
		complainant.FullName, actionDescription, now,
	)
	emailTo := ""
	if useActualUserMail {
		emailTo = s.getEmployeeEmail(ctx, grievance.ComplainantStaffID)
	}

	if err := s.emailSvc.SendEmail(ctx, emailTo, subject, body); err != nil {
		s.log.Warn().Err(err).Str("action", action).Msg("failed to send resolution email to complainant")
	}

	// Email to HRD if configured.
	if hrdGrievanceNotificationMail != "" {
		hrdBody := fmt.Sprintf(
			`<p>Dear Sir/Madam,</p><br/><p>This is to notify you that grievance on the above subject %s.`+
				`<br/> Complainant: %s (%s)`+
				`<br/> Respondent: %s (%s)`+
				`<br/> Date: %s </p>`+
				`<br/><p>Regards, <br/>Performance Management System (PMS)</p>`,
			actionDescription,
			complainant.FullName, grievance.ComplainantStaffID,
			respondent.FullName, grievance.RespondentStaffID,
			now,
		)
		if err := s.emailSvc.SendEmail(ctx, hrdGrievanceNotificationMail, subject, hrdBody); err != nil {
			s.log.Warn().Err(err).Str("action", action).Msg("failed to send resolution email to HRD")
		}
	}
}

// sendGrievanceEscalationEmail sends notification emails when a grievance is
// escalated. Mirrors .NET's inline email sending inside LogGrievanceResolution
// for the escalation case (both complainant and HRD).
func (s *grievanceManagementService) sendGrievanceEscalationEmail(
	ctx context.Context,
	grievance *performance.Grievance,
	useActualUserMail bool,
	hrdGrievanceNotificationMail string,
) {
	if s.emailSvc == nil {
		s.log.Warn().Msg("email service not configured, skipping grievance escalation notification")
		return
	}

	complainant, err := s.getEmployeeData(ctx, grievance.ComplainantStaffID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to fetch complainant for escalation email")
		return
	}

	respondent, err := s.getEmployeeData(ctx, grievance.RespondentStaffID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to fetch respondent for escalation email")
		return
	}

	action := "PMS_GRIEVANCE_ESCALATION_TRIGGER"
	levelName := resolutionLevelName(grievance.CurrentResolutionLevel)
	subject := fmt.Sprintf("GRIEVANCE ESCALATION ON: %s (ID: %s)", grievance.Subject, grievance.SubjectID)
	now := time.Now().Format("02 Jan 2006, 03:04:05 PM")

	// Email to complainant.
	body := fmt.Sprintf(
		`<p>Dear %s,</p><br/><p>This is to notify you that the grievance has been on the above subject has been escalated to %s`+
			`<br/> Date: %s </p>`+
			`<br/><p>Regards, <br/>Performance Management System (PMS)</p>`,
		complainant.FullName, levelName, now,
	)
	emailTo := ""
	if useActualUserMail {
		emailTo = s.getEmployeeEmail(ctx, grievance.ComplainantStaffID)
	}

	if err := s.emailSvc.SendEmail(ctx, emailTo, subject, body); err != nil {
		s.log.Warn().Err(err).Str("action", action).Msg("failed to send escalation email to complainant")
	}

	// Email to HRD if configured.
	if hrdGrievanceNotificationMail != "" {
		hrdBody := fmt.Sprintf(
			`<p>Dear Sir/Madam,</p><br/><p>This is to notify you that grievance on the above subject has been escalated to %s.`+
				`<br/> Complainant: %s (%s)`+
				`<br/> Respondent: %s (%s)`+
				`<br/> Subject: %s`+
				`<br/> Subject ID: %s`+
				`<br/> Description: %s`+
				`<br/> Date: %s </p>`+
				`<br/><p>Regards, <br/>Performance Management System (PMS)</p>`,
			levelName,
			complainant.FullName, grievance.ComplainantStaffID,
			respondent.FullName, grievance.RespondentStaffID,
			grievance.Subject,
			grievance.SubjectID,
			grievance.Description,
			now,
		)
		if err := s.emailSvc.SendEmail(ctx, hrdGrievanceNotificationMail, subject, hrdBody); err != nil {
			s.log.Warn().Err(err).Str("action", action).Msg("failed to send escalation email to HRD")
		}
	}
}

// getEmployeeEmail fetches an employee's email address by staff ID.
// Returns empty string on failure (best-effort).
func (s *grievanceManagementService) getEmployeeEmail(ctx context.Context, staffID string) string {
	result, err := s.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
	if err != nil || result == nil {
		return ""
	}
	switch v := result.(type) {
	case *erp.EmployeeData:
		return v.EmailAddress
	case erp.EmployeeData:
		return v.EmailAddress
	case *erp.EmployeeErpDetailsDTO:
		return v.EmailAddress
	case erp.EmployeeErpDetailsDTO:
		return v.EmailAddress
	default:
		return ""
	}
}

// ---------------------------------------------------------------------------
// LogRequestAsync creates or updates a FeedbackRequestLog and sends a
// notification email to the assignee.
// Mirrors the .NET GrievanceManagementService.LogRequestAsync method.
// ---------------------------------------------------------------------------
func (s *grievanceManagementService) LogRequestAsync(
	ctx context.Context,
	referenceID string,
	assignerStaffID string,
	assigneeStaffID string,
	feedbackRequestType enums.FeedbackRequestType,
	hasSLA bool,
) error {
	if assignerStaffID == "" {
		return nil
	}

	assigner, err := s.getEmployeeData(ctx, assignerStaffID)
	if err != nil || assigner == nil {
		s.log.Warn().Err(err).Str("assignerStaffID", assignerStaffID).Msg("assigner not found, skipping log request")
		return nil
	}

	assignee, err := s.getEmployeeData(ctx, assigneeStaffID)
	if err != nil || assignee == nil {
		s.log.Warn().Err(err).Str("assigneeStaffID", assigneeStaffID).Msg("assignee not found, skipping log request")
		return nil
	}

	// Check if a log already exists for this reference + assignee + type.
	var existingLog performance.FeedbackRequestLog
	err = s.db.WithContext(ctx).
		Where("reference_id = ? AND assigned_staff_id = ? AND feedback_request_type = ? AND soft_deleted = false",
			referenceID, assigneeStaffID, int(feedbackRequestType)).
		First(&existingLog).Error

	if err == nil {
		// Existing log found: reset it (mirrors .NET update path).
		now := time.Now().UTC()
		existingLog.TimeInitiated = now
		existingLog.TimeCompleted = nil
		existingLog.RecordStatus = enums.StatusActive.String()

		if err := s.db.WithContext(ctx).Save(&existingLog).Error; err != nil {
			return fmt.Errorf("updating existing feedback request log: %w", err)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new log.
		logID, err := s.generateCode(ctx, enums.SeqFeedbackRequest, 15)
		if err != nil {
			return fmt.Errorf("generating feedback request log ID: %w", err)
		}

		now := time.Now().UTC()
		newLog := performance.FeedbackRequestLog{
			FeedbackRequestLogID: logID,
			FeedbackRequestType:  feedbackRequestType,
			ReferenceID:          referenceID,
			TimeInitiated:        now,
			AssignedStaffID:      assigneeStaffID,
			RequestOwnerStaffID:  assignerStaffID,
			HasSLA:               hasSLA,
		}
		newLog.RecordStatus = enums.StatusActive.String()
		newLog.IsActive = true

		// Check vacation rule: if assignee has a vacation rule, redirect to delegate.
		vacationResp, vacationErr := s.HasVacationRule(ctx, assigneeStaffID, now)
		if vacationErr == nil && vacationResp.IsSuccess {
			assigneeStaffID = vacationResp.DelegateStaffID
			newLog.AssignedStaffID = vacationResp.DelegateStaffID
			// Re-fetch assignee data for the delegate.
			if delegateData, dErr := s.getEmployeeData(ctx, assigneeStaffID); dErr == nil && delegateData != nil {
				assignee = delegateData
			}
		}

		if assigner.FullName != "" {
			newLog.RequestOwnerStaffName = strings.TrimSpace(assigner.FullName)
		}
		if assignee.FullName != "" {
			newLog.AssignedStaffName = strings.TrimSpace(assignee.FullName)
		}

		// Determine review period ID.
		reviewPeriodID := s.resolveReviewPeriodID(ctx, assigneeStaffID, referenceID, feedbackRequestType)
		if reviewPeriodID != "" {
			newLog.ReviewPeriodID = reviewPeriodID
		}

		if err := s.db.WithContext(ctx).Create(&newLog).Error; err != nil {
			return fmt.Errorf("creating feedback request log: %w", err)
		}
	} else {
		return fmt.Errorf("querying feedback request log: %w", err)
	}

	// Send notification email to assignee.
	s.sendAssigneeNotificationEmail(ctx, assigneeStaffID, assignee.FullName, feedbackRequestType)

	return nil
}

// resolveReviewPeriodID determines the review period ID for a feedback request log.
// Mirrors the .NET logic of trying GetActiveReviewPeriod first, then
// GetStaffActiveReviewPeriod, then checking ReviewPeriodExtensions.
func (s *grievanceManagementService) resolveReviewPeriodID(
	ctx context.Context,
	assigneeStaffID string,
	referenceID string,
	feedbackRequestType enums.FeedbackRequestType,
) string {
	if s.reviewPeriodSvc == nil {
		return ""
	}

	// Try active review period first.
	rpResult, err := s.reviewPeriodSvc.GetActiveReviewPeriod(ctx)
	if err == nil && rpResult != nil {
		if !rpResult.HasError && rpResult.PerformanceReviewPeriod != nil {
			return rpResult.PerformanceReviewPeriod.PeriodID
		}
	}

	// Fall back to staff-specific active period.
	rpResult, err = s.reviewPeriodSvc.GetStaffActiveReviewPeriod(ctx, assigneeStaffID)
	if err == nil && rpResult != nil {
		if !rpResult.HasError && rpResult.PerformanceReviewPeriod != nil {
			return rpResult.PerformanceReviewPeriod.PeriodID
		}
	}

	// For ReviewPeriodExtension type, look up the extension record.
	if feedbackRequestType == enums.FeedbackRequestReviewPeriodExtension {
		var ext performance.ReviewPeriodExtension
		if err := s.db.WithContext(ctx).
			Where("review_period_extension_id = ?", referenceID).
			First(&ext).Error; err == nil {
			return ext.ReviewPeriodID
		}
	}

	return ""
}

// sendAssigneeNotificationEmail sends an email notification to the assigned
// staff member about a new feedback request.
// Mirrors the .NET email sending at the end of LogRequestAsync.
func (s *grievanceManagementService) sendAssigneeNotificationEmail(
	ctx context.Context,
	assigneeStaffID string,
	assigneeName string,
	feedbackRequestType enums.FeedbackRequestType,
) {
	if s.emailSvc == nil {
		return
	}

	useActualUserMail, _ := s.getGlobalBool(ctx, "USE_ACTUAL_USER_MAIL")

	requestTypeName := feedbackRequestTypeName(feedbackRequestType)
	action := fmt.Sprintf("PMS_%s", feedbackRequestTypeEnumName(feedbackRequestType))
	subject := fmt.Sprintf("NEW ASSIGNED REQUEST ON: %s", strings.ToUpper(requestTypeName))
	now := time.Now().Format("02 Jan 2006, 03:04:05 PM")

	body := fmt.Sprintf(
		`<p>Dear %s,</p><br/><p>You have been assigned a request on Performance Management System on %s.<br/>`+
			`Kindly logon to your account using this link <a href="https://pms.cbn.gov.ng/">PMS Login</a> and treat the assigned request.</p>`+
			`<br/><p>Regards, <br/>Performance Management System (PMS)</p>`,
		strings.TrimSpace(assigneeName), now,
	)

	emailTo := ""
	if useActualUserMail {
		emailTo = s.getEmployeeEmail(ctx, assigneeStaffID)
	}

	if err := s.emailSvc.SendEmail(ctx, emailTo, subject, body); err != nil {
		s.log.Warn().Err(err).Str("action", action).Msg("failed to send assignee notification email")
	}
}

// ---------------------------------------------------------------------------
// HasVacationRule checks whether a staff member has an active vacation rule
// delegation. If so, returns the delegate's staff ID.
// Mirrors the .NET GrievanceManagementService.HasVacationRule method.
// ---------------------------------------------------------------------------

// VacationRuleResult is the response from HasVacationRule.
type VacationRuleResult struct {
	IsSuccess      bool
	DelegateStaffID string
	Message        string
	Errors         []string
}

func (s *grievanceManagementService) HasVacationRule(
	ctx context.Context,
	staffID string,
	startDate time.Time,
) (*VacationRuleResult, error) {
	result := &VacationRuleResult{
		IsSuccess: false,
		Message:   "An error occurred",
	}

	if s.erpDataDB == nil {
		return result, fmt.Errorf("ERP data database not configured")
	}

	// Get the employee's username for matching vacation rules.
	empData, err := s.getEmployeeData(ctx, staffID)
	if err != nil || empData == nil {
		return result, fmt.Errorf("employee not found for vacation rule check: %s", staffID)
	}

	// Query vacation rules where the rule owner matches and the date range
	// covers the start date. Mirrors the .NET LINQ query.
	type vacationRuleRow struct {
		RuleOwner  *string    `db:"rule_owner"`
		AssignedTo *string    `db:"assigned_to"`
		BeginDate  time.Time  `db:"begin_date"`
		EndDate    *time.Time `db:"end_date"`
	}

	var rules []vacationRuleRow
	query := `SELECT rule_owner, assigned_to, begin_date, end_date
		FROM dbo.VACATIONSRULE_DATA
		WHERE rule_owner IS NOT NULL
		  AND begin_date <= @p1
		  AND (end_date IS NULL OR end_date >= @p2)`

	if err := s.erpDataDB.SelectContext(ctx, &rules, query, startDate, startDate); err != nil {
		s.log.Error().Err(err).Msg("failed to query vacation rules")
		return result, fmt.Errorf("querying vacation rules: %w", err)
	}

	// Find the matching rule by username (case-insensitive, mirrors .NET).
	userName := strings.ToLower(strings.TrimSpace(empData.FullName))
	// The .NET code matches on user.UserName -- extract it from EmployeeData.
	if data, ok := s.extractUserName(ctx, staffID); ok {
		userName = strings.ToLower(data)
	}

	digitRegex := regexp.MustCompile(`\d+`)

	for _, rule := range rules {
		if rule.RuleOwner == nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(*rule.RuleOwner), userName) {
			if rule.AssignedTo != nil {
				// Extract numeric staff ID from the assigned_to field (mirrors .NET Regex.Match).
				match := digitRegex.FindString(*rule.AssignedTo)
				if match != "" {
					result.DelegateStaffID = match
					result.IsSuccess = true
					result.Message = "Vacation rule is configured"
					return result, nil
				}
			}
		}
	}

	result.Errors = append(result.Errors, "Vacation rule is currently not configured")
	return result, nil
}

// extractUserName attempts to extract the UserName field from the ERP employee
// detail response. Returns the username and true on success.
func (s *grievanceManagementService) extractUserName(ctx context.Context, staffID string) (string, bool) {
	result, err := s.erpEmployeeSvc.GetEmployeeDetail(ctx, staffID)
	if err != nil || result == nil {
		return "", false
	}
	switch v := result.(type) {
	case *erp.EmployeeData:
		return v.UserName, true
	case erp.EmployeeData:
		return v.UserName, true
	case *erp.EmployeeErpDetailsDTO:
		return v.UserName, true
	case erp.EmployeeErpDetailsDTO:
		return v.UserName, true
	default:
		return "", false
	}
}

// ---------------------------------------------------------------------------
// Enum name helpers for FeedbackRequestType
// ---------------------------------------------------------------------------

// feedbackRequestTypeName returns a human-readable name for a
// FeedbackRequestType. Mirrors .NET's Humanize() on the enum.
func feedbackRequestTypeName(t enums.FeedbackRequestType) string {
	switch t {
	case enums.FeedbackRequestWorkProductEvaluation:
		return "Work Product Evaluation"
	case enums.FeedbackRequestObjectivePlanning:
		return "Objective Planning"
	case enums.FeedbackRequestProjectPlanning:
		return "Project Planning"
	case enums.FeedbackRequestCommitteePlanning:
		return "Committee Planning"
	case enums.FeedbackRequestWorkProductFeedback:
		return "Work Product Feedback"
	case enums.FeedbackRequest360ReviewFeedback:
		return "360 Review Feedback"
	case enums.FeedbackRequestWorkProductPlanning:
		return "Work Product Planning"
	case enums.FeedbackRequestCompetencyReview:
		return "Competency Review"
	case enums.FeedbackRequestReviewPeriod:
		return "Review Period"
	case enums.FeedbackRequestReviewPeriodExtension:
		return "Review Period Extension"
	case enums.FeedbackRequestProjectMemberAssignment:
		return "Project Member Assignment"
	case enums.FeedbackRequestCommitteeMemberAssignment:
		return "Committee Member Assignment"
	case enums.FeedbackRequestPeriodObjectiveOutcome:
		return "Period Objective Outcome"
	case enums.FeedbackRequestDeptObjectiveOutcome:
		return "Department Objective Outcome"
	case enums.FeedbackRequestReviewPeriod360Review:
		return "Review Period 360 Review"
	case enums.FeedbackRequestProjectWorkProductDef:
		return "Project Work Product Definition"
	case enums.FeedbackRequestCommitteeWorkProductDef:
		return "Committee Work Product Definition"
	default:
		return fmt.Sprintf("FeedbackRequestType(%d)", int(t))
	}
}

// feedbackRequestTypeEnumName returns the Go enum constant name for a
// FeedbackRequestType. Mirrors .NET's Enum.GetName().
func feedbackRequestTypeEnumName(t enums.FeedbackRequestType) string {
	switch t {
	case enums.FeedbackRequestWorkProductEvaluation:
		return "WorkProductEvaluation"
	case enums.FeedbackRequestObjectivePlanning:
		return "ObjectivePlanning"
	case enums.FeedbackRequestProjectPlanning:
		return "ProjectPlanning"
	case enums.FeedbackRequestCommitteePlanning:
		return "CommitteePlanning"
	case enums.FeedbackRequestWorkProductFeedback:
		return "WorkProductFeedback"
	case enums.FeedbackRequest360ReviewFeedback:
		return "_360ReviewFeedback"
	case enums.FeedbackRequestWorkProductPlanning:
		return "WorkProductPlanning"
	case enums.FeedbackRequestCompetencyReview:
		return "CompetencyReview"
	case enums.FeedbackRequestReviewPeriod:
		return "ReviewPeriod"
	case enums.FeedbackRequestReviewPeriodExtension:
		return "ReviewPeriodExtension"
	case enums.FeedbackRequestProjectMemberAssignment:
		return "ProjectMemberAssignment"
	case enums.FeedbackRequestCommitteeMemberAssignment:
		return "CommitteeMemberAssignment"
	case enums.FeedbackRequestPeriodObjectiveOutcome:
		return "PeriodObjectiveOutcome"
	case enums.FeedbackRequestDeptObjectiveOutcome:
		return "DeptObjectiveOutcome"
	case enums.FeedbackRequestReviewPeriod360Review:
		return "ReviewPeriod360Review"
	case enums.FeedbackRequestProjectWorkProductDef:
		return "ProjectWorkProductDef"
	case enums.FeedbackRequestCommitteeWorkProductDef:
		return "CommitteeWorkProductDef"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}
