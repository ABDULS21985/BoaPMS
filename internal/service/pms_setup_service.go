package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// pmsSetupService implements the PmsSetupService interface.
// Converted from the .NET PmsSetupService class which handles CRUD for
// Settings (pms.settings) and PmsConfigurations (pms.pms_configurations).
// ---------------------------------------------------------------------------

type pmsSetupService struct {
	settingRepo *repository.PMSRepository[performance.Setting]
	configRepo  *repository.PMSRepository[performance.PmsConfiguration]
	seqGen      *sequenceGenerator
	encryption  EncryptionService
	log         zerolog.Logger
}

func newPmsSetupService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	encryption EncryptionService,
) PmsSetupService {
	return &pmsSetupService{
		settingRepo: repository.NewPMSRepository[performance.Setting](repos.GormDB),
		configRepo:  repository.NewPMSRepository[performance.PmsConfiguration](repos.GormDB),
		seqGen:      newSequenceGenerator(repos.GormDB, log),
		encryption:  encryption,
		log:         log.With().Str("service", "pms_setup").Logger(),
	}
}

// ==========================================================================
// Settings CRUD
// ==========================================================================

// AddSetting creates a new setting record after validating uniqueness, type,
// and optionally encrypting the value. Mirrors .NET PmsSetupService.AddSetting.
func (s *pmsSetupService) AddSetting(ctx context.Context, req interface{}) (interface{}, error) {
	addReq, ok := req.(*performance.AddSettingRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *AddSettingRequestModel")
	}

	// Check for duplicate name
	existing, err := s.settingRepo.FirstOrDefault(ctx, "LOWER(name) = ?", strings.ToLower(addReq.Name))
	if err != nil {
		s.log.Error().Err(err).Str("action", "ADD_SETTING").Msg("error checking existing setting")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_SETTING",
		}, nil
	}
	if existing != nil {
		s.log.Warn().Str("name", addReq.Name).Msg("setting already exists")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Setting: %s already exists", addReq.Name)},
			StatusCode:      "SET_EXT",
			ActionCall:      "ADD_SETTING",
		}, nil
	}

	// Validate setting type
	if !isValidSettingType(addReq.Type) {
		s.log.Warn().Str("type", addReq.Type).Msg("invalid setting type")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Type %s", addReq.Type)},
			StatusCode:      "SET_TYP_EXT_FL",
			ActionCall:      "ADD_SETTING",
		}, nil
	}

	// Encrypt value if required
	value := addReq.Value
	if addReq.IsEncrypted && s.encryption != nil {
		encrypted, encErr := s.encryption.Encrypt(value)
		if encErr != nil {
			s.log.Error().Err(encErr).Str("action", "ADD_SETTING").Msg("encryption failed")
			return &performance.SettingResponse{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "ADD_SETTING",
			}, nil
		}
		value = encrypted
	}

	// Generate setting ID from sequence (mirrors SequenceNumberTypes.Settings)
	settingID, err := s.seqGen.GenerateCode(ctx, enums.SeqGlobalSetting, 3, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Str("action", "ADD_SETTING").Msg("sequence generation failed")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_SETTING",
		}, nil
	}

	newSetting := &performance.Setting{
		SettingID:   settingID,
		Name:        addReq.Name,
		Value:       value,
		Type:        addReq.Type,
		IsEncrypted: addReq.IsEncrypted,
	}
	newSetting.IsActive = true

	if err := s.settingRepo.InsertAndSave(ctx, newSetting); err != nil {
		s.log.Error().Err(err).Str("action", "ADD_SETTING").Msg("failed to insert setting")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_SETTING",
		}, nil
	}

	detail := &performance.SettingResponseDetail{
		BaseEntityVm: toBaseEntityVm(newSetting.BaseEntity),
		SettingID:    newSetting.SettingID,
		Name:         newSetting.Name,
		Value:        addReq.Value, // return unencrypted value
		Type:         newSetting.Type,
		IsEncrypted:  newSetting.IsEncrypted,
	}

	s.log.Info().Str("settingId", settingID).Str("name", addReq.Name).Msg("setting created")
	return &performance.SettingResponse{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "SET_ADD",
		ActionCall:      "ADD_SETTING",
	}, nil
}

