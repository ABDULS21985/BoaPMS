package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// competencyService implements CompetencyService.
type competencyService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	reviewAgent *reviewAgentService // handles population & calculation

	competencyRepo          *repository.Repository[competency.Competency]
	categoryRepo            *repository.Repository[competency.CompetencyCategory]
	categoryGradingRepo     *repository.Repository[competency.CompetencyCategoryGrading]
	ratingDefRepo           *repository.Repository[competency.CompetencyRatingDefinition]
	reviewRepo              *repository.Repository[competency.CompetencyReview]
	reviewProfileRepo       *repository.Repository[competency.CompetencyReviewProfile]
	developmentPlanRepo     *repository.Repository[competency.DevelopmentPlan]
	jobRoleRepo             *repository.Repository[competency.JobRole]
	officeJobRoleRepo       *repository.Repository[competency.OfficeJobRole]
	jobRoleCompetencyRepo   *repository.Repository[competency.JobRoleCompetency]
	behavioralCompetencyRepo *repository.Repository[competency.BehavioralCompetency]
	jobRoleGradeRepo        *repository.Repository[competency.JobRoleGrade]
	jobGradeRepo            *repository.Repository[competency.JobGrade]
	jobGradeGroupRepo       *repository.Repository[competency.JobGradeGroup]
	assignJobGradeGroupRepo *repository.Repository[competency.AssignJobGradeGroup]
	ratingRepo              *repository.Repository[competency.Rating]
	reviewPeriodRepo        *repository.Repository[competency.ReviewPeriod]
	reviewTypeRepo          *repository.Repository[competency.ReviewType]
	trainingTypeRepo        *repository.Repository[competency.TrainingType]
	bankYearRepo            *repository.Repository[identity.BankYear]
}

func newCompetencyService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) CompetencyService {
	return &competencyService{
		db:                       repos.GormDB,
		cfg:                      cfg,
		log:                      log.With().Str("service", "competency").Logger(),
		reviewAgent:              newReviewAgentService(repos, cfg, log),
		competencyRepo:           repository.NewRepository[competency.Competency](repos.GormDB),
		categoryRepo:             repository.NewRepository[competency.CompetencyCategory](repos.GormDB),
		categoryGradingRepo:      repository.NewRepository[competency.CompetencyCategoryGrading](repos.GormDB),
		ratingDefRepo:            repository.NewRepository[competency.CompetencyRatingDefinition](repos.GormDB),
		reviewRepo:               repository.NewRepository[competency.CompetencyReview](repos.GormDB),
		reviewProfileRepo:        repository.NewRepository[competency.CompetencyReviewProfile](repos.GormDB),
		developmentPlanRepo:      repository.NewRepository[competency.DevelopmentPlan](repos.GormDB),
		jobRoleRepo:              repository.NewRepository[competency.JobRole](repos.GormDB),
		officeJobRoleRepo:        repository.NewRepository[competency.OfficeJobRole](repos.GormDB),
		jobRoleCompetencyRepo:    repository.NewRepository[competency.JobRoleCompetency](repos.GormDB),
		behavioralCompetencyRepo: repository.NewRepository[competency.BehavioralCompetency](repos.GormDB),
		jobRoleGradeRepo:         repository.NewRepository[competency.JobRoleGrade](repos.GormDB),
		jobGradeRepo:             repository.NewRepository[competency.JobGrade](repos.GormDB),
		jobGradeGroupRepo:        repository.NewRepository[competency.JobGradeGroup](repos.GormDB),
		assignJobGradeGroupRepo:  repository.NewRepository[competency.AssignJobGradeGroup](repos.GormDB),
		ratingRepo:               repository.NewRepository[competency.Rating](repos.GormDB),
		reviewPeriodRepo:         repository.NewRepository[competency.ReviewPeriod](repos.GormDB),
		reviewTypeRepo:           repository.NewRepository[competency.ReviewType](repos.GormDB),
		trainingTypeRepo:         repository.NewRepository[competency.TrainingType](repos.GormDB),
		bankYearRepo:             repository.NewRepository[identity.BankYear](repos.GormDB),
	}
}

// ---------------------------------------------------------------------------
// responseVm is a local response envelope matching the .NET ResponseVm pattern.
// ---------------------------------------------------------------------------

type responseVm struct {
	IsSuccess bool   `json:"isSuccess"`
	ID        string `json:"id,omitempty"`
	Message   string `json:"message"`
}

// ============================= Competencies =================================

func (s *competencyService) GetCompetencies(ctx context.Context, req interface{}) (interface{}, error) {
	search, ok := req.(*competency.SearchCompetencyVm)
	if !ok || search == nil {
		search = &competency.SearchCompetencyVm{BasePagedData: competency.BasePagedData{PageSize: 50}}
	}

	q := s.competencyRepo.Query(ctx).
		Joins("LEFT JOIN \"CoreSchema\".\"competency_categories\" ON \"CoreSchema\".\"competency_categories\".\"competency_category_id\" = \"CoreSchema\".\"competencies\".\"competency_category_id\"")

	if search.CategoryID != nil && *search.CategoryID > 0 {
		q = q.Where("\"CoreSchema\".\"competencies\".\"competency_category_id\" = ?", *search.CategoryID)
	}
	if search.IsTechnical != nil {
		q = q.Where("\"CoreSchema\".\"competency_categories\".\"is_technical\" = ?", *search.IsTechnical)
	}
	if search.IsApproved != nil {
		q = q.Where("\"CoreSchema\".\"competencies\".\"is_approved\" = ?", *search.IsApproved)
	}
	if search.IsRejected != nil {
		q = q.Where("\"CoreSchema\".\"competencies\".\"is_rejected\" = ?", *search.IsRejected)
	}
	if search.SearchString != "" {
		q = q.Where("UPPER(\"CoreSchema\".\"competencies\".\"competency_name\") LIKE ?", "%"+strings.ToUpper(search.SearchString)+"%")
	}

	var totalRecords int64
	q.Count(&totalRecords)

	var entities []competency.Competency
	err := q.Preload("CompetencyCategory").
		Order("competency_name").
		Offset(search.Skip).Limit(pageSize(search.PageSize)).
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get competencies: %w", err)
	}

	vms := make([]competency.CompetencyVm, 0, len(entities))
	for _, e := range entities {
		catName := ""
		if e.CompetencyCategory != nil {
			catName = e.CompetencyCategory.CategoryName
		}
		vms = append(vms, competency.CompetencyVm{
			CompetencyID:           e.CompetencyID,
			CompetencyName:         e.CompetencyName,
			CompetencyCategoryID:   e.CompetencyCategoryID,
			CompetencyCategoryName: catName,
			Description:            e.Description,
			IsApproved:             e.IsApproved,
			IsRejected:             e.IsRejected,
			RejectedBy:             e.RejectedBy,
			RejectionReason:        e.RejectionReason,
			BaseAuditVm:           toBaseAuditVm(e.BaseWorkFlowData.BaseAudit),
		})
	}

	return &competency.CompetencyListVm{
		Competencies: vms,
		TotalRecord:  int(totalRecords),
	}, nil
}

func (s *competencyService) SaveCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.SaveCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetency")
	}

	var message string
	var entity competency.Competency

	if vm.CompetencyID > 0 {
		existing, err := s.competencyRepo.GetByID(ctx, vm.CompetencyID)
		if err != nil {
			return nil, fmt.Errorf("save competency: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Competency not found"}, nil
		}
		existing.CompetencyCategoryID = vm.CompetencyCategoryID
		existing.CompetencyName = strings.ToUpper(vm.CompetencyName)
		existing.Description = vm.Description
		existing.IsActive = vm.IsActive
		existing.IsRejected = false

		if err := s.competencyRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: "Duplicate competency"}, nil
		}
		entity = *existing
		message = "Updated successfully"
	} else {
		entity = competency.Competency{
			CompetencyCategoryID: vm.CompetencyCategoryID,
			CompetencyName:       strings.ToUpper(vm.CompetencyName),
			Description:          vm.Description,
		}
		entity.IsActive = vm.IsActive
		entity.IsRejected = false

		if err := s.competencyRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: "Duplicate competency"}, nil
		}
		message = "Created successfully"
	}

	return &responseVm{
		IsSuccess: true,
		ID:        fmt.Sprintf("%d", entity.CompetencyID),
		Message:   message,
	}, nil
}

func (s *competencyService) ApproveCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.ApproveCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for ApproveCompetency")
	}

	existing, err := s.competencyRepo.GetByID(ctx, vm.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("approve competency: %w", err)
	}
	if existing == nil {
		return &responseVm{IsSuccess: false, Message: "Invalid Competency ID"}, nil
	}

	existing.ApprovedBy = vm.ApprovedBy
	existing.IsApproved = vm.IsApproved
	existing.DateApproved = vm.DateApproved
	existing.IsRejected = false

	if err := s.competencyRepo.Update(ctx, existing); err != nil {
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}

	return &responseVm{IsSuccess: true, Message: "Approved successfully"}, nil
}

func (s *competencyService) RejectCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.RejectCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for RejectCompetency")
	}

	existing, err := s.competencyRepo.GetByID(ctx, vm.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("reject competency: %w", err)
	}
	if existing == nil {
		return &responseVm{IsSuccess: false, Message: "Invalid Competency ID"}, nil
	}

	existing.IsApproved = false
	existing.RejectedBy = vm.RejectedBy
	existing.DateRejected = vm.DateRejected
	existing.IsRejected = vm.IsRejected
	existing.RejectionReason = vm.RejectionReason

	if err := s.competencyRepo.Update(ctx, existing); err != nil {
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}

	return &responseVm{IsSuccess: true, Message: "Rejected successfully"}, nil
}

// ======================== Competency Categories =============================

