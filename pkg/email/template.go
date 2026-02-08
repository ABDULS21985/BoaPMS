package email

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// TemplateLoader loads HTML email templates from the filesystem.
// Templates are cached after the first read so that repeated lookups
// do not hit the disk.
type TemplateLoader struct {
	basePath string
	cache    map[string]string
	mu       sync.RWMutex
}

// NewTemplateLoader creates a new template loader with the given base path.
// The basePath should point to the directory containing the HTML template files
// (equivalent to the .NET wwwroot/EmailTemplates directory).
func NewTemplateLoader(basePath string) *TemplateLoader {
	return &TemplateLoader{
		basePath: basePath,
		cache:    make(map[string]string),
	}
}

// GetTemplate loads and caches a template from the basePath directory.
// Subsequent calls for the same template name return the cached version.
func (tl *TemplateLoader) GetTemplate(name string) (string, error) {
	// Fast path: check cache with read lock.
	tl.mu.RLock()
	if content, ok := tl.cache[name]; ok {
		tl.mu.RUnlock()
		return content, nil
	}
	tl.mu.RUnlock()

	// Slow path: read from disk and populate cache.
	tl.mu.Lock()
	defer tl.mu.Unlock()

	// Double-check after acquiring write lock.
	if content, ok := tl.cache[name]; ok {
		return content, nil
	}

	path := filepath.Join(tl.basePath, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("email template %q: %w", name, err)
	}

	content := string(data)
	tl.cache[name] = content
	return content, nil
}

// GetAccountCreationTemplate returns the account creation email template.
func (tl *TemplateLoader) GetAccountCreationTemplate() (string, error) {
	return tl.GetTemplate("AccountCreation.html")
}

// GetSignupTemplate returns the signup email template.
func (tl *TemplateLoader) GetSignupTemplate() (string, error) {
	return tl.GetTemplate("Signup.html")
}

// GetResetPasswordTemplate returns the password reset email template.
func (tl *TemplateLoader) GetResetPasswordTemplate() (string, error) {
	return tl.GetTemplate("ResetPassword.html")
}
