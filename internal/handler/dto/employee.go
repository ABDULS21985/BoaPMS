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

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// ---------------------------------------------------------------------------
// Employee Detail / List View Models
// ---------------------------------------------------------------------------

// EmployeeDetailVm represents a detailed view of an employee from the ERP system.
type EmployeeDetailVm struct {
	EmployeeNumber string     `json:"employee_number"`
	FullName       string     `json:"full_name"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `json:"email"`
	Grade          string     `json:"grade"`
	JobName        string     `json:"job_name"`
	JobTitle       string     `json:"job_title"`
	DepartmentID   string     `json:"department_id"`
	DepartmentName string     `json:"department_name"`
	DivisionID     string     `json:"division_id"`
	DivisionName   string     `json:"division_name"`
	OfficeID       string     `json:"office_id"`
	OfficeName     string     `json:"office_name"`
	SupervisorID   string     `json:"supervisor_id"`
	SupervisorName string     `json:"supervisor_name"`
	HeadOfOfficeID string     `json:"head_of_office_id"`
	HeadOfDivID    string     `json:"head_of_div_id"`
	HeadOfDeptID   string     `json:"head_of_dept_id"`
	HireDate       *time.Time `json:"hire_date,omitempty"`
	PhotoURL       string     `json:"photo_url,omitempty"`
}

// EmployeeListVm wraps a list of employee detail view models.
type EmployeeListVm struct {
	GenericListResponseVm
	Employees []EmployeeDetailVm `json:"employees"`
}

// ---------------------------------------------------------------------------
// Subordinate View Model
// ---------------------------------------------------------------------------

// SubordinateVm represents a subordinate employee for supervisory views.
type SubordinateVm struct {
	EmployeeNumber string `json:"employee_number"`
	FullName       string `json:"full_name"`
	Grade          string `json:"grade"`
	JobName        string `json:"job_name"`
	DepartmentName string `json:"department_name"`
	DivisionName   string `json:"division_name"`
	OfficeName     string `json:"office_name"`
}

// ---------------------------------------------------------------------------
// Staff Job Role View Models
// ---------------------------------------------------------------------------

// StaffJobRoleVm represents a staff member's current job role assignment.
type StaffJobRoleVm struct {
	StaffJobRoleID string `json:"staff_job_role_id"`
	EmployeeID     string `json:"employee_id"`
	FullName       string `json:"full_name"`
	DepartmentID   string `json:"department_id"`
	DivisionID     string `json:"division_id"`
	OfficeID       string `json:"office_id"`
	SupervisorID   string `json:"supervisor_id"`
	JobRoleID      string `json:"job_role_id"`
	JobRoleName    string `json:"job_role_name"`
	SoaStatus      bool   `json:"soa_status"`
	SoaResponse    string `json:"soa_response,omitempty"`
}

// UpdateStaffJobRoleVm is the request body for updating a staff member's job role.
type UpdateStaffJobRoleVm struct {
	StaffJobRoleID string `json:"staff_job_role_id"`
	EmployeeID     string `json:"employee_id"`
	JobRoleID      string `json:"job_role_id"`
	OfficeID       string `json:"office_id"`
}

// OfficeJobRoleRequestVm is the request body for managing office-to-job-role assignments.
type OfficeJobRoleRequestVm struct {
	OfficeJobRoleID string `json:"office_job_role_id,omitempty"`
	OfficeID        string `json:"office_id"`
	JobRoleID       string `json:"job_role_id"`
}
