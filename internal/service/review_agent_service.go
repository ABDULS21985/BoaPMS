package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// pmGrade is the special grade code for Permanent Member (Grade 41).
// Employees at this grade have special subordinate/superior filtering rules.
const pmGrade = "41"

// reviewAgentService is an internal helper used by CompetencyService.
// It handles random reviewer selection (360-degree reviews), review population
// (creating review records for employees), and review average calculation.
//
// Converted from .NET ReviewAgentService (1,415 lines).
type reviewAgentService struct {
	repos       *repository.Container
	cfg         *config.Config
	log         zerolog.Logger
	db          *gorm.DB // shortcut to repos.GormDB (competency core DB)
	reviewRepo  *repository.Repository[competency.CompetencyReview]
	profileRepo *repository.Repository[competency.CompetencyReviewProfile]
	gradingRepo *repository.Repository[competency.CompetencyCategoryGrading]
}

// newReviewAgentService creates a reviewAgentService wired up with repository
// access for competency reviews, profiles, and grading weights.
func newReviewAgentService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) *reviewAgentService {
	return &reviewAgentService{
		repos:       repos,
		cfg:         cfg,
		log:         log.With().Str("service", "review_agent").Logger(),
		db:          repos.GormDB,
		reviewRepo:  repository.NewRepository[competency.CompetencyReview](repos.GormDB),
		profileRepo: repository.NewRepository[competency.CompetencyReviewProfile](repos.GormDB),
		gradingRepo: repository.NewRepository[competency.CompetencyCategoryGrading](repos.GormDB),
	}
}

// ---------------------------------------------------------------------------
// ERP data helpers — queries against the ERP SQL Server (repos.ErpSQL).
// These replace the .NET ErpEmployeeService calls used by ReviewAgentService.
// ---------------------------------------------------------------------------

// getEmployeeDetail retrieves a single employee's details from the ERP database.
func (s *reviewAgentService) getEmployeeDetail(ctx context.Context, employeeNumber string) (*erp.EmployeeErpDetailsDTO, error) {
	if strings.TrimSpace(employeeNumber) == "" {
		return nil, nil
	}
	const q = `SELECT TOP 1
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position, PersonId AS personId
	FROM dbo.EmployeeDetails
	WHERE EmployeeNumber = @p1 AND PersonTypeId = 1120`

	emp, err := repository.RawQuerySingle[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, employeeNumber)
	if err != nil || emp == nil {
		return nil, err
	}
	// Mirror .NET: Position = Position.Split(".").FirstOrDefault() ?? JobName
	if strings.TrimSpace(emp.Position) != "" {
		parts := strings.SplitN(emp.Position, ".", 2)
		emp.Position = parts[0]
	} else {
		emp.Position = emp.JobName
	}
	return emp, nil
}

// getGovernorDGs retrieves employees whose JobName contains "GOVERNOR".
func (s *reviewAgentService) getGovernorDGs(ctx context.Context) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE JobName LIKE '%GOVERNOR%' AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q)
}

// getAllHeadDepartments retrieves employees who are heads of their own department
// (HeadOfDeptId == EmployeeNumber) excluding governors.
func (s *reviewAgentService) getAllHeadDepartments(ctx context.Context) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE PersonTypeId = 1120 AND HeadOfDeptId = EmployeeNumber
		AND JobName NOT LIKE '%GOVERNOR%'`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q)
}

// getSubordinatesByOfficeAndGrades returns employees in the same office with a
// higher grade number (lower rank) than the given grade.
func (s *reviewAgentService) getSubordinatesByOfficeAndGrades(ctx context.Context, empNumber, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE OfficeId = @p1 AND EmployeeNumber != @p2
		AND CAST(Grade AS INT) > CAST(@p3 AS INT)
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, officeID, empNumber, grade)
}

// getSubordinatesByDivisionAndGrades returns employees in the same division with
// a higher grade number than the given grade.
func (s *reviewAgentService) getSubordinatesByDivisionAndGrades(ctx context.Context, empNumber, grade string, divisionID *int) ([]erp.EmployeeDetails, error) {
	if divisionID == nil {
		return nil, nil
	}
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DivisionId = @p1 AND EmployeeNumber != @p2
		AND CAST(Grade AS INT) > CAST(@p3 AS INT)
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, *divisionID, empNumber, grade)
}

// getSubordinatesByDepartmentAndGrades returns employees in the same department
// with a higher grade number.
func (s *reviewAgentService) getSubordinatesByDepartmentAndGrades(ctx context.Context, empNumber, grade string, departmentID *int) ([]erp.EmployeeDetails, error) {
	if departmentID == nil {
		return nil, nil
	}
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DepartmentId = @p1 AND EmployeeNumber != @p2
		AND CAST(Grade AS INT) > CAST(@p3 AS INT)
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, *departmentID, empNumber, grade)
}

// getPeersByOfficeAndGrades returns employees in the same office with the exact
// same grade.
func (s *reviewAgentService) getPeersByOfficeAndGrades(ctx context.Context, empNumber, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE OfficeId = @p1 AND Grade = @p2
		AND EmployeeNumber != @p3 AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, officeID, grade, empNumber)
}

// getPeersByDivisionAndGrades returns employees in the same division with the
// exact same grade.
func (s *reviewAgentService) getPeersByDivisionAndGrades(ctx context.Context, empNumber, grade string, divisionID int) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DivisionId = @p1 AND Grade = @p2
		AND EmployeeNumber != @p3 AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, divisionID, grade, empNumber)
}

// getPeersByDepartmentAndGrades returns employees in the same department with
// the exact same grade.
func (s *reviewAgentService) getPeersByDepartmentAndGrades(ctx context.Context, empNumber, grade string, departmentID int) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DepartmentId = @p1 AND Grade = @p2
		AND EmployeeNumber != @p3 AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, departmentID, grade, empNumber)
}

// getSuperiorsByOfficeAndGrades returns employees in the same office with a lower
// grade number (higher rank).
func (s *reviewAgentService) getSuperiorsByOfficeAndGrades(ctx context.Context, empNumber, grade string, officeID int) ([]erp.EmployeeDetails, error) {
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE OfficeId = @p1 AND EmployeeNumber != @p2
		AND CAST(Grade AS INT) < CAST(@p3 AS INT)
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, officeID, empNumber, grade)
}

// getSuperiorsByDivisionAndGrades returns employees in the same division with a
// lower grade number.
func (s *reviewAgentService) getSuperiorsByDivisionAndGrades(ctx context.Context, empNumber, grade string, divisionID *int) ([]erp.EmployeeDetails, error) {
	if divisionID == nil {
		return nil, nil
	}
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DivisionId = @p1 AND EmployeeNumber != @p2
		AND CAST(Grade AS INT) < CAST(@p3 AS INT)
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, *divisionID, empNumber, grade)
}

// getSuperiorsByDepartmentAndGrades returns employees in the same department with
// a lower grade number or grade "41" (PM grade).
func (s *reviewAgentService) getSuperiorsByDepartmentAndGrades(ctx context.Context, empNumber, grade string, departmentID *int) ([]erp.EmployeeDetails, error) {
	if departmentID == nil {
		return nil, nil
	}
	const q = `SELECT * FROM dbo.EmployeeDetails
		WHERE DepartmentId = @p1 AND EmployeeNumber != @p2
		AND (CAST(Grade AS INT) < CAST(@p3 AS INT) OR Grade = '41')
		AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeDetails](s.repos.ErpSQL, ctx, q, *departmentID, empNumber, grade)
}

// getAllEmployees returns all active employees from the ERP database.
func (s *reviewAgentService) getAllEmployees(ctx context.Context) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE EmployeeNumber IS NOT NULL AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q)
}

// getEmployeesByOfficeID returns all active employees for a given office.
func (s *reviewAgentService) getEmployeesByOfficeID(ctx context.Context, officeID int) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE OfficeId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, officeID)
}

// getEmployeesByDivisionID returns all active employees for a given division.
func (s *reviewAgentService) getEmployeesByDivisionID(ctx context.Context, divisionID int) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE DivisionId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, divisionID)
}

