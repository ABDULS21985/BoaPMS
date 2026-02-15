package dto

// ===========================================================================
// Staff Request VMs
// ===========================================================================

// AddStaffRequest is the request body for creating a new staff user.
type AddStaffRequest struct {
	ID        string `json:"id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// ===========================================================================
// Role VMs
// ===========================================================================

// RoleVm represents a role in the system.
type RoleVm struct {
	Name string `json:"name"`
}

// AddStaffToRoleVm is the request to assign a staff member to a role.
type AddStaffToRoleVm struct {
	UserID   string `json:"user_id"`
	RoleName string `json:"role_name"`
}

// ===========================================================================
// Permission VMs
// ===========================================================================

// AddPermissionToRoleVm is the request to assign a permission to a role.
type AddPermissionToRoleVm struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

// PermissionVm represents a permission in the system.
type PermissionVm struct {
	PermissionID string `json:"permission_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

// ===========================================================================
// Staff & Role Response VMs
// ===========================================================================

// StaffData holds staff data for responses.
type StaffData struct {
	ID        string `json:"id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// StaffResponseVm wraps a single staff record.
type StaffResponseVm struct {
	BaseAPIResponse
	Staff StaffData `json:"staff"`
}

// StaffListResponseVm wraps a list of staff records.
type StaffListResponseVm struct {
	GenericListResponseVm
	Staff []StaffData `json:"staff"`
}

// RoleData holds role data for responses.
type RoleData struct {
	RoleID      string           `json:"role_id"`
	Name        string           `json:"name"`
	Permissions []PermissionVm   `json:"permissions,omitempty"`
}

// RoleResponseVm wraps a single role record.
type RoleResponseVm struct {
	BaseAPIResponse
	Role RoleData `json:"role"`
}

// RoleListResponseVm wraps a list of roles.
type RoleListResponseVm struct {
	GenericListResponseVm
	Roles []RoleData `json:"roles"`
}

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// StaffVm represents a staff member with their assigned roles.
type StaffVm struct {
	ID        string   `json:"id"`
	UserName  string   `json:"user_name"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsActive  bool     `json:"is_active"`
	Roles     []string `json:"roles,omitempty"`
}

// AddStaffVm is the request body for creating a new staff user with all required fields.
type AddStaffVm struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// StaffRoleVm represents a role with its associated permissions for display.
type StaffRoleVm struct {
	RoleID      string       `json:"role_id"`
	RoleName    string       `json:"role_name"`
	Permissions []PermissionVm `json:"permissions,omitempty"`
}

// AddRoleVm is the request body for creating a new role.
type AddRoleVm struct {
	RoleName string `json:"role_name"`
}

// AssignRoleVm is the request body for assigning one or more roles to a staff member.
type AssignRoleVm struct {
	UserID    string   `json:"user_id"`
	RoleNames []string `json:"role_names"`
}