func (s *competencyService) GetCompetencyCategories(ctx context.Context) (interface{}, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get competency categories: %w", err)
	}

	vms := make([]competency.CompetencyCategoryVm, 0, len(categories))
	for _, c := range categories {
		vms = append(vms, competency.CompetencyCategoryVm{
			CompetencyCategoryID: c.CompetencyCategoryID,
			CategoryName:         c.CategoryName,
			IsTechnical:          c.IsTechnical,
			BaseAuditVm:          toBaseAuditVm(c.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveCompetencyCategory(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyCategoryVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetencyCategory")
	}

	var message string
	if vm.CompetencyCategoryID > 0 {
		existing, err := s.categoryRepo.GetByID(ctx, vm.CompetencyCategoryID)
		if err != nil {
			return nil, fmt.Errorf("save competency category: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Category not found"}, nil
		}
		existing.CategoryName = strings.ToUpper(vm.CategoryName)
		existing.IsTechnical = vm.IsTechnical
		existing.IsActive = vm.IsActive

		if err := s.categoryRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Competency Category has been updated successfully", existing.CategoryName)
	} else {
		entity := competency.CompetencyCategory{
			CategoryName: strings.ToUpper(vm.CategoryName),
			IsTechnical:  vm.IsTechnical,
		}
		entity.IsActive = vm.IsActive

		if err := s.categoryRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Competency Category has been created successfully", entity.CategoryName)
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ==================== Competency Category Gradings ==========================

func (s *competencyService) GetCompetencyCategoryGradings(ctx context.Context) (interface{}, error) {
	var entities []competency.CompetencyCategoryGrading
	err := s.categoryGradingRepo.Query(ctx).
		Preload("CompetencyCategory").
		Preload("ReviewType").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get category gradings: %w", err)
	}

	vms := make([]competency.CompetencyCategoryGradingVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.CompetencyCategoryGradingVm{
			CompetencyCategoryGradingID: e.CompetencyCategoryGradingID,
			CompetencyCategoryID:        e.CompetencyCategoryID,
			ReviewTypeID:                e.ReviewTypeID,
			WeightPercentage:            e.WeightPercentage,
			BaseAuditVm:                 toBaseAuditVm(e.BaseAudit),
		}
		if e.CompetencyCategory != nil {
			vm.CompetencyCategoryName = e.CompetencyCategory.CategoryName
		}
		if e.ReviewType != nil {
			vm.ReviewTypeName = e.ReviewType.ReviewTypeName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) SaveCompetencyCategoryGrading(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyCategoryGradingVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetencyCategoryGrading")
	}

	var message string
	if vm.CompetencyCategoryGradingID > 0 {
		existing, err := s.categoryGradingRepo.GetByID(ctx, vm.CompetencyCategoryGradingID)
		if err != nil {
			return nil, fmt.Errorf("save category grading: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Category Grading not found"}, nil
		}
		existing.CompetencyCategoryID = vm.CompetencyCategoryID
		existing.ReviewTypeID = vm.ReviewTypeID
		existing.WeightPercentage = vm.WeightPercentage
		existing.IsActive = vm.IsActive

		if err := s.categoryGradingRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Category Grading has been updated successfully"
	} else {
		entity := competency.CompetencyCategoryGrading{
			CompetencyCategoryID: vm.CompetencyCategoryID,
			ReviewTypeID:         vm.ReviewTypeID,
			WeightPercentage:     vm.WeightPercentage,
		}
		entity.IsActive = vm.IsActive

		if err := s.categoryGradingRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Category Grading has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// =================== Competency Rating Definitions ==========================

func (s *competencyService) GetCompetencyRatingDefinitions(ctx context.Context, competencyId *int) (interface{}, error) {
	q := s.ratingDefRepo.Query(ctx).
		Preload("Rating").
		Preload("Competency")

	if competencyId != nil && *competencyId > 0 {
		q = q.Where("competency_id = ?", *competencyId)
	}

	var entities []competency.CompetencyRatingDefinition
	err := q.Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get rating definitions: %w", err)
	}

	vms := make([]competency.CompetencyRatingDefinitionVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.CompetencyRatingDefinitionVm{
			CompetencyRatingDefinitionID: e.CompetencyRatingDefinitionID,
			CompetencyID:                 e.CompetencyID,
			RatingID:                     e.RatingID,
			Definition:                   e.Definition,
			BaseAuditVm:                  toBaseAuditVm(e.BaseAudit),
		}
		if e.Rating != nil {
			vm.RatingName = e.Rating.Name
			vm.RatingValue = e.Rating.Value
		}
		if e.Competency != nil {
			vm.CompetencyName = e.Competency.CompetencyName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) SaveCompetencyRatingDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyRatingDefinitionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetencyRatingDefinition")
	}

	var message string
	if vm.CompetencyRatingDefinitionID > 0 {
		existing, err := s.ratingDefRepo.GetByID(ctx, vm.CompetencyRatingDefinitionID)
		if err != nil {
			return nil, fmt.Errorf("save rating definition: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Rating Definition not found"}, nil
		}
		existing.CompetencyID = vm.CompetencyID
		existing.RatingID = vm.RatingID
		existing.Definition = vm.Definition
		existing.IsActive = vm.IsActive

		if err := s.ratingDefRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Rating definition has been updated successfully"
	} else {
		entity := competency.CompetencyRatingDefinition{
			CompetencyID: vm.CompetencyID,
			RatingID:     vm.RatingID,
			Definition:   vm.Definition,
		}
		entity.IsActive = vm.IsActive

		if err := s.ratingDefRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Rating definition has been Created Successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ======================== Competency Reviews =================================

func (s *competencyService) GetCompetencyReviews(ctx context.Context) (interface{}, error) {
	var entities []competency.CompetencyReview
	err := s.reviewRepo.Query(ctx).
		Preload("Competency").
		Preload("ReviewType").
		Preload("ReviewPeriod").
		Preload("ExpectedRating").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get competency reviews: %w", err)
	}

	return mapReviewsToVms(entities), nil
}

func (s *competencyService) GetCompetencyReviewByReviewer(ctx context.Context, reviewerId string, reviewPeriodId *int) (interface{}, error) {
	q := s.reviewRepo.Query(ctx).
		Where("reviewer_id = ?", reviewerId).
		Preload("Competency").
		Preload("Competency.CompetencyCategory").
		Preload("ReviewType").
		Preload("ReviewPeriod").
		Preload("ExpectedRating")

	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}

	var entities []competency.CompetencyReview
	if err := q.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("get reviews by reviewer: %w", err)
	}

	return mapReviewsToVms(entities), nil
}

func (s *competencyService) GetCompetencyReviewForEmployee(ctx context.Context, employeeNumber string, reviewPeriodId *int) (interface{}, error) {
	q := s.reviewRepo.Query(ctx).
		Where("employee_number = ?", employeeNumber).
		Preload("Competency").
		Preload("ReviewType").
		Preload("ReviewPeriod").
		Preload("ExpectedRating")

	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}

	var entities []competency.CompetencyReview
	if err := q.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("get reviews for employee: %w", err)
	}

	return mapReviewsToVms(entities), nil
}

func (s *competencyService) GetCompetencyReviewDetail(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.SearchForReviewDetailVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for GetCompetencyReviewDetail")
	}

	var entities []competency.CompetencyReview
	err := s.reviewRepo.Query(ctx).
		Where("reviewer_id = ? AND is_technical = ? AND employee_number = ? AND review_period_id = ? AND review_type_id = ?",
			vm.ReviewerID, vm.IsTechnical, vm.EmployeeID, vm.ReviewPeriodID, vm.ReviewTypeID).
		Preload("Competency").
		Preload("Competency.CompetencyCategory").
		Preload("Competency.CompetencyRatingDefinitions").
		Preload("Competency.CompetencyRatingDefinitions.Rating").
		Preload("ReviewType").
		Preload("ReviewPeriod").
		Preload("ExpectedRating").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get review detail: %w", err)
	}

	reviews := make([]competency.CompetencyReviewVm, 0, len(entities))
	for _, e := range entities {
		rv := mapSingleReviewToVm(e)
		// Attach rating definitions
		if e.Competency != nil {
			rv.CompetencyDefinition = e.Competency.Description
			if e.Competency.CompetencyCategory != nil {
				rv.CompetencyCategoryName = e.Competency.CompetencyCategory.CategoryName
			}
			defs := make([]competency.CompetencyRatingDefinitionVm, 0, len(e.Competency.CompetencyRatingDefinitions))
			for _, d := range e.Competency.CompetencyRatingDefinitions {
				dvm := competency.CompetencyRatingDefinitionVm{
					Definition: d.Definition,
				}
				if d.Rating != nil {
					dvm.RatingName = d.Rating.Name
					dvm.RatingValue = d.Rating.Value
				}
				defs = append(defs, dvm)
			}
			rv.CompetencyRatingDefinitions = defs
		}
		reviews = append(reviews, rv)
	}

	return &competency.CompetencyReviewDetailVm{
		CompetencyReviews: reviews,
	}, nil
}

func (s *competencyService) SaveCompetencyReview(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyReviewVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetencyReview")
	}

	var message string
	if vm.CompetencyReviewID > 0 {
		existing, err := s.reviewRepo.GetByID(ctx, vm.CompetencyReviewID)
		if err != nil {
			return nil, fmt.Errorf("save competency review: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Review not found"}, nil
		}
		existing.CompetencyID = vm.CompetencyID
		existing.ReviewDate = vm.ReviewDate
		existing.ExpectedRatingID = vm.ExpectedRatingID
		existing.ActualRatingID = vm.ActualRatingID
		existing.ActualRatingName = vm.ActualRatingName
		existing.ActualRatingValue = vm.ActualRatingValue
		existing.IsActive = vm.IsActive

		if err := s.reviewRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Review has been updated successfully"
	} else {
		entity := competency.CompetencyReview{
			EmployeeNumber: vm.EmployeeNumber,
			CompetencyID:   vm.CompetencyID,
			ReviewDate:     vm.ReviewDate,
			ReviewPeriodID: vm.ReviewPeriodID,
			ReviewTypeID:   vm.ReviewTypeID,
			ReviewerID:     vm.ReviewerID,
			ReviewerName:   vm.ReviewerName,
			ExpectedRatingID: vm.ExpectedRatingID,
		}
		entity.IsActive = vm.IsActive

		if err := s.reviewRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Review has been Created Successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ==================== Competency Review Profiles ============================

func (s *competencyService) GetCompetencyReviewProfiles(ctx context.Context, employeeNumber string, reviewPeriodId *int) (interface{}, error) {
	q := s.reviewProfileRepo.Query(ctx).
		Where("employee_number = ?", employeeNumber).
		Preload("DevelopmentPlans")

	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}

	var entities []competency.CompetencyReviewProfile
	if err := q.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("get review profiles: %w", err)
	}

	return mapProfilesToVms(entities), nil
}

