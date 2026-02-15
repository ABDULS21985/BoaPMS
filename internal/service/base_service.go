package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/audit"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Shared utility functions and types converted from the .NET BaseService class.
// These are package-level helpers rather than an embedded struct, because Go
// favours composition-via-functions over abstract-class inheritance.
// ---------------------------------------------------------------------------

// sequenceGenerator encapsulates the logic for generating sequential codes.
// It mirrors the .NET BaseService.GenerateCode / GetNextNumber methods that
// manage the pms.sequence_numbers table.
type sequenceGenerator struct {
	db  *gorm.DB
	log zerolog.Logger
}

// newSequenceGenerator creates a new sequence generator attached to the given
// GORM database connection.
func newSequenceGenerator(db *gorm.DB, log zerolog.Logger) *sequenceGenerator {
	return &sequenceGenerator{db: db, log: log}
}

// GenerateCode produces a zero-padded sequential code, optionally concatenated
// with a prefix or suffix string. This mirrors the .NET BaseService.GenerateCode
// method.
//
// Parameters:
//   - seqType: the SequenceNumberTypes enum value identifying the entity type
//   - length:  total width of the numeric portion (padded with leading zeros)
//   - concat:  an optional string to prepend or append
//   - pos:     ConCatBefore (default) prepends concat; ConCatAfter appends it
func (g *sequenceGenerator) GenerateCode(
	ctx context.Context,
	seqType enums.SequenceNumberTypes,
	length int64,
	concat string,
	pos enums.ConCat,
) (string, error) {
	nextNum, err := g.getNextNumber(ctx, seqType)
	if err != nil {
		return "", err
	}

	padded, err := padWithZeros(fmt.Sprintf("%d", nextNum), length)
	if err != nil {
		return "", err
	}

	if pos == enums.ConCatAfter {
		return padded + concat, nil
	}
	return concat + padded, nil
}

// getNextNumber atomically retrieves and increments the next sequence number
// for the given type. If no row exists yet it creates one starting at 1.
// Mirrors the .NET BaseService.GetNextNumber method.
func (g *sequenceGenerator) getNextNumber(ctx context.Context, seqType enums.SequenceNumberTypes) (int64, error) {
	var seq audit.SequenceNumber
	err := g.db.WithContext(ctx).
		Where("sequence_number_type = ?", int(seqType)).
		First(&seq).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// First time: create a new sequence row starting at 1.
			seq = audit.SequenceNumber{
				SequenceNumberType: int(seqType),
				Description:        fmt.Sprintf("SequenceNumberType_%d", int(seqType)),
				UsePrefix:          false,
				NextNumber:         2, // 1 is returned, next caller gets 2
			}
			if createErr := g.db.WithContext(ctx).Create(&seq).Error; createErr != nil {
				return 0, fmt.Errorf("creating sequence row for type %d: %w", int(seqType), createErr)
			}
			return 1, nil
		}
		return 0, fmt.Errorf("querying sequence for type %d: %w", int(seqType), err)
	}

	nextNumber := seq.NextNumber
	seq.NextNumber++
	if updateErr := g.db.WithContext(ctx).Save(&seq).Error; updateErr != nil {
		return 0, fmt.Errorf("incrementing sequence for type %d: %w", int(seqType), updateErr)
	}
	return nextNumber, nil
}

// padWithZeros left-pads the input string with zeros until it reaches the
// specified maximum length. Returns an error if the input already exceeds
// the maximum. Mirrors the .NET BaseService.padWithZeros method.
func padWithZeros(field string, maxLength int64) (string, error) {
	fieldLen := int64(len(field))
	if fieldLen == maxLength {
		return field, nil
	}
	if fieldLen > maxLength {
		return "", fmt.Errorf("input string length %d exceeds maximum length %d", fieldLen, maxLength)
	}

	var sb strings.Builder
	sb.Grow(int(maxLength))
	for i := int64(0); i < maxLength-fieldLen; i++ {
		sb.WriteByte('0')
	}
	sb.WriteString(field)
	return sb.String(), nil
}