// UpdateSetting modifies an existing setting. Validates that the new name is
// unique (excluding the current record) and that the type is valid.
// Mirrors .NET PmsSetupService.UpdateSetting.
func (s *pmsSetupService) UpdateSetting(ctx context.Context, req interface{}) (interface{}, error) {
	updReq, ok := req.(*performance.SettingRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *SettingRequestModel")
	}

	// Fetch existing record
	setting, err := s.settingRepo.GetByStringID(ctx, "setting_id", updReq.SettingID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_SETTING").Msg("error fetching setting")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}
	if setting == nil {
		s.log.Warn().Str("settingId", updReq.SettingID).Msg("setting not found")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Id: %s provided", updReq.SettingID)},
			StatusCode:      "SET_EXT_FL",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}

	// Check name uniqueness (exclude current record)
	duplicate, err := s.settingRepo.FirstOrDefault(ctx,
		"LOWER(name) = ? AND setting_id != ?", strings.ToLower(updReq.Name), updReq.SettingID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_SETTING").Msg("error checking duplicate name")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}
	if duplicate != nil {
		s.log.Warn().Str("name", updReq.Name).Msg("duplicate setting name")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Setting: %s already exists", strings.ToUpper(updReq.Name))},
			StatusCode:      "SET_EXT",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}

	// Validate type
	if !isValidSettingType(updReq.Type) {
		s.log.Warn().Str("type", updReq.Type).Msg("invalid setting type")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Type %s", updReq.Type)},
			StatusCode:      "SET_TYP_EXT_FL",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}

	// Encrypt if required
	value := updReq.Value
	if updReq.IsEncrypted && s.encryption != nil {
		encrypted, encErr := s.encryption.Encrypt(value)
		if encErr != nil {
			s.log.Error().Err(encErr).Str("action", "UPDATE_SETTING").Msg("encryption failed")
			return &performance.SettingResponse{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "UPDATE_SETTING",
			}, nil
		}
		value = encrypted
	}

	// Apply changes
	if updReq.Name != "" {
		setting.Name = updReq.Name
	}
	if updReq.Type != "" {
		setting.Type = updReq.Type
	}
	setting.Value = value
	setting.IsEncrypted = updReq.IsEncrypted

	if err := s.settingRepo.UpdateAndSave(ctx, setting); err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_SETTING").Msg("failed to update setting")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_SETTING",
		}, nil
	}

	detail := &performance.SettingResponseDetail{
		BaseEntityVm: toBaseEntityVm(setting.BaseEntity),
		SettingID:    setting.SettingID,
		Name:         setting.Name,
		Value:        updReq.Value, // return unencrypted value
		Type:         setting.Type,
		IsEncrypted:  setting.IsEncrypted,
	}

	s.log.Info().Str("settingId", updReq.SettingID).Msg("setting updated")
	return &performance.SettingResponse{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "UPD_SET",
		ActionCall:      "UPDATE_SETTING",
	}, nil
}

// GetSettingDetails retrieves a single setting by its SettingID.
// Decrypts the value if the setting is marked as encrypted.
// Mirrors .NET PmsSetupService.GetSettingDetails.
func (s *pmsSetupService) GetSettingDetails(ctx context.Context, settingID string) (interface{}, error) {
	setting, err := s.settingRepo.GetByStringID(ctx, "setting_id", settingID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "GET_SETTING_DETAILS").Msg("error fetching setting")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "GET_SETTING_DETAILS",
		}, nil
	}
	if setting == nil {
		s.log.Warn().Str("settingId", settingID).Msg("setting not found")
		return &performance.SettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Id: %s provided", settingID)},
			StatusCode:      "SET_EXT_FL",
			ActionCall:      "GET_SETTING_DETAILS",
		}, nil
	}

	value := setting.Value
	if setting.IsEncrypted && s.encryption != nil {
		decrypted, decErr := s.encryption.Decrypt(value)
		if decErr != nil {
			s.log.Error().Err(decErr).Str("action", "GET_SETTING_DETAILS").Msg("decryption failed")
			return &performance.SettingResponse{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "GET_SETTING_DETAILS",
			}, nil
		}
		value = decrypted
	}

	detail := &performance.SettingResponseDetail{
		BaseEntityVm: toBaseEntityVm(setting.BaseEntity),
		SettingID:    setting.SettingID,
		Name:         setting.Name,
		Value:        value,
		Type:         setting.Type,
		IsEncrypted:  setting.IsEncrypted,
	}

	return &performance.SettingResponse{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "SET_EXT",
		ActionCall:      "GET_SETTING_DETAILS",
	}, nil
}

