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

// PaginationRequest carries standard pagination parameters for list endpoints.
type PaginationRequest struct {
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	SearchTerm string `json:"search_term,omitempty"`
}

// PaginatedResponse is a generic paginated response wrapper.
type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// FileUploadDto carries file data for upload endpoints.
type FileUploadDto struct {
	FileName    string `json:"file_name"`
	FileData    string `json:"file_data"`
	ContentType string `json:"content_type"`
}

// ApproveRejectRequestVm is the request for batch-approving or rejecting records.
type ApproveRejectRequestVm struct {
	EntityType      string   `json:"entity_type"`
	RecordIds       []string `json:"record_ids"`
	RejectionReason string   `json:"rejection_reason,omitempty"`
}

// ApproveRejectRequestSingleVm is the request for approving or rejecting a single record.
type ApproveRejectRequestSingleVm struct {
	EntityType      string `json:"entity_type"`
	RecordId        string `json:"record_id"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

// CommonExceptionResponse represents a standard error response from the API.
type CommonExceptionResponse struct {
	IsSuccess  bool   `json:"is_success"`
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	ActionCall string `json:"action_call"`
}

// AssignedRequestModel represents an assigned feedback request for SLA tracking.
type AssignedRequestModel struct {
	RequestID          string `json:"request_id"          validate:"required"`
	ReviewPeriodID     string `json:"review_period_id"    validate:"required"`
	ReferenceID        string `json:"reference_id"        validate:"required"`
	FeedBackRequestType int   `json:"feedback_request_type" validate:"required"`
	Requestor          string `json:"requestor,omitempty"`
}

// ResetPasswordRequestModel is the request body for resetting a user password.
type ResetPasswordRequestModel struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IPAddress  string `json:"ip_address,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
}
