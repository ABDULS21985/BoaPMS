package sas

import "time"

// AbsenceMode maps to the XXCBN_SAS_AbsenceMode table.
type AbsenceMode struct {
	AbsenceID       int        `json:"absenceId"       gorm:"column:AbsenceId;primaryKey"`
	AbsenceModeName *string    `json:"absenceModeName" gorm:"column:AbsenceModeName;size:50"`
	CreateDate      *time.Time `json:"createDate"      gorm:"column:CreateDate"`
	LastUpdateDate  *time.Time `json:"lastUpdateDate"  gorm:"column:LastUpdateDate"`
}

// TableName returns the SQL Server table name for AbsenceMode.
func (AbsenceMode) TableName() string { return "dbo.XXCBN_SAS_AbsenceMode" }

// StaffLunchAttendance maps to the XXCBN_SAS_StaffLunchAttendance table.
type StaffLunchAttendance struct {
	LunchID                  int        `json:"lunchId"                  gorm:"column:LunchId;primaryKey"`
	EmployeeNumber           string     `json:"employeeNumber"           gorm:"column:EmployeeNumber;size:50;not null"`
	FirstName                *string    `json:"firstName"                gorm:"column:FirstName;size:200"`
	LastName                 *string    `json:"lastName"                 gorm:"column:LastName;size:200"`
	Office                   string     `json:"office"                   gorm:"column:Office;size:200;not null"`
	OfficeID                 *int       `json:"officeId"                 gorm:"column:OfficeId"`
	Division                 string     `json:"division"                 gorm:"column:Division;size:200;not null"`
	DivisionID               *int       `json:"divisionId"               gorm:"column:DivisionId"`
	Department               string     `json:"department"               gorm:"column:Department;size:200;not null"`
	DepartmentID             *int       `json:"departmentId"             gorm:"column:DepartmentId"`
	Grade                    string     `json:"grade"                    gorm:"column:Grade;size:10;not null"`
	TurnstileClockedIn       *string    `json:"turnstileClockedIn"       gorm:"column:TurnstileClockedIn;size:5"`
	YearRef                  *int       `json:"yearRef"                  gorm:"column:YearRef"`
	WeekRef                  *int       `json:"weekRef"                  gorm:"column:WeekRef"`
	DayRef                   *int       `json:"dayRef"                   gorm:"column:DayRef"`
	Status                   *string    `json:"status"                   gorm:"column:Status;size:50"`
	PreviousAttendanceStatus *string    `json:"previousAttendanceStatus" gorm:"column:PreviousAttendanceStatus;size:50"`
	AttendanceStatus         int        `json:"attendanceStatus"         gorm:"column:AttendanceStatus"`
	CreateDate               *time.Time `json:"createDate"               gorm:"column:CreateDate"`
	LastUpdateBy             *string    `json:"lastUpdateBy"             gorm:"column:LastUpdateBy;size:10"`
	LastUpdateDate           *time.Time `json:"lastUpdateDate"           gorm:"column:LastUpdateDate"`
	CreatedBy                *string    `json:"createdBy"                gorm:"column:CreatedBy;size:50"`
	RejectedBy               *string    `json:"rejectedBy"               gorm:"column:RejectedBy;size:10"`
	VerifiedBy               *string    `json:"verifiedBy"               gorm:"column:VerifiedBy;size:10"`
	ApprovedBy               *string    `json:"approvedBy"               gorm:"column:ApprovedBy;size:10"`
	ApprovedDate             *time.Time `json:"approvedDate"             gorm:"column:ApprovedDate"`
	Remarks                  *string    `json:"remarks"                  gorm:"column:Remarks;size:50"`
	Location                 *string    `json:"location"                 gorm:"column:Location;size:100"`
	LocationID               *int       `json:"locationId"               gorm:"column:LocationId"`
}

// TableName returns the SQL Server table name for StaffLunchAttendance.
func (StaffLunchAttendance) TableName() string { return "dbo.XXCBN_SAS_StaffLunchAttendance" }
