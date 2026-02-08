package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Result represents the outcome of a file storage operation.
type Result struct {
	Status   bool   `json:"status"`
	FullPath string `json:"full_path"`
	FileName string `json:"file_name"`
	Message  string `json:"message"`
}

// FileStorage handles local file operations rooted at a base directory.
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new file storage with the given base path.
// The base directory is created if it does not already exist.
func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

// SaveFile saves an uploaded file with the given filename to the base directory.
func (fs *FileStorage) SaveFile(reader io.Reader, fileName string) (*Result, error) {
	return fs.saveToDir(reader, fileName, fs.basePath)
}

// SaveToSubdir saves a file to a subdirectory under the base path.
// The subdirectory is created if it does not already exist.
func (fs *FileStorage) SaveToSubdir(reader io.Reader, fileName, subdir string) (*Result, error) {
	dir := filepath.Join(fs.basePath, subdir)
	return fs.saveToDir(reader, fileName, dir)
}

// DeleteFile removes a file from the base path.
func (fs *FileStorage) DeleteFile(fileName string) error {
	fullPath := filepath.Join(fs.basePath, fileName)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // already gone; treat as success
		}
		return fmt.Errorf("delete file %q: %w", fileName, err)
	}
	return nil
}

// UpdateFile replaces an existing file (deletes old, saves new).
// If the old file cannot be deleted the operation still proceeds with the new file.
func (fs *FileStorage) UpdateFile(reader io.Reader, oldFileName, newFileName string) (*Result, error) {
	// Best-effort removal of the old file.
	_ = fs.DeleteFile(oldFileName)
	return fs.SaveFile(reader, newFileName)
}

// FileExists checks if a file exists in the base path.
func (fs *FileStorage) FileExists(fileName string) bool {
	fullPath := filepath.Join(fs.basePath, fileName)
	info, err := os.Stat(fullPath)
	return err == nil && !info.IsDir()
}

// saveToDir is the shared helper that writes reader content to dir/fileName.
func (fs *FileStorage) saveToDir(reader io.Reader, fileName, dir string) (*Result, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &Result{
			Status:  false,
			Message: fmt.Sprintf("failed to create directory: %v", err),
		}, fmt.Errorf("create directory %q: %w", dir, err)
	}

	fullPath := filepath.Join(dir, fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		return &Result{
			Status:  false,
			Message: fmt.Sprintf("failed to create file: %v", err),
		}, fmt.Errorf("create file %q: %w", fullPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		// Clean up the partially written file.
		os.Remove(fullPath)
		return &Result{
			Status:  false,
			Message: fmt.Sprintf("failed to write file: %v", err),
		}, fmt.Errorf("write file %q: %w", fullPath, err)
	}

	return &Result{
		Status:   true,
		FullPath: fullPath,
		FileName: fileName,
		Message:  "File saved successfully",
	}, nil
}