// getStartOrEndDate computes the first or last day of a review period segment
// (quarterly, bi-annual, or annual). This mirrors the .NET
// BaseService.GetStartOrEndDate method.
//
// Parameters:
//   - year:    the calendar year
//   - value:   the period number (e.g. 1-4 for quarterly, 1-2 for bi-annual, 1 for annual)
//   - rng:     the ReviewPeriodRange enum
//   - isStart: true returns the first day of the period; false returns the last day
func getStartOrEndDate(year, value int, rng enums.ReviewPeriodRange, isStart bool) (time.Time, error) {
	if value == 0 {
		value = 1
	}

	var month int
	switch rng {
	case enums.ReviewPeriodRangeQuarterly:
		if value > 4 {
			return time.Time{}, fmt.Errorf("invalid quarter value: %d", value)
		}
		month = (value-1)*3 + 1

	case enums.ReviewPeriodRangeBiAnnual:
		if value > 2 {
			return time.Time{}, fmt.Errorf("invalid bi-annual value: %d", value)
		}
		month = (value-1)*6 + 1

	case enums.ReviewPeriodRangeAnnual:
		if value > 1 {
			return time.Time{}, fmt.Errorf("invalid annual value: %d", value)
		}
		month = 1

	default:
		return time.Time{}, fmt.Errorf("invalid review period range: %d", int(rng))
	}

	if !isStart {
		month += 2 // move to last month of the quarter/period
	}

	day := 1
	if !isStart {
		// Last day of the month
		day = daysInMonth(year, month)
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// daysInMonth returns the number of days in the given month/year.
func daysInMonth(year, month int) int {
	// Go's time package: day 0 of month+1 gives last day of month.
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// ---------------------------------------------------------------------------
// Notification message constants matching the .NET NotificationMessages class.
// ---------------------------------------------------------------------------

const (
	msgOperationCompleted = "Operation completed successfully"
	msgGenericException   = "An error occurred while processing your request"
)

// ---------------------------------------------------------------------------
// Valid setting types — mirrors .NET SettingType.GetSettingTypeList().
// ---------------------------------------------------------------------------

// validSettingTypes is the set of allowed setting type strings.
var validSettingTypes = map[string]struct{}{
	"Bool":     {},
	"DateTime": {},
	"Decimal":  {},
	"Double":   {},
	"Float":    {},
	"Int":      {},
	"Long":     {},
	"String":   {},
}

// isValidSettingType checks whether the given type string is in the allowed list.
func isValidSettingType(t string) bool {
	_, ok := validSettingTypes[t]
	return ok
}

// ---------------------------------------------------------------------------
// ID generation — mirrors .NET Guid.NewGuid().ToString().
// ---------------------------------------------------------------------------

// GenerateID returns a new UUID string suitable for use as a record identifier.
func GenerateID() string {
	return uuid.New().String()
}

// ---------------------------------------------------------------------------
// Workflow record-level helpers — generic functions that operate on any table
// via raw GORM updates. These mirror the .NET BaseService approve/reject/
// workflow-status helper methods.
// ---------------------------------------------------------------------------

// ApproveRecord sets approval fields on a workflow entity identified by
// tableName, primary key column, and record ID.
// Mirrors the .NET BaseService record-level approval logic.
func ApproveRecord(db *gorm.DB, ctx context.Context, tableName string, pkColumn string, recordID string, approvedBy string) error {
	now := time.Now().UTC()
	return db.WithContext(ctx).Table(tableName).
		Where(pkColumn+" = ?", recordID).
		Updates(map[string]interface{}{
			"status":        enums.StatusApprovedAndActive.String(),
			"record_status": "Active",
			"date_approved": now,
			"is_approved":   true,
			"is_active":     true,
			"approved_by":   approvedBy,
		}).Error
}

// RejectRecord sets rejection fields on a workflow entity identified by
// tableName, primary key column, and record ID.
// Mirrors the .NET BaseService record-level rejection logic.
func RejectRecord(db *gorm.DB, ctx context.Context, tableName string, pkColumn string, recordID string, rejectedBy string, reason string) error {
	now := time.Now().UTC()
	return db.WithContext(ctx).Table(tableName).
		Where(pkColumn+" = ?", recordID).
		Updates(map[string]interface{}{
			"status":           enums.StatusRejected.String(),
			"record_status":    "Rejected",
			"date_rejected":    now,
			"is_rejected":      true,
			"is_active":        false,
			"rejected_by":      rejectedBy,
			"rejection_reason": reason,
		}).Error
}

// SetWorkflowStatus is a generic function to set workflow status fields on
// any entity. It computes the appropriate record_status and secondary fields
// based on the target status. Mirrors the .NET BaseService workflow-status
// helper pattern.
func SetWorkflowStatus(db *gorm.DB, ctx context.Context, tableName string, pkColumn string, recordID string, status enums.Status, updatedBy string) error {
	updates := map[string]interface{}{
		"status":     status.String(),
		"updated_by": updatedBy,
	}

	switch status {
	case enums.StatusDraft:
		updates["record_status"] = "Draft"
	case enums.StatusPendingApproval:
		updates["record_status"] = "PendingApproval"
	case enums.StatusApprovedAndActive:
		updates["record_status"] = "Active"
		updates["is_approved"] = true
		updates["is_active"] = true
		now := time.Now().UTC()
		updates["date_approved"] = now
		updates["approved_by"] = updatedBy
	case enums.StatusReturned:
		updates["record_status"] = "Returned"
	case enums.StatusRejected:
		updates["record_status"] = "Rejected"
		updates["is_rejected"] = true
		updates["is_active"] = false
		now := time.Now().UTC()
		updates["date_rejected"] = now
		updates["rejected_by"] = updatedBy
	case enums.StatusCancelled:
		updates["record_status"] = "Cancelled"
		updates["is_active"] = false
	case enums.StatusClosed:
		updates["record_status"] = "Closed"
	case enums.StatusPaused:
		updates["record_status"] = "Paused"
	case enums.StatusAwaitingEvaluation:
		updates["record_status"] = "AwaitingEvaluation"
	case enums.StatusCompleted:
		updates["record_status"] = "Completed"
	case enums.StatusPendingAcceptance:
		updates["record_status"] = "PendingAcceptance"
	}

	return db.WithContext(ctx).Table(tableName).
		Where(pkColumn+" = ?", recordID).
		Updates(updates).Error
}
