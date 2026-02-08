package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// organogramService manages the organisational hierarchy:
// Directorates -> Departments -> Divisions -> Offices.
// Mirrors the .NET OrganogramHandlers (SaveDirectorateHandler, GetDirectorateQueryHandler, etc.)
type organogramService struct {
	directorateRepo *repository.Repository[organogram.Directorate]
	departmentRepo  *repository.Repository[organogram.Department]
	divisionRepo    *repository.Repository[organogram.Division]
	officeRepo      *repository.Repository[organogram.Office]
	db              *gorm.DB
	log             zerolog.Logger
}

// newOrganogramService creates an OrganogramService with all required repositories.
func newOrganogramService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) OrganogramService {
	return &organogramService{
		directorateRepo: repository.NewRepository[organogram.Directorate](repos.GormDB),
		departmentRepo:  repository.NewRepository[organogram.Department](repos.GormDB),
		divisionRepo:    repository.NewRepository[organogram.Division](repos.GormDB),
		officeRepo:      repository.NewRepository[organogram.Office](repos.GormDB),
		db:              repos.GormDB,
		log:             log.With().Str("service", "organogram").Logger(),
	}
}

// ---------------------------------------------------------------------------
// Directorates
// ---------------------------------------------------------------------------

// GetDirectorates retrieves all non-deleted directorates.
// Mirrors .NET GetDirectorateQueryHandler which projects to DirectorateVm.
func (s *organogramService) GetDirectorates(ctx context.Context) (interface{}, error) {
	directorates, err := s.directorateRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve directorates")
		return nil, fmt.Errorf("retrieving directorates: %w", err)
	}

	vms := make([]organogram.DirectorateVm, 0, len(directorates))
	for _, d := range directorates {
		vms = append(vms, organogram.DirectorateVm{
			DirectorateID:   d.DirectorateID,
			DirectorateName: d.DirectorateName,
			DirectorateCode: d.DirectorateCode,
			IsActive:        d.IsActive,
		})
	}

	s.log.Debug().Int("count", len(vms)).Msg("directorates retrieved")
	return vms, nil
}

// SaveDirectorate creates or updates a directorate.
// Mirrors .NET SaveDirectorateHandler: if DirectorateId > 0 it updates, otherwise creates.
func (s *organogramService) SaveDirectorate(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*organogram.DirectorateVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *organogram.DirectorateVm")
	}

	if vm.DirectorateID > 0 {
		// Update existing directorate
		existing, err := s.directorateRepo.GetByID(ctx, vm.DirectorateID)
		if err != nil {
			return nil, fmt.Errorf("retrieving directorate for update: %w", err)
		}
		if existing == nil {
			return nil, fmt.Errorf("directorate with id %d not found", vm.DirectorateID)
		}

		existing.DirectorateCode = strings.ToUpper(vm.DirectorateCode)
		existing.DirectorateName = vm.DirectorateName

		if err := s.directorateRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", vm.DirectorateID).Msg("failed to update directorate")
			return nil, fmt.Errorf("updating directorate: %w", err)
		}

		s.log.Info().Int("id", vm.DirectorateID).Str("name", existing.DirectorateName).Msg("directorate updated")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   fmt.Sprintf("%s Directorate has been Updated Successfully", existing.DirectorateName),
		}, nil
	}

	// Create new directorate
	entity := &organogram.Directorate{
		DirectorateCode: strings.ToUpper(vm.DirectorateCode),
		DirectorateName: vm.DirectorateName,
	}

	if err := s.directorateRepo.Create(ctx, entity); err != nil {
		s.log.Error().Err(err).Str("name", vm.DirectorateName).Msg("failed to create directorate")
		return nil, fmt.Errorf("creating directorate: %w", err)
	}

	s.log.Info().Str("name", entity.DirectorateName).Msg("directorate created")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Directorate has been Created Successfully", entity.DirectorateName),
	}, nil
}

