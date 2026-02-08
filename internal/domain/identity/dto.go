package identity

// PermissionVm is the view model for a single permission.
type PermissionVm struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RolePermissionVm represents a role together with its assigned permissions.
type RolePermissionVm struct {
	RoleId      string         `json:"roleId"`
	RoleName    string         `json:"roleName"`
	Permissions []PermissionVm `json:"permissions"`
}

// GetRolePermissionVm is used to display all permissions alongside a specific role's assignments.
type GetRolePermissionVm struct {
	AllPermissions    []PermissionVm    `json:"allPermissions"`
	RolesAndPermissions *RolePermissionVm `json:"rolesAndPermissions"`
}

// AddPermissionToRoleVm is the payload for assigning a permission to a role.
type AddPermissionToRoleVm struct {
	RoleId       string `json:"roleId"`
	PermissionId int    `json:"permissionId"`
}

// RoleSelectList is a lightweight role reference used in dropdowns.
type RoleSelectList struct {
	RoleId   string `json:"roleId"`
	RoleName string `json:"roleName"`
}

// AddStaffToRoleVm is the payload for assigning a staff member to a role.
type AddStaffToRoleVm struct {
	UserId    string `json:"userId"    validate:"required"`
	RoleId    string `json:"roleId"`
	RoleName  string `json:"roleName"  validate:"required"`
	StaffName string `json:"staffName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// RoleVm is the view model for creating or displaying a role.
// RoleName must contain only letters, underscores, or hyphens.
type RoleVm struct {
	ID       string `json:"id"`
	RoleName string `json:"roleName" validate:"required,alphaunicode_underscore_hyphen"`
}

// AssignRoleToUserVm is the payload for assigning one or more roles to a user.
type AssignRoleToUserVm struct {
	UserId    string   `json:"userId"`
	RoleNames []string `json:"roleNames"`
}

// RemoveRoleFromUserVm is the payload for removing a role from a user.
type RemoveRoleFromUserVm struct {
	UserId   string `json:"userId"`
	RoleName string `json:"roleName"`
}

// PersonVm is the base view model for a person with name and contact fields.
// Embed this struct in concrete types that represent people.
type PersonVm struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	MiddleName  string `json:"middleName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Passport    []byte `json:"passport"`
}

// FullName returns the person's display name (first + last).
func (p *PersonVm) FullName() string {
	if p.MiddleName != "" {
		return p.FirstName + " " + p.MiddleName + " " + p.LastName
	}
	return p.FirstName + " " + p.LastName
}

// BankYearVm is the view model for a fiscal/bank year.
type BankYearVm struct {
	BankYearId int    `json:"bankYearId"`
	YearName   string `json:"yearName"`
	IsActive   bool   `json:"isActive"`
}

// BaseAPIResponse is a generic API envelope returned by endpoints.
type BaseAPIResponse struct {
	IsSuccess    bool     `json:"isSuccess"`
	Message      string   `json:"message"`
	TotalRecords int      `json:"totalRecords"`
	Data         any      `json:"data"`
	Errors       []string `json:"errors"`
}
