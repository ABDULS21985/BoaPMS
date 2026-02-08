package service

// strategy_service_setup.go contains all CRUD setup methods for the
// PerformanceManagementService interface. This mirrors the .NET
// SMDService.cs #region setup (lines 1224-1961).

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
)

// ---------------------------------------------------------------------------
// Strategy CRUD
// ---------------------------------------------------------------------------

// CreateStrategy creates a new corporate strategy. Mirrors .NET
// SMDService.StrategySetup with OperationTypes.Add.
func (s *strategyService) CreateStrategy(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CreateNewStrategyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateStrategy")
	}

	resp := &performance.GenericResponseVm{}

	// Duplicate check by name
	existing, err := s.strategyRepo.FirstOrDefault(ctx,
		"LOWER(name) = ?", strings.ToLower(strings.TrimSpace(vm.Name)),
	)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateStrategy: error checking duplicate")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if existing != nil {
		resp.Errors = append(resp.Errors, "Strategy name already exists")
		return resp, nil
	}

	// Generate strategy ID
	strategyID, err := s.seqGen.GenerateCode(ctx, enums.SeqStrategy, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateStrategy: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	strategy := performance.Strategy{
		StrategyID:  strategyID,
		Name:        vm.Name,
		Description: vm.Description,
		BankYearID:  vm.BankYearID,
		StartDate:   vm.StartDate,
		EndDate:     vm.EndDate,
		FileImage:   vm.ImageFile,
		BaseWorkFlow: domain.BaseWorkFlow{
			BaseEntity: domain.BaseEntity{
				RecordStatus: enums.StatusActive.String(),
				Status:       enums.StatusApprovedAndActive.String(),
				IsActive:     true,
				CreatedBy:    s.userCtx.GetUserID(ctx),
			},
		},
	}

	if err := s.db.WithContext(ctx).Create(&strategy).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateStrategy: error saving strategy")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateStrategy updates an existing corporate strategy. Mirrors .NET
// SMDService.StrategySetup with OperationTypes.Update.
func (s *strategyService) UpdateStrategy(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.StrategyVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateStrategy")
	}

	resp := &performance.GenericResponseVm{}

	// Fetch by strategy ID (direct DB query, matching .NET _coreDbContext.Strategies.FirstOrDefault)
	var strategy performance.Strategy
	if err := s.db.WithContext(ctx).
		Where("strategy_id = ?", vm.StrategyID).
		First(&strategy).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateStrategy: strategy not found")
		resp.Errors = append(resp.Errors, "Strategy not found")
		return resp, nil
	}

	// Duplicate name check (exclude current)
	var nameCheck performance.Strategy
	nameErr := s.db.WithContext(ctx).
		Where("LOWER(name) = ? AND strategy_id != ?",
			strings.ToLower(strings.TrimSpace(vm.Name)), vm.StrategyID).
		First(&nameCheck).Error
	if nameErr == nil {
		// Found a duplicate
		resp.Errors = append(resp.Errors, "Strategy name already exists")
		return resp, nil
	}

	strategy.Name = vm.Name
	strategy.StartDate = vm.StartDate
	strategy.EndDate = vm.EndDate
	strategy.Description = vm.Description
	strategy.BankYearID = vm.BankYearID

	// Handle image file if provided
	if vm.ImageFile != "" {
		strategy.FileImage = vm.ImageFile
	}

	strategy.RecordStatus = enums.StatusActive.String()
	strategy.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(&strategy).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateStrategy: error saving strategy")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Strategic Theme CRUD
// ---------------------------------------------------------------------------

