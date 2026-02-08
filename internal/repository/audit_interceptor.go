package repository

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/audit"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AuditInterceptor provides GORM callbacks that automatically log entity changes.
// This mirrors the .NET SaveAuditLog() / UpdateAuditLogRecordId() pattern in SaveChangesAsync.
type AuditInterceptor struct {
	db  *gorm.DB
	log zerolog.Logger
}

// NewAuditInterceptor creates an interceptor and registers GORM callbacks.
func NewAuditInterceptor(db *gorm.DB, log zerolog.Logger) *AuditInterceptor {
	ai := &AuditInterceptor{db: db, log: log}
	ai.register()
	return ai
}

func (ai *AuditInterceptor) register() {
	ai.db.Callback().Create().After("gorm:create").Register("audit:after_create", ai.afterCreate)
	ai.db.Callback().Update().After("gorm:update").Register("audit:after_update", ai.afterUpdate)
	ai.db.Callback().Delete().After("gorm:delete").Register("audit:after_delete", ai.afterDelete)
}

func (ai *AuditInterceptor) afterCreate(db *gorm.DB) {
	if db.Error != nil || db.Statement == nil || db.Statement.Model == nil {
		return
	}
	ai.logAudit(db, enums.AuditEventAdded)
}

func (ai *AuditInterceptor) afterUpdate(db *gorm.DB) {
	if db.Error != nil || db.Statement == nil || db.Statement.Model == nil {
		return
	}
	ai.logAudit(db, enums.AuditEventModified)
}

func (ai *AuditInterceptor) afterDelete(db *gorm.DB) {
	if db.Error != nil || db.Statement == nil || db.Statement.Model == nil {
		return
	}
	ai.logAudit(db, enums.AuditEventDeleted)
}

func (ai *AuditInterceptor) logAudit(db *gorm.DB, eventType enums.AuditEventType) {
	stmt := db.Statement
	if stmt.Schema == nil {
		return
	}

	tableName := stmt.Schema.Table
	model := stmt.Model
	if model == nil {
		return
	}

	// Extract the user from context if available
	username := "SYSTEM"
	if ctx := db.Statement.Context; ctx != nil {
		if u, ok := ctx.Value(auditUserKey{}).(string); ok && u != "" {
			username = u
		}
	}

	now := time.Now().UTC()

	switch eventType {
	case enums.AuditEventAdded:
		// Log all non-zero fields for new records
		ai.logNewRecord(tableName, model, username, now)

	case enums.AuditEventModified:
		// Log changed fields from GORM's changed map
		ai.logModifiedRecord(db, tableName, model, username, now)

	case enums.AuditEventDeleted:
		ai.saveAuditEntry(audit.AuditLog{
			UserName:          username,
			AuditEventDateUTC: now,
			AuditEventType:    eventType,
			AuditTableName:    tableName,
			RecordID:          extractRecordID(model),
			FieldName:         "soft_deleted",
			OriginalValue:     "false",
			NewValue:          "true",
		})
	}
}

func (ai *AuditInterceptor) logNewRecord(tableName string, model interface{}, username string, now time.Time) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	recordID := extractRecordID(model)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if shouldSkipAuditField(field.Name) {
			continue
		}
		fieldVal := v.Field(i)
		if fieldVal.IsZero() {
			continue
		}
		newVal := formatValue(fieldVal)
		ai.saveAuditEntry(audit.AuditLog{
			UserName:          username,
			AuditEventDateUTC: now,
			AuditEventType:    enums.AuditEventAdded,
			AuditTableName:    tableName,
			RecordID:          recordID,
			FieldName:         field.Name,
			OriginalValue:     "",
			NewValue:          newVal,
		})
	}
}

func (ai *AuditInterceptor) logModifiedRecord(db *gorm.DB, tableName string, model interface{}, username string, now time.Time) {
	recordID := extractRecordID(model)

	// GORM stores changed columns in Statement.Changed
	if db.Statement.Changed() {
		for col, change := range db.Statement.Clauses {
			_ = change // iterate changed values
			ai.saveAuditEntry(audit.AuditLog{
				UserName:          username,
				AuditEventDateUTC: now,
				AuditEventType:    enums.AuditEventModified,
				AuditTableName:    tableName,
				RecordID:          recordID,
				FieldName:         col,
			})
		}
	}
}

func (ai *AuditInterceptor) saveAuditEntry(entry audit.AuditLog) {
	if err := ai.db.Create(&entry).Error; err != nil {
		ai.log.Error().Err(err).
			Str("table", entry.AuditTableName).
			Str("record_id", entry.RecordID).
			Msg("Failed to save audit log entry")
	}
}

// extractRecordID gets the primary key value from a model.
func extractRecordID(model interface{}) string {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Try common PK field names
	for _, pkName := range []string{"ID", "Id"} {
		f := v.FieldByName(pkName)
		if f.IsValid() && !f.IsZero() {
			return fmt.Sprintf("%v", f.Interface())
		}
	}

	// Scan for any field ending in "ID" that is tagged as primaryKey
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("gorm")
		if contains(tag, "primaryKey") {
			return fmt.Sprintf("%v", v.Field(i).Interface())
		}
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func shouldSkipAuditField(name string) bool {
	skip := map[string]bool{
		"ID": true, "Id": true, "CreatedAt": true, "UpdatedAt": true,
		"CreatedBy": true, "UpdatedBy": true, "SoftDeleted": true,
		"DateCreated": true, "DateUpdated": true, "RecordStatus": true,
	}
	return skip[name]
}

func formatValue(v reflect.Value) string {
	if !v.IsValid() || v.IsZero() {
		return ""
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		b, err := json.Marshal(v.Interface())
		if err != nil {
			return fmt.Sprintf("%v", v.Interface())
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// auditUserKey is a context key for passing the current user to audit callbacks.
type auditUserKey struct{}