// DeleteDirectorate soft-deletes or hard-deletes a directorate.
// Mirrors .NET DeleteDirectorateHandler: supports IsSoftDelete flag.
func (s *organogramService) DeleteDirectorate(ctx context.Context, id int, isSoftDelete bool) (interface{}, error) {
	existing, err := s.directorateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving directorate for delete: %w", err)
	}
	if existing == nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Directorate with id %d not found", id),
		}, nil
	}

	if isSoftDelete {
		existing.SoftDeleted = true
		existing.IsActive = false
		if err := s.directorateRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to soft-delete directorate")
			return nil, fmt.Errorf("soft-deleting directorate: %w", err)
		}
	} else {
		if err := s.directorateRepo.HardDelete(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to hard-delete directorate")
			return nil, fmt.Errorf("hard-deleting directorate: %w", err)
		}
	}

	s.log.Info().Int("id", id).Bool("softDelete", isSoftDelete).Msg("directorate deleted")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Directorate has been deleted successfully", existing.DirectorateName),
	}, nil
}

// ---------------------------------------------------------------------------
// Departments
// ---------------------------------------------------------------------------

// GetDepartments retrieves all non-deleted departments with their parent directorate name.
// Optionally filtered by directorateId.
// Mirrors .NET GetDepartmentQueryHandler which includes Directorate via Include and
// optionally filters by DirectorateId.
func (s *organogramService) GetDepartments(ctx context.Context, directorateId *int) (interface{}, error) {
	q := s.departmentRepo.Preload(ctx, "Directorate")

	if directorateId != nil {
		q = q.Where("directorate_id = ?", *directorateId)
	}

	var departments []organogram.Department
	if err := q.Find(&departments).Error; err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve departments")
		return nil, fmt.Errorf("retrieving departments: %w", err)
	}

	vms := make([]organogram.DepartmentVm, 0, len(departments))
	for _, d := range departments {
		vm := organogram.DepartmentVm{
			DepartmentID:   d.DepartmentID,
			DepartmentName: d.DepartmentName,
			DepartmentCode: d.DepartmentCode,
			DirectorateID:  d.DirectorateID,
			IsBranch:       d.IsBranch,
			IsActive:       d.IsActive,
		}
		if d.Directorate != nil {
			vm.DirectorateName = d.Directorate.DirectorateName
		}
		vms = append(vms, vm)
	}

	s.log.Debug().Int("count", len(vms)).Msg("departments retrieved")
	return vms, nil
}

// SaveDepartment creates or updates a department.
// Mirrors .NET SaveDepartmentHandler: checks duplicate DepartmentCode on create,
// updates fields on existing record when DepartmentId > 0.
func (s *organogramService) SaveDepartment(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*organogram.DepartmentVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *organogram.DepartmentVm")
	}

	if vm.DepartmentID > 0 {
		// Update existing department
		existing, err := s.departmentRepo.GetByID(ctx, vm.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("retrieving department for update: %w", err)
		}
		if existing == nil {
			return nil, fmt.Errorf("department with id %d not found", vm.DepartmentID)
		}

		existing.DepartmentCode = strings.ToUpper(vm.DepartmentCode)
		existing.DepartmentName = vm.DepartmentName
		existing.DirectorateID = vm.DirectorateID
		existing.IsBranch = vm.IsBranch

		if err := s.departmentRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", vm.DepartmentID).Msg("failed to update department")
			return nil, fmt.Errorf("updating department: %w", err)
		}

		s.log.Info().Int("id", vm.DepartmentID).Str("name", existing.DepartmentName).Msg("department updated")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   fmt.Sprintf("%s Department has been Updated Successfully", existing.DepartmentName),
		}, nil
	}

	// Check for duplicate department code before creating
	duplicate, err := s.departmentRepo.FirstOrDefault(ctx, "department_code = ?", vm.DepartmentCode)
	if err != nil {
		return nil, fmt.Errorf("checking duplicate department code: %w", err)
	}
	if duplicate != nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Department Code: %s already exists", vm.DepartmentCode),
		}, nil
	}

	entity := &organogram.Department{
		DepartmentCode: strings.ToUpper(vm.DepartmentCode),
		DepartmentName: vm.DepartmentName,
		DirectorateID:  vm.DirectorateID,
		IsBranch:       vm.IsBranch,
	}
	entity.IsActive = vm.IsActive

	if err := s.departmentRepo.Create(ctx, entity); err != nil {
		s.log.Error().Err(err).Str("name", vm.DepartmentName).Msg("failed to create department")
		return nil, fmt.Errorf("creating department: %w", err)
	}

	s.log.Info().Str("name", entity.DepartmentName).Msg("department created")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Department has been Created Successfully", entity.DepartmentName),
	}, nil
}

