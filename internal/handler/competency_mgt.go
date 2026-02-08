package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// CompetencyMgtHandler handles competency management HTTP endpoints.
// Mirrors the .NET CompetencyMgtController — route base: api/competencyMgt.
// All endpoints require [Authorize].
type CompetencyMgtHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewCompetencyMgtHandler creates a new competency management handler.
func NewCompetencyMgtHandler(svc *service.Container, log zerolog.Logger) *CompetencyMgtHandler {
	return &CompetencyMgtHandler{svc: svc, log: log}
}

// ---------------------------------------------------------------------------
// 1. GetCompetencies — POST (search)
// Mirrors .NET CompetencyMgtController.GetCompetencies
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencies(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.GetCompetencies(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competencies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 2. SaveCompetency — POST
// Mirrors .NET CompetencyMgtController.SaveCompetency
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetency(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 3. ApproveCompetency — POST
// Mirrors .NET CompetencyMgtController.ApproveCompetency
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) ApproveCompetency(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.ApproveCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to approve competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 4. RejectCompetency — POST
// Mirrors .NET CompetencyMgtController.RejectCompetency
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) RejectCompetency(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.RejectCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to reject competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 5. GetCompetencyCategories — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyCategories
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetCompetencyCategories(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency categories")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 6. SaveCompetencyCategory — POST
// Mirrors .NET CompetencyMgtController.SaveCompetencyCategory
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetencyCategory(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetencyCategory(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency category")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 7. GetCompetencyCategoryGradings — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyCategoryGradings
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyCategoryGradings(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetCompetencyCategoryGradings(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency category gradings")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 8. SaveCompetencyCategoryGrading — POST
// Mirrors .NET CompetencyMgtController.SaveCompetencyCategoryGrading
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetencyCategoryGrading(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetencyCategoryGrading(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency category grading")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 9. GetCompetencyRatingDefinitions — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyRatingDefinitions
// Query: competencyId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyRatingDefinitions(w http.ResponseWriter, r *http.Request) {
	var competencyId *int
	if cidStr := r.URL.Query().Get("competencyId"); cidStr != "" {
		cid, err := strconv.Atoi(cidStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid competencyId")
			return
		}
		competencyId = &cid
	}
	result, err := h.svc.Competency.GetCompetencyRatingDefinitions(r.Context(), competencyId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency rating definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 10. SaveCompetencyRatingDefinition — POST
// Mirrors .NET CompetencyMgtController.SaveCompetencyRatingDefinition
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetencyRatingDefinition(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetencyRatingDefinition(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency rating definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 11. GetCompetencyReviews — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyReviews
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyReviews(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetCompetencyReviews(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 12. GetCompetencyReviewByReviewer — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyReviewByReviewer
// Query: reviewerId (string), reviewPeriodId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyReviewByReviewer(w http.ResponseWriter, r *http.Request) {
	reviewerId := r.URL.Query().Get("reviewerId")
	if reviewerId == "" {
		response.Error(w, http.StatusBadRequest, "reviewerId is required")
		return
	}

	var reviewPeriodId *int
	if rpStr := r.URL.Query().Get("reviewPeriodId"); rpStr != "" {
		rp, err := strconv.Atoi(rpStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &rp
	}

	result, err := h.svc.Competency.GetCompetencyReviewByReviewer(r.Context(), reviewerId, reviewPeriodId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency review by reviewer")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 13. GetCompetencyReviewForEmployee — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyReviewForEmployee
// Query: employeeNumber (string), reviewPeriodId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyReviewForEmployee(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}

	var reviewPeriodId *int
	if rpStr := r.URL.Query().Get("reviewPeriodId"); rpStr != "" {
		rp, err := strconv.Atoi(rpStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &rp
	}

	result, err := h.svc.Competency.GetCompetencyReviewForEmployee(r.Context(), employeeNumber, reviewPeriodId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency review for employee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 14. GetCompetencyReviewDetail — POST (search)
// Mirrors .NET CompetencyMgtController.GetCompetencyReviewDetail
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyReviewDetail(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.GetCompetencyReviewDetail(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency review detail")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 15. SaveCompetencyReview — POST
// Mirrors .NET CompetencyMgtController.SaveCompetencyReview
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetencyReview(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetencyReview(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 16. GetCompetencyReviewProfiles — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyReviewProfiles
// Query: employeeNumber (string), reviewPeriodId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyReviewProfiles(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}

	var reviewPeriodId *int
	if rpStr := r.URL.Query().Get("reviewPeriodId"); rpStr != "" {
		rp, err := strconv.Atoi(rpStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &rp
	}

	result, err := h.svc.Competency.GetCompetencyReviewProfiles(r.Context(), employeeNumber, reviewPeriodId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency review profiles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 17. GetOfficeCompetencyReviews — GET
// Mirrors .NET CompetencyMgtController.GetOfficeCompetencyReviews
// Query: officeId (int), reviewPeriodId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetOfficeCompetencyReviews(w http.ResponseWriter, r *http.Request) {
	officeIdStr := r.URL.Query().Get("officeId")
	if officeIdStr == "" {
		response.Error(w, http.StatusBadRequest, "officeId is required")
		return
	}
	officeId, err := strconv.Atoi(officeIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid officeId")
		return
	}

	var reviewPeriodId *int
	if rpStr := r.URL.Query().Get("reviewPeriodId"); rpStr != "" {
		rp, err := strconv.Atoi(rpStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &rp
	}

	result, err := h.svc.Competency.GetOfficeCompetencyReviews(r.Context(), officeId, reviewPeriodId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get office competency reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 18. GetGroupCompetencyReviewProfiles — GET
// Mirrors .NET CompetencyMgtController.GetGroupCompetencyReviewProfiles
// Query: reviewPeriodId (opt int), officeId (opt int), divisionId (opt int), departmentId (opt int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetGroupCompetencyReviewProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var reviewPeriodId *int
	if v := q.Get("reviewPeriodId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &id
	}

	var officeId *int
	if v := q.Get("officeId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid officeId")
			return
		}
		officeId = &id
	}

	var divisionId *int
	if v := q.Get("divisionId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid divisionId")
			return
		}
		divisionId = &id
	}

	var departmentId *int
	if v := q.Get("departmentId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid departmentId")
			return
		}
		departmentId = &id
	}

	result, err := h.svc.Competency.GetGroupCompetencyReviewProfiles(r.Context(), reviewPeriodId, officeId, divisionId, departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get group competency review profiles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 19. GetCompetencyMatrixReviewProfiles — GET
// Mirrors .NET CompetencyMgtController.GetCompetencyMatrixReviewProfiles
// Query: reviewPeriodId (opt int), officeId (opt int), divisionId (opt int), departmentId (opt int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetCompetencyMatrixReviewProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var reviewPeriodId *int
	if v := q.Get("reviewPeriodId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &id
	}

	var officeId *int
	if v := q.Get("officeId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid officeId")
			return
		}
		officeId = &id
	}

	var divisionId *int
	if v := q.Get("divisionId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid divisionId")
			return
		}
		divisionId = &id
	}

	var departmentId *int
	if v := q.Get("departmentId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid departmentId")
			return
		}
		departmentId = &id
	}

	result, err := h.svc.Competency.GetCompetencyMatrixReviewProfiles(r.Context(), reviewPeriodId, officeId, divisionId, departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get competency matrix review profiles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 20. GetTechnicalCompetencyMatrixReviewProfiles — GET
// Mirrors .NET CompetencyMgtController.GetTechnicalCompetencyMatrixReviewProfiles
// Query: reviewPeriodId (opt int), jobRoleId (int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetTechnicalCompetencyMatrixReviewProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var reviewPeriodId *int
	if v := q.Get("reviewPeriodId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid reviewPeriodId")
			return
		}
		reviewPeriodId = &id
	}

	jobRoleIdStr := q.Get("jobRoleId")
	if jobRoleIdStr == "" {
		response.Error(w, http.StatusBadRequest, "jobRoleId is required")
		return
	}
	jobRoleId, err := strconv.Atoi(jobRoleIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid jobRoleId")
		return
	}

	result, err := h.svc.Competency.GetTechnicalCompetencyMatrixReviewProfiles(r.Context(), reviewPeriodId, jobRoleId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get technical competency matrix review profiles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 21. SaveCompetencyReviewProfile — POST
// Mirrors .NET CompetencyMgtController.SaveCompetencyReviewProfile
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveCompetencyReviewProfile(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveCompetencyReviewProfile(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save competency review profile")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 22. GetDevelopmentPlans — GET
// Mirrors .NET CompetencyMgtController.GetDevelopmentPlans
// Query: competencyProfileReviewId (optional int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetDevelopmentPlans(w http.ResponseWriter, r *http.Request) {
	var competencyProfileReviewId *int
	if v := r.URL.Query().Get("competencyProfileReviewId"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid competencyProfileReviewId")
			return
		}
		competencyProfileReviewId = &id
	}
	result, err := h.svc.Competency.GetDevelopmentPlans(r.Context(), competencyProfileReviewId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get development plans")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 23. SaveDevelopmentPlan — POST
// Mirrors .NET CompetencyMgtController.SaveDevelopmentPlan
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveDevelopmentPlan(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveDevelopmentPlan(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save development plan")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 24. GetJobRoles — GET
// Mirrors .NET CompetencyMgtController.GetJobRoles
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetJobRoles(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetJobRoles(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job roles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 25. SaveJobRole — POST
// Mirrors .NET CompetencyMgtController.SaveJobRole
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveJobRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveJobRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save job role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 26. GetOfficeJobRoles — POST (search)
// Mirrors .NET CompetencyMgtController.GetOfficeJobRoles
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetOfficeJobRoles(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.GetOfficeJobRoles(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get office job roles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 27. SaveOfficeJobRole — POST
// Mirrors .NET CompetencyMgtController.SaveOfficeJobRole
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveOfficeJobRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveOfficeJobRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save office job role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 28. GetJobRoleCompetencies — POST (search)
// Mirrors .NET CompetencyMgtController.GetJobRoleCompetencies
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetJobRoleCompetencies(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.GetJobRoleCompetencies(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job role competencies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 29. SaveJobRoleCompetency — POST
// Mirrors .NET CompetencyMgtController.SaveJobRoleCompetency
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveJobRoleCompetency(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveJobRoleCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save job role competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 30. GetBehavioralCompetencies — GET
// Mirrors .NET CompetencyMgtController.GetBehavioralCompetencies
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetBehavioralCompetencies(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetBehavioralCompetencies(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get behavioral competencies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 31. SaveBehavioralCompetency — POST
// Mirrors .NET CompetencyMgtController.SaveBehavioralCompetency
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveBehavioralCompetency(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveBehavioralCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save behavioral competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 32. GetJobRoleGrades — GET
// Mirrors .NET CompetencyMgtController.GetJobRoleGrades
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetJobRoleGrades(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetJobRoleGrades(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job role grades")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 33. SaveJobRoleGrade — POST
// Mirrors .NET CompetencyMgtController.SaveJobRoleGrade
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveJobRoleGrade(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveJobRoleGrade(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save job role grade")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 34. GetJobGrades — GET
// Mirrors .NET CompetencyMgtController.GetJobGrades
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetJobGrades(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetJobGrades(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job grades")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 35. SaveJobGrade — POST
// Mirrors .NET CompetencyMgtController.SaveJobGrade
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveJobGrade(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveJobGrade(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save job grade")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 36. GetJobGradeGroups — GET
// Mirrors .NET CompetencyMgtController.GetJobGradeGroups
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetJobGradeGroups(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetJobGradeGroups(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job grade groups")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 37. SaveJobGradeGroup — POST
// Mirrors .NET CompetencyMgtController.SaveJobGradeGroup
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveJobGradeGroup(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveJobGradeGroup(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save job grade group")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 38. GetAssignJobGradeGroups — GET
// Mirrors .NET CompetencyMgtController.GetAssignJobGradeGroups
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetAssignJobGradeGroups(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetAssignJobGradeGroups(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get assign job grade groups")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 39. SaveAssignJobGradeGroup — POST
// Mirrors .NET CompetencyMgtController.SaveAssignJobGradeGroup
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveAssignJobGradeGroup(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveAssignJobGradeGroup(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save assign job grade group")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 40. GetRatings — GET
// Mirrors .NET CompetencyMgtController.GetRatings
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetRatings(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetRatings(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get ratings")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 41. SaveRating — POST
// Mirrors .NET CompetencyMgtController.SaveRating
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveRating(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveRating(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save rating")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 42. GetReviewPeriods — GET
// Mirrors .NET CompetencyMgtController.GetReviewPeriods
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetReviewPeriods(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetReviewPeriods(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get review periods")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 43. SaveReviewPeriod — POST
// Mirrors .NET CompetencyMgtController.SaveReviewPeriod
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveReviewPeriod(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 44. ApproveReviewPeriod — POST
// Mirrors .NET CompetencyMgtController.ApproveReviewPeriod
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) ApproveReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.ApproveReviewPeriod(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to approve review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 45. GetReviewTypes — GET
// Mirrors .NET CompetencyMgtController.GetReviewTypes
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetReviewTypes(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetReviewTypes(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get review types")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 46. SaveReviewType — POST
// Mirrors .NET CompetencyMgtController.SaveReviewType
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveReviewType(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveReviewType(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save review type")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 47. GetBankYears — GET
// Mirrors .NET CompetencyMgtController.GetBankYears
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetBankYears(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.GetBankYears(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get bank years")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 48. SaveBankYear — POST
// Mirrors .NET CompetencyMgtController.SaveBankYear
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveBankYear(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveBankYear(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save bank year")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 49. GetTrainingTypes — GET
// Mirrors .NET CompetencyMgtController.GetTrainingTypes
// Query: isActive (optional bool)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) GetTrainingTypes(w http.ResponseWriter, r *http.Request) {
	var isActive *bool
	if v := r.URL.Query().Get("isActive"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid isActive")
			return
		}
		isActive = &b
	}
	result, err := h.svc.Competency.GetTrainingTypes(r.Context(), isActive)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get training types")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 50. SaveTrainingType — POST
// Mirrors .NET CompetencyMgtController.SaveTrainingType
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SaveTrainingType(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SaveTrainingType(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save training type")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 51. PopulateAllReviews — GET
// Mirrors .NET CompetencyMgtController.PopulateAllReviews
// (background job in .NET, regular call in Go)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) PopulateAllReviews(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Competency.PopulateAllReviews(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to populate all reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 52. PopulateOfficeReviews — GET
// Mirrors .NET CompetencyMgtController.PopulateOfficeReviews
// Query: officeId (int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) PopulateOfficeReviews(w http.ResponseWriter, r *http.Request) {
	officeIdStr := r.URL.Query().Get("officeId")
	if officeIdStr == "" {
		response.Error(w, http.StatusBadRequest, "officeId is required")
		return
	}
	officeId, err := strconv.Atoi(officeIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid officeId")
		return
	}
	result, err := h.svc.Competency.PopulateOfficeReviews(r.Context(), officeId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to populate office reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 53. PopulateDivisionReviews — GET
// Mirrors .NET CompetencyMgtController.PopulateDivisionReviews
// Query: divisionId (int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) PopulateDivisionReviews(w http.ResponseWriter, r *http.Request) {
	divisionIdStr := r.URL.Query().Get("divisionId")
	if divisionIdStr == "" {
		response.Error(w, http.StatusBadRequest, "divisionId is required")
		return
	}
	divisionId, err := strconv.Atoi(divisionIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid divisionId")
		return
	}
	result, err := h.svc.Competency.PopulateDivisionReviews(r.Context(), divisionId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to populate division reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 54. PopulateDepartmentReviews — GET
// Mirrors .NET CompetencyMgtController.PopulateDepartmentReviews
// Query: departmentId (int)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) PopulateDepartmentReviews(w http.ResponseWriter, r *http.Request) {
	departmentIdStr := r.URL.Query().Get("departmentId")
	if departmentIdStr == "" {
		response.Error(w, http.StatusBadRequest, "departmentId is required")
		return
	}
	departmentId, err := strconv.Atoi(departmentIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid departmentId")
		return
	}
	result, err := h.svc.Competency.PopulateDepartmentReviews(r.Context(), departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to populate department reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 55. PopulateReviewsByEmployeeId — GET
// Mirrors .NET CompetencyMgtController.PopulateReviewsByEmployeeId
// Query: employeeNumber (string)
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) PopulateReviewsByEmployeeId(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.svc.Competency.PopulateReviewsByEmployeeId(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to populate reviews by employee id")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 56. CalculateReviews — POST
// Mirrors .NET CompetencyMgtController.CalculateReviews
// Body: {isTechnical bool, employeeNumber string, reviewPeriodId int}
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) CalculateReviews(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.CalculateReviews(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to calculate reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 57. RecalculateReviewsProfiles — POST
// Mirrors .NET CompetencyMgtController.RecalculateReviewsProfiles
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) RecalculateReviewsProfiles(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.RecalculateReviewsProfiles(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to recalculate reviews profiles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 58. EmailService — POST
// Mirrors .NET CompetencyMgtController.EmailService
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) EmailService(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.EmailService(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to send email")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// 59. SyncJobRoleUpdateSOA — POST
// Mirrors .NET CompetencyMgtController.SyncJobRoleUpdateSOA
// ---------------------------------------------------------------------------

func (h *CompetencyMgtHandler) SyncJobRoleUpdateSOA(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.svc.Competency.SyncJobRoleUpdateSOA(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to sync job role update SOA")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}