func (s *competencyService) GetOfficeCompetencyReviews(ctx context.Context, officeId int, reviewPeriodId *int) (interface{}, error) {
	q := s.reviewProfileRepo.Query(ctx).
		Where("office_id = ?", fmt.Sprintf("%d", officeId))

	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}

	var profiles []competency.CompetencyReviewProfile
	if err := q.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("get office reviews: %w", err)
	}

	// Build overview report matching .NET ManagerOverviewReportVm structure
	type basicEmployeeData struct {
		EmployeeNumber       string `json:"employeeNumber"`
		FullName             string `json:"fullName"`
		Department           string `json:"department"`
		Grade                string `json:"grade"`
		Office               string `json:"office"`
		Position             string `json:"position"`
		NoOfCompletedReviews int    `json:"noOfCompletedReviews"`
		NoOfNotCompleted     int    `json:"noOfNotCompletedReviews"`
	}

	type managerOverview struct {
		NoNotStartedReviews int                 `json:"noNotStartedReviews"`
		NoStartedReviews    int                 `json:"noStartedReviews"`
		NoOfCompletedReviews int                `json:"noOfCompletedReviews"`
		BasicEmployeeDatas  []basicEmployeeData `json:"basicEmployeeDatas"`
	}

	result := managerOverview{BasicEmployeeDatas: []basicEmployeeData{}}

	employeeMap := make(map[string]bool)
	for _, p := range profiles {
		if p.AverageRatingID == 0 {
			result.NoNotStartedReviews++
		}
		if p.AverageRatingName == "" {
			result.NoStartedReviews++
		} else {
			result.NoOfCompletedReviews++
		}

		if !employeeMap[p.EmployeeNumber] {
			employeeMap[p.EmployeeNumber] = true
			completed := 0
			notCompleted := 0
			for _, ep := range profiles {
				if ep.EmployeeNumber == p.EmployeeNumber {
					if ep.AverageRatingName != "" {
						completed++
					} else {
						notCompleted++
					}
				}
			}
			result.BasicEmployeeDatas = append(result.BasicEmployeeDatas, basicEmployeeData{
				EmployeeNumber:       p.EmployeeNumber,
				FullName:             p.EmployeeName,
				Department:           p.DepartmentName,
				Grade:                p.GradeName,
				Office:               p.OfficeName,
				Position:             p.JobRoleName,
				NoOfCompletedReviews: completed,
				NoOfNotCompleted:     notCompleted,
			})
		}
	}

	return result, nil
}

func (s *competencyService) GetGroupCompetencyReviewProfiles(ctx context.Context, reviewPeriodId, officeId, divisionId, departmentId *int) (interface{}, error) {
	q := s.reviewProfileRepo.Query(ctx)
	q = applyOrgFilter(q, reviewPeriodId, officeId, divisionId, departmentId)

	var profiles []competency.CompetencyReviewProfile
	if err := q.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("get group review profiles: %w", err)
	}

	// Fetch ratings for stats
	var ratings []competency.Rating
	s.ratingRepo.Query(ctx).Where("is_active = ?", true).Find(&ratings)

	type chartData struct {
		Label    string  `json:"label"`
		Actual   float64 `json:"actual"`
		Expected float64 `json:"expected"`
	}
	type competencyRatingStat struct {
		RatingOrder     int     `json:"ratingOrder"`
		RatingName      string  `json:"ratingName"`
		RatingValue     int     `json:"ratingValue"`
		NumberOfStaff   int     `json:"numberOfStaff"`
		StaffPercentage float64 `json:"staffPercentage"`
	}
	type categoryDetailStat struct {
		CategoryName          string                 `json:"categoryName"`
		AverageRating         float64                `json:"averageRating"`
		HighestRating         int                    `json:"highestRating"`
		LowestRating          int                    `json:"lowestRating"`
		MostCommonRating      float64                `json:"mostCommonRating"`
		GroupCompetencyRatings []chartData           `json:"groupCompetencyRatings"`
		CompetencyRatingStat  []competencyRatingStat `json:"competencyRatingStat"`
	}
	type categoryStat struct {
		CategoryName string  `json:"categoryName"`
		Actual       float64 `json:"actual"`
		Expected     float64 `json:"expected"`
	}
	type grouped struct {
		CategoryCompetencyStats       []categoryStat       `json:"categoryCompetencyStats"`
		CategoryCompetencyDetailStats []categoryDetailStat `json:"categoryCompetencyDetailStats"`
	}

	result := grouped{
		CategoryCompetencyStats:       []categoryStat{},
		CategoryCompetencyDetailStats: []categoryDetailStat{},
	}

	categories := uniqueStrings(profiles, func(p competency.CompetencyReviewProfile) string { return p.CompetencyCategoryName })

	for _, cat := range categories {
		catProfiles := filterProfiles(profiles, func(p competency.CompetencyReviewProfile) bool { return p.CompetencyCategoryName == cat })

		result.CategoryCompetencyStats = append(result.CategoryCompetencyStats, categoryStat{
			CategoryName: cat,
			Actual:       avgField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
			Expected:     avgField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.ExpectedRatingValue }),
		})

		detail := categoryDetailStat{
			CategoryName:          cat,
			AverageRating:         avgField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
			HighestRating:         maxField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
			LowestRating:          minField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
			MostCommonRating:      avgField(catProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
			GroupCompetencyRatings: []chartData{},
			CompetencyRatingStat:  []competencyRatingStat{},
		}

		compNames := uniqueStrings(catProfiles, func(p competency.CompetencyReviewProfile) string { return p.CompetencyName })
		for _, cn := range compNames {
			compProfiles := filterProfiles(catProfiles, func(p competency.CompetencyReviewProfile) bool { return p.CompetencyName == cn })
			detail.GroupCompetencyRatings = append(detail.GroupCompetencyRatings, chartData{
				Label:    cn,
				Actual:   avgField(compProfiles, func(p competency.CompetencyReviewProfile) int { return p.AverageRatingValue }),
				Expected: avgField(compProfiles, func(p competency.CompetencyReviewProfile) int { return p.ExpectedRatingValue }),
			})
		}

		for _, r := range ratings {
			detail.CompetencyRatingStat = append(detail.CompetencyRatingStat, competencyRatingStat{
				RatingOrder: r.RatingID,
				RatingName:  r.Name,
				RatingValue: r.Value,
			})
		}

		result.CategoryCompetencyDetailStats = append(result.CategoryCompetencyDetailStats, detail)
	}

	return result, nil
}

func (s *competencyService) GetCompetencyMatrixReviewProfiles(ctx context.Context, reviewPeriodId, officeId, divisionId, departmentId *int) (interface{}, error) {
	q := s.reviewProfileRepo.Query(ctx)
	q = applyOrgFilter(q, reviewPeriodId, officeId, divisionId, departmentId)
	q = q.Where("LOWER(competency_category_name) != ?", "technical")

	var profiles []competency.CompetencyReviewProfile
	if err := q.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("get matrix review profiles: %w", err)
	}

	return buildMatrixResult(profiles), nil
}

func (s *competencyService) GetTechnicalCompetencyMatrixReviewProfiles(ctx context.Context, reviewPeriodId *int, jobRoleId int) (interface{}, error) {
	q := s.reviewProfileRepo.Query(ctx).
		Where("LOWER(competency_category_name) = ?", "technical").
		Where("job_role_id = ?", fmt.Sprintf("%d", jobRoleId))

	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}

	var profiles []competency.CompetencyReviewProfile
	if err := q.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("get technical matrix profiles: %w", err)
	}

	return buildMatrixResult(profiles), nil
}

func (s *competencyService) SaveCompetencyReviewProfile(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyReviewProfileVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveCompetencyReviewProfile")
	}

	var message string
	if vm.CompetencyReviewProfileID > 0 {
		existing, err := s.reviewProfileRepo.GetByID(ctx, vm.CompetencyReviewProfileID)
		if err != nil {
			return nil, fmt.Errorf("save review profile: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Review Profile not found"}, nil
		}
		existing.EmployeeNumber = vm.EmployeeNumber
		existing.CompetencyID = vm.CompetencyID
		existing.CompetencyName = vm.CompetencyName
		existing.ExpectedRatingID = vm.ExpectedRatingID
		existing.ExpectedRatingName = vm.ExpectedRatingName
		existing.ExpectedRatingValue = vm.ExpectedRatingValue
		existing.ReviewPeriodID = vm.ReviewPeriodID
		existing.ReviewPeriodName = vm.ReviewPeriodName
		existing.AverageRatingID = vm.AverageRatingID
		existing.AverageRatingName = vm.AverageRatingName
		existing.AverageRatingValue = vm.AverageRatingValue
		existing.AverageScore = vm.AverageScore
		existing.CompetencyCategoryName = vm.CompetencyCategoryName
		existing.OfficeID = vm.OfficeID
		existing.OfficeName = vm.OfficeName
		existing.DivisionID = vm.DivisionID
		existing.DivisionName = vm.DivisionName
		existing.DepartmentID = vm.DepartmentID
		existing.DepartmentName = vm.DepartmentName
		existing.JobRoleID = vm.JobRoleID
		existing.JobRoleName = vm.JobRoleName
		existing.GradeName = vm.GradeName
		existing.IsActive = vm.IsActive

		if err := s.reviewProfileRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Review Profile has been updated successfully"
	} else {
		entity := competency.CompetencyReviewProfile{
			EmployeeNumber:     vm.EmployeeNumber,
			CompetencyID:       vm.CompetencyID,
			CompetencyName:     vm.CompetencyName,
			ExpectedRatingID:   vm.ExpectedRatingID,
			ExpectedRatingName: vm.ExpectedRatingName,
			ExpectedRatingValue: vm.ExpectedRatingValue,
			ReviewPeriodID:     vm.ReviewPeriodID,
			ReviewPeriodName:   vm.ReviewPeriodName,
			AverageRatingID:    vm.AverageRatingID,
			AverageRatingName:  vm.AverageRatingName,
			AverageRatingValue: vm.AverageRatingValue,
			AverageScore:       vm.AverageScore,
			CompetencyCategoryName: vm.CompetencyCategoryName,
			OfficeID:           vm.OfficeID,
			OfficeName:         vm.OfficeName,
			DivisionID:         vm.DivisionID,
			DivisionName:       vm.DivisionName,
			DepartmentID:       vm.DepartmentID,
			DepartmentName:     vm.DepartmentName,
			JobRoleID:          vm.JobRoleID,
			JobRoleName:        vm.JobRoleName,
			GradeName:          vm.GradeName,
		}
		entity.IsActive = vm.IsActive

		if err := s.reviewProfileRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Review Profile has been Created Successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ======================== Competency Gaps ====================================

func (s *competencyService) GetCompetencyGaps(ctx context.Context, employeeNumber string) (interface{}, error) {
	var entities []competency.CompetencyReviewProfile
	err := s.reviewProfileRepo.Query(ctx).
		Where("employee_number = ? AND have_gap = ?", employeeNumber, true).
		Preload("DevelopmentPlans").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get competency gaps: %w", err)
	}

	return mapProfilesToVms(entities), nil
}