// ListAllSettings returns all settings, decrypting encrypted values.
// Mirrors .NET PmsSetupService.ListAllSettings.
func (s *pmsSetupService) ListAllSettings(ctx context.Context) (interface{}, error) {
	settings, err := s.settingRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Str("action", "LIST_ALL_SETTINGS").Msg("error listing settings")
		return &performance.ListSettingResponse{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "LIST_ALL_SETTINGS",
		}, nil
	}

	list := make([]performance.SettingResponseDetail, 0, len(settings))
	for _, st := range settings {
		value := st.Value
		if st.IsEncrypted && s.encryption != nil {
			decrypted, decErr := s.encryption.Decrypt(value)
			if decErr != nil {
				s.log.Warn().Err(decErr).Str("settingId", st.SettingID).Msg("failed to decrypt setting value, returning raw")
			} else {
				value = decrypted
			}
		}
		list = append(list, performance.SettingResponseDetail{
			BaseEntityVm: toBaseEntityVm(st.BaseEntity),
			SettingID:    st.SettingID,
			Name:         st.Name,
			Value:        value,
			Type:         st.Type,
			IsEncrypted:  st.IsEncrypted,
		})
	}

	return &performance.ListSettingResponse{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            list,
		TotalSettings:   len(list),
		StatusCode:      "LIS_SET",
		ActionCall:      "LIST_ALL_SETTINGS",
	}, nil
}

// ==========================================================================
// PmsConfiguration CRUD
// ==========================================================================

// AddPmsConfiguration creates a new PMS configuration record.
// Mirrors .NET PmsSetupService.AddPmsConfiguration.
func (s *pmsSetupService) AddPmsConfiguration(ctx context.Context, req interface{}) (interface{}, error) {
	addReq, ok := req.(*performance.AddPmsConfigurationRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *AddPmsConfigurationRequestModel")
	}

	// Check for duplicate name
	existing, err := s.configRepo.FirstOrDefault(ctx, "LOWER(name) = ?", strings.ToLower(addReq.Name))
	if err != nil {
		s.log.Error().Err(err).Str("action", "ADD_PMS_SETTING").Msg("error checking existing configuration")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_PMS_SETTING",
		}, nil
	}
	if existing != nil {
		s.log.Warn().Str("name", addReq.Name).Msg("PMS configuration already exists")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Setting: %s already exists", addReq.Name)},
			StatusCode:      "SET_EXT",
			ActionCall:      "ADD_PMS_SETTING",
		}, nil
	}

	// Validate type
	if !isValidSettingType(addReq.Type) {
		s.log.Warn().Str("type", addReq.Type).Msg("invalid setting type")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Type %s", addReq.Type)},
			StatusCode:      "SET_TYP_EXT_FL",
			ActionCall:      "ADD_PMS_SETTING",
		}, nil
	}

	// Encrypt if required
	value := addReq.Value
	if addReq.IsEncrypted && s.encryption != nil {
		encrypted, encErr := s.encryption.Encrypt(value)
		if encErr != nil {
			s.log.Error().Err(encErr).Str("action", "ADD_PMS_SETTING").Msg("encryption failed")
			return &performance.PmsConfigurationResponseVm{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "ADD_PMS_SETTING",
			}, nil
		}
		value = encrypted
	}

	// Generate PMS configuration ID from sequence (mirrors SequenceNumberTypes.PmsConfigurations)
	configID, err := s.seqGen.GenerateCode(ctx, enums.SeqConfigItem, 3, "", enums.ConCatBefore)
	if err != nil {
		s.log.Error().Err(err).Str("action", "ADD_PMS_SETTING").Msg("sequence generation failed")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_PMS_SETTING",
		}, nil
	}

	newConfig := &performance.PmsConfiguration{
		PmsConfigurationID: configID,
		Name:               addReq.Name,
		Value:              value,
		Type:               addReq.Type,
		IsEncrypted:        addReq.IsEncrypted,
	}
	newConfig.IsActive = true

	if err := s.configRepo.InsertAndSave(ctx, newConfig); err != nil {
		s.log.Error().Err(err).Str("action", "ADD_PMS_SETTING").Msg("failed to insert PMS configuration")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "ADD_PMS_SETTING",
		}, nil
	}

	detail := &performance.PmsConfigurationVm{
		BaseEntityVm:       toBaseEntityVm(newConfig.BaseEntity),
		PmsConfigurationID: newConfig.PmsConfigurationID,
		Name:               newConfig.Name,
		Value:              addReq.Value, // return unencrypted value
		Type:               newConfig.Type,
		IsEncrypted:        newConfig.IsEncrypted,
	}

	s.log.Info().Str("pmsConfigurationId", configID).Str("name", addReq.Name).Msg("PMS configuration created")
	return &performance.PmsConfigurationResponseVm{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "SET_ADD",
		ActionCall:      "ADD_PMS_SETTING",
	}, nil
}

