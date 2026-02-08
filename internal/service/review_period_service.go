package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// reviewPeriodService manages review periods and objectives planning.
// Mirrors the .NET ReviewPeriodService.
type reviewPeriodService struct {
	reviewPeriodRepo    *repository.PMSRepository[performance.PerformanceReviewPeriod]
	periodObjectiveRepo *repository.PMSRepository[performance.PeriodObjective]
	categoryDefRepo     *repository.PMSRepository[performance.CategoryDefinition]
	objectiveCategoryRepo *repository.PMSRepository[performance.ObjectiveCategory]
	enterpriseObjRepo   *repository.PMSRepository[performance.EnterpriseObjective]
	departmentObjRepo   *repository.PMSRepository[performance.DepartmentObjective]
	divisionObjRepo     *repository.PMSRepository[performance.DivisionObjective]
	officeObjRepo       *repository.PMSRepository[performance.OfficeObjective]
	plannedObjRepo      *repository.PMSRepository[performance.ReviewPeriodIndividualPlannedObjective]
	extensionRepo       *repository.PMSRepository[performance.ReviewPeriodExtension]
	review360Repo       *repository.PMSRepository[performance.ReviewPeriod360Review]
	periodScoreRepo     *repository.PMSRepository[performance.PeriodScore]
	periodObjEvalRepo   *repository.PMSRepository[performance.PeriodObjectiveEvaluation]
	periodObjDeptEvalRepo *repository.PMSRepository[performance.PeriodObjectiveDepartmentEvaluation]
	strategyRepo        *repository.PMSRepository[performance.Strategy]
	db                  *gorm.DB
	cfg                 *config.Config
	log                 zerolog.Logger
}

// newReviewPeriodService creates a ReviewPeriodService with all required repositories.
func newReviewPeriodService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) ReviewPeriodService {
	return &reviewPeriodService{
		reviewPeriodRepo:      repository.NewPMSRepository[performance.PerformanceReviewPeriod](repos.GormDB),
		periodObjectiveRepo:   repository.NewPMSRepository[performance.PeriodObjective](repos.GormDB),
		categoryDefRepo:       repository.NewPMSRepository[performance.CategoryDefinition](repos.GormDB),
		objectiveCategoryRepo: repository.NewPMSRepository[performance.ObjectiveCategory](repos.GormDB),
		enterpriseObjRepo:     repository.NewPMSRepository[performance.EnterpriseObjective](repos.GormDB),
		departmentObjRepo:     repository.NewPMSRepository[performance.DepartmentObjective](repos.GormDB),
		divisionObjRepo:       repository.NewPMSRepository[performance.DivisionObjective](repos.GormDB),
		officeObjRepo:         repository.NewPMSRepository[performance.OfficeObjective](repos.GormDB),
		plannedObjRepo:        repository.NewPMSRepository[performance.ReviewPeriodIndividualPlannedObjective](repos.GormDB),
		extensionRepo:         repository.NewPMSRepository[performance.ReviewPeriodExtension](repos.GormDB),
		review360Repo:         repository.NewPMSRepository[performance.ReviewPeriod360Review](repos.GormDB),
		periodScoreRepo:       repository.NewPMSRepository[performance.PeriodScore](repos.GormDB),
		periodObjEvalRepo:     repository.NewPMSRepository[performance.PeriodObjectiveEvaluation](repos.GormDB),
		periodObjDeptEvalRepo: repository.NewPMSRepository[performance.PeriodObjectiveDepartmentEvaluation](repos.GormDB),
		strategyRepo:          repository.NewPMSRepository[performance.Strategy](repos.GormDB),
		db:                    repos.GormDB,
		cfg:                   cfg,
		log:                   log.With().Str("service", "review_period").Logger(),
	}
}

// ---------------------------------------------------------------------------
// Helper: getReviewPeriod retrieves a review period by ID (internal use).
// ---------------------------------------------------------------------------

func (s *reviewPeriodService) getReviewPeriod(ctx context.Context, reviewPeriodID string) (*performance.ReviewPeriodResponseVm, error) {
	response := &performance.ReviewPeriodResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", reviewPeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	response.PerformanceReviewPeriod = rp
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ---------------------------------------------------------------------------
// Helper: getExtendedActiveReviewPeriod checks extensions for the period.
// ---------------------------------------------------------------------------

func (s *reviewPeriodService) getExtendedActiveReviewPeriod(ctx context.Context, rp *performance.PerformanceReviewPeriod, staffID string) (*performance.ReviewPeriodResponseVm, error) {
	response := &performance.ReviewPeriodResponseVm{
		PerformanceReviewPeriod: rp,
	}
	response.HasError = true

	now := time.Now().UTC()

	// Check Bankwide extension
	bankwideExt, err := s.extensionRepo.FirstOrDefault(ctx,
		"review_period_id = ? AND target_type = ? AND record_status = ? AND start_date <= ? AND end_date >= ?",
		rp.PeriodID, enums.ExtensionTargetBankwide, enums.StatusActive.String(), now, now)
	if err != nil {
		return response, err
	}
	if bankwideExt != nil {
		response.HasExtension = true
		response.ExtensionEndDate = bankwideExt.EndDate
		response.HasError = false
		response.Message = "Operation completed"
		return response, nil
	}

	if staffID == "" {
		response.HasError = false
		response.Message = "Operation completed"
		return response, nil
	}

	// Check Staff extension
	staffExt, err := s.extensionRepo.FirstOrDefault(ctx,
		"review_period_id = ? AND target_type = ? AND target_reference = ? AND record_status = ? AND start_date <= ? AND end_date >= ?",
		rp.PeriodID, enums.ExtensionTargetStaff, staffID, enums.StatusActive.String(), now, now)
	if err != nil {
		return response, err
	}
	if staffExt != nil {
		response.HasExtension = true
		response.ExtensionEndDate = staffExt.EndDate
		response.HasError = false
		response.Message = "Operation completed"
		return response, nil
	}

	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ---------------------------------------------------------------------------
// Helper: validateRangeValue validates range value for the period range.
// ---------------------------------------------------------------------------

func (s *reviewPeriodService) validateRangeValue(rangeType enums.ReviewPeriodRange, rangeValue int) error {
	switch rangeType {
	case enums.ReviewPeriodRangeQuarterly:
		if rangeValue < 1 || rangeValue > 4 {
			return fmt.Errorf("quarterly range value must be between 1 and 4")
		}
	case enums.ReviewPeriodRangeBiAnnual:
		if rangeValue < 1 || rangeValue > 2 {
			return fmt.Errorf("bi-annual range value must be between 1 and 2")
		}
	case enums.ReviewPeriodRangeAnnual:
		if rangeValue != 1 {
			return fmt.Errorf("annual range value must be 1")
		}
	default:
		return fmt.Errorf("invalid review period range")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper: getStartOrEndDate calculates start/end dates from range.
// ---------------------------------------------------------------------------

func (s *reviewPeriodService) getStartDate(year int, rangeType enums.ReviewPeriodRange, rangeValue int) time.Time {
	switch rangeType {
	case enums.ReviewPeriodRangeQuarterly:
		month := time.Month((rangeValue-1)*3 + 1)
		return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	case enums.ReviewPeriodRangeBiAnnual:
		month := time.Month((rangeValue-1)*6 + 1)
		return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	case enums.ReviewPeriodRangeAnnual:
		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func (s *reviewPeriodService) getEndDate(year int, rangeType enums.ReviewPeriodRange, rangeValue int) time.Time {
	switch rangeType {
	case enums.ReviewPeriodRangeQuarterly:
		month := time.Month(rangeValue * 3)
		return time.Date(year, month+1, 0, 23, 59, 59, 0, time.UTC)
	case enums.ReviewPeriodRangeBiAnnual:
		month := time.Month(rangeValue * 6)
		return time.Date(year, month+1, 0, 23, 59, 59, 0, time.UTC)
	case enums.ReviewPeriodRangeAnnual:
		return time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	}
	return time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
}

// ---------------------------------------------------------------------------
// Helper: generateCode generates a unique code for entities.
// ---------------------------------------------------------------------------

func (s *reviewPeriodService) generateCode(ctx context.Context, seqType enums.SequenceNumberTypes, length int) (string, error) {
	var count int64
	tableName := ""
	prefix := ""
	switch seqType {
	case enums.SeqReviewPeriod:
		tableName = "pms.performance_review_periods"
		prefix = "RP"
	case enums.SeqObjectivePeriodMapping:
		tableName = "pms.period_objectives"
		prefix = "PO"
	case enums.SeqCategoryDefinitions:
		tableName = "pms.category_definitions"
		prefix = "CD"
	case enums.SeqReviewPeriodExtension:
		tableName = "pms.review_period_extensions"
		prefix = "RE"
	case enums.SeqReviewPeriod360:
		tableName = "pms.review_period_360_reviews"
		prefix = "R3"
	case enums.SeqObjective:
		tableName = "pms.review_period_individual_planned_objectives"
		prefix = "IP"
	case enums.SeqObjectiveOutcomePeriodMapping:
		tableName = "pms.period_objective_evaluations"
		prefix = "PE"
	case enums.SeqDeptObjectiveOutcomePeriodMapping:
		tableName = "pms.period_objective_department_evaluations"
		prefix = "DE"
	default:
		tableName = "pms.performance_review_periods"
		prefix = "XX"
	}

	err := s.db.WithContext(ctx).Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count).Error
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%s%0*d", prefix, length-len(prefix), count+1)
	return code, nil
}

// ===========================================================================
// Review Period Lifecycle Methods
// ===========================================================================

// SaveDraftReviewPeriod creates a new review period in Draft status.
func (s *reviewPeriodService) SaveDraftReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CreateNewReviewPeriodVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CreateNewReviewPeriodVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	// Validate name uniqueness (case-insensitive, same year)
	existing, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"LOWER(name) = ? AND year = ? AND record_status != ?",
		strings.ToLower(vm.Name), vm.Year, enums.StatusCancelled.String())
	if err != nil {
		s.log.Error().Err(err).Msg("failed to check name uniqueness")
		return response, err
	}
	if existing != nil {
		response.Message = "A review period with this name already exists for the selected year"
		return response, nil
	}

	// Validate short name uniqueness
	if vm.ShortName != "" {
		existingShort, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
			"LOWER(short_name) = ? AND year = ? AND record_status != ?",
			strings.ToLower(vm.ShortName), vm.Year, enums.StatusCancelled.String())
		if err != nil {
			s.log.Error().Err(err).Msg("failed to check short name uniqueness")
			return response, err
		}
		if existingShort != nil {
			response.Message = "A review period with this short name already exists for the selected year"
			return response, nil
		}
	}

	// Validate min/max objectives
	if vm.MinNoOfObjectives < 0 {
		response.Message = "Minimum number of objectives must be at least 0"
		return response, nil
	}
	if vm.MinNoOfObjectives > vm.MaxNoOfObjectives {
		response.Message = "Minimum number of objectives must not exceed maximum"
		return response, nil
	}

	// Validate strategy
	strategy, err := s.strategyRepo.GetByStringID(ctx, "strategy_id", vm.StrategyID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve strategy")
		return response, err
	}
	if strategy == nil {
		response.Message = "Strategy record not found"
		return response, nil
	}
	if !strategy.IsApproved {
		response.Message = "The selected strategy has not been approved"
		return response, nil
	}

	// Validate range
	rangeType := enums.ReviewPeriodRange(vm.Range)
	if err := s.validateRangeValue(rangeType, vm.RangeValue); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	// Validate period uniqueness for range+year
	existingPeriod, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"year = ? AND \"range\" = ? AND range_value = ? AND record_status != ?",
		vm.Year, vm.Range, vm.RangeValue, enums.StatusCancelled.String())
	if err != nil {
		s.log.Error().Err(err).Msg("failed to check period uniqueness")
		return response, err
	}
	if existingPeriod != nil {
		response.Message = "A review period already exists for the selected range and year"
		return response, nil
	}

	// Calculate dates
	startDate := s.getStartDate(vm.Year, rangeType, vm.RangeValue)
	endDate := s.getEndDate(vm.Year, rangeType, vm.RangeValue)

	// Generate code
	periodID, err := s.generateCode(ctx, enums.SeqReviewPeriod, 15)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate period code")
		return response, err
	}

	entity := &performance.PerformanceReviewPeriod{
		PeriodID:               periodID,
		Year:                   vm.Year,
		Range:                  rangeType,
		RangeValue:             vm.RangeValue,
		Name:                   vm.Name,
		Description:            vm.Description,
		ShortName:              vm.ShortName,
		StartDate:              startDate,
		EndDate:                endDate,
		AllowObjectivePlanning: false,
		MaxPoints:              vm.MaxPoints,
		MinNoOfObjectives:      vm.MinNoOfObjectives,
		MaxNoOfObjectives:      vm.MaxNoOfObjectives,
		StrategyID:             vm.StrategyID,
	}
	entity.RecordStatus = enums.StatusDraft.String()

	if err := s.reviewPeriodRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to save draft review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = periodID
	s.log.Info().Str("periodID", periodID).Msg("review period draft saved")
	return response, nil
}

