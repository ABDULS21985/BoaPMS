package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// erpEmployeeService implements ErpEmployeeService.
// It retrieves employee and organisational data from the ERP SQL Server database
// and manages staff job roles via the primary PostgreSQL database.
// Mirrors the .NET EmployeeInformationController service dependencies.
type erpEmployeeService struct {
	erpRepo  *repository.ErpRepository
	staffRepo *repository.StaffRepository
	db       *gorm.DB
	cfg      *config.Config
	log      zerolog.Logger
}

// newErpEmployeeService creates a new ErpEmployeeService with all required dependencies.
func newErpEmployeeService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) ErpEmployeeService {
	return &erpEmployeeService{
		erpRepo:   repos.Erp,
		staffRepo: repos.Staff,
		db:        repos.GormDB,
		cfg:       cfg,
		log:       log.With().Str("service", "erp_employee").Logger(),
	}
}

// ---------------------------------------------------------------------------
// Organogram lookups -- ERP SQL Server queries
// ---------------------------------------------------------------------------

// GetAllDepartments retrieves all distinct departments from the ERP database.
// Mirrors .NET GetAllDepartmentsQuery handler.
func (s *erpEmployeeService) GetAllDepartments(ctx context.Context) (interface{}, error) {
	s.log.Debug().Msg("fetching all departments from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	depts, err := s.erpRepo.AllDepartments(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve departments from ERP")
		return nil, fmt.Errorf("retrieving ERP departments: %w", err)
	}

	s.log.Debug().Int("count", len(depts)).Msg("ERP departments retrieved")
	return depts, nil
}

// GetAllDivisions retrieves all distinct divisions from the ERP database,
// optionally filtered by department ID.
// Mirrors .NET GetAllDivisionsQuery handler.
func (s *erpEmployeeService) GetAllDivisions(ctx context.Context, departmentId *int) (interface{}, error) {
	s.log.Debug().Msg("fetching all divisions from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	allDivisions, err := s.erpRepo.AllDivisions(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve divisions from ERP")
		return nil, fmt.Errorf("retrieving ERP divisions: %w", err)
	}

	// Filter by department ID if provided.
	if departmentId != nil {
		filtered := make([]erp.ErpOrganizationVm, 0)
		for _, d := range allDivisions {
			if d.DepartmentID != nil && *d.DepartmentID == *departmentId {
				filtered = append(filtered, d)
			}
		}
		s.log.Debug().Int("count", len(filtered)).Int("departmentId", *departmentId).Msg("ERP divisions filtered by department")
		return filtered, nil
	}

	s.log.Debug().Int("count", len(allDivisions)).Msg("ERP divisions retrieved")
	return allDivisions, nil
}

// GetAllOffices retrieves all distinct offices from the ERP database,
// optionally filtered by division ID.
// Mirrors .NET GetAllOfficesQuery handler.
func (s *erpEmployeeService) GetAllOffices(ctx context.Context, divisionId *int) (interface{}, error) {
	s.log.Debug().Msg("fetching all offices from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	allOffices, err := s.erpRepo.AllOffices(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve offices from ERP")
		return nil, fmt.Errorf("retrieving ERP offices: %w", err)
	}

	// Filter by division ID if provided.
	if divisionId != nil {
		filtered := make([]erp.ErpOrganizationVm, 0)
		for _, o := range allOffices {
			if o.DivisionID != nil && *o.DivisionID == *divisionId {
				filtered = append(filtered, o)
			}
		}
		s.log.Debug().Int("count", len(filtered)).Int("divisionId", *divisionId).Msg("ERP offices filtered by division")
		return filtered, nil
	}

	s.log.Debug().Int("count", len(allOffices)).Msg("ERP offices retrieved")
	return allOffices, nil
}

// ---------------------------------------------------------------------------
// Employee lookups -- ERP SQL Server queries
// ---------------------------------------------------------------------------

