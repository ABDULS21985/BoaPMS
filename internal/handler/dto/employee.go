package dto

import "time"

// ===========================================================================
// Employee ERP Details DTO
// ===========================================================================

// EmployeeErpDetailsDTO holds the employee details sourced from the ERP system.
type EmployeeErpDetailsDTO struct {
	EmployeeNumber string `json:"employee_number"`
	FullName       string `json:"full_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Grade          string `json:"grade"`
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	DivisionID     string `json:"division_id"`
	DivisionName   string `json:"division_name"`
	OfficeID       string `json:"office_id"`
	OfficeName     string `json:"office_name"`
	SupervisorID   string `json:"supervisor_id"`
	SupervisorName string `json:"supervisor_name"`
	PhotoURL       string `json:"photo_url"`
}

// ===========================================================================
// Staff Job Roles VM
// ===========================================================================

// StaffJobRolesVm represents the job role assignment for a staff member.
type StaffJobRolesVm struct {
	StaffJobRoleID string `json:"staff_job_role_id"`
	EmployeeID     string `json:"employee_id"`
	FullName       string `json:"full_name"`
	DepartmentID   string `json:"department_id"`
	DivisionID     string `json:"division_id"`
	OfficeID       string `json:"office_id"`
	SupervisorID   string `json:"supervisor_id"`
	JobRoleID      string `json:"job_role_id"`
	JobRoleName    string `json:"job_role_name"`
}

// ===========================================================================
// Employee Response VMs
// ===========================================================================

// EmployeeResponseVm wraps a single employee record.
type EmployeeResponseVm struct {
	BaseAPIResponse
	Employee EmployeeErpDetailsDTO `json:"employee"`
}

// EmployeeListResponseVm wraps a list of employee records.
type EmployeeListResponseVm struct {
	GenericListResponseVm
	Employees []EmployeeErpDetailsDTO `json:"employees"`
}

// StaffJobRolesResponseVm wraps a single staff job role.
type StaffJobRolesResponseVm struct {
	BaseAPIResponse
	JobRole StaffJobRolesVm `json:"job_role"`
}

// StaffJobRolesListResponseVm wraps a list of staff job roles.
type StaffJobRolesListResponseVm struct {
	GenericListResponseVm
	JobRoles []StaffJobRolesVm `json:"job_roles"`
}