// CreateStrategicTheme creates a new strategic theme. Mirrors .NET
// SMDService.StrategicThemeSetup with OperationTypes.Add.
func (s *strategyService) CreateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CreateStrategicThemeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateStrategicTheme")
	}

	resp := &performance.GenericResponseVm{}

	// Validate strategy exists
	var strategy performance.Strategy
	if err := s.db.WithContext(ctx).
		Where("LOWER(strategy_id) = ?", strings.ToLower(vm.StrategyID)).
		First(&strategy).Error; err != nil {
		resp.Errors = append(resp.Errors, "Strategy is invalid")
		return resp, nil
	}

	// Duplicate name check
	existing, err := s.themeRepo.FirstOrDefault(ctx,
		"LOWER(name) = ?", strings.ToLower(strings.TrimSpace(vm.Name)),
	)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateStrategicTheme: error checking duplicate")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if existing != nil {
		resp.Errors = append(resp.Errors, "Strategic theme name already exists")
		return resp, nil
	}

	// Generate theme ID
	themeID, err := s.seqGen.GenerateCode(ctx, enums.SeqStrategicTheme, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateStrategicTheme: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	theme := performance.StrategicTheme{
		StrategicThemeID: themeID,
		Name:             vm.Name,
		Description:      vm.Description,
		StrategyID:       vm.StrategyID,
		FileImage:        vm.ImageFile,
		BaseWorkFlow: domain.BaseWorkFlow{
			BaseEntity: domain.BaseEntity{
				RecordStatus: enums.StatusActive.String(),
				Status:       enums.StatusApprovedAndActive.String(),
				IsActive:     true,
				CreatedBy:    s.userCtx.GetUserID(ctx),
			},
		},
	}

	if err := s.db.WithContext(ctx).Create(&theme).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateStrategicTheme: error saving theme")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateStrategicTheme updates an existing strategic theme. Mirrors .NET
// SMDService.StrategicThemeSetup with OperationTypes.Update.
func (s *strategyService) UpdateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.StrategicThemeVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateStrategicTheme")
	}

	resp := &performance.GenericResponseVm{}

	// Validate strategy exists
	var strategy performance.Strategy
	if err := s.db.WithContext(ctx).
		Where("LOWER(strategy_id) = ?", strings.ToLower(vm.StrategyID)).
		First(&strategy).Error; err != nil {
		resp.Errors = append(resp.Errors, "Strategy is invalid")
		return resp, nil
	}

	// Fetch existing theme
	var theme performance.StrategicTheme
	if err := s.db.WithContext(ctx).
		Where("strategic_theme_id = ?", vm.StrategicThemeID).
		First(&theme).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateStrategicTheme: theme not found")
		resp.Errors = append(resp.Errors, "Strategic theme not found")
		return resp, nil
	}

	// Duplicate name check (exclude current)
	var nameCheck performance.StrategicTheme
	nameErr := s.db.WithContext(ctx).
		Where("LOWER(name) = ? AND strategic_theme_id != ?",
			strings.ToLower(strings.TrimSpace(vm.Name)), vm.StrategicThemeID).
		First(&nameCheck).Error
	if nameErr == nil {
		resp.Errors = append(resp.Errors, "Strategic theme name already exists")
		return resp, nil
	}

	theme.Name = vm.Name
	theme.Description = vm.Description

	// Handle image file if provided
	if vm.ImageFile != "" {
		theme.FileImage = vm.ImageFile
	}

	theme.RecordStatus = enums.StatusActive.String()
	theme.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(&theme).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateStrategicTheme: error saving theme")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Enterprise Objective CRUD
// ---------------------------------------------------------------------------

