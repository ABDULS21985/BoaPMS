package organogram

// DirectorateVm is the API representation of a Directorate.
type DirectorateVm struct {
	DirectorateID   int    `json:"directorateId"`
	DirectorateName string `json:"directorateName"`
	DirectorateCode string `json:"directorateCode"`
	IsActive        bool   `json:"isActive"`
}

// DepartmentVm is the API representation of a Department.
type DepartmentVm struct {
	DepartmentID   int    `json:"departmentId"`
	DirectorateID  *int   `json:"directorateId"`
	DepartmentName string `json:"departmentName"`
	DepartmentCode string `json:"departmentCode"`
	DirectorateName string `json:"directorateName"`
	IsBranch       bool   `json:"isBranch"`
	IsActive       bool   `json:"isActive"`
}

// DivisionVm is the API representation of a Division.
type DivisionVm struct {
	DivisionID     int    `json:"divisionId"`
	DepartmentID   int    `json:"departmentId"`
	DivisionName   string `json:"divisionName"`
	DivisionCode   string `json:"divisionCode"`
	DepartmentName string `json:"departmentName"`
	IsActive       bool   `json:"isActive"`
}

// OfficeVm is the API representation of an Office.
type OfficeVm struct {
	OfficeID     int    `json:"officeId"`
	DivisionID   int    `json:"divisionId"`
	OfficeName   string `json:"officeName"`
	OfficeCode   string `json:"officeCode"`
	DivisionName string `json:"divisionName"`
	IsActive     bool   `json:"isActive"`
}
