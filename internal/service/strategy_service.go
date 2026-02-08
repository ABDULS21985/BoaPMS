package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// strategyService implements PerformanceManagementService.
// It is the Go equivalent of .NET SMDService.cs.
type strategyService struct {
	strategyRepo      *repository.PMSRepository[performance.Strategy]
	themeRepo         *repository.PMSRepository[performance.StrategicTheme]
	entObjRepo        *repository.PMSRepository[performance.EnterpriseObjective]
	divObjRepo        *repository.PMSRepository[performance.DivisionObjective]
	deptObjRepo       *repository.PMSRepository[performance.DepartmentObjective]
	offObjRepo        *repository.PMSRepository[performance.OfficeObjective]
	objCatRepo        *repository.PMSRepository[performance.ObjectiveCategory]
	catDefRepo        *repository.PMSRepository[performance.CategoryDefinition]
	deptRepo          *repository.Repository[organogram.Department]
	divRepo           *repository.Repository[organogram.Division]
	officeRepo        *repository.Repository[organogram.Office]
	evalOptionRepo    *repository.PMSRepository[performance.EvaluationOption]
	feedbackRepo      *repository.PMSRepository[performance.FeedbackQuestionaire]
	feedbackOptRepo   *repository.PMSRepository[performance.FeedbackQuestionaireOption]
	pmsCompRepo       *repository.PMSRepository[performance.PmsCompetency]
	wpDefRepo         *repository.PMSRepository[performance.WorkProductDefinition]
	plannedObjRepo    *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	seqGen            *sequenceGenerator
	db                *gorm.DB
	userCtx           UserContextService
	log               zerolog.Logger
}

func newStrategyService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	userCtx UserContextService,
) *strategyService {
	return &strategyService{
		strategyRepo:   repository.NewPMSRepository[performance.Strategy](repos.GormDB),
		themeRepo:      repository.NewPMSRepository[performance.StrategicTheme](repos.GormDB),
		entObjRepo:     repository.NewPMSRepository[performance.EnterpriseObjective](repos.GormDB),
		divObjRepo:     repository.NewPMSRepository[performance.DivisionObjective](repos.GormDB),
		deptObjRepo:    repository.NewPMSRepository[performance.DepartmentObjective](repos.GormDB),
		offObjRepo:     repository.NewPMSRepository[performance.OfficeObjective](repos.GormDB),
		objCatRepo:     repository.NewPMSRepository[performance.ObjectiveCategory](repos.GormDB),
		catDefRepo:     repository.NewPMSRepository[performance.CategoryDefinition](repos.GormDB),
		deptRepo:       repository.NewRepository[organogram.Department](repos.GormDB),
		divRepo:        repository.NewRepository[organogram.Division](repos.GormDB),
		officeRepo:     repository.NewRepository[organogram.Office](repos.GormDB),
		evalOptionRepo: repository.NewPMSRepository[performance.EvaluationOption](repos.GormDB),
		feedbackRepo:   repository.NewPMSRepository[performance.FeedbackQuestionaire](repos.GormDB),
		feedbackOptRepo: repository.NewPMSRepository[performance.FeedbackQuestionaireOption](repos.GormDB),
		pmsCompRepo:    repository.NewPMSRepository[performance.PmsCompetency](repos.GormDB),
		wpDefRepo:      repository.NewPMSRepository[performance.WorkProductDefinition](repos.GormDB),
		plannedObjRepo: repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](repos.GormDB),
		seqGen:         newSequenceGenerator(repos.GormDB, log),
		db:             repos.GormDB,
		userCtx:        userCtx,
		log:            log,
	}
}

// ---------------------------------------------------------------------------
// Objective level display names (mirrors .NET Humanize())
// ---------------------------------------------------------------------------

func objectiveLevelName(level enums.ObjectiveLevel) string {
	switch level {
	case enums.ObjectiveLevelEnterprise:
		return "Enterprise"
	case enums.ObjectiveLevelDepartment:
		return "Department"
	case enums.ObjectiveLevelDivision:
		return "Division"
	case enums.ObjectiveLevelOffice:
		return "Office"
	default:
		return "Unknown"
	}
}

