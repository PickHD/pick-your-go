// Package template provides template management and GitHub integration
package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PickHD/pick-your-go/internal/cache"
	"github.com/PickHD/pick-your-go/internal/config"
)

// Template represents a project template
type Template struct {
	Type        config.ArchitectureType `json:"type"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Repository  string                  `json:"repository"`
	Branch      string                  `json:"branch"`
}

// Manager handles template operations
type Manager struct {
	cacheManager *cache.Manager
	templates    []*Template
}

// NewManager creates a new template manager
func NewManager() *Manager {
	m := &Manager{
		cacheManager: cache.NewManager(),
		templates:    getDefaultTemplates(),
	}
	return m
}

// getDefaultTemplates returns the default template definitions
func getDefaultTemplates() []*Template {
	return []*Template{
		{
			Type:        config.LayeredArchitecture,
			Name:        "Layered Architecture Template",
			Description: "Traditional layered architecture with clear separation between presentation, business logic, and data layers",
			Repository:  "https://github.com/PickHD/go-layered-template.git",
			Branch:      "main",
		},
		{
			Type:        config.ModularArchitecture,
			Name:        "Modular Architecture Template",
			Description: "Modular monolithic architecture with domain-driven design principles",
			Repository:  "https://github.com/PickHD/go-modular-template.git",
			Branch:      "main",
		},
		{
			Type:        config.HexagonalArchitecture,
			Name:        "Hexagonal Architecture Template",
			Description: "Hexagonal architecture (ports and adapters) with isolation of core logic from external concerns",
			Repository:  "https://github.com/PickHD/go-hexagonal-template.git",
			Branch:      "main",
		},
	}
}

// GetTemplates returns all available templates
func (m *Manager) GetTemplates() ([]*Template, error) {
	return m.templates, nil
}

// GetTemplate returns a template by architecture type
func (m *Manager) GetTemplate(archType config.ArchitectureType) (*Template, error) {
	for _, tmpl := range m.templates {
		if tmpl.Type == archType {
			return tmpl, nil
		}
	}
	return nil, fmt.Errorf("template not found for architecture type: %s", archType)
}

// IsCached checks if a template is cached
func (m *Manager) IsCached(archType config.ArchitectureType) bool {
	return m.cacheManager.IsCached(archType)
}

// GetTemplatePath returns the path to a cached template
func (m *Manager) GetTemplatePath(archType config.ArchitectureType) (string, error) {
	if !m.IsCached(archType) {
		return "", fmt.Errorf("template not cached: %s", archType)
	}

	return m.cacheManager.GetTemplateCachePath(archType), nil
}

// UpdateTemplate downloads or updates a template from GitHub
func (m *Manager) UpdateTemplate(archType config.ArchitectureType, token string) error {
	template, err := m.GetTemplate(archType)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	cachePath := m.cacheManager.GetTemplateCachePath(archType)

	// Check if template directory already exists
	if _, err := os.Stat(cachePath); err == nil {
		// Directory exists, pull latest changes
		if err := m.pullTemplate(cachePath, token); err != nil {
			// If pull fails, try cloning fresh
			if err := os.RemoveAll(cachePath); err != nil {
				return fmt.Errorf("failed to remove old cache: %w", err)
			}
			if err := m.cloneTemplate(template, cachePath, token); err != nil {
				return err
			}
		}
	} else {
		// Directory doesn't exist, clone it
		if err := m.cloneTemplate(template, cachePath, token); err != nil {
			return err
		}
	}

	// Update cache metadata AFTER successful clone/pull
	return m.cacheManager.UpdateCacheTime(archType)
}

// cloneTemplate clones a template repository from GitHub
func (m *Manager) cloneTemplate(template *Template, cachePath string, token string) error {
	// Ensure parent directory exists
	parentDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Build git clone command with token authentication
	repoURL := m.buildAuthenticatedURL(template.Repository, token)

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", template.Branch, repoURL, cachePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Remove .git directory to save space
	gitDir := filepath.Join(cachePath, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		// Not a critical error, just log it
		fmt.Printf("Warning: failed to remove .git directory: %v\n", err)
	}

	return nil
}

// pullTemplate pulls latest changes for a cached template
func (m *Manager) pullTemplate(cachePath string, token string) error {
	// We need to re-initialize git to pull, since we removed .git
	// So it's easier to just re-clone
	return fmt.Errorf("pull not supported, please re-clone")
}

// buildAuthenticatedURL creates a GitHub URL with token authentication
func (m *Manager) buildAuthenticatedURL(repoURL string, token string) string {
	// Parse the URL to insert the token
	// Format: https://TOKEN@github.com/user/repo.git

	if token == "" {
		return repoURL
	}

	// Remove https:// prefix if present
	url := strings.TrimPrefix(repoURL, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Build authenticated URL
	return fmt.Sprintf("https://%s@%s", token, url)
}

// EnsureTemplateCached ensures a template is cached, downloading if necessary
func (m *Manager) EnsureTemplateCached(archType config.ArchitectureType, token string) error {
	// Check if already cached and valid
	if m.IsCached(archType) {
		return nil
	}

	// Download the template
	return m.UpdateTemplate(archType, token)
}

// GetTemplateFiles returns a list of files in a cached template
func (m *Manager) GetTemplateFiles(archType config.ArchitectureType) ([]string, error) {
	cachePath, err := m.GetTemplatePath(archType)
	if err != nil {
		return nil, err
	}

	var files []string

	err = filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the cache directory itself and .git
		if path == cachePath {
			return nil
		}

		if info.IsDir() {
			// Skip .git directory
			if filepath.Base(path) == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(cachePath, path)
		if err != nil {
			return err
		}

		files = append(files, relPath)
		return nil
	})

	return files, err
}

// CopyTemplateToDestination copies a template to a destination directory
func (m *Manager) CopyTemplateToDestination(archType config.ArchitectureType, destPath string) error {
	// CRITICAL: Ensure destPath is absolute to avoid path resolution issues
	if !filepath.IsAbs(destPath) {
		return fmt.Errorf("BUG: destPath is not absolute: %s", destPath)
	}

	cachePath, err := m.GetTemplatePath(archType)
	if err != nil {
		return fmt.Errorf("failed to get template path: %w", err)
	}


	// Copy all files from cache to destination
	return filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the cache root directory
		if path == cachePath {
			return nil
		}

		// Skip .git directory
		if filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}

		// Calculate destination path
		relPath, err := filepath.Rel(cachePath, path)
		if err != nil {
			return err
		}
		// BUG FIX: Use different variable name to avoid shadowing the destPath parameter
		// This was causing incorrect path resolution
		targetPath := filepath.Join(destPath, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Copy file
		return copyFile(path, targetPath, info.Mode())
	})
}

// copyFile copies a file from src to dst
func copyFile(src, dst string, mode os.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", src, err)
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	if err := os.WriteFile(dst, data, mode); err != nil {
		return fmt.Errorf("failed to write destination file %s: %w", dst, err)
	}

	return nil
}

// ClearAllCache clears all cached templates
func (m *Manager) ClearAllCache() error {
	return m.cacheManager.ClearCache()
}

// ClearTemplateCache clears cache for a specific template
func (m *Manager) ClearTemplateCache(archType config.ArchitectureType) error {
	return m.cacheManager.ClearTemplateCache(archType)
}
