package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// objectiveService implements objective-related methods delegated from
// performanceManagementService: enterprise/department/division/office
// objectives CRUD, category definitions, consolidated objectives, evaluation
// options, feedback questionnaires, PMS competencies, work product
// definitions, and approval/rejection workflows.
// ---------------------------------------------------------------------------

type objectiveService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	// Back-reference to parent for shared helpers / peer services
	parent *performanceManagementService

	// Repositories
	enterpriseObjRepo *repository.PMSRepository[performance.EnterpriseObjective]
	departmentObjRepo *repository.PMSRepository[performance.DepartmentObjective]
	divisionObjRepo   *repository.PMSRepository[performance.DivisionObjective]
	officeObjRepo     *repository.PMSRepository[performance.OfficeObjective]
	objCategoryRepo   *repository.PMSRepository[performance.ObjectiveCategory]
	categoryDefRepo   *repository.PMSRepository[performance.CategoryDefinition]
	periodObjRepo     *repository.PMSRepository[performance.PeriodObjective]
	plannedObjRepo    *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	evalOptionRepo    *repository.PMSRepository[performance.EvaluationOption]
	feedbackQRepo     *repository.PMSRepository[performance.FeedbackQuestionaire]
	feedbackQOptRepo  *repository.PMSRepository[performance.FeedbackQuestionaireOption]
	pmsCompetencyRepo *repository.PMSRepository[performance.PmsCompetency]
	wpDefinitionRepo  *repository.PMSRepository[performance.WorkProductDefinition]
	cascadedWPRepo    *repository.PMSRepository[performance.CascadedWorkProduct]
}

func newObjectiveService(
	db *gorm.DB,
	cfg *config.Config,
	log zerolog.Logger,
	parent *performanceManagementService,
) *objectiveService {
	return &objectiveService{
		db:     db,
		cfg:    cfg,
		log:    log.With().Str("sub", "objectives").Logger(),
		parent: parent,

		enterpriseObjRepo: repository.NewPMSRepository[performance.EnterpriseObjective](db),
		departmentObjRepo: repository.NewPMSRepository[performance.DepartmentObjective](db),
		divisionObjRepo:   repository.NewPMSRepository[performance.DivisionObjective](db),
		officeObjRepo:     repository.NewPMSRepository[performance.OfficeObjective](db),
		objCategoryRepo:   repository.NewPMSRepository[performance.ObjectiveCategory](db),
		categoryDefRepo:   repository.NewPMSRepository[performance.CategoryDefinition](db),
		periodObjRepo:     repository.NewPMSRepository[performance.PeriodObjective](db),
		plannedObjRepo:    repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](db),
		evalOptionRepo:    repository.NewPMSRepository[performance.EvaluationOption](db),
		feedbackQRepo:     repository.NewPMSRepository[performance.FeedbackQuestionaire](db),
		feedbackQOptRepo:  repository.NewPMSRepository[performance.FeedbackQuestionaireOption](db),
		pmsCompetencyRepo: repository.NewPMSRepository[performance.PmsCompetency](db),
		wpDefinitionRepo:  repository.NewPMSRepository[performance.WorkProductDefinition](db),
		cascadedWPRepo:    repository.NewPMSRepository[performance.CascadedWorkProduct](db),
	}
}

// =========================================================================
// Enterprise Objectives
// =========================================================================

