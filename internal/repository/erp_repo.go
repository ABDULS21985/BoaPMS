package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/jmoiron/sqlx"
)

// ActiveStaffPersonType is the ERP person type ID for active staff.
const ActiveStaffPersonType = 1120

// ErpRepository provides data access for the external ERP SQL Server database.
// All queries use sqlx (Dapper-style) since the ERP database is read-only.
type ErpRepository struct {
	db *sqlx.DB // ErpSQL connection
}

// NewErpRepository creates a new ERP repository.
func NewErpRepository(db *sqlx.DB) *ErpRepository {
	if db == nil {
		return nil
	}
	return &ErpRepository{db: db}
}

// ─── Employee Queries ────────────────────────────────────────────────────────

// GetEmployeeByID retrieves a single employee by employee number.
func (r *ErpRepository) GetEmployeeByID(ctx context.Context, employeeID string) (*erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeeByID: ERP database not configured")
	}
	var emp erp.EmployeeDetails
	err := r.db.GetContext(ctx, &emp,
		`SELECT * FROM dbo.EmployeeDetails WHERE EmployeeNumber = @p1`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeeByID: %w", err)
	}
	return &emp, nil
}

// GetEmployeeByUserName retrieves a single employee by username.
func (r *ErpRepository) GetEmployeeByUserName(ctx context.Context, userName string) (*erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeeByUserName: ERP database not configured")
	}
	var emp erp.EmployeeDetails
	err := r.db.GetContext(ctx, &emp,
		`SELECT * FROM dbo.EmployeeDetails WHERE UserName = @p1`, userName)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeeByUserName: %w", err)
	}
	return &emp, nil
}

// GetAllActiveEmployees retrieves all active employees.
func (r *ErpRepository) GetAllActiveEmployees(ctx context.Context) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetAllActiveEmployees: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails WHERE EmployeeNumber IS NOT NULL AND PersonTypeId = @p1`,
		ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetAllActiveEmployees: %w", err)
	}
	return results, nil
}

// GetSubordinates retrieves direct subordinates of an employee.
func (r *ErpRepository) GetSubordinates(ctx context.Context, supervisorID string) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinates: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails WHERE SupervisorId = @p1 AND PersonTypeId = @p2`,
		supervisorID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinates: %w", err)
	}
	return results, nil
}

// GetEmployeesByOfficeAndGrade retrieves peers in the same office and grade.
func (r *ErpRepository) GetEmployeesByOfficeAndGrade(ctx context.Context, excludeID string, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeesByOfficeAndGrade: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE OfficeId = @p1 AND Grade = @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		officeID, grade, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetEmployeesByOfficeAndGrade: %w", err)
	}
	return results, nil
}

// GetSubordinatesByOfficeAndGrades retrieves employees in the same office with lower grade (higher number).
func (r *ErpRepository) GetSubordinatesByOfficeAndGrades(ctx context.Context, excludeID string, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByOfficeAndGrades: ERP database not configured")
	}
	gradeNum, err := strconv.Atoi(grade)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByOfficeAndGrades: invalid grade %q: %w", grade, err)
	}
	var results []erp.EmployeeDetails
	// In ERP, higher grade number = lower rank (subordinate)
	err = r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE OfficeId = @p1 AND CAST(Grade AS INT) > @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		officeID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByOfficeAndGrades: %w", err)
	}
	return results, nil
}

// GetSuperiorsByOfficeAndGrades retrieves employees in the same office with higher grade (lower number).
func (r *ErpRepository) GetSuperiorsByOfficeAndGrades(ctx context.Context, excludeID string, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByOfficeAndGrades: ERP database not configured")
	}
	gradeNum, err := strconv.Atoi(grade)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByOfficeAndGrades: invalid grade: %w", err)
	}
	var results []erp.EmployeeDetails
	err = r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE OfficeId = @p1 AND CAST(Grade AS INT) < @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		officeID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByOfficeAndGrades: %w", err)
	}
	return results, nil
}