// AddReviewPeriod creates a new review period in PendingApproval status.
func (s *reviewPeriodService) AddReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CreateNewReviewPeriodVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CreateNewReviewPeriodVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	// Validate name uniqueness
	existing, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"LOWER(name) = ? AND year = ? AND record_status != ?",
		strings.ToLower(vm.Name), vm.Year, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "A review period with this name already exists for the selected year"
		return response, nil
	}

	// Validate short name uniqueness
	if vm.ShortName != "" {
		existingShort, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
			"LOWER(short_name) = ? AND year = ? AND record_status != ?",
			strings.ToLower(vm.ShortName), vm.Year, enums.StatusCancelled.String())
		if err != nil {
			return response, err
		}
		if existingShort != nil {
			response.Message = "A review period with this short name already exists for the selected year"
			return response, nil
		}
	}

	// Validate objectives
	if vm.MinNoOfObjectives < 0 {
		response.Message = "Minimum number of objectives must be at least 0"
		return response, nil
	}
	if vm.MinNoOfObjectives > vm.MaxNoOfObjectives {
		response.Message = "Minimum number of objectives must not exceed maximum"
		return response, nil
	}

	// Validate strategy
	strategy, err := s.strategyRepo.GetByStringID(ctx, "strategy_id", vm.StrategyID)
	if err != nil {
		return response, err
	}
	if strategy == nil {
		response.Message = "Strategy record not found"
		return response, nil
	}
	if !strategy.IsApproved {
		response.Message = "The selected strategy has not been approved"
		return response, nil
	}

	// Validate range
	rangeType := enums.ReviewPeriodRange(vm.Range)
	if err := s.validateRangeValue(rangeType, vm.RangeValue); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	// Period uniqueness
	existingPeriod, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"year = ? AND \"range\" = ? AND range_value = ? AND record_status != ?",
		vm.Year, vm.Range, vm.RangeValue, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existingPeriod != nil {
		response.Message = "A review period already exists for the selected range and year"
		return response, nil
	}

	startDate := s.getStartDate(vm.Year, rangeType, vm.RangeValue)
	endDate := s.getEndDate(vm.Year, rangeType, vm.RangeValue)

	periodID, err := s.generateCode(ctx, enums.SeqReviewPeriod, 15)
	if err != nil {
		return response, err
	}

	entity := &performance.PerformanceReviewPeriod{
		PeriodID:               periodID,
		Year:                   vm.Year,
		Range:                  rangeType,
		RangeValue:             vm.RangeValue,
		Name:                   vm.Name,
		Description:            vm.Description,
		ShortName:              vm.ShortName,
		StartDate:              startDate,
		EndDate:                endDate,
		AllowObjectivePlanning: false,
		MaxPoints:              vm.MaxPoints,
		MinNoOfObjectives:      vm.MinNoOfObjectives,
		MaxNoOfObjectives:      vm.MaxNoOfObjectives,
		StrategyID:             vm.StrategyID,
	}
	entity.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.reviewPeriodRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to add review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = periodID
	s.log.Info().Str("periodID", periodID).Msg("review period added")
	return response, nil
}

// SubmitDraftReviewPeriod transitions a Draft review period to PendingApproval.
func (s *reviewPeriodService) SubmitDraftReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusDraft.String() {
		response.Message = "Review period cannot be submitted; it is not in Draft status"
		return response, nil
	}

	// Validate objectives
	if vm.MinNoOfObjectives < 0 {
		response.Message = "Minimum number of objectives must be at least 0"
		return response, nil
	}
	if vm.MinNoOfObjectives > vm.MaxNoOfObjectives {
		response.Message = "Minimum number of objectives must not exceed maximum"
		return response, nil
	}

	rp.Name = vm.Name
	rp.Description = vm.Description
	rp.ShortName = vm.ShortName
	rp.MaxPoints = vm.MaxPoints
	rp.MinNoOfObjectives = vm.MinNoOfObjectives
	rp.MaxNoOfObjectives = vm.MaxNoOfObjectives
	rp.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to submit draft review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// UpdateReviewPeriod updates a Draft or Returned review period.
func (s *reviewPeriodService) UpdateReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusDraft.String() && rp.RecordStatus != enums.StatusReturned.String() {
		response.Message = "Review period cannot be modified"
		return response, nil
	}

	if vm.MinNoOfObjectives < 0 {
		response.Message = "Minimum number of objectives must be at least 0"
		return response, nil
	}
	if vm.MinNoOfObjectives > vm.MaxNoOfObjectives {
		response.Message = "Minimum number of objectives must not exceed maximum"
		return response, nil
	}

	rp.Name = vm.Name
	rp.Description = vm.Description
	rp.ShortName = vm.ShortName
	rp.MaxPoints = vm.MaxPoints
	rp.MinNoOfObjectives = vm.MinNoOfObjectives
	rp.MaxNoOfObjectives = vm.MaxNoOfObjectives

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to update review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// ApproveReviewPeriod approves a PendingApproval review period, making it Active.
func (s *reviewPeriodService) ApproveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Review period cannot be approved; it is not pending approval"
		return response, nil
	}

	// Check no existing active period for the same year
	activePeriod, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"year = ? AND record_status = ? AND period_id != ?",
		rp.Year, enums.StatusActive.String(), rp.PeriodID)
	if err != nil {
		return response, err
	}
	if activePeriod != nil {
		response.Message = "An active review period already exists for the selected year"
		return response, nil
	}

	// Validate period objectives exist
	periodObjs, err := s.periodObjectiveRepo.Where(ctx,
		"review_period_id = ? AND record_status != ?",
		rp.PeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if len(periodObjs) == 0 {
		response.Message = "Review period must have at least one period objective before approval"
		return response, nil
	}

	// Validate category definitions exist
	catDefs, err := s.categoryDefRepo.Where(ctx,
		"review_period_id = ? AND record_status != ?",
		rp.PeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if len(catDefs) == 0 {
		response.Message = "Review period must have at least one category definition before approval"
		return response, nil
	}

	now := time.Now().UTC()

	rp.RecordStatus = enums.StatusActive.String()
	rp.IsActive = true
	rp.IsApproved = true
	rp.ApprovedBy = vm.ApprovedBy
	rp.DateApproved = &now
	rp.IsRejected = false
	rp.RejectedBy = ""
	rp.RejectionReason = ""
	rp.AllowObjectivePlanning = true

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to approve review period")
		return response, err
	}

	// Cascade approval to period objectives
	for i := range periodObjs {
		periodObjs[i].RecordStatus = enums.StatusActive.String()
		periodObjs[i].IsActive = true
		periodObjs[i].IsApproved = true
		periodObjs[i].ApprovedBy = vm.ApprovedBy
		periodObjs[i].DateApproved = &now
		if err := s.periodObjectiveRepo.UpdateAndSave(ctx, &periodObjs[i]); err != nil {
			s.log.Error().Err(err).Str("periodObjID", periodObjs[i].PeriodObjectiveID).Msg("failed to approve period objective")
		}
	}

	// Cascade approval to category definitions
	for i := range catDefs {
		catDefs[i].RecordStatus = enums.StatusActive.String()
		catDefs[i].IsActive = true
		catDefs[i].IsApproved = true
		catDefs[i].ApprovedBy = vm.ApprovedBy
		catDefs[i].DateApproved = &now
		if err := s.categoryDefRepo.UpdateAndSave(ctx, &catDefs[i]); err != nil {
			s.log.Error().Err(err).Str("definitionID", catDefs[i].DefinitionID).Msg("failed to approve category definition")
		}
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	s.log.Info().Str("periodID", rp.PeriodID).Msg("review period approved")
	return response, nil
}

// RejectReviewPeriod rejects a PendingApproval review period.
func (s *reviewPeriodService) RejectReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Review period cannot be rejected; it is not pending approval"
		return response, nil
	}

	now := time.Now().UTC()

	rp.RecordStatus = enums.StatusRejected.String()
	rp.IsActive = false
	rp.IsRejected = true
	rp.RejectedBy = vm.RejectedBy
	rp.DateRejected = &now
	rp.RejectionReason = vm.RejectionReason
	rp.IsApproved = false
	rp.ApprovedBy = ""

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to reject review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// ReturnReviewPeriod returns a PendingApproval review period for revisions.
func (s *reviewPeriodService) ReturnReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Review period cannot be returned; it is not pending approval"
		return response, nil
	}

	now := time.Now().UTC()

	rp.RecordStatus = enums.StatusReturned.String()
	rp.IsActive = false
	rp.IsRejected = true
	rp.RejectedBy = vm.RejectedBy
	rp.DateRejected = &now
	rp.RejectionReason = vm.RejectionReason
	rp.IsApproved = false
	rp.ApprovedBy = ""

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to return review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// ReSubmitReviewPeriod re-submits a Returned review period for approval.
func (s *reviewPeriodService) ReSubmitReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusReturned.String() {
		response.Message = "Review period cannot be re-submitted; it is not in Returned status"
		return response, nil
	}

	if vm.MinNoOfObjectives < 0 {
		response.Message = "Minimum number of objectives must be at least 0"
		return response, nil
	}
	if vm.MinNoOfObjectives > vm.MaxNoOfObjectives {
		response.Message = "Minimum number of objectives must not exceed maximum"
		return response, nil
	}

	rp.Name = vm.Name
	rp.Description = vm.Description
	rp.ShortName = vm.ShortName
	rp.MaxPoints = vm.MaxPoints
	rp.MinNoOfObjectives = vm.MinNoOfObjectives
	rp.MaxNoOfObjectives = vm.MaxNoOfObjectives
	rp.RecordStatus = enums.StatusPendingApproval.String()
	rp.IsRejected = false
	rp.IsApproved = false

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to resubmit review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// CancelReviewPeriod cancels a Draft review period.
func (s *reviewPeriodService) CancelReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusDraft.String() {
		response.Message = "Review period cannot be cancelled; it is not in Draft status"
		return response, nil
	}

	rp.RecordStatus = enums.StatusCancelled.String()
	rp.IsActive = false

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to cancel review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// CloseReviewPeriod closes an Active review period.
func (s *reviewPeriodService) CloseReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	if rp.RecordStatus != enums.StatusActive.String() && rp.RecordStatus != enums.StatusApprovedAndActive.String() {
		response.Message = "Review period cannot be closed"
		return response, nil
	}

	rp.RecordStatus = enums.StatusClosed.String()
	rp.IsActive = false
	rp.AllowObjectivePlanning = false
	rp.AllowWorkProductPlanning = false
	rp.AllowWorkProductEvaluation = false

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Msg("failed to close review period")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	s.log.Info().Str("periodID", rp.PeriodID).Msg("review period closed")
	return response, nil
}