// getEmployeesByDepartmentID returns all active employees for a given department.
func (s *reviewAgentService) getEmployeesByDepartmentID(ctx context.Context, departmentID int) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE DepartmentId = @p1 AND EmployeeNumber IS NOT NULL AND PersonTypeId = 1120`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, departmentID)
}

// ---------------------------------------------------------------------------
// Competency data helpers — queries against the GORM (PostgreSQL) database.
// ---------------------------------------------------------------------------

// getCurrentReviewPeriod returns the currently active and approved review period.
func (s *reviewAgentService) getCurrentReviewPeriod(ctx context.Context) (*competency.ReviewPeriod, error) {
	var rp competency.ReviewPeriod
	err := s.db.WithContext(ctx).
		Where("is_active = ? AND is_approved = ? AND soft_deleted = ?", true, true, false).
		First(&rp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &rp, err
}

// getReviewTypes returns all active review types.
func (s *reviewAgentService) getReviewTypes(ctx context.Context) ([]competency.ReviewType, error) {
	var rts []competency.ReviewType
	err := s.db.WithContext(ctx).
		Where("is_active = ? AND soft_deleted = ?", true, false).
		Find(&rts).Error
	return rts, err
}

// getAssignedJobGradeGroup returns the grade group name for a given employee grade code.
func (s *reviewAgentService) getAssignedJobGradeGroup(ctx context.Context, gradeCode string) (*competency.AssignJobGradeGroup, error) {
	var agg competency.AssignJobGradeGroup
	err := s.db.WithContext(ctx).
		Joins("JOIN \"CoreSchema\".\"job_grades\" jg ON jg.job_grade_id = \"CoreSchema\".\"assign_job_grade_groups\".job_grade_id").
		Joins("JOIN \"CoreSchema\".\"job_grade_groups\" jgg ON jgg.job_grade_group_id = \"CoreSchema\".\"assign_job_grade_groups\".job_grade_group_id").
		Where("jg.grade_code = ? AND \"CoreSchema\".\"assign_job_grade_groups\".soft_deleted = ?", gradeCode, false).
		Preload("JobGradeGroup").
		First(&agg).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &agg, err
}

// getOfficeByID returns an organogram office by its ID.
func (s *reviewAgentService) getOfficeByID(ctx context.Context, officeID int) (*competency.OfficeJobRole, error) {
	// We only need the office ID itself for downstream use; the .NET code
	// fetches the office just to confirm it exists and get its ID.
	var office struct {
		OfficeID int `gorm:"column:office_id"`
	}
	err := s.db.WithContext(ctx).
		Table("\"CoreSchema\".\"offices\"").
		Where("office_id = ? AND soft_deleted = ?", officeID, false).
		First(&office).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Return a sentinel to indicate found.
	return &competency.OfficeJobRole{OfficeID: office.OfficeID}, nil
}

// getJobRoleByName returns a job role whose name matches (LIKE %name%).
func (s *reviewAgentService) getJobRoleByName(ctx context.Context, name string) (*competency.JobRole, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}
	var jr competency.JobRole
	err := s.db.WithContext(ctx).
		Where("LOWER(job_role_name) LIKE ? AND soft_deleted = ?", "%"+strings.ToLower(name)+"%", false).
		First(&jr).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &jr, err
}

// getBehavioralCompetencies returns all behavioral competency mappings for a
// given grade group name, with their associated rating IDs.
func (s *reviewAgentService) getBehavioralCompetencies(ctx context.Context, gradeGroupName string) ([]competency.BehavioralCompetency, error) {
	var results []competency.BehavioralCompetency
	err := s.db.WithContext(ctx).
		Joins("JOIN \"CoreSchema\".\"job_grade_groups\" jgg ON jgg.job_grade_group_id = \"CoreSchema\".\"behavioral_competencies\".job_grade_group_id").
		Where("jgg.group_name = ? AND \"CoreSchema\".\"behavioral_competencies\".soft_deleted = ?", gradeGroupName, false).
		Find(&results).Error
	return results, err
}

// getRatings returns all active ratings, optionally ordered by value.
func (s *reviewAgentService) getRatings(ctx context.Context) ([]competency.Rating, error) {
	var ratings []competency.Rating
	err := s.db.WithContext(ctx).
		Where("is_active = ? AND soft_deleted = ?", true, false).
		Order("value ASC").
		Find(&ratings).Error
	return ratings, err
}

// getJobRoleCompetencies returns technical competencies for a given office and
// job role.
func (s *reviewAgentService) getJobRoleCompetencies(ctx context.Context, officeID, jobRoleID int) ([]competency.JobRoleCompetency, error) {
	var results []competency.JobRoleCompetency
	err := s.db.WithContext(ctx).
		Preload("JobRole").
		Where("office_id = ? AND job_role_id = ? AND soft_deleted = ?", officeID, jobRoleID, false).
		Find(&results).Error
	return results, err
}

// getJobRoleCompetenciesByOffice returns all technical competencies for an
// office (fallback when exact job role has none).
func (s *reviewAgentService) getJobRoleCompetenciesByOffice(ctx context.Context, officeID int) ([]competency.JobRoleCompetency, error) {
	var results []competency.JobRoleCompetency
	err := s.db.WithContext(ctx).
		Preload("JobRole").
		Where("office_id = ? AND soft_deleted = ?", officeID, false).
		Find(&results).Error
	return results, err
}

// ---------------------------------------------------------------------------
// Deduplication / filtering helpers
// ---------------------------------------------------------------------------

// deduplicateEmployeeDetails removes nils, blanks, self, and duplicates from a
// list of ERP EmployeeDetails.
func deduplicateEmployeeDetails(list []erp.EmployeeDetails, excludeEmpNum string) []erp.EmployeeDetails {
	seen := make(map[string]bool)
	var out []erp.EmployeeDetails
	for i := range list {
		num := strings.TrimSpace(list[i].EmployeeNumber)
		if num == "" {
			continue
		}
		if strings.EqualFold(num, excludeEmpNum) {
			continue
		}
		if seen[strings.ToUpper(num)] {
			continue
		}
		seen[strings.ToUpper(num)] = true
		out = append(out, list[i])
	}
	return out
}

// filterSubordinates applies the .NET subordinate filtering rules:
//   - Remove nulls, blanks, self, head-of-office, PM-grade employees
//   - Deduplicate by employee number
//   - If the employee IS PM grade, only keep subordinates with grade > 4
func filterSubordinates(list []erp.EmployeeDetails, empNumber, headOfOfficeID, empGrade string) []erp.EmployeeDetails {
	seen := make(map[string]bool)
	var out []erp.EmployeeDetails
	hoID := strings.ToUpper(strings.TrimSpace(headOfOfficeID))

	for i := range list {
		num := strings.TrimSpace(list[i].EmployeeNumber)
		if num == "" {
			continue
		}
		upper := strings.ToUpper(num)
		if strings.EqualFold(num, empNumber) {
			continue
		}
		if hoID != "" && upper == hoID {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(list[i].Grade), pmGrade) {
			continue
		}
		if seen[upper] {
			continue
		}
		seen[upper] = true
		out = append(out, list[i])
	}

	// PM-grade employees only see subordinates with numeric grade > 4
	if empGrade == pmGrade {
		var filtered []erp.EmployeeDetails
		for i := range out {
			g, err := strconv.Atoi(strings.TrimSpace(out[i].Grade))
			if err == nil && g > 4 {
				filtered = append(filtered, out[i])
			}
		}
		return filtered
	}
	return out
}

// filterPeers applies the .NET peer filtering rules:
//   - Remove nulls, blanks, self, head-of-office
//   - Only keep employees with the exact same grade
//   - Deduplicate
func filterPeers(list []erp.EmployeeDetails, empNumber, headOfOfficeID, empGrade string) []erp.EmployeeDetails {
	seen := make(map[string]bool)
	var out []erp.EmployeeDetails
	hoID := strings.ToUpper(strings.TrimSpace(headOfOfficeID))

	for i := range list {
		num := strings.TrimSpace(list[i].EmployeeNumber)
		if num == "" {
			continue
		}
		upper := strings.ToUpper(num)
		if strings.EqualFold(num, empNumber) {
			continue
		}
		if hoID != "" && upper == hoID {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(list[i].Grade), empGrade) {
			continue
		}
		if seen[upper] {
			continue
		}
		seen[upper] = true
		out = append(out, list[i])
	}
	return out
}

// filterSuperiors applies the .NET superior filtering rules:
//   - Remove nulls, blanks, self, head-of-office, supervisor
//   - Deduplicate
//   - If the employee IS PM grade, only keep superiors with grade < 4
func filterSuperiors(list []erp.EmployeeDetails, empNumber, headOfOfficeID, supervisorID, empGrade string) []erp.EmployeeDetails {
	seen := make(map[string]bool)
	var out []erp.EmployeeDetails
	hoID := strings.ToUpper(strings.TrimSpace(headOfOfficeID))
	supID := strings.ToUpper(strings.TrimSpace(supervisorID))

	for i := range list {
		num := strings.TrimSpace(list[i].EmployeeNumber)
		if num == "" {
			continue
		}
		upper := strings.ToUpper(num)
		if strings.EqualFold(num, empNumber) {
			continue
		}
		if hoID != "" && upper == hoID {
			continue
		}
		if supID != "" && upper == supID {
			continue
		}
		if seen[upper] {
			continue
		}
		seen[upper] = true
		out = append(out, list[i])
	}

	if empGrade == pmGrade {
		var filtered []erp.EmployeeDetails
		for i := range out {
			g, err := strconv.Atoi(strings.TrimSpace(out[i].Grade))
			if err == nil && g < 4 {
				filtered = append(filtered, out[i])
			}
		}
		return filtered
	}
	return out
}

// pickRandom selects a random element from a non-empty slice.
func pickRandom[T any](list []T) *T {
	if len(list) == 0 {
		return nil
	}
	if len(list) == 1 {
		return &list[0]
	}
	return &list[rand.Intn(len(list))]
}

// ---------------------------------------------------------------------------
// 1. Random Reviewer Selection (for 360-degree reviews)
// ---------------------------------------------------------------------------

// GetRandomEmployeeSubordinate selects a random subordinate for the given
// employee with a multi-level fallback: office -> division -> department.
// Governor/DG employees use a special path (head department subordinates).
func (s *reviewAgentService) GetRandomEmployeeSubordinate(ctx context.Context, employeeNumber string) (*erp.EmployeeDetails, error) {
	if strings.TrimSpace(employeeNumber) == "" {
		return nil, nil
	}

	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil || strings.TrimSpace(emp.Grade) == "" {
		return nil, err
	}

	var subordinates []erp.EmployeeDetails

	// Check if this employee is a Governor/DG
	govDGs, err := s.getGovernorDGs(ctx)
	if err != nil {
		return nil, fmt.Errorf("getGovernorDGs: %w", err)
	}
	isGovernorDG := false
	for _, g := range govDGs {
		if strings.EqualFold(g.EmployeeNumber, employeeNumber) {
			isGovernorDG = true
			break
		}
	}

	if isGovernorDG {
		allHeadDepts, err := s.getAllHeadDepartments(ctx)
		if err != nil {
			return nil, fmt.Errorf("getAllHeadDepartments: %w", err)
		}
		// Filter to those reporting to this employee (HeadOfOfficeId or SupervisorId)
		// Mirrors .NET: x.HeadOfOfficeId == employeeNumber || x.SupervisorId == employeeNumber
		for _, hd := range allHeadDepts {
			if strings.EqualFold(hd.HeadOfOfficeID, employeeNumber) || strings.EqualFold(hd.SupervisorID, employeeNumber) {
				subordinates = append(subordinates, hd)
			}
		}
		subordinates = deduplicateEmployeeDetails(subordinates, employeeNumber)
	} else {
		// 1. Try office
		if emp.OfficeID > 0 {
			officeSubs, err := s.getSubordinatesByOfficeAndGrades(ctx, employeeNumber, emp.Grade, emp.OfficeID)
			if err != nil {
				return nil, fmt.Errorf("getSubordinatesByOfficeAndGrades: %w", err)
			}
			subordinates = filterSubordinates(officeSubs, employeeNumber, emp.HeadOfOfficeID, emp.Grade)
		}

		// 2. Fallback to division
		if len(subordinates) == 0 && emp.DivisionID != nil {
			divSubs, err := s.getSubordinatesByDivisionAndGrades(ctx, employeeNumber, emp.Grade, emp.DivisionID)
			if err != nil {
				return nil, fmt.Errorf("getSubordinatesByDivisionAndGrades: %w", err)
			}
			subordinates = filterSubordinates(divSubs, employeeNumber, emp.HeadOfOfficeID, emp.Grade)
		}

		// 3. Fallback to department
		if len(subordinates) == 0 && emp.DepartmentID != nil {
			deptSubs, err := s.getSubordinatesByDepartmentAndGrades(ctx, employeeNumber, emp.Grade, emp.DepartmentID)
			if err != nil {
				return nil, fmt.Errorf("getSubordinatesByDepartmentAndGrades: %w", err)
			}
			subordinates = filterSubordinates(deptSubs, employeeNumber, emp.HeadOfOfficeID, emp.Grade)
		}
	}

	return pickRandom(subordinates), nil
}

// GetRandomEmployeePeers selects a random peer for the given employee with a
// multi-level fallback: office -> division -> department.
// Governor/DG employees select among other governors/DGs.
func (s *reviewAgentService) GetRandomEmployeePeers(ctx context.Context, employeeNumber string) (*erp.EmployeeDetails, error) {
	if strings.TrimSpace(employeeNumber) == "" {
		return nil, nil
	}

	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil || strings.TrimSpace(emp.Grade) == "" {
		return nil, err
	}

	var peers []erp.EmployeeDetails

	govDGs, err := s.getGovernorDGs(ctx)
	if err != nil {
		return nil, fmt.Errorf("getGovernorDGs: %w", err)
	}
	isGovernorDG := false
	for _, g := range govDGs {
		if strings.EqualFold(g.EmployeeNumber, employeeNumber) {
			isGovernorDG = true
			break
		}
	}

	if isGovernorDG {
		peers = deduplicateEmployeeDetails(govDGs, employeeNumber)
	} else {
		// 1. Try office peers
		officePeers, err := s.getPeersByOfficeAndGrades(ctx, employeeNumber, emp.Grade, emp.OfficeID)
		if err != nil {
			return nil, fmt.Errorf("getPeersByOfficeAndGrades: %w", err)
		}
		peers = filterPeers(officePeers, employeeNumber, emp.HeadOfOfficeID, emp.Grade)

		// 2. Fallback to division
		if len(peers) == 0 && emp.DivisionID != nil {
			divPeers, err := s.getPeersByDivisionAndGrades(ctx, employeeNumber, emp.Grade, *emp.DivisionID)
			if err != nil {
				return nil, fmt.Errorf("getPeersByDivisionAndGrades: %w", err)
			}
			peers = filterPeers(divPeers, employeeNumber, emp.HeadOfOfficeID, emp.Grade)
		}

		// 3. Fallback to department
		if len(peers) == 0 && emp.DepartmentID != nil {
			deptPeers, err := s.getPeersByDepartmentAndGrades(ctx, employeeNumber, emp.Grade, *emp.DepartmentID)
			if err != nil {
				return nil, fmt.Errorf("getPeersByDepartmentAndGrades: %w", err)
			}
			peers = filterPeers(deptPeers, employeeNumber, emp.HeadOfOfficeID, emp.Grade)
		}
	}

	return pickRandom(peers), nil
}

// GetRandomEmployeeSuperior selects a random superior for the given employee
// with a multi-level fallback: office -> division -> department.
func (s *reviewAgentService) GetRandomEmployeeSuperior(ctx context.Context, employeeNumber string) (*erp.EmployeeDetails, error) {
	if strings.TrimSpace(employeeNumber) == "" {
		return nil, nil
	}

	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil || strings.TrimSpace(emp.Grade) == "" {
		return nil, err
	}

	var superiors []erp.EmployeeDetails

	// 1. Try office superiors
	officeSups, err := s.getSuperiorsByOfficeAndGrades(ctx, employeeNumber, emp.Grade, emp.OfficeID)
	if err != nil {
		return nil, fmt.Errorf("getSuperiorsByOfficeAndGrades: %w", err)
	}
	if len(officeSups) > 0 {
		superiors = filterSuperiors(officeSups, employeeNumber, emp.HeadOfOfficeID, emp.SupervisorID, emp.Grade)
	}

	// 2. Fallback to division
	if len(superiors) == 0 {
		divSups, err := s.getSuperiorsByDivisionAndGrades(ctx, employeeNumber, emp.Grade, emp.DivisionID)
		if err != nil {
			return nil, fmt.Errorf("getSuperiorsByDivisionAndGrades: %w", err)
		}
		if len(divSups) > 0 {
			superiors = filterSuperiors(divSups, employeeNumber, emp.HeadOfOfficeID, emp.SupervisorID, emp.Grade)
		}
	}

	// 3. Fallback to department
	if len(superiors) == 0 {
		deptSups, err := s.getSuperiorsByDepartmentAndGrades(ctx, employeeNumber, emp.Grade, emp.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("getSuperiorsByDepartmentAndGrades: %w", err)
		}
		if len(deptSups) > 0 {
			superiors = filterSuperiors(deptSups, employeeNumber, emp.HeadOfOfficeID, emp.SupervisorID, emp.Grade)
		}
	}

	return pickRandom(superiors), nil
}

// ---------------------------------------------------------------------------
// 2. Review Population — creates review records for employees
// ---------------------------------------------------------------------------

// PopulateReviewsForEmployee creates all review type records for a single employee.
func (s *reviewAgentService) PopulateReviewsForEmployee(ctx context.Context, employeeNumber string) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil {
		return err
	}
	reviewTypes, err := s.getReviewTypes(ctx)
	if err != nil || len(reviewTypes) == 0 {
		return err
	}
	employees := []erp.EmployeeErpDetailsDTO{*emp}
	return s.createEmployeesReviews(ctx, period, employees, reviewTypes)
}

// PopulateOfficeEmployeeReviews creates reviews for all employees in an office.
func (s *reviewAgentService) PopulateOfficeEmployeeReviews(ctx context.Context, officeID int) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getEmployeesByOfficeID(ctx, officeID)
	if err != nil || len(employees) == 0 {
		return err
	}
	reviewTypes, err := s.getReviewTypes(ctx)
	if err != nil || len(reviewTypes) == 0 {
		return err
	}
	return s.createEmployeesReviews(ctx, period, employees, reviewTypes)
}

// PopulateDivisionEmployeeReviews creates reviews for all employees in a division.
func (s *reviewAgentService) PopulateDivisionEmployeeReviews(ctx context.Context, divisionID int) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getEmployeesByDivisionID(ctx, divisionID)
	if err != nil || len(employees) == 0 {
		return err
	}
	reviewTypes, err := s.getReviewTypes(ctx)
	if err != nil || len(reviewTypes) == 0 {
		return err
	}
	return s.createEmployeesReviews(ctx, period, employees, reviewTypes)
}

// PopulateDepartmentEmployeeReviews creates reviews for all employees in a department.
func (s *reviewAgentService) PopulateDepartmentEmployeeReviews(ctx context.Context, departmentID int) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getEmployeesByDepartmentID(ctx, departmentID)
	if err != nil || len(employees) == 0 {
		return err
	}
	reviewTypes, err := s.getReviewTypes(ctx)
	if err != nil || len(reviewTypes) == 0 {
		return err
	}
	return s.createEmployeesReviews(ctx, period, employees, reviewTypes)
}

// PopulateAllEmployeeReviews creates reviews for every active employee.
func (s *reviewAgentService) PopulateAllEmployeeReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}
	reviewTypes, err := s.getReviewTypes(ctx)
	if err != nil || len(reviewTypes) == 0 {
		return err
	}
	return s.createEmployeesReviews(ctx, period, employees, reviewTypes)
}

// createEmployeesReviews iterates through employees and creates reviews for each.
// In the .NET version this enqueues Hangfire background jobs; in Go we process
// sequentially within the same goroutine context.
func (s *reviewAgentService) createEmployeesReviews(
	ctx context.Context,
	period *competency.ReviewPeriod,
	employees []erp.EmployeeErpDetailsDTO,
	reviewTypes []competency.ReviewType,
) error {
	for i := range employees {
		if err := s.ProcessEmployeeReviewCreation(ctx, period, reviewTypes, &employees[i]); err != nil {
			s.log.Error().Err(err).
				Str("employee", employees[i].EmployeeNumber).
				Msg("failed to process employee review creation")
			// Continue processing other employees, do not abort batch.
		}
	}
	return nil
}

// ProcessEmployeeReviewCreation creates all review-type records for a single
// employee. It resolves grade group, office, and job role, then dispatches to
// the per-review-type population methods.
func (s *reviewAgentService) ProcessEmployeeReviewCreation(
	ctx context.Context,
	period *competency.ReviewPeriod,
	reviewTypes []competency.ReviewType,
	employee *erp.EmployeeErpDetailsDTO,
) error {
	if employee == nil || period == nil || len(reviewTypes) == 0 {
		return nil
	}

	gradeGroup, err := s.getAssignedJobGradeGroup(ctx, employee.Grade)
	if err != nil {
		return fmt.Errorf("getAssignedJobGradeGroup: %w", err)
	}
	officeCheck, err := s.getOfficeByID(ctx, employee.OfficeID)
	if err != nil {
		return fmt.Errorf("getOfficeByID: %w", err)
	}

	// Clean position
	if strings.TrimSpace(employee.Position) != "" {
		parts := strings.SplitN(employee.Position, ".", 2)
		employee.Position = parts[0]
	}

	jobRole, err := s.getJobRoleByName(ctx, employee.Position)
	if err != nil {
		return fmt.Errorf("getJobRoleByName: %w", err)
	}

	if officeCheck == nil || gradeGroup == nil {
		return nil
	}
	officeID := officeCheck.OfficeID

	if jobRole == nil {
		jobRole = &competency.JobRole{
			JobRoleID:   0,
			Description: employee.Position,
		}
	}

	gradeGroupName := ""
	if gradeGroup.JobGradeGroup != nil {
		gradeGroupName = gradeGroup.JobGradeGroup.GroupName
	}

	for _, rt := range reviewTypes {
		if strings.TrimSpace(rt.ReviewTypeName) == "" {
			continue
		}

		var processErr error
		switch rt.ReviewTypeName {
		case "Self":
			processErr = s.PopulateReviewForSelfReviewType(ctx, employee, gradeGroupName, period.ReviewPeriodID, rt.ReviewTypeID, officeID, jobRole.JobRoleID)
		case "Supervisor":
			processErr = s.PopulateReviewForSupervisorReviewType(ctx, employee, gradeGroupName, period.ReviewPeriodID, rt.ReviewTypeID, officeID, jobRole.JobRoleID)
		case "Peers":
			processErr = s.PopulateReviewForPeerReviewType(ctx, employee, gradeGroupName, period.ReviewPeriodID, rt.ReviewTypeID)
		case "Subordinates":
			processErr = s.PopulateReviewForSubordinateReviewType(ctx, employee, gradeGroupName, period.ReviewPeriodID, rt.ReviewTypeID)
		case "Superior":
			processErr = s.PopulateReviewForSuperiorReviewType(ctx, employee, gradeGroupName, period.ReviewPeriodID, rt.ReviewTypeID)
		}

		if processErr != nil {
			s.log.Error().Err(processErr).
				Str("reviewType", rt.ReviewTypeName).
				Str("employee", employee.EmployeeNumber).
				Msg("failed to process review type")
			// Continue with other review types
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// 3. Per-Review-Type Population Methods
// ---------------------------------------------------------------------------

// PopulateReviewForSelfReviewType creates self-review records (both behavioral
// and technical) for an employee.
func (s *reviewAgentService) PopulateReviewForSelfReviewType(
	ctx context.Context,
	employee *erp.EmployeeErpDetailsDTO,
	gradeGroupName string,
	reviewPeriodID, reviewTypeID, officeID, jobRoleID int,
) error {
	if employee == nil {
		return nil
	}

	// --- Behavioral ---
	alreadyBehavioral, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND is_technical = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, false, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check behavioral exists: %w", err)
	}

	if !alreadyBehavioral {
		behaviorals, err := s.getBehavioralCompetencies(ctx, gradeGroupName)
		if err != nil {
			return fmt.Errorf("getBehavioralCompetencies: %w", err)
		}
		if len(behaviorals) > 0 {
			reviews := make([]competency.CompetencyReview, 0, len(behaviorals))
			for _, bc := range behaviorals {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       bc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   bc.RatingID,
					ReviewerID:         employee.EmployeeNumber,
					ReviewerName:       employee.FullName(),
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        false,
				})
			}
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create behavioral reviews: %w", err)
			}
		}
	}

	// --- Technical ---
	alreadyTechnical, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND is_technical = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, true, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check technical exists: %w", err)
	}

	if !alreadyTechnical {
		techComps, err := s.getJobRoleCompetencies(ctx, officeID, jobRoleID)
		if err != nil {
			return fmt.Errorf("getJobRoleCompetencies: %w", err)
		}

		var reviews []competency.CompetencyReview
		if len(techComps) > 0 {
			for _, tc := range techComps {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       tc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   tc.RatingID,
					ReviewerID:         employee.EmployeeNumber,
					ReviewerName:       employee.FullName(),
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        true,
				})
			}
		} else {
			// Fallback: partial match against all office competencies
			allOfficeComps, err := s.getJobRoleCompetenciesByOffice(ctx, officeID)
			if err != nil {
				return fmt.Errorf("getJobRoleCompetenciesByOffice: %w", err)
			}
			empJobDescription := employee.Position + "." + employee.OfficeName + "." + employee.JobName
			seen := make(map[int]bool)
			for _, tc := range allOfficeComps {
				desc := ""
				if tc.JobRole != nil {
					desc = tc.JobRole.Description
				}
				if !isPartialMatch(desc, empJobDescription) {
					continue
				}
				if seen[tc.CompetencyID] {
					continue
				}
				seen[tc.CompetencyID] = true

				// Check if this specific technical review already exists
				exists, err := s.reviewRepo.Exists(ctx,
					"employee_number = ? AND is_technical = ? AND competency_id = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
					employee.EmployeeNumber, true, tc.CompetencyID, reviewPeriodID, reviewTypeID, false)
				if err != nil {
					return fmt.Errorf("check technical review exists: %w", err)
				}
				if !exists {
					reviews = append(reviews, competency.CompetencyReview{
						CompetencyID:       tc.CompetencyID,
						EmployeeNumber:     employee.EmployeeNumber,
						ExpectedRatingID:   tc.RatingID,
						ReviewerID:         employee.EmployeeNumber,
						ReviewerName:       employee.FullName(),
						ReviewPeriodID:     reviewPeriodID,
						ReviewTypeID:       reviewTypeID,
						EmployeeName:       employee.FullName(),
						EmployeeInitial:    employee.NameInitial(),
						EmployeeGrade:      employee.Grade,
						EmployeeDepartment: employee.DepartmentName,
						IsTechnical:        true,
					})
				}
			}
		}

		if len(reviews) > 0 {
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create technical reviews: %w", err)
			}
		}
	}

	return nil
}

// getEmployeeSupervisor resolves the supervisor for an employee.
// Currently returns the employee's direct supervisor (SupervisorID).
func (s *reviewAgentService) getEmployeeSupervisor(ctx context.Context, employee *erp.EmployeeErpDetailsDTO) (*erp.EmployeeErpDetailsDTO, error) {
	if employee == nil {
		return nil, nil
	}
	return s.getEmployeeDetail(ctx, employee.SupervisorID)
}

// PopulateReviewForSupervisorReviewType creates supervisor-review records (both
// behavioral and technical) for an employee.
func (s *reviewAgentService) PopulateReviewForSupervisorReviewType(
	ctx context.Context,
	employee *erp.EmployeeErpDetailsDTO,
	gradeGroupName string,
	reviewPeriodID, reviewTypeID, officeID, jobRoleID int,
) error {
	if employee == nil {
		return nil
	}

	supervisor, err := s.getEmployeeSupervisor(ctx, employee)
	if err != nil {
		return fmt.Errorf("getEmployeeSupervisor: %w", err)
	}

	reviewerID := ""
	reviewerName := ""
	if supervisor != nil {
		reviewerID = supervisor.EmployeeNumber
		reviewerName = supervisor.FullName()
	}

	// --- Behavioral ---
	alreadyBehavioral, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND is_technical = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, false, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check behavioral exists: %w", err)
	}

	if !alreadyBehavioral {
		behaviorals, err := s.getBehavioralCompetencies(ctx, gradeGroupName)
		if err != nil {
			return fmt.Errorf("getBehavioralCompetencies: %w", err)
		}
		if reviewerID != "" && len(behaviorals) > 0 {
			reviews := make([]competency.CompetencyReview, 0, len(behaviorals))
			for _, bc := range behaviorals {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       bc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   bc.RatingID,
					ReviewerID:         reviewerID,
					ReviewerName:       reviewerName,
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        false,
				})
			}
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create supervisor behavioral reviews: %w", err)
			}
		}
	}

	// --- Technical ---
	alreadyTechnical, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND is_technical = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, true, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check technical exists: %w", err)
	}

	if !alreadyTechnical {
		techComps, err := s.getJobRoleCompetencies(ctx, officeID, jobRoleID)
		if err != nil {
			return fmt.Errorf("getJobRoleCompetencies: %w", err)
		}

		var reviews []competency.CompetencyReview
		if len(techComps) > 0 {
			for _, tc := range techComps {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       tc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   tc.RatingID,
					ReviewerID:         reviewerID,
					ReviewerName:       reviewerName,
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        true,
				})
			}
		} else {
			// Fallback: partial match
			allOfficeComps, err := s.getJobRoleCompetenciesByOffice(ctx, officeID)
			if err != nil {
				return fmt.Errorf("getJobRoleCompetenciesByOffice: %w", err)
			}
			empJobDescription := employee.Position + "." + employee.OfficeName + "." + employee.JobName
			seen := make(map[int]bool)
			for _, tc := range allOfficeComps {
				desc := ""
				if tc.JobRole != nil {
					desc = tc.JobRole.Description
				}
				if !isPartialMatch(desc, empJobDescription) {
					continue
				}
				if seen[tc.CompetencyID] {
					continue
				}
				seen[tc.CompetencyID] = true

				// Check individual competency existence
				exists, err := s.reviewRepo.Exists(ctx,
					"employee_number = ? AND is_technical = ? AND competency_id = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
					employee.EmployeeNumber, true, tc.CompetencyID, reviewPeriodID, reviewTypeID, false)
				if err != nil {
					return fmt.Errorf("check technical review exists: %w", err)
				}
				if !exists {
					// Additional duplicate check matching .NET's second AnyAsync
					var dupCount int64
					s.db.WithContext(ctx).Model(&competency.CompetencyReview{}).
						Where("employee_number = ? AND review_period_id = ? AND competency_id = ? AND is_technical = ? AND review_type_id = ? AND reviewer_id = ?",
							employee.EmployeeNumber, reviewPeriodID, tc.CompetencyID, true, reviewTypeID, reviewerID).
						Count(&dupCount)
					if dupCount == 0 {
						reviews = append(reviews, competency.CompetencyReview{
							CompetencyID:       tc.CompetencyID,
							EmployeeNumber:     employee.EmployeeNumber,
							ExpectedRatingID:   tc.RatingID,
							ReviewerID:         reviewerID,
							ReviewerName:       reviewerName,
							ReviewPeriodID:     reviewPeriodID,
							ReviewTypeID:       reviewTypeID,
							EmployeeName:       employee.FullName(),
							EmployeeInitial:    employee.NameInitial(),
							EmployeeGrade:      employee.Grade,
							EmployeeDepartment: employee.DepartmentName,
							IsTechnical:        true,
						})
					}
				}
			}
		}

		if len(reviews) > 0 {
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create supervisor technical reviews: %w", err)
			}
		}
	}

	return nil
}

// PopulateReviewForPeerReviewType creates peer-review records (behavioral only)
// for an employee. A random peer is selected.
func (s *reviewAgentService) PopulateReviewForPeerReviewType(
	ctx context.Context,
	employee *erp.EmployeeErpDetailsDTO,
	gradeGroupName string,
	reviewPeriodID, reviewTypeID int,
) error {
	if employee == nil {
		return nil
	}

	alreadyCreated, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check peer review exists: %w", err)
	}

	if !alreadyCreated {
		behaviorals, err := s.getBehavioralCompetencies(ctx, gradeGroupName)
		if err != nil {
			return fmt.Errorf("getBehavioralCompetencies: %w", err)
		}

		peer, err := s.GetRandomEmployeePeers(ctx, employee.EmployeeNumber)
		if err != nil {
			return fmt.Errorf("GetRandomEmployeePeers: %w", err)
		}

		if peer != nil && len(behaviorals) > 0 {
			reviews := make([]competency.CompetencyReview, 0, len(behaviorals))
			for _, bc := range behaviorals {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       bc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   bc.RatingID,
					ReviewerID:         peer.EmployeeNumber,
					ReviewerName:       peer.FullName,
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        false,
				})
			}
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create peer reviews: %w", err)
			}
		}
	}

	return nil
}

// PopulateReviewForSubordinateReviewType creates subordinate-review records
// (behavioral only) for an employee. A random subordinate is selected.
func (s *reviewAgentService) PopulateReviewForSubordinateReviewType(
	ctx context.Context,
	employee *erp.EmployeeErpDetailsDTO,
	gradeGroupName string,
	reviewPeriodID, reviewTypeID int,
) error {
	if employee == nil {
		return nil
	}

	alreadyCreated, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check subordinate review exists: %w", err)
	}

	if !alreadyCreated {
		behaviorals, err := s.getBehavioralCompetencies(ctx, gradeGroupName)
		if err != nil {
			return fmt.Errorf("getBehavioralCompetencies: %w", err)
		}

		subordinate, err := s.GetRandomEmployeeSubordinate(ctx, employee.EmployeeNumber)
		if err != nil {
			return fmt.Errorf("GetRandomEmployeeSubordinate: %w", err)
		}

		if subordinate != nil && !strings.EqualFold(subordinate.EmployeeNumber, employee.EmployeeNumber) && len(behaviorals) > 0 {
			reviews := make([]competency.CompetencyReview, 0, len(behaviorals))
			for _, bc := range behaviorals {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       bc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   bc.RatingID,
					ReviewerID:         subordinate.EmployeeNumber,
					ReviewerName:       subordinate.FullName,
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        false,
				})
			}
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create subordinate reviews: %w", err)
			}
		}
	}

	return nil
}

// PopulateReviewForSuperiorReviewType creates superior-review records
// (behavioral only) for an employee. A random superior is selected.
func (s *reviewAgentService) PopulateReviewForSuperiorReviewType(
	ctx context.Context,
	employee *erp.EmployeeErpDetailsDTO,
	gradeGroupName string,
	reviewPeriodID, reviewTypeID int,
) error {
	if employee == nil {
		return nil
	}

	alreadyCreated, err := s.reviewRepo.Exists(ctx,
		"employee_number = ? AND review_period_id = ? AND review_type_id = ? AND soft_deleted = ?",
		employee.EmployeeNumber, reviewPeriodID, reviewTypeID, false)
	if err != nil {
		return fmt.Errorf("check superior review exists: %w", err)
	}

	if !alreadyCreated {
		behaviorals, err := s.getBehavioralCompetencies(ctx, gradeGroupName)
		if err != nil {
			return fmt.Errorf("getBehavioralCompetencies: %w", err)
		}

		superior, err := s.GetRandomEmployeeSuperior(ctx, employee.EmployeeNumber)
		if err != nil {
			return fmt.Errorf("GetRandomEmployeeSuperior: %w", err)
		}

		if superior != nil && len(behaviorals) > 0 {
			reviews := make([]competency.CompetencyReview, 0, len(behaviorals))
			for _, bc := range behaviorals {
				reviews = append(reviews, competency.CompetencyReview{
					CompetencyID:       bc.CompetencyID,
					EmployeeNumber:     employee.EmployeeNumber,
					ExpectedRatingID:   bc.RatingID,
					ReviewerID:         superior.EmployeeNumber,
					ReviewerName:       superior.FullName,
					ReviewPeriodID:     reviewPeriodID,
					ReviewTypeID:       reviewTypeID,
					EmployeeName:       employee.FullName(),
					EmployeeInitial:    employee.NameInitial(),
					EmployeeGrade:      employee.Grade,
					EmployeeDepartment: employee.DepartmentName,
					IsTechnical:        false,
				})
			}
			if err := s.reviewRepo.CreateBatch(ctx, reviews); err != nil {
				return fmt.Errorf("create superior reviews: %w", err)
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// 4. Review Calculation — computes weighted average ratings
// ---------------------------------------------------------------------------

// roundRating replicates the .NET IntExtensions.Round(double) method:
// if fractional part >= 0.5, round up; otherwise floor.
func roundRating(value float64) int {
	decimal := math.Abs(value - math.Floor(value))
	if decimal >= 0.5 {
		return int(math.Round(value))
	}
	return int(math.Floor(value))
}

// computeCompetencyGap returns the non-negative gap between expected and actual.
func computeCompetencyGap(expected, actual int) int {
	gap := expected - actual
	if gap < 0 {
		return 0
	}
	return gap
}

// computeHaveGap returns true if expected > actual, false otherwise.
func computeHaveGap(expected, actual int) bool {
	return expected > actual
}

// findRatingByValue returns the first rating whose Value equals the given value.
func findRatingByValue(ratings []competency.Rating, value int) *competency.Rating {
	for i := range ratings {
		if ratings[i].Value == value {
			return &ratings[i]
		}
	}
	return nil
}

// CalculateBehavioralReviewAverage computes the average behavioral competency
// review score for an employee in a given review period. It creates or updates
// CompetencyReviewProfile records per competency.
func (s *reviewAgentService) CalculateBehavioralReviewAverage(ctx context.Context, employeeNumber string, reviewPeriodID int) error {
	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil {
		return err
	}

	// Clean position for job role lookup
	if strings.TrimSpace(emp.Position) != "" {
		parts := strings.SplitN(emp.Position, ".", 2)
		emp.Position = parts[0]
	}
	jobRole, err := s.getJobRoleByName(ctx, emp.Position)
	if err != nil {
		return fmt.Errorf("getJobRoleByName: %w", err)
	}

	ratings, err := s.getRatings(ctx)
	if err != nil || len(ratings) == 0 {
		return err
	}

	// Get all behavioral reviews with actual ratings for this employee/period
	var allReviews []competency.CompetencyReview
	err = s.db.WithContext(ctx).
		Where("employee_number = ? AND actual_rating_value != 0 AND is_technical = ? AND review_period_id = ? AND soft_deleted = ?",
			employeeNumber, false, reviewPeriodID, false).
		Preload("Competency.CompetencyCategory").
		Preload("ReviewPeriod").
		Preload("ExpectedRating").
		Find(&allReviews).Error
	if err != nil {
		return fmt.Errorf("fetch behavioral reviews: %w", err)
	}

	// Identify distinct competencies
	competencyMap := make(map[int]*competency.CompetencyReview)
	for i := range allReviews {
		if _, exists := competencyMap[allReviews[i].CompetencyID]; !exists {
			competencyMap[allReviews[i].CompetencyID] = &allReviews[i]
		}
	}

	for compID, compReview := range competencyMap {
		// Calculate average rating for this competency
		var sum float64
		var count int
		for i := range allReviews {
			if allReviews[i].CompetencyID == compID {
				sum += float64(allReviews[i].ActualRatingValue)
				count++
			}
		}
		var avgScore float64
		if count > 0 {
			avgScore = sum / float64(count)
		}
		roundedScore := roundRating(avgScore)

		// Look up existing profile
		profile, err := s.profileRepo.FirstOrDefault(ctx,
			"employee_number = ? AND review_period_id = ? AND competency_id = ? AND soft_deleted = ?",
			employeeNumber, reviewPeriodID, compID, false)
		if err != nil {
			return fmt.Errorf("fetch profile: %w", err)
		}

		matchedRating := findRatingByValue(ratings, roundedScore)

		if profile != nil {
			// Update existing
			profile.EmployeeName = emp.FullName()
			profile.OfficeID = fmt.Sprintf("%d", emp.OfficeID)
			profile.OfficeName = emp.OfficeName
			profile.DepartmentID = intPtrToString(emp.DepartmentID)
			profile.DepartmentName = emp.DepartmentName
			profile.DivisionID = intPtrToString(emp.DivisionID)
			profile.DivisionName = emp.DivisionName
			if jobRole != nil {
				profile.JobRoleID = fmt.Sprintf("%d", jobRole.JobRoleID)
				profile.JobRoleName = jobRole.JobRoleName
			}
			profile.GradeName = emp.Grade

			if matchedRating != nil {
				profile.AverageRatingID = matchedRating.RatingID
				profile.AverageRatingName = matchedRating.Name
				profile.AverageRatingValue = matchedRating.Value
			} else {
				profile.AverageRatingID = 1
			}
			profile.AverageScore = float64(roundedScore)

			profile.CompetencyGap = computeCompetencyGap(profile.ExpectedRatingValue, profile.AverageRatingValue)
			profile.HaveGap = computeHaveGap(profile.ExpectedRatingValue, profile.AverageRatingValue)

			if err := s.profileRepo.Update(ctx, profile); err != nil {
				return fmt.Errorf("update profile: %w", err)
			}
		} else {
			// Create new profile
			newProfile := competency.CompetencyReviewProfile{
				EmployeeNumber:     employeeNumber,
				EmployeeName:       emp.FullName(),
				OfficeID:           fmt.Sprintf("%d", emp.OfficeID),
				OfficeName:         emp.OfficeName,
				DepartmentID:       intPtrToString(emp.DepartmentID),
				DepartmentName:     emp.DepartmentName,
				DivisionID:         intPtrToString(emp.DivisionID),
				DivisionName:       emp.DivisionName,
				GradeName:          emp.Grade,
				ReviewPeriodID:     reviewPeriodID,
				CompetencyID:       compID,
				ExpectedRatingID:   compReview.ExpectedRatingID,
			}
			if jobRole != nil {
				newProfile.JobRoleID = fmt.Sprintf("%d", jobRole.JobRoleID)
				newProfile.JobRoleName = jobRole.JobRoleName
			}
			if compReview.ReviewPeriod != nil {
				newProfile.ReviewPeriodName = compReview.ReviewPeriod.Name
			}
			if compReview.Competency != nil {
				newProfile.CompetencyName = compReview.Competency.CompetencyName
				if compReview.Competency.CompetencyCategory != nil {
					newProfile.CompetencyCategoryName = compReview.Competency.CompetencyCategory.CategoryName
				}
			}
			if compReview.ExpectedRating != nil {
				newProfile.ExpectedRatingValue = compReview.ExpectedRating.Value
				newProfile.ExpectedRatingName = compReview.ExpectedRating.Name
			}

			if matchedRating != nil {
				newProfile.AverageRatingID = matchedRating.RatingID
				newProfile.AverageRatingName = matchedRating.Name
				newProfile.AverageRatingValue = matchedRating.Value
			} else {
				newProfile.AverageRatingID = 1
			}
			newProfile.AverageScore = float64(roundedScore)

			newProfile.CompetencyGap = computeCompetencyGap(newProfile.ExpectedRatingValue, newProfile.AverageRatingValue)
			newProfile.HaveGap = computeHaveGap(newProfile.ExpectedRatingValue, newProfile.AverageRatingValue)

			if err := s.profileRepo.Create(ctx, &newProfile); err != nil {
				return fmt.Errorf("create profile: %w", err)
			}
		}
	}

	return nil
}

// CalculateTechnicalReviewAverage computes the weighted average technical
// competency review score for an employee. The weighting uses self-review and
// supervisor-review percentages from CompetencyCategoryGrading.
func (s *reviewAgentService) CalculateTechnicalReviewAverage(ctx context.Context, employeeNumber string, reviewPeriodID int) error {
	ratings, err := s.getRatings(ctx)
	if err != nil || len(ratings) == 0 {
		return err
	}

	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil {
		return err
	}

	// Clean position
	if strings.TrimSpace(emp.Position) != "" {
		parts := strings.SplitN(emp.Position, ".", 2)
		emp.Position = parts[0]
	}
	jobRole, err := s.getJobRoleByName(ctx, emp.Position)
	if err != nil {
		return fmt.Errorf("getJobRoleByName: %w", err)
	}

	// Get all technical reviews for this employee/period (including those with 0 rating)
	var allReviews []competency.CompetencyReview
	err = s.db.WithContext(ctx).
		Where("employee_number = ? AND is_technical = ? AND review_period_id = ? AND soft_deleted = ?",
			employeeNumber, true, reviewPeriodID, false).
		Preload("Competency.CompetencyCategory").
		Preload("ReviewPeriod").
		Preload("ExpectedRating").
		Find(&allReviews).Error
	if err != nil {
		return fmt.Errorf("fetch technical reviews: %w", err)
	}

	// Identify distinct competencies
	competencyMap := make(map[int]*competency.CompetencyReview)
	for i := range allReviews {
		if _, exists := competencyMap[allReviews[i].CompetencyID]; !exists {
			competencyMap[allReviews[i].CompetencyID] = &allReviews[i]
		}
	}

	if jobRole == nil || emp == nil || len(ratings) == 0 {
		return nil
	}

	for compID, compReview := range competencyMap {
		if compReview == nil || compReview.Competency == nil {
			continue
		}

		// Get self and supervisor weight percentages from grading table
		selfWeight := 30.0  // default
		supWeight := 70.0   // default

		var selfGrading competency.CompetencyCategoryGrading
		err := s.db.WithContext(ctx).
			Joins("JOIN \"CoreSchema\".\"review_types\" rt ON rt.review_type_id = \"CoreSchema\".\"competency_category_gradings\".review_type_id").
			Where("rt.review_type_name = ? AND \"CoreSchema\".\"competency_category_gradings\".competency_category_id = ? AND \"CoreSchema\".\"competency_category_gradings\".soft_deleted = ?",
				"Self", compReview.Competency.CompetencyCategoryID, false).
			First(&selfGrading).Error
		if err == nil {
			selfWeight = selfGrading.WeightPercentage
		}

		var supGrading competency.CompetencyCategoryGrading
		err = s.db.WithContext(ctx).
			Joins("JOIN \"CoreSchema\".\"review_types\" rt ON rt.review_type_id = \"CoreSchema\".\"competency_category_gradings\".review_type_id").
			Where("rt.review_type_name = ? AND \"CoreSchema\".\"competency_category_gradings\".competency_category_id = ? AND \"CoreSchema\".\"competency_category_gradings\".soft_deleted = ?",
				"Supervisor", compReview.Competency.CompetencyCategoryID, false).
			First(&supGrading).Error
		if err == nil {
			supWeight = supGrading.WeightPercentage
		}

		// Calculate self-review average for this competency
		var selfSum float64
		var selfCount int
		for i := range allReviews {
			if allReviews[i].CompetencyID == compID && strings.EqualFold(allReviews[i].ReviewerID, employeeNumber) {
				selfSum += float64(allReviews[i].ActualRatingValue)
				selfCount++
			}
		}
		var allSelfScore float64
		if selfCount > 0 {
			allSelfScore = selfSum / float64(selfCount)
		}

		// Calculate supervisor/non-self review average
		var supSum float64
		var supCount int
		for i := range allReviews {
			if allReviews[i].CompetencyID == compID && !strings.EqualFold(allReviews[i].ReviewerID, employeeNumber) {
				supSum += float64(allReviews[i].ActualRatingValue)
				supCount++
			}
		}
		var allSupervisorScore float64
		if supCount > 0 {
			allSupervisorScore = supSum / float64(supCount)
		}

		// Weighted average
		selfScore := (allSelfScore * selfWeight) / 100.0
		supervisorScore := (allSupervisorScore * supWeight) / 100.0
		averageRatingSum := selfScore + supervisorScore
		roundedScore := roundRating(averageRatingSum)

		// Look up existing profile
		profile, err := s.profileRepo.FirstOrDefault(ctx,
			"employee_number = ? AND review_period_id = ? AND competency_id = ? AND soft_deleted = ?",
			employeeNumber, reviewPeriodID, compID, false)
		if err != nil {
			return fmt.Errorf("fetch profile: %w", err)
		}

		matchedRating := findRatingByValue(ratings, roundedScore)

		if profile != nil {
			// Update existing
			profile.EmployeeName = emp.FullName()
			profile.OfficeID = fmt.Sprintf("%d", emp.OfficeID)
			profile.OfficeName = emp.OfficeName
			profile.DepartmentID = intPtrToString(emp.DepartmentID)
			profile.DepartmentName = emp.DepartmentName
			profile.DivisionID = intPtrToString(emp.DivisionID)
			profile.DivisionName = emp.DivisionName
			if jobRole != nil {
				profile.JobRoleID = fmt.Sprintf("%d", jobRole.JobRoleID)
				profile.JobRoleName = jobRole.JobRoleName
			}
			profile.GradeName = emp.Grade

			if matchedRating != nil {
				profile.AverageRatingID = matchedRating.RatingID
				profile.AverageRatingName = matchedRating.Name
				profile.AverageRatingValue = matchedRating.Value
			} else {
				profile.AverageRatingID = 1
			}
			profile.AverageScore = float64(roundedScore)

			profile.CompetencyGap = computeCompetencyGap(profile.ExpectedRatingValue, profile.AverageRatingValue)
			profile.HaveGap = computeHaveGap(profile.ExpectedRatingValue, profile.AverageRatingValue)

			if err := s.profileRepo.Update(ctx, profile); err != nil {
				return fmt.Errorf("update profile: %w", err)
			}
		} else {
			// Create new
			newProfile := competency.CompetencyReviewProfile{
				EmployeeNumber:     employeeNumber,
				EmployeeName:       emp.FullName(),
				OfficeID:           fmt.Sprintf("%d", emp.OfficeID),
				OfficeName:         emp.OfficeName,
				DepartmentID:       intPtrToString(emp.DepartmentID),
				DepartmentName:     emp.DepartmentName,
				DivisionID:         intPtrToString(emp.DivisionID),
				DivisionName:       emp.DivisionName,
				GradeName:          emp.Grade,
				ReviewPeriodID:     reviewPeriodID,
				CompetencyID:       compID,
				ExpectedRatingID:   compReview.ExpectedRatingID,
			}
			if jobRole != nil {
				newProfile.JobRoleID = fmt.Sprintf("%d", jobRole.JobRoleID)
				newProfile.JobRoleName = jobRole.JobRoleName
			}
			if compReview.ReviewPeriod != nil {
				newProfile.ReviewPeriodName = compReview.ReviewPeriod.Name
			}
			if compReview.Competency != nil {
				newProfile.CompetencyName = compReview.Competency.CompetencyName
				if compReview.Competency.CompetencyCategory != nil {
					newProfile.CompetencyCategoryName = compReview.Competency.CompetencyCategory.CategoryName
				}
			}
			if compReview.ExpectedRating != nil {
				newProfile.ExpectedRatingValue = compReview.ExpectedRating.Value
				newProfile.ExpectedRatingName = compReview.ExpectedRating.Name
			}

			if matchedRating != nil {
				newProfile.AverageRatingID = matchedRating.RatingID
				newProfile.AverageRatingName = matchedRating.Name
				newProfile.AverageRatingValue = matchedRating.Value
			} else {
				newProfile.AverageRatingID = 1
			}
			// .NET stores averageRatingSum (not rounded) for technical new profiles
			newProfile.AverageScore = averageRatingSum

			newProfile.CompetencyGap = computeCompetencyGap(newProfile.ExpectedRatingValue, newProfile.AverageRatingValue)
			newProfile.HaveGap = computeHaveGap(newProfile.ExpectedRatingValue, newProfile.AverageRatingValue)

			if err := s.profileRepo.Create(ctx, &newProfile); err != nil {
				return fmt.Errorf("create profile: %w", err)
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// 5. Batch population / calculation methods
// ---------------------------------------------------------------------------

// PopulateAllEmployeeSelfReviews creates self-review records for all employees.
func (s *reviewAgentService) PopulateAllEmployeeSelfReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for i := range employees {
		emp := &employees[i]
		gradeGroup, err := s.getAssignedJobGradeGroup(ctx, emp.Grade)
		if err != nil || gradeGroup == nil {
			continue
		}
		if strings.TrimSpace(emp.Position) != "" {
			parts := strings.SplitN(emp.Position, ".", 2)
			emp.Position = parts[0]
		}
		jobRole, err := s.getJobRoleByName(ctx, emp.Position)
		if err != nil {
			continue
		}
		office, err := s.getOfficeByID(ctx, emp.OfficeID)
		if err != nil || office == nil {
			continue
		}
		if jobRole == nil {
			continue
		}

		gradeGroupName := ""
		if gradeGroup.JobGradeGroup != nil {
			gradeGroupName = gradeGroup.JobGradeGroup.GroupName
		}

		// Self review type ID = 3 (hardcoded in .NET)
		selfReviewTypeID := 3
		if err := s.PopulateReviewForSelfReviewType(ctx, emp, gradeGroupName, period.ReviewPeriodID, selfReviewTypeID, office.OfficeID, jobRole.JobRoleID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to populate self review")
		}
	}
	return nil
}

// PopulateAllEmployeeSupervisorReviews creates supervisor-review records for all employees.
func (s *reviewAgentService) PopulateAllEmployeeSupervisorReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for i := range employees {
		emp := &employees[i]
		gradeGroup, err := s.getAssignedJobGradeGroup(ctx, emp.Grade)
		if err != nil || gradeGroup == nil {
			continue
		}
		if strings.TrimSpace(emp.Position) != "" {
			parts := strings.SplitN(emp.Position, ".", 2)
			emp.Position = parts[0]
		}
		jobRole, err := s.getJobRoleByName(ctx, emp.Position)
		if err != nil {
			continue
		}
		office, err := s.getOfficeByID(ctx, emp.OfficeID)
		if err != nil || office == nil {
			continue
		}
		if jobRole == nil {
			continue
		}

		gradeGroupName := ""
		if gradeGroup.JobGradeGroup != nil {
			gradeGroupName = gradeGroup.JobGradeGroup.GroupName
		}

		// Supervisor review type ID = 1 (hardcoded in .NET)
		supervisorReviewTypeID := 1
		if err := s.PopulateReviewForSupervisorReviewType(ctx, emp, gradeGroupName, period.ReviewPeriodID, supervisorReviewTypeID, office.OfficeID, jobRole.JobRoleID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to populate supervisor review")
		}
	}
	return nil
}

// PopulateAllEmployeePeersReviews creates peer-review records for all employees.
func (s *reviewAgentService) PopulateAllEmployeePeersReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for i := range employees {
		emp := &employees[i]
		gradeGroup, err := s.getAssignedJobGradeGroup(ctx, emp.Grade)
		if err != nil || gradeGroup == nil {
			continue
		}
		if strings.TrimSpace(emp.Position) != "" {
			parts := strings.SplitN(emp.Position, ".", 2)
			emp.Position = parts[0]
		}
		jobRole, err := s.getJobRoleByName(ctx, emp.Position)
		if err != nil || jobRole == nil {
			continue
		}

		gradeGroupName := ""
		if gradeGroup.JobGradeGroup != nil {
			gradeGroupName = gradeGroup.JobGradeGroup.GroupName
		}

		// Peers review type ID = 2 (hardcoded in .NET)
		peersReviewTypeID := 2
		if err := s.PopulateReviewForPeerReviewType(ctx, emp, gradeGroupName, period.ReviewPeriodID, peersReviewTypeID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to populate peers review")
		}
	}
	return nil
}

// PopulateAllEmployeeSubordinatesReviews creates subordinate-review records for all employees.
func (s *reviewAgentService) PopulateAllEmployeeSubordinatesReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for i := range employees {
		emp := &employees[i]
		gradeGroup, err := s.getAssignedJobGradeGroup(ctx, emp.Grade)
		if err != nil || gradeGroup == nil {
			continue
		}
		if strings.TrimSpace(emp.Position) != "" {
			parts := strings.SplitN(emp.Position, ".", 2)
			emp.Position = parts[0]
		}
		jobRole, err := s.getJobRoleByName(ctx, emp.Position)
		if err != nil || jobRole == nil {
			continue
		}

		gradeGroupName := ""
		if gradeGroup.JobGradeGroup != nil {
			gradeGroupName = gradeGroup.JobGradeGroup.GroupName
		}

		// Subordinate review type ID = 4 (hardcoded in .NET)
		subordinateReviewTypeID := 4
		if err := s.PopulateReviewForSubordinateReviewType(ctx, emp, gradeGroupName, period.ReviewPeriodID, subordinateReviewTypeID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to populate subordinate review")
		}
	}
	return nil
}

// PopulateAllEmployeeSuperiorReviews creates superior-review records for all employees.
func (s *reviewAgentService) PopulateAllEmployeeSuperiorReviews(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for i := range employees {
		emp := &employees[i]
		gradeGroup, err := s.getAssignedJobGradeGroup(ctx, emp.Grade)
		if err != nil || gradeGroup == nil {
			continue
		}
		if strings.TrimSpace(emp.Position) != "" {
			parts := strings.SplitN(emp.Position, ".", 2)
			emp.Position = parts[0]
		}
		jobRole, err := s.getJobRoleByName(ctx, emp.Position)
		if err != nil || jobRole == nil {
			continue
		}

		gradeGroupName := ""
		if gradeGroup.JobGradeGroup != nil {
			gradeGroupName = gradeGroup.JobGradeGroup.GroupName
		}

		// Superior review type ID = 5 (hardcoded in .NET)
		superiorReviewTypeID := 5
		if err := s.PopulateReviewForSuperiorReviewType(ctx, emp, gradeGroupName, period.ReviewPeriodID, superiorReviewTypeID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to populate superior review")
		}
	}
	return nil
}

// CalculateReviewsProfileForAllEmployees computes both behavioral and technical
// review averages for every active employee.
func (s *reviewAgentService) CalculateReviewsProfileForAllEmployees(ctx context.Context) error {
	period, err := s.getCurrentReviewPeriod(ctx)
	if err != nil || period == nil {
		return err
	}
	employees, err := s.getAllEmployees(ctx)
	if err != nil || len(employees) == 0 {
		return err
	}

	for _, emp := range employees {
		if err := s.CalculateTechnicalReviewAverage(ctx, emp.EmployeeNumber, period.ReviewPeriodID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to calculate technical review average")
		}
		if err := s.CalculateBehavioralReviewAverage(ctx, emp.EmployeeNumber, period.ReviewPeriodID); err != nil {
			s.log.Error().Err(err).Str("employee", emp.EmployeeNumber).Msg("failed to calculate behavioral review average")
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// 6. Partial Match (Job Role Matching)
// ---------------------------------------------------------------------------

// isPartialMatch replicates the .NET IsPartialMatch method. It splits two
// dot-delimited strings, extracts middle and last segments, and compares
// word overlap using bidirectional scoring.
func isPartialMatch(input1, input2 string) bool {
	parts1 := splitAndTrim(input1, '.')
	parts2 := splitAndTrim(input2, '.')

	middle1, last1 := extractMiddleLast(parts1)
	middle2, last2 := extractMiddleLast(parts2)

	middleScore := compareSegments(middle1, middle2)
	lastScore := compareSegments(last1, last2)

	return middleScore >= 20 && lastScore >= 10
}

// splitAndTrim splits a string by separator, removes empty entries, and trims.
func splitAndTrim(s string, sep byte) []string {
	raw := strings.Split(s, string(sep))
	var out []string
	for _, part := range raw {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// extractMiddleLast mirrors the .NET logic for extracting middle and last
// segments from a dot-split list.
func extractMiddleLast(parts []string) (middle, last string) {
	switch {
	case len(parts) >= 3:
		middle = parts[1]
		last = parts[len(parts)-1]
	case len(parts) == 2:
		middle = parts[0]
		last = parts[1]
	case len(parts) == 1:
		middle = parts[0]
		last = parts[0]
	default:
		middle = ""
		last = ""
	}
	return
}

// compareSegments computes a bidirectional word-overlap percentage between two
// space-delimited segments. Returns the average of both directional scores.
func compareSegments(seg1, seg2 string) float64 {
	words1 := splitWords(seg1)
	words2 := splitWords(seg2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}

	// Build set from words2
	set2 := make(map[string]bool, len(words2))
	for _, w := range words2 {
		set2[w] = true
	}

	// Count common words
	var commonCount int
	for _, w := range words1 {
		if set2[w] {
			commonCount++
		}
	}

	score1 := float64(commonCount) / float64(len(words1)) * 100
	score2 := float64(commonCount) / float64(len(words2)) * 100

	return (score1 + score2) / 2
}

// splitWords splits a string by spaces, removes empties, and uppercases.
func splitWords(s string) []string {
	raw := strings.Fields(s)
	out := make([]string, 0, len(raw))
	for _, w := range raw {
		upper := strings.ToUpper(strings.TrimSpace(w))
		if upper != "" {
			out = append(out, upper)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// 7. Head Subordinates Helper
// ---------------------------------------------------------------------------

// GetHeadSubordinates returns the subordinates of an employee based on their
// organisational head role (head-of-office, head-of-division, head-of-department).
func (s *reviewAgentService) GetHeadSubordinates(ctx context.Context, employeeNumber string) ([]erp.EmployeeErpDetailsDTO, error) {
	emp, err := s.getEmployeeDetail(ctx, employeeNumber)
	if err != nil || emp == nil {
		return nil, err
	}

	if strings.EqualFold(emp.EmployeeNumber, emp.HeadOfOfficeID) {
		return s.getEmployeesByOfficeID(ctx, emp.OfficeID)
	}
	if strings.EqualFold(emp.EmployeeNumber, emp.HeadOfDivID) && emp.DivisionID != nil {
		return s.getHeadOfDivisionSubordinates(ctx, *emp.DivisionID)
	}
	if strings.EqualFold(emp.EmployeeNumber, emp.HeadOfDeptID) && emp.DepartmentID != nil {
		return s.getHeadOfDepartmentSubordinates(ctx, *emp.DepartmentID)
	}
	return nil, nil
}

// getHeadOfDivisionSubordinates returns the head-of-office employees for each
// office in the given division.
func (s *reviewAgentService) getHeadOfDivisionSubordinates(ctx context.Context, divisionID int) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE DivisionId = @p1 AND PersonTypeId = 1120
	AND EmployeeNumber = HeadOfOfficeId`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, divisionID)
}