// GetSubordinatesByDivisionAndGrades retrieves subordinates in a division.
func (r *ErpRepository) GetSubordinatesByDivisionAndGrades(ctx context.Context, excludeID string, grade string, divisionID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByDivisionAndGrades: ERP database not configured")
	}
	gradeNum, _ := strconv.Atoi(grade)
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DivisionId = @p1 AND CAST(Grade AS INT) > @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		divisionID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByDivisionAndGrades: %w", err)
	}
	return results, nil
}

// GetSuperiorsByDivisionAndGrades retrieves superiors in a division.
func (r *ErpRepository) GetSuperiorsByDivisionAndGrades(ctx context.Context, excludeID string, grade string, divisionID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByDivisionAndGrades: ERP database not configured")
	}
	gradeNum, _ := strconv.Atoi(grade)
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DivisionId = @p1 AND CAST(Grade AS INT) < @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		divisionID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByDivisionAndGrades: %w", err)
	}
	return results, nil
}

// GetPeersByDivisionAndGrade retrieves peers in a division with the same grade.
func (r *ErpRepository) GetPeersByDivisionAndGrade(ctx context.Context, excludeID string, grade string, divisionID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetPeersByDivisionAndGrade: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DivisionId = @p1 AND Grade = @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		divisionID, grade, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetPeersByDivisionAndGrade: %w", err)
	}
	return results, nil
}

// GetPeersByDepartmentAndGrade retrieves peers in a department with the same grade.
func (r *ErpRepository) GetPeersByDepartmentAndGrade(ctx context.Context, excludeID string, grade string, deptID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetPeersByDepartmentAndGrade: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DepartmentId = @p1 AND Grade = @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		deptID, grade, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetPeersByDepartmentAndGrade: %w", err)
	}
	return results, nil
}

// GetSubordinatesByDepartmentAndGrades retrieves subordinates in a department.
func (r *ErpRepository) GetSubordinatesByDepartmentAndGrades(ctx context.Context, excludeID string, grade string, deptID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByDepartmentAndGrades: ERP database not configured")
	}
	gradeNum, _ := strconv.Atoi(grade)
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DepartmentId = @p1 AND CAST(Grade AS INT) > @p2 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		deptID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSubordinatesByDepartmentAndGrades: %w", err)
	}
	return results, nil
}

// GetSuperiorsByDepartmentAndGrades retrieves superiors in a department (including grade 41).
func (r *ErpRepository) GetSuperiorsByDepartmentAndGrades(ctx context.Context, excludeID string, grade string, deptID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByDepartmentAndGrades: ERP database not configured")
	}
	gradeNum, _ := strconv.Atoi(grade)
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DepartmentId = @p1 AND (CAST(Grade AS INT) < @p2 OR Grade = '41')
		 AND EmployeeNumber != @p3 AND PersonTypeId = @p4`,
		deptID, gradeNum, excludeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetSuperiorsByDepartmentAndGrades: %w", err)
	}
	return results, nil
}

// ─── Filtered Employee Lists ─────────────────────────────────────────────────

// GetByDepartmentID retrieves all active employees in a department.
func (r *ErpRepository) GetByDepartmentID(ctx context.Context, deptID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetByDepartmentID: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DepartmentId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = @p2`,
		deptID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetByDepartmentID: %w", err)
	}
	return results, nil
}

// GetByDivisionID retrieves all active employees in a division.
func (r *ErpRepository) GetByDivisionID(ctx context.Context, divisionID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetByDivisionID: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE DivisionId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = @p2`,
		divisionID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetByDivisionID: %w", err)
	}
	return results, nil
}

// GetByOfficeID retrieves all active employees in an office.
func (r *ErpRepository) GetByOfficeID(ctx context.Context, officeID int) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetByOfficeID: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE OfficeId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = @p2`,
		officeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetByOfficeID: %w", err)
	}
	return results, nil
}

// ─── Organization Queries ────────────────────────────────────────────────────

