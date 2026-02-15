package service

import (
	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
)

// Container holds all service implementations.
// This is the Go equivalent of the .NET DI container for services.
type Container struct {
	Performance  PerformanceManagementService
	Competency   CompetencyService
	PmsSetup     PmsSetupService
	ReviewPeriod ReviewPeriodService
	Grievance    GrievanceManagementService
	RoleMgt      RoleManagementService
	StaffMgt     StaffManagementService
	Organogram   OrganogramService
	ErpEmployee  ErpEmployeeService
	GlobalSetting GlobalSettingService
	Auth         AuthService
	Email        EmailService
	FileStorage  FileStorageService
	Notification NotificationService
	Encryption   EncryptionService
	AD           ActiveDirectoryService
	UserContext   UserContextService
	PasswordGen  PasswordGenerator
}

// New creates the service container with all dependencies wired up.
// This replicates .NET's AddScoped/AddTransient service registrations.
func New(repos *repository.Container, cfg *config.Config, log zerolog.Logger) *Container {
	// --- Foundation services (no service dependencies) ---
	gsSvc := newGlobalSettingService(repos, log)
	adSvc := newActiveDirectoryService(cfg.ActiveDirectory, log)
	authSvc := newAuthService(repos, cfg, log)
	ucSvc := newUserContextService(log)
	pwGen := newPasswordGenerator()
	emailSvc := newEmailService(repos, cfg, log, gsSvc)
	encSvc := newEncryptionService(cfg, log)
	fsSvc := newFileStorageService(cfg, log)
	pmsSetupSvc := newPmsSetupService(repos, cfg, log, encSvc)
	userMgr := NewUserManagementService(repos, log)

	// --- Domain services ---
	rpSvc := newReviewPeriodService(repos, cfg, log)
	competencySvc := newCompetencyService(repos, cfg, log)
	staffMgtSvc := newStaffManagementService(repos, cfg, log, userMgr)
	erpSvc := newErpEmployeeService(repos, cfg, log)
	perfSvc := newPerformanceManagementService(repos, cfg, log, rpSvc, erpSvc, gsSvc, ucSvc)

	// Grievance depends on several other services (mirrors .NET DI graph).
	grievanceSvc := newGrievanceManagementService(repos, cfg, log,
		erpSvc,   // ErpEmployeeService
		gsSvc,    // GlobalSettingService
		ucSvc,    // UserContextService
		emailSvc, // EmailService
		rpSvc,    // ReviewPeriodService
	)

	return &Container{
		Performance:   perfSvc,
		PmsSetup:      pmsSetupSvc,
		ReviewPeriod:  rpSvc,
		Grievance:     grievanceSvc,
		RoleMgt:       newRoleManagementService(repos, cfg, log),
		StaffMgt:      staffMgtSvc,
		Competency:    competencySvc,
		Organogram:    newOrganogramService(repos, cfg, log),
		ErpEmployee:   erpSvc,
		GlobalSetting: gsSvc,
		Auth:          authSvc,
		Email:         emailSvc,
		FileStorage:   fsSvc,
		Notification:  newNotificationService(emailSvc, cfg, log),
		Encryption:    encSvc,
		AD:            adSvc,
		UserContext:    ucSvc,
		PasswordGen:   pwGen,
	}
}