// getHeadOfDepartmentSubordinates returns the head-of-division employees for
// each division in the given department.
func (s *reviewAgentService) getHeadOfDepartmentSubordinates(ctx context.Context, departmentID int) ([]erp.EmployeeErpDetailsDTO, error) {
	const q = `SELECT
		UserName AS userName, EmailAddress AS emailAddress,
		FirstName AS firstName, MiddleNames AS middleNames, LastName AS lastName,
		EmployeeNumber AS employeeNumber, JobName AS jobName,
		DepartmentName AS departmentName, DivisionName AS divisionName,
		HeadOfDivName AS headOfDivName, OfficeName AS officeName,
		ISNULL(SupervisorId,'') AS supervisorId,
		ISNULL(HeadOfOfficeId,'') AS headOfOfficeId,
		ISNULL(HeadOfDivId,'') AS headOfDivId,
		ISNULL(HeadOfDeptId,'') AS headOfDeptId,
		DepartmentId AS departmentId, OfficeId AS officeId,
		Grade AS grade, DivisionId AS divisionId,
		Position AS position
	FROM dbo.EmployeeDetails
	WHERE DepartmentId = @p1 AND PersonTypeId = 1120
	AND EmployeeNumber = HeadOfDivId`
	return repository.RawQuery[erp.EmployeeErpDetailsDTO](s.repos.ErpSQL, ctx, q, departmentID)
}

// ---------------------------------------------------------------------------
// Utility
// ---------------------------------------------------------------------------

// intPtrToString converts an *int to its string representation, or "" if nil.
func intPtrToString(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

// Ensure the unused import for time is used (DateCreated in CompetencyReview uses it).
var _ = time.Now
