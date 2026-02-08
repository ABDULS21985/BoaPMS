package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// globalSettingService reads typed configuration values from the pms.settings table.
// This mirrors .NET's GlobalSetting service which provides GetBoolValue, GetStringValue, etc.
type globalSettingService struct {
	db  *gorm.DB
	log zerolog.Logger
}

func newGlobalSettingService(repos *repository.Container, log zerolog.Logger) GlobalSettingService {
	return &globalSettingService{
		db:  repos.GormDB,
		log: log,
	}
}

func (s *globalSettingService) GetBoolValue(ctx context.Context, key string) (bool, error) {
	val, err := s.getRawValue(ctx, key)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(val, "true") || val == "1", nil
}

func (s *globalSettingService) GetStringValue(ctx context.Context, key string) (string, error) {
	return s.getRawValue(ctx, key)
}

func (s *globalSettingService) GetIntValue(ctx context.Context, key string) (int, error) {
	val, err := s.getRawValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (s *globalSettingService) GetFloatValue(ctx context.Context, key string) (float64, error) {
	val, err := s.getRawValue(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}

func (s *globalSettingService) getRawValue(ctx context.Context, key string) (string, error) {
	var setting performance.Setting
	err := s.db.WithContext(ctx).
		Where("name = ? AND soft_deleted = false", key).
		First(&setting).Error
	if err != nil {
		return "", fmt.Errorf("setting %q not found: %w", key, err)
	}
	return setting.Value, nil
}