// DeleteDepartment soft-deletes or hard-deletes a department.
// Mirrors .NET DeleteDepartmentHandler: supports IsSoftDelete flag.
func (s *organogramService) DeleteDepartment(ctx context.Context, id int, isSoftDelete bool) (interface{}, error) {
	existing, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving department for delete: %w", err)
	}
	if existing == nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Department with id %d not found", id),
		}, nil
	}

	if isSoftDelete {
		existing.SoftDeleted = true
		existing.IsActive = false
		if err := s.departmentRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to soft-delete department")
			return nil, fmt.Errorf("soft-deleting department: %w", err)
		}
	} else {
		if err := s.departmentRepo.HardDelete(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to hard-delete department")
			return nil, fmt.Errorf("hard-deleting department: %w", err)
		}
	}

	s.log.Info().Int("id", id).Bool("softDelete", isSoftDelete).Msg("department deleted")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Department has been deleted successfully", existing.DepartmentName),
	}, nil
}

// ---------------------------------------------------------------------------
// Divisions
// ---------------------------------------------------------------------------

// GetDivisions retrieves all non-deleted divisions, optionally filtered by departmentId.
// Mirrors .NET GetDivisionQueryHandler which includes Department via Include and filters by DepartmentId.
func (s *organogramService) GetDivisions(ctx context.Context, departmentId *int) (interface{}, error) {
	q := s.divisionRepo.Preload(ctx, "Department")

	if departmentId != nil {
		q = q.Where("department_id = ?", *departmentId)
	}

	var divisions []organogram.Division
	if err := q.Find(&divisions).Error; err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve divisions")
		return nil, fmt.Errorf("retrieving divisions: %w", err)
	}

	vms := make([]organogram.DivisionVm, 0, len(divisions))
	for _, d := range divisions {
		vm := organogram.DivisionVm{
			DivisionID:   d.DivisionID,
			DivisionName: d.DivisionName,
			DivisionCode: d.DivisionCode,
			DepartmentID: d.DepartmentID,
			IsActive:     d.IsActive,
		}
		if d.Department != nil {
			vm.DepartmentName = d.Department.DepartmentName
		}
		vms = append(vms, vm)
	}

	s.log.Debug().Int("count", len(vms)).Msg("divisions retrieved")
	return vms, nil
}