// ===========================================================================
// Toggle Methods
// ===========================================================================

func (s *reviewPeriodService) toggleReviewPeriodFlag(ctx context.Context, req interface{}, flagName string, value bool) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.PeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = rpResp.Message
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod

	// Check period is Active or has active extension
	isActive := rp.RecordStatus == enums.StatusActive.String() || rp.RecordStatus == enums.StatusApprovedAndActive.String()
	if !isActive {
		extResp, err := s.getExtendedActiveReviewPeriod(ctx, rp, "")
		if err != nil {
			return response, err
		}
		if !extResp.HasExtension {
			response.Message = "Review period is not active"
			return response, nil
		}
	}

	switch flagName {
	case "AllowObjectivePlanning":
		rp.AllowObjectivePlanning = value
	case "AllowWorkProductPlanning":
		rp.AllowWorkProductPlanning = value
	case "AllowWorkProductEvaluation":
		rp.AllowWorkProductEvaluation = value
	}

	if err := s.reviewPeriodRepo.UpdateAndSave(ctx, rp); err != nil {
		s.log.Error().Err(err).Str("flag", flagName).Msg("failed to toggle review period flag")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = rp.PeriodID
	return response, nil
}

// EnableObjectivePlanning enables objective planning for a review period.
func (s *reviewPeriodService) EnableObjectivePlanning(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowObjectivePlanning", true)
}

// DisableObjectivePlanning disables objective planning for a review period.
func (s *reviewPeriodService) DisableObjectivePlanning(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowObjectivePlanning", false)
}

// EnableWorkProductPlanning enables work product planning for a review period.
func (s *reviewPeriodService) EnableWorkProductPlanning(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowWorkProductPlanning", true)
}

// DisableWorkProductPlanning disables work product planning for a review period.
func (s *reviewPeriodService) DisableWorkProductPlanning(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowWorkProductPlanning", false)
}

// EnableWorkProductEvaluation enables work product evaluation for a review period.
func (s *reviewPeriodService) EnableWorkProductEvaluation(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowWorkProductEvaluation", true)
}

// DisableWorkProductEvaluation disables work product evaluation for a review period.
func (s *reviewPeriodService) DisableWorkProductEvaluation(ctx context.Context, req interface{}) (interface{}, error) {
	return s.toggleReviewPeriodFlag(ctx, req, "AllowWorkProductEvaluation", false)
}

// ===========================================================================
// Retrieval Methods
// ===========================================================================

// GetActiveReviewPeriod retrieves the currently active review period.
// Checks extensions if the period is Closed.
func (s *reviewPeriodService) GetActiveReviewPeriod(ctx context.Context) (interface{}, error) {
	response := &performance.ReviewPeriodResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
		"record_status = ?", enums.StatusActive.String())
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve active review period")
		return response, err
	}

	if rp == nil {
		// Check for a closed period with active extension
		closedRP, err := s.reviewPeriodRepo.FirstOrDefault(ctx,
			"record_status = ?", enums.StatusClosed.String())
		if err != nil {
			return response, err
		}
		if closedRP != nil {
			extResp, err := s.getExtendedActiveReviewPeriod(ctx, closedRP, "")
			if err != nil {
				return response, err
			}
			if extResp.HasExtension {
				extResp.HasError = false
				extResp.Message = "Operation completed"
				return extResp, nil
			}
		}
		response.Message = "No active review period found"
		return response, nil
	}

	response.PerformanceReviewPeriod = rp
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// GetStaffActiveReviewPeriod retrieves the active review period for a specific staff.
func (s *reviewPeriodService) GetStaffActiveReviewPeriod(ctx context.Context, staffID string) (interface{}, error) {
	response := &performance.ReviewPeriodResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	// First check if staff has any planned objectives
	plannedObjs, err := s.plannedObjRepo.Where(ctx,
		"staff_id = ? AND record_status != ?",
		staffID, enums.StatusCancelled.String())
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve staff planned objectives")
		return response, err
	}

	if len(plannedObjs) > 0 {
		// Get the review period from the first planned objective
		reviewPeriodID := plannedObjs[0].ReviewPeriodID
		rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", reviewPeriodID)
		if err != nil {
			return response, err
		}
		if rp != nil && (rp.RecordStatus == enums.StatusActive.String() || rp.RecordStatus == enums.StatusApprovedAndActive.String()) {
			response.PerformanceReviewPeriod = rp
			response.HasError = false
			response.Message = "Operation completed"
			return response, nil
		}

		// Check extensions for this staff
		if rp != nil && rp.RecordStatus == enums.StatusClosed.String() {
			extResp, err := s.getExtendedActiveReviewPeriod(ctx, rp, staffID)
			if err != nil {
				return response, err
			}
			if extResp.HasExtension {
				extResp.HasError = false
				extResp.Message = "Operation completed"
				return extResp, nil
			}
		}
	}

	// Fallback: get the general active review period
	return s.GetActiveReviewPeriod(ctx)
}

// GetReviewPeriodDetails retrieves complete details for a review period.
func (s *reviewPeriodService) GetReviewPeriodDetails(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.PerformanceReviewPeriodResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rp, err := s.reviewPeriodRepo.GetByStringID(ctx, "period_id", reviewPeriodID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve review period details")
		return response, err
	}
	if rp == nil {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rpVm := &performance.PerformanceReviewPeriodVm{
		PeriodID:                   rp.PeriodID,
		Year:                       rp.Year,
		Range:                      int(rp.Range),
		RangeValue:                 rp.RangeValue,
		Name:                       rp.Name,
		Description:                rp.Description,
		ShortName:                  rp.ShortName,
		StartDate:                  rp.StartDate,
		EndDate:                    rp.EndDate,
		AllowObjectivePlanning:     rp.AllowObjectivePlanning,
		AllowWorkProductPlanning:   rp.AllowWorkProductPlanning,
		AllowWorkProductEvaluation: rp.AllowWorkProductEvaluation,
		MaxPoints:                  rp.MaxPoints,
		MinNoOfObjectives:          rp.MinNoOfObjectives,
		MaxNoOfObjectives:          rp.MaxNoOfObjectives,
		StrategyID:                 rp.StrategyID,
	}
	rpVm.RecordStatus = rp.RecordStatus
	rpVm.IsActive = rp.IsActive
	rpVm.IsApproved = rp.IsApproved
	rpVm.CreatedBy = rp.CreatedBy
	rpVm.CreatedAt = rp.CreatedAt

	response.PerformanceReviewPeriod = rpVm
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// Period Objectives
// ===========================================================================

// SaveDraftReviewPeriodObjective saves a period objective in Draft status.
func (s *reviewPeriodService) SaveDraftReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.SaveDraftPeriodObjectiveVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.SaveDraftPeriodObjectiveVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod
	if rp.RecordStatus != enums.StatusDraft.String() && rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Period objectives can only be added to Draft or PendingApproval review periods"
		return response, nil
	}

	var lastID string
	for _, objID := range vm.ObjectiveIDs {
		// Check duplicate
		existing, err := s.periodObjectiveRepo.FirstOrDefault(ctx,
			"review_period_id = ? AND objective_id = ? AND record_status != ?",
			vm.ReviewPeriodID, objID, enums.StatusCancelled.String())
		if err != nil {
			return response, err
		}
		if existing != nil {
			continue // Skip duplicates
		}

		periodObjID, err := s.generateCode(ctx, enums.SeqObjectivePeriodMapping, 15)
		if err != nil {
			return response, err
		}

		po := &performance.PeriodObjective{
			PeriodObjectiveID: periodObjID,
			ObjectiveID:       objID,
			ReviewPeriodID:    vm.ReviewPeriodID,
		}
		po.RecordStatus = rp.RecordStatus

		if err := s.periodObjectiveRepo.InsertAndSave(ctx, po); err != nil {
			s.log.Error().Err(err).Msg("failed to save draft period objective")
			return response, err
		}
		lastID = periodObjID
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = lastID
	return response, nil
}

// AddReviewPeriodObjective adds period objectives in PendingApproval status.
func (s *reviewPeriodService) AddReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.AddPeriodObjectiveVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.AddPeriodObjectiveVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod
	if rp.RecordStatus != enums.StatusDraft.String() && rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Period objectives can only be added to Draft or PendingApproval review periods"
		return response, nil
	}

	var lastID string
	for _, objID := range vm.ObjectiveIDs {
		existing, err := s.periodObjectiveRepo.FirstOrDefault(ctx,
			"review_period_id = ? AND objective_id = ? AND record_status != ?",
			vm.ReviewPeriodID, objID, enums.StatusCancelled.String())
		if err != nil {
			return response, err
		}
		if existing != nil {
			continue
		}

		periodObjID, err := s.generateCode(ctx, enums.SeqObjectivePeriodMapping, 15)
		if err != nil {
			return response, err
		}

		po := &performance.PeriodObjective{
			PeriodObjectiveID: periodObjID,
			ObjectiveID:       objID,
			ReviewPeriodID:    vm.ReviewPeriodID,
		}
		po.RecordStatus = rp.RecordStatus

		if err := s.periodObjectiveRepo.InsertAndSave(ctx, po); err != nil {
			s.log.Error().Err(err).Msg("failed to add period objective")
			return response, err
		}
		lastID = periodObjID
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = lastID
	return response, nil
}

// SubmitDraftReviewPeriodObjective submits a draft period objective (CommitDraft).
func (s *reviewPeriodService) SubmitDraftReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.PeriodObjectiveRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.PeriodObjectiveRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.periodObjectiveRepo.GetByStringID(ctx, "period_objective_id", vm.PeriodObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Period objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusDraft.String() {
		response.Message = "Period objective cannot be submitted; it is not in Draft status"
		return response, nil
	}

	po.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.periodObjectiveRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to submit draft period objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PeriodObjectiveID
	return response, nil
}

