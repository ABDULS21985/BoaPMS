package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// PmsSetupHandler handles PMS setup and configuration HTTP endpoints.
// Mirrors the .NET PmsSetupController with settings and configuration CRUD endpoints.
type PmsSetupHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewPmsSetupHandler creates a new PMS setup handler.
func NewPmsSetupHandler(svc *service.Container, log zerolog.Logger) *PmsSetupHandler {
	return &PmsSetupHandler{svc: svc, log: log}
}

// ============================================================
// Settings Endpoints
// ============================================================

// AddSetting handles POST /api/v1/setup/settings
// Mirrors .NET PerformanceMgtController.AddSetting -- creates a new global setting.
func (h *PmsSetupHandler) AddSetting(w http.ResponseWriter, r *http.Request) {
	var req performance.AddSettingRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Value == "" {
		response.Error(w, http.StatusBadRequest, "Name and value are required")
		return
	}

	if req.Type == "" {
		response.Error(w, http.StatusBadRequest, "Setting type is required")
		return
	}

	result, err := h.svc.PmsSetup.AddSetting(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddSetting").Msg("Failed to add setting")
		response.Error(w, http.StatusInternalServerError, "Failed to add setting")
		return
	}

	response.Created(w, result)
}

// UpdateSetting handles PUT /api/v1/setup/settings
// Mirrors .NET PerformanceMgtController.UpdateSetting -- updates an existing global setting.
func (h *PmsSetupHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	var req performance.SettingRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SettingID == "" {
		response.Error(w, http.StatusBadRequest, "Setting ID is required")
		return
	}

	if req.Name == "" || req.Value == "" {
		response.Error(w, http.StatusBadRequest, "Name and value are required")
		return
	}

	if req.Type == "" {
		response.Error(w, http.StatusBadRequest, "Setting type is required")
		return
	}

	result, err := h.svc.PmsSetup.UpdateSetting(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateSetting").Msg("Failed to update setting")
		response.Error(w, http.StatusInternalServerError, "Failed to update setting")
		return
	}

	response.OK(w, result)
}

// GetSettingDetails handles GET /api/v1/setup/settings/{settingId}
// Mirrors .NET PerformanceMgtController.GetSettingDetails -- returns a single setting by ID.
func (h *PmsSetupHandler) GetSettingDetails(w http.ResponseWriter, r *http.Request) {
	settingID := r.PathValue("settingId")
	if settingID == "" {
		response.Error(w, http.StatusBadRequest, "Setting ID is required")
		return
	}

	result, err := h.svc.PmsSetup.GetSettingDetails(r.Context(), settingID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetSettingDetails").Str("settingId", settingID).Msg("Failed to get setting details")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve setting details")
		return
	}

	response.OK(w, result)
}

// ListAllSettings handles GET /api/v1/setup/settings
// Mirrors .NET PerformanceMgtController.ListAllSettings -- returns all global settings.
func (h *PmsSetupHandler) ListAllSettings(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.PmsSetup.ListAllSettings(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "ListAllSettings").Msg("Failed to list settings")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve settings")
		return
	}

	response.OK(w, result)
}

// ============================================================
// PMS Configuration Endpoints
// ============================================================

// AddPmsConfiguration handles POST /api/v1/setup/pms-configurations
// Mirrors .NET PerformanceMgtController.AddPmsConfiguration -- creates a new PMS configuration entry.
func (h *PmsSetupHandler) AddPmsConfiguration(w http.ResponseWriter, r *http.Request) {
	var req performance.AddPmsConfigurationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Value == "" {
		response.Error(w, http.StatusBadRequest, "Name and value are required")
		return
	}

	if req.Type == "" {
		response.Error(w, http.StatusBadRequest, "Configuration type is required")
		return
	}

	result, err := h.svc.PmsSetup.AddPmsConfiguration(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddPmsConfiguration").Msg("Failed to add PMS configuration")
		response.Error(w, http.StatusInternalServerError, "Failed to add PMS configuration")
		return
	}

	response.Created(w, result)
}

// UpdatePmsConfiguration handles PUT /api/v1/setup/pms-configurations
// Mirrors .NET PerformanceMgtController.UpdatePmsConfiguration -- updates an existing PMS configuration.
func (h *PmsSetupHandler) UpdatePmsConfiguration(w http.ResponseWriter, r *http.Request) {
	var req performance.PmsConfigurationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PmsConfigurationID == "" {
		response.Error(w, http.StatusBadRequest, "Configuration ID is required")
		return
	}

	if req.Name == "" || req.Value == "" {
		response.Error(w, http.StatusBadRequest, "Name and value are required")
		return
	}

	if req.Type == "" {
		response.Error(w, http.StatusBadRequest, "Configuration type is required")
		return
	}

	result, err := h.svc.PmsSetup.UpdatePmsConfiguration(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdatePmsConfiguration").Msg("Failed to update PMS configuration")
		response.Error(w, http.StatusInternalServerError, "Failed to update PMS configuration")
		return
	}

	response.OK(w, result)
}

// GetPmsConfigurationDetails handles GET /api/v1/setup/pms-configurations/{configId}
// Mirrors .NET PerformanceMgtController.GetPmsConfigurationDetails -- returns a single configuration by ID.
func (h *PmsSetupHandler) GetPmsConfigurationDetails(w http.ResponseWriter, r *http.Request) {
	configID := r.PathValue("configId")
	if configID == "" {
		response.Error(w, http.StatusBadRequest, "Configuration ID is required")
		return
	}

	result, err := h.svc.PmsSetup.GetPmsConfigurationDetails(r.Context(), configID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPmsConfigurationDetails").Str("configId", configID).Msg("Failed to get PMS configuration details")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve PMS configuration details")
		return
	}

	response.OK(w, result)
}

// ListAllPmsConfigurations handles GET /api/v1/setup/pms-configurations
// Mirrors .NET PerformanceMgtController.ListAllPmsConfigurations -- returns all PMS configurations.
func (h *PmsSetupHandler) ListAllPmsConfigurations(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.PmsSetup.ListAllPmsConfigurations(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "ListAllPmsConfigurations").Msg("Failed to list PMS configurations")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve PMS configurations")
		return
	}

	response.OK(w, result)
}
