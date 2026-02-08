package audit

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// AuditLog records property-level change history for audited entities.
type AuditLog struct {
	domain.BaseEntity
	UserName          string         `json:"user_name"            gorm:"column:user_name"`
	AuditEventDateUTC time.Time      `json:"audit_event_date_utc" gorm:"column:audit_event_date_utc;index"`
	AuditEventType    enums.AuditEventType `json:"audit_event_type" gorm:"column:audit_event_type"`
	AuditTableName    string         `json:"table_name"           gorm:"column:table_name;index"`
	RecordID          string         `json:"record_id"            gorm:"column:record_id"`
	FieldName         string         `json:"field_name"           gorm:"column:field_name"`
	OriginalValue     string         `json:"original_value"       gorm:"column:original_value;type:text"`
	NewValue          string         `json:"new_value"            gorm:"column:new_value;type:text"`
}

func (AuditLog) TableName() string { return "pmsaudit.audit_logs" }

// AuditableEntity defines which entities are audit-tracked.
type AuditableEntity struct {
	domain.BaseEntity
	EntityName         string             `json:"entity_name"  gorm:"column:entity_name;primaryKey"`
	EnableAudit        bool               `json:"enable_audit" gorm:"column:enable_audit;default:true"`
	AuditableAttributes []AuditableAttribute `json:"auditable_attributes" gorm:"foreignKey:AuditableEntityID"`
}

func (AuditableEntity) TableName() string { return "pmsaudit.auditable_entities" }

// AuditableAttribute defines which attributes of an entity are audit-tracked.
type AuditableAttribute struct {
	domain.BaseEntity
	AuditableEntityID int    `json:"auditable_entity_id" gorm:"column:auditable_entity_id"`
	AttributeName     string `json:"attribute_name"      gorm:"column:attribute_name"`
	EnableAudit       bool   `json:"enable_audit"        gorm:"column:enable_audit;default:true"`
	AuditableEntity   *AuditableEntity `json:"auditable_entity" gorm:"foreignKey:AuditableEntityID"`
}

func (AuditableAttribute) TableName() string { return "pmsaudit.auditable_attributes" }

// SequenceNumber generates unique sequential IDs with optional prefixes.
type SequenceNumber struct {
	domain.BaseEntity
	SequenceNumberType int    `json:"sequence_number_type" gorm:"column:sequence_number_type"`
	Description        string `json:"description"          gorm:"column:description;primaryKey"`
	Prefix             string `json:"prefix"               gorm:"column:prefix"`
	NextNumber         int64  `json:"next_number"          gorm:"column:next_number"`
	UsePrefix          bool   `json:"use_prefix"           gorm:"column:use_prefix;default:false"`
}

func (SequenceNumber) TableName() string { return "pms.sequence_numbers" }