// GetEmployeeDetail retrieves detailed employee information by employee number.
// Returns an *erp.EmployeeData populated from the ERP EmployeeDetails view.
// Mirrors .NET GetEmployeeDetailQuery handler.
func (s *erpEmployeeService) GetEmployeeDetail(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching employee detail from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emp, err := s.erpRepo.GetEmployeeByID(ctx, employeeNumber)
	if err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve employee from ERP")
		return nil, fmt.Errorf("retrieving employee %s from ERP: %w", employeeNumber, err)
	}
	if emp == nil {
		return nil, nil
	}

	// Map EmployeeDetails to EmployeeData (the DTO used by downstream services).
	data := s.mapEmployeeDetailsToData(emp)

	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("employee detail retrieved")
	return data, nil
}

// GetHeadSubordinates retrieves all employees who report to the given employee
// as head of office, division, or department.
// Mirrors .NET GetHeadSubordinatesQuery handler.
func (s *erpEmployeeService) GetHeadSubordinates(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching head subordinates from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	// First get the employee to find their org position.
	emp, err := s.erpRepo.GetEmployeeByID(ctx, employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("retrieving employee %s: %w", employeeNumber, err)
	}
	if emp == nil {
		return nil, fmt.Errorf("employee not found: %s", employeeNumber)
	}

	// Collect all employees where this person is their head.
	// The .NET approach gathers subordinates from office, division, and department
	// where the employee is the head of that unit.
	allEmployees, err := s.erpRepo.GetAllActiveEmployees(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving all employees: %w", err)
	}

	seen := make(map[string]bool)
	subordinates := make([]erp.EmployeeData, 0)

	for i := range allEmployees {
		e := &allEmployees[i]
		if e.EmployeeNumber == employeeNumber {
			continue
		}
		if seen[e.EmployeeNumber] {
			continue
		}

		isSubordinate := false

		// Check if this employee's head of office/division/department matches.
		if e.HeadOfOfficeID == employeeNumber ||
			e.HeadOfDivID == employeeNumber ||
			e.HeadOfDeptID == employeeNumber {
			isSubordinate = true
		}

		if isSubordinate {
			seen[e.EmployeeNumber] = true
			subordinates = append(subordinates, *s.mapEmployeeDetailsToData(e))
		}
	}

	s.log.Debug().
		Str("employeeNumber", employeeNumber).
		Int("count", len(subordinates)).
		Msg("head subordinates retrieved")

	return subordinates, nil
}

// GetEmployeeSubordinates retrieves direct subordinates (by supervisor ID) for a given employee.
// Mirrors .NET GetEmployeeSubordinatesQuery handler.
func (s *erpEmployeeService) GetEmployeeSubordinates(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching employee subordinates from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emps, err := s.erpRepo.GetSubordinates(ctx, employeeNumber)
	if err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve subordinates from ERP")
		return nil, fmt.Errorf("retrieving subordinates for %s: %w", employeeNumber, err)
	}

	results := make([]erp.EmployeeData, 0, len(emps))
	for i := range emps {
		results = append(results, *s.mapEmployeeDetailsToData(&emps[i]))
	}

	s.log.Debug().
		Str("employeeNumber", employeeNumber).
		Int("count", len(results)).
		Msg("employee subordinates retrieved")

	return results, nil
}

// GetEmployeePeers retrieves peers for a given employee (same office and grade).
// Mirrors .NET GetEmployeePeersQuery handler.
func (s *erpEmployeeService) GetEmployeePeers(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching employee peers from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	// First get the employee to determine their office and grade.
	emp, err := s.erpRepo.GetEmployeeByID(ctx, employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("retrieving employee %s: %w", employeeNumber, err)
	}
	if emp == nil {
		return nil, fmt.Errorf("employee not found: %s", employeeNumber)
	}

	officeID := 0
	if emp.OfficeID != nil {
		officeID = *emp.OfficeID
	}

	if officeID == 0 || emp.Grade == "" {
		s.log.Warn().Str("employeeNumber", employeeNumber).Msg("employee has no office or grade, returning empty peers list")
		return []erp.EmployeeData{}, nil
	}

	peers, err := s.erpRepo.GetEmployeesByOfficeAndGrade(ctx, employeeNumber, emp.Grade, officeID)
	if err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve peers from ERP")
		return nil, fmt.Errorf("retrieving peers for %s: %w", employeeNumber, err)
	}

	results := make([]erp.EmployeeData, 0, len(peers))
	for i := range peers {
		results = append(results, *s.mapEmployeeDetailsToData(&peers[i]))
	}

	s.log.Debug().
		Str("employeeNumber", employeeNumber).
		Int("count", len(results)).
		Msg("employee peers retrieved")

	return results, nil
}