// AllDepartments returns distinct departments from ERP.
func (r *ErpRepository) AllDepartments(ctx context.Context) ([]erp.ErpOrganizationVm, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.AllDepartments: ERP database not configured")
	}
	var results []erp.ErpOrganizationVm
	err := r.db.SelectContext(ctx, &results,
		`SELECT DISTINCT DepartmentId, DepartmentName FROM dbo.EmployeeDetails WHERE DepartmentName IS NOT NULL`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.AllDepartments: %w", err)
	}
	return results, nil
}

// AllDivisions returns distinct divisions from ERP.
func (r *ErpRepository) AllDivisions(ctx context.Context) ([]erp.ErpOrganizationVm, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.AllDivisions: ERP database not configured")
	}
	var results []erp.ErpOrganizationVm
	err := r.db.SelectContext(ctx, &results,
		`SELECT DISTINCT DepartmentId, DepartmentName, DivisionId, DivisionName
		 FROM dbo.EmployeeDetails WHERE DivisionName IS NOT NULL`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.AllDivisions: %w", err)
	}
	return results, nil
}

// AllOffices returns distinct offices from ERP.
func (r *ErpRepository) AllOffices(ctx context.Context) ([]erp.ErpOrganizationVm, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.AllOffices: ERP database not configured")
	}
	var results []erp.ErpOrganizationVm
	err := r.db.SelectContext(ctx, &results,
		`SELECT DISTINCT DepartmentId, DepartmentName, DivisionId, DivisionName, OfficeId, OfficeName
		 FROM dbo.EmployeeDetails WHERE OfficeName IS NOT NULL`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.AllOffices: %w", err)
	}
	return results, nil
}

// AllJobGrades returns distinct job grades from ERP.
func (r *ErpRepository) AllJobGrades(ctx context.Context) ([]erp.EROJobGradeVm, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.AllJobGrades: ERP database not configured")
	}
	var results []erp.EROJobGradeVm
	err := r.db.SelectContext(ctx, &results,
		`SELECT DISTINCT GradeId AS grade_id, Grade AS grade_name
		 FROM dbo.EmployeeDetails WHERE Grade IS NOT NULL AND Grade != ''`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.AllJobGrades: %w", err)
	}
	return results, nil
}

// AllOfficeJobRoles returns distinct job roles per office.
func (r *ErpRepository) AllOfficeJobRoles(ctx context.Context) ([]erp.ERPOfficeJobRoleVm, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.AllOfficeJobRoles: ERP database not configured")
	}
	var raw []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &raw,
		`SELECT * FROM dbo.EmployeeDetails`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.AllOfficeJobRoles: %w", err)
	}

	seen := make(map[string]bool)
	var results []erp.ERPOfficeJobRoleVm
	for _, e := range raw {
		pos := e.JobTitle
		if pos == "" {
			continue
		}
		roleName := pos
		if idx := strings.Index(pos, "."); idx > 0 {
			roleName = pos[:idx]
		}
		if seen[roleName] {
			continue
		}
		seen[roleName] = true
		oid := 0
		if e.OfficeID != nil {
			oid = *e.OfficeID
		}
		results = append(results, erp.ERPOfficeJobRoleVm{
			OfficeFullName: e.Office,
			OfficeID:       oid,
			OfficeName:     pos,
			JobRoleName:    roleName,
		})
	}
	return results, nil
}

// ─── Head Queries ────────────────────────────────────────────────────────────

// GetHeadOfOfficeIDs returns head-of-office IDs for a given office.
func (r *ErpRepository) GetHeadOfOfficeIDs(ctx context.Context, officeID int) ([]string, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfOfficeIDs: ERP database not configured")
	}
	var ids []string
	err := r.db.SelectContext(ctx, &ids,
		`SELECT DISTINCT HeadOfOfficeId FROM dbo.EmployeeDetails
		 WHERE OfficeId = @p1 AND PersonTypeId = @p2`,
		officeID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfOfficeIDs: %w", err)
	}
	return ids, nil
}