// CancelReviewPeriodObjective cancels a period objective.
func (s *reviewPeriodService) CancelReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.PeriodObjectiveRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.PeriodObjectiveRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.periodObjectiveRepo.GetByStringID(ctx, "period_objective_id", vm.PeriodObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Period objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusDraft.String() && po.RecordStatus != enums.StatusReturned.String() {
		response.Message = "Period objective cannot be cancelled"
		return response, nil
	}

	po.RecordStatus = enums.StatusCancelled.String()
	po.IsActive = false

	if err := s.periodObjectiveRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to cancel period objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PeriodObjectiveID
	return response, nil
}

// GetReviewPeriodObjectives retrieves all objectives for a review period.
func (s *reviewPeriodService) GetReviewPeriodObjectives(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.ReviewPeriodObjectivesResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, reviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	// Get period objectives with preloaded objective data
	var periodObjs []performance.PeriodObjective
	err = s.db.WithContext(ctx).
		Where("review_period_id = ? AND soft_deleted = ?", reviewPeriodID, false).
		Preload("Objective").
		Preload("Objective.Category").
		Preload("Objective.Category.CategoryDefinitions").
		Find(&periodObjs).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve period objectives")
		return response, err
	}

	objectives := make([]performance.EnterpriseObjectiveData, 0, len(periodObjs))
	for _, po := range periodObjs {
		if po.Objective == nil {
			continue
		}
		obj := po.Objective
		data := performance.EnterpriseObjectiveData{
			EnterpriseObjectiveID:          obj.EnterpriseObjectiveID,
			PeriodObjectiveID:              po.PeriodObjectiveID,
			Name:                           obj.Name,
			Description:                    obj.Description,
			Kpi:                            obj.Kpi,
			Target:                         obj.Target,
			EnterpriseObjectivesCategoryID: obj.EnterpriseObjectivesCategoryID,
			StrategyID:                     obj.StrategyID,
		}
		data.RecordStatus = po.RecordStatus
		data.IsActive = obj.IsActive
		data.IsApproved = obj.IsApproved
		data.IsRejected = obj.IsRejected
		data.CreatedBy = obj.CreatedBy
		data.CreatedAt = obj.CreatedAt
		data.ApprovedBy = obj.ApprovedBy
		data.DateApproved = obj.DateApproved
		data.RejectedBy = obj.RejectedBy
		data.RejectionReason = obj.RejectionReason
		data.DateRejected = obj.DateRejected
		data.UpdatedBy = obj.UpdatedBy
		data.UpdatedAt = obj.UpdatedAt

		// Check for evaluation
		var eval performance.PeriodObjectiveEvaluation
		evalErr := s.db.WithContext(ctx).
			Joins("JOIN pms.period_objectives ON pms.period_objectives.period_objective_id = pms.period_objective_evaluations.period_objective_id").
			Where("pms.period_objectives.objective_id = ? AND pms.period_objectives.review_period_id = ?", obj.EnterpriseObjectiveID, reviewPeriodID).
			First(&eval).Error
		if evalErr == nil {
			data.HasEvaluation = true
			data.OutcomeScore = math.Round(eval.OutcomeScore*100) / 100
			data.TotalOutcomeScore = math.Round(eval.TotalOutcomeScore*100) / 100
			data.Evaluator = eval.CreatedBy
			if eval.CreatedAt != nil {
				data.EvaluationDate = *eval.CreatedAt
			}
		}

		objectives = append(objectives, data)
	}

	response.TotalRecords = len(objectives)
	response.Objectives = objectives
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// Category Definitions
// ===========================================================================

func (s *reviewPeriodService) validateCategoryDefinition(ctx context.Context, request *performance.CategoryDefinitionRequestVm, rp *performance.PerformanceReviewPeriod) (string, error) {
	// Validate weight
	if request.Weight <= 0 || request.Weight > 100 {
		return "Weight must be between 0 - 100%", nil
	}

	if request.MaxPoints > int(rp.MaxPoints) {
		return fmt.Sprintf("Selected maximum points of %d cannot be more than review period maximum points of %.0f", request.MaxPoints, rp.MaxPoints), nil
	}

	if request.MaxNoObjectives > rp.MaxNoOfObjectives {
		return fmt.Sprintf("Selected maximum objectives of %d cannot be more than review period maximum objectives of %d", request.MaxNoObjectives, rp.MaxNoOfObjectives), nil
	}

	return "", nil
}

// SaveDraftCategoryDefinition saves a category definition in Draft status.
func (s *reviewPeriodService) SaveDraftCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CategoryDefinitionRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}
	rp := rpResp.PerformanceReviewPeriod

	// Validate category exists
	cat, err := s.objectiveCategoryRepo.GetByStringID(ctx, "objective_category_id", vm.ObjectiveCategoryID)
	if err != nil {
		return response, err
	}
	if cat == nil {
		response.Message = "Objective Category record not found"
		return response, nil
	}

	// Check uniqueness
	existing, err := s.categoryDefRepo.FirstOrDefault(ctx,
		"review_period_id = ? AND objective_category_id = ? AND grade_group_id = ? AND record_status != ?",
		vm.ReviewPeriodID, vm.ObjectiveCategoryID, vm.GradeGroupID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "Objective Category Definition already exists, provide another definition"
		return response, nil
	}

	// Validate
	if msg, _ := s.validateCategoryDefinition(ctx, vm, rp); msg != "" {
		response.Message = msg
		return response, nil
	}

	maxPoint := rp.MaxPoints * (vm.Weight / 100)
	maxNoWorkProduct := vm.MaxNoWorkProduct

	defID, err := s.generateCode(ctx, enums.SeqCategoryDefinitions, 15)
	if err != nil {
		return response, err
	}

	entity := &performance.CategoryDefinition{
		DefinitionID:        defID,
		ObjectiveCategoryID: vm.ObjectiveCategoryID,
		ReviewPeriodID:      vm.ReviewPeriodID,
		Weight:              vm.Weight,
		MaxNoObjectives:     vm.MaxNoObjectives,
		MaxNoWorkProduct:    maxNoWorkProduct,
		MaxPoints:           maxPoint,
		IsCompulsory:        vm.IsCompulsory,
		EnforceWorkProductLimit: vm.EnforceWorkProductLimit,
		Description:         vm.Description,
		GradeGroupID:        vm.GradeGroupID,
	}
	entity.RecordStatus = rp.RecordStatus

	if err := s.categoryDefRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to save draft category definition")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = defID
	return response, nil
}

// AddCategoryDefinition adds a category definition.
func (s *reviewPeriodService) AddCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CategoryDefinitionRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}
	rp := rpResp.PerformanceReviewPeriod

	if rp.RecordStatus != enums.StatusDraft.String() && rp.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Objective Category Definition cannot be added to the review period"
		return response, nil
	}

	// Check uniqueness
	existing, err := s.categoryDefRepo.FirstOrDefault(ctx,
		"review_period_id = ? AND objective_category_id = ? AND grade_group_id = ? AND record_status != ?",
		vm.ReviewPeriodID, vm.ObjectiveCategoryID, vm.GradeGroupID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "Objective Category Definition already exists, provide another definition"
		return response, nil
	}

	if msg, _ := s.validateCategoryDefinition(ctx, vm, rp); msg != "" {
		response.Message = msg
		return response, nil
	}

	maxPoint := rp.MaxPoints * (vm.Weight / 100)

	defID, err := s.generateCode(ctx, enums.SeqCategoryDefinitions, 15)
	if err != nil {
		return response, err
	}

	entity := &performance.CategoryDefinition{
		DefinitionID:        defID,
		ObjectiveCategoryID: vm.ObjectiveCategoryID,
		ReviewPeriodID:      vm.ReviewPeriodID,
		Weight:              vm.Weight,
		MaxNoObjectives:     vm.MaxNoObjectives,
		MaxNoWorkProduct:    vm.MaxNoWorkProduct,
		MaxPoints:           maxPoint,
		IsCompulsory:        vm.IsCompulsory,
		EnforceWorkProductLimit: vm.EnforceWorkProductLimit,
		Description:         vm.Description,
		GradeGroupID:        vm.GradeGroupID,
	}
	entity.RecordStatus = rp.RecordStatus

	if err := s.categoryDefRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to add category definition")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = defID
	return response, nil
}

// SubmitDraftCategoryDefinition submits a Draft category definition (CommitDraft).
func (s *reviewPeriodService) SubmitDraftCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CategoryDefinitionRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	def, err := s.categoryDefRepo.GetByStringID(ctx, "definition_id", vm.DefinitionID)
	if err != nil {
		return response, err
	}
	if def == nil {
		response.Message = "Objective Category Definition record is not found"
		return response, nil
	}

	if def.RecordStatus != enums.StatusDraft.String() {
		response.Message = "Objective Category Definition cannot be modified"
		return response, nil
	}

	rpResp, err := s.getReviewPeriod(ctx, def.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	rp := rpResp.PerformanceReviewPeriod

	if vm.Weight < 0 || vm.Weight > 100 {
		response.Message = "Weight must be between 0 - 100%"
		return response, nil
	}

	maxPoint := rp.MaxPoints * (vm.Weight / 100)

	def.Weight = vm.Weight
	def.MaxPoints = maxPoint
	def.MaxNoObjectives = vm.MaxNoObjectives
	def.IsCompulsory = vm.IsCompulsory
	def.GradeGroupID = vm.GradeGroupID
	def.Description = vm.Description
	def.RecordStatus = enums.StatusPendingApproval.String()
	def.EnforceWorkProductLimit = vm.EnforceWorkProductLimit

	if err := s.categoryDefRepo.UpdateAndSave(ctx, def); err != nil {
		s.log.Error().Err(err).Msg("failed to submit draft category definition")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = def.DefinitionID
	return response, nil
}

// ApproveCategoryDefinition approves a PendingApproval category definition.
func (s *reviewPeriodService) ApproveCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CategoryDefinitionRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	def, err := s.categoryDefRepo.GetByStringID(ctx, "definition_id", vm.DefinitionID)
	if err != nil {
		return response, err
	}
	if def == nil {
		response.Message = "Objective Category Definition record is not found"
		return response, nil
	}

	if def.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Objective Category Definition cannot be modified"
		return response, nil
	}

	now := time.Now().UTC()

	def.RecordStatus = enums.StatusActive.String()
	def.IsActive = true
	def.IsApproved = true
	def.ApprovedBy = vm.ApprovedBy
	def.DateApproved = &now
	def.IsRejected = false
	def.RejectedBy = ""
	def.RejectionReason = ""

	if err := s.categoryDefRepo.UpdateAndSave(ctx, def); err != nil {
		s.log.Error().Err(err).Msg("failed to approve category definition")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = def.DefinitionID
	return response, nil
}