func (s *objectiveService) GetEnterpriseObjectives(ctx context.Context) (interface{}, error) {
	resp := performance.ReviewPeriodObjectivesResponseVm{}

	objs, err := s.enterpriseObjRepo.GetAllIncluding(ctx, "Category", "Strategy", "StrategicTheme")
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get enterprise objectives")
		resp.HasError = true
		resp.Message = "failed to retrieve enterprise objectives"
		return resp, err
	}

	var data []performance.EnterpriseObjectiveData
	for _, obj := range objs {
		if obj.SoftDeleted {
			continue
		}
		d := performance.EnterpriseObjectiveData{
			EnterpriseObjectiveID:          obj.EnterpriseObjectiveID,
			Name:                           obj.Name,
			Description:                    obj.Description,
			Kpi:                            obj.Kpi,
			Target:                         obj.Target,
			EnterpriseObjectivesCategoryID: obj.EnterpriseObjectivesCategoryID,
			StrategyID:                     obj.StrategyID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.EnterpriseObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateEnterpriseObjective")
	}

	// Validate strategy exists
	var strategy performance.Strategy
	if err := s.db.WithContext(ctx).
		Where("strategy_id = ?", r.StrategyID).
		First(&strategy).Error; err != nil {
		return nil, fmt.Errorf("strategy not found: %w", err)
	}

	// Validate category exists
	var category performance.ObjectiveCategory
	if err := s.db.WithContext(ctx).
		Where("objective_category_id = ?", r.EnterpriseObjectivesCategoryID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("objective category not found: %w", err)
	}

	// Check duplicate name
	var existing performance.EnterpriseObjective
	err := s.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND strategy_id = ? AND record_status != ?",
			r.Name, r.StrategyID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("enterprise objective with this name already exists")
	}

	obj := performance.EnterpriseObjective{
		EnterpriseObjectivesCategoryID: r.EnterpriseObjectivesCategoryID,
		StrategyID:                     r.StrategyID,
		Type:                           enums.ObjectiveTypeEnterprise,
	}
	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target
	obj.RecordStatus = enums.StatusPendingApproval.String()
	obj.IsActive = true

	if err := s.db.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, fmt.Errorf("creating enterprise objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.EnterpriseObjectiveID
	resp.Message = "enterprise objective created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.EnterpriseObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateEnterpriseObjective")
	}

	var obj performance.EnterpriseObjective
	if err := s.db.WithContext(ctx).
		Where("enterprise_objective_id = ?", r.EnterpriseObjectiveID).
		First(&obj).Error; err != nil {
		return nil, fmt.Errorf("enterprise objective not found: %w", err)
	}

	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target
	obj.EnterpriseObjectivesCategoryID = r.EnterpriseObjectivesCategoryID

	if err := s.db.WithContext(ctx).Save(&obj).Error; err != nil {
		return nil, fmt.Errorf("updating enterprise objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.EnterpriseObjectiveID
	resp.Message = "enterprise objective updated successfully"
	return resp, nil
}

// =========================================================================
// Department Objectives
// =========================================================================

func (s *objectiveService) GetDepartmentObjectives(ctx context.Context) (interface{}, error) {
	resp := performance.CascadedObjectiveDataListResponseVm{}

	objs, err := s.departmentObjRepo.GetAllIncluding(ctx, "EnterpriseObjective", "Department")
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get department objectives")
		resp.HasError = true
		resp.Message = "failed to retrieve department objectives"
		return resp, err
	}

	var data []performance.CascadedObjectiveData
	for _, obj := range objs {
		if obj.SoftDeleted {
			continue
		}
		d := performance.CascadedObjectiveData{
			ObjectiveID:          obj.DepartmentObjectiveID,
			Name:                 obj.Name,
			Description:          obj.Description,
			Kpi:                  obj.Kpi,
			Target:               obj.Target,
			ObjectivesCategoryID: obj.EnterpriseObjectiveID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateDepartmentObjective")
	}

	// Validate enterprise objective exists
	var enterpriseObj performance.EnterpriseObjective
	if err := s.db.WithContext(ctx).
		Where("enterprise_objective_id = ?", r.ObjectivesCategoryID).
		First(&enterpriseObj).Error; err != nil {
		return nil, fmt.Errorf("enterprise objective not found: %w", err)
	}

	obj := performance.DepartmentObjective{
		EnterpriseObjectiveID: r.ObjectivesCategoryID,
	}
	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target
	obj.RecordStatus = enums.StatusPendingApproval.String()
	obj.IsActive = true

	if err := s.db.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, fmt.Errorf("creating department objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.DepartmentObjectiveID
	resp.Message = "department objective created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateDepartmentObjective")
	}

	var obj performance.DepartmentObjective
	if err := s.db.WithContext(ctx).
		Where("department_objective_id = ?", r.ObjectiveID).
		First(&obj).Error; err != nil {
		return nil, fmt.Errorf("department objective not found: %w", err)
	}

	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target

	if err := s.db.WithContext(ctx).Save(&obj).Error; err != nil {
		return nil, fmt.Errorf("updating department objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.DepartmentObjectiveID
	resp.Message = "department objective updated successfully"
	return resp, nil
}

// =========================================================================
// Division Objectives
// =========================================================================

func (s *objectiveService) GetDivisionObjectives(ctx context.Context) (interface{}, error) {
	resp := performance.CascadedObjectiveDataListResponseVm{}

	objs, err := s.divisionObjRepo.GetAllIncluding(ctx, "DepartmentObjective", "Division")
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get division objectives")
		resp.HasError = true
		resp.Message = "failed to retrieve division objectives"
		return resp, err
	}

	var data []performance.CascadedObjectiveData
	for _, obj := range objs {
		if obj.SoftDeleted {
			continue
		}
		d := performance.CascadedObjectiveData{
			ObjectiveID:          obj.DivisionObjectiveID,
			Name:                 obj.Name,
			Description:          obj.Description,
			Kpi:                  obj.Kpi,
			Target:               obj.Target,
			ObjectivesCategoryID: obj.DepartmentObjectiveID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) GetDivisionObjectivesByDivisionId(ctx context.Context, divisionID int) (interface{}, error) {
	resp := performance.CascadedObjectiveDataListResponseVm{}

	var objs []performance.DivisionObjective
	err := s.db.WithContext(ctx).
		Where("division_id = ? AND record_status != ?", divisionID, enums.StatusCancelled.String()).
		Preload("DepartmentObjective").
		Preload("Division").
		Find(&objs).Error
	if err != nil {
		s.log.Error().Err(err).Int("divisionID", divisionID).Msg("failed to get division objectives by division")
		resp.HasError = true
		resp.Message = "failed to retrieve division objectives"
		return resp, err
	}

	var data []performance.CascadedObjectiveData
	for _, obj := range objs {
		d := performance.CascadedObjectiveData{
			ObjectiveID:          obj.DivisionObjectiveID,
			Name:                 obj.Name,
			Description:          obj.Description,
			Kpi:                  obj.Kpi,
			Target:               obj.Target,
			ObjectivesCategoryID: obj.DepartmentObjectiveID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateDivisionObjective")
	}

	var deptObj performance.DepartmentObjective
	if err := s.db.WithContext(ctx).
		Where("department_objective_id = ?", r.ObjectivesCategoryID).
		First(&deptObj).Error; err != nil {
		return nil, fmt.Errorf("department objective not found: %w", err)
	}

	obj := performance.DivisionObjective{
		DepartmentObjectiveID: r.ObjectivesCategoryID,
	}
	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target
	obj.RecordStatus = enums.StatusPendingApproval.String()
	obj.IsActive = true

	if err := s.db.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, fmt.Errorf("creating division objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.DivisionObjectiveID
	resp.Message = "division objective created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateDivisionObjective")
	}

	var obj performance.DivisionObjective
	if err := s.db.WithContext(ctx).
		Where("division_objective_id = ?", r.ObjectiveID).
		First(&obj).Error; err != nil {
		return nil, fmt.Errorf("division objective not found: %w", err)
	}

	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target

	if err := s.db.WithContext(ctx).Save(&obj).Error; err != nil {
		return nil, fmt.Errorf("updating division objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.DivisionObjectiveID
	resp.Message = "division objective updated successfully"
	return resp, nil
}

// =========================================================================
// Office Objectives
// =========================================================================

func (s *objectiveService) GetOfficeObjectives(ctx context.Context) (interface{}, error) {
	resp := performance.CascadedObjectiveDataListResponseVm{}

	objs, err := s.officeObjRepo.GetAllIncluding(ctx, "DivisionObjective", "Office")
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get office objectives")
		resp.HasError = true
		resp.Message = "failed to retrieve office objectives"
		return resp, err
	}

	var data []performance.CascadedObjectiveData
	for _, obj := range objs {
		if obj.SoftDeleted {
			continue
		}
		d := performance.CascadedObjectiveData{
			ObjectiveID:          obj.OfficeObjectiveID,
			Name:                 obj.Name,
			Description:          obj.Description,
			Kpi:                  obj.Kpi,
			Target:               obj.Target,
			ObjectivesCategoryID: obj.DivisionObjectiveID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) GetOfficeObjectivesByOfficeId(ctx context.Context, officeID int) (interface{}, error) {
	resp := performance.CascadedObjectiveDataListResponseVm{}

	var objs []performance.OfficeObjective
	err := s.db.WithContext(ctx).
		Where("office_id = ? AND record_status != ?", officeID, enums.StatusCancelled.String()).
		Preload("DivisionObjective").
		Preload("Office").
		Find(&objs).Error
	if err != nil {
		s.log.Error().Err(err).Int("officeID", officeID).Msg("failed to get office objectives by office")
		resp.HasError = true
		resp.Message = "failed to retrieve office objectives"
		return resp, err
	}

	var data []performance.CascadedObjectiveData
	for _, obj := range objs {
		d := performance.CascadedObjectiveData{
			ObjectiveID:          obj.OfficeObjectiveID,
			Name:                 obj.Name,
			Description:          obj.Description,
			Kpi:                  obj.Kpi,
			Target:               obj.Target,
			ObjectivesCategoryID: obj.DivisionObjectiveID,
		}
		data = append(data, d)
	}

	resp.Objectives = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateOfficeObjective")
	}

	var divObj performance.DivisionObjective
	if err := s.db.WithContext(ctx).
		Where("division_objective_id = ?", r.ObjectivesCategoryID).
		First(&divObj).Error; err != nil {
		return nil, fmt.Errorf("division objective not found: %w", err)
	}

	obj := performance.OfficeObjective{
		DivisionObjectiveID: r.ObjectivesCategoryID,
	}
	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target
	obj.RecordStatus = enums.StatusPendingApproval.String()
	obj.IsActive = true

	if err := s.db.WithContext(ctx).Create(&obj).Error; err != nil {
		return nil, fmt.Errorf("creating office objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.OfficeObjectiveID
	resp.Message = "office objective created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CascadedObjectiveData)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateOfficeObjective")
	}

	var obj performance.OfficeObjective
	if err := s.db.WithContext(ctx).
		Where("office_objective_id = ?", r.ObjectiveID).
		First(&obj).Error; err != nil {
		return nil, fmt.Errorf("office objective not found: %w", err)
	}

	obj.Name = r.Name
	obj.Description = r.Description
	obj.Kpi = r.Kpi
	obj.Target = r.Target

	if err := s.db.WithContext(ctx).Save(&obj).Error; err != nil {
		return nil, fmt.Errorf("updating office objective: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.OfficeObjectiveID
	resp.Message = "office objective updated successfully"
	return resp, nil
}

// =========================================================================
// Objective Categories & Category Definitions
// =========================================================================

func (s *objectiveService) GetObjectiveCategories(ctx context.Context) (interface{}, error) {
	resp := performance.GenericListVm{}

	categories, err := s.objCategoryRepo.GetAllIncluding(ctx, "CategoryDefinitions")
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get objective categories")
		resp.HasError = true
		resp.Message = "failed to retrieve objective categories"
		return resp, err
	}

	resp.ListData = categories
	resp.TotalRecord = len(categories)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.ObjectiveCategory)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateObjectiveCategory")
	}

	// Check duplicate name
	var existing performance.ObjectiveCategory
	err := s.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND record_status != ?", r.Name, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("objective category with this name already exists")
	}

	r.RecordStatus = enums.StatusActive.String()
	r.IsActive = true

	if err := s.db.WithContext(ctx).Create(r).Error; err != nil {
		return nil, fmt.Errorf("creating objective category: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = r.ObjectiveCategoryID
	resp.Message = "objective category created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.ObjectiveCategory)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateObjectiveCategory")
	}

	var obj performance.ObjectiveCategory
	if err := s.db.WithContext(ctx).
		Where("objective_category_id = ?", r.ObjectiveCategoryID).
		First(&obj).Error; err != nil {
		return nil, fmt.Errorf("objective category not found: %w", err)
	}

	obj.Name = r.Name
	obj.Description = r.Description

	if err := s.db.WithContext(ctx).Save(&obj).Error; err != nil {
		return nil, fmt.Errorf("updating objective category: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = obj.ObjectiveCategoryID
	resp.Message = "objective category updated successfully"
	return resp, nil
}

func (s *objectiveService) GetCategoryDefinitions(ctx context.Context, categoryID string) (interface{}, error) {
	resp := performance.ReviewPeriodCategoryDefinitionResponseVm{}

	var defs []performance.CategoryDefinition
	err := s.db.WithContext(ctx).
		Where("objective_category_id = ?", categoryID).
		Preload("Category").
		Preload("ReviewPeriod").
		Find(&defs).Error
	if err != nil {
		s.log.Error().Err(err).Str("categoryID", categoryID).Msg("failed to get category definitions")
		resp.HasError = true
		resp.Message = "failed to retrieve category definitions"
		return resp, err
	}

	var data []performance.CategoryDefinitionData
	for _, def := range defs {
		d := performance.CategoryDefinitionData{
			DefinitionID:            def.DefinitionID,
			ObjectiveCategoryID:     def.ObjectiveCategoryID,
			ReviewPeriodID:          def.ReviewPeriodID,
			Weight:                  def.Weight,
			MaxNoObjectives:         def.MaxNoObjectives,
			MaxNoWorkProduct:        def.MaxNoWorkProduct,
			MaxPoints:               def.MaxPoints,
			IsCompulsory:            def.IsCompulsory,
			EnforceWorkProductLimit: def.EnforceWorkProductLimit,
			Description:             def.Description,
			GradeGroupID:            def.GradeGroupID,
		}
		if def.Category != nil {
			d.CategoryName = def.Category.Name
		}
		data = append(data, d)
	}

	resp.CategoryDefinitions = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CategoryDefinitionData2)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateCategoryDefinition")
	}

	// Validate category exists
	var category performance.ObjectiveCategory
	if err := s.db.WithContext(ctx).
		Where("objective_category_id = ?", r.ObjectiveCategoryID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("objective category not found: %w", err)
	}

	// Check for duplicate definition in the same review period and grade group
	var existing performance.CategoryDefinition
	err := s.db.WithContext(ctx).
		Where("objective_category_id = ? AND review_period_id = ? AND grade_group_id = ? AND record_status != ?",
			r.ObjectiveCategoryID, r.ReviewPeriodID, r.GradeGroupID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("category definition already exists for this category, review period and grade group")
	}

	def := performance.CategoryDefinition{
		ObjectiveCategoryID:     r.ObjectiveCategoryID,
		ReviewPeriodID:          r.ReviewPeriodID,
		Weight:                  r.Weight,
		MaxNoObjectives:         r.MaxNoObjectives,
		MaxNoWorkProduct:        r.MaxNoWorkProduct,
		MaxPoints:               r.MaxPoints,
		IsCompulsory:            r.IsCompulsory,
		EnforceWorkProductLimit: r.EnforceWorkProductLimit,
		Description:             r.Description,
		GradeGroupID:            r.GradeGroupID,
	}
	def.RecordStatus = enums.StatusPendingApproval.String()
	def.IsActive = true

	if err := s.db.WithContext(ctx).Create(&def).Error; err != nil {
		return nil, fmt.Errorf("creating category definition: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = def.DefinitionID
	resp.Message = "category definition created successfully"
	return resp, nil
}

func (s *objectiveService) UpdateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.CategoryDefinitionData2)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateCategoryDefinition")
	}

	var def performance.CategoryDefinition
	if err := s.db.WithContext(ctx).
		Where("definition_id = ?", r.DefinitionID).
		First(&def).Error; err != nil {
		return nil, fmt.Errorf("category definition not found: %w", err)
	}

	def.Weight = r.Weight
	def.MaxNoObjectives = r.MaxNoObjectives
	def.MaxNoWorkProduct = r.MaxNoWorkProduct
	def.MaxPoints = r.MaxPoints
	def.IsCompulsory = r.IsCompulsory
	def.EnforceWorkProductLimit = r.EnforceWorkProductLimit
	def.Description = r.Description

	if err := s.db.WithContext(ctx).Save(&def).Error; err != nil {
		return nil, fmt.Errorf("updating category definition: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = def.DefinitionID
	resp.Message = "category definition updated successfully"
	return resp, nil
}

// =========================================================================
// Consolidated Objectives
// =========================================================================

func (s *objectiveService) GetConsolidatedObjectives(ctx context.Context) (interface{}, error) {
	resp := performance.ConsolidatedObjectiveListResponseVm{}

	var allObjs []performance.ConsolidatedObjectiveVm

	// Enterprise objectives
	var enterprises []performance.EnterpriseObjective
	s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("Category").
		Find(&enterprises)

	for _, obj := range enterprises {
		vm := performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{
					RecordStatus: obj.RecordStatus,
					IsActive:     obj.IsActive,
				},
				ObjectiveID:    obj.EnterpriseObjectiveID,
				Name:           obj.Name,
				Description:    obj.Description,
				Kpi:            obj.Kpi,
				Target:         obj.Target,
				ObjectiveLevel: "Enterprise",
			},
		}
		if obj.Category != nil {
			vm.SBUName = obj.Category.Name
		}
		allObjs = append(allObjs, vm)
	}

	// Department objectives
	var departments []performance.DepartmentObjective
	s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("EnterpriseObjective").
		Preload("Department").
		Find(&departments)

	for _, obj := range departments {
		vm := performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{
					RecordStatus: obj.RecordStatus,
					IsActive:     obj.IsActive,
				},
				ObjectiveID:    obj.DepartmentObjectiveID,
				Name:           obj.Name,
				Description:    obj.Description,
				Kpi:            obj.Kpi,
				Target:         obj.Target,
				ObjectiveLevel: "Department",
			},
		}
		allObjs = append(allObjs, vm)
	}

	// Division objectives
	var divisions []performance.DivisionObjective
	s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("DepartmentObjective").
		Preload("Division").
		Find(&divisions)

	for _, obj := range divisions {
		vm := performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{
					RecordStatus: obj.RecordStatus,
					IsActive:     obj.IsActive,
				},
				ObjectiveID:    obj.DivisionObjectiveID,
				Name:           obj.Name,
				Description:    obj.Description,
				Kpi:            obj.Kpi,
				Target:         obj.Target,
				ObjectiveLevel: "Division",
			},
		}
		allObjs = append(allObjs, vm)
	}

	// Office objectives
	var offices []performance.OfficeObjective
	s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("DivisionObjective").
		Preload("Office").
		Find(&offices)

	for _, obj := range offices {
		vm := performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{
					RecordStatus: obj.RecordStatus,
					IsActive:     obj.IsActive,
				},
				ObjectiveID:    obj.OfficeObjectiveID,
				Name:           obj.Name,
				Description:    obj.Description,
				Kpi:            obj.Kpi,
				Target:         obj.Target,
				ObjectiveLevel: "Office",
			},
		}
		allObjs = append(allObjs, vm)
	}

	resp.Objectives = allObjs
	resp.TotalRecords = len(allObjs)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) GetConsolidatedObjectivesPaginated(ctx context.Context, params interface{}) (interface{}, error) {
	// Reuse the full list and paginate in-memory
	allResp, err := s.GetConsolidatedObjectives(ctx)
	if err != nil {
		return allResp, err
	}

	resp, ok := allResp.(performance.ConsolidatedObjectiveListResponseVm)
	if !ok {
		return allResp, nil
	}

	// Apply pagination from params if provided
	if p, ok := params.(*performance.BasePagedData); ok {
		total := len(resp.Objectives)
		start := p.Skip
		if start > total {
			start = total
		}
		end := start + p.PageSize
		if end > total {
			end = total
		}
		resp.Objectives = resp.Objectives[start:end]
		resp.TotalRecords = total
	}

	return resp, nil
}

