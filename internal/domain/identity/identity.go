package identity

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
)

// ApplicationUser mirrors .NET ApplicationUser (extends IdentityUser).
type ApplicationUser struct {
	ID            string     `json:"id"             gorm:"column:id;primaryKey;size:450"`
	UserName      string     `json:"user_name"      gorm:"column:user_name;uniqueIndex;size:256"`
	NormalizedUserName string `json:"-"             gorm:"column:normalized_user_name;uniqueIndex;size:256"`
	Email         string     `json:"email"          gorm:"column:email;size:256"`
	NormalizedEmail string   `json:"-"              gorm:"column:normalized_email;size:256;index"`
	EmailConfirmed bool      `json:"email_confirmed" gorm:"column:email_confirmed;default:false"`
	PasswordHash  string     `json:"-"              gorm:"column:password_hash"`
	SecurityStamp string     `json:"-"              gorm:"column:security_stamp"`
	ConcurrencyStamp string  `json:"-"              gorm:"column:concurrency_stamp"`
	PhoneNumber   string     `json:"phone_number"   gorm:"column:phone_number"`
	PhoneNumberConfirmed bool `json:"-"             gorm:"column:phone_number_confirmed;default:false"`
	TwoFactorEnabled bool   `json:"-"              gorm:"column:two_factor_enabled;default:false"`
	LockoutEnd    *time.Time `json:"-"              gorm:"column:lockout_end"`
	LockoutEnabled bool     `json:"-"              gorm:"column:lockout_enabled;default:true"`
	AccessFailedCount int   `json:"-"              gorm:"column:access_failed_count;default:0"`
	// Custom fields
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name"  gorm:"column:last_name"`
	IsActive  bool   `json:"is_active"  gorm:"column:is_active;default:true"`
}

func (ApplicationUser) TableName() string { return "CoreSchema.asp_net_users" }

// FullName returns the user's display name.
func (u *ApplicationUser) FullName() string {
	return u.FirstName + " " + u.LastName
}

// ApplicationRole mirrors .NET ApplicationRole (extends IdentityRole).
type ApplicationRole struct {
	ID               string `json:"id"                gorm:"column:id;primaryKey;size:450"`
	Name             string `json:"name"              gorm:"column:name;size:256"`
	NormalizedName   string `json:"-"                 gorm:"column:normalized_name;uniqueIndex;size:256"`
	ConcurrencyStamp string `json:"-"                 gorm:"column:concurrency_stamp"`
	RolePermissions  []RolePermission `json:"role_permissions" gorm:"foreignKey:RoleID"`
}

func (ApplicationRole) TableName() string { return "CoreSchema.asp_net_roles" }

// Permission represents an application permission.
type Permission struct {
	PermissionID int    `json:"permission_id" gorm:"column:permission_id;primaryKey;autoIncrement"`
	Name         string `json:"name"          gorm:"column:name;not null"`
	Description  string `json:"description"   gorm:"column:description"`
	domain.BaseAudit
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:PermissionID"`
}

func (Permission) TableName() string { return "CoreSchema.permissions" }

// RolePermission maps a permission to a role.
type RolePermission struct {
	RolePermissionID int    `json:"role_permission_id" gorm:"column:role_permission_id;primaryKey;autoIncrement"`
	PermissionID     int    `json:"permission_id"      gorm:"column:permission_id;not null"`
	RoleID           string `json:"role_id"            gorm:"column:role_id;not null;size:450"`
	domain.BaseAudit
	Role       *ApplicationRole `json:"role"       gorm:"foreignKey:RoleID"`
	Permission *Permission      `json:"permission" gorm:"foreignKey:PermissionID"`
}

func (RolePermission) TableName() string { return "CoreSchema.role_permissions" }

// BankYear represents a fiscal year.
type BankYear struct {
	BankYearID int    `json:"bank_year_id" gorm:"column:bank_year_id;primaryKey;autoIncrement"`
	YearName   string `json:"year_name"    gorm:"column:year_name;uniqueIndex;size:10;not null"`
	domain.BaseAudit
}

func (BankYear) TableName() string { return "CoreSchema.bank_years" }
