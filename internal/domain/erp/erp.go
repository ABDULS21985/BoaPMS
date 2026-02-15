package erp

import "time"

// EmployeeDetails is a read-only model from the ERP SQL Server database.
// The struct maps the dbo.EmployeeDetails view used by ReviewAgentService
// for 360-degree review processing (subordinate/peer/superior selection).
type EmployeeDetails struct {
	EmployeeNumber  string     `json:"employee_number"   db:"EmployeeNumber"   gorm:"column:EmployeeNumber;primaryKey"`
	FirstName       string     `json:"first_name"        db:"FirstName"        gorm:"column:FirstName"`
	LastName        string     `json:"last_name"         db:"LastName"         gorm:"column:LastName"`
	FullName        string     `json:"full_name"         db:"FullName"         gorm:"column:FullName"`
	Email           string     `json:"email"             db:"EmailAddress"     gorm:"column:EmailAddress"`
	Grade           string     `json:"grade"             db:"Grade"            gorm:"column:Grade"`
	Department      string     `json:"department"        db:"DepartmentName"   gorm:"column:DepartmentName"`
	DepartmentID    *int       `json:"department_id"     db:"DepartmentId"     gorm:"column:DepartmentId"`
	Division        string     `json:"division"          db:"DivisionName"     gorm:"column:DivisionName"`
	DivisionID      *int       `json:"division_id"       db:"DivisionId"       gorm:"column:DivisionId"`
	Office          string     `json:"office"            db:"OfficeName"       gorm:"column:OfficeName"`
	OfficeID        *int       `json:"office_id"         db:"OfficeId"         gorm:"column:OfficeId"`
	LocationID      *int       `json:"location_id"       db:"LocationId"       gorm:"column:LocationId"`
	SupervisorID    string     `json:"supervisor_id"     db:"SupervisorId"     gorm:"column:SupervisorId"`
	SupervisorName  string     `json:"supervisor_name"   db:"SupervisorName"   gorm:"column:SupervisorName"`
	HeadOfOfficeID  string     `json:"head_of_office_id" db:"HeadOfOfficeId"   gorm:"column:HeadOfOfficeId"`
	HeadOfDivID     string     `json:"head_of_div_id"    db:"HeadOfDivId"      gorm:"column:HeadOfDivId"`
	HeadOfDeptID    string     `json:"head_of_dept_id"   db:"HeadOfDeptId"     gorm:"column:HeadOfDeptId"`
	JobName         string     `json:"job_name"          db:"JobName"          gorm:"column:JobName"`
	JobTitle        string     `json:"job_title"         db:"JobTitle"         gorm:"column:JobTitle"`
	PersonTypeID    int        `json:"person_type_id"    db:"PersonTypeId"     gorm:"column:PersonTypeId"`
	HireDate        *time.Time `json:"hire_date"         db:"HireDate"         gorm:"column:HireDate"`
}

func (EmployeeDetails) TableName() string { return "dbo.EmployeeDetails" }

// ErpLocationDetail provides branch/location info.
type ErpLocationDetail struct {
	LocationID   int    `json:"location_id"   db:"LOCATION_ID"   gorm:"column:LOCATION_ID;primaryKey"`
	LocationName string `json:"location_name" db:"LOCATION_NAME" gorm:"column:LOCATION_NAME"`
	LocationCode string `json:"location_code" db:"LOCATION_CODE" gorm:"column:LOCATION_CODE"`
}

func (ErpLocationDetail) TableName() string { return "dbo.ErpLocationDetails" }

// StaffIDMaskDetails is a read-only model from the StaffIDMask SQL Server database.
type StaffIDMaskDetails struct {
	ID           int    `json:"id"            db:"Id"            gorm:"column:Id;primaryKey"`
	StaffID      string `json:"staff_id"      db:"StaffId"       gorm:"column:StaffId"`
	StaffName    string `json:"staff_name"    db:"StaffName"     gorm:"column:StaffName"`
	PictureData  []byte `json:"-"             db:"PictureData"   gorm:"column:PictureData;type:bytea"`
	BloodGroup   string `json:"blood_group"   db:"BloodGroup"    gorm:"column:BloodGroup"`
}

func (StaffIDMaskDetails) TableName() string { return "dbo.StaffIDMaskDetails" }