// CreateEnterpriseObjective creates a new enterprise-level objective. Mirrors
// .NET SMDService.EnterpriseObjectiveSetup with OperationTypes.Add.
func (s *strategyService) CreateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.EnterpriseObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateEnterpriseObjective")
	}

	resp := &performance.GenericResponseVm{}

	// Duplicate check by name + KPI
	var existing performance.EnterpriseObjective
	err := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ? AND LOWER(TRIM(kpi)) = ?",
			strings.ToLower(strings.TrimSpace(vm.Name)),
			strings.ToLower(strings.TrimSpace(vm.Kpi))).
		First(&existing).Error
	if err == nil {
		resp.Errors = append(resp.Errors, "Enterprise Objective already exists")
		return resp, nil
	}

	// Generate enterprise objective ID
	objID, err := s.seqGen.GenerateCode(ctx, enums.SeqEnterpriseObjective, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateEnterpriseObjective: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	obj := performance.EnterpriseObjective{
		EnterpriseObjectiveID:          objID,
		EnterpriseObjectivesCategoryID: vm.EnterpriseObjectivesCategoryID,
		StrategyID:                     vm.StrategyID,
		ObjectiveBase: domain.ObjectiveBase{
			BaseWorkFlow: domain.BaseWorkFlow{
				BaseEntity: domain.BaseEntity{
					RecordStatus: enums.StatusActive.String(),
					Status:       enums.StatusApprovedAndActive.String(),
					IsActive:     true,
					CreatedBy:    s.userCtx.GetUserID(ctx),
				},
			},
			Name:        vm.Name,
			Description: vm.Description,
			Kpi:         vm.Kpi,
			Target:      vm.Target,
		},
	}

	if err := s.db.WithContext(ctx).Create(&obj).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateEnterpriseObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateEnterpriseObjective updates an existing enterprise objective. Mirrors
// .NET SMDService.EnterpriseObjectiveSetup with OperationTypes.Update.
func (s *strategyService) UpdateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.EnterpriseObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateEnterpriseObjective")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getEnterpriseObjective(ctx, vm.EnterpriseObjectiveID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateEnterpriseObjective: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Enterprise Objective not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.Kpi = vm.Kpi
	obj.Target = vm.Target
	obj.StrategyID = vm.StrategyID
	obj.EnterpriseObjectivesCategoryID = vm.EnterpriseObjectivesCategoryID
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateEnterpriseObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Department Objective CRUD
// ---------------------------------------------------------------------------

// CreateDepartmentObjective creates a new department-level objective. Mirrors
// .NET SMDService.DepartmentObjectiveSetup with OperationTypes.Add.
func (s *strategyService) CreateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.DepartmentObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateDepartmentObjective")
	}

	resp := &performance.GenericResponseVm{}

	// In .NET, duplicate check does NOT throw - it just overwrites. Match that behavior.
	var existing performance.DepartmentObjective
	dupErr := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ? AND department_id = ?",
			strings.ToLower(strings.TrimSpace(vm.Name)), vm.DepartmentID).
		First(&existing).Error

	var deptObj performance.DepartmentObjective
	if dupErr == nil {
		// Duplicate exists: update it in place (matches .NET behavior)
		deptObj = existing
		deptObj.Name = vm.Name
		deptObj.Description = vm.Description
		deptObj.Kpi = vm.Kpi
		deptObj.Target = vm.Target
		deptObj.DepartmentID = vm.DepartmentID
		deptObj.EnterpriseObjectiveID = vm.EnterpriseObjectiveID
	} else {
		// Generate department objective ID
		objID, err := s.seqGen.GenerateCode(ctx, enums.SeqDepartmentObjective, 15, "", enums.ConCatBefore)
		if err != nil {
			s.log.Error().Err(err).Msg("CreateDepartmentObjective: error generating code")
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}

		deptObj = performance.DepartmentObjective{
			DepartmentObjectiveID: objID,
			DepartmentID:          vm.DepartmentID,
			EnterpriseObjectiveID: vm.EnterpriseObjectiveID,
			ObjectiveBase: domain.ObjectiveBase{
				Name:        vm.Name,
				Description: vm.Description,
				Kpi:         vm.Kpi,
				Target:      vm.Target,
			},
		}
	}

	deptObj.RecordStatus = enums.StatusActive.String()
	deptObj.Status = enums.StatusApprovedAndActive.String()
	deptObj.IsActive = true
	deptObj.CreatedBy = s.userCtx.GetUserID(ctx)

	if err := s.db.WithContext(ctx).Save(&deptObj).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateDepartmentObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	// Save work product if provided (mirrors .NET logic)
	wpName := ""
	if vm.WorkProductName != nil {
		wpName = *vm.WorkProductName
	}
	if strings.TrimSpace(wpName) != "" {
		wpDesc := ""
		if vm.WorkProductDescription != nil {
			wpDesc = *vm.WorkProductDescription
		}
		wpDeliv := ""
		if vm.WorkProductDeliverable != nil {
			wpDeliv = *vm.WorkProductDeliverable
		}
		sbuName := ""
		if vm.SBUName != nil {
			sbuName = *vm.SBUName
		}
		grade := ""
		if vm.JobGradeGroup != nil {
			grade = *vm.JobGradeGroup
		}
		wpReq := &performance.WorkProductDefinitionVm{
			BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
			Name:           wpName,
			Description:    wpDesc,
			ObjectiveID:    deptObj.DepartmentObjectiveID,
			Deliverables:   wpDeliv,
			SBUName:        sbuName,
			Grade:          grade,
			ObjectiveName:  deptObj.Name,
			ObjectiveLevel: "Department",
		}
		if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
			s.log.Error().Err(wpErr).Msg("CreateDepartmentObjective: error saving work product")
			resp.Errors = append(resp.Errors, wpErr.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateDepartmentObjective updates an existing department objective. Mirrors
// .NET SMDService.DepartmentObjectiveSetup with OperationTypes.Update.
func (s *strategyService) UpdateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.DepartmentObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateDepartmentObjective")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getDepartmentObjective(ctx, vm.DepartmentObjectiveID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateDepartmentObjective: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Department Objective not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.Kpi = vm.Kpi
	obj.Target = vm.Target
	obj.DepartmentID = vm.DepartmentID
	obj.EnterpriseObjectiveID = vm.EnterpriseObjectiveID
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateDepartmentObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	// Save work product if provided
	wpName := ""
	if vm.WorkProductName != nil {
		wpName = *vm.WorkProductName
	}
	if strings.TrimSpace(wpName) != "" {
		wpDesc := ""
		if vm.WorkProductDescription != nil {
			wpDesc = *vm.WorkProductDescription
		}
		wpDeliv := ""
		if vm.WorkProductDeliverable != nil {
			wpDeliv = *vm.WorkProductDeliverable
		}
		sbuName := ""
		if vm.SBUName != nil {
			sbuName = *vm.SBUName
		}
		grade := ""
		if vm.JobGradeGroup != nil {
			grade = *vm.JobGradeGroup
		}
		wpReq := &performance.WorkProductDefinitionVm{
			BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
			Name:           wpName,
			Description:    wpDesc,
			ObjectiveID:    obj.DepartmentObjectiveID,
			Deliverables:   wpDeliv,
			SBUName:        sbuName,
			Grade:          grade,
			ObjectiveName:  obj.Name,
			ObjectiveLevel: "Department",
		}
		if wpErr := s.workProductDefinitionSetup(ctx, wpReq, false); wpErr != nil {
			s.log.Error().Err(wpErr).Msg("UpdateDepartmentObjective: error saving work product")
			resp.Errors = append(resp.Errors, wpErr.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Division Objective CRUD
// ---------------------------------------------------------------------------

// CreateDivisionObjective creates a new division-level objective. Mirrors
// .NET SMDService.DivisionObjectiveSetup with OperationTypes.Add.
func (s *strategyService) CreateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.DivisionObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateDivisionObjective")
	}

	resp := &performance.GenericResponseVm{}

	// Check if duplicate exists by name + division
	var existing performance.DivisionObjective
	dupErr := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ? AND division_id = ?",
			strings.ToLower(strings.TrimSpace(vm.Name)), vm.DivisionID).
		First(&existing).Error

	var divObj performance.DivisionObjective
	if dupErr == nil {
		// Duplicate exists: update in place (mirrors .NET behavior which does NOT throw)
		divObj = existing
		divObj.Name = vm.Name
		divObj.Description = vm.Description
		divObj.Kpi = vm.Kpi
		divObj.Target = vm.Target
		divObj.DepartmentObjectiveID = vm.DepartmentObjectiveID
	} else {
		// Generate division objective ID
		objID, err := s.seqGen.GenerateCode(ctx, enums.SeqDivisionObjective, 15, "", enums.ConCatBefore)
		if err != nil {
			s.log.Error().Err(err).Msg("CreateDivisionObjective: error generating code")
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}

		divObj = performance.DivisionObjective{
			DivisionObjectiveID:   objID,
			DivisionID:            vm.DivisionID,
			DepartmentObjectiveID: vm.DepartmentObjectiveID,
			ObjectiveBase: domain.ObjectiveBase{
				Name:        vm.Name,
				Description: vm.Description,
				Kpi:         vm.Kpi,
				Target:      vm.Target,
			},
		}
	}

	divObj.RecordStatus = enums.StatusActive.String()
	divObj.Status = enums.StatusApprovedAndActive.String()
	divObj.IsActive = true
	divObj.CreatedBy = s.userCtx.GetUserID(ctx)

	if err := s.db.WithContext(ctx).Save(&divObj).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateDivisionObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	// Save work product if provided
	wpName := ""
	if vm.WorkProductName != nil {
		wpName = *vm.WorkProductName
	}
	if strings.TrimSpace(wpName) != "" {
		wpDesc := ""
		if vm.WorkProductDescription != nil {
			wpDesc = *vm.WorkProductDescription
		}
		wpDeliv := ""
		if vm.WorkProductDeliverable != nil {
			wpDeliv = *vm.WorkProductDeliverable
		}
		sbuName := ""
		if vm.SBUName != nil {
			sbuName = *vm.SBUName
		}
		wpReq := &performance.WorkProductDefinitionVm{
			BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
			Name:           wpName,
			Description:    wpDesc,
			ObjectiveID:    divObj.DivisionObjectiveID,
			Deliverables:   wpDeliv,
			SBUName:        sbuName,
			Grade:          vm.JobGradeGroup,
			ObjectiveName:  divObj.Name,
			ObjectiveLevel: "Division",
		}
		if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
			s.log.Error().Err(wpErr).Msg("CreateDivisionObjective: error saving work product")
			resp.Errors = append(resp.Errors, wpErr.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateDivisionObjective updates an existing division objective. Mirrors
// .NET SMDService.DivisionObjectiveSetup with OperationTypes.Update.
func (s *strategyService) UpdateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.DivisionObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateDivisionObjective")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getDivisionObjective(ctx, vm.DivisionObjectiveID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateDivisionObjective: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Division Objective not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.Kpi = vm.Kpi
	obj.Target = vm.Target
	obj.DepartmentObjectiveID = vm.DepartmentObjectiveID
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateDivisionObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	// Save work product if provided
	wpName := ""
	if vm.WorkProductName != nil {
		wpName = *vm.WorkProductName
	}
	if strings.TrimSpace(wpName) != "" {
		wpDesc := ""
		if vm.WorkProductDescription != nil {
			wpDesc = *vm.WorkProductDescription
		}
		wpDeliv := ""
		if vm.WorkProductDeliverable != nil {
			wpDeliv = *vm.WorkProductDeliverable
		}
		sbuName := ""
		if vm.SBUName != nil {
			sbuName = *vm.SBUName
		}
		wpReq := &performance.WorkProductDefinitionVm{
			BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
			Name:           wpName,
			Description:    wpDesc,
			ObjectiveID:    obj.DivisionObjectiveID,
			Deliverables:   wpDeliv,
			SBUName:        sbuName,
			Grade:          vm.JobGradeGroup,
			ObjectiveName:  obj.Name,
			ObjectiveLevel: "Division",
		}
		if wpErr := s.workProductDefinitionSetup(ctx, wpReq, false); wpErr != nil {
			s.log.Error().Err(wpErr).Msg("UpdateDivisionObjective: error saving work product")
			resp.Errors = append(resp.Errors, wpErr.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Office Objective CRUD
// ---------------------------------------------------------------------------

// CreateOfficeObjective creates a new office-level objective. Mirrors .NET
// SMDService.OfficeObjectiveSetup with OperationTypes.Add.
//
// This is the most complex setup method because it auto-creates a parent
// division objective if one does not already exist for the given parent
// division + objective name.
func (s *strategyService) CreateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.OfficeObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateOfficeObjective")
	}

	resp := &performance.GenericResponseVm{}

	// ---- Step 1: ensure parent division objective exists ----
	var divisionObjective *performance.DivisionObjective
	if vm.ParentDivisionObjectiveName != "" {
		var divObj performance.DivisionObjective
		err := s.db.WithContext(ctx).
			Where("LOWER(TRIM(name)) = ? AND division_id = ?",
				strings.ToLower(strings.TrimSpace(vm.ParentDivisionObjectiveName)),
				vm.ParentDivisionID).
			First(&divObj).Error
		if err != nil {
			// Not found: auto-create the division objective (mirrors .NET behavior)
			var division organogram.Division
			if divErr := s.db.WithContext(ctx).
				Where("division_id = ?", vm.ParentDivisionID).
				First(&division).Error; divErr != nil {
				resp.Errors = append(resp.Errors, "Parent division not found")
				return resp, nil
			}

			// Find first department objective for this department
			var deptObj performance.DepartmentObjective
			if deptErr := s.db.WithContext(ctx).
				Where("department_id = ?", division.DepartmentID).
				First(&deptObj).Error; deptErr != nil {
				resp.Errors = append(resp.Errors, "No department objective found for this division's department")
				return resp, nil
			}

			wpName := ""
			if vm.WorkProductName != nil {
				wpName = *vm.WorkProductName
			}
			wpDesc := ""
			if vm.WorkProductDescription != nil {
				wpDesc = *vm.WorkProductDescription
			}
			wpDeliv := ""
			if vm.WorkProductDeliverable != nil {
				wpDeliv = *vm.WorkProductDeliverable
			}

			divReq := &performance.DivisionObjectiveVmDTO{
				Name:                  vm.Name,
				Description:           vm.Description,
				Kpi:                   vm.Kpi,
				Target:                vm.Target,
				DivisionID:            vm.ParentDivisionID,
				DepartmentObjectiveID: deptObj.DepartmentObjectiveID,
				DepartmentID:          division.DepartmentID,
				JobGradeGroup:         vm.JobGradeGroupName,
				SBUName:               ptrStr(division.DivisionName),
				WorkProductName:       ptrStr(wpName),
				WorkProductDescription: ptrStr(wpDesc),
				WorkProductDeliverable: ptrStr(wpDeliv),
			}
			// Recursively create the division objective
			if _, divSetupErr := s.CreateDivisionObjective(ctx, divReq); divSetupErr != nil {
				resp.Errors = append(resp.Errors, divSetupErr.Error())
				return resp, nil
			}
		} else {
			divisionObjective = &divObj
		}

		// Re-fetch the division objective (it may have just been created)
		if divisionObjective == nil {
			var freshDivObj performance.DivisionObjective
			if err := s.db.WithContext(ctx).
				Where("LOWER(TRIM(name)) = ? AND division_id = ?",
					strings.ToLower(strings.TrimSpace(vm.ParentDivisionObjectiveName)),
					vm.ParentDivisionID).
				First(&freshDivObj).Error; err != nil {
				resp.Errors = append(resp.Errors, "Failed to find or create parent division objective")
				return resp, nil
			}
			divisionObjective = &freshDivObj
		}
	}

	// ---- Step 2: check if office objective already exists ----
	var offObj performance.OfficeObjective
	dupErr := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ? AND office_id = ? AND job_grade_group_id = ?",
			strings.ToLower(strings.TrimSpace(vm.Name)), vm.OfficeID, vm.JobGradeGroupID).
		First(&offObj).Error

	if dupErr == nil {
		// Already exists: update in place
		offObj.Name = vm.Name
		offObj.Description = vm.Description
		offObj.Kpi = vm.Kpi
		offObj.Target = vm.Target
		offObj.JobGradeGroupID = vm.JobGradeGroupID
		if divisionObjective != nil {
			offObj.DivisionObjectiveID = divisionObjective.DivisionObjectiveID
		}
		offObj.RecordStatus = enums.StatusActive.String()
		offObj.Status = enums.StatusApprovedAndActive.String()

		if err := s.db.WithContext(ctx).Save(&offObj).Error; err != nil {
			s.log.Error().Err(err).Msg("CreateOfficeObjective: error updating existing")
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	} else {
		// Create new office objective
		objID, err := s.seqGen.GenerateCode(ctx, enums.SeqOfficeObjective, 15, "", enums.ConCatBefore)
		if err != nil {
			s.log.Error().Err(err).Msg("CreateOfficeObjective: error generating code")
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}

		divObjID := vm.DivisionObjectiveID
		if divisionObjective != nil {
			divObjID = divisionObjective.DivisionObjectiveID
		}

		offObj = performance.OfficeObjective{
			OfficeObjectiveID:   objID,
			OfficeID:            vm.OfficeID,
			DivisionObjectiveID: divObjID,
			JobGradeGroupID:     vm.JobGradeGroupID,
			ObjectiveBase: domain.ObjectiveBase{
				BaseWorkFlow: domain.BaseWorkFlow{
					BaseEntity: domain.BaseEntity{
						RecordStatus: enums.StatusActive.String(),
						Status:       enums.StatusApprovedAndActive.String(),
						IsActive:     true,
						CreatedBy:    s.userCtx.GetUserID(ctx),
					},
				},
				Name:        vm.Name,
				Description: vm.Description,
				Kpi:         vm.Kpi,
				Target:      vm.Target,
			},
		}

		if err := s.db.WithContext(ctx).Create(&offObj).Error; err != nil {
			s.log.Error().Err(err).Msg("CreateOfficeObjective: error creating")
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	// Save work product if provided
	wpName := ""
	if vm.WorkProductName != nil {
		wpName = *vm.WorkProductName
	}
	if strings.TrimSpace(wpName) != "" {
		wpDesc := ""
		if vm.WorkProductDescription != nil {
			wpDesc = *vm.WorkProductDescription
		}
		wpDeliv := ""
		if vm.WorkProductDeliverable != nil {
			wpDeliv = *vm.WorkProductDeliverable
		}
		sbuName := ""
		if vm.SBUName != nil {
			sbuName = *vm.SBUName
		}
		wpReq := &performance.WorkProductDefinitionVm{
			BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
			Name:           wpName,
			Description:    wpDesc,
			ObjectiveID:    offObj.OfficeObjectiveID,
			Deliverables:   wpDeliv,
			SBUName:        sbuName,
			Grade:          vm.JobGradeGroupName,
			ObjectiveName:  vm.Name,
			ObjectiveLevel: "Office",
		}
		if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
			s.log.Error().Err(wpErr).Msg("CreateOfficeObjective: error saving work product")
			resp.Errors = append(resp.Errors, wpErr.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateOfficeObjective updates an existing office objective. Mirrors
// .NET SMDService.OfficeObjectiveSetup with OperationTypes.Update.
func (s *strategyService) UpdateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.OfficeObjectiveVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateOfficeObjective")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getOfficeObjective(ctx, vm.OfficeObjectiveID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateOfficeObjective: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Office Objective not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.Kpi = vm.Kpi
	obj.Target = vm.Target
	obj.JobGradeGroupID = vm.JobGradeGroupID
	obj.DivisionObjectiveID = vm.DivisionObjectiveID
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateOfficeObjective: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Objective Category CRUD
// ---------------------------------------------------------------------------

// CreateObjectiveCategory creates a new objective category. Mirrors
// .NET SMDService.ObjectiveCategorySetup with OperationTypes.Add.
func (s *strategyService) CreateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ObjectiveCategoryVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateObjectiveCategory")
	}

	resp := &performance.GenericResponseVm{}

	// Duplicate check by name
	var existing performance.ObjectiveCategory
	if err := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ?", strings.ToLower(strings.TrimSpace(vm.Name))).
		First(&existing).Error; err == nil {
		resp.Errors = append(resp.Errors, "Objective Category already exists")
		return resp, nil
	}

	catID, err := s.seqGen.GenerateCode(ctx, enums.SeqObjectiveCategory, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateObjectiveCategory: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	cat := performance.ObjectiveCategory{
		ObjectiveCategoryID: catID,
		Name:                vm.Name,
		Description:         vm.Description,
		BaseWorkFlow: domain.BaseWorkFlow{
			BaseEntity: domain.BaseEntity{
				RecordStatus: enums.StatusActive.String(),
				Status:       enums.StatusApprovedAndActive.String(),
				IsActive:     true,
				CreatedBy:    s.userCtx.GetUserID(ctx),
			},
		},
	}

	if err := s.db.WithContext(ctx).Create(&cat).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateObjectiveCategory: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateObjectiveCategory updates an existing objective category. Mirrors
// .NET SMDService.ObjectiveCategorySetup with OperationTypes.Update.
func (s *strategyService) UpdateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ObjectiveCategoryVmDTO)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateObjectiveCategory")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getObjectiveCategory(ctx, vm.ObjectiveCategoryID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateObjectiveCategory: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Objective Category not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateObjectiveCategory: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Category Definition CRUD
// ---------------------------------------------------------------------------

// CreateCategoryDefinition creates a new category definition. Mirrors
// .NET SMDService.SetupCategoryDefinition with OperationTypes.Add.
func (s *strategyService) CreateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreateCategoryDefinition")
	}

	resp := &performance.GenericResponseVm{}

	defID, err := s.seqGen.GenerateCode(ctx, enums.SeqCategoryDefinitions, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreateCategoryDefinition: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	catDef := performance.CategoryDefinition{
		DefinitionID:            defID,
		ObjectiveCategoryID:     vm.ObjectiveCategoryID,
		Weight:                  vm.Weight,
		MaxNoObjectives:         vm.MaxNoObjectives,
		MaxNoWorkProduct:        vm.MaxNoWorkProduct,
		MaxPoints:               float64(vm.MaxPoints),
		IsCompulsory:            vm.IsCompulsory,
		EnforceWorkProductLimit: vm.EnforceWorkProductLimit,
		Description:             vm.Description,
		GradeGroupID:            vm.GradeGroupID,
		BaseWorkFlow: domain.BaseWorkFlow{
			BaseEntity: domain.BaseEntity{
				RecordStatus: enums.StatusActive.String(),
				Status:       enums.StatusApprovedAndActive.String(),
				IsActive:     true,
				CreatedBy:    s.userCtx.GetUserID(ctx),
			},
		},
	}

	if err := s.db.WithContext(ctx).Create(&catDef).Error; err != nil {
		s.log.Error().Err(err).Msg("CreateCategoryDefinition: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdateCategoryDefinition updates an existing category definition. Mirrors
// .NET SMDService.SetupCategoryDefinition with OperationTypes.Update.
func (s *strategyService) UpdateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.CategoryDefinitionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdateCategoryDefinition")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getCategoryDefinition(ctx, vm.DefinitionID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdateCategoryDefinition: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "Category Definition not found")
		return resp, nil
	}

	obj.Description = vm.Description
	obj.EnforceWorkProductLimit = vm.EnforceWorkProductLimit
	obj.MaxNoWorkProduct = vm.MaxNoWorkProduct
	obj.GradeGroupID = vm.GradeGroupID
	obj.MaxNoObjectives = vm.MaxNoObjectives
	obj.Weight = vm.Weight
	obj.ObjectiveCategoryID = vm.ObjectiveCategoryID
	obj.IsCompulsory = vm.IsCompulsory
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdateCategoryDefinition: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// PMS Competency CRUD
// ---------------------------------------------------------------------------

// CreatePmsCompetency creates a new PMS competency. Mirrors
// .NET SMDService.PmsCompetencySetup with OperationTypes.Add.
func (s *strategyService) CreatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.PmsCompetencyRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for CreatePmsCompetency")
	}

	resp := &performance.GenericResponseVm{}

	// Duplicate check by name
	var existing performance.PmsCompetency
	if err := s.db.WithContext(ctx).
		Where("LOWER(TRIM(name)) = ?", strings.ToLower(strings.TrimSpace(vm.Name))).
		First(&existing).Error; err == nil {
		resp.Errors = append(resp.Errors, "PMS Competency already exists")
		return resp, nil
	}

	compID, err := s.seqGen.GenerateCode(ctx, enums.SeqPmsCompetency, 15, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Msg("CreatePmsCompetency: error generating code")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	comp := performance.PmsCompetency{
		PmsCompetencyID:  compID,
		Name:             vm.Name,
		Description:      vm.Description,
		ObjectCategoryID: vm.ObjectCategoryID,
		BaseWorkFlow: domain.BaseWorkFlow{
			BaseEntity: domain.BaseEntity{
				RecordStatus: enums.StatusActive.String(),
				Status:       enums.StatusApprovedAndActive.String(),
				IsActive:     true,
				CreatedBy:    s.userCtx.GetUserID(ctx),
			},
		},
	}

	if err := s.db.WithContext(ctx).Create(&comp).Error; err != nil {
		s.log.Error().Err(err).Msg("CreatePmsCompetency: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// UpdatePmsCompetency updates an existing PMS competency. Mirrors
// .NET SMDService.PmsCompetencySetup with OperationTypes.Update.
func (s *strategyService) UpdatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.PmsCompetencyRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for UpdatePmsCompetency")
	}

	resp := &performance.GenericResponseVm{}

	obj, err := s.getPmsCompetency(ctx, vm.PmsCompetencyID)
	if err != nil {
		s.log.Error().Err(err).Msg("UpdatePmsCompetency: error fetching")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	if obj == nil {
		resp.Errors = append(resp.Errors, "PMS Competency not found")
		return resp, nil
	}

	obj.Name = vm.Name
	obj.Description = vm.Description
	obj.ObjectCategoryID = vm.ObjectCategoryID
	obj.RecordStatus = enums.StatusActive.String()
	obj.Status = enums.StatusApprovedAndActive.String()

	if err := s.db.WithContext(ctx).Save(obj).Error; err != nil {
		s.log.Error().Err(err).Msg("UpdatePmsCompetency: error saving")
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Helper: convert string to *string
// ---------------------------------------------------------------------------

func ptrStr(s string) *string {
	return &s
}