// RejectCategoryDefinition rejects a PendingApproval category definition.
func (s *reviewPeriodService) RejectCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CategoryDefinitionRequestVm")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	def, err := s.categoryDefRepo.GetByStringID(ctx, "definition_id", vm.DefinitionID)
	if err != nil {
		return response, err
	}
	if def == nil {
		response.Message = "Objective Category Definition record is not found"
		return response, nil
	}

	if def.RecordStatus != enums.StatusPendingApproval.String() {
		response.Message = "Objective Category Definition cannot be rejected"
		return response, nil
	}

	now := time.Now().UTC()

	def.IsActive = false
	def.IsRejected = true
	def.RejectedBy = vm.RejectedBy
	def.DateRejected = &now
	def.RecordStatus = enums.StatusRejected.String()
	def.RejectionReason = vm.RejectionReason
	def.IsApproved = false
	def.ApprovedBy = ""

	if err := s.categoryDefRepo.UpdateAndSave(ctx, def); err != nil {
		s.log.Error().Err(err).Msg("failed to reject category definition")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = def.DefinitionID
	return response, nil
}

// ===========================================================================
// Extensions
// ===========================================================================

// AddReviewPeriodExtension adds a review period extension.
func (s *reviewPeriodService) AddReviewPeriodExtension(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodExtensionRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodExtensionRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod
	if rp.RecordStatus != enums.StatusClosed.String() {
		response.Message = "Review Period extension can only be done if its closed"
		return response, nil
	}

	// Validate dates
	if !vm.StartDate.IsZero() && !vm.EndDate.IsZero() {
		if vm.StartDate.After(vm.EndDate) || vm.StartDate.Equal(vm.EndDate) {
			response.Message = "The start date must not exceed the end date"
			return response, nil
		}
		if vm.StartDate.Before(rp.EndDate) {
			response.Message = "The extension start date must exceed the review period end date"
			return response, nil
		}
	}

	// Check for existing active extension
	targetType := enums.ReviewPeriodExtensionTargetType(vm.TargetType)
	existing, err := s.extensionRepo.FirstOrDefault(ctx,
		"review_period_id = ? AND target_reference = ? AND target_type = ? AND record_status != ? AND record_status != ?",
		vm.ReviewPeriodID, vm.TargetReference, vm.TargetType, enums.StatusCancelled.String(), enums.StatusClosed.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "Record already exists"
		return response, nil
	}

	extID, err := s.generateCode(ctx, enums.SeqReviewPeriodExtension, 10)
	if err != nil {
		return response, err
	}

	targetRef := vm.TargetReference
	if targetType == enums.ExtensionTargetBankwide {
		targetRef = "All Staff"
	}

	entity := &performance.ReviewPeriodExtension{
		ReviewPeriodExtensionID: extID,
		ReviewPeriodID:          vm.ReviewPeriodID,
		TargetType:              targetType,
		TargetReference:         targetRef,
		Description:             vm.Description,
		StartDate:               vm.StartDate,
		EndDate:                 vm.EndDate,
	}
	entity.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.extensionRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to add review period extension")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = extID
	s.log.Info().Str("extensionID", extID).Msg("review period extension added")
	return response, nil
}

// GetReviewPeriodExtensions retrieves all extensions for a review period.
func (s *reviewPeriodService) GetReviewPeriodExtensions(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.ReviewPeriodExtensionListResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var extensions []performance.ReviewPeriodExtension
	err := s.db.WithContext(ctx).
		Where("review_period_id = ? AND soft_deleted = ?", reviewPeriodID, false).
		Preload("ReviewPeriod").
		Find(&extensions).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve review period extensions")
		return response, err
	}

	data := make([]performance.ReviewPeriodExtensionData, 0, len(extensions))
	for _, ext := range extensions {
		d := performance.ReviewPeriodExtensionData{
			ReviewPeriodExtensionID: ext.ReviewPeriodExtensionID,
			ReviewPeriodID:          ext.ReviewPeriodID,
			TargetType:              int(ext.TargetType),
			TargetReference:         ext.TargetReference,
			Description:             ext.Description,
			StartDate:               ext.StartDate,
			EndDate:                 ext.EndDate,
			RecordStatus:            enums.Status(0), // Will use string comparison
		}
		if ext.ReviewPeriod != nil {
			d.ReviewPeriod = ext.ReviewPeriod.Name
		}
		data = append(data, d)
	}

	response.TotalRecords = len(data)
	response.ReviewPeriodExtensions = data
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// 360 Reviews
// ===========================================================================

// AddReviewPeriod360Review adds a 360 review configuration for a review period.
func (s *reviewPeriodService) AddReviewPeriod360Review(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CreateReviewPeriod360ReviewRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.CreateReviewPeriod360ReviewRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	// Check duplicate
	existing, err := s.review360Repo.FirstOrDefault(ctx,
		"review_period_id = ? AND target_type = ? AND target_reference = ? AND record_status != ?",
		vm.ReviewPeriodID, vm.TargetType, vm.TargetReference, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "360 Review record already exists"
		return response, nil
	}

	reviewID, err := s.generateCode(ctx, enums.SeqReviewPeriod360, 10)
	if err != nil {
		return response, err
	}

	entity := &performance.ReviewPeriod360Review{
		ReviewPeriod360ReviewID: reviewID,
		ReviewPeriodID:          vm.ReviewPeriodID,
		TargetType:              enums.ReviewPeriod360TargetType(vm.TargetType),
		TargetReference:         vm.TargetReference,
	}
	entity.RecordStatus = enums.StatusActive.String()
	entity.IsActive = true

	if err := s.review360Repo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to add 360 review")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = reviewID
	return response, nil
}

// GetReviewPeriod360Reviews retrieves all 360 reviews for a review period.
func (s *reviewPeriodService) GetReviewPeriod360Reviews(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	reviews, err := s.review360Repo.WhereWithPreload(ctx,
		[]string{"ReviewPeriod"},
		"review_period_id = ? AND record_status != ?",
		reviewPeriodID, enums.StatusCancelled.String())
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve 360 reviews")
		return nil, err
	}

	return map[string]interface{}{
		"hasError":     false,
		"message":      "Operation completed",
		"totalRecords": len(reviews),
		"reviews":      reviews,
	}, nil
}

// ===========================================================================
// Individual Planned Objectives
// ===========================================================================

// SaveDraftIndividualPlannedObjective saves a planned objective in Draft status.
func (s *reviewPeriodService) SaveDraftIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod
	if !rp.AllowObjectivePlanning {
		response.Message = "Objective planning is not enabled for this review period"
		return response, nil
	}

	// Check for existing planned objective with the same objectiveID + staffID
	existing, err := s.plannedObjRepo.FirstOrDefault(ctx,
		"objective_id = ? AND staff_id = ? AND review_period_id = ? AND record_status != ?",
		vm.ObjectiveID, vm.StaffID, vm.ReviewPeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "Individual planned objective already exists"
		return response, nil
	}

	// Default to Office level
	objectiveLevel := enums.ObjectiveLevelOffice

	plannedObjID, err := s.generateCode(ctx, enums.SeqObjective, 15)
	if err != nil {
		return response, err
	}

	entity := &performance.ReviewPeriodIndividualPlannedObjective{
		PlannedObjectiveID: plannedObjID,
		ObjectiveID:        vm.ObjectiveID,
		StaffID:            vm.StaffID,
		ObjectiveLevel:     objectiveLevel,
		StaffJobRole:       "",
		ReviewPeriodID:     vm.ReviewPeriodID,
	}
	entity.RecordStatus = enums.StatusDraft.String()

	if err := s.plannedObjRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to save draft planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = plannedObjID
	return response, nil
}

// AddIndividualPlannedObjective adds a planned objective in PendingApproval status.
func (s *reviewPeriodService) AddIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	rpResp, err := s.getReviewPeriod(ctx, vm.ReviewPeriodID)
	if err != nil {
		return response, err
	}
	if rpResp.HasError {
		response.Message = "Review Period record not found"
		return response, nil
	}

	rp := rpResp.PerformanceReviewPeriod
	if !rp.AllowObjectivePlanning {
		response.Message = "Objective planning is not enabled for this review period"
		return response, nil
	}

	// Check for existing
	existing, err := s.plannedObjRepo.FirstOrDefault(ctx,
		"objective_id = ? AND staff_id = ? AND review_period_id = ? AND record_status != ?",
		vm.ObjectiveID, vm.StaffID, vm.ReviewPeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if existing != nil {
		response.Message = "Individual planned objective already exists"
		return response, nil
	}

	objectiveLevel := enums.ObjectiveLevelOffice

	plannedObjID, err := s.generateCode(ctx, enums.SeqObjective, 15)
	if err != nil {
		return response, err
	}

	entity := &performance.ReviewPeriodIndividualPlannedObjective{
		PlannedObjectiveID: plannedObjID,
		ObjectiveID:        vm.ObjectiveID,
		StaffID:            vm.StaffID,
		ObjectiveLevel:     objectiveLevel,
		StaffJobRole:       "",
		ReviewPeriodID:     vm.ReviewPeriodID,
	}
	entity.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.plannedObjRepo.InsertAndSave(ctx, entity); err != nil {
		s.log.Error().Err(err).Msg("failed to add planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = plannedObjID
	s.log.Info().Str("plannedObjID", plannedObjID).Msg("individual planned objective added")
	return response, nil
}

// SubmitDraftIndividualPlannedObjective submits a draft planned objective.
func (s *reviewPeriodService) SubmitDraftIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.GetByStringID(ctx, "planned_objective_id", vm.PlannedObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Planned objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusDraft.String() {
		response.Message = "Planned objective cannot be submitted; it is not in Draft status"
		return response, nil
	}

	po.RecordStatus = enums.StatusPendingApproval.String()

	if err := s.plannedObjRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to submit draft planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PlannedObjectiveID
	return response, nil
}

// ApproveIndividualPlannedObjective approves a planned objective.
func (s *reviewPeriodService) ApproveIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.GetByStringID(ctx, "planned_objective_id", vm.PlannedObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Planned objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusPendingApproval.String() &&
		po.RecordStatus != enums.StatusPendingAcceptance.String() &&
		po.RecordStatus != enums.StatusSuspensionPendingApproval.String() {
		response.Message = "Planned objective cannot be approved"
		return response, nil
	}

	now := time.Now().UTC()

	// If SuspensionPendingApproval, transition to Paused
	if po.RecordStatus == enums.StatusSuspensionPendingApproval.String() {
		po.RecordStatus = enums.StatusPaused.String()
	} else {
		po.RecordStatus = enums.StatusActive.String()
	}

	po.IsActive = true
	po.IsApproved = true
	po.ApprovedBy = vm.ApprovedBy
	po.DateApproved = &now
	po.IsRejected = false
	po.RejectedBy = ""
	po.RejectionReason = ""

	if err := s.plannedObjRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to approve planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PlannedObjectiveID
	return response, nil
}

// RejectIndividualPlannedObjective rejects a planned objective.
func (s *reviewPeriodService) RejectIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.GetByStringID(ctx, "planned_objective_id", vm.PlannedObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Planned objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusPendingApproval.String() &&
		po.RecordStatus != enums.StatusPendingAcceptance.String() {
		response.Message = "Planned objective cannot be rejected"
		return response, nil
	}

	now := time.Now().UTC()

	po.RecordStatus = enums.StatusRejected.String()
	po.IsActive = false
	po.IsRejected = true
	po.RejectedBy = vm.RejectedBy
	po.DateRejected = &now
	po.RejectionReason = vm.RejectionReason
	po.IsApproved = false
	po.ApprovedBy = ""

	if err := s.plannedObjRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to reject planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PlannedObjectiveID
	return response, nil
}

// ReturnIndividualPlannedObjective returns a planned objective for revisions.
// Auto-triggers grievance if max return count is reached.
func (s *reviewPeriodService) ReturnIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.GetByStringID(ctx, "planned_objective_id", vm.PlannedObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Planned objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusPendingApproval.String() &&
		po.RecordStatus != enums.StatusPendingAcceptance.String() {
		response.Message = "Planned objective cannot be returned"
		return response, nil
	}

	now := time.Now().UTC()

	po.NoReturned++
	po.RecordStatus = enums.StatusReturned.String()
	po.IsActive = false
	po.IsRejected = true
	po.RejectedBy = vm.RejectedBy
	po.DateRejected = &now
	po.RejectionReason = vm.RejectionReason
	po.IsApproved = false
	po.ApprovedBy = ""
	po.Remark = vm.Remark

	if err := s.plannedObjRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to return planned objective")
		return response, err
	}

	// NOTE: In .NET, auto-grievance is triggered when NoReturned >= MAX_RETURN_NO.
	// This is handled by the grievance service integration which is out of scope
	// for this conversion but the NoReturned counter is properly incremented above.

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PlannedObjectiveID
	return response, nil
}

// CancelIndividualPlannedObjective cancels a planned objective.
func (s *reviewPeriodService) CancelIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ReviewPeriodIndividualPlannedObjectiveRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.ReviewPeriodIndividualPlannedObjectiveRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.GetByStringID(ctx, "planned_objective_id", vm.PlannedObjectiveID)
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Planned objective record not found"
		return response, nil
	}

	if po.RecordStatus != enums.StatusDraft.String() && po.RecordStatus != enums.StatusReturned.String() {
		response.Message = "Planned objective cannot be cancelled"
		return response, nil
	}

	po.RecordStatus = enums.StatusCancelled.String()
	po.IsActive = false

	if err := s.plannedObjRepo.UpdateAndSave(ctx, po); err != nil {
		s.log.Error().Err(err).Msg("failed to cancel planned objective")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = po.PlannedObjectiveID
	return response, nil
}

