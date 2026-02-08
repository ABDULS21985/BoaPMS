package organogram

import (
	"github.com/enterprise-pms/pms-api/internal/domain"
)

// Directorate is the top level of the organisational hierarchy.
type Directorate struct {
	DirectorateID   int    `json:"directorate_id"   gorm:"column:directorate_id;primaryKey;autoIncrement"`
	DirectorateName string `json:"directorate_name" gorm:"column:directorate_name;not null"`
	DirectorateCode string `json:"directorate_code" gorm:"column:directorate_code"`
	domain.BaseAudit
	Departments []Department `json:"departments" gorm:"foreignKey:DirectorateID"`
}

func (Directorate) TableName() string { return "CoreSchema.directorates" }

// Department belongs to a Directorate.
type Department struct {
	DepartmentID   int    `json:"department_id"   gorm:"column:department_id;primaryKey;autoIncrement"`
	DirectorateID  *int   `json:"directorate_id"  gorm:"column:directorate_id"`
	DepartmentName string `json:"department_name" gorm:"column:department_name;not null"`
	DepartmentCode string `json:"department_code" gorm:"column:department_code"`
	IsBranch       bool   `json:"is_branch"       gorm:"column:is_branch;default:false"`
	domain.BaseAudit
	Directorate *Directorate `json:"directorate" gorm:"foreignKey:DirectorateID"`
	Divisions   []Division   `json:"divisions"   gorm:"foreignKey:DepartmentID"`
}

func (Department) TableName() string { return "CoreSchema.departments" }

// Division belongs to a Department.
type Division struct {
	DivisionID   int    `json:"division_id"   gorm:"column:division_id;primaryKey;autoIncrement"`
	DepartmentID int    `json:"department_id" gorm:"column:department_id;not null"`
	DivisionName string `json:"division_name" gorm:"column:division_name;not null"`
	DivisionCode string `json:"division_code" gorm:"column:division_code"`
	domain.BaseAudit
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Offices    []Office    `json:"offices"    gorm:"foreignKey:DivisionID"`
}

func (Division) TableName() string { return "CoreSchema.divisions" }

// Office is the lowest organisational unit, belongs to a Division.
type Office struct {
	OfficeID   int    `json:"office_id"   gorm:"column:office_id;primaryKey;autoIncrement"`
	DivisionID int    `json:"division_id" gorm:"column:division_id;not null"`
	OfficeName string `json:"office_name" gorm:"column:office_name;not null"`
	OfficeCode string `json:"office_code" gorm:"column:office_code"`
	domain.BaseAudit
	Division *Division `json:"division" gorm:"foreignKey:DivisionID"`
}

func (Office) TableName() string { return "CoreSchema.offices" }