// GetHeadOfDivisionIDs returns head-of-division IDs for a given division.
func (r *ErpRepository) GetHeadOfDivisionIDs(ctx context.Context, divisionID int) ([]string, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfDivisionIDs: ERP database not configured")
	}
	var ids []string
	err := r.db.SelectContext(ctx, &ids,
		`SELECT DISTINCT HeadOfDivId FROM dbo.EmployeeDetails
		 WHERE DivisionId = @p1 AND PersonTypeId = @p2`,
		divisionID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfDivisionIDs: %w", err)
	}
	return ids, nil
}

// GetHeadOfDepartmentIDs returns head-of-department IDs for a given department.
func (r *ErpRepository) GetHeadOfDepartmentIDs(ctx context.Context, deptID int) ([]string, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfDepartmentIDs: ERP database not configured")
	}
	var ids []string
	err := r.db.SelectContext(ctx, &ids,
		`SELECT DISTINCT HeadOfDeptId FROM dbo.EmployeeDetails
		 WHERE DepartmentId = @p1 AND PersonTypeId = @p2`,
		deptID, ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetHeadOfDepartmentIDs: %w", err)
	}
	return ids, nil
}

// GetGovernorDGEmails returns email addresses of Governor/Deputy Governors.
func (r *ErpRepository) GetGovernorDGEmails(ctx context.Context) ([]string, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetGovernorDGEmails: ERP database not configured")
	}
	var emails []string
	err := r.db.SelectContext(ctx, &emails,
		`SELECT EmailAddress FROM dbo.EmployeeDetails
		 WHERE JobName LIKE '%GOVERNOR%' AND PersonTypeId = @p1`,
		ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetGovernorDGEmails: %w", err)
	}
	return emails, nil
}

// GetGovernorDGs returns Governor/Deputy Governor employee records.
func (r *ErpRepository) GetGovernorDGs(ctx context.Context) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetGovernorDGs: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE JobName LIKE '%GOVERNOR%' AND PersonTypeId = @p1`,
		ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetGovernorDGs: %w", err)
	}
	return results, nil
}

// GetAllHeadDepartments returns employees who are heads of their departments (excluding governors).
func (r *ErpRepository) GetAllHeadDepartments(ctx context.Context) ([]erp.EmployeeDetails, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetAllHeadDepartments: ERP database not configured")
	}
	var results []erp.EmployeeDetails
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmployeeDetails
		 WHERE PersonTypeId = @p1 AND HeadOfDeptId = EmployeeNumber AND JobName NOT LIKE '%GOVERNOR%'`,
		ActiveStaffPersonType)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetAllHeadDepartments: %w", err)
	}
	return results, nil
}

// ─── Location Data ───────────────────────────────────────────────────────────

// GetAllLocations retrieves all ERP location records.
func (r *ErpRepository) GetAllLocations(ctx context.Context) ([]erp.ErpLocationDetail, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetAllLocations: ERP database not configured")
	}
	var results []erp.ErpLocationDetail
	err := r.db.SelectContext(ctx, &results, `SELECT * FROM dbo.ErpLocationDetails`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetAllLocations: %w", err)
	}
	return results, nil
}

// ─── Holiday / Vacation Data ─────────────────────────────────────────────────

// GetPublicHolidays retrieves public holiday records.
func (r *ErpRepository) GetPublicHolidays(ctx context.Context) ([]erp.PublicHolidayData, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetPublicHolidays: ERP database not configured")
	}
	var results []erp.PublicHolidayData
	err := r.db.SelectContext(ctx, &results, `SELECT * FROM dbo.HOLIDAYS_T24`)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetPublicHolidays: %w", err)
	}
	return results, nil
}

// GetVacationRules retrieves vacation rules for an employee.
func (r *ErpRepository) GetVacationRules(ctx context.Context, employeeID string) ([]erp.VacationRuleData, error) {
	if r.db == nil {
		return nil, fmt.Errorf("erpRepo.GetVacationRules: ERP database not configured")
	}
	var results []erp.VacationRuleData
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.VACATIONSRULE_DATA WHERE EMPLOYEE_ID = @p1`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("erpRepo.GetVacationRules: %w", err)
	}
	return results, nil
}
