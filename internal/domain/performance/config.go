package performance

import (
	"github.com/enterprise-pms/pms-api/internal/domain"
)

// PmsConfiguration stores system-level PMS configuration key-value pairs.
type PmsConfiguration struct {
	PmsConfigurationID string `json:"pms_configuration_id" gorm:"column:pms_configuration_id;primaryKey"`
	Name               string `json:"name"                 gorm:"column:name;not null"`
	Value              string `json:"value"                gorm:"column:value"`
	Type               string `json:"type"                 gorm:"column:type"`
	IsEncrypted        bool   `json:"is_encrypted"         gorm:"column:is_encrypted;default:false"`
	domain.BaseEntity
}

func (PmsConfiguration) TableName() string { return "pms.pms_configurations" }

// Setting stores application settings as typed key-value pairs.
type Setting struct {
	SettingID   string `json:"setting_id"   gorm:"column:setting_id;primaryKey"`
	Name        string `json:"name"         gorm:"column:name;not null"`
	Value       string `json:"value"        gorm:"column:value"`
	Type        string `json:"type"         gorm:"column:type"`
	IsEncrypted bool   `json:"is_encrypted" gorm:"column:is_encrypted;default:false"`
	domain.BaseEntity
}

func (Setting) TableName() string { return "pms.settings" }

// SettingType constants for typed setting values.
const (
	SettingTypeBool     = "Bool"
	SettingTypeDateTime = "DateTime"
	SettingTypeDecimal  = "Decimal"
	SettingTypeDouble   = "Double"
	SettingTypeFloat    = "Float"
	SettingTypeInt      = "Int"
	SettingTypeLong     = "Long"
	SettingTypeString   = "String"
)