func (s *objectiveService) ProcessObjectivesUpload(ctx context.Context, req interface{}) (interface{}, error) {
	resp := performance.GenericResponseVm{
		IsSuccess: false,
	}

	r, ok := req.(*performance.ObjectivesUploadRequestModel)
	if !ok {
		resp.Message = "invalid request type for ProcessObjectivesUpload"
		return resp, fmt.Errorf("invalid request type for ProcessObjectivesUpload")
	}

	if len(r.Objectives) == 0 {
		resp.Message = "no objectives provided"
		return resp, nil
	}

	// Track created counts per level
	var enterpriseCount, departmentCount, divisionCount, officeCount int

	// Caches to avoid duplicate lookups and re-creation within the same batch.
	// Keys are normalised "name|parentRef" to match existing objectives.
	enterpriseCache := make(map[string]string) // key -> EnterpriseObjectiveID
	departmentCache := make(map[string]string) // key -> DepartmentObjectiveID
	divisionCache := make(map[string]string)   // key -> DivisionObjectiveID

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Message = "failed to start transaction"
		return resp, fmt.Errorf("starting transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, row := range r.Objectives {
		// -----------------------------------------------------------------
		// 1. Enterprise Objective
		// -----------------------------------------------------------------
		var enterpriseObjID string

		if row.EObjName != "" {
			cacheKey := strings.ToLower(row.EObjName) + "|" + row.StrategyID
			if cachedID, found := enterpriseCache[cacheKey]; found {
				enterpriseObjID = cachedID
			} else {
				// Try to find existing enterprise objective with same name + strategy
				var existing performance.EnterpriseObjective
				err := tx.Where("LOWER(name) = LOWER(?) AND strategy_id = ? AND record_status != ?",
					row.EObjName, row.StrategyID, enums.StatusCancelled.String()).
					First(&existing).Error
				if err == nil {
					enterpriseObjID = existing.EnterpriseObjectiveID
				} else {
					// Resolve category ID from name
					categoryID := s.resolveObjectiveCategoryID(tx, row.EObjCategory)

					obj := performance.EnterpriseObjective{
						EnterpriseObjectivesCategoryID: categoryID,
						StrategyID:                     row.StrategyID,
						StrategicThemeID:               row.StrategicThemeID,
						Type:                           enums.ObjectiveTypeEnterprise,
					}
					obj.Name = row.EObjName
					obj.Description = row.EObjDesc
					obj.Kpi = row.EObjKPI
					obj.Target = row.EObjTarget
					obj.RecordStatus = enums.StatusActive.String()
					obj.IsActive = true
					obj.CreatedBy = r.CreatedBy

					if err := tx.Create(&obj).Error; err != nil {
						tx.Rollback()
						resp.Message = fmt.Sprintf("failed to create enterprise objective at row %d: %s", i+1, err.Error())
						resp.Errors = append(resp.Errors, resp.Message)
						return resp, fmt.Errorf("creating enterprise objective row %d: %w", i+1, err)
					}
					enterpriseObjID = obj.EnterpriseObjectiveID
					enterpriseCount++
				}
				enterpriseCache[cacheKey] = enterpriseObjID
			}
		}

		// -----------------------------------------------------------------
		// 2. Department Objective
		// -----------------------------------------------------------------
		var departmentObjID string

		if row.DeptObjName != "" && enterpriseObjID != "" {
			cacheKey := strings.ToLower(row.DeptObjName) + "|" + enterpriseObjID + "|" + fmt.Sprintf("%d", row.DepartmentID)
			if cachedID, found := departmentCache[cacheKey]; found {
				departmentObjID = cachedID
			} else {
				var existing performance.DepartmentObjective
				err := tx.Where("LOWER(name) = LOWER(?) AND enterprise_objective_id = ? AND department_id = ? AND record_status != ?",
					row.DeptObjName, enterpriseObjID, row.DepartmentID, enums.StatusCancelled.String()).
					First(&existing).Error
				if err == nil {
					departmentObjID = existing.DepartmentObjectiveID
				} else {
					obj := performance.DepartmentObjective{
						EnterpriseObjectiveID: enterpriseObjID,
						DepartmentID:          row.DepartmentID,
					}
					obj.Name = row.DeptObjName
					obj.Description = row.DeptObjDesc
					obj.Kpi = row.DeptObjKPI
					obj.Target = row.DeptObjTarget
					obj.RecordStatus = enums.StatusActive.String()
					obj.IsActive = true
					obj.CreatedBy = r.CreatedBy

					if err := tx.Create(&obj).Error; err != nil {
						tx.Rollback()
						resp.Message = fmt.Sprintf("failed to create department objective at row %d: %s", i+1, err.Error())
						resp.Errors = append(resp.Errors, resp.Message)
						return resp, fmt.Errorf("creating department objective row %d: %w", i+1, err)
					}
					departmentObjID = obj.DepartmentObjectiveID
					departmentCount++
				}
				departmentCache[cacheKey] = departmentObjID
			}
		}

		// -----------------------------------------------------------------
		// 3. Division Objective
		// -----------------------------------------------------------------
		var divisionObjID string

		if row.DivObjName != "" && departmentObjID != "" {
			cacheKey := strings.ToLower(row.DivObjName) + "|" + departmentObjID + "|" + fmt.Sprintf("%d", row.DivisionID)
			if cachedID, found := divisionCache[cacheKey]; found {
				divisionObjID = cachedID
			} else {
				var existing performance.DivisionObjective
				err := tx.Where("LOWER(name) = LOWER(?) AND department_objective_id = ? AND division_id = ? AND record_status != ?",
					row.DivObjName, departmentObjID, row.DivisionID, enums.StatusCancelled.String()).
					First(&existing).Error
				if err == nil {
					divisionObjID = existing.DivisionObjectiveID
				} else {
					obj := performance.DivisionObjective{
						DepartmentObjectiveID: departmentObjID,
						DivisionID:            row.DivisionID,
					}
					obj.Name = row.DivObjName
					obj.Description = row.DivObjDesc
					obj.Kpi = row.DivObjKPI
					obj.Target = row.DivObjTarget
					obj.RecordStatus = enums.StatusActive.String()
					obj.IsActive = true
					obj.CreatedBy = r.CreatedBy

					if err := tx.Create(&obj).Error; err != nil {
						tx.Rollback()
						resp.Message = fmt.Sprintf("failed to create division objective at row %d: %s", i+1, err.Error())
						resp.Errors = append(resp.Errors, resp.Message)
						return resp, fmt.Errorf("creating division objective row %d: %w", i+1, err)
					}
					divisionObjID = obj.DivisionObjectiveID
					divisionCount++
				}
				divisionCache[cacheKey] = divisionObjID
			}
		}

		// -----------------------------------------------------------------
		// 4. Office Objective
		// -----------------------------------------------------------------
		if row.OffObjName != "" && divisionObjID != "" {
			// Office objectives are not cached since each row may have
			// a unique office + job-grade-group combination.
			var existing performance.OfficeObjective
			err := tx.Where("LOWER(name) = LOWER(?) AND division_objective_id = ? AND office_id = ? AND job_grade_group_id = ? AND record_status != ?",
				row.OffObjName, divisionObjID, row.OfficeID, row.JobGradeGroupID, enums.StatusCancelled.String()).
				First(&existing).Error
			if err != nil {
				// Does not exist – create it
				obj := performance.OfficeObjective{
					DivisionObjectiveID: divisionObjID,
					OfficeID:            row.OfficeID,
					JobGradeGroupID:     row.JobGradeGroupID,
				}
				obj.Name = row.OffObjName
				obj.Description = row.OffObjDesc
				obj.Kpi = row.OffObjKPI
				obj.Target = row.OffObjTarget
				obj.RecordStatus = enums.StatusActive.String()
				obj.IsActive = true
				obj.CreatedBy = r.CreatedBy

				if err := tx.Create(&obj).Error; err != nil {
					tx.Rollback()
					resp.Message = fmt.Sprintf("failed to create office objective at row %d: %s", i+1, err.Error())
					resp.Errors = append(resp.Errors, resp.Message)
					return resp, fmt.Errorf("creating office objective row %d: %w", i+1, err)
				}
				officeCount++
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		resp.Message = "failed to commit transaction"
		return resp, fmt.Errorf("committing transaction: %w", err)
	}

	totalCreated := enterpriseCount + departmentCount + divisionCount + officeCount
	resp.IsSuccess = true
	resp.TotalRecords = totalCreated
	resp.Message = fmt.Sprintf("objectives upload completed successfully: %d enterprise, %d department, %d division, %d office objectives created (%d total)",
		enterpriseCount, departmentCount, divisionCount, officeCount, totalCreated)

	s.log.Info().
		Int("enterprise", enterpriseCount).
		Int("department", departmentCount).
		Int("division", divisionCount).
		Int("office", officeCount).
		Int("total", totalCreated).
		Msg("ProcessObjectivesUpload completed")

	return resp, nil
}

// resolveObjectiveCategoryID looks up an ObjectiveCategory by name and returns
// its ID. If not found, it returns the name as-is (the caller may have passed
// the ID directly).
func (s *objectiveService) resolveObjectiveCategoryID(tx *gorm.DB, nameOrID string) string {
	if nameOrID == "" {
		return ""
	}
	var cat performance.ObjectiveCategory
	if err := tx.Where("LOWER(name) = LOWER(?) AND record_status != ?",
		nameOrID, enums.StatusCancelled.String()).
		First(&cat).Error; err == nil {
		return cat.ObjectiveCategoryID
	}
	// Fall back: treat the input as an ID
	return nameOrID
}

func (s *objectiveService) DeActivateOrReactivateObjectives(ctx context.Context, req interface{}, deactivate bool) (interface{}, error) {
	r, ok := req.(*performance.ApprovalRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for DeActivateOrReactivateObjectives")
	}

	newStatus := enums.StatusActive.String()
	if deactivate {
		newStatus = enums.StatusDeactivated.String()
	}

	for _, id := range r.RecordIDs {
		switch strings.ToLower(r.EntityType) {
		case "enterprise":
			s.db.WithContext(ctx).Model(&performance.EnterpriseObjective{}).
				Where("enterprise_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": newStatus,
					"is_active":     !deactivate,
				})
		case "department":
			s.db.WithContext(ctx).Model(&performance.DepartmentObjective{}).
				Where("department_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": newStatus,
					"is_active":     !deactivate,
				})
		case "division":
			s.db.WithContext(ctx).Model(&performance.DivisionObjective{}).
				Where("division_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": newStatus,
					"is_active":     !deactivate,
				})
		case "office":
			s.db.WithContext(ctx).Model(&performance.OfficeObjective{}).
				Where("office_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": newStatus,
					"is_active":     !deactivate,
				})
		}
	}

	action := "reactivated"
	if deactivate {
		action = "deactivated"
	}

	resp := performance.ResponseVm{}
	resp.Message = fmt.Sprintf("objectives %s successfully", action)
	return resp, nil
}