// GetAllByDepartmentId retrieves all active employees in a department.
// Mirrors .NET GetAllByDepartmentIdQuery handler.
func (s *erpEmployeeService) GetAllByDepartmentId(ctx context.Context, departmentId int) (interface{}, error) {
	s.log.Debug().Int("departmentId", departmentId).Msg("fetching employees by department from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emps, err := s.erpRepo.GetByDepartmentID(ctx, departmentId)
	if err != nil {
		s.log.Error().Err(err).Int("departmentId", departmentId).Msg("failed to retrieve employees by department")
		return nil, fmt.Errorf("retrieving employees for department %d: %w", departmentId, err)
	}

	results := make([]erp.EmployeeData, 0, len(emps))
	for i := range emps {
		results = append(results, *s.mapEmployeeDetailsToData(&emps[i]))
	}

	s.log.Debug().Int("departmentId", departmentId).Int("count", len(results)).Msg("employees by department retrieved")
	return results, nil
}

// GetAllByDivisionId retrieves all active employees in a division.
// Mirrors .NET GetAllByDivisionIdQuery handler.
func (s *erpEmployeeService) GetAllByDivisionId(ctx context.Context, divisionId int) (interface{}, error) {
	s.log.Debug().Int("divisionId", divisionId).Msg("fetching employees by division from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emps, err := s.erpRepo.GetByDivisionID(ctx, divisionId)
	if err != nil {
		s.log.Error().Err(err).Int("divisionId", divisionId).Msg("failed to retrieve employees by division")
		return nil, fmt.Errorf("retrieving employees for division %d: %w", divisionId, err)
	}

	results := make([]erp.EmployeeData, 0, len(emps))
	for i := range emps {
		results = append(results, *s.mapEmployeeDetailsToData(&emps[i]))
	}

	s.log.Debug().Int("divisionId", divisionId).Int("count", len(results)).Msg("employees by division retrieved")
	return results, nil
}

// GetAllByOfficeId retrieves all active employees in an office.
// Mirrors .NET GetAllByOfficeIdQuery handler.
func (s *erpEmployeeService) GetAllByOfficeId(ctx context.Context, officeId int) (interface{}, error) {
	s.log.Debug().Int("officeId", officeId).Msg("fetching employees by office from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emps, err := s.erpRepo.GetByOfficeID(ctx, officeId)
	if err != nil {
		s.log.Error().Err(err).Int("officeId", officeId).Msg("failed to retrieve employees by office")
		return nil, fmt.Errorf("retrieving employees for office %d: %w", officeId, err)
	}

	results := make([]erp.EmployeeData, 0, len(emps))
	for i := range emps {
		results = append(results, *s.mapEmployeeDetailsToData(&emps[i]))
	}

	s.log.Debug().Int("officeId", officeId).Int("count", len(results)).Msg("employees by office retrieved")
	return results, nil
}

// GetAllEmployees retrieves all active employees from the ERP database.
// Mirrors .NET GetAllEmployeesQuery handler.
func (s *erpEmployeeService) GetAllEmployees(ctx context.Context) (interface{}, error) {
	s.log.Debug().Msg("fetching all employees from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	emps, err := s.erpRepo.GetAllActiveEmployees(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve all employees from ERP")
		return nil, fmt.Errorf("retrieving all employees: %w", err)
	}

	results := make([]erp.EmployeeData, 0, len(emps))
	for i := range emps {
		results = append(results, *s.mapEmployeeDetailsToData(&emps[i]))
	}

	s.log.Debug().Int("count", len(results)).Msg("all employees retrieved")
	return results, nil
}

// SeedOrganizationData seeds the local PMS organogram tables (departments, divisions,
// offices) from the ERP database. This ensures the PMS has an up-to-date copy of
// the organisational structure.
// Mirrors .NET SeedOrganizationDataCommand handler.
func (s *erpEmployeeService) SeedOrganizationData(ctx context.Context) (interface{}, error) {
	s.log.Info().Msg("seeding organization data from ERP")

	if s.erpRepo == nil {
		return nil, fmt.Errorf("ERP database not configured")
	}

	// Fetch organisation units from ERP.
	erpDepts, err := s.erpRepo.AllDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching ERP departments for seeding: %w", err)
	}

	erpDivisions, err := s.erpRepo.AllDivisions(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching ERP divisions for seeding: %w", err)
	}

	erpOffices, err := s.erpRepo.AllOffices(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching ERP offices for seeding: %w", err)
	}

	deptCount := 0
	divCount := 0
	officeCount := 0

	// Seed departments: upsert by department_id.
	for _, d := range erpDepts {
		if d.DepartmentID == nil {
			continue
		}
		var existing organogram.Department
		err := s.db.WithContext(ctx).
			Where("department_id = ?", *d.DepartmentID).
			First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			dept := organogram.Department{
				DepartmentID:   *d.DepartmentID,
				DepartmentName: d.DepartmentName,
				DepartmentCode: strings.ToUpper(d.DepartmentName[:min(5, len(d.DepartmentName))]),
			}
			dept.IsActive = true
			if createErr := s.db.WithContext(ctx).Create(&dept).Error; createErr != nil {
				s.log.Warn().Err(createErr).Int("departmentId", *d.DepartmentID).Msg("failed to seed department")
				continue
			}
			deptCount++
		} else if err == nil {
			// Update name if changed.
			if existing.DepartmentName != d.DepartmentName {
				existing.DepartmentName = d.DepartmentName
				now := time.Now().UTC()
				existing.DateUpdated = &now
				s.db.WithContext(ctx).Save(&existing)
			}
		}
	}

	// Seed divisions: upsert by division_id.
	for _, d := range erpDivisions {
		if d.DivisionID == nil {
			continue
		}
		deptID := 0
		if d.DepartmentID != nil {
			deptID = *d.DepartmentID
		}
		var existing organogram.Division
		err := s.db.WithContext(ctx).
			Where("division_id = ?", *d.DivisionID).
			First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			div := organogram.Division{
				DivisionID:   *d.DivisionID,
				DepartmentID: deptID,
				DivisionName: d.DivisionName,
				DivisionCode: strings.ToUpper(d.DivisionName[:min(5, len(d.DivisionName))]),
			}
			div.IsActive = true
			if createErr := s.db.WithContext(ctx).Create(&div).Error; createErr != nil {
				s.log.Warn().Err(createErr).Int("divisionId", *d.DivisionID).Msg("failed to seed division")
				continue
			}
			divCount++
		} else if err == nil {
			if existing.DivisionName != d.DivisionName {
				existing.DivisionName = d.DivisionName
				now := time.Now().UTC()
				existing.DateUpdated = &now
				s.db.WithContext(ctx).Save(&existing)
			}
		}
	}

	// Seed offices: upsert by office_id.
	for _, o := range erpOffices {
		if o.OfficeID == 0 {
			continue
		}
		divID := 0
		if o.DivisionID != nil {
			divID = *o.DivisionID
		}
		var existing organogram.Office
		err := s.db.WithContext(ctx).
			Where("office_id = ?", o.OfficeID).
			First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			office := organogram.Office{
				OfficeID:   o.OfficeID,
				DivisionID: divID,
				OfficeName: o.OfficeName,
				OfficeCode: strings.ToUpper(o.OfficeName[:min(5, len(o.OfficeName))]),
			}
			office.IsActive = true
			if createErr := s.db.WithContext(ctx).Create(&office).Error; createErr != nil {
				s.log.Warn().Err(createErr).Int("officeId", o.OfficeID).Msg("failed to seed office")
				continue
			}
			officeCount++
		} else if err == nil {
			if existing.OfficeName != o.OfficeName {
				existing.OfficeName = o.OfficeName
				now := time.Now().UTC()
				existing.DateUpdated = &now
				s.db.WithContext(ctx).Save(&existing)
			}
		}
	}

	s.log.Info().
		Int("departments", deptCount).
		Int("divisions", divCount).
		Int("offices", officeCount).
		Msg("organization data seeded from ERP")

	return map[string]interface{}{
		"isSuccess":   true,
		"message":     "Organization data seeded successfully",
		"departments": deptCount,
		"divisions":   divCount,
		"offices":     officeCount,
	}, nil
}