// GetStaffIndividualPlannedObjectives retrieves planned objectives for a staff member.
func (s *reviewPeriodService) GetStaffIndividualPlannedObjectives(ctx context.Context, staffID, reviewPeriodID string) (interface{}, error) {
	response := &performance.PlannedOperationalObjectivesResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var plannedObjs []performance.ReviewPeriodIndividualPlannedObjective
	err := s.db.WithContext(ctx).
		Where("staff_id = ? AND review_period_id = ? AND soft_deleted = ?", staffID, reviewPeriodID, false).
		Preload("ReviewPeriod").
		Find(&plannedObjs).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve staff planned objectives")
		return response, err
	}

	plannedData := make([]performance.PlannedObjectiveData, 0, len(plannedObjs))
	for _, po := range plannedObjs {
		data := performance.PlannedObjectiveData{
			PlannedObjectiveID: po.PlannedObjectiveID,
			ReviewPeriodID:     po.ReviewPeriodID,
			ObjectiveLevel:     fmt.Sprintf("%d", po.ObjectiveLevel),
			ObjectiveID:        po.ObjectiveID,
			RecordStatus:       po.RecordStatus,
			StaffID:            po.StaffID,
			CreatedBy:          po.CreatedBy,
			CreatedDate:        po.CreatedAt,
			IsApproved:         po.IsApproved,
			IsRejected:         po.IsRejected,
			IsActive:           po.IsActive,
			Approver:           po.ApprovedBy,
			Comment:            po.Remark,
		}
		if po.ReviewPeriod != nil {
			data.ReviewPeriod = po.ReviewPeriod.Name
			data.Year = po.ReviewPeriod.Year
		}

		// Resolve objective name based on level
		switch po.ObjectiveLevel {
		case enums.ObjectiveLevelOffice:
			var offObj performance.OfficeObjective
			if err := s.db.WithContext(ctx).Where("office_objective_id = ?", po.ObjectiveID).First(&offObj).Error; err == nil {
				data.Objective = offObj.Name
				data.Kpi = offObj.Kpi
				data.Target = offObj.Target
				data.Description = offObj.Description
			}
		case enums.ObjectiveLevelDivision:
			var divObj performance.DivisionObjective
			if err := s.db.WithContext(ctx).Where("division_objective_id = ?", po.ObjectiveID).First(&divObj).Error; err == nil {
				data.Objective = divObj.Name
				data.Kpi = divObj.Kpi
				data.Target = divObj.Target
				data.Description = divObj.Description
			}
		case enums.ObjectiveLevelDepartment:
			var deptObj performance.DepartmentObjective
			if err := s.db.WithContext(ctx).Where("department_objective_id = ?", po.ObjectiveID).First(&deptObj).Error; err == nil {
				data.Objective = deptObj.Name
				data.Kpi = deptObj.Kpi
				data.Target = deptObj.Target
				data.Description = deptObj.Description
			}
		}

		plannedData = append(plannedData, data)
	}

	response.TotalRecords = len(plannedData)
	response.PlannedObjectives = plannedData
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// Period Objective Evaluations
// ===========================================================================

// CreatePeriodObjectiveEvaluation creates or updates a period objective evaluation.
func (s *reviewPeriodService) CreatePeriodObjectiveEvaluation(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.AddPeriodObjectiveEvaluationRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.AddPeriodObjectiveEvaluationRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	// Find the period objective
	po, err := s.periodObjectiveRepo.FirstOrDefault(ctx,
		"objective_id = ? AND review_period_id = ? AND record_status != ?",
		vm.EnterpriseObjectiveID, vm.ReviewPeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Period objective record not found"
		return response, nil
	}

	// Check if evaluation already exists
	existingEval, err := s.periodObjEvalRepo.FirstOrDefault(ctx,
		"period_objective_id = ?", po.PeriodObjectiveID)
	if err != nil {
		return response, err
	}

	if existingEval != nil {
		// Update existing
		existingEval.TotalOutcomeScore = vm.TotalOutcomeScore
		existingEval.OutcomeScore = vm.OutcomeScore
		if err := s.periodObjEvalRepo.UpdateAndSave(ctx, existingEval); err != nil {
			return response, err
		}
		response.HasError = false
		response.Message = "Operation completed"
		response.ID = existingEval.PeriodObjectiveEvaluationID
		return response, nil
	}

	// Create new
	evalID, err := s.generateCode(ctx, enums.SeqObjectiveOutcomePeriodMapping, 15)
	if err != nil {
		return response, err
	}

	eval := &performance.PeriodObjectiveEvaluation{
		PeriodObjectiveEvaluationID: evalID,
		TotalOutcomeScore:           vm.TotalOutcomeScore,
		OutcomeScore:                vm.OutcomeScore,
		PeriodObjectiveID:           po.PeriodObjectiveID,
	}
	eval.RecordStatus = enums.StatusActive.String()
	eval.IsActive = true

	if err := s.periodObjEvalRepo.InsertAndSave(ctx, eval); err != nil {
		s.log.Error().Err(err).Msg("failed to create period objective evaluation")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = evalID
	return response, nil
}

// CreatePeriodObjectiveDepartmentEvaluation creates or updates a department evaluation.
func (s *reviewPeriodService) CreatePeriodObjectiveDepartmentEvaluation(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.AddPeriodObjectiveDepartmentEvaluationRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.AddPeriodObjectiveDepartmentEvaluationRequestModel")
	}

	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	// Find the period objective
	po, err := s.periodObjectiveRepo.FirstOrDefault(ctx,
		"objective_id = ? AND review_period_id = ? AND record_status != ?",
		vm.EnterpriseObjectiveID, vm.ReviewPeriodID, enums.StatusCancelled.String())
	if err != nil {
		return response, err
	}
	if po == nil {
		response.Message = "Period objective record not found"
		return response, nil
	}

	// Check if department evaluation already exists
	existingEval, err := s.periodObjDeptEvalRepo.FirstOrDefault(ctx,
		"period_objective_id = ? AND department_id = ?",
		po.PeriodObjectiveID, vm.DepartmentID)
	if err != nil {
		return response, err
	}

	if existingEval != nil {
		existingEval.OverallOutcomeScored = vm.OverallOutcomeScored
		existingEval.AllocatedOutcome = vm.AllocatedOutcome
		existingEval.OutcomeScore = vm.OutcomeScore
		if err := s.periodObjDeptEvalRepo.UpdateAndSave(ctx, existingEval); err != nil {
			return response, err
		}
		response.HasError = false
		response.Message = "Operation completed"
		response.ID = existingEval.PeriodObjectiveDepartmentEvaluationID
		return response, nil
	}

	evalID, err := s.generateCode(ctx, enums.SeqDeptObjectiveOutcomePeriodMapping, 15)
	if err != nil {
		return response, err
	}

	eval := &performance.PeriodObjectiveDepartmentEvaluation{
		PeriodObjectiveDepartmentEvaluationID: evalID,
		OverallOutcomeScored:                  vm.OverallOutcomeScored,
		AllocatedOutcome:                      vm.AllocatedOutcome,
		OutcomeScore:                          vm.OutcomeScore,
		DepartmentID:                          vm.DepartmentID,
		PeriodObjectiveID:                     po.PeriodObjectiveID,
	}
	eval.RecordStatus = enums.StatusActive.String()
	eval.IsActive = true

	if err := s.periodObjDeptEvalRepo.InsertAndSave(ctx, eval); err != nil {
		s.log.Error().Err(err).Msg("failed to create department evaluation")
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	response.ID = evalID
	return response, nil
}

// GetPeriodObjectiveEvaluations retrieves evaluations for a review period.
func (s *reviewPeriodService) GetPeriodObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.PeriodObjectiveEvaluationListResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var evals []performance.PeriodObjectiveEvaluation
	err := s.db.WithContext(ctx).
		Joins("JOIN pms.period_objectives ON pms.period_objectives.period_objective_id = pms.period_objective_evaluations.period_objective_id").
		Where("pms.period_objectives.review_period_id = ? AND pms.period_objective_evaluations.soft_deleted = ?", reviewPeriodID, false).
		Preload("PeriodObjective").
		Preload("PeriodObjective.Objective").
		Preload("PeriodObjective.ReviewPeriod").
		Find(&evals).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve period objective evaluations")
		return response, err
	}

	data := make([]performance.PeriodObjectiveEvaluationData, 0, len(evals))
	for _, e := range evals {
		d := performance.PeriodObjectiveEvaluationData{
			PeriodObjectiveEvaluationID: e.PeriodObjectiveEvaluationID,
			TotalOutcomeScore:           e.TotalOutcomeScore,
			OutcomeScore:                e.OutcomeScore,
			PeriodObjectiveID:           e.PeriodObjectiveID,
		}
		d.RecordStatus = e.RecordStatus
		d.IsActive = e.IsActive
		d.CreatedBy = e.CreatedBy
		d.CreatedAt = e.CreatedAt

		if e.PeriodObjective != nil {
			if e.PeriodObjective.Objective != nil {
				d.EnterpriseObjectiveID = e.PeriodObjective.Objective.EnterpriseObjectiveID
				d.EnterpriseObjective = e.PeriodObjective.Objective.Name
			}
			if e.PeriodObjective.ReviewPeriod != nil {
				d.ReviewPeriodID = e.PeriodObjective.ReviewPeriod.PeriodID
				d.ReviewPeriod = e.PeriodObjective.ReviewPeriod.Name
			}
		}
		data = append(data, d)
	}

	response.TotalRecords = len(data)
	response.ObjectiveEvaluations = data
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// GetPeriodObjectiveDepartmentEvaluations retrieves department evaluations.
func (s *reviewPeriodService) GetPeriodObjectiveDepartmentEvaluations(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.PeriodObjectiveDepartmentEvaluationListResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var evals []performance.PeriodObjectiveDepartmentEvaluation
	err := s.db.WithContext(ctx).
		Joins("JOIN pms.period_objectives ON pms.period_objectives.period_objective_id = pms.period_objective_department_evaluations.period_objective_id").
		Where("pms.period_objectives.review_period_id = ? AND pms.period_objective_department_evaluations.soft_deleted = ?", reviewPeriodID, false).
		Preload("PeriodObjective").
		Preload("PeriodObjective.Objective").
		Preload("PeriodObjective.ReviewPeriod").
		Preload("Department").
		Find(&evals).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve department evaluations")
		return response, err
	}

	data := make([]performance.PeriodObjectiveDepartmentEvaluationData, 0, len(evals))
	for _, e := range evals {
		d := performance.PeriodObjectiveDepartmentEvaluationData{
			PeriodObjectiveDepartmentEvaluationID: e.PeriodObjectiveDepartmentEvaluationID,
			OverallOutcomeScored:                  e.OverallOutcomeScored,
			AllocatedOutcome:                      e.AllocatedOutcome,
			OutcomeScore:                           e.OutcomeScore,
			DepartmentID:                          e.DepartmentID,
			PeriodObjectiveID:                     e.PeriodObjectiveID,
		}
		d.RecordStatus = e.RecordStatus
		d.IsActive = e.IsActive
		d.CreatedBy = e.CreatedBy
		d.CreatedAt = e.CreatedAt

		if e.PeriodObjective != nil {
			if e.PeriodObjective.Objective != nil {
				d.EnterpriseObjectiveID = e.PeriodObjective.Objective.EnterpriseObjectiveID
			}
			if e.PeriodObjective.ReviewPeriod != nil {
				d.ReviewPeriodID = e.PeriodObjective.ReviewPeriod.PeriodID
				d.ReviewPeriod = e.PeriodObjective.ReviewPeriod.Name
			}
		}
		if e.Department != nil {
			d.DepartmentName = e.Department.DepartmentName
		}

		data = append(data, d)
	}

	response.TotalRecords = len(data)
	response.DepartmentObjectiveEvaluations = data
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// Period Scores
// ===========================================================================

