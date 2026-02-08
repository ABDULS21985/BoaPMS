package performance

import "time"

// BaseAuditVm is the base DTO type carrying audit trail fields.
type BaseAuditVm struct {
	CreatedBy   string     `json:"createdBy"`
	DateCreated *time.Time `json:"dateCreated"`
	IsActive    bool       `json:"isActive"`
	Status      string     `json:"status"`
	DateUpdated *time.Time `json:"dateUpdated"`
	UpdatedBy   string     `json:"updatedBy"`
}

// BaseAPIResponse is the standard API envelope.
type BaseAPIResponse struct {
	HasError bool   `json:"hasError"`
	Message  string `json:"message"`
}

// BasePagedData carries pagination parameters.
type BasePagedData struct {
	Skip       int `json:"skip"`
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

// ===========================================================================
// Common Response / Utility Types (ResponseVm.cs)
// ===========================================================================

// ResponseVm extends BaseAPIResponse with an entity ID.
type ResponseVm struct {
	BaseAPIResponse
	ID string `json:"id"`
}

// ApiErrorResponse is the standard error response envelope.
type ApiErrorResponse struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	TraceID string `json:"traceId"`
}

// EnumList is a simple name/ID pair used for enum dropdown lists.
type EnumList struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// WebAppAPIConfig carries configuration for an external web API.
type WebAppAPIConfig struct {
	BaseURL string `json:"baseUrl"`
	APIKey  string `json:"apiKey"`
	URL     string `json:"url"`
}

// EmailRequest is the payload for triggering an email notification.
type EmailRequest struct {
	UserID           string `json:"userId"`
	Title            string `json:"title"`
	RecordCount      int    `json:"recordCount"`
	EmailDescription string `json:"emailDescription"`
}

// SoaJobRoleVm is the payload for updating an employee's job role via SOA.
type SoaJobRoleVm struct {
	PersonID int    `json:"p_PERSON_ID"`
	JobRole  string `json:"p_JOB_ROLE"`
}

// SoaOutputParameters wraps the SOA output for job role updates.
type SoaOutputParameters struct {
	BaseAPIResponse
	UpdateEmployeeJobRole string `json:"update_EMPLOYEE_JOB_ROLE"`
}

// GenericResponseVm is a generic API response carrying a typed data payload.
type GenericResponseVm struct {
	Data         interface{} `json:"data"`
	TotalRecords int         `json:"totalRecords"`
	Message      string      `json:"message"`
	Errors       []string    `json:"errors"`
	IsSuccess    bool        `json:"isSuccess"`
}