// GetStaffIDMaskDetail retrieves staff ID mask data for a given employee number.
// Mirrors .NET GetStaffIDMaskDetailQuery handler. Uses the StaffIDSQL database.
func (s *erpEmployeeService) GetStaffIDMaskDetail(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching staff ID mask detail")

	if s.staffRepo == nil {
		return nil, fmt.Errorf("StaffIDMask database not configured")
	}

	mask, err := s.staffRepo.GetStaffIDMask(ctx, employeeNumber)
	if err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve staff ID mask detail")
		return nil, fmt.Errorf("retrieving staff ID mask for %s: %w", employeeNumber, err)
	}

	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("staff ID mask detail retrieved")
	return mask, nil
}

// ---------------------------------------------------------------------------
// Job role management -- PostgreSQL GORM queries on StaffJobRoles
// ---------------------------------------------------------------------------

// UpdateStaffJobRole creates or updates a staff job role record.
// The request is expected to be a *competency.StaffJobRoles or a map with the
// relevant fields. Mirrors .NET UpdateStaffJobRoleCommand handler.
func (s *erpEmployeeService) UpdateStaffJobRole(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Info().Msg("updating staff job role")

	sjr, ok := req.(*competency.StaffJobRoles)
	if !ok {
		// Attempt map-based request for handler flexibility.
		if m, mok := req.(map[string]interface{}); mok {
			sjr = s.mapToStaffJobRole(m)
		} else {
			return nil, fmt.Errorf("invalid request type: expected *competency.StaffJobRoles or map[string]interface{}")
		}
	}

	if sjr.EmployeeID == "" {
		return nil, fmt.Errorf("employee_id is required")
	}

	// Check if a record already exists for this employee.
	var existing competency.StaffJobRoles
	err := s.db.WithContext(ctx).
		Where("employee_id = ? AND soft_deleted = ?", sjr.EmployeeID, false).
		First(&existing).Error

	if err == nil {
		// Update existing record.
		existing.FullName = sjr.FullName
		existing.DepartmentID = sjr.DepartmentID
		existing.DivisionID = sjr.DivisionID
		existing.OfficeID = sjr.OfficeID
		existing.SupervisorID = sjr.SupervisorID
		existing.JobRoleID = sjr.JobRoleID
		existing.JobRoleName = sjr.JobRoleName
		// Reset approval status on update (pending re-approval).
		existing.IsApproved = false
		existing.IsRejected = false
		existing.Status = "PENDING"

		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			s.log.Error().Err(err).Str("employeeId", sjr.EmployeeID).Msg("failed to update staff job role")
			return nil, fmt.Errorf("updating staff job role: %w", err)
		}

		s.log.Info().Str("employeeId", sjr.EmployeeID).Msg("staff job role updated")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   "Staff job role has been updated successfully",
		}, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("querying existing staff job role: %w", err)
	}

	// Create new record.
	sjr.Status = "PENDING"
	sjr.IsActive = true
	if err := s.db.WithContext(ctx).Create(sjr).Error; err != nil {
		s.log.Error().Err(err).Str("employeeId", sjr.EmployeeID).Msg("failed to create staff job role")
		return nil, fmt.Errorf("creating staff job role: %w", err)
	}

	s.log.Info().Str("employeeId", sjr.EmployeeID).Msg("staff job role created")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   "Staff job role has been created successfully",
	}, nil
}