func parseObjectiveLevel(s string) enums.ObjectiveLevel {
	switch strings.ToLower(s) {
	case "enterprise":
		return enums.ObjectiveLevelEnterprise
	case "department":
		return enums.ObjectiveLevelDepartment
	case "division":
		return enums.ObjectiveLevelDivision
	case "office":
		return enums.ObjectiveLevelOffice
	default:
		return enums.ObjectiveLevelOffice
	}
}

// ---------------------------------------------------------------------------
// Internal data-retrieval helpers (mirrors .NET #region data retrieval)
// ---------------------------------------------------------------------------

func (s *strategyService) getAllEnterpriseObjectives(ctx context.Context, status enums.Status) ([]performance.EnterpriseObjective, error) {
	return s.entObjRepo.GetRecordsWithStatus(ctx, status)
}

func (s *strategyService) getAllDepartmentObjectives(ctx context.Context, status enums.Status) ([]performance.DepartmentObjective, error) {
	return s.deptObjRepo.GetRecordsWithStatus(ctx, status)
}

func (s *strategyService) getAllDivisionObjectives(ctx context.Context, status enums.Status) ([]performance.DivisionObjective, error) {
	return s.divObjRepo.GetRecordsWithStatus(ctx, status)
}

func (s *strategyService) getAllOfficeObjectives(ctx context.Context, status enums.Status) ([]performance.OfficeObjective, error) {
	return s.offObjRepo.GetRecordsWithStatus(ctx, status)
}

func (s *strategyService) getAllObjectiveCategories(ctx context.Context, status enums.Status) ([]performance.ObjectiveCategory, error) {
	return s.objCatRepo.GetRecordsWithStatus(ctx, status)
}

func (s *strategyService) getAllCategoryDefinitions(ctx context.Context, categoryID string) ([]performance.CategoryDefinition, error) {
	catDefs, err := s.catDefRepo.GetRecordsWithStatus(ctx, enums.StatusApprovedAndActive)
	if err != nil {
		return nil, err
	}
	var filtered []performance.CategoryDefinition
	for _, cd := range catDefs {
		if cd.ObjectiveCategoryID == categoryID {
			filtered = append(filtered, cd)
		}
	}
	return filtered, nil
}

func (s *strategyService) getEnterpriseObjective(ctx context.Context, id string) (*performance.EnterpriseObjective, error) {
	return s.entObjRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"Category", "Strategy"},
		"enterprise_objective_id = ?", id,
	)
}

func (s *strategyService) getDepartmentObjective(ctx context.Context, id string) (*performance.DepartmentObjective, error) {
	return s.deptObjRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"Department", "DivisionObjectives"},
		"department_objective_id = ?", id,
	)
}

func (s *strategyService) getDivisionObjective(ctx context.Context, id string) (*performance.DivisionObjective, error) {
	return s.divObjRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"Division", "OfficeObjectives"},
		"division_objective_id = ?", id,
	)
}

func (s *strategyService) getOfficeObjective(ctx context.Context, id string) (*performance.OfficeObjective, error) {
	return s.offObjRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"Office"},
		"office_objective_id = ?", id,
	)
}

func (s *strategyService) getObjectiveCategory(ctx context.Context, id string) (*performance.ObjectiveCategory, error) {
	return s.objCatRepo.FirstOrDefaultWithPreload(ctx, nil, "objective_category_id = ?", id)
}

func (s *strategyService) getStrategy(ctx context.Context, id string) (*performance.Strategy, error) {
	return s.strategyRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"EnterpriseObjectives", "StrategicThemes"},
		"strategy_id = ?", id,
	)
}

func (s *strategyService) getStrategicTheme(ctx context.Context, id string) (*performance.StrategicTheme, error) {
	return s.themeRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"EnterpriseObjectives"},
		"strategic_theme_id = ?", id,
	)
}

func (s *strategyService) getCategoryDefinition(ctx context.Context, id string) (*performance.CategoryDefinition, error) {
	return s.catDefRepo.FirstOrDefaultWithPreload(ctx, nil, "definition_id = ?", id)
}

func (s *strategyService) getPmsCompetency(ctx context.Context, id string) (*performance.PmsCompetency, error) {
	return s.pmsCompRepo.FirstOrDefaultWithPreload(ctx, nil, "pms_competency_id = ?", id)
}

// ---------------------------------------------------------------------------
// Work product objective helpers (mirrors .NET GetWorkProductObjective, etc.)
// ---------------------------------------------------------------------------