// UpdatePmsConfiguration modifies an existing PMS configuration.
// Mirrors .NET PmsSetupService.UpdatePmsConfiguration.
func (s *pmsSetupService) UpdatePmsConfiguration(ctx context.Context, req interface{}) (interface{}, error) {
	updReq, ok := req.(*performance.PmsConfigurationRequestModel)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *PmsConfigurationRequestModel")
	}

	// Fetch existing
	cfg, err := s.configRepo.GetByStringID(ctx, "pms_configuration_id", updReq.PmsConfigurationID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_PMS_SETTING").Msg("error fetching configuration")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}
	if cfg == nil {
		s.log.Warn().Str("pmsConfigurationId", updReq.PmsConfigurationID).Msg("PMS configuration not found")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Id: %s provided", updReq.PmsConfigurationID)},
			StatusCode:      "SET_EXT_FL",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}

	// Check name uniqueness (exclude current record)
	duplicate, err := s.configRepo.FirstOrDefault(ctx,
		"LOWER(name) = ? AND pms_configuration_id != ?", strings.ToLower(updReq.Name), updReq.PmsConfigurationID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_PMS_SETTING").Msg("error checking duplicate name")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}
	if duplicate != nil {
		s.log.Warn().Str("name", updReq.Name).Msg("duplicate PMS configuration name")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Setting: %s already exists", strings.ToUpper(updReq.Name))},
			StatusCode:      "SET_EXT",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}

	// Validate type
	if !isValidSettingType(updReq.Type) {
		s.log.Warn().Str("type", updReq.Type).Msg("invalid setting type")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Type %s", updReq.Type)},
			StatusCode:      "SET_TYP_EXT_FL",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}

	// Encrypt if required
	value := updReq.Value
	if updReq.IsEncrypted && s.encryption != nil {
		encrypted, encErr := s.encryption.Encrypt(value)
		if encErr != nil {
			s.log.Error().Err(encErr).Str("action", "UPDATE_PMS_SETTING").Msg("encryption failed")
			return &performance.PmsConfigurationResponseVm{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "UPDATE_PMS_SETTING",
			}, nil
		}
		value = encrypted
	}

	// Apply changes
	if updReq.Name != "" {
		cfg.Name = updReq.Name
	}
	if updReq.Type != "" {
		cfg.Type = updReq.Type
	}
	cfg.Value = value
	cfg.IsEncrypted = updReq.IsEncrypted

	if err := s.configRepo.UpdateAndSave(ctx, cfg); err != nil {
		s.log.Error().Err(err).Str("action", "UPDATE_PMS_SETTING").Msg("failed to update PMS configuration")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "UPDATE_PMS_SETTING",
		}, nil
	}

	detail := &performance.PmsConfigurationVm{
		BaseEntityVm:       toBaseEntityVm(cfg.BaseEntity),
		PmsConfigurationID: cfg.PmsConfigurationID,
		Name:               cfg.Name,
		Value:              updReq.Value, // return unencrypted value
		Type:               cfg.Type,
		IsEncrypted:        cfg.IsEncrypted,
	}

	s.log.Info().Str("pmsConfigurationId", updReq.PmsConfigurationID).Msg("PMS configuration updated")
	return &performance.PmsConfigurationResponseVm{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "UPD_SET",
		ActionCall:      "UPDATE_PMS_SETTING",
	}, nil
}

