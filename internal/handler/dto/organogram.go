package dto

// ===========================================================================
// Organogram Request VMs
// ===========================================================================

// DirectorateVm represents a directorate in the organisational structure.
type DirectorateVm struct {
	DirectorateID   string `json:"directorate_id"`
	Name            string `json:"name"`
	DirectorStaffID string `json:"director_staff_id"`
}

// DepartmentVm represents a department under a directorate.
type DepartmentVm struct {
	DepartmentID  string `json:"department_id"`
	Name          string `json:"name"`
	HeadStaffID   string `json:"head_staff_id"`
	DirectorateID string `json:"directorate_id"`
}

// DivisionVm represents a division under a department.
type DivisionVm struct {
	DivisionID   string `json:"division_id"`
	Name         string `json:"name"`
	HeadStaffID  string `json:"head_staff_id"`
	DepartmentID string `json:"department_id"`
}

// OfficeVm represents an office under a division.
type OfficeVm struct {
	OfficeID    string `json:"office_id"`
	Name        string `json:"name"`
	HeadStaffID string `json:"head_staff_id"`
	DivisionID  string `json:"division_id"`
}

// ===========================================================================
// Organogram Response VMs & Data Structs
// ===========================================================================

// DirectorateData holds directorate data for responses.
type DirectorateData struct {
	DirectorateID   string `json:"directorate_id"`
	Name            string `json:"name"`
	DirectorStaffID string `json:"director_staff_id"`
	DirectorName    string `json:"director_name"`
	RecordStatus    string `json:"record_status"`
	DepartmentCount int    `json:"department_count"`
}

// DirectorateResponseVm wraps a single directorate in a standard response.
type DirectorateResponseVm struct {
	BaseAPIResponse
	Directorate DirectorateData `json:"directorate"`
}

// DirectorateListResponseVm wraps a list of directorates.
type DirectorateListResponseVm struct {
	GenericListResponseVm
	Directorates []DirectorateData `json:"directorates"`
}

// DepartmentData holds department data for responses.
type DepartmentData struct {
	DepartmentID  string `json:"department_id"`
	Name          string `json:"name"`
	HeadStaffID   string `json:"head_staff_id"`
	HeadName      string `json:"head_name"`
	DirectorateID string `json:"directorate_id"`
	DirectorateName string `json:"directorate_name"`
	RecordStatus  string `json:"record_status"`
	DivisionCount int    `json:"division_count"`
}

// DepartmentResponseVm wraps a single department in a standard response.
type DepartmentResponseVm struct {
	BaseAPIResponse
	Department DepartmentData `json:"department"`
}

// DepartmentListResponseVm wraps a list of departments.
type DepartmentListResponseVm struct {
	GenericListResponseVm
	Departments []DepartmentData `json:"departments"`
}

// DivisionData holds division data for responses.
type DivisionData struct {
	DivisionID     string `json:"division_id"`
	Name           string `json:"name"`
	HeadStaffID    string `json:"head_staff_id"`
	HeadName       string `json:"head_name"`
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	RecordStatus   string `json:"record_status"`
	OfficeCount    int    `json:"office_count"`
}

// DivisionResponseVm wraps a single division in a standard response.
type DivisionResponseVm struct {
	BaseAPIResponse
	Division DivisionData `json:"division"`
}

// DivisionListResponseVm wraps a list of divisions.
type DivisionListResponseVm struct {
	GenericListResponseVm
	Divisions []DivisionData `json:"divisions"`
}

// OfficeData holds office data for responses.
type OfficeData struct {
	OfficeID     string `json:"office_id"`
	Name         string `json:"name"`
	HeadStaffID  string `json:"head_staff_id"`
	HeadName     string `json:"head_name"`
	DivisionID   string `json:"division_id"`
	DivisionName string `json:"division_name"`
	RecordStatus string `json:"record_status"`
	StaffCount   int    `json:"staff_count"`
}

// OfficeResponseVm wraps a single office in a standard response.
type OfficeResponseVm struct {
	BaseAPIResponse
	Office OfficeData `json:"office"`
}

// OfficeListResponseVm wraps a list of offices.
type OfficeListResponseVm struct {
	GenericListResponseVm
	Offices []OfficeData `json:"offices"`
}

// ===========================================================================
// Save (Create/Update) View Models
// ===========================================================================

// SaveDirectorateVm is the request body for creating or updating a directorate.
type SaveDirectorateVm struct {
	DirectorateID   string `json:"directorate_id,omitempty"`
	Name            string `json:"name"`
	DirectorateCode string `json:"directorate_code"`
	DirectorStaffID string `json:"director_staff_id"`
}

// SaveDepartmentVm is the request body for creating or updating a department.
type SaveDepartmentVm struct {
	DepartmentID   string `json:"department_id,omitempty"`
	Name           string `json:"name"`
	DepartmentCode string `json:"department_code"`
	DirectorateID  string `json:"directorate_id"`
	HeadStaffID    string `json:"head_staff_id"`
	IsBranch       bool   `json:"is_branch"`
}

// SaveDivisionVm is the request body for creating or updating a division.
type SaveDivisionVm struct {
	DivisionID   string `json:"division_id,omitempty"`
	Name         string `json:"name"`
	DivisionCode string `json:"division_code"`
	DepartmentID string `json:"department_id"`
	HeadStaffID  string `json:"head_staff_id"`
}

// SaveOfficeVm is the request body for creating or updating an office.
type SaveOfficeVm struct {
	OfficeID    string `json:"office_id,omitempty"`
	Name        string `json:"name"`
	OfficeCode  string `json:"office_code"`
	DivisionID  string `json:"division_id"`
	HeadStaffID string `json:"head_staff_id"`
}
