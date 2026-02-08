package service

// strategy_service_ops.go contains approval, bulk upload, deactivation,
// and save methods for evaluation/feedback/work-product entities.
// This mirrors .NET SMDService.cs #region Approval (lines 1966-2068),
// ProcessObectivesUpload (2070-2331), DeActivateOrReactivateObjectives
// (2335-2518), and ProcessUpload (2521-2585).

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
)

// ---------------------------------------------------------------------------
// Approval
// ---------------------------------------------------------------------------

// ApproveRecords approves one or more records by entity type. Mirrors .NET
// SMDService.AprroveOrRejectRecord with Approval.Approved.
func (s *strategyService) ApproveRecords(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.ApprovalRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for ApproveRecords")
	}
	return s.processApproval(ctx, &vm.ApprovalBase, enums.ApprovalApproved, "")
}

// RejectRecords rejects one or more records by entity type. Mirrors .NET
// SMDService.AprroveOrRejectRecord with Approval.Rejected.
func (s *strategyService) RejectRecords(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*performance.RejectionRequestVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for RejectRecords")
	}
	return s.processApproval(ctx, &vm.ApprovalBase, enums.ApprovalRejected, vm.RejectionReason)
}

// processApproval iterates over RecordIDs and approves or rejects each.
func (s *strategyService) processApproval(ctx context.Context, base *performance.ApprovalBase, approval enums.Approval, reason string) (interface{}, error) {
	resp := &performance.GenericResponseVm{}

	for _, recordID := range base.RecordIDs {
		var entity interface{}
		var err error

		switch base.EntityType {
		case "EnterpriseObjective":
			entity, err = s.getEnterpriseObjective(ctx, recordID)
		case "DepartmentObjective":
			entity, err = s.getDepartmentObjective(ctx, recordID)
		case "DivisionObjective":
			entity, err = s.getDivisionObjective(ctx, recordID)
		case "OfficeObjective":
			entity, err = s.getOfficeObjective(ctx, recordID)
		case "Strategy":
			entity, err = s.getStrategy(ctx, recordID)
		case "StrategicTheme":
			entity, err = s.getStrategicTheme(ctx, recordID)
		case "ObjectiveCategory":
			entity, err = s.getObjectiveCategory(ctx, recordID)
		default:
			resp.Errors = append(resp.Errors, fmt.Sprintf("unsupported entity type: %s", base.EntityType))
			return resp, nil
		}

		if err != nil {
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
		if entity == nil {
			resp.Errors = append(resp.Errors, fmt.Sprintf("%s not found", base.EntityType))
			return resp, nil
		}

		if err := s.approveOrRejectEntity(ctx, entity, approval, reason); err != nil {
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// ProcessObjectivesUpload – bulk cascaded objective upload with transaction
// ---------------------------------------------------------------------------

// ProcessObjectivesUpload processes a bulk upload of cascaded objectives.
// It groups rows by enterprise objective, then cascades through department,
// division, and office levels creating any missing objectives along the way.
// Uses a database transaction with rollback on error.
// Mirrors .NET SMDService.ProcessObectivesUpload (lines 2070-2331).
func (s *strategyService) ProcessObjectivesUpload(ctx context.Context, req interface{}) (interface{}, error) {
	rows, ok := req.([]performance.CascadedObjectiveUploadVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for ProcessObjectivesUpload")
	}

	resp := &performance.GenericResponseVm{}

	if len(rows) == 0 {
		resp.IsSuccess = true
		resp.Message = msgOperationCompleted
		return resp, nil
	}

	// Begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Errors = append(resp.Errors, tx.Error.Error())
		return resp, nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Group by enterprise objective name
	entGroups := groupByField(rows, func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.EObjName) })

	for eObjName, entries := range entGroups {
		if eObjName == "" || len(entries) == 0 {
			continue
		}
		first := entries[0]

		// Validate category
		var category performance.ObjectiveCategory
		if err := tx.Where("LOWER(TRIM(name)) = ?", strings.ToLower(strings.TrimSpace(first.EObjCategory))).
			First(&category).Error; err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Category: %s", strings.TrimSpace(first.EObjCategory)))
			return resp, nil
		}

		// Validate strategy
		var strategy performance.Strategy
		if err := tx.Where("strategy_id = ?", first.StrategyID).
			First(&strategy).Error; err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Strategy ID: %s", first.StrategyID))
			return resp, nil
		}

		// Validate strategic theme
		var theme performance.StrategicTheme
		if err := tx.Where("strategic_theme_id = ?", first.StrategicThemeID).
			First(&theme).Error; err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Strategic Theme ID: %s", first.StrategicThemeID))
			return resp, nil
		}

		// Find or create enterprise objective
		var entObj performance.EnterpriseObjective
		err := tx.Where("LOWER(TRIM(name)) = ? AND strategy_id = ? AND strategic_theme_id = ? AND enterprise_objectives_category_id = ?",
			strings.ToLower(strings.TrimSpace(eObjName)),
			strategy.StrategyID, theme.StrategicThemeID, category.ObjectiveCategoryID).
			First(&entObj).Error
		if err != nil {
			// Create new enterprise objective
			entObjID, genErr := s.seqGen.GenerateCode(ctx, enums.SeqEnterpriseObjective, 15, "", enums.ConCatBefore)
			if genErr != nil {
				tx.Rollback()
				resp.Errors = append(resp.Errors, genErr.Error())
				return resp, nil
			}
			entObj = performance.EnterpriseObjective{
				EnterpriseObjectiveID:          entObjID,
				StrategyID:                     strategy.StrategyID,
				StrategicThemeID:               theme.StrategicThemeID,
				EnterpriseObjectivesCategoryID: category.ObjectiveCategoryID,
				ObjectiveBase: domain.ObjectiveBase{
					Name:        strings.TrimSpace(first.EObjName),
					Description: strings.TrimSpace(first.EObjDesc),
					Kpi:         strings.TrimSpace(first.EObjKPI),
					Target:      strings.TrimSpace(first.EObjTarget),
					BaseWorkFlow: domain.BaseWorkFlow{
						BaseEntity: domain.BaseEntity{
							RecordStatus: enums.StatusActive.String(),
							Status:       enums.StatusApprovedAndActive.String(),
							IsActive:     true,
							CreatedBy:    s.userCtx.GetUserID(ctx),
						},
					},
				},
			}
			if createErr := tx.Create(&entObj).Error; createErr != nil {
				tx.Rollback()
				resp.Errors = append(resp.Errors, createErr.Error())
				return resp, nil
			}
		}

		// Group by department
		deptGroups := groupByTwoFields(entries,
			func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.DeptObjName) },
			func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.Dept) },
		)

		for deptKey, deptEntries := range deptGroups {
			dFirst := deptEntries[0]

			// Find department
			var dept organogram.Department
			if err := tx.Where("LOWER(TRIM(department_name)) = ?",
				strings.ToLower(strings.TrimSpace(dFirst.Dept))).
				First(&dept).Error; err != nil {
				tx.Rollback()
				resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Department: %s", dFirst.Dept))
				return resp, nil
			}

			_ = deptKey // used for grouping

			// Find or create department objective
			var deptObj performance.DepartmentObjective
			deptObjErr := tx.Where("LOWER(TRIM(name)) = ? AND department_id = ?",
				strings.ToLower(strings.TrimSpace(dFirst.DeptObjName)), dept.DepartmentID).
				First(&deptObj).Error
			if deptObjErr != nil {
				deptObjID, genErr := s.seqGen.GenerateCode(ctx, enums.SeqDepartmentObjective, 15, "", enums.ConCatBefore)
				if genErr != nil {
					tx.Rollback()
					resp.Errors = append(resp.Errors, genErr.Error())
					return resp, nil
				}
				deptObj = performance.DepartmentObjective{
					DepartmentObjectiveID: deptObjID,
					DepartmentID:          dept.DepartmentID,
					EnterpriseObjectiveID: entObj.EnterpriseObjectiveID,
					ObjectiveBase: domain.ObjectiveBase{
						Name:        strings.TrimSpace(dFirst.DeptObjName),
						Description: strings.TrimSpace(dFirst.DeptObjDesc),
						Kpi:         strings.TrimSpace(dFirst.DeptObjKPI),
						Target:      strings.TrimSpace(dFirst.DeptObjTarget),
						BaseWorkFlow: domain.BaseWorkFlow{
							BaseEntity: domain.BaseEntity{
								RecordStatus: enums.StatusActive.String(),
								Status:       enums.StatusApprovedAndActive.String(),
								IsActive:     true,
								CreatedBy:    s.userCtx.GetUserID(ctx),
							},
						},
					},
				}
				if createErr := tx.Create(&deptObj).Error; createErr != nil {
					tx.Rollback()
					resp.Errors = append(resp.Errors, createErr.Error())
					return resp, nil
				}
			}

			// Save work product for department objective if provided
			wpName := ""
			if dFirst.WorkProductName != nil {
				wpName = *dFirst.WorkProductName
			}
			if strings.TrimSpace(wpName) != "" {
				wpDesc := ""
				if dFirst.WorkProductDescription != nil {
					wpDesc = *dFirst.WorkProductDescription
				}
				wpDeliv := ""
				if dFirst.WorkProductDeliverable != nil {
					wpDeliv = *dFirst.WorkProductDeliverable
				}
				wpReq := &performance.WorkProductDefinitionVm{
					BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
					Name:           wpName,
					Description:    wpDesc,
					ObjectiveID:    deptObj.DepartmentObjectiveID,
					Deliverables:   wpDeliv,
					SBUName:        strings.TrimSpace(dFirst.Office),
					Grade:          strings.TrimSpace(dFirst.JobGradeGroup),
					ObjectiveName:  strings.TrimSpace(dFirst.DeptObjName),
					ObjectiveLevel: "Department",
				}
				if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
					s.log.Warn().Err(wpErr).Msg("ProcessObjectivesUpload: dept wp error")
				}
			}

			// Group by division
			divGroups := groupByTwoFields(deptEntries,
				func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.DivObjName) },
				func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.Division) },
			)

			for _, divEntries := range divGroups {
				divFirst := divEntries[0]

				// Find division
				var division organogram.Division
				if err := tx.Where("LOWER(TRIM(division_name)) = ?",
					strings.ToLower(strings.TrimSpace(divFirst.Division))).
					First(&division).Error; err != nil {
					tx.Rollback()
					resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Division: %s", divFirst.Division))
					return resp, nil
				}

				// Find or create division objective
				var divObj performance.DivisionObjective
				divObjErr := tx.Where("LOWER(TRIM(name)) = ? AND division_id = ?",
					strings.ToLower(strings.TrimSpace(divFirst.DivObjName)), division.DivisionID).
					First(&divObj).Error
				if divObjErr != nil {
					divObjID, genErr := s.seqGen.GenerateCode(ctx, enums.SeqDivisionObjective, 15, "", enums.ConCatBefore)
					if genErr != nil {
						tx.Rollback()
						resp.Errors = append(resp.Errors, genErr.Error())
						return resp, nil
					}
					divObj = performance.DivisionObjective{
						DivisionObjectiveID:   divObjID,
						DivisionID:            division.DivisionID,
						DepartmentObjectiveID: deptObj.DepartmentObjectiveID,
						ObjectiveBase: domain.ObjectiveBase{
							Name:        strings.TrimSpace(divFirst.DivObjName),
							Description: strings.TrimSpace(divFirst.DivObjDesc),
							Kpi:         strings.TrimSpace(divFirst.DivObjKPI),
							Target:      strings.TrimSpace(divFirst.DivObjTarget),
							BaseWorkFlow: domain.BaseWorkFlow{
								BaseEntity: domain.BaseEntity{
									RecordStatus: enums.StatusActive.String(),
									Status:       enums.StatusApprovedAndActive.String(),
									IsActive:     true,
									CreatedBy:    s.userCtx.GetUserID(ctx),
								},
							},
						},
					}
					if createErr := tx.Create(&divObj).Error; createErr != nil {
						tx.Rollback()
						resp.Errors = append(resp.Errors, createErr.Error())
						return resp, nil
					}
				}

				// Save work product for division objective if provided
				divWpName := ""
				if divFirst.WorkProductName != nil {
					divWpName = *divFirst.WorkProductName
				}
				if strings.TrimSpace(divWpName) != "" {
					divWpDesc := ""
					if divFirst.WorkProductDescription != nil {
						divWpDesc = *divFirst.WorkProductDescription
					}
					divWpDeliv := ""
					if divFirst.WorkProductDeliverable != nil {
						divWpDeliv = *divFirst.WorkProductDeliverable
					}
					wpReq := &performance.WorkProductDefinitionVm{
						BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
						Name:           divWpName,
						Description:    divWpDesc,
						ObjectiveID:    divObj.DivisionObjectiveID,
						Deliverables:   divWpDeliv,
						SBUName:        strings.TrimSpace(divFirst.Office),
						Grade:          strings.TrimSpace(divFirst.JobGradeGroup),
						ObjectiveName:  strings.TrimSpace(divFirst.DivObjName),
						ObjectiveLevel: "Division",
					}
					if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
						s.log.Warn().Err(wpErr).Msg("ProcessObjectivesUpload: div wp error")
					}
				}

				// Group by office
				offGroups := groupByThreeFields(divEntries,
					func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.OffObjName) },
					func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.JobGradeGroup) },
					func(r performance.CascadedObjectiveUploadVm) string { return strings.TrimSpace(r.Office) },
				)

				for _, offEntries := range offGroups {
					offFirst := offEntries[0]

					// Find office
					var office organogram.Office
					if err := tx.Where("LOWER(TRIM(office_name)) = ?",
						strings.ToLower(strings.TrimSpace(offFirst.Office))).
						First(&office).Error; err != nil {
						tx.Rollback()
						resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for Office: %s", offFirst.Office))
						return resp, nil
					}

					// Find job grade group
					var jobGrade competency.JobGradeGroup
					if err := tx.Where("LOWER(group_name) = ?",
						strings.ToLower(strings.TrimSpace(offFirst.JobGradeGroup))).
						First(&jobGrade).Error; err != nil {
						tx.Rollback()
						resp.Errors = append(resp.Errors, fmt.Sprintf("No record found for JobGrade: %s", offFirst.JobGradeGroup))
						return resp, nil
					}

					// Find or create office objective
					var offObj performance.OfficeObjective
					offObjErr := tx.Where("LOWER(TRIM(name)) = ? AND job_grade_group_id = ? AND office_id = ?",
						strings.ToLower(strings.TrimSpace(offFirst.OffObjName)),
						jobGrade.JobGradeGroupID, office.OfficeID).
						First(&offObj).Error
					if offObjErr != nil {
						offObjID, genErr := s.seqGen.GenerateCode(ctx, enums.SeqOfficeObjective, 15, "", enums.ConCatBefore)
						if genErr != nil {
							tx.Rollback()
							resp.Errors = append(resp.Errors, genErr.Error())
							return resp, nil
						}
						offObj = performance.OfficeObjective{
							OfficeObjectiveID:   offObjID,
							OfficeID:            office.OfficeID,
							DivisionObjectiveID: divObj.DivisionObjectiveID,
							JobGradeGroupID:     jobGrade.JobGradeGroupID,
							ObjectiveBase: domain.ObjectiveBase{
								Name:        strings.TrimSpace(offFirst.OffObjName),
								Description: strings.TrimSpace(offFirst.OffObjDesc),
								Kpi:         strings.TrimSpace(offFirst.OffObjKPI),
								Target:      strings.TrimSpace(offFirst.OffObjTarget),
								BaseWorkFlow: domain.BaseWorkFlow{
									BaseEntity: domain.BaseEntity{
										RecordStatus: enums.StatusActive.String(),
										Status:       enums.StatusApprovedAndActive.String(),
										IsActive:     true,
										CreatedBy:    s.userCtx.GetUserID(ctx),
									},
								},
							},
						}
						if createErr := tx.Create(&offObj).Error; createErr != nil {
							tx.Rollback()
							resp.Errors = append(resp.Errors, createErr.Error())
							return resp, nil
						}
					}

					// Save work product for office objective if provided
					offWpName := ""
					if offFirst.WorkProductName != nil {
						offWpName = *offFirst.WorkProductName
					}
					if strings.TrimSpace(offWpName) != "" {
						offWpDesc := ""
						if offFirst.WorkProductDescription != nil {
							offWpDesc = *offFirst.WorkProductDescription
						}
						offWpDeliv := ""
						if offFirst.WorkProductDeliverable != nil {
							offWpDeliv = *offFirst.WorkProductDeliverable
						}
						wpReq := &performance.WorkProductDefinitionVm{
							BaseAuditVm:    performance.BaseAuditVm{Status: enums.StatusApprovedAndActive.String()},
							Name:           offWpName,
							Description:    offWpDesc,
							ObjectiveID:    offObj.OfficeObjectiveID,
							Deliverables:   offWpDeliv,
							SBUName:        strings.TrimSpace(offFirst.Office),
							Grade:          strings.TrimSpace(offFirst.JobGradeGroup),
							ObjectiveName:  strings.TrimSpace(offFirst.OffObjName),
							ObjectiveLevel: "Office",
						}
						if wpErr := s.workProductDefinitionSetup(ctx, wpReq, true); wpErr != nil {
							s.log.Warn().Err(wpErr).Msg("ProcessObjectivesUpload: off wp error")
						}
					}
				}
			}
		}
	}

	// Commit transaction
	if commitErr := tx.Commit().Error; commitErr != nil {
		resp.Errors = append(resp.Errors, commitErr.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// DeActivateOrReactivateObjectives
// ---------------------------------------------------------------------------

// DeActivateOrReactivateObjectives deactivates or reactivates objectives.
// For deactivation, it checks whether the objective is used in an active
// review period first. Mirrors .NET SMDService.DeActivateOrReactivateObjectives.
func (s *strategyService) DeActivateOrReactivateObjectives(ctx context.Context, req interface{}, deactivate bool) (interface{}, error) {
	objectives, ok := req.([]performance.ConsolidatedObjectiveVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for DeActivateOrReactivateObjectives")
	}

	resp := &performance.GenericResponseVm{}

	for _, obj := range objectives {
		objLevel := strings.ToLower(obj.ObjectiveLevel)

		switch objLevel {
		case "enterprise":
			entObj, err := s.getEnterpriseObjective(ctx, obj.ObjectiveID)
			if err != nil || entObj == nil {
				resp.Errors = append(resp.Errors, "Enterprise Objective not found")
				return resp, nil
			}
			if deactivate {
				if used, _ := s.isObjectiveUsedInActiveReviewPeriod(ctx, obj.ObjectiveID); used {
					resp.Errors = append(resp.Errors, "Objective cannot be deactivated, because it has been used in the current review period")
					return resp, nil
				}
				entObj.RecordStatus = enums.StatusDeactivated.String()
				entObj.Status = enums.StatusDeactivated.String()
			} else {
				entObj.RecordStatus = enums.StatusApprovedAndActive.String()
				entObj.Status = enums.StatusApprovedAndActive.String()
			}
			if err := s.db.WithContext(ctx).Save(entObj).Error; err != nil {
				resp.Errors = append(resp.Errors, err.Error())
				return resp, nil
			}

		case "department":
			deptObj, err := s.getDepartmentObjective(ctx, obj.ObjectiveID)
			if err != nil || deptObj == nil {
				resp.Errors = append(resp.Errors, "Department Objective not found")
				return resp, nil
			}
			if deactivate {
				if used, _ := s.isObjectiveUsedInActiveReviewPeriod(ctx, obj.ObjectiveID); used {
					resp.Errors = append(resp.Errors, "Objective cannot be deactivated, because it has been used in the current review period")
					return resp, nil
				}
				deptObj.RecordStatus = enums.StatusDeactivated.String()
				deptObj.Status = enums.StatusDeactivated.String()
			} else {
				deptObj.RecordStatus = enums.StatusApprovedAndActive.String()
				deptObj.Status = enums.StatusApprovedAndActive.String()
			}
			if err := s.db.WithContext(ctx).Save(deptObj).Error; err != nil {
				resp.Errors = append(resp.Errors, err.Error())
				return resp, nil
			}

		case "division":
			divObj, err := s.getDivisionObjective(ctx, obj.ObjectiveID)
			if err != nil || divObj == nil {
				resp.Errors = append(resp.Errors, "Division Objective not found")
				return resp, nil
			}
			if deactivate {
				if used, _ := s.isObjectiveUsedInActiveReviewPeriod(ctx, obj.ObjectiveID); used {
					resp.Errors = append(resp.Errors, "Objective cannot be deactivated, because it has been used in the current review period")
					return resp, nil
				}
				divObj.RecordStatus = enums.StatusDeactivated.String()
				divObj.Status = enums.StatusDeactivated.String()
			} else {
				divObj.RecordStatus = enums.StatusApprovedAndActive.String()
				divObj.Status = enums.StatusApprovedAndActive.String()
			}
			if err := s.db.WithContext(ctx).Save(divObj).Error; err != nil {
				resp.Errors = append(resp.Errors, err.Error())
				return resp, nil
			}

		case "office":
			offObj, err := s.getOfficeObjective(ctx, obj.ObjectiveID)
			if err != nil || offObj == nil {
				resp.Errors = append(resp.Errors, "Office Objective not found")
				return resp, nil
			}
			if deactivate {
				if used, _ := s.isObjectiveUsedInActiveReviewPeriod(ctx, obj.ObjectiveID); used {
					resp.Errors = append(resp.Errors, "Objective cannot be deactivated, because it has been used in the current review period")
					return resp, nil
				}
				offObj.RecordStatus = enums.StatusDeactivated.String()
				offObj.Status = enums.StatusDeactivated.String()
			} else {
				offObj.RecordStatus = enums.StatusApprovedAndActive.String()
				offObj.Status = enums.StatusApprovedAndActive.String()
			}
			if err := s.db.WithContext(ctx).Save(offObj).Error; err != nil {
				resp.Errors = append(resp.Errors, err.Error())
				return resp, nil
			}
		}
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// isObjectiveUsedInActiveReviewPeriod checks whether the given objective ID
// is associated with a planned objective in an active review period.
// Mirrors the .NET check using ReviewPeriodIndividualPlannedObjective.
func (s *strategyService) isObjectiveUsedInActiveReviewPeriod(ctx context.Context, objectiveID string) (bool, error) {
	var planned performance.ReviewPeriodIndividualPlannedObjective
	err := s.db.WithContext(ctx).
		Joins("JOIN pms.performance_review_periods rp ON rp.period_id = pms.review_period_individual_planned_objectives.review_period_id").
		Where("pms.review_period_individual_planned_objectives.objective_id = ? AND rp.record_status = ? AND pms.review_period_individual_planned_objectives.record_status NOT IN (?, ?, ?)",
			objectiveID,
			enums.StatusActive.String(),
			enums.StatusCancelled.String(),
			enums.StatusRejected.String(),
			enums.StatusReturned.String(),
		).
		First(&planned).Error
	if err != nil {
		return false, nil // not found = not used
	}
	return true, nil // found = used
}

// ---------------------------------------------------------------------------
// Save methods for Evaluation Options, Feedback, Work Product Definitions
// ---------------------------------------------------------------------------

// SaveEvaluationOptions saves (creates or updates) evaluation options.
// Mirrors .NET SMDService.ProcessUpload with EvaluationOptionVm items.
func (s *strategyService) SaveEvaluationOptions(ctx context.Context, req interface{}) (interface{}, error) {
	items, ok := req.([]performance.EvaluationOptionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveEvaluationOptions")
	}

	resp := &performance.GenericResponseVm{}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Errors = append(resp.Errors, tx.Error.Error())
		return resp, nil
	}

	for _, item := range items {
		isAdd := item.ID == 0
		if err := s.evaluationOptionsSetup(ctx, &item, isAdd); err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	if err := tx.Commit().Error; err != nil {
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// SaveFeedbackQuestionnaires saves (creates or updates) feedback questionnaires.
// Mirrors .NET SMDService.ProcessUpload with FeedbackQuestionaireVm items.
func (s *strategyService) SaveFeedbackQuestionnaires(ctx context.Context, req interface{}) (interface{}, error) {
	items, ok := req.([]performance.FeedbackQuestionaireVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveFeedbackQuestionnaires")
	}

	resp := &performance.GenericResponseVm{}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Errors = append(resp.Errors, tx.Error.Error())
		return resp, nil
	}

	for _, item := range items {
		isAdd := item.FeedbackQuestionaireID == ""
		if err := s.feedbackQuestionaireSetup(ctx, &item, isAdd); err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	if err := tx.Commit().Error; err != nil {
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// SaveFeedbackQuestionnaireOptions saves (creates or updates) feedback options.
// Mirrors .NET SMDService.ProcessUpload with FeedbackQuestionaireOptionVm items.
func (s *strategyService) SaveFeedbackQuestionnaireOptions(ctx context.Context, req interface{}) (interface{}, error) {
	items, ok := req.([]performance.FeedbackQuestionaireOptionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveFeedbackQuestionnaireOptions")
	}

	resp := &performance.GenericResponseVm{}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Errors = append(resp.Errors, tx.Error.Error())
		return resp, nil
	}

	for _, item := range items {
		isAdd := item.FeedbackQuestionaireOptionID == ""
		if err := s.feedbackQuestionaireOptionsSetup(ctx, &item, isAdd); err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	if err := tx.Commit().Error; err != nil {
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// SaveWorkProductDefinitions saves (creates or updates) work product definitions.
// Mirrors .NET SMDService.ProcessUpload with WorkProductDefinitionVm items.
func (s *strategyService) SaveWorkProductDefinitions(ctx context.Context, req interface{}) (interface{}, error) {
	items, ok := req.([]performance.WorkProductDefinitionVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type for SaveWorkProductDefinitions")
	}

	resp := &performance.GenericResponseVm{}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		resp.Errors = append(resp.Errors, tx.Error.Error())
		return resp, nil
	}

	for _, item := range items {
		isAdd := item.WorkProductDefinitionID == ""
		if err := s.workProductDefinitionSetup(ctx, &item, isAdd); err != nil {
			tx.Rollback()
			resp.Errors = append(resp.Errors, err.Error())
			return resp, nil
		}
	}

	if err := tx.Commit().Error; err != nil {
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.IsSuccess = true
	resp.Message = msgOperationCompleted
	return resp, nil
}

// ---------------------------------------------------------------------------
// Grouping helpers for bulk upload
// ---------------------------------------------------------------------------

// groupByField groups a slice by a single key extracted from each element.
func groupByField[T any](items []T, keyFn func(T) string) map[string][]T {
	result := make(map[string][]T)
	for _, item := range items {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}
	return result
}

// groupByTwoFields groups by a compound key of two fields.
func groupByTwoFields[T any](items []T, key1Fn, key2Fn func(T) string) map[string][]T {
	result := make(map[string][]T)
	for _, item := range items {
		key := key1Fn(item) + "|" + key2Fn(item)
		result[key] = append(result[key], item)
	}
	return result
}

// groupByThreeFields groups by a compound key of three fields.
func groupByThreeFields[T any](items []T, key1Fn, key2Fn, key3Fn func(T) string) map[string][]T {
	result := make(map[string][]T)
	for _, item := range items {
		key := key1Fn(item) + "|" + key2Fn(item) + "|" + key3Fn(item)
		result[key] = append(result[key], item)
	}
	return result
}