// =========================================================================
// Evaluation Options
// =========================================================================

func (s *objectiveService) GetEvaluationOptions(ctx context.Context) (interface{}, error) {
	resp := performance.EvaluationOptionResponseVm{}

	var options []performance.EvaluationOption
	err := s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Find(&options).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get evaluation options")
		resp.HasError = true
		resp.Message = "failed to retrieve evaluation options"
		return resp, err
	}

	var vms []performance.EvaluationOptionVm
	for _, opt := range options {
		vm := performance.EvaluationOptionVm{
			EvaluationOptionID: opt.EvaluationOptionID,
			Name:               opt.Name,
			Description:        opt.Description,
			Score:              opt.Score,
		}
		vms = append(vms, vm)
	}

	resp.EvaluationOptions = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) SaveEvaluationOptions(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.EvaluationOptionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveEvaluationOptions")
	}

	if r.EvaluationOptionID != "" {
		// Update
		var existing performance.EvaluationOption
		if err := s.db.WithContext(ctx).
			Where("evaluation_option_id = ?", r.EvaluationOptionID).
			First(&existing).Error; err != nil {
			return nil, fmt.Errorf("evaluation option not found: %w", err)
		}
		existing.Name = r.Name
		existing.Description = r.Description
		existing.Score = r.Score
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("updating evaluation option: %w", err)
		}
	} else {
		// Create
		opt := performance.EvaluationOption{
			Name:        r.Name,
			Description: r.Description,
			Score:       r.Score,
		}
		opt.RecordStatus = enums.StatusActive
		opt.IsActive = true
		if err := s.db.WithContext(ctx).Create(&opt).Error; err != nil {
			return nil, fmt.Errorf("creating evaluation option: %w", err)
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "evaluation option saved successfully"
	return resp, nil
}

// =========================================================================
// Feedback Questionnaires
// =========================================================================

func (s *objectiveService) GetFeedbackQuestionnaires(ctx context.Context) (interface{}, error) {
	resp := performance.FeedbackQuestionaireListResponseVm{}

	var questionnaires []performance.FeedbackQuestionaire
	err := s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("Options").
		Preload("PmsCompetency").
		Find(&questionnaires).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get feedback questionnaires")
		resp.HasError = true
		resp.Message = "failed to retrieve feedback questionnaires"
		return resp, err
	}

	var vms []performance.FeedbackQuestionaireVm
	for _, q := range questionnaires {
		vm := performance.FeedbackQuestionaireVm{
			FeedbackQuestionaireID: q.FeedbackQuestionaireID,
			Question:               q.Question,
			Description:            q.Description,
			PmsCompetencyID:        q.PmsCompetencyID,
		}
		if q.PmsCompetency != nil {
			vm.PmsCompetencyName = q.PmsCompetency.Name
		}

		var opts []performance.FeedbackQuestionaireOptionVm
		for _, o := range q.Options {
			optVm := performance.FeedbackQuestionaireOptionVm{
				FeedbackQuestionaireOptionID: o.FeedbackQuestionaireOptionID,
				OptionStatement:              o.OptionStatement,
				Description:                  o.Description,
				Score:                        o.Score,
				QuestionID:                   o.QuestionID,
			}
			opts = append(opts, optVm)
		}
		vm.Options = opts
		vms = append(vms, vm)
	}

	resp.FeedbackQuestionaires = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) SaveFeedbackQuestionnaires(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.FeedbackQuestionaireVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveFeedbackQuestionnaires")
	}

	if r.FeedbackQuestionaireID != "" {
		// Update existing
		var existing performance.FeedbackQuestionaire
		if err := s.db.WithContext(ctx).
			Where("feedback_questionaire_id = ?", r.FeedbackQuestionaireID).
			First(&existing).Error; err != nil {
			return nil, fmt.Errorf("feedback questionnaire not found: %w", err)
		}
		existing.Question = r.Question
		existing.Description = r.Description
		existing.PmsCompetencyID = r.PmsCompetencyID
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("updating feedback questionnaire: %w", err)
		}
	} else {
		// Create new
		q := performance.FeedbackQuestionaire{
			Question:        r.Question,
			Description:     r.Description,
			PmsCompetencyID: r.PmsCompetencyID,
		}
		q.RecordStatus = enums.StatusActive.String()
		q.IsActive = true
		if err := s.db.WithContext(ctx).Create(&q).Error; err != nil {
			return nil, fmt.Errorf("creating feedback questionnaire: %w", err)
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "feedback questionnaire saved successfully"
	return resp, nil
}

func (s *objectiveService) SaveFeedbackQuestionnaireOptions(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.FeedbackQuestionaireOptionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveFeedbackQuestionnaireOptions")
	}

	if r.FeedbackQuestionaireOptionID != "" {
		// Update existing
		var existing performance.FeedbackQuestionaireOption
		if err := s.db.WithContext(ctx).
			Where("feedback_questionaire_option_id = ?", r.FeedbackQuestionaireOptionID).
			First(&existing).Error; err != nil {
			return nil, fmt.Errorf("feedback questionnaire option not found: %w", err)
		}
		existing.OptionStatement = r.OptionStatement
		existing.Description = r.Description
		existing.Score = r.Score
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("updating feedback questionnaire option: %w", err)
		}
	} else {
		// Create new
		opt := performance.FeedbackQuestionaireOption{
			OptionStatement: r.OptionStatement,
			Description:     r.Description,
			Score:           r.Score,
			QuestionID:      r.QuestionID,
		}
		opt.RecordStatus = enums.StatusActive.String()
		opt.IsActive = true
		if err := s.db.WithContext(ctx).Create(&opt).Error; err != nil {
			return nil, fmt.Errorf("creating feedback questionnaire option: %w", err)
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "feedback questionnaire option saved successfully"
	return resp, nil
}

// =========================================================================
// PMS Competencies
// =========================================================================

func (s *objectiveService) GetPmsCompetencies(ctx context.Context) (interface{}, error) {
	resp := performance.PmsCompetencyListResponseVm{}

	var competencies []performance.PmsCompetency
	err := s.db.WithContext(ctx).
		Where("record_status != ?", enums.StatusCancelled.String()).
		Preload("ObjectiveCategory").
		Find(&competencies).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get PMS competencies")
		resp.HasError = true
		resp.Message = "failed to retrieve PMS competencies"
		return resp, err
	}

	var vms []performance.PmsCompetencyVm
	for _, c := range competencies {
		vm := performance.PmsCompetencyVm{
			PmsCompetencyID:  c.PmsCompetencyID,
			Name:             c.Name,
			Description:      c.Description,
			ObjectCategoryID: c.ObjectCategoryID,
		}
		vms = append(vms, vm)
	}

	resp.PmsCompetencies = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) CreatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.PmsCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreatePmsCompetency")
	}

	// Validate category
	var category performance.ObjectiveCategory
	if err := s.db.WithContext(ctx).
		Where("objective_category_id = ?", r.ObjectCategoryID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("objective category not found: %w", err)
	}

	// Check duplicate
	var existing performance.PmsCompetency
	err := s.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND object_category_id = ? AND record_status != ?",
			r.Name, r.ObjectCategoryID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("PMS competency with this name already exists in this category")
	}

	c := performance.PmsCompetency{
		Name:             r.Name,
		Description:      r.Description,
		ObjectCategoryID: r.ObjectCategoryID,
	}
	c.RecordStatus = enums.StatusActive.String()
	c.IsActive = true

	if err := s.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, fmt.Errorf("creating PMS competency: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = c.PmsCompetencyID
	resp.Message = "PMS competency created successfully"
	return resp, nil
}

