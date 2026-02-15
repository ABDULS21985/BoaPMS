package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// fileStorageService implements the FileStorageService interface.
// It stores files on the local filesystem under a configurable base directory
// read from Config.Storage.BasePath. Sub-directories are created automatically
// when saving files.
// ---------------------------------------------------------------------------

type fileStorageService struct {
	basePath string
	log      zerolog.Logger
}

// newFileStorageService creates a new FileStorageService backed by the local
// filesystem. The base storage directory is read from cfg.Storage.BasePath
// and created if it does not already exist.
func newFileStorageService(cfg *config.Config, log zerolog.Logger) FileStorageService {
	l := log.With().Str("service", "file_storage").Logger()

	basePath := cfg.Storage.BasePath
	if basePath == "" {
		basePath = "./uploads"
	}

	// Resolve to an absolute path for deterministic behaviour.
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		l.Error().Err(err).Str("basePath", basePath).Msg("failed to resolve absolute storage path, using raw value")
		absPath = basePath
	}

	// Ensure the base directory exists.
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		l.Error().Err(err).Str("basePath", absPath).Msg("failed to create storage base directory")
	}

	l.Info().Str("basePath", absPath).Msg("file storage service initialised")
	return &fileStorageService{basePath: absPath, log: l}
}

// SaveFile writes data to disk under the base storage directory and returns
// the full path of the persisted file. Parent directories are created
// automatically if the fileName contains path separators (e.g. "2024/report.pdf").
func (s *fileStorageService) SaveFile(_ context.Context, fileName string, data []byte) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file_storage: file name must not be empty")
	}

	// Clean the file name to prevent directory traversal.
	cleaned := filepath.Clean(fileName)
	fullPath := filepath.Join(s.basePath, cleaned)

	// Ensure the parent directory exists.
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		s.log.Error().Err(err).Str("dir", dir).Msg("failed to create parent directory")
		return "", fmt.Errorf("file_storage: creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		s.log.Error().Err(err).Str("path", fullPath).Msg("failed to write file")
		return "", fmt.Errorf("file_storage: writing file %s: %w", fullPath, err)
	}

	s.log.Info().Str("path", fullPath).Int("size", len(data)).Msg("file saved")
	return fullPath, nil
}

// DeleteFile removes a file from disk. The path may be absolute or relative
// to the base storage directory.
func (s *fileStorageService) DeleteFile(_ context.Context, path string) error {
	if path == "" {
		return fmt.Errorf("file_storage: path must not be empty")
	}

	fullPath := s.resolvePath(path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			s.log.Warn().Str("path", fullPath).Msg("file does not exist, nothing to delete")
			return nil
		}
		s.log.Error().Err(err).Str("path", fullPath).Msg("failed to delete file")
		return fmt.Errorf("file_storage: deleting file %s: %w", fullPath, err)
	}

	s.log.Info().Str("path", fullPath).Msg("file deleted")
	return nil
}

// GetFile reads and returns the contents of a file from disk.
func (s *fileStorageService) GetFile(_ context.Context, path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("file_storage: path must not be empty")
	}

	fullPath := s.resolvePath(path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.log.Warn().Str("path", fullPath).Msg("file not found")
			return nil, fmt.Errorf("file_storage: file not found %s: %w", fullPath, err)
		}
		s.log.Error().Err(err).Str("path", fullPath).Msg("failed to read file")
		return nil, fmt.Errorf("file_storage: reading file %s: %w", fullPath, err)
	}

	return data, nil
}

// resolvePath returns an absolute path. If the input is already absolute it is
// returned as-is; otherwise it is joined with the base storage directory.
func (s *fileStorageService) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(s.basePath, path)
}

func init() {
	// Compile-time interface compliance check.
	var _ FileStorageService = (*fileStorageService)(nil)
}