// GetStaffJobRoleById retrieves the staff job role for a specific employee.
// Mirrors .NET GetStaffJobRoleByIdQuery handler.
func (s *erpEmployeeService) GetStaffJobRoleById(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching staff job role by employee ID")

	var sjr competency.StaffJobRoles
	err := s.db.WithContext(ctx).
		Where("employee_id = ? AND soft_deleted = ?", employeeNumber, false).
		First(&sjr).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Debug().Str("employeeNumber", employeeNumber).Msg("no staff job role found")
			return nil, nil
		}
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve staff job role")
		return nil, fmt.Errorf("retrieving staff job role for %s: %w", employeeNumber, err)
	}

	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("staff job role retrieved")
	return &sjr, nil
}

// GetJobRolesByOffice retrieves all staff job roles for a given office.
// The request is expected to contain an officeId field.
// Mirrors .NET GetJobRolesByOfficeQuery handler.
func (s *erpEmployeeService) GetJobRolesByOffice(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Debug().Msg("fetching job roles by office")

	officeID := 0
	switch v := req.(type) {
	case map[string]interface{}:
		if id, ok := v["officeId"]; ok {
			switch oid := id.(type) {
			case float64:
				officeID = int(oid)
			case int:
				officeID = oid
			}
		}
	case *struct{ OfficeID int }:
		officeID = v.OfficeID
	}

	if officeID == 0 {
		return nil, fmt.Errorf("officeId is required")
	}

	var roles []competency.StaffJobRoles
	err := s.db.WithContext(ctx).
		Where("office_id = ? AND soft_deleted = ?", officeID, false).
		Find(&roles).Error
	if err != nil {
		s.log.Error().Err(err).Int("officeId", officeID).Msg("failed to retrieve job roles by office")
		return nil, fmt.Errorf("retrieving job roles for office %d: %w", officeID, err)
	}

	s.log.Debug().Int("officeId", officeID).Int("count", len(roles)).Msg("job roles by office retrieved")
	return roles, nil
}