// GetStaffPeriodScore retrieves the period score for a staff member.
func (s *reviewPeriodService) GetStaffPeriodScore(ctx context.Context, staffID, reviewPeriodID string) (interface{}, error) {
	response := &performance.PeriodScoreResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	score, err := s.periodScoreRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"ReviewPeriod", "Strategy"},
		"staff_id = ? AND review_period_id = ?",
		staffID, reviewPeriodID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve staff period score")
		return response, err
	}
	if score == nil {
		response.Message = "Period score record not found"
		return response, nil
	}

	data := &performance.PeriodScoreData{
		PeriodScoreID:     score.PeriodScoreID,
		ReviewPeriodID:    score.ReviewPeriodID,
		StaffID:           score.StaffID,
		FinalScore:        score.FinalScore,
		ScorePercentage:   score.ScorePercentage,
		FinalGrade:        int(score.FinalGrade),
		FinalGradeName:    score.FinalGrade.String(),
		EndDate:           score.EndDate,
		OfficeID:          score.OfficeID,
		MinNoOfObjectives: score.MinNoOfObjectives,
		MaxNoOfObjectives: score.MaxNoOfObjectives,
		StrategyID:        score.StrategyID,
		StaffGrade:        score.StaffGrade,
		LocationID:        score.LocationID,
		HRDDeductedPoints: score.HRDDeductedPoints,
		IsUnderPerforming: score.IsUnderPerforming,
	}
	data.IsActive = score.IsActive
	data.CreatedBy = score.CreatedBy
	data.CreatedAt = score.CreatedAt

	if score.ReviewPeriod != nil {
		data.ReviewPeriod = score.ReviewPeriod.Name
		data.Year = score.ReviewPeriod.Year
		data.MaxPoint = score.ReviewPeriod.MaxPoints
	}
	if score.Strategy != nil {
		data.StrategyName = score.Strategy.Name
	}

	response.PeriodScore = data
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ===========================================================================
// Additional Retrieval Methods
// (mirrors remaining .NET IReviewPeriodService methods)
// ===========================================================================

