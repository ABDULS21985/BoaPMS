package service

import (
	"context"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
)

// ===========================================================================
// ListVm methods — mirrors .NET SMDService #region ListVm
// ===========================================================================

// GetStrategies returns all strategies with base64 image prefix.
func (s *strategyService) GetStrategies(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	strategies, err := s.strategyRepo.GetAll(ctx)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.StrategyVm
	for _, st := range strategies {
		vm := performance.StrategyVm{
			StrategyID:  st.StrategyID,
			Name:        st.Name,
			Description: st.Description,
			BankYearID:  st.BankYearID,
			StartDate:   st.StartDate,
			EndDate:     st.EndDate,
		}
		vm.ID = st.ID
		vm.Status = st.Status
		vm.IsActive = st.IsActive
		vm.CreatedBy = st.CreatedBy
		vm.CreatedAt = st.CreatedAt
		vm.UpdatedAt = st.UpdatedAt
		vm.UpdatedBy = st.UpdatedBy
		vm.IsApproved = st.IsApproved
		vm.ApprovedBy = st.ApprovedBy
		vm.DateApproved = st.DateApproved
		vm.IsRejected = st.IsRejected
		vm.RejectedBy = st.RejectedBy
		vm.RejectionReason = st.RejectionReason
		vm.DateRejected = st.DateRejected
		if st.FileImage != "" {
			vm.ImageFile = "data:image/jpeg;base64," + st.FileImage
		}
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetStrategicThemes returns all strategic themes with strategy name.
func (s *strategyService) GetStrategicThemes(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	themes, err := s.themeRepo.GetAllIncluding(ctx, "Strategy")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.StrategicThemeVm
	for _, t := range themes {
		vm := performance.StrategicThemeVm{
			StrategicThemeID: t.StrategicThemeID,
			Name:             t.Name,
			Description:      t.Description,
			StrategyID:       t.StrategyID,
		}
		vm.ID = t.ID
		vm.Status = t.Status
		vm.IsActive = t.IsActive
		vm.CreatedBy = t.CreatedBy
		vm.CreatedAt = t.CreatedAt
		vm.UpdatedAt = t.UpdatedAt
		vm.IsApproved = t.IsApproved
		vm.ApprovedBy = t.ApprovedBy
		vm.DateApproved = t.DateApproved
		vm.IsRejected = t.IsRejected
		vm.RejectedBy = t.RejectedBy
		vm.RejectionReason = t.RejectionReason
		vm.DateRejected = t.DateRejected
		if t.FileImage != "" {
			vm.ImageFile = "data:image/jpeg;base64," + t.FileImage
		}
		if t.Strategy != nil {
			vm.StrategyName = t.Strategy.Name
		}
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetStrategicThemesById returns strategic themes filtered by strategyID.
func (s *strategyService) GetStrategicThemesById(ctx context.Context, strategyID string) (interface{}, error) {
	resp := &performance.GenericListVm{}
	themes, err := s.themeRepo.WhereWithPreload(ctx,
		[]string{"Strategy"},
		"strategy_id = ?", strategyID,
	)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.StrategicThemeVm
	for _, t := range themes {
		vm := performance.StrategicThemeVm{
			StrategicThemeID: t.StrategicThemeID,
			Name:             t.Name,
			Description:      t.Description,
			StrategyID:       t.StrategyID,
		}
		vm.ID = t.ID
		vm.Status = t.Status
		vm.IsActive = t.IsActive
		if t.FileImage != "" {
			vm.ImageFile = "data:image/jpeg;base64," + t.FileImage
		}
		if t.Strategy != nil {
			vm.StrategyName = t.Strategy.Name
		}
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetEvaluationOptions returns all evaluation options.
func (s *strategyService) GetEvaluationOptions(ctx context.Context) (interface{}, error) {
	resp := &performance.EvaluationOptionResponseVm{}
	options, err := s.evalOptionRepo.GetAll(ctx)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	for _, opt := range options {
		vm := performance.EvaluationOptionVm{
			EvaluationOptionID: opt.EvaluationOptionID,
			Name:               opt.Name,
			Description:        opt.Description,
			Score:              opt.Score,
			EvaluationType:     int(opt.EvaluationType),
			RecordStatus:       int(opt.RecordStatus),
		}
		vm.ID = opt.ID
		vm.Status = opt.Status
		resp.EvaluationOptions = append(resp.EvaluationOptions, vm)
	}

	resp.TotalRecords = len(resp.EvaluationOptions)
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetObjectiveWorkProductDefinitions returns work product definitions for an objective.
func (s *strategyService) GetObjectiveWorkProductDefinitions(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error) {
	resp := &performance.WorkProductDefinitionResponseVm{}
	level := enums.ObjectiveLevel(objectiveLevel)

	owp, err := s.getWorkProductObjectiveByID(ctx, objectiveID, level)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	if owp == nil {
		resp.HasError = true
		resp.Message = "record not found for objective (" + objectiveID + ") in " + objectiveLevelName(level) + " objectives"
		return resp, nil
	}

	for _, wp := range owp.workProducts {
		vm := performance.WorkProductDefinitionVm{
			WorkProductDefinitionID: wp.WorkProductDefinitionID,
			ReferenceNo:             wp.ReferenceNo,
			Name:                    wp.Name,
			Description:             wp.Description,
			Deliverables:            wp.Deliverables,
			ObjectiveID:             wp.ObjectiveID,
			ObjectiveLevel:          wp.ObjectiveLevel,
		}
		resp.WorkProductDefinitions = append(resp.WorkProductDefinitions, vm)
	}

	resp.TotalRecords = len(resp.WorkProductDefinitions)
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetConsolidatedObjectives returns objectives from all levels consolidated.
func (s *strategyService) GetConsolidatedObjectives(ctx context.Context) (interface{}, error) {
	resp := &performance.ConsolidatedObjectiveListResponseVm{}

	// Enterprise objectives
	entObjs, err := s.entObjRepo.GetAllIncluding(ctx, "Strategy")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	for _, o := range entObjs {
		sbuName := ""
		if o.Strategy != nil {
			sbuName = o.Strategy.Name
		}
		resp.Objectives = append(resp.Objectives, performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{ID: o.ID, RecordStatus: o.RecordStatus},
				ObjectiveID:    o.EnterpriseObjectiveID,
				ObjectiveLevel: objectiveLevelName(enums.ObjectiveLevelEnterprise),
				Name:           o.Name,
				Description:    o.Description,
				Kpi:            o.Kpi,
				Target:         o.Target,
				SBUName:        sbuName,
			},
		})
	}

	// Department objectives
	deptObjs, err := s.deptObjRepo.GetAllIncluding(ctx, "Department")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	for _, o := range deptObjs {
		sbuName := ""
		if o.Department != nil {
			sbuName = o.Department.DepartmentName
		}
		resp.Objectives = append(resp.Objectives, performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{ID: o.ID, RecordStatus: o.RecordStatus},
				ObjectiveID:    o.DepartmentObjectiveID,
				ObjectiveLevel: objectiveLevelName(enums.ObjectiveLevelDepartment),
				Name:           o.Name,
				Description:    o.Description,
				Kpi:            o.Kpi,
				Target:         o.Target,
				SBUName:        sbuName,
			},
		})
	}

	// Division objectives
	divObjs, err := s.divObjRepo.GetAllIncluding(ctx, "Division")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	for _, o := range divObjs {
		sbuName := ""
		if o.Division != nil {
			sbuName = o.Division.DivisionName
		}
		resp.Objectives = append(resp.Objectives, performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{ID: o.ID, RecordStatus: o.RecordStatus},
				ObjectiveID:    o.DivisionObjectiveID,
				ObjectiveLevel: objectiveLevelName(enums.ObjectiveLevelDivision),
				Name:           o.Name,
				Description:    o.Description,
				Kpi:            o.Kpi,
				Target:         o.Target,
				SBUName:        sbuName,
			},
		})
	}

	// Office objectives
	offObjs, err := s.offObjRepo.GetAllIncluding(ctx, "Office")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	for _, o := range offObjs {
		sbuName := ""
		if o.Office != nil {
			sbuName = o.Office.OfficeName
		}
		resp.Objectives = append(resp.Objectives, performance.ConsolidatedObjectiveVm{
			ObjectiveBaseVm: performance.ObjectiveBaseVm{
				BaseWorkFlowVm: performance.BaseWorkFlowVm{ID: o.ID, RecordStatus: o.RecordStatus},
				ObjectiveID:    o.OfficeObjectiveID,
				ObjectiveLevel: objectiveLevelName(enums.ObjectiveLevelOffice),
				Name:           o.Name,
				Description:    o.Description,
				Kpi:            o.Kpi,
				Target:         o.Target,
				SBUName:        sbuName,
			},
		})
	}

	resp.TotalRecords = len(resp.Objectives)
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetConsolidatedObjectivesPaginated returns paginated consolidated objectives with filters.
func (s *strategyService) GetConsolidatedObjectivesPaginated(ctx context.Context, params interface{}) (interface{}, error) {
	vm, ok := params.(*performance.SearchObjectiveVm)
	if !ok {
		return &performance.GenericResponseVm{Message: "invalid parameters", IsSuccess: false}, nil
	}

	// Get all consolidated objectives first
	allResult, err := s.GetConsolidatedObjectives(ctx)
	if err != nil {
		return &performance.GenericResponseVm{Message: err.Error(), IsSuccess: false}, nil
	}
	consolidated, ok := allResult.(*performance.ConsolidatedObjectiveListResponseVm)
	if !ok {
		return &performance.GenericResponseVm{Message: "unexpected result type", IsSuccess: false}, nil
	}

	objectives := consolidated.Objectives

	// Apply status filter
	if vm.Status != "" {
		var filtered []performance.ConsolidatedObjectiveVm
		for _, o := range objectives {
			if strings.EqualFold(o.RecordStatus, vm.Status) {
				filtered = append(filtered, o)
			}
		}
		objectives = filtered
	}

	// Apply target reference filter
	if vm.TargetReference != "" {
		search := strings.ToLower(strings.TrimSpace(vm.TargetReference))
		var filtered []performance.ConsolidatedObjectiveVm
		for _, o := range objectives {
			if strings.Contains(strings.ToLower(o.Target), strings.ToLower(vm.TargetTypeRaw)) ||
				strings.Contains(strings.ToLower(o.SBUName), search) {
				filtered = append(filtered, o)
			}
		}
		objectives = filtered
	}

	// Apply search string filter
	if vm.SearchString != "" {
		search := strings.ToLower(strings.TrimSpace(vm.SearchString))
		var filtered []performance.ConsolidatedObjectiveVm
		for _, o := range objectives {
			if strings.Contains(strings.ToLower(o.Name), search) ||
				strings.Contains(strings.ToLower(o.Description), search) ||
				strings.Contains(strings.ToLower(o.Kpi), search) ||
				strings.Contains(strings.ToLower(o.SBUName), search) {
				filtered = append(filtered, o)
			}
		}
		objectives = filtered
	}

	totalRecords := len(objectives)

	// Paginate
	pageSize := vm.PageSize
	pageNumber := vm.PageNumber
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}
	skip := (pageNumber - 1) * pageSize
	if skip > len(objectives) {
		skip = len(objectives)
	}
	end := skip + pageSize
	if end > len(objectives) {
		end = len(objectives)
	}
	paged := objectives[skip:end]

	return &performance.GenericResponseVm{
		Data:         paged,
		TotalRecords: totalRecords,
		Message:      msgOperationCompleted,
		IsSuccess:    true,
	}, nil
}

// GetAllWorkProductDefinitions returns all work product definitions across all levels.
func (s *strategyService) GetAllWorkProductDefinitions(ctx context.Context) (interface{}, error) {
	resp := &performance.WorkProductDefinitionResponseVm{}

	// Enterprise
	entObjs, _ := s.entObjRepo.GetAllIncluding(ctx, "Strategy")
	for _, o := range entObjs {
		sbuName := ""
		if o.Strategy != nil {
			sbuName = o.Strategy.Name
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", o.EnterpriseObjectiveID, false).Find(&wps)
		for _, wp := range wps {
			resp.WorkProductDefinitions = append(resp.WorkProductDefinitions, performance.WorkProductDefinitionVm{
				WorkProductDefinitionID: wp.WorkProductDefinitionID,
				ReferenceNo:             wp.ReferenceNo,
				Name:                    wp.Name,
				Description:             wp.Description,
				Deliverables:            wp.Deliverables,
				ObjectiveID:             wp.ObjectiveID,
				ObjectiveLevel:          objectiveLevelName(enums.ObjectiveLevelEnterprise),
				ObjectiveName:           o.Name,
				SBUName:                 sbuName,
				Grade:                   "N/A",
			})
		}
	}

	// Department
	deptObjs, _ := s.deptObjRepo.GetAllIncluding(ctx, "Department")
	for _, o := range deptObjs {
		sbuName := ""
		if o.Department != nil {
			sbuName = o.Department.DepartmentName
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", o.DepartmentObjectiveID, false).Find(&wps)
		for _, wp := range wps {
			resp.WorkProductDefinitions = append(resp.WorkProductDefinitions, performance.WorkProductDefinitionVm{
				WorkProductDefinitionID: wp.WorkProductDefinitionID,
				ReferenceNo:             wp.ReferenceNo,
				Name:                    wp.Name,
				Description:             wp.Description,
				Deliverables:            wp.Deliverables,
				ObjectiveID:             wp.ObjectiveID,
				ObjectiveLevel:          objectiveLevelName(enums.ObjectiveLevelDepartment),
				ObjectiveName:           o.Name,
				SBUName:                 sbuName,
				Grade:                   "N/A",
			})
		}
	}

	// Division
	divObjs, _ := s.divObjRepo.GetAllIncluding(ctx, "Division")
	for _, o := range divObjs {
		sbuName := ""
		if o.Division != nil {
			sbuName = o.Division.DivisionName
		}
		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", o.DivisionObjectiveID, false).Find(&wps)
		for _, wp := range wps {
			resp.WorkProductDefinitions = append(resp.WorkProductDefinitions, performance.WorkProductDefinitionVm{
				WorkProductDefinitionID: wp.WorkProductDefinitionID,
				ReferenceNo:             wp.ReferenceNo,
				Name:                    wp.Name,
				Description:             wp.Description,
				Deliverables:            wp.Deliverables,
				ObjectiveID:             wp.ObjectiveID,
				ObjectiveLevel:          objectiveLevelName(enums.ObjectiveLevelDivision),
				ObjectiveName:           o.Name,
				SBUName:                 sbuName,
				Grade:                   "N/A",
			})
		}
	}

	// Office (with JobGradeGroup lookup)
	offObjs, _ := s.offObjRepo.GetAllIncluding(ctx, "Office")
	for _, o := range offObjs {
		sbuName := ""
		if o.Office != nil {
			sbuName = o.Office.OfficeName
		}
		grade := "N/A"
		// Lookup JobGradeGroup name
		var gradeGroup struct {
			GroupName string `gorm:"column:group_name"`
		}
		if err := s.db.WithContext(ctx).Table("pms.job_grade_groups").
			Where("job_grade_group_id = ?", o.JobGradeGroupID).
			First(&gradeGroup).Error; err == nil {
			grade = gradeGroup.GroupName
		}

		var wps []performance.WorkProductDefinition
		s.db.WithContext(ctx).Where("objective_id = ? AND soft_deleted = ?", o.OfficeObjectiveID, false).Find(&wps)
		for _, wp := range wps {
			resp.WorkProductDefinitions = append(resp.WorkProductDefinitions, performance.WorkProductDefinitionVm{
				WorkProductDefinitionID: wp.WorkProductDefinitionID,
				ReferenceNo:             wp.ReferenceNo,
				Name:                    wp.Name,
				Description:             wp.Description,
				Deliverables:            wp.Deliverables,
				ObjectiveID:             wp.ObjectiveID,
				ObjectiveLevel:          objectiveLevelName(enums.ObjectiveLevelOffice),
				ObjectiveName:           o.Name,
				SBUName:                 sbuName,
				Grade:                   grade,
			})
		}
	}

	resp.TotalRecords = len(resp.WorkProductDefinitions)
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetAllPaginatedWorkProductDefinitions returns paginated work product definitions with search.
func (s *strategyService) GetAllPaginatedWorkProductDefinitions(ctx context.Context, pageIndex, pageSize int, search string) (interface{}, error) {
	resp := &performance.PaginatedWorkProductDefinitionResponseVm{}

	allResult, err := s.GetAllWorkProductDefinitions(ctx)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}
	wpResp, ok := allResult.(*performance.WorkProductDefinitionResponseVm)
	if !ok {
		resp.HasError = true
		resp.Message = "unexpected result type"
		return resp, nil
	}

	wps := wpResp.WorkProductDefinitions

	// Apply search filter
	if search != "" {
		if strings.ToLower(search) == "all_enterprise_only" {
			var filtered []performance.WorkProductDefinitionVm
			for _, wp := range wps {
				if strings.EqualFold(wp.ObjectiveLevel, "enterprise") {
					filtered = append(filtered, wp)
				}
			}
			wps = filtered
		} else {
			lowerSearch := strings.ToLower(search)
			var filtered []performance.WorkProductDefinitionVm
			for _, wp := range wps {
				if strings.Contains(strings.ToLower(wp.SBUName), lowerSearch) ||
					strings.Contains(strings.ToLower(wp.ObjectiveLevel), lowerSearch) ||
					strings.Contains(strings.ToLower(wp.ObjectiveName), lowerSearch) ||
					strings.Contains(strings.ToLower(wp.Name), lowerSearch) {
					filtered = append(filtered, wp)
				}
			}
			wps = filtered
		}
	}

	totalItems := len(wps)

	// Paginate
	skip := pageIndex * pageSize
	if skip > len(wps) {
		skip = len(wps)
	}
	end := skip + pageSize
	if end > len(wps) {
		end = len(wps)
	}
	paged := wps[skip:end]

	resp.WorkProductDefinitions = &performance.PaginatedResult[performance.WorkProductDefinitionVm]{
		Items:      paged,
		TotalCount: totalItems,
	}
	return resp, nil
}

// GetFeedbackQuestionnaires returns all feedback questionnaires with options.
func (s *strategyService) GetFeedbackQuestionnaires(ctx context.Context) (interface{}, error) {
	resp := &performance.FeedbackQuestionaireResponseVm{}
	questionnaires, err := s.feedbackRepo.GetAllIncluding(ctx, "Options")
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	for _, fb := range questionnaires {
		vm := performance.FeedbackQuestionaireVm{
			FeedbackQuestionaireID: fb.FeedbackQuestionaireID,
			Question:               fb.Question,
			Description:            fb.Description,
			PmsCompetencyID:        fb.PmsCompetencyID,
			RecordStatus:           enums.Status(fb.ID), // placeholder
		}
		for _, opt := range fb.Options {
			vm.Options = append(vm.Options, performance.FeedbackQuestionaireOptionVm{
				FeedbackQuestionaireOptionID: opt.FeedbackQuestionaireOptionID,
				OptionStatement:              opt.OptionStatement,
				Description:                  opt.Description,
				Score:                        opt.Score,
				QuestionID:                   opt.QuestionID,
			})
		}
		resp.FeedbackQuestionaires = append(resp.FeedbackQuestionaires, vm)
	}

	resp.TotalRecords = len(resp.FeedbackQuestionaires)
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetDepartmentObjectives returns all department objectives.
func (s *strategyService) GetDepartmentObjectives(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.getAllDepartmentObjectives(ctx, enums.StatusAll)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.DepartmentObjectiveVmDTO
	for _, o := range objs {
		vm := performance.DepartmentObjectiveVmDTO{
			DepartmentObjectiveID: o.DepartmentObjectiveID,
			Name:                  o.Name,
			Description:           o.Description,
			Kpi:                   o.Kpi,
			Target:                o.Target,
			DepartmentID:          o.DepartmentID,
			EnterpriseObjectiveID: o.EnterpriseObjectiveID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		vm.IsActive = o.IsActive
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetEnterpriseObjectives returns all enterprise objectives (non-empty names, sorted).
func (s *strategyService) GetEnterpriseObjectives(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.getAllEnterpriseObjectives(ctx, enums.StatusAll)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.EnterpriseObjectiveVmDTO
	for _, o := range objs {
		if o.Name == "" {
			continue
		}
		vm := performance.EnterpriseObjectiveVmDTO{
			EnterpriseObjectiveID:          o.EnterpriseObjectiveID,
			Name:                           o.Name,
			Description:                    o.Description,
			Kpi:                            o.Kpi,
			Target:                         o.Target,
			EnterpriseObjectivesCategoryID: o.EnterpriseObjectivesCategoryID,
			StrategyID:                     o.StrategyID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		vm.IsActive = o.IsActive
		list = append(list, vm)
	}

	// Sort by name (simple insertion sort for stability)
	for i := 1; i < len(list); i++ {
		for j := i; j > 0 && strings.ToLower(list[j].Name) < strings.ToLower(list[j-1].Name); j-- {
			list[j], list[j-1] = list[j-1], list[j]
		}
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetDivisionObjectives returns all division objectives.
func (s *strategyService) GetDivisionObjectives(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.getAllDivisionObjectives(ctx, enums.StatusAll)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.DivisionObjectiveVmDTO
	for _, o := range objs {
		vm := performance.DivisionObjectiveVmDTO{
			DivisionObjectiveID:   o.DivisionObjectiveID,
			Name:                  o.Name,
			Description:           o.Description,
			Kpi:                   o.Kpi,
			Target:                o.Target,
			DivisionID:            o.DivisionID,
			DepartmentObjectiveID: o.DepartmentObjectiveID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		vm.IsActive = o.IsActive
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetDivisionObjectivesByDivisionId returns division objectives for a specific division.
func (s *strategyService) GetDivisionObjectivesByDivisionId(ctx context.Context, divisionID int) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.divObjRepo.Where(ctx, "division_id = ?", divisionID)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.DivisionObjectiveVmDTO
	for _, o := range objs {
		vm := performance.DivisionObjectiveVmDTO{
			DivisionObjectiveID:   o.DivisionObjectiveID,
			Name:                  o.Name,
			Description:           o.Description,
			Kpi:                   o.Kpi,
			Target:                o.Target,
			DivisionID:            o.DivisionID,
			DepartmentObjectiveID: o.DepartmentObjectiveID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetOfficeObjectives returns all office objectives with job grade group name.
func (s *strategyService) GetOfficeObjectives(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.getAllOfficeObjectives(ctx, enums.StatusAll)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.OfficeObjectiveVmDTO
	for _, o := range objs {
		vm := performance.OfficeObjectiveVmDTO{
			OfficeObjectiveID:   o.OfficeObjectiveID,
			Name:                o.Name,
			Description:         o.Description,
			Kpi:                 o.Kpi,
			Target:              o.Target,
			OfficeID:            o.OfficeID,
			DivisionObjectiveID: o.DivisionObjectiveID,
			JobGradeGroupID:     o.JobGradeGroupID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		vm.IsActive = o.IsActive
		// Lookup grade group name
		var gradeGroup struct {
			GroupName string `gorm:"column:group_name"`
		}
		if err := s.db.WithContext(ctx).Table("pms.job_grade_groups").
			Where("job_grade_group_id = ?", o.JobGradeGroupID).
			First(&gradeGroup).Error; err == nil {
			vm.JobGradeGroupName = gradeGroup.GroupName
		}
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetOfficeObjectivesByOfficeId returns office objectives for a specific office.
func (s *strategyService) GetOfficeObjectivesByOfficeId(ctx context.Context, officeID int) (interface{}, error) {
	resp := &performance.GenericListVm{}
	objs, err := s.offObjRepo.Where(ctx, "office_id = ?", officeID)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.OfficeObjectiveVmDTO
	for _, o := range objs {
		vm := performance.OfficeObjectiveVmDTO{
			OfficeObjectiveID:   o.OfficeObjectiveID,
			Name:                o.Name,
			Description:         o.Description,
			Kpi:                 o.Kpi,
			Target:              o.Target,
			OfficeID:            o.OfficeID,
			DivisionObjectiveID: o.DivisionObjectiveID,
			JobGradeGroupID:     o.JobGradeGroupID,
		}
		vm.ID = o.ID
		vm.Status = o.Status
		var gradeGroup struct {
			GroupName string `gorm:"column:group_name"`
		}
		if err := s.db.WithContext(ctx).Table("pms.job_grade_groups").
			Where("job_grade_group_id = ?", o.JobGradeGroupID).
			First(&gradeGroup).Error; err == nil {
			vm.JobGradeGroupName = gradeGroup.GroupName
		}
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetObjectiveCategories returns all objective categories.
func (s *strategyService) GetObjectiveCategories(ctx context.Context) (interface{}, error) {
	resp := &performance.GenericListVm{}
	cats, err := s.getAllObjectiveCategories(ctx, enums.StatusAll)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.ObjectiveCategoryVmDTO
	for _, c := range cats {
		vm := performance.ObjectiveCategoryVmDTO{
			ObjectiveCategoryID: c.ObjectiveCategoryID,
			Name:                c.Name,
			Description:         c.Description,
		}
		vm.ID = c.ID
		vm.Status = c.Status
		vm.IsActive = c.IsActive
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetCategoryDefinitions returns category definitions for a given category.
func (s *strategyService) GetCategoryDefinitions(ctx context.Context, categoryID string) (interface{}, error) {
	resp := &performance.GenericListVm{}
	defs, err := s.getAllCategoryDefinitions(ctx, categoryID)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	var list []performance.CategoryDefinitionVm
	for _, d := range defs {
		vm := performance.CategoryDefinitionVm{
			DefinitionID:            d.DefinitionID,
			ObjectiveCategoryID:     d.ObjectiveCategoryID,
			Weight:                  d.Weight,
			MaxNoObjectives:         d.MaxNoObjectives,
			MaxNoWorkProduct:        d.MaxNoWorkProduct,
			MaxPoints:               int(d.MaxPoints),
			IsCompulsory:            d.IsCompulsory,
			EnforceWorkProductLimit: d.EnforceWorkProductLimit,
			Description:             d.Description,
			GradeGroupID:            d.GradeGroupID,
		}
		vm.ID = d.ID
		vm.Status = d.Status
		list = append(list, vm)
	}

	resp.TotalRecord = len(list)
	resp.ListData = list
	resp.Message = msgOperationCompleted
	return resp, nil
}

// GetPmsCompetencies returns all PMS competencies.
func (s *strategyService) GetPmsCompetencies(ctx context.Context) (interface{}, error) {
	resp := &performance.PmsCompetencyListResponseVm{}
	comps, err := s.pmsCompRepo.GetAllIncluding(ctx)
	if err != nil {
		resp.HasError = true
		resp.Message = err.Error()
		return resp, nil
	}

	for _, c := range comps {
		vm := performance.PmsCompetencyVm{
			PmsCompetencyID:  c.PmsCompetencyID,
			Name:             c.Name,
			Description:      c.Description,
			ObjectCategoryID: c.ObjectCategoryID,
		}
		resp.PmsCompetencies = append(resp.PmsCompetencies, vm)
	}

	resp.TotalRecords = len(resp.PmsCompetencies)
	resp.Message = msgOperationCompleted
	return resp, nil
}
