package erp

import (
	"encoding/base64"
	"fmt"
	"time"
)

// EmployeeErpDetailsDTO is the API representation of employee details from ERP.
type EmployeeErpDetailsDTO struct {
	UserName       string `json:"userName"`
	EmailAddress   string `json:"emailAddress"`
	FirstName      string `json:"firstName"`
	MiddleNames    string `json:"middleNames"`
	LastName       string `json:"lastName"`
	EmployeeNumber string `json:"employeeNumber"`
	JobName        string `json:"jobName"`
	DepartmentName string `json:"departmentName"`
	DivisionName   string `json:"divisionName"`
	HeadOfDivName  string `json:"headOfDivName"`
	OfficeName     string `json:"officeName"`
	SupervisorID   string `json:"supervisorId"`
	HeadOfOfficeID string `json:"headOfOfficeId"`
	HeadOfDivID    string `json:"headOfDivId"`
	HeadOfDeptID   string `json:"headOfDeptId"`
	DepartmentID   *int   `json:"departmentId"`
	OfficeID       int    `json:"officeId"`
	Grade          string `json:"grade"`
	DivisionID     *int   `json:"divisionId"`
	Position       string `json:"position"`
	PersonID       int    `json:"personId"`
}

// FullName returns the formatted full name: "LastName, FirstName MiddleNames".
func (e EmployeeErpDetailsDTO) FullName() string {
	return fmt.Sprintf("%s, %s %s", e.LastName, e.FirstName, e.MiddleNames)
}

// NameInitial returns the first letter of LastName and FirstName combined.
func (e EmployeeErpDetailsDTO) NameInitial() string {
	if len(e.LastName) > 0 && len(e.FirstName) > 0 {
		return string(e.LastName[0]) + string(e.FirstName[0])
	}
	return ""
}

// EmployeeData extends EmployeeErpDetailsDTO with additional ERP status fields.
type EmployeeData struct {
	EmployeeErpDetailsDTO
	Status       string `json:"status"`
	PersonTypeID int    `json:"personTypeId"`
	LocationID   int    `json:"locationId"`
	LocationCode string `json:"locationCode"`
}

// ErpOrganizationVm is the API representation of an ERP organisational unit.
type ErpOrganizationVm struct {
	DepartmentID   *int   `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	DivisionID     *int   `json:"divisionId"`
	DivisionName   string `json:"divisionName"`
	OfficeID       int    `json:"officeId"`
	OfficeName     string `json:"officeName"`
}

// ERPOfficeJobRoleVm is the API representation of an office job role from ERP.
type ERPOfficeJobRoleVm struct {
	DivisionCode   string `json:"divisionCode"`
	JobRoleName    string `json:"jobRoleName"`
	OfficeFullName string `json:"officeFullName"`
	OfficeID       int    `json:"officeId"`
	OfficeName     string `json:"officeName"`
}

// EROJobGradeVm is the API representation of a job grade from ERP.
type EROJobGradeVm struct {
	GradeID   string `json:"gradeId"`
	GradeName string `json:"gradeName"`
}

// StaffIDMaskDetailsDTO is the API representation of staff ID mask details.
type StaffIDMaskDetailsDTO struct {
	StaffIDMaskID  int        `json:"staffIdMaskId"`
	Name           string     `json:"name"`
	EmployeeNumber string     `json:"employeeNumber"`
	CurrentPicture []byte     `json:"currentPicture"`
	NewPicture     []byte     `json:"newPicture"`
	BloodGroup     string     `json:"bloodGroup"`
	ApprovedBy     string     `json:"approvedBy"`
	ApprovalDate   *time.Time `json:"approvalDate"`
	RejectReason   string     `json:"rejectReason"`
	RejectedBy     string     `json:"rejectedby"`
	RejectionDate  time.Time  `json:"rejectionDate"`
	Status         string     `json:"status"`
	CreateDate     time.Time  `json:"createDate"`
	CreatedBy      string     `json:"createdBy"`
	LastUpdatedBy  string     `json:"lastUpdatedBy"`
	LastUpdateDate time.Time  `json:"lastUpdateDate"`
	MessageStatus  bool       `json:"messageStatus"`
}

// CurrentStaffPhoto returns the current picture as a base64 data URI, or empty string if nil.
func (s StaffIDMaskDetailsDTO) CurrentStaffPhoto() string {
	if len(s.CurrentPicture) > 0 {
		return fmt.Sprintf("data:image/gif;base64,%s", base64.StdEncoding.EncodeToString(s.CurrentPicture))
	}
	return ""
}