// GetReviewPeriods returns all review periods ordered by year (desc) and start
// date, checking each for an active extension. Mirrors .NET GetReviewPeriods().
func (s *reviewPeriodService) GetReviewPeriods(ctx context.Context) (interface{}, error) {
	response := &performance.GetAllReviewPeriodResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var reviewPeriods []performance.PerformanceReviewPeriod
	err := s.db.WithContext(ctx).
		Where("soft_deleted = ?", false).
		Order("year DESC, start_date ASC").
		Find(&reviewPeriods).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve review periods")
		return response, err
	}

	rpList := make([]performance.PerformanceReviewPeriodVm, 0, len(reviewPeriods))
	for _, rp := range reviewPeriods {
		rpVm := performance.PerformanceReviewPeriodVm{
			PeriodID:                   rp.PeriodID,
			Year:                       rp.Year,
			Range:                      int(rp.Range),
			RangeValue:                 rp.RangeValue,
			Name:                       rp.Name,
			Description:                rp.Description,
			ShortName:                  rp.ShortName,
			StartDate:                  rp.StartDate,
			EndDate:                    rp.EndDate,
			AllowObjectivePlanning:     rp.AllowObjectivePlanning,
			AllowWorkProductPlanning:   rp.AllowWorkProductPlanning,
			AllowWorkProductEvaluation: rp.AllowWorkProductEvaluation,
			MaxPoints:                  rp.MaxPoints,
			MinNoOfObjectives:          rp.MinNoOfObjectives,
			MaxNoOfObjectives:          rp.MaxNoOfObjectives,
			StrategyID:                 rp.StrategyID,
		}
		rpVm.RecordStatus = rp.RecordStatus
		rpVm.IsActive = rp.IsActive
		rpVm.IsApproved = rp.IsApproved
		rpVm.CreatedBy = rp.CreatedBy
		rpVm.CreatedAt = rp.CreatedAt

		// Check for an active extension on this period.
		ext, _ := s.extensionRepo.FirstOrDefault(ctx,
			"review_period_id = ? AND record_status = ?",
			rp.PeriodID, enums.StatusActive.String())
		if ext != nil {
			rpVm.HasExtension = true
		}

		rpList = append(rpList, rpVm)
	}

	response.TotalRecords = len(rpList)
	response.PerformanceReviewPeriods = rpList
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// GetReviewPeriodCategoryDefinitions retrieves category definitions for a
// review period, preloading Category and ReviewPeriod associations.
// Mirrors .NET GetReviewPeriodCategoryDefinitions(reviewPeriodId).
func (s *reviewPeriodService) GetReviewPeriodCategoryDefinitions(ctx context.Context, reviewPeriodID string) (interface{}, error) {
	response := &performance.ReviewPeriodCategoryDefinitionResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	var catDefs []performance.CategoryDefinition
	err := s.db.WithContext(ctx).
		Where("review_period_id = ? AND record_status != ? AND soft_deleted = ?",
			reviewPeriodID, enums.StatusCancelled.String(), false).
		Preload("Category").
		Preload("ReviewPeriod").
		Find(&catDefs).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve category definitions")
		return response, err
	}

	catDefData := make([]performance.CategoryDefinitionData, 0, len(catDefs))
	for _, def := range catDefs {
		data := performance.CategoryDefinitionData{
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
		data.RecordStatus = def.RecordStatus
		data.IsActive = def.IsActive
		data.IsApproved = def.IsApproved
		data.IsRejected = def.IsRejected
		data.CreatedBy = def.CreatedBy
		data.CreatedAt = def.CreatedAt
		data.UpdatedBy = def.UpdatedBy
		data.UpdatedAt = def.UpdatedAt
		data.ApprovedBy = def.ApprovedBy
		data.DateApproved = def.DateApproved
		data.RejectedBy = def.RejectedBy
		data.RejectionReason = def.RejectionReason
		data.DateRejected = def.DateRejected

		if def.Category != nil {
			data.CategoryName = def.Category.Name
		}

		// GradeGroupName: look up from competency schema.
		if def.GradeGroupID > 0 {
			var groupName string
			if err := s.db.WithContext(ctx).
				Table("\"CoreSchema\".\"job_grade_groups\"").
				Select("group_name").
				Where("job_grade_group_id = ?", def.GradeGroupID).
				Scan(&groupName).Error; err == nil && groupName != "" {
				data.GradeGroupName = groupName
			}
		}

		catDefData = append(catDefData, data)
	}

	response.TotalRecords = len(catDefData)
	response.CategoryDefinitions = catDefData
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// GetPlannedObjective retrieves a single planned objective by ID with its
// ReviewPeriod preloaded, then resolves the objective name based on the
// objective level. Mirrors .NET GetPlannedObjectiveAsync(plannedObjectiveId).
//
// NOTE: The .NET version also fetches ERP employee data to determine the
// objective level dynamically. In Go the objective level is stored directly on
// the entity, so we use that instead.
func (s *reviewPeriodService) GetPlannedObjective(ctx context.Context, plannedObjectiveID string) (interface{}, error) {
	response := &performance.PlannedObjectiveResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	po, err := s.plannedObjRepo.FirstOrDefaultWithPreload(ctx,
		[]string{"ReviewPeriod"},
		"planned_objective_id = ?", plannedObjectiveID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve planned objective")
		return response, err
	}
	if po == nil {
		response.Message = "Objective record not found"
		return response, nil
	}

	data := &performance.PlannedObjectiveData{
		PlannedObjectiveID: po.PlannedObjectiveID,
		ReviewPeriodID:     po.ReviewPeriodID,
		ObjectiveLevel:     fmt.Sprintf("%d", po.ObjectiveLevel),
		ObjectiveID:        po.ObjectiveID,
		RecordStatus:       po.RecordStatus,
		StaffID:            po.StaffID,
		CreatedBy:          po.CreatedBy,
		CreatedDate:        po.CreatedAt,
		IsApproved:         po.IsApproved,
		IsRejected:         po.IsRejected,
		IsActive:           po.IsActive,
		Approver:           po.ApprovedBy,
		Comment:            po.Remark,
	}
	if po.ReviewPeriod != nil {
		data.ReviewPeriod = po.ReviewPeriod.Name
		data.Year = po.ReviewPeriod.Year
	}

	// Resolve objective name/details based on level.
	switch po.ObjectiveLevel {
	case enums.ObjectiveLevelOffice:
		var offObj performance.OfficeObjective
		if err := s.db.WithContext(ctx).Where("office_objective_id = ?", po.ObjectiveID).First(&offObj).Error; err == nil {
			data.ObjectiveLevel = "Office"
			data.Objective = offObj.Name
			data.Kpi = offObj.Kpi
			data.Target = offObj.Target
			data.Description = offObj.Description
		}
	case enums.ObjectiveLevelDivision:
		var divObj performance.DivisionObjective
		if err := s.db.WithContext(ctx).Where("division_objective_id = ?", po.ObjectiveID).First(&divObj).Error; err == nil {
			data.ObjectiveLevel = "Division"
			data.Objective = divObj.Name
			data.Kpi = divObj.Kpi
			data.Target = divObj.Target
			data.Description = divObj.Description
		}
	case enums.ObjectiveLevelDepartment:
		var deptObj performance.DepartmentObjective
		if err := s.db.WithContext(ctx).Where("department_objective_id = ?", po.ObjectiveID).First(&deptObj).Error; err == nil {
			data.ObjectiveLevel = "Department"
			data.Objective = deptObj.Name
			data.Kpi = deptObj.Kpi
			data.Target = deptObj.Target
			data.Description = deptObj.Description
		}
	case enums.ObjectiveLevelEnterprise:
		var entObj performance.EnterpriseObjective
		if err := s.db.WithContext(ctx).Where("enterprise_objective_id = ?", po.ObjectiveID).First(&entObj).Error; err == nil {
			data.ObjectiveLevel = "Enterprise"
			data.Objective = entObj.Name
			data.Kpi = entObj.Kpi
			data.Target = entObj.Target
			data.Description = entObj.Description
		}
	}

	response.PlannedObjective = data
	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// GetEnterpriseObjectiveByLevel traces an objective at the given hierarchy
// level back to its root enterprise objective and returns the aggregated data.
// Mirrors .NET GetEnterpriseObjectiveByLevelAsync(objectiveId, objectiveLevel).
func (s *reviewPeriodService) GetEnterpriseObjectiveByLevel(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error) {
	response := &performance.EnterpriseObjectiveResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	level := enums.ObjectiveLevel(objectiveLevel)

	switch level {
	case enums.ObjectiveLevelOffice:
		// Office -> Division -> Department -> Enterprise
		var offObj performance.OfficeObjective
		err := s.db.WithContext(ctx).
			Where("office_objective_id = ?", objectiveID).
			Preload("DivisionObjective.DepartmentObjective.EnterpriseObjective").
			First(&offObj).Error
		if err != nil {
			s.log.Error().Err(err).Msg("failed to retrieve office objective")
			return response, err
		}
		if offObj.DivisionObjective != nil &&
			offObj.DivisionObjective.DepartmentObjective != nil &&
			offObj.DivisionObjective.DepartmentObjective.EnterpriseObjective != nil {
			ent := offObj.DivisionObjective.DepartmentObjective.EnterpriseObjective
			data := &performance.EnterpriseObjectiveData{
				EnterpriseObjectiveID:          ent.EnterpriseObjectiveID,
				EnterpriseObjectivesCategoryID: ent.EnterpriseObjectivesCategoryID,
				StrategyID:                     ent.StrategyID,
				Name:                           offObj.Name,
				Description:                    offObj.Description,
				Kpi:                            offObj.Kpi,
				Target:                         offObj.Target,
			}
			data.IsActive = offObj.IsActive
			data.IsApproved = offObj.IsApproved
			data.IsRejected = offObj.IsRejected
			data.CreatedBy = offObj.CreatedBy
			data.CreatedAt = offObj.CreatedAt
			data.UpdatedBy = offObj.UpdatedBy
			data.UpdatedAt = offObj.UpdatedAt
			data.ApprovedBy = offObj.ApprovedBy
			data.DateApproved = offObj.DateApproved
			data.RejectedBy = offObj.RejectedBy
			data.RejectionReason = offObj.RejectionReason
			data.DateRejected = offObj.DateRejected
			data.RecordStatus = offObj.RecordStatus
			response.EnterpriseObjective = data
		}

	case enums.ObjectiveLevelDivision:
		// Division -> Department -> Enterprise
		var divObj performance.DivisionObjective
		err := s.db.WithContext(ctx).
			Where("division_objective_id = ?", objectiveID).
			Preload("DepartmentObjective.EnterpriseObjective").
			First(&divObj).Error
		if err != nil {
			s.log.Error().Err(err).Msg("failed to retrieve division objective")
			return response, err
		}
		if divObj.DepartmentObjective != nil &&
			divObj.DepartmentObjective.EnterpriseObjective != nil {
			ent := divObj.DepartmentObjective.EnterpriseObjective
			data := &performance.EnterpriseObjectiveData{
				EnterpriseObjectiveID:          ent.EnterpriseObjectiveID,
				EnterpriseObjectivesCategoryID: ent.EnterpriseObjectivesCategoryID,
				StrategyID:                     ent.StrategyID,
				Name:                           divObj.Name,
				Description:                    divObj.Description,
				Kpi:                            divObj.Kpi,
				Target:                         divObj.Target,
			}
			data.IsActive = divObj.IsActive
			data.IsApproved = divObj.IsApproved
			data.IsRejected = divObj.IsRejected
			data.CreatedBy = divObj.CreatedBy
			data.CreatedAt = divObj.CreatedAt
			data.UpdatedBy = divObj.UpdatedBy
			data.UpdatedAt = divObj.UpdatedAt
			data.ApprovedBy = divObj.ApprovedBy
			data.DateApproved = divObj.DateApproved
			data.RejectedBy = divObj.RejectedBy
			data.RejectionReason = divObj.RejectionReason
			data.DateRejected = divObj.DateRejected
			data.RecordStatus = divObj.RecordStatus
			response.EnterpriseObjective = data
		}

	case enums.ObjectiveLevelDepartment:
		// Department -> Enterprise
		var deptObj performance.DepartmentObjective
		err := s.db.WithContext(ctx).
			Where("department_objective_id = ?", objectiveID).
			Preload("EnterpriseObjective").
			First(&deptObj).Error
		if err != nil {
			s.log.Error().Err(err).Msg("failed to retrieve department objective")
			return response, err
		}
		if deptObj.EnterpriseObjective != nil {
			ent := deptObj.EnterpriseObjective
			data := &performance.EnterpriseObjectiveData{
				EnterpriseObjectiveID:          ent.EnterpriseObjectiveID,
				EnterpriseObjectivesCategoryID: ent.EnterpriseObjectivesCategoryID,
				StrategyID:                     ent.StrategyID,
				Name:                           deptObj.Name,
				Description:                    deptObj.Description,
				Kpi:                            deptObj.Kpi,
				Target:                         deptObj.Target,
			}
			data.IsActive = deptObj.IsActive
			data.IsApproved = deptObj.IsApproved
			data.IsRejected = deptObj.IsRejected
			data.CreatedBy = deptObj.CreatedBy
			data.CreatedAt = deptObj.CreatedAt
			data.UpdatedBy = deptObj.UpdatedBy
			data.UpdatedAt = deptObj.UpdatedAt
			data.ApprovedBy = deptObj.ApprovedBy
			data.DateApproved = deptObj.DateApproved
			data.RejectedBy = deptObj.RejectedBy
			data.RejectionReason = deptObj.RejectionReason
			data.DateRejected = deptObj.DateRejected
			data.RecordStatus = deptObj.RecordStatus
			response.EnterpriseObjective = data
		}

	case enums.ObjectiveLevelEnterprise:
		// Already at enterprise level.
		var entObj performance.EnterpriseObjective
		err := s.db.WithContext(ctx).
			Where("enterprise_objective_id = ?", objectiveID).
			First(&entObj).Error
		if err != nil {
			s.log.Error().Err(err).Msg("failed to retrieve enterprise objective")
			return response, err
		}
		data := &performance.EnterpriseObjectiveData{
			EnterpriseObjectiveID:          entObj.EnterpriseObjectiveID,
			EnterpriseObjectivesCategoryID: entObj.EnterpriseObjectivesCategoryID,
			StrategyID:                     entObj.StrategyID,
			Name:                           entObj.Name,
			Description:                    entObj.Description,
			Kpi:                            entObj.Kpi,
			Target:                         entObj.Target,
		}
		data.IsActive = entObj.IsActive
		data.IsApproved = entObj.IsApproved
		data.IsRejected = entObj.IsRejected
		data.CreatedBy = entObj.CreatedBy
		data.CreatedAt = entObj.CreatedAt
		data.UpdatedBy = entObj.UpdatedBy
		data.UpdatedAt = entObj.UpdatedAt
		data.ApprovedBy = entObj.ApprovedBy
		data.DateApproved = entObj.DateApproved
		data.RejectedBy = entObj.RejectedBy
		data.RejectionReason = entObj.RejectionReason
		data.DateRejected = entObj.DateRejected
		data.RecordStatus = entObj.RecordStatus
		response.EnterpriseObjective = data

	default:
		response.Message = fmt.Sprintf("Unsupported objective level: %d", objectiveLevel)
		return response, nil
	}

	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ArchiveCancelledObjectives soft-deletes all cancelled planned objectives for
// a given staff member and review period by calling the stored procedure.
// Mirrors .NET ArchiveCancelledObjectiveAsync(StaffId, ReviewperiodId).
func (s *reviewPeriodService) ArchiveCancelledObjectives(ctx context.Context, staffID string, reviewPeriodID string) (interface{}, error) {
	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	err := s.db.WithContext(ctx).
		Exec("CALL pms.archive_cancelled_objectives(?, ?)", staffID, reviewPeriodID).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to archive cancelled objectives")
		response.Message = "An unexpected error occurred, try again or report to our support with code: ERR:00RR"
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// ArchiveCancelledWorkProducts soft-deletes all cancelled work products for a
// given staff member and review period by calling the stored procedure.
// Mirrors .NET ArchiveCancelledWorkProductAsync(StaffId, ReviewperiodId).
func (s *reviewPeriodService) ArchiveCancelledWorkProducts(ctx context.Context, staffID string, reviewPeriodID string) (interface{}, error) {
	response := &performance.ResponseVm{}
	response.HasError = true
	response.Message = "An error occurred"

	err := s.db.WithContext(ctx).
		Exec("CALL pms.archive_cancelled_workproducts(?, ?)", staffID, reviewPeriodID).Error
	if err != nil {
		s.log.Error().Err(err).Msg("failed to archive cancelled work products")
		response.Message = "An unexpected error occurred, try again or report to our support with code: ERR:00RR"
		return response, err
	}

	response.HasError = false
	response.Message = "Operation completed"
	return response, nil
}

// Compile-time interface compliance check.
var _ ReviewPeriodService = (*reviewPeriodService)(nil)