func (s *objectiveService) UpdatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.PmsCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdatePmsCompetency")
	}

	var c performance.PmsCompetency
	if err := s.db.WithContext(ctx).
		Where("pms_competency_id = ?", r.PmsCompetencyID).
		First(&c).Error; err != nil {
		return nil, fmt.Errorf("PMS competency not found: %w", err)
	}

	c.Name = r.Name
	c.Description = r.Description
	c.ObjectCategoryID = r.ObjectCategoryID

	if err := s.db.WithContext(ctx).Save(&c).Error; err != nil {
		return nil, fmt.Errorf("updating PMS competency: %w", err)
	}

	resp := performance.ResponseVm{}
	resp.ID = c.PmsCompetencyID
	resp.Message = "PMS competency updated successfully"
	return resp, nil
}

// =========================================================================
// Work Product Definitions
// =========================================================================

func (s *objectiveService) GetObjectiveWorkProductDefinitions(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error) {
	resp := performance.WorkProductDefinitionResponseVm{}

	var defs []performance.WorkProductDefinition
	err := s.db.WithContext(ctx).
		Where("objective_id = ? AND objective_level = ?", objectiveID, objectiveLevel).
		Find(&defs).Error
	if err != nil {
		s.log.Error().Err(err).
			Str("objectiveID", objectiveID).
			Int("objectiveLevel", objectiveLevel).
			Msg("failed to get work product definitions")
		resp.HasError = true
		resp.Message = "failed to retrieve work product definitions"
		return resp, err
	}

	var vms []performance.WorkProductDefinitionVm
	for _, d := range defs {
		vm := performance.WorkProductDefinitionVm{
			WorkProductDefinitionID: d.WorkProductDefinitionID,
			ReferenceNo:             d.ReferenceNo,
			Name:                    d.Name,
			Description:             d.Description,
			Deliverables:            d.Deliverables,
			ObjectiveID:             d.ObjectiveID,
			ObjectiveLevel:          d.ObjectiveLevel,
		}
		vms = append(vms, vm)
	}

	resp.WorkProductDefinitions = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) GetAllWorkProductDefinitions(ctx context.Context) (interface{}, error) {
	resp := performance.WorkProductDefinitionResponseVm{}

	var defs []performance.WorkProductDefinition
	err := s.db.WithContext(ctx).Find(&defs).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get all work product definitions")
		resp.HasError = true
		resp.Message = "failed to retrieve work product definitions"
		return resp, err
	}

	var vms []performance.WorkProductDefinitionVm
	for _, d := range defs {
		vm := performance.WorkProductDefinitionVm{
			WorkProductDefinitionID: d.WorkProductDefinitionID,
			ReferenceNo:             d.ReferenceNo,
			Name:                    d.Name,
			Description:             d.Description,
			Deliverables:            d.Deliverables,
			ObjectiveID:             d.ObjectiveID,
			ObjectiveLevel:          d.ObjectiveLevel,
		}
		vms = append(vms, vm)
	}

	resp.WorkProductDefinitions = vms
	resp.TotalRecords = len(vms)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) GetAllPaginatedWorkProductDefinitions(ctx context.Context, pageIndex, pageSize int, search string) (interface{}, error) {
	resp := performance.PaginatedWorkProductDefinitionResponseVm{}

	query := s.db.WithContext(ctx).Model(&performance.WorkProductDefinition{})

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(reference_no) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	var totalCount int64
	query.Count(&totalCount)

	offset := pageIndex * pageSize
	var defs []performance.WorkProductDefinition
	err := query.Offset(offset).Limit(pageSize).Find(&defs).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get paginated work product definitions")
		resp.HasError = true
		resp.Message = "failed to retrieve work product definitions"
		return resp, err
	}

	var vms []performance.WorkProductDefinitionVm
	for _, d := range defs {
		vm := performance.WorkProductDefinitionVm{
			WorkProductDefinitionID: d.WorkProductDefinitionID,
			ReferenceNo:             d.ReferenceNo,
			Name:                    d.Name,
			Description:             d.Description,
			Deliverables:            d.Deliverables,
			ObjectiveID:             d.ObjectiveID,
			ObjectiveLevel:          d.ObjectiveLevel,
		}
		vms = append(vms, vm)
	}

	resp.WorkProductDefinitions = &performance.PaginatedResult[performance.WorkProductDefinitionVm]{
		Items:      vms,
		TotalCount: int(totalCount),
	}
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *objectiveService) SaveWorkProductDefinitions(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.WorkProductDefinitionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveWorkProductDefinitions")
	}

	if r.WorkProductDefinitionID != "" {
		// Update existing
		var existing performance.WorkProductDefinition
		if err := s.db.WithContext(ctx).
			Where("work_product_definition_id = ?", r.WorkProductDefinitionID).
			First(&existing).Error; err != nil {
			return nil, fmt.Errorf("work product definition not found: %w", err)
		}
		existing.ReferenceNo = r.ReferenceNo
		existing.Name = r.Name
		existing.Description = r.Description
		existing.Deliverables = r.Deliverables
		existing.ObjectiveID = r.ObjectiveID
		existing.ObjectiveLevel = r.ObjectiveLevel
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("updating work product definition: %w", err)
		}
	} else {
		// Create new
		def := performance.WorkProductDefinition{
			ReferenceNo:    r.ReferenceNo,
			Name:           r.Name,
			Description:    r.Description,
			Deliverables:   r.Deliverables,
			ObjectiveID:    r.ObjectiveID,
			ObjectiveLevel: r.ObjectiveLevel,
		}
		def.RecordStatus = enums.StatusActive.String()
		def.IsActive = true
		if err := s.db.WithContext(ctx).Create(&def).Error; err != nil {
			return nil, fmt.Errorf("creating work product definition: %w", err)
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "work product definition saved successfully"
	return resp, nil
}

// =========================================================================
// Approval / Rejection
// =========================================================================

func (s *objectiveService) ApproveRecords(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.ApprovalRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for ApproveRecords")
	}

	for _, id := range r.RecordIDs {
		switch strings.ToLower(r.EntityType) {
		case "enterprise_objective", "enterpriseobjective":
			s.db.WithContext(ctx).Model(&performance.EnterpriseObjective{}).
				Where("enterprise_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": enums.StatusActive.String(),
					"is_active":     true,
					"is_approved":   true,
				})
		case "department_objective", "departmentobjective":
			s.db.WithContext(ctx).Model(&performance.DepartmentObjective{}).
				Where("department_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": enums.StatusActive.String(),
					"is_active":     true,
					"is_approved":   true,
				})
		case "division_objective", "divisionobjective":
			s.db.WithContext(ctx).Model(&performance.DivisionObjective{}).
				Where("division_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": enums.StatusActive.String(),
					"is_active":     true,
					"is_approved":   true,
				})
		case "office_objective", "officeobjective":
			s.db.WithContext(ctx).Model(&performance.OfficeObjective{}).
				Where("office_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": enums.StatusActive.String(),
					"is_active":     true,
					"is_approved":   true,
				})
		case "category_definition", "categorydefinition":
			s.db.WithContext(ctx).Model(&performance.CategoryDefinition{}).
				Where("definition_id = ?", id).
				Updates(map[string]interface{}{
					"record_status": enums.StatusActive.String(),
					"is_active":     true,
					"is_approved":   true,
				})
		default:
			s.log.Warn().Str("entityType", r.EntityType).Msg("unknown entity type for approval")
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "records approved successfully"
	return resp, nil
}

