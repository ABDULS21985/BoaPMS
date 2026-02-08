package auth

// Role constants mirror .NET RoleName static class.
const (
	RoleSuperAdmin          = "SuperAdmin"
	RoleAdmin               = "Admin"
	RoleStaff               = "Staff"
	RoleHeadOfOffice        = "HeadOfOffice"
	RoleHeadOfDivision      = "HeadOfDivision"
	RoleHeadOfDepartment    = "HeadOfDepartment"
	RoleSupervisor          = "Supervisor"
	RoleHRD                 = "HRD"
	RoleReviewer            = "Reviewer"
	RoleApprover            = "Approver"
	RoleDeputyDirector      = "DeputyDirector"
	RoleDirector            = "Director"
	RoleHrAdmin             = "HrAdmin"
	RoleHrApprover          = "HrApprover"
	RoleHrReportAdmin       = "HrReportAdmin"
	RoleGeneralReportAdmin  = "GeneralReportAdmin"
	RoleSmd                 = "Smd"
	RoleSmdApprover         = "SmdApprover"
	RoleSmdOutcomeEvaluator = "SmdOutcomeEvaluator"
	RoleSecurityAdmin       = "SecurityAdmin"
)

// AllRoles returns all defined role names.
func AllRoles() []string {
	return []string{
		RoleSuperAdmin,
		RoleAdmin,
		RoleStaff,
		RoleHeadOfOffice,
		RoleHeadOfDivision,
		RoleHeadOfDepartment,
		RoleSupervisor,
		RoleHRD,
		RoleReviewer,
		RoleApprover,
		RoleDeputyDirector,
		RoleDirector,
		RoleHrAdmin,
		RoleHrApprover,
		RoleHrReportAdmin,
		RoleGeneralReportAdmin,
		RoleSmd,
		RoleSmdApprover,
		RoleSmdOutcomeEvaluator,
		RoleSecurityAdmin,
	}
}

// GetRoleList returns the subset of roles typically displayed in admin UIs.
func GetRoleList() []string {
	return []string{
		RoleSuperAdmin,
		RoleAdmin,
		RoleStaff,
		RoleHeadOfOffice,
		RoleHeadOfDivision,
		RoleHeadOfDepartment,
		RoleSupervisor,
		RoleHRD,
		RoleReviewer,
		RoleApprover,
		RoleDeputyDirector,
		RoleDirector,
		RoleHrAdmin,
		RoleHrApprover,
		RoleHrReportAdmin,
		RoleGeneralReportAdmin,
		RoleSmd,
		RoleSmdApprover,
		RoleSmdOutcomeEvaluator,
		RoleSecurityAdmin,
	}
}

// GetStaffRoles returns the roles that apply to regular staff members
// (excluding administrative and security roles).
func GetStaffRoles() []string {
	return []string{
		RoleStaff,
		RoleHeadOfOffice,
		RoleHeadOfDivision,
		RoleHeadOfDepartment,
		RoleSupervisor,
		RoleReviewer,
		RoleApprover,
		RoleDeputyDirector,
		RoleDirector,
	}
}

// GlobalSettingKeys mirrors the .NET GlobalSetting key constants.
const (
	SettingEnableADAuth         = "ENABLE_AD_AUTHENTICATION"
	SettingTokenExpiryMinutes   = "TOKEN_EXPIRY_IN_MINUTES"
	SettingDefaultPassword      = "DEFAULT_PASSWORD"
	SettingMaxFailedAttempts     = "MAX_FAILED_ACCESS_ATTEMPTS"
	SettingLockoutDuration      = "LOCKOUT_DURATION_MINUTES"
)
