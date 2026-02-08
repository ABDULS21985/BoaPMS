package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/sas"
	"github.com/jmoiron/sqlx"
)

// SasRepository provides data access for the SAS (Staff Attendance System) SQL Server database.
type SasRepository struct {
	db *sqlx.DB // SasSQL connection
}

// NewSasRepository creates a new SAS repository.
func NewSasRepository(db *sqlx.DB) *SasRepository {
	if db == nil {
		return nil
	}
	return &SasRepository{db: db}
}

// GetAbsenceModes retrieves all absence mode records.
func (r *SasRepository) GetAbsenceModes(ctx context.Context) ([]sas.AbsenceMode, error) {
	if r.db == nil {
		return nil, fmt.Errorf("sasRepo.GetAbsenceModes: SAS database not configured")
	}
	var results []sas.AbsenceMode
	err := r.db.SelectContext(ctx, &results, `SELECT * FROM dbo.XXCBN_SAS_AbsenceMode`)
	if err != nil {
		return nil, fmt.Errorf("sasRepo.GetAbsenceModes: %w", err)
	}
	return results, nil
}

// GetStaffAttendanceByEmployee retrieves lunch attendance records for an employee.
func (r *SasRepository) GetStaffAttendanceByEmployee(ctx context.Context, employeeNumber string) ([]sas.StaffLunchAttendance, error) {
	if r.db == nil {
		return nil, fmt.Errorf("sasRepo.GetStaffAttendanceByEmployee: SAS database not configured")
	}
	var results []sas.StaffLunchAttendance
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.XXCBN_SAS_StaffLunchAttendance WHERE EmployeeNumber = @p1 ORDER BY CreateDate DESC`,
		employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("sasRepo.GetStaffAttendanceByEmployee: %w", err)
	}
	return results, nil
}

// GetStaffAttendanceByDepartment retrieves lunch attendance for an entire department.
func (r *SasRepository) GetStaffAttendanceByDepartment(ctx context.Context, deptID int) ([]sas.StaffLunchAttendance, error) {
	if r.db == nil {
		return nil, fmt.Errorf("sasRepo.GetStaffAttendanceByDepartment: SAS database not configured")
	}
	var results []sas.StaffLunchAttendance
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.XXCBN_SAS_StaffLunchAttendance WHERE DepartmentId = @p1 ORDER BY CreateDate DESC`,
		deptID)
	if err != nil {
		return nil, fmt.Errorf("sasRepo.GetStaffAttendanceByDepartment: %w", err)
	}
	return results, nil
}
