package domain

import (
	"time"

	"gorm.io/gorm"
)

// BaseAudit mirrors the .NET BaseAudit abstract class.
// Embedded in all entities from the legacy CompetencyApp schema.
type BaseAudit struct {
	CreatedBy   string     `json:"created_by"   gorm:"column:created_by;size:75;default:'SYSTEM'"`
	DateCreated time.Time  `json:"date_created" gorm:"column:date_created;autoCreateTime"`
	IsActive    bool       `json:"is_active"    gorm:"column:is_active;default:true"`
	Status      string     `json:"status"       gorm:"column:status;size:25"`
	SoftDeleted bool       `json:"-"            gorm:"column:soft_deleted;default:false;index"`
	DateUpdated *time.Time `json:"date_updated" gorm:"column:date_updated"`
	UpdatedBy   string     `json:"updated_by"   gorm:"column:updated_by"`
}

// BaseWorkFlowData extends BaseAudit with single-level approval workflow.
type BaseWorkFlowData struct {
	BaseAudit
	ApprovedBy      string     `json:"approved_by"       gorm:"column:approved_by"`
	DateApproved    *time.Time `json:"date_approved"     gorm:"column:date_approved"`
	IsApproved      bool       `json:"is_approved"       gorm:"column:is_approved;default:false"`
	IsRejected      bool       `json:"is_rejected"       gorm:"column:is_rejected;default:false"`
	RejectedBy      string     `json:"rejected_by"       gorm:"column:rejected_by"`
	RejectionReason string     `json:"rejection_reason"  gorm:"column:rejection_reason"`
	DateRejected    *time.Time `json:"date_rejected"     gorm:"column:date_rejected"`
}

// HrdWorkFlowData extends BaseWorkFlowData with HRD-level approval.
type HrdWorkFlowData struct {
	BaseWorkFlowData
	HrdApprovedBy      string     `json:"hrd_approved_by"       gorm:"column:hrd_approved_by"`
	HrdDateApproved    *time.Time `json:"hrd_date_approved"     gorm:"column:hrd_date_approved"`
	HrdIsApproved      bool       `json:"hrd_is_approved"       gorm:"column:hrd_is_approved;default:false"`
	HrdIsRejected      bool       `json:"hrd_is_rejected"       gorm:"column:hrd_is_rejected;default:false"`
	HrdRejectedBy      string     `json:"hrd_rejected_by"       gorm:"column:hrd_rejected_by"`
	HrdRejectionReason string     `json:"hrd_rejection_reason"  gorm:"column:hrd_rejection_reason"`
	HrdDateRejected    *time.Time `json:"hrd_date_rejected"     gorm:"column:hrd_date_rejected"`
}

// BaseEntity mirrors the .NET BaseEntity for the PMS schema.
// Uses int auto-increment PK (matching the .NET [Key] pattern).
type BaseEntity struct {
	ID          int        `json:"id"           gorm:"column:id;primaryKey;autoIncrement"`
	RecordStatus string   `json:"record_status" gorm:"column:record_status;default:'Active'"`
	CreatedAt   *time.Time `json:"created_at"   gorm:"column:created_at;autoCreateTime"`
	SoftDeleted bool       `json:"-"            gorm:"column:soft_deleted;default:false;index"`
	Status      string     `json:"status"       gorm:"column:status"`
	UpdatedAt   *time.Time `json:"updated_at"   gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy   string     `json:"created_by"   gorm:"column:created_by;size:100"`
	UpdatedBy   string     `json:"updated_by"   gorm:"column:updated_by;size:100"`
	IsActive    bool       `json:"is_active"    gorm:"column:is_active;default:true"`
}

// BaseWorkFlow extends BaseEntity with single-level approval for PMS entities.
type BaseWorkFlow struct {
	BaseEntity
	ApprovedBy      string     `json:"approved_by"       gorm:"column:approved_by"`
	DateApproved    *time.Time `json:"date_approved"     gorm:"column:date_approved"`
	IsApproved      bool       `json:"is_approved"       gorm:"column:is_approved;default:false"`
	IsRejected      bool       `json:"is_rejected"       gorm:"column:is_rejected;default:false"`
	RejectedBy      string     `json:"rejected_by"       gorm:"column:rejected_by"`
	RejectionReason string     `json:"rejection_reason"  gorm:"column:rejection_reason"`
	DateRejected    *time.Time `json:"date_rejected"     gorm:"column:date_rejected"`
}

// HrdWorkFlow extends BaseWorkFlow with HRD-level approval for PMS entities.
type HrdWorkFlow struct {
	BaseWorkFlow
	HrdApprovedBy      string     `json:"hrd_approved_by"       gorm:"column:hrd_approved_by"`
	HrdDateApproved    *time.Time `json:"hrd_date_approved"     gorm:"column:hrd_date_approved"`
	HrdIsApproved      bool       `json:"hrd_is_approved"       gorm:"column:hrd_is_approved;default:false"`
	HrdIsRejected      bool       `json:"hrd_is_rejected"       gorm:"column:hrd_is_rejected;default:false"`
	HrdRejectedBy      string     `json:"hrd_rejected_by"       gorm:"column:hrd_rejected_by"`
	HrdRejectionReason string     `json:"hrd_rejection_reason"  gorm:"column:hrd_rejection_reason"`
	HrdDateRejected    *time.Time `json:"hrd_date_rejected"     gorm:"column:hrd_date_rejected"`
}

// ObjectiveBase is the shared base for all objective types.
type ObjectiveBase struct {
	BaseWorkFlow
	Name             string `json:"name"               gorm:"column:name;not null"`
	SmdReferenceCode string `json:"smd_reference_code" gorm:"column:smd_reference_code"`
	Description      string `json:"description"        gorm:"column:description"`
	Kpi              string `json:"kpi"                gorm:"column:kpi;not null"`
	Target           string `json:"target"             gorm:"column:target"`
}

// --- GORM Hooks ---

// BeforeCreate sets audit fields on new BaseAudit records.
func (b *BaseAudit) BeforeCreate(tx *gorm.DB) error {
	b.DateCreated = time.Now().UTC()
	b.Status = "CREATE"
	if b.CreatedBy == "" {
		b.CreatedBy = "SYSTEM"
	}
	b.IsActive = true
	return nil
}

// BeforeUpdate sets audit fields on modified BaseAudit records.
func (b *BaseAudit) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now().UTC()
	b.DateUpdated = &now
	if b.SoftDeleted {
		b.Status = "DELETED"
	} else {
		b.Status = "UPDATE"
	}
	return nil
}

// BeforeCreate sets audit fields on new BaseEntity records.
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	b.CreatedAt = &now
	if b.CreatedBy == "" {
		b.CreatedBy = "SYSTEM"
	}
	b.IsActive = true
	return nil
}

// BeforeUpdate sets audit fields on modified BaseEntity records.
func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now().UTC()
	b.UpdatedAt = &now
	return nil
}