// GetStaffJobRoleRequests retrieves all job role update requests for a given employee
// (including pending approval records).
// Mirrors .NET GetStaffJobRoleRequestsQuery handler.
func (s *erpEmployeeService) GetStaffJobRoleRequests(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Debug().Str("employeeNumber", employeeNumber).Msg("fetching staff job role requests")

	var roles []competency.StaffJobRoles
	err := s.db.WithContext(ctx).
		Where("(employee_id = ? OR supervisor_id = ?) AND soft_deleted = ?",
			employeeNumber, employeeNumber, false).
		Order("date_created DESC").
		Find(&roles).Error
	if err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("failed to retrieve staff job role requests")
		return nil, fmt.Errorf("retrieving job role requests for %s: %w", employeeNumber, err)
	}

	s.log.Debug().Str("employeeNumber", employeeNumber).Int("count", len(roles)).Msg("staff job role requests retrieved")
	return roles, nil
}

// ApproveRejectStaffJobRole approves or rejects a staff job role update request.
// The request is expected to contain staffJobRoleId, isApproved (bool), and optionally
// rejectionReason. Mirrors .NET ApproveRejectStaffJobRoleCommand handler.
func (s *erpEmployeeService) ApproveRejectStaffJobRole(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Info().Msg("approving/rejecting staff job role")

	m, ok := req.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected map[string]interface{}")
	}

	staffJobRoleID := 0
	if id, exists := m["staffJobRoleId"]; exists {
		switch v := id.(type) {
		case float64:
			staffJobRoleID = int(v)
		case int:
			staffJobRoleID = v
		}
	}
	if staffJobRoleID == 0 {
		return nil, fmt.Errorf("staffJobRoleId is required")
	}

	isApproved := false
	if val, exists := m["isApproved"]; exists {
		if b, bOk := val.(bool); bOk {
			isApproved = b
		}
	}

	rejectionReason := ""
	if val, exists := m["rejectionReason"]; exists {
		if rs, sOk := val.(string); sOk {
			rejectionReason = rs
		}
	}

	approvedBy := ""
	if val, exists := m["approvedBy"]; exists {
		if ab, sOk := val.(string); sOk {
			approvedBy = ab
		}
	}

	var sjr competency.StaffJobRoles
	err := s.db.WithContext(ctx).
		Where("staff_job_role_id = ? AND soft_deleted = ?", staffJobRoleID, false).
		First(&sjr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("staff job role with id %d not found", staffJobRoleID)
		}
		return nil, fmt.Errorf("retrieving staff job role: %w", err)
	}

	now := time.Now().UTC()

	if isApproved {
		sjr.IsApproved = true
		sjr.IsRejected = false
		sjr.ApprovedBy = approvedBy
		sjr.DateApproved = &now
		sjr.Status = "APPROVED"

		if err := s.db.WithContext(ctx).Save(&sjr).Error; err != nil {
			return nil, fmt.Errorf("approving staff job role: %w", err)
		}

		s.log.Info().Int("staffJobRoleId", staffJobRoleID).Msg("staff job role approved")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   "Staff job role has been approved successfully",
		}, nil
	}

	// Reject
	sjr.IsRejected = true
	sjr.IsApproved = false
	sjr.RejectedBy = approvedBy
	sjr.RejectionReason = rejectionReason
	sjr.DateRejected = &now
	sjr.Status = "REJECTED"

	if err := s.db.WithContext(ctx).Save(&sjr).Error; err != nil {
		return nil, fmt.Errorf("rejecting staff job role: %w", err)
	}

	s.log.Info().Int("staffJobRoleId", staffJobRoleID).Msg("staff job role rejected")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   "Staff job role has been rejected",
	}, nil
}