// PublicHolidayData represents public holiday entries from ERP.
// Maps to dbo.HOLIDAYS_T24 table. Mirrors .NET HOLIDAYS_T24 entity.
type PublicHolidayData struct {
	UniqueID     int        `json:"unique_id"     db:"unique_id"     gorm:"column:unique_id;primaryKey"`
	HName        *string    `json:"h_name"        db:"HName"         gorm:"column:HName"`
	HDate        *time.Time `json:"h_date"        db:"HDate"         gorm:"column:HDate"`
	HType        *string    `json:"h_type"        db:"HType"         gorm:"column:HType"`
	EventStatus  *string    `json:"event_status"  db:"event_status"  gorm:"column:event_status"`
	EventKey     *string    `json:"event_key"     db:"event_key"     gorm:"column:event_key"`
	CreationDate *time.Time `json:"creation_date" db:"creation_date" gorm:"column:creation_date"`
}

func (PublicHolidayData) TableName() string { return "dbo.HOLIDAYS_T24" }

// VacationRuleData represents leave/vacation delegation rules from ERP.
// Maps to dbo.VACATIONSRULE_DATA table. Mirrors .NET VACATIONSRULE_DATA entity.
type VacationRuleData struct {
	UniqueID     int        `json:"unique_id"     db:"unique_id"     gorm:"column:unique_id;primaryKey"`
	RuleID       int        `json:"rule_id"       db:"rule_id"       gorm:"column:rule_id"`
	RuleOwner    *string    `json:"rule_owner"    db:"rule_owner"    gorm:"column:rule_owner"`
	Action       *string    `json:"action"        db:"action"        gorm:"column:action"`
	BeginDate    time.Time  `json:"begin_date"    db:"begin_date"    gorm:"column:begin_date"`
	EndDate      *time.Time `json:"end_date"      db:"end_date"      gorm:"column:end_date"`
	MessageType  *string    `json:"message_type"  db:"message_type"  gorm:"column:message_type"`
	MessageName  *string    `json:"message_name"  db:"message_name"  gorm:"column:message_name"`
	AssignedTo   *string    `json:"assigned_to"   db:"assigned_to"   gorm:"column:assigned_to"`
	RuleComment  *string    `json:"rule_comment"  db:"rule_comment"  gorm:"column:rule_comment"`
	EventStatus  *string    `json:"event_status"  db:"event_status"  gorm:"column:event_status"`
	EventKey     *string    `json:"event_key"     db:"event_key"     gorm:"column:event_key"`
	CreationDate *time.Time `json:"creation_date" db:"creation_date" gorm:"column:creation_date"`
}

func (VacationRuleData) TableName() string { return "dbo.VACATIONSRULE_DATA" }

// EmailObject represents email records from the email service database.
// Mirrors the .NET CompetencyApp.Models.Core.EmailObjects entity.
type EmailObject struct {
	ID               int        `json:"id"               db:"Id"              gorm:"column:Id;primaryKey;autoIncrement"`
	From             string     `json:"from"             db:"From"            gorm:"column:From"`
	To               string     `json:"to"               db:"To"              gorm:"column:To"`
	CC               string     `json:"cc"               db:"CC"              gorm:"column:CC"`
	BCC              string     `json:"bcc"              db:"BCC"             gorm:"column:BCC"`
	Subject          string     `json:"subject"          db:"Subject"         gorm:"column:Subject"`
	Body             string     `json:"body"             db:"Body"            gorm:"column:Body;type:text"`
	Status           string     `json:"status"           db:"Status"          gorm:"column:Status"`
	NoOfRetry        int        `json:"noOfRetry"        db:"NoOfRetry"       gorm:"column:NoOfRetry;default:0"`
	ExpectedSendDate *time.Time `json:"expectedSendDate" db:"ExpectedSendDate" gorm:"column:ExpectedSendDate"`
	ActualSendDate   *time.Time `json:"actualSendDate"   db:"ActualSendDate"  gorm:"column:ActualSendDate"`
	Action           string     `json:"action"           db:"Action"          gorm:"column:Action"`
	AppSource        string     `json:"appSource"        db:"AppSource"       gorm:"column:AppSource"`
	CreatedBy        string     `json:"createdBy"        db:"CreatedBy"       gorm:"column:CreatedBy"`
	DateCreated      *time.Time `json:"dateCreated"      db:"DateCreated"     gorm:"column:DateCreated"`
	LastUpdatedBy    string     `json:"lastUpdatedBy"    db:"LastUpdatedBy"   gorm:"column:LastUpdatedBy"`
	LastUpdatedDate  *time.Time `json:"lastUpdatedDate"  db:"LastUpdatedDate" gorm:"column:LastUpdatedDate"`
	EmailGUID        string     `json:"emailGuid"        db:"EmailGuid"       gorm:"column:EmailGuid"`
}

func (EmailObject) TableName() string { return "dbo.EmailObjects" }
