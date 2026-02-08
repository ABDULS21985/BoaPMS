package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/jmoiron/sqlx"
)

// StaffRepository provides data access for the Staff ID Mask SQL Server database.
type StaffRepository struct {
	db *sqlx.DB // StaffIDSQL connection
}

// NewStaffRepository creates a new Staff ID Mask repository.
func NewStaffRepository(db *sqlx.DB) *StaffRepository {
	if db == nil {
		return nil
	}
	return &StaffRepository{db: db}
}

// GetStaffIDMask retrieves staff ID mask details by employee number.
func (r *StaffRepository) GetStaffIDMask(ctx context.Context, employeeID string) (*erp.StaffIDMaskDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("staffRepo.GetStaffIDMask: StaffIDMask database not configured")
	}
	var result erp.StaffIDMaskDetails
	err := r.db.GetContext(ctx, &result,
		`SELECT * FROM dbo.StaffIDMaskDetails WHERE StaffId = @p1`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("staffRepo.GetStaffIDMask: %w", err)
	}
	return &result, nil
}

// GetAllStaffIDMasks retrieves all staff ID mask records.
func (r *StaffRepository) GetAllStaffIDMasks(ctx context.Context) ([]erp.StaffIDMaskDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("staffRepo.GetAllStaffIDMasks: StaffIDMask database not configured")
	}
	var results []erp.StaffIDMaskDetails
	err := r.db.SelectContext(ctx, &results, `SELECT * FROM dbo.StaffIDMaskDetails`)
	if err != nil {
		return nil, fmt.Errorf("staffRepo.GetAllStaffIDMasks: %w", err)
	}
	return results, nil
}