func (s *competencyService) CloseCompetencyGap(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CompetencyReviewProfileVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CloseCompetencyGap")
	}

	if vm.CompetencyReviewProfileID <= 0 {
		return &responseVm{IsSuccess: false, Message: "Competency Profile does not exist!"}, nil
	}

	existing, err := s.reviewProfileRepo.GetByID(ctx, vm.CompetencyReviewProfileID)
	if err != nil {
		return nil, fmt.Errorf("close competency gap: %w", err)
	}
	if existing == nil {
		return &responseVm{IsSuccess: false, Message: "Competency Profile does not exist!"}, nil
	}

	message := ""
	if existing.HaveGap {
		existing.AverageRatingID = existing.ExpectedRatingID
		existing.AverageRatingName = existing.ExpectedRatingName
		existing.AverageRatingValue = existing.ExpectedRatingValue
		existing.AverageScore = float64(existing.ExpectedRatingID)
		existing.HaveGap = false
		existing.CompetencyGap = 0
		existing.UpdatedBy = vm.EmployeeNumber

		if err := s.reviewProfileRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Competency Gap has been closed successfully!"
	} else {
		message = "Competency Gap cannot not be closed when you still have a gap!!!"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ======================== Development Plans ==================================

func (s *competencyService) GetDevelopmentPlans(ctx context.Context, competencyProfileReviewId *int) (interface{}, error) {
	q := s.developmentPlanRepo.Query(ctx).
		Preload("CompetencyReviewProfile")

	if competencyProfileReviewId != nil && *competencyProfileReviewId > 0 {
		q = q.Where("competency_review_profile_id = ?", *competencyProfileReviewId)
	} else {
		// .NET returns empty when no filter: Where(x => x.EmployeeNumber == null)
		q = q.Where("employee_number IS NULL")
	}

	var entities []competency.DevelopmentPlan
	if err := q.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("get development plans: %w", err)
	}

	vms := make([]competency.DevelopmentPlanVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.DevelopmentPlanVm{
			DevelopmentPlanID:         e.DevelopmentPlanID,
			EmployeeNumber:            e.EmployeeNumber,
			CompetencyReviewProfileID: e.CompetencyReviewProfileID,
			Activity:                  e.Activity,
			CompletionDate:            e.CompletionDate,
			LearningResource:          e.LearningResource,
			TargetDate:                e.TargetDate,
			TaskStatus:                e.TaskStatus,
			TrainingTypeName:          e.TrainingTypeName,
			BaseAuditVm:              toBaseAuditVm(e.BaseAudit),
		}
		if e.CompetencyReviewProfile != nil {
			vm.CompetencyCategoryName = e.CompetencyReviewProfile.CompetencyCategoryName
			vm.CompetencyName = e.CompetencyReviewProfile.CompetencyName
			vm.CurrentGap = e.CompetencyReviewProfile.CompetencyGap
			vm.ReviewPeriod = e.CompetencyReviewProfile.ReviewPeriodName
			vm.EmployeeName = e.CompetencyReviewProfile.EmployeeName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) SaveDevelopmentPlan(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.DevelopmentPlanVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveDevelopmentPlan")
	}

	var message string
	if vm.DevelopmentPlanID > 0 {
		existing, err := s.developmentPlanRepo.GetByID(ctx, vm.DevelopmentPlanID)
		if err != nil {
			return nil, fmt.Errorf("save development plan: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Development Plan not found"}, nil
		}
		existing.EmployeeNumber = vm.EmployeeNumber
		existing.CompetencyReviewProfileID = vm.CompetencyReviewProfileID
		existing.Activity = vm.Activity
		existing.CompletionDate = vm.CompletionDate
		existing.LearningResource = vm.LearningResource
		existing.TrainingTypeName = vm.TrainingTypeName
		existing.TargetDate = vm.TargetDate
		existing.TaskStatus = vm.TaskStatus
		existing.IsActive = vm.IsActive

		if err := s.developmentPlanRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Development Plan has been updated successfully"
	} else {
		entity := competency.DevelopmentPlan{
			EmployeeNumber:            vm.EmployeeNumber,
			CompetencyReviewProfileID: vm.CompetencyReviewProfileID,
			Activity:                  vm.Activity,
			CompletionDate:            vm.CompletionDate,
			LearningResource:          vm.LearningResource,
			TargetDate:                vm.TargetDate,
			TaskStatus:                vm.TaskStatus,
			TrainingTypeName:          vm.TrainingTypeName,
		}
		entity.IsActive = vm.IsActive

		if err := s.developmentPlanRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Development Plan has been Created Successfully"
	}

	// Handle gap closure when task status is ClosedGap
	if strings.EqualFold(vm.TaskStatus, "ClosedGap") {
		profile, err := s.reviewProfileRepo.GetByID(ctx, vm.CompetencyReviewProfileID)
		if err == nil && profile != nil {
			profile.HaveGap = false
			profile.AverageRatingID = profile.ExpectedRatingID
			profile.AverageRatingValue = profile.ExpectedRatingValue
			profile.AverageScore = float64(profile.ExpectedRatingValue)
			profile.AverageRatingName = profile.ExpectedRatingName
			_ = s.reviewProfileRepo.Update(ctx, profile)
			message = "Competency Gap has been Closed successfully"
		}
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ============================== Job Roles ====================================

func (s *competencyService) GetJobRoles(ctx context.Context) (interface{}, error) {
	roles, err := s.jobRoleRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get job roles: %w", err)
	}

	vms := make([]competency.JobRoleVm, 0, len(roles))
	for _, r := range roles {
		vms = append(vms, competency.JobRoleVm{
			JobRoleID:   r.JobRoleID,
			JobRoleName: r.JobRoleName,
			Description: r.Description,
			BaseAuditVm: toBaseAuditVm(r.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) GetJobRoleByName(ctx context.Context, jobRoleName string) (interface{}, error) {
	var entity competency.JobRole
	err := s.jobRoleRepo.Query(ctx).
		Where("LOWER(job_role_name) = LOWER(?)", jobRoleName).
		First(&entity).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("get job role by name: %w", err)
	}

	return &competency.JobRoleVm{
		JobRoleID:   entity.JobRoleID,
		JobRoleName: entity.JobRoleName,
		Description: entity.Description,
		BaseAuditVm: toBaseAuditVm(entity.BaseAudit),
	}, nil
}

func (s *competencyService) SaveJobRole(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.JobRoleVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveJobRole")
	}

	var message string
	if vm.JobRoleID > 0 {
		existing, err := s.jobRoleRepo.GetByID(ctx, vm.JobRoleID)
		if err != nil {
			return nil, fmt.Errorf("save job role: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Job Role not found"}, nil
		}
		existing.JobRoleName = vm.JobRoleName
		existing.Description = vm.Description
		existing.IsActive = vm.IsActive

		if err := s.jobRoleRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Job Role has been updated successfully"
	} else {
		entity := competency.JobRole{
			JobRoleName: vm.JobRoleName,
			Description: vm.Description,
		}
		entity.IsActive = vm.IsActive

		if err := s.jobRoleRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Job Role has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ========================== Office Job Roles =================================

func (s *competencyService) GetOfficeJobRoles(ctx context.Context, req interface{}) (interface{}, error) {
	search, ok := req.(*competency.SearchOfficeJobRoleVm)
	if !ok {
		search = &competency.SearchOfficeJobRoleVm{PageSize: 50}
	}

	q := s.officeJobRoleRepo.Query(ctx).
		Preload("Office").
		Preload("JobRole")

	if search.OfficeID != nil {
		q = q.Where("office_id = ?", *search.OfficeID)
	}
	if search.SearchString != "" {
		q = q.Joins("LEFT JOIN \"CoreSchema\".\"job_roles\" ON \"CoreSchema\".\"job_roles\".\"job_role_id\" = \"CoreSchema\".\"office_job_roles\".\"job_role_id\"").
			Joins("LEFT JOIN \"CoreSchema\".\"offices\" ON \"CoreSchema\".\"offices\".\"office_id\" = \"CoreSchema\".\"office_job_roles\".\"office_id\"").
			Where("\"CoreSchema\".\"job_roles\".\"job_role_name\" LIKE ? OR \"CoreSchema\".\"offices\".\"office_name\" LIKE ?",
				"%"+search.SearchString+"%", "%"+search.SearchString+"%")
	}

	var totalRecords int64
	q.Count(&totalRecords)

	var entities []competency.OfficeJobRole
	err := q.Offset(search.Skip).Limit(pageSize(search.PageSize)).Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get office job roles: %w", err)
	}

	vms := make([]competency.OfficeJobRoleVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.OfficeJobRoleVm{
			OfficeJobRoleID: e.OfficeJobRoleID,
			OfficeID:        e.OfficeID,
			JobRoleID:       e.JobRoleID,
			BaseAuditVm:     toBaseAuditVm(e.BaseAudit),
		}
		if e.Office != nil {
			vm.OfficeName = e.Office.OfficeName
		}
		if e.JobRole != nil {
			vm.JobRoleName = e.JobRole.JobRoleName
		}
		vms = append(vms, vm)
	}

	return &competency.OfficeJobRoleListVm{
		OfficeJobRoles: vms,
		TotalRecord:    int(totalRecords),
	}, nil
}

func (s *competencyService) SaveOfficeJobRole(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.OfficeJobRoleVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveOfficeJobRole")
	}

	var message string
	if vm.OfficeJobRoleID > 0 {
		existing, err := s.officeJobRoleRepo.GetByID(ctx, vm.OfficeJobRoleID)
		if err != nil {
			return nil, fmt.Errorf("save office job role: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Office Job Role not found"}, nil
		}
		existing.OfficeID = vm.OfficeID
		existing.JobRoleID = vm.JobRoleID
		existing.IsActive = vm.IsActive

		if err := s.officeJobRoleRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Office Job Role has been updated successfully"
	} else {
		entity := competency.OfficeJobRole{
			OfficeID:  vm.OfficeID,
			JobRoleID: vm.JobRoleID,
		}
		entity.IsActive = vm.IsActive

		if err := s.officeJobRoleRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Office Job Role has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ======================= Job Role Competencies ===============================

func (s *competencyService) GetJobRoleCompetencies(ctx context.Context, req interface{}) (interface{}, error) {
	search, ok := req.(*competency.SearchJobRoleCompetencyVm)
	if !ok {
		search = &competency.SearchJobRoleCompetencyVm{BasePagedData: competency.BasePagedData{PageSize: 50}}
	}

	q := s.jobRoleCompetencyRepo.Query(ctx).
		Preload("Competency").
		Preload("JobRole").
		Preload("Rating").
		Preload("Office")

	if search.OfficeID != nil {
		q = q.Where("office_id = ?", *search.OfficeID)
	} else if search.DivisionID != nil {
		q = q.Joins("LEFT JOIN \"CoreSchema\".\"offices\" o ON o.\"office_id\" = \"CoreSchema\".\"job_role_competencies\".\"office_id\"").
			Where("o.\"division_id\" = ?", *search.DivisionID)
	} else if search.DepartmentID != nil {
		q = q.Joins("LEFT JOIN \"CoreSchema\".\"offices\" o ON o.\"office_id\" = \"CoreSchema\".\"job_role_competencies\".\"office_id\"").
			Joins("LEFT JOIN \"CoreSchema\".\"divisions\" d ON d.\"division_id\" = o.\"division_id\"").
			Where("d.\"department_id\" = ?", *search.DepartmentID)
	}
	if search.JobRoleID != nil {
		q = q.Where("job_role_id = ?", *search.JobRoleID)
	}
	if search.SearchString != "" {
		upper := strings.ToUpper(search.SearchString)
		q = q.Joins("LEFT JOIN \"CoreSchema\".\"competencies\" c ON c.\"competency_id\" = \"CoreSchema\".\"job_role_competencies\".\"competency_id\"").
			Where("UPPER(c.\"competency_name\") LIKE ?", "%"+upper+"%")
	}

	var totalRecords int64
	q.Count(&totalRecords)

	var entities []competency.JobRoleCompetency
	err := q.Offset(search.Skip).Limit(pageSize(search.PageSize)).Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get job role competencies: %w", err)
	}

	vms := make([]competency.JobRoleCompetencyVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.JobRoleCompetencyVm{
			JobRoleCompetencyID: e.JobRoleCompetencyID,
			CompetencyID:        e.CompetencyID,
			JobRoleID:           e.JobRoleID,
			OfficeID:            e.OfficeID,
			RatingID:            e.RatingID,
			BaseAuditVm:         toBaseAuditVm(e.BaseAudit),
		}
		if e.Competency != nil {
			vm.CompetencyName = e.Competency.CompetencyName
		}
		if e.JobRole != nil {
			vm.JobRoleName = e.JobRole.JobRoleName
		}
		if e.Rating != nil {
			vm.RatingName = e.Rating.Name
		}
		if e.Office != nil {
			vm.OfficeName = e.Office.OfficeName
		}
		vms = append(vms, vm)
	}

	return &competency.PagedJobRoleCompetencyVm{
		JobRoleCompetencies: vms,
		TotalRecords:        int(totalRecords),
	}, nil
}

func (s *competencyService) SaveJobRoleCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.JobRoleCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveJobRoleCompetency")
	}

	var message string
	if vm.JobRoleCompetencyID > 0 {
		existing, err := s.jobRoleCompetencyRepo.GetByID(ctx, vm.JobRoleCompetencyID)
		if err != nil {
			return nil, fmt.Errorf("save job role competency: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Job Role Competency not found"}, nil
		}
		existing.CompetencyID = vm.CompetencyID
		existing.JobRoleID = vm.JobRoleID
		existing.OfficeID = vm.OfficeID
		existing.RatingID = vm.RatingID
		existing.IsActive = vm.IsActive

		if err := s.jobRoleCompetencyRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: fmt.Sprintf("%s already exist for %s", vm.CompetencyName, vm.JobRoleName)}, nil
		}
		message = fmt.Sprintf("%s has been added to the %s successfully", vm.CompetencyName, vm.JobRoleName)
	} else {
		entity := competency.JobRoleCompetency{
			CompetencyID: vm.CompetencyID,
			JobRoleID:    vm.JobRoleID,
			OfficeID:     vm.OfficeID,
			RatingID:     vm.RatingID,
		}
		entity.IsActive = vm.IsActive

		if err := s.jobRoleCompetencyRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: fmt.Sprintf("%s already exist for %s", vm.CompetencyName, vm.JobRoleName)}, nil
		}
		message = "Added successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ===================== Behavioral Competencies ===============================

func (s *competencyService) GetBehavioralCompetencies(ctx context.Context) (interface{}, error) {
	var entities []competency.BehavioralCompetency
	err := s.behavioralCompetencyRepo.Query(ctx).
		Preload("Competency").
		Preload("JobGradeGroup").
		Preload("Rating").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get behavioral competencies: %w", err)
	}

	vms := make([]competency.BehavioralCompetencyVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.BehavioralCompetencyVm{
			BehavioralCompetencyID: e.BehavioralCompetencyID,
			CompetencyID:           e.CompetencyID,
			JobGradeGroupID:        e.JobGradeGroupID,
			RatingID:               e.RatingID,
			BaseAuditVm:            toBaseAuditVm(e.BaseAudit),
		}
		if e.Competency != nil {
			vm.CompetencyName = e.Competency.CompetencyName
		}
		if e.Rating != nil {
			vm.RatingName = e.Rating.Name
		}
		if e.JobGradeGroup != nil {
			vm.JobGradeGroupName = e.JobGradeGroup.GroupName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) GetBehavioralCompetenciesByGradeName(ctx context.Context, gradeName string) (interface{}, error) {
	var entities []competency.BehavioralCompetency
	q := s.behavioralCompetencyRepo.Query(ctx).
		Preload("Competency").
		Preload("JobGradeGroup").
		Preload("Rating")

	if gradeName != "" {
		q = q.Joins("LEFT JOIN \"CoreSchema\".\"job_grade_groups\" ON \"CoreSchema\".\"job_grade_groups\".\"job_grade_group_id\" = \"CoreSchema\".\"behavioral_competencies\".\"job_grade_group_id\"").
			Where("\"CoreSchema\".\"job_grade_groups\".\"group_name\" = ?", gradeName)
	}

	if err := q.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("get behavioral competencies by grade: %w", err)
	}

	vms := make([]competency.BehavioralCompetencyVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.BehavioralCompetencyVm{
			BehavioralCompetencyID: e.BehavioralCompetencyID,
			CompetencyID:           e.CompetencyID,
			JobGradeGroupID:        e.JobGradeGroupID,
			RatingID:               e.RatingID,
			BaseAuditVm:            toBaseAuditVm(e.BaseAudit),
		}
		if e.Competency != nil {
			vm.CompetencyName = e.Competency.CompetencyName
		}
		if e.Rating != nil {
			vm.RatingName = e.Rating.Name
		}
		if e.JobGradeGroup != nil {
			vm.JobGradeGroupName = e.JobGradeGroup.GroupName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) SaveBehavioralCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.BehavioralCompetencyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveBehavioralCompetency")
	}

	var message string
	if vm.BehavioralCompetencyID > 0 {
		existing, err := s.behavioralCompetencyRepo.GetByID(ctx, vm.BehavioralCompetencyID)
		if err != nil {
			return nil, fmt.Errorf("save behavioral competency: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Behavioral Competency not found"}, nil
		}
		existing.CompetencyID = vm.CompetencyID
		existing.JobGradeGroupID = vm.JobGradeGroupID
		existing.RatingID = vm.RatingID
		existing.IsActive = vm.IsActive

		if err := s.behavioralCompetencyRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: fmt.Sprintf("%s already exist for %s", vm.CompetencyName, vm.JobGradeGroupName)}, nil
		}
		message = fmt.Sprintf("%s has been added to the %s successfully", vm.CompetencyName, vm.JobGradeGroupName)
	} else {
		entity := competency.BehavioralCompetency{
			CompetencyID:    vm.CompetencyID,
			JobGradeGroupID: vm.JobGradeGroupID,
			RatingID:        vm.RatingID,
		}
		entity.IsActive = vm.IsActive

		if err := s.behavioralCompetencyRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: fmt.Sprintf("%s already exist for %s", vm.CompetencyName, vm.JobGradeGroupName)}, nil
		}
		message = "Added successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ========================= Job Role Grades ===================================