func (s *objectiveService) RejectRecords(ctx context.Context, req interface{}) (interface{}, error) {
	r, ok := req.(*performance.RejectionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for RejectRecords")
	}

	for _, id := range r.RecordIDs {
		switch strings.ToLower(r.EntityType) {
		case "enterprise_objective", "enterpriseobjective":
			s.db.WithContext(ctx).Model(&performance.EnterpriseObjective{}).
				Where("enterprise_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status":    enums.StatusRejected.String(),
					"is_rejected":      true,
					"rejection_reason": r.RejectionReason,
				})
		case "department_objective", "departmentobjective":
			s.db.WithContext(ctx).Model(&performance.DepartmentObjective{}).
				Where("department_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status":    enums.StatusRejected.String(),
					"is_rejected":      true,
					"rejection_reason": r.RejectionReason,
				})
		case "division_objective", "divisionobjective":
			s.db.WithContext(ctx).Model(&performance.DivisionObjective{}).
				Where("division_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status":    enums.StatusRejected.String(),
					"is_rejected":      true,
					"rejection_reason": r.RejectionReason,
				})
		case "office_objective", "officeobjective":
			s.db.WithContext(ctx).Model(&performance.OfficeObjective{}).
				Where("office_objective_id = ?", id).
				Updates(map[string]interface{}{
					"record_status":    enums.StatusRejected.String(),
					"is_rejected":      true,
					"rejection_reason": r.RejectionReason,
				})
		case "category_definition", "categorydefinition":
			s.db.WithContext(ctx).Model(&performance.CategoryDefinition{}).
				Where("definition_id = ?", id).
				Updates(map[string]interface{}{
					"record_status":    enums.StatusRejected.String(),
					"is_rejected":      true,
					"rejection_reason": r.RejectionReason,
				})
		default:
			s.log.Warn().Str("entityType", r.EntityType).Msg("unknown entity type for rejection")
		}
	}

	resp := performance.ResponseVm{}
	resp.Message = "records rejected successfully"
	return resp, nil
}