// getWorkProductObjectiveByID finds the objective containing work products by its ID.
// Returns the objective's ObjectiveBase fields, the work products list, and the objective
// entity itself as an interface{} for saving.
type objectiveWithWorkProducts struct {
	base         *domain.ObjectiveBase
	workProducts []performance.WorkProductDefinition
	raw          interface{} // the actual entity for saving
}

func (s *strategyService) getWorkProductObjectiveByID(ctx context.Context, objectiveID string, level enums.ObjectiveLevel) (*objectiveWithWorkProducts, error) {
	switch level {
	case enums.ObjectiveLevelOffice:
		obj, err := s.offObjRepo.FirstOrDefaultWithPreload(ctx, nil, "office_objective_id = ?", objectiveID)
		if err != nil || obj == nil {
			return nil, err
		}
		// Load work products via direct query
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", objectiveID, false).Find(&wps)
		return &objectiveWithWorkProducts{base: &obj.ObjectiveBase, workProducts: wps, raw: obj}, nil

	case enums.ObjectiveLevelDivision:
		obj, err := s.divObjRepo.FirstOrDefaultWithPreload(ctx, nil, "division_objective_id = ?", objectiveID)
		if err != nil || obj == nil {
			return nil, err
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", objectiveID, false).Find(&wps)
		return &objectiveWithWorkProducts{base: &obj.ObjectiveBase, workProducts: wps, raw: obj}, nil

	case enums.ObjectiveLevelDepartment:
		obj, err := s.deptObjRepo.FirstOrDefaultWithPreload(ctx, nil, "department_objective_id = ?", objectiveID)
		if err != nil || obj == nil {
			return nil, err
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", objectiveID, false).Find(&wps)
		return &objectiveWithWorkProducts{base: &obj.ObjectiveBase, workProducts: wps, raw: obj}, nil

	case enums.ObjectiveLevelEnterprise:
		obj, err := s.entObjRepo.FirstOrDefaultWithPreload(ctx, nil, "enterprise_objective_id = ?", objectiveID)
		if err != nil || obj == nil {
			return nil, err
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", objectiveID, false).Find(&wps)
		return &objectiveWithWorkProducts{base: &obj.ObjectiveBase, workProducts: wps, raw: obj}, nil

	default:
		return nil, fmt.Errorf("unsupported objective level: %d", int(level))
	}
}

// getWorkProductObjectiveByName finds an objective by name at a given level.
func (s *strategyService) getWorkProductObjectiveByName(ctx context.Context, name string, level enums.ObjectiveLevel) (string, error) {
	lowerName := strings.ToLower(strings.TrimSpace(name))
	switch level {
	case enums.ObjectiveLevelOffice:
		obj, err := s.offObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ?", lowerName)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.OfficeObjectiveID, nil
	case enums.ObjectiveLevelDivision:
		obj, err := s.divObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ?", lowerName)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.DivisionObjectiveID, nil
	case enums.ObjectiveLevelDepartment:
		obj, err := s.deptObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ?", lowerName)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.DepartmentObjectiveID, nil
	case enums.ObjectiveLevelEnterprise:
		obj, err := s.entObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ?", lowerName)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.EnterpriseObjectiveID, nil
	default:
		return "", fmt.Errorf("unsupported objective level: %d", int(level))
	}
}

// getWorkProductObjectiveByNameByReference finds by name + SBU reference.
func (s *strategyService) getWorkProductObjectiveByNameByReference(ctx context.Context, name string, level enums.ObjectiveLevel, sbu string) (string, error) {
	lowerName := strings.ToLower(strings.TrimSpace(name))
	lowerSBU := strings.ToLower(strings.TrimSpace(sbu))

	switch level {
	case enums.ObjectiveLevelOffice:
		var office organogram.Office
		if err := s.db.WithContext(ctx).Where("LOWER(TRIM(office_name)) = ?", lowerSBU).First(&office).Error; err != nil {
			return "", nil
		}
		obj, err := s.offObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ? AND office_id = ?", lowerName, office.OfficeID)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.OfficeObjectiveID, nil

	case enums.ObjectiveLevelDivision:
		var division organogram.Division
		if err := s.db.WithContext(ctx).Where("LOWER(TRIM(division_name)) = ?", lowerSBU).First(&division).Error; err != nil {
			return "", nil
		}
		obj, err := s.divObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ? AND division_id = ?", lowerName, division.DivisionID)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.DivisionObjectiveID, nil

	case enums.ObjectiveLevelDepartment:
		var dept organogram.Department
		if err := s.db.WithContext(ctx).Where("LOWER(TRIM(department_name)) = ?", lowerSBU).First(&dept).Error; err != nil {
			return "", nil
		}
		obj, err := s.deptObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ? AND department_id = ?", lowerName, dept.DepartmentID)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.DepartmentObjectiveID, nil

	case enums.ObjectiveLevelEnterprise:
		obj, err := s.entObjRepo.FirstOrDefault(ctx, "LOWER(TRIM(name)) = ?", lowerName)
		if err != nil || obj == nil {
			return "", err
		}
		return obj.EnterpriseObjectiveID, nil

	default:
		return "", fmt.Errorf("unsupported objective level: %d", int(level))
	}
}

// saveObjectiveWorkProduct saves a work product to the objective, deduplicating by name+level.
func (s *strategyService) saveObjectiveWorkProduct(ctx context.Context, objectiveID string, wp *performance.WorkProductDefinition) error {
	level := parseObjectiveLevel(wp.ObjectiveLevel)
	owp, err := s.getWorkProductObjectiveByID(ctx, objectiveID, level)
	if err != nil {
		return err
	}
	if owp == nil {
		return fmt.Errorf("record not found for %s objective %s", objectiveLevelName(level), wp.ObjectiveID)
	}

	// Check for duplicates by name + objective level
	for _, existing := range owp.workProducts {
		if strings.EqualFold(strings.TrimSpace(existing.Name), strings.TrimSpace(wp.Name)) &&
			existing.ObjectiveLevel == wp.ObjectiveLevel {
			return nil // already exists, skip
		}
	}

	// Save the work product
	if err := s.db.WithContext(ctx).Create(wp).Error; err != nil {
		return err
	}

	// Update parent objective status
	switch obj := owp.raw.(type) {
	case *performance.OfficeObjective:
		obj.Status = enums.StatusApprovedAndActive.String()
		return s.db.WithContext(ctx).Save(obj).Error
	case *performance.DivisionObjective:
		obj.Status = enums.StatusApprovedAndActive.String()
		return s.db.WithContext(ctx).Save(obj).Error
	case *performance.DepartmentObjective:
		obj.Status = enums.StatusApprovedAndActive.String()
		return s.db.WithContext(ctx).Save(obj).Error
	case *performance.EnterpriseObjective:
		obj.Status = enums.StatusApprovedAndActive.String()
		return s.db.WithContext(ctx).Save(obj).Error
	}
	return nil
}

// workProductDefinitionSetup creates or updates a work product definition.
func (s *strategyService) workProductDefinitionSetup(ctx context.Context, req *performance.WorkProductDefinitionVm, isAdd bool) error {
	// Resolve objective ID by name if not provided
	if req.ObjectiveID == "" {
		level := parseObjectiveLevel(req.ObjectiveLevel)
		objID, err := s.getWorkProductObjectiveByNameByReference(ctx, req.ObjectiveName, level, req.SBUName)
		if err != nil {
			return err
		}
		if objID == "" {
			return fmt.Errorf("record not found for objective (%s) in %s objectives", req.ObjectiveName, objectiveLevelName(level))
		}
		req.ObjectiveID = objID
	}

	if isAdd {
		nextNum, err := s.seqGen.getNextNumber(ctx, enums.SeqWorkProductDefinition)
		if err != nil {
			return err
		}
		wpDef := performance.WorkProductDefinition{
			BaseEntity: domain.BaseEntity{
				ID:        int(nextNum),
				IsActive:  true,
				Status:    enums.StatusApprovedAndActive.String(),
				CreatedBy: s.userCtx.GetUserID(ctx),
			},
			WorkProductDefinitionID: fmt.Sprintf("%s-%d", req.ObjectiveID, nextNum),
			ReferenceNo:             req.ReferenceNo,
			Name:                    req.Name,
			Description:             req.Description,
			Deliverables:            req.Deliverables,
			ObjectiveID:             req.ObjectiveID,
			ObjectiveLevel:          req.ObjectiveLevel,
		}
		now := time.Now().UTC()
		wpDef.CreatedAt = &now
		return s.saveObjectiveWorkProduct(ctx, req.ObjectiveID, &wpDef)
	}

	// Update
	wpDef := performance.WorkProductDefinition{
		WorkProductDefinitionID: req.WorkProductDefinitionID,
		ReferenceNo:             req.ReferenceNo,
		Name:                    req.Name,
		Description:             req.Description,
		Deliverables:            req.Deliverables,
		ObjectiveID:             req.ObjectiveID,
		ObjectiveLevel:          req.ObjectiveLevel,
	}
	now := time.Now().UTC()
	wpDef.UpdatedAt = &now
	wpDef.UpdatedBy = s.userCtx.GetUserID(ctx)
	return s.saveObjectiveWorkProduct(ctx, req.ObjectiveID, &wpDef)
}

// feedbackQuestionaireSetup creates or updates a feedback questionnaire.
func (s *strategyService) feedbackQuestionaireSetup(ctx context.Context, req *performance.FeedbackQuestionaireVm, isAdd bool) error {
	if isAdd {
		// Validate duplicate question per competency
		existing, err := s.feedbackRepo.FirstOrDefaultWithPreload(ctx,
			[]string{"PmsCompetency"},
			"LOWER(question) = ? AND pms_competency_id = ? AND record_status != ?",
			strings.ToLower(strings.TrimSpace(req.Question)), req.PmsCompetencyID, enums.StatusCancelled.String(),
		)
		if err != nil {
			return err
		}
		if existing != nil {
			compName := ""
			if existing.PmsCompetency != nil {
				compName = existing.PmsCompetency.Name
			}
			return fmt.Errorf("question: %s already exists for %s", req.Question, compName)
		}

		fbID, err := s.seqGen.GenerateCode(ctx, enums.SeqFeedbackQuestionaire, 15, "", enums.ConCatBefore)
		if err != nil {
			return err
		}

		fb := performance.FeedbackQuestionaire{
			FeedbackQuestionaireID: fbID,
			Question:               req.Question,
			Description:            req.Description,
			PmsCompetencyID:        req.PmsCompetencyID,
			BaseWorkFlow: domain.BaseWorkFlow{
				BaseEntity: domain.BaseEntity{
					Status: enums.StatusApprovedAndActive.String(),
				},
			},
		}

		// Generate options
		if req.Options != nil {
			for _, opt := range req.Options {
				optID, err := s.seqGen.GenerateCode(ctx, enums.SeqFeedbackQuestionaireOption, 8, "", enums.ConCatBefore)
				if err != nil {
					return err
				}
				fb.Options = append(fb.Options, performance.FeedbackQuestionaireOption{
					FeedbackQuestionaireOptionID: optID,
					OptionStatement:              opt.OptionStatement,
					Description:                  opt.Description,
					Score:                        opt.Score,
					QuestionID:                   fbID,
				})
			}
		}

		return s.db.WithContext(ctx).Create(&fb).Error
	}

	// Update
	var fb performance.FeedbackQuestionaire
	if err := s.db.WithContext(ctx).Where("feedback_questionaire_id = ? AND soft_deleted = ?", req.FeedbackQuestionaireID, false).First(&fb).Error; err != nil {
		return fmt.Errorf("record not found")
	}
	fb.Description = req.Description
	fb.Question = req.Question
	fb.Status = enums.StatusApprovedAndActive.String()
	return s.db.WithContext(ctx).Save(&fb).Error
}

// feedbackQuestionaireOptionsSetup creates or updates a feedback questionnaire option.
func (s *strategyService) feedbackQuestionaireOptionsSetup(ctx context.Context, req *performance.FeedbackQuestionaireOptionVm, isAdd bool) error {
	if isAdd {
		// Validate duplicate option per question
		existing, err := s.feedbackOptRepo.FirstOrDefaultWithPreload(ctx,
			[]string{"Question"},
			"LOWER(option_statement) = ? AND question_id = ? AND record_status != ?",
			strings.ToLower(strings.TrimSpace(req.OptionStatement)), req.QuestionID, enums.StatusCancelled.String(),
		)
		if err != nil {
			return err
		}
		if existing != nil {
			qName := ""
			if existing.Question != nil {
				qName = existing.Question.Question
			}
			return fmt.Errorf("option statement: %s already exists for the question: %s", req.OptionStatement, qName)
		}

		optID, err := s.seqGen.GenerateCode(ctx, enums.SeqFeedbackQuestionaire, 15, "", enums.ConCatBefore)
		if err != nil {
			return err
		}

		opt := performance.FeedbackQuestionaireOption{
			FeedbackQuestionaireOptionID: optID,
			OptionStatement:              req.OptionStatement,
			Description:                  req.Description,
			Score:                        req.Score,
			QuestionID:                   req.QuestionID,
			BaseEntity: domain.BaseEntity{
				Status: enums.StatusApprovedAndActive.String(),
			},
		}
		return s.db.WithContext(ctx).Create(&opt).Error
	}

	// Update
	var opt performance.FeedbackQuestionaireOption
	if err := s.db.WithContext(ctx).Where("feedback_questionaire_option_id = ? AND soft_deleted = ?", req.FeedbackQuestionaireOptionID, false).First(&opt).Error; err != nil {
		return fmt.Errorf("record not found")
	}
	opt.Description = req.Description
	opt.OptionStatement = req.OptionStatement
	opt.Score = req.Score
	opt.Status = enums.StatusApprovedAndActive.String()
	return s.db.WithContext(ctx).Save(&opt).Error
}

// evaluationOptionsSetup creates or updates an evaluation option.
func (s *strategyService) evaluationOptionsSetup(ctx context.Context, req *performance.EvaluationOptionVm, isAdd bool) error {
	if isAdd {
		// Validate duplicate by name + evaluation type
		existing, err := s.evalOptionRepo.FirstOrDefault(ctx,
			"LOWER(name) = ? AND evaluation_type = ? AND record_status != ?",
			strings.ToLower(strings.TrimSpace(req.Name)), req.EvaluationType, enums.StatusCancelled.String(),
		)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("evaluation option: %s already exists", req.Name)
		}

		optID, err := s.seqGen.GenerateCode(ctx, enums.SeqEvaluationOption, 15, "", enums.ConCatBefore)
		if err != nil {
			return err
		}

		opt := performance.EvaluationOption{
			EvaluationOptionID: optID,
			Name:               req.Name,
			Description:        req.Description,
			Score:              req.Score,
			EvaluationType:     enums.EvaluationType(req.EvaluationType),
			RecordStatus:       enums.StatusActive,
			BaseWorkFlow: domain.BaseWorkFlow{
				BaseEntity: domain.BaseEntity{
					Status: enums.StatusApprovedAndActive.String(),
				},
			},
		}
		return s.db.WithContext(ctx).Create(&opt).Error
	}

	// Update
	var opt performance.EvaluationOption
	if err := s.db.WithContext(ctx).Where("evaluation_option_id = ? AND soft_deleted = ?", req.EvaluationOptionID, false).First(&opt).Error; err != nil {
		return fmt.Errorf("record not found")
	}
	opt.Description = req.Description
	opt.EvaluationType = enums.EvaluationType(req.EvaluationType)
	opt.Name = req.Name
	opt.Status = enums.StatusApprovedAndActive.String()
	return s.db.WithContext(ctx).Save(&opt).Error
}

// approveOrRejectEntity applies approval/rejection to a workflow entity.
func (s *strategyService) approveOrRejectEntity(ctx context.Context, entity interface{}, approval enums.Approval, reason string) error {
	now := time.Now().UTC()
	userID := s.userCtx.GetUserID(ctx)

	switch e := entity.(type) {
	case *performance.EnterpriseObjective:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.DepartmentObjective:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.DivisionObjective:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.OfficeObjective:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.Strategy:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.StrategicTheme:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	case *performance.ObjectiveCategory:
		if approval == enums.ApprovalApproved {
			e.IsApproved = true
			e.ApprovedBy = userID
			e.DateApproved = &now
			e.RecordStatus = enums.StatusApprovedAndActive.String()
			e.Status = enums.StatusApprovedAndActive.String()
		} else {
			e.IsRejected = true
			e.RejectedBy = userID
			e.DateRejected = &now
			e.RejectionReason = reason
			e.RecordStatus = enums.StatusRejected.String()
			e.Status = enums.StatusRejected.String()
		}
		return s.db.WithContext(ctx).Save(e).Error

	default:
		return fmt.Errorf("unsupported entity type for approval")
	}
}