// SaveDivision creates or updates a division.
// Mirrors .NET AddOrUpdatDivisionHandler: checks duplicate DivisionCode on create,
// updates fields on existing record when DivisionId > 0.
func (s *organogramService) SaveDivision(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*organogram.DivisionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *organogram.DivisionVm")
	}

	if vm.DivisionID > 0 {
		// Update existing division
		existing, err := s.divisionRepo.GetByID(ctx, vm.DivisionID)
		if err != nil {
			return nil, fmt.Errorf("retrieving division for update: %w", err)
		}
		if existing == nil {
			return nil, fmt.Errorf("division with id %d not found", vm.DivisionID)
		}

		existing.DivisionCode = strings.ToUpper(vm.DivisionCode)
		existing.DivisionName = vm.DivisionName
		existing.DepartmentID = vm.DepartmentID
		existing.IsActive = vm.IsActive

		if err := s.divisionRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", vm.DivisionID).Msg("failed to update division")
			return nil, fmt.Errorf("updating division: %w", err)
		}

		s.log.Info().Int("id", vm.DivisionID).Str("name", existing.DivisionName).Msg("division updated")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   fmt.Sprintf("%s Division has been Updated Successfully", existing.DivisionName),
		}, nil
	}

	// Check for duplicate division code before creating
	duplicate, err := s.divisionRepo.FirstOrDefault(ctx, "division_code = ?", vm.DivisionCode)
	if err != nil {
		return nil, fmt.Errorf("checking duplicate division code: %w", err)
	}
	if duplicate != nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Division Code: %s already exists", vm.DivisionCode),
		}, nil
	}

	entity := &organogram.Division{
		DivisionCode: strings.ToUpper(vm.DivisionCode),
		DivisionName: vm.DivisionName,
		DepartmentID: vm.DepartmentID,
	}
	entity.IsActive = vm.IsActive

	if err := s.divisionRepo.Create(ctx, entity); err != nil {
		s.log.Error().Err(err).Str("name", vm.DivisionName).Msg("failed to create division")
		return nil, fmt.Errorf("creating division: %w", err)
	}

	s.log.Info().Str("name", entity.DivisionName).Msg("division created")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Division has been Created Successfully", entity.DivisionName),
	}, nil
}

// DeleteDivision soft-deletes or hard-deletes a division.
// Mirrors .NET DeleteDivisionHandler: supports IsSoftDelete flag.
func (s *organogramService) DeleteDivision(ctx context.Context, id int, isSoftDelete bool) (interface{}, error) {
	existing, err := s.divisionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving division for delete: %w", err)
	}
	if existing == nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Division with id %d not found", id),
		}, nil
	}

	if isSoftDelete {
		existing.SoftDeleted = true
		existing.IsActive = false
		if err := s.divisionRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to soft-delete division")
			return nil, fmt.Errorf("soft-deleting division: %w", err)
		}
	} else {
		if err := s.divisionRepo.HardDelete(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to hard-delete division")
			return nil, fmt.Errorf("hard-deleting division: %w", err)
		}
	}

	s.log.Info().Int("id", id).Bool("softDelete", isSoftDelete).Msg("division deleted")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Division has been deleted successfully", existing.DivisionName),
	}, nil
}

// ---------------------------------------------------------------------------
// Offices
// ---------------------------------------------------------------------------

// GetOffices retrieves all non-deleted offices, optionally filtered by divisionId.
// Mirrors .NET GetOfficeQueryHandler which includes Division via Include and filters by DivisionId.
func (s *organogramService) GetOffices(ctx context.Context, divisionId *int) (interface{}, error) {
	q := s.officeRepo.Preload(ctx, "Division")

	if divisionId != nil {
		q = q.Where("division_id = ?", *divisionId)
	}

	var offices []organogram.Office
	if err := q.Find(&offices).Error; err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve offices")
		return nil, fmt.Errorf("retrieving offices: %w", err)
	}

	vms := make([]organogram.OfficeVm, 0, len(offices))
	for _, o := range offices {
		vm := organogram.OfficeVm{
			OfficeID:   o.OfficeID,
			OfficeName: o.OfficeName,
			OfficeCode: o.OfficeCode,
			DivisionID: o.DivisionID,
			IsActive:   o.IsActive,
		}
		if o.Division != nil {
			vm.DivisionName = o.Division.DivisionName
		}
		vms = append(vms, vm)
	}

	s.log.Debug().Int("count", len(vms)).Msg("offices retrieved")
	return vms, nil
}

// GetOfficeByCode retrieves a single office by its office code.
// Mirrors .NET GetOfficeByIdQueryHandler which queries by OfficeCode and includes Division.
func (s *organogramService) GetOfficeByCode(ctx context.Context, officeCode string) (interface{}, error) {
	var office organogram.Office
	err := s.officeRepo.Preload(ctx, "Division").
		Where("office_code = ? AND soft_deleted = ?", officeCode, false).
		First(&office).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		s.log.Error().Err(err).Str("officeCode", officeCode).Msg("failed to retrieve office by code")
		return nil, fmt.Errorf("retrieving office by code: %w", err)
	}

	vm := organogram.OfficeVm{
		OfficeID:   office.OfficeID,
		OfficeName: office.OfficeName,
		OfficeCode: office.OfficeCode,
		DivisionID: office.DivisionID,
		IsActive:   office.IsActive,
	}
	if office.Division != nil {
		vm.DivisionName = office.Division.DivisionName
	}

	return &vm, nil
}