// ===========================================================================
// Private helpers
// ===========================================================================

// mapEmployeeDetailsToData converts an ERP EmployeeDetails record into an
// EmployeeData DTO suitable for API consumption and downstream service usage.
func (s *erpEmployeeService) mapEmployeeDetailsToData(e *erp.EmployeeDetails) *erp.EmployeeData {
	officeID := 0
	if e.OfficeID != nil {
		officeID = *e.OfficeID
	}
	locationID := 0
	if e.LocationID != nil {
		locationID = *e.LocationID
	}

	status := "Active"
	if e.PersonTypeID != repository.ActiveStaffPersonType {
		status = "Inactive"
	}

	return &erp.EmployeeData{
		EmployeeErpDetailsDTO: erp.EmployeeErpDetailsDTO{
			UserName:       e.EmployeeNumber, // ERP view does not expose UserName; use EmployeeNumber as fallback
			EmailAddress:   e.Email,
			FirstName:      e.FirstName,
			LastName:       e.LastName,
			EmployeeNumber: e.EmployeeNumber,
			JobName:        e.JobName,
			DepartmentName: e.Department,
			DivisionName:   e.Division,
			OfficeName:     e.Office,
			SupervisorID:   e.SupervisorID,
			HeadOfOfficeID: e.HeadOfOfficeID,
			HeadOfDivID:    e.HeadOfDivID,
			HeadOfDeptID:   e.HeadOfDeptID,
			DepartmentID:   e.DepartmentID,
			OfficeID:       officeID,
			Grade:          e.Grade,
			DivisionID:     e.DivisionID,
			Position:       e.JobTitle,
		},
		Status:       status,
		PersonTypeID: e.PersonTypeID,
		LocationID:   locationID,
	}
}

// mapToStaffJobRole creates a StaffJobRoles from a map[string]interface{} request body.
func (s *erpEmployeeService) mapToStaffJobRole(m map[string]interface{}) *competency.StaffJobRoles {
	sjr := &competency.StaffJobRoles{}

	if v, ok := m["employeeId"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.EmployeeID = s
		}
	}
	if v, ok := m["employee_id"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.EmployeeID = s
		}
	}
	if v, ok := m["fullName"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.FullName = s
		}
	}
	if v, ok := m["full_name"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.FullName = s
		}
	}
	if v, ok := m["departmentId"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.DepartmentID = int(f)
		}
	}
	if v, ok := m["department_id"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.DepartmentID = int(f)
		}
	}
	if v, ok := m["divisionId"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.DivisionID = int(f)
		}
	}
	if v, ok := m["division_id"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.DivisionID = int(f)
		}
	}
	if v, ok := m["officeId"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.OfficeID = int(f)
		}
	}
	if v, ok := m["office_id"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.OfficeID = int(f)
		}
	}
	if v, ok := m["supervisorId"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.SupervisorID = s
		}
	}
	if v, ok := m["supervisor_id"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.SupervisorID = s
		}
	}
	if v, ok := m["jobRoleId"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.JobRoleID = int(f)
		}
	}
	if v, ok := m["job_role_id"]; ok {
		if f, fOk := v.(float64); fOk {
			sjr.JobRoleID = int(f)
		}
	}
	if v, ok := m["jobRoleName"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.JobRoleName = s
		}
	}
	if v, ok := m["job_role_name"]; ok {
		if s, sOk := v.(string); sOk {
			sjr.JobRoleName = s
		}
	}

	return sjr
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	// Compile-time interface compliance check.
	var _ ErpEmployeeService = (*erpEmployeeService)(nil)
}