// GetPmsConfigurationDetails retrieves a single PMS configuration by its ID.
// Decrypts the value if the configuration is marked as encrypted.
// Mirrors .NET PmsSetupService.GetPmsConfigurationDetails.
func (s *pmsSetupService) GetPmsConfigurationDetails(ctx context.Context, configID string) (interface{}, error) {
	cfg, err := s.configRepo.GetByStringID(ctx, "pms_configuration_id", configID)
	if err != nil {
		s.log.Error().Err(err).Str("action", "GET_PMS_SETTING_DETAILS").Msg("error fetching configuration")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "GET_PMS_SETTING_DETAILS",
		}, nil
	}
	if cfg == nil {
		s.log.Warn().Str("pmsConfigurationId", configID).Msg("PMS configuration not found")
		return &performance.PmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: fmt.Sprintf("Invalid Setting Id: %s provided", configID)},
			StatusCode:      "SET_EXT_FL",
			ActionCall:      "GET_PMS_SETTING_DETAILS",
		}, nil
	}

	value := cfg.Value
	if cfg.IsEncrypted && s.encryption != nil {
		decrypted, decErr := s.encryption.Decrypt(value)
		if decErr != nil {
			s.log.Error().Err(decErr).Str("action", "GET_PMS_SETTING_DETAILS").Msg("decryption failed")
			return &performance.PmsConfigurationResponseVm{
				BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
				StatusCode:      "EXPT",
				ActionCall:      "GET_PMS_SETTING_DETAILS",
			}, nil
		}
		value = decrypted
	}

	detail := &performance.PmsConfigurationVm{
		BaseEntityVm:       toBaseEntityVm(cfg.BaseEntity),
		PmsConfigurationID: cfg.PmsConfigurationID,
		Name:               cfg.Name,
		Value:              value,
		Type:               cfg.Type,
		IsEncrypted:        cfg.IsEncrypted,
	}

	return &performance.PmsConfigurationResponseVm{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            detail,
		StatusCode:      "SET_EXT",
		ActionCall:      "GET_PMS_SETTING_DETAILS",
	}, nil
}

// ListAllPmsConfigurations returns all PMS configurations, decrypting encrypted values.
// Mirrors .NET PmsSetupService.ListAllPmsConfigurations.
func (s *pmsSetupService) ListAllPmsConfigurations(ctx context.Context) (interface{}, error) {
	configs, err := s.configRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Str("action", "LIST_ALL_PMS_SETTINGS").Msg("error listing configurations")
		return &performance.ListPmsConfigurationResponseVm{
			BaseAPIResponse: performance.BaseAPIResponse{HasError: true, Message: msgGenericException},
			StatusCode:      "EXPT",
			ActionCall:      "LIST_ALL_PMS_SETTINGS",
		}, nil
	}

	list := make([]performance.PmsConfigurationVm, 0, len(configs))
	for _, c := range configs {
		value := c.Value
		if c.IsEncrypted && s.encryption != nil {
			decrypted, decErr := s.encryption.Decrypt(value)
			if decErr != nil {
				s.log.Warn().Err(decErr).Str("pmsConfigurationId", c.PmsConfigurationID).Msg("failed to decrypt config value, returning raw")
			} else {
				value = decrypted
			}
		}
		list = append(list, performance.PmsConfigurationVm{
			BaseEntityVm:       toBaseEntityVm(c.BaseEntity),
			PmsConfigurationID: c.PmsConfigurationID,
			Name:               c.Name,
			Value:              value,
			Type:               c.Type,
			IsEncrypted:        c.IsEncrypted,
		})
	}

	return &performance.ListPmsConfigurationResponseVm{
		BaseAPIResponse: performance.BaseAPIResponse{Message: msgOperationCompleted},
		Data:            list,
		TotalSettings:   len(list),
		StatusCode:      "LIS__PMS_SET",
		ActionCall:      "LIST_ALL_PMS_SETTINGS",
	}, nil
}

// ==========================================================================
// Helpers
// ==========================================================================

// toBaseEntityVm converts a domain.BaseEntity to the DTO BaseEntityVm.
func toBaseEntityVm(be domain.BaseEntity) performance.BaseEntityVm {
	return performance.BaseEntityVm{
		ID:           be.ID,
		RecordStatus: be.RecordStatus,
		CreatedAt:    be.CreatedAt,
		UpdatedAt:    be.UpdatedAt,
		CreatedBy:    be.CreatedBy,
		UpdatedBy:    be.UpdatedBy,
		IsActive:     be.IsActive,
	}
}

func init() {
	// Compile-time interface compliance check.
	var _ PmsSetupService = (*pmsSetupService)(nil)
}