// SaveOffice creates or updates an office.
// Mirrors .NET AddOrUpdateOfficeHandler: checks duplicate OfficeCode on create,
// updates fields on existing record when OfficeId > 0.
func (s *organogramService) SaveOffice(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*organogram.OfficeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *organogram.OfficeVm")
	}

	if vm.OfficeID > 0 {
		// Update existing office
		existing, err := s.officeRepo.GetByID(ctx, vm.OfficeID)
		if err != nil {
			return nil, fmt.Errorf("retrieving office for update: %w", err)
		}
		if existing == nil {
			return nil, fmt.Errorf("office with id %d not found", vm.OfficeID)
		}

		existing.OfficeCode = strings.ToUpper(vm.OfficeCode)
		existing.OfficeName = vm.OfficeName
		existing.DivisionID = vm.DivisionID
		existing.IsActive = vm.IsActive

		if err := s.officeRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", vm.OfficeID).Msg("failed to update office")
			return nil, fmt.Errorf("updating office: %w", err)
		}

		s.log.Info().Int("id", vm.OfficeID).Str("name", existing.OfficeName).Msg("office updated")
		return map[string]interface{}{
			"isSuccess": true,
			"message":   fmt.Sprintf("%s Office has been Updated Successfully", existing.OfficeName),
		}, nil
	}

	// Check for duplicate office code before creating
	duplicate, err := s.officeRepo.FirstOrDefault(ctx, "office_code = ?", vm.OfficeCode)
	if err != nil {
		return nil, fmt.Errorf("checking duplicate office code: %w", err)
	}
	if duplicate != nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Office Code: %s already exists", vm.OfficeCode),
		}, nil
	}

	entity := &organogram.Office{
		OfficeCode: strings.ToUpper(vm.OfficeCode),
		OfficeName: vm.OfficeName,
		DivisionID: vm.DivisionID,
	}
	entity.IsActive = vm.IsActive

	if err := s.officeRepo.Create(ctx, entity); err != nil {
		s.log.Error().Err(err).Str("name", vm.OfficeName).Msg("failed to create office")
		return nil, fmt.Errorf("creating office: %w", err)
	}

	s.log.Info().Str("name", entity.OfficeName).Msg("office created")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Office has been Created Successfully", entity.OfficeName),
	}, nil
}

// DeleteOffice soft-deletes or hard-deletes an office.
// Mirrors .NET DeleteOfficeHandler: supports IsSoftDelete flag.
func (s *organogramService) DeleteOffice(ctx context.Context, id int, isSoftDelete bool) (interface{}, error) {
	existing, err := s.officeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving office for delete: %w", err)
	}
	if existing == nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   fmt.Sprintf("Office with id %d not found", id),
		}, nil
	}

	if isSoftDelete {
		existing.SoftDeleted = true
		existing.IsActive = false
		if err := s.officeRepo.Update(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to soft-delete office")
			return nil, fmt.Errorf("soft-deleting office: %w", err)
		}
	} else {
		if err := s.officeRepo.HardDelete(ctx, existing); err != nil {
			s.log.Error().Err(err).Int("id", id).Msg("failed to hard-delete office")
			return nil, fmt.Errorf("hard-deleting office: %w", err)
		}
	}

	s.log.Info().Int("id", id).Bool("softDelete", isSoftDelete).Msg("office deleted")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   fmt.Sprintf("%s Office has been deleted successfully", existing.OfficeName),
	}, nil
}

func init() {
	// Compile-time interface compliance check.
	var _ OrganogramService = (*organogramService)(nil)
}
