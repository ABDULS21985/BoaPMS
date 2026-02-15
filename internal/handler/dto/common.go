package dto

// BaseAPIResponse represents the standard API response wrapper.
type BaseAPIResponse struct {
	IsSuccess bool     `json:"is_success"`
	Message   string   `json:"message"`
	Errors    []string `json:"errors,omitempty"`
}

// ResponseVm extends BaseAPIResponse with an entity identifier.
type ResponseVm struct {
	BaseAPIResponse
	ID string `json:"id"`
}

// GenericListResponseVm extends BaseAPIResponse with a total record count.
type GenericListResponseVm struct {
	BaseAPIResponse
	TotalRecords int `json:"total_records"`
}

// GenericResponseVM is a generic typed response carrying data of type T.
type GenericResponseVM[T any] struct {
	Data         T        `json:"data"`
	TotalRecords int      `json:"total_records"`
	Message      string   `json:"message"`
	Errors       []string `json:"errors,omitempty"`
	IsSuccess    bool     `json:"is_success"`
}

// PaginatedResult holds a page of items together with the total count.
type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
}

// EnumList represents a name/ID pair typically used for dropdown lists.
type EnumList struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// BasePagedData carries pagination parameters for paged queries.
type BasePagedData struct {
	Skip       int `json:"skip"`
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

// ApprovalBase contains the common fields for approval and rejection requests.
type ApprovalBase struct {
	RecordIds  []string `json:"record_ids"`
	ApprovedBy string   `json:"approved_by"`
	Comment    string   `json:"comment,omitempty"`
}

// ApprovalRequestVm represents a request to approve one or more records.
type ApprovalRequestVm struct {
	RecordIds  []string `json:"record_ids"`
	ApprovedBy string   `json:"approved_by"`
}

// RejectionRequestVm represents a request to reject one or more records.
type RejectionRequestVm struct {
	RecordIds  []string `json:"record_ids"`
	RejectedBy string   `json:"rejected_by"`
	Reason     string   `json:"reason"`
}

// EmailRequest carries the data needed to dispatch an email notification.
type EmailRequest struct {
	UserID           string `json:"user_id"`
	Title            string `json:"title"`
	RecordCount      int    `json:"record_count"`
	EmailDescription string `json:"email_description"`
}