func (s *competencyService) GetJobRoleGrades(ctx context.Context) (interface{}, error) {
	grades, err := s.jobRoleGradeRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get job role grades: %w", err)
	}

	vms := make([]competency.JobRoleGradeVm, 0, len(grades))
	for _, g := range grades {
		vms = append(vms, competency.JobRoleGradeVm{
			JobRoleGradeID: g.JobRoleGradeID,
			JobRoleID:      g.JobRoleID,
			GradeID:        g.GradeID,
			GradeName:      g.GradeName,
			BaseAuditVm:    toBaseAuditVm(g.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveJobRoleGrade(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.JobRoleGradeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveJobRoleGrade")
	}

	var message string
	if vm.JobRoleGradeID > 0 {
		existing, err := s.jobRoleGradeRepo.GetByID(ctx, vm.JobRoleGradeID)
		if err != nil {
			return nil, fmt.Errorf("save job role grade: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Job Role Grade not found"}, nil
		}
		existing.JobRoleID = vm.JobRoleID
		existing.GradeID = vm.GradeID
		existing.GradeName = vm.GradeName
		existing.IsActive = vm.IsActive

		if err := s.jobRoleGradeRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Job Role Grade has been updated successfully"
	} else {
		entity := competency.JobRoleGrade{
			JobRoleID: vm.JobRoleID,
			GradeID:   vm.GradeID,
			GradeName: vm.GradeName,
		}
		entity.IsActive = vm.IsActive

		if err := s.jobRoleGradeRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Job Role Grade has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ============================= Job Grades ====================================

func (s *competencyService) GetJobGrades(ctx context.Context) (interface{}, error) {
	var grades []competency.JobGrade
	err := s.jobGradeRepo.Query(ctx).Order("grade_name").Find(&grades).Error
	if err != nil {
		return nil, fmt.Errorf("get job grades: %w", err)
	}

	vms := make([]competency.JobGradeVm, 0, len(grades))
	for _, g := range grades {
		vms = append(vms, competency.JobGradeVm{
			JobGradeID:  g.JobGradeID,
			GradeCode:   g.GradeCode,
			GradeName:   g.GradeName,
			BaseAuditVm: toBaseAuditVm(g.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveJobGrade(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.JobGradeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveJobGrade")
	}

	var message string
	if vm.JobGradeID > 0 {
		existing, err := s.jobGradeRepo.GetByID(ctx, vm.JobGradeID)
		if err != nil {
			return nil, fmt.Errorf("save job grade: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Job Grade not found"}, nil
		}
		existing.GradeName = strings.ToUpper(vm.GradeName)
		existing.GradeCode = vm.GradeCode
		existing.IsActive = vm.IsActive

		if err := s.jobGradeRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Grade has been updated successfully", existing.GradeName)
	} else {
		entity := competency.JobGrade{
			GradeName: strings.ToUpper(vm.GradeName),
			GradeCode: vm.GradeCode,
		}
		entity.IsActive = vm.IsActive

		if err := s.jobGradeRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Grade has been created successfully", entity.GradeName)
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ========================= Job Grade Groups ==================================

func (s *competencyService) GetJobGradeGroups(ctx context.Context) (interface{}, error) {
	var groups []competency.JobGradeGroup
	err := s.jobGradeGroupRepo.Query(ctx).Order("\"order\"").Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("get job grade groups: %w", err)
	}

	vms := make([]competency.JobGradeGroupVm, 0, len(groups))
	for _, g := range groups {
		vms = append(vms, competency.JobGradeGroupVm{
			JobGradeGroupID: g.JobGradeGroupID,
			GroupName:        g.GroupName,
			Order:            g.Order,
			BaseAuditVm:     toBaseAuditVm(g.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveJobGradeGroup(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.JobGradeGroupVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveJobGradeGroup")
	}

	var message string
	if vm.JobGradeGroupID > 0 {
		existing, err := s.jobGradeGroupRepo.GetByID(ctx, vm.JobGradeGroupID)
		if err != nil {
			return nil, fmt.Errorf("save job grade group: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Grade Group not found"}, nil
		}
		existing.GroupName = strings.ToUpper(vm.GroupName)
		existing.Order = vm.Order
		existing.IsActive = vm.IsActive

		if err := s.jobGradeGroupRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Grade Group has been updated successfully", existing.GroupName)
	} else {
		entity := competency.JobGradeGroup{
			GroupName: strings.ToUpper(vm.GroupName),
			Order:     vm.Order,
		}
		entity.IsActive = vm.IsActive

		if err := s.jobGradeGroupRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Grade Group has been created successfully", entity.GroupName)
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ====================== Assign Job Grade Groups ==============================

func (s *competencyService) GetAssignJobGradeGroups(ctx context.Context) (interface{}, error) {
	var entities []competency.AssignJobGradeGroup
	err := s.assignJobGradeGroupRepo.Query(ctx).
		Preload("JobGrade").
		Preload("JobGradeGroup").
		Order("job_grade_id").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get assign job grade groups: %w", err)
	}

	vms := make([]competency.AssignJobGradeGroupVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.AssignJobGradeGroupVm{
			AssignJobGradeGroupID: e.AssignJobGradeGroupID,
			JobGradeGroupID:       e.JobGradeGroupID,
			JobGradeID:            e.JobGradeID,
			BaseAuditVm:           toBaseAuditVm(e.BaseAudit),
		}
		if e.JobGrade != nil {
			vm.JobGradeName = e.JobGrade.GradeName
		}
		if e.JobGradeGroup != nil {
			vm.JobGradeGroupName = e.JobGradeGroup.GroupName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) GetAssignJobGradeGroupByGradeName(ctx context.Context, gradeName string) (interface{}, error) {
	var entity competency.AssignJobGradeGroup
	err := s.assignJobGradeGroupRepo.Query(ctx).
		Joins("LEFT JOIN \"CoreSchema\".\"job_grades\" ON \"CoreSchema\".\"job_grades\".\"job_grade_id\" = \"CoreSchema\".\"assign_job_grade_groups\".\"job_grade_id\"").
		Where("\"CoreSchema\".\"job_grades\".\"grade_name\" = ?", gradeName).
		Preload("JobGrade").
		Preload("JobGradeGroup").
		First(&entity).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("get assign job grade group by grade name: %w", err)
	}

	vm := competency.AssignJobGradeGroupVm{
		AssignJobGradeGroupID: entity.AssignJobGradeGroupID,
		JobGradeGroupID:       entity.JobGradeGroupID,
		JobGradeID:            entity.JobGradeID,
		BaseAuditVm:           toBaseAuditVm(entity.BaseAudit),
	}
	if entity.JobGrade != nil {
		vm.JobGradeName = entity.JobGrade.GradeName
	}
	if entity.JobGradeGroup != nil {
		vm.JobGradeGroupName = entity.JobGradeGroup.GroupName
	}
	return &vm, nil
}

func (s *competencyService) SaveAssignJobGradeGroup(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.AssignJobGradeGroupVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveAssignJobGradeGroup")
	}

	var message string
	if vm.AssignJobGradeGroupID > 0 {
		existing, err := s.assignJobGradeGroupRepo.GetByID(ctx, vm.AssignJobGradeGroupID)
		if err != nil {
			return nil, fmt.Errorf("save assign job grade group: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Assignment not found"}, nil
		}
		existing.JobGradeID = vm.JobGradeID
		existing.JobGradeGroupID = vm.JobGradeGroupID
		existing.IsActive = vm.IsActive

		if err := s.assignJobGradeGroupRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Grade group assigned successfully"
	} else {
		entity := competency.AssignJobGradeGroup{
			JobGradeGroupID: vm.JobGradeGroupID,
			JobGradeID:      vm.JobGradeID,
		}
		entity.IsActive = vm.IsActive

		if err := s.assignJobGradeGroupRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Grade group assigned successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ================================ Ratings ====================================

func (s *competencyService) GetRatings(ctx context.Context) (interface{}, error) {
	var ratings []competency.Rating
	err := s.ratingRepo.Query(ctx).Order("value").Find(&ratings).Error
	if err != nil {
		return nil, fmt.Errorf("get ratings: %w", err)
	}

	vms := make([]competency.RatingVm, 0, len(ratings))
	for _, r := range ratings {
		vms = append(vms, competency.RatingVm{
			RatingID:    r.RatingID,
			Name:        r.Name,
			Value:       r.Value,
			BaseAuditVm: toBaseAuditVm(r.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveRating(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.RatingVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveRating")
	}

	var message string
	if vm.RatingID > 0 {
		existing, err := s.ratingRepo.GetByID(ctx, vm.RatingID)
		if err != nil {
			return nil, fmt.Errorf("save rating: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Rating not found"}, nil
		}
		existing.Name = vm.Name
		existing.Value = vm.Value
		existing.IsActive = vm.IsActive

		if err := s.ratingRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Rating has been updated successfully"
	} else {
		entity := competency.Rating{
			Name:  vm.Name,
			Value: vm.Value,
		}
		entity.IsActive = vm.IsActive

		if err := s.ratingRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Rating has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ============================ Review Periods =================================

func (s *competencyService) GetReviewPeriods(ctx context.Context) (interface{}, error) {
	var entities []competency.ReviewPeriod
	err := s.reviewPeriodRepo.Query(ctx).
		Preload("BankYear").
		Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("get review periods: %w", err)
	}

	vms := make([]competency.ReviewPeriodVm, 0, len(entities))
	for _, e := range entities {
		vm := competency.ReviewPeriodVm{
			ReviewPeriodID: e.ReviewPeriodID,
			BankYearID:     e.BankYearID,
			Name:           e.Name,
			StartDate:      e.StartDate,
			EndDate:        e.EndDate,
			IsApproved:     e.IsApproved,
			ApprovedBy:     e.ApprovedBy,
			DateApproved:   e.DateApproved,
			BaseAuditVm:    toBaseAuditVm(e.BaseWorkFlowData.BaseAudit),
		}
		if e.BankYear != nil {
			vm.BankYearName = e.BankYear.YearName
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

func (s *competencyService) GetCurrentReviewPeriod(ctx context.Context) (interface{}, error) {
	var entity competency.ReviewPeriod
	err := s.reviewPeriodRepo.Query(ctx).
		Where("is_approved = ? AND is_active = ?", true, true).
		Preload("BankYear").
		First(&entity).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("get current review period: %w", err)
	}

	vm := competency.ReviewPeriodVm{
		ReviewPeriodID: entity.ReviewPeriodID,
		BankYearID:     entity.BankYearID,
		Name:           entity.Name,
		StartDate:      entity.StartDate,
		EndDate:        entity.EndDate,
		IsApproved:     entity.IsApproved,
		BaseAuditVm:    toBaseAuditVm(entity.BaseWorkFlowData.BaseAudit),
	}
	if entity.BankYear != nil {
		vm.BankYearName = entity.BankYear.YearName
	}
	return &vm, nil
}

func (s *competencyService) SaveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.ReviewPeriodVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveReviewPeriod")
	}

	// Deactivate currently active periods if the new/updated one is active
	if vm.IsActive {
		s.db.WithContext(ctx).
			Model(&competency.ReviewPeriod{}).
			Where("is_active = ? AND soft_deleted = ?", true, false).
			Update("is_active", false)
	}

	var message string
	if vm.ReviewPeriodID > 0 {
		existing, err := s.reviewPeriodRepo.GetByID(ctx, vm.ReviewPeriodID)
		if err != nil {
			return nil, fmt.Errorf("save review period: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Review Period not found"}, nil
		}
		existing.Name = vm.Name
		existing.BankYearID = vm.BankYearID
		existing.StartDate = vm.StartDate
		existing.EndDate = vm.EndDate
		existing.IsActive = vm.IsActive
		existing.IsApproved = vm.IsApproved
		existing.ApprovedBy = vm.ApprovedBy
		existing.DateApproved = vm.DateApproved

		if err := s.reviewPeriodRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Review Period has been updated successfully"
	} else {
		entity := competency.ReviewPeriod{
			Name:       vm.Name,
			BankYearID: vm.BankYearID,
			StartDate:  vm.StartDate,
			EndDate:    vm.EndDate,
		}
		entity.IsActive = vm.IsActive
		entity.IsApproved = vm.IsApproved
		entity.ApprovedBy = vm.ApprovedBy
		entity.DateApproved = vm.DateApproved

		if err := s.reviewPeriodRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Review Period has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

func (s *competencyService) ApproveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.ReviewPeriodVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for ApproveReviewPeriod")
	}

	existing, err := s.reviewPeriodRepo.GetByID(ctx, vm.ReviewPeriodID)
	if err != nil {
		return nil, fmt.Errorf("approve review period: %w", err)
	}
	if existing == nil {
		return &responseVm{IsSuccess: false, Message: "Record with Id not found"}, nil
	}

	existing.IsActive = vm.IsActive
	existing.IsApproved = vm.IsApproved
	existing.ApprovedBy = vm.ApprovedBy
	existing.DateApproved = vm.DateApproved

	if err := s.reviewPeriodRepo.Update(ctx, existing); err != nil {
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}

	return &responseVm{IsSuccess: true, Message: "Review Period has been approved successfully"}, nil
}

// ============================= Review Types ==================================

func (s *competencyService) GetReviewTypes(ctx context.Context) (interface{}, error) {
	types, err := s.reviewTypeRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get review types: %w", err)
	}

	vms := make([]competency.ReviewTypeVm, 0, len(types))
	for _, t := range types {
		vms = append(vms, competency.ReviewTypeVm{
			ReviewTypeID:   t.ReviewTypeID,
			ReviewTypeName: t.ReviewTypeName,
			BaseAuditVm:    toBaseAuditVm(t.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveReviewType(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.ReviewTypeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveReviewType")
	}

	var message string
	if vm.ReviewTypeID > 0 {
		existing, err := s.reviewTypeRepo.GetByID(ctx, vm.ReviewTypeID)
		if err != nil {
			return nil, fmt.Errorf("save review type: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Review Type not found"}, nil
		}
		existing.ReviewTypeName = vm.ReviewTypeName
		existing.IsActive = vm.IsActive

		if err := s.reviewTypeRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Review Type has been updated successfully"
	} else {
		entity := competency.ReviewType{
			ReviewTypeName: vm.ReviewTypeName,
		}
		entity.IsActive = vm.IsActive

		if err := s.reviewTypeRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Review Type has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ============================== Bank Years ===================================

func (s *competencyService) GetBankYears(ctx context.Context) (interface{}, error) {
	var years []identity.BankYear
	err := s.bankYearRepo.Query(ctx).Find(&years).Error
	if err != nil {
		return nil, fmt.Errorf("get bank years: %w", err)
	}

	vms := make([]competency.BankYearVm, 0, len(years))
	for _, y := range years {
		vms = append(vms, competency.BankYearVm{
			BankYearID: y.BankYearID,
			YearName:   y.YearName,
		})
	}
	return vms, nil
}

func (s *competencyService) SaveBankYear(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.BankYearVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveBankYear")
	}

	var message string
	if vm.BankYearID > 0 {
		existing, err := s.bankYearRepo.GetByID(ctx, vm.BankYearID)
		if err != nil {
			return nil, fmt.Errorf("save bank year: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Bank Year not found"}, nil
		}
		existing.YearName = vm.YearName

		if err := s.bankYearRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Bank Year has been Updated Successfully", existing.YearName)
	} else {
		entity := identity.BankYear{
			YearName: vm.YearName,
		}
		entity.IsActive = vm.IsActive

		if err := s.bankYearRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = fmt.Sprintf("%s Bank Year has been Created Successfully", entity.YearName)
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ============================ Training Types =================================

func (s *competencyService) GetTrainingTypes(ctx context.Context, isActive *bool) (interface{}, error) {
	q := s.trainingTypeRepo.Query(ctx)
	if isActive != nil {
		q = q.Where("is_active = ?", *isActive)
	}

	var types []competency.TrainingType
	if err := q.Find(&types).Error; err != nil {
		return nil, fmt.Errorf("get training types: %w", err)
	}

	vms := make([]competency.TrainingTypeVm, 0, len(types))
	for _, t := range types {
		vms = append(vms, competency.TrainingTypeVm{
			TrainingTypeID:   t.TrainingTypeID,
			TrainingTypeName: t.TrainingTypeName,
			BaseAuditVm:      toBaseAuditVm(t.BaseAudit),
		})
	}
	return vms, nil
}

func (s *competencyService) SaveTrainingType(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.TrainingTypeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveTrainingType")
	}

	var message string
	if vm.TrainingTypeID > 0 {
		existing, err := s.trainingTypeRepo.GetByID(ctx, vm.TrainingTypeID)
		if err != nil {
			return nil, fmt.Errorf("save training type: %w", err)
		}
		if existing == nil {
			return &responseVm{IsSuccess: false, Message: "Training Type not found"}, nil
		}
		existing.TrainingTypeName = vm.TrainingTypeName
		existing.IsActive = vm.IsActive

		if err := s.trainingTypeRepo.Update(ctx, existing); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Development Intervention Type has been updated successfully"
	} else {
		entity := competency.TrainingType{
			TrainingTypeName: vm.TrainingTypeName,
		}
		entity.IsActive = vm.IsActive

		if err := s.trainingTypeRepo.Create(ctx, &entity); err != nil {
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
		message = "Development Intervention Type has been created successfully"
	}

	return &responseVm{IsSuccess: true, Message: message}, nil
}

// ================ Population / Calculation (ReviewAgentService) ==============
// These methods delegate to the internal reviewAgentService which contains the
// full implementation ported from the .NET ReviewAgentService (~1,415 lines).

func (s *competencyService) PopulateAllReviews(ctx context.Context) (interface{}, error) {
	s.log.Info().Msg("PopulateAllReviews: populating reviews for all employees")
	if err := s.reviewAgent.PopulateAllEmployeeReviews(ctx); err != nil {
		s.log.Error().Err(err).Msg("PopulateAllReviews failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: "All employee reviews populated successfully"}, nil
}

func (s *competencyService) PopulateOfficeReviews(ctx context.Context, officeId int) (interface{}, error) {
	s.log.Info().Int("officeId", officeId).Msg("PopulateOfficeReviews: populating reviews for office employees")
	if err := s.reviewAgent.PopulateOfficeEmployeeReviews(ctx, officeId); err != nil {
		s.log.Error().Err(err).Int("officeId", officeId).Msg("PopulateOfficeReviews failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: fmt.Sprintf("Reviews for office %d populated successfully", officeId)}, nil
}

func (s *competencyService) PopulateDivisionReviews(ctx context.Context, divisionId int) (interface{}, error) {
	s.log.Info().Int("divisionId", divisionId).Msg("PopulateDivisionReviews: populating reviews for division employees")
	if err := s.reviewAgent.PopulateDivisionEmployeeReviews(ctx, divisionId); err != nil {
		s.log.Error().Err(err).Int("divisionId", divisionId).Msg("PopulateDivisionReviews failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: fmt.Sprintf("Reviews for division %d populated successfully", divisionId)}, nil
}

func (s *competencyService) PopulateDepartmentReviews(ctx context.Context, departmentId int) (interface{}, error) {
	s.log.Info().Int("departmentId", departmentId).Msg("PopulateDepartmentReviews: populating reviews for department employees")
	if err := s.reviewAgent.PopulateDepartmentEmployeeReviews(ctx, departmentId); err != nil {
		s.log.Error().Err(err).Int("departmentId", departmentId).Msg("PopulateDepartmentReviews failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: fmt.Sprintf("Reviews for department %d populated successfully", departmentId)}, nil
}

func (s *competencyService) PopulateReviewsByEmployeeId(ctx context.Context, employeeNumber string) (interface{}, error) {
	s.log.Info().Str("employeeNumber", employeeNumber).Msg("PopulateReviewsByEmployeeId: populating reviews for single employee")
	if err := s.reviewAgent.PopulateReviewsForEmployee(ctx, employeeNumber); err != nil {
		s.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("PopulateReviewsByEmployeeId failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: fmt.Sprintf("Reviews for employee %s populated successfully", employeeNumber)}, nil
}

func (s *competencyService) CalculateReviews(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*competency.CalculateReviewProfileVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CalculateReviews")
	}

	s.log.Info().
		Str("employeeNumber", vm.EmployeeNumber).
		Int("reviewPeriodId", vm.ReviewPeriodID).
		Bool("isTechnical", vm.IsTechnical).
		Msg("CalculateReviews: computing review averages")

	if vm.IsTechnical {
		if err := s.reviewAgent.CalculateTechnicalReviewAverage(ctx, vm.EmployeeNumber, vm.ReviewPeriodID); err != nil {
			s.log.Error().Err(err).Msg("CalculateTechnicalReviewAverage failed")
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
	} else {
		if err := s.reviewAgent.CalculateBehavioralReviewAverage(ctx, vm.EmployeeNumber, vm.ReviewPeriodID); err != nil {
			s.log.Error().Err(err).Msg("CalculateBehavioralReviewAverage failed")
			return &responseVm{IsSuccess: false, Message: err.Error()}, nil
		}
	}
	return &responseVm{IsSuccess: true, Message: fmt.Sprintf("Reviews calculated for %s successfully", vm.EmployeeNumber)}, nil
}

func (s *competencyService) RecalculateReviewsProfiles(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Info().Msg("RecalculateReviewsProfiles: recalculating all employee profiles")
	if err := s.reviewAgent.CalculateReviewsProfileForAllEmployees(ctx); err != nil {
		s.log.Error().Err(err).Msg("RecalculateReviewsProfiles failed")
		return &responseVm{IsSuccess: false, Message: err.Error()}, nil
	}
	return &responseVm{IsSuccess: true, Message: "All review profiles recalculated successfully"}, nil
}

// ========================= Email / Sync ======================================

func (s *competencyService) EmailService(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Info().Msg("EmailService: processing email request")
	// TODO: integrate with the Email service to send competency-related notifications.
	return &responseVm{IsSuccess: true, Message: "Email request processed"}, nil
}

func (s *competencyService) SyncJobRoleUpdateSOA(ctx context.Context, req interface{}) (interface{}, error) {
	s.log.Info().Msg("SyncJobRoleUpdateSOA: syncing job role update to SOA ERP")
	// TODO: post job role update to the SOA API endpoint configured in cfg.
	return &responseVm{IsSuccess: true, Message: "SOA sync request processed"}, nil
}

// ===========================================================================
// Helper functions
// ===========================================================================

// toBaseAuditVm converts domain.BaseAudit to the DTO BaseAuditVm.
func toBaseAuditVm(a domain.BaseAudit) competency.BaseAuditVm {
	return competency.BaseAuditVm{
		CreatedBy:   a.CreatedBy,
		DateCreated: a.DateCreated,
		IsActive:    a.IsActive,
		Status:      a.Status,
		DateUpdated: a.DateUpdated,
		UpdatedBy:   a.UpdatedBy,
	}
}

// pageSize returns a safe page size, defaulting to 50 if zero.
func pageSize(ps int) int {
	if ps <= 0 {
		return 50
	}
	return ps
}

// mapReviewsToVms converts a slice of CompetencyReview entities to VMs.
func mapReviewsToVms(entities []competency.CompetencyReview) []competency.CompetencyReviewVm {
	vms := make([]competency.CompetencyReviewVm, 0, len(entities))
	for _, e := range entities {
		vms = append(vms, mapSingleReviewToVm(e))
	}
	return vms
}

func mapSingleReviewToVm(e competency.CompetencyReview) competency.CompetencyReviewVm {
	vm := competency.CompetencyReviewVm{
		CompetencyReviewID: e.CompetencyReviewID,
		EmployeeNumber:     e.EmployeeNumber,
		CompetencyID:       e.CompetencyID,
		ReviewDate:         e.ReviewDate,
		ReviewerID:         e.ReviewerID,
		ReviewTypeID:       e.ReviewTypeID,
		ReviewPeriodID:     e.ReviewPeriodID,
		ReviewerName:       e.ReviewerName,
		ActualRatingID:     e.ActualRatingID,
		ActualRatingName:   e.ActualRatingName,
		ActualRatingValue:  e.ActualRatingValue,
		ExpectedRatingID:   e.ExpectedRatingID,
		IsTechnical:        e.IsTechnical,
		EmployeeName:       e.EmployeeName,
		EmployeeDepartment: e.EmployeeDepartment,
		EmployeeInitial:    e.EmployeeInitial,
		EmployeeGrade:      e.EmployeeGrade,
		BaseAuditVm:        toBaseAuditVm(e.BaseAudit),
	}
	if e.Competency != nil {
		vm.CompetencyName = e.Competency.CompetencyName
		if e.Competency.CompetencyCategory != nil {
			vm.CompetencyCategoryName = e.Competency.CompetencyCategory.CategoryName
		}
	}
	if e.ReviewType != nil {
		vm.ReviewTypeName = e.ReviewType.ReviewTypeName
	}
	if e.ReviewPeriod != nil {
		vm.ReviewPeriodName = e.ReviewPeriod.Name
	}
	if e.ExpectedRating != nil {
		vm.ExpectedRatingName = e.ExpectedRating.Name
		vm.ExpectedRatingValue = e.ExpectedRating.Value
	}
	return vm
}

// mapProfilesToVms converts review profile entities to VMs with development plan counts.
func mapProfilesToVms(entities []competency.CompetencyReviewProfile) []competency.CompetencyReviewProfileVm {
	vms := make([]competency.CompetencyReviewProfileVm, 0, len(entities))
	for _, e := range entities {
		progressCount := 0
		completedCount := 0
		for _, dp := range e.DevelopmentPlans {
			if strings.EqualFold(dp.TaskStatus, "InProgress") {
				progressCount++
			}
			if strings.EqualFold(dp.TaskStatus, "Completed") {
				completedCount++
			}
		}

		vms = append(vms, competency.CompetencyReviewProfileVm{
			CompetencyReviewProfileID: e.CompetencyReviewProfileID,
			EmployeeNumber:            e.EmployeeNumber,
			CompetencyID:              e.CompetencyID,
			CompetencyName:            e.CompetencyName,
			ExpectedRatingID:          e.ExpectedRatingID,
			ExpectedRatingName:        e.ExpectedRatingName,
			ExpectedRatingValue:       e.ExpectedRatingValue,
			ReviewPeriodID:            e.ReviewPeriodID,
			ReviewPeriodName:          e.ReviewPeriodName,
			AverageRatingID:           e.AverageRatingID,
			AverageRatingName:         e.AverageRatingName,
			AverageRatingValue:        e.AverageRatingValue,
			AverageScore:              e.AverageScore,
			CompetencyCategoryName:    e.CompetencyCategoryName,
			EmployeeFullName:          e.EmployeeName,
			OfficeID:                  e.OfficeID,
			OfficeName:                e.OfficeName,
			DivisionID:                e.DivisionID,
			DivisionName:              e.DivisionName,
			DepartmentID:              e.DepartmentID,
			DepartmentName:            e.DepartmentName,
			GradeName:                 e.GradeName,
			JobRoleName:               e.JobRoleName,
			JobRoleID:                 e.JobRoleID,
			NumberOfDevelopmentPlans:   len(e.DevelopmentPlans),
			ProgressCount:             progressCount,
			CompletedCount:            completedCount,
			BaseAuditVm:               toBaseAuditVm(e.BaseAudit),
		})
	}
	return vms
}

// applyOrgFilter adds review_period_id and office/division/department filters.
func applyOrgFilter(q *gorm.DB, reviewPeriodId, officeId, divisionId, departmentId *int) *gorm.DB {
	if reviewPeriodId != nil {
		q = q.Where("review_period_id = ?", *reviewPeriodId)
	}
	if officeId != nil && *officeId > 0 {
		q = q.Where("office_id = ?", fmt.Sprintf("%d", *officeId))
	} else if divisionId != nil && *divisionId > 0 {
		q = q.Where("division_id = ?", fmt.Sprintf("%d", *divisionId))
	} else if departmentId != nil && *departmentId > 0 {
		q = q.Where("department_id = ?", fmt.Sprintf("%d", *departmentId))
	}
	return q
}

// buildMatrixResult constructs the competency matrix overview from profiles.
func buildMatrixResult(profiles []competency.CompetencyReviewProfile) interface{} {
	type matrixDetail struct {
		CompetencyName     string  `json:"competencyName"`
		AverageScore       int     `json:"averageScore"`
		ExpectedRatingValue int    `json:"expectedRatingValue"`
	}
	type matrixProfile struct {
		EmployeeId             string         `json:"employeeId"`
		EmployeeName           string         `json:"employeeName"`
		OfficeName             string         `json:"officeName"`
		DivisionName           string         `json:"divisionName"`
		DepartmentName         string         `json:"departmentName"`
		Grade                  string         `json:"grade"`
		Position               string         `json:"position"`
		GapCount               int            `json:"gapCount"`
		NoOfCompetent          int            `json:"noOfCompetent"`
		NoOfCompetencies       int            `json:"noOfCompetencies"`
		OverallAverage         float64        `json:"overallAverage"`
		CompetencyMatrixDetails []matrixDetail `json:"competencyMatrixDetails"`
	}
	type matrixOverview struct {
		CompetencyNames               []string        `json:"competencyNames"`
		CompetencyMatrixReviewProfiles []matrixProfile `json:"competencyMatrixReviewProfiles"`
	}

	result := matrixOverview{
		CompetencyNames:               []string{},
		CompetencyMatrixReviewProfiles: []matrixProfile{},
	}

	nameSet := make(map[string]bool)
	for _, p := range profiles {
		if !nameSet[p.CompetencyName] {
			nameSet[p.CompetencyName] = true
			result.CompetencyNames = append(result.CompetencyNames, p.CompetencyName)
		}
	}

	employeeMap := make(map[string][]competency.CompetencyReviewProfile)
	for _, p := range profiles {
		employeeMap[p.EmployeeNumber] = append(employeeMap[p.EmployeeNumber], p)
	}

	for _, empProfiles := range employeeMap {
		if len(empProfiles) == 0 {
			continue
		}
		first := empProfiles[0]
		gapCount := 0
		noCompetent := 0
		var totalAvg float64
		details := make([]matrixDetail, 0, len(empProfiles))

		for _, ep := range empProfiles {
			if ep.HaveGap {
				gapCount++
			} else {
				noCompetent++
			}
			totalAvg += float64(ep.AverageRatingValue)
			details = append(details, matrixDetail{
				CompetencyName:      ep.CompetencyName,
				AverageScore:        ep.AverageRatingValue,
				ExpectedRatingValue: ep.ExpectedRatingValue,
			})
		}

		avg := 0.0
		if len(empProfiles) > 0 {
			avg = totalAvg / float64(len(empProfiles))
		}

		result.CompetencyMatrixReviewProfiles = append(result.CompetencyMatrixReviewProfiles, matrixProfile{
			EmployeeId:              first.EmployeeNumber,
			EmployeeName:            first.EmployeeName,
			OfficeName:              first.OfficeName,
			DivisionName:            first.DivisionName,
			DepartmentName:          first.DepartmentName,
			Grade:                   first.GradeName,
			Position:                first.JobRoleName,
			GapCount:                gapCount,
			NoOfCompetent:           noCompetent,
			NoOfCompetencies:        len(empProfiles),
			OverallAverage:          avg,
			CompetencyMatrixDetails: details,
		})
	}

	return result
}

// uniqueStrings extracts unique string values from a slice using a selector.
func uniqueStrings(profiles []competency.CompetencyReviewProfile, selector func(competency.CompetencyReviewProfile) string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, p := range profiles {
		v := selector(p)
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// filterProfiles filters profiles using a predicate.
func filterProfiles(profiles []competency.CompetencyReviewProfile, pred func(competency.CompetencyReviewProfile) bool) []competency.CompetencyReviewProfile {
	var result []competency.CompetencyReviewProfile
	for _, p := range profiles {
		if pred(p) {
			result = append(result, p)
		}
	}
	return result
}

// avgField computes the average of an int field extracted from profiles.
func avgField(profiles []competency.CompetencyReviewProfile, field func(competency.CompetencyReviewProfile) int) float64 {
	if len(profiles) == 0 {
		return 0
	}
	var sum float64
	for _, p := range profiles {
		sum += float64(field(p))
	}
	return sum / float64(len(profiles))
}

// maxField returns the maximum value of an int field extracted from profiles.
func maxField(profiles []competency.CompetencyReviewProfile, field func(competency.CompetencyReviewProfile) int) int {
	if len(profiles) == 0 {
		return 0
	}
	m := field(profiles[0])
	for _, p := range profiles[1:] {
		if v := field(p); v > m {
			m = v
		}
	}
	return m
}

// minField returns the minimum value of an int field extracted from profiles.
func minField(profiles []competency.CompetencyReviewProfile, field func(competency.CompetencyReviewProfile) int) int {
	if len(profiles) == 0 {
		return 0
	}
	m := field(profiles[0])
	for _, p := range profiles[1:] {
		if v := field(p); v < m {
			m = v
		}
	}
	return m
}
