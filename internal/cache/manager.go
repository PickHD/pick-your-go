// Package cache provides template caching functionality with TTL
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/PickHD/pick-your-go/internal/config"
)

const (
	// CacheTTL is the time-to-live for cache entries (24 hours)
	CacheTTL = 24 * time.Hour
	// CacheDirName is the name of the cache directory
	CacheDirName = ".pick-your-go"
	// CacheMetadataFile is the name of the cache metadata file
	CacheMetadataFile = "cache-metadata.json"
)

// CacheMetadata represents metadata for cached templates
type CacheMetadata struct {
	Templates map[string]TemplateCacheInfo `json:"templates"`
}

// TemplateCacheInfo stores cache information for a template
type TemplateCacheInfo struct {
	CachedAt    time.Time `json:"cached_at"`
	LastChecked time.Time `json:"last_checked"`
	Path        string    `json:"path"`
	Version     string    `json:"version,omitempty"`
}

// Manager handles template caching
type Manager struct {
	cacheDir string
	metadata *CacheMetadata
}

// NewManager creates a new cache manager
func NewManager() *Manager {
	cacheDir, _ := os.UserCacheDir()
	cacheDir = filepath.Join(cacheDir, CacheDirName)

	// Ensure cache directory exists
	os.MkdirAll(cacheDir, 0755)

	return &Manager{
		cacheDir: cacheDir,
		metadata: &CacheMetadata{
			Templates: make(map[string]TemplateCacheInfo),
		},
	}
}

// GetCacheDir returns the cache directory path
func (m *Manager) GetCacheDir() string {
	return m.cacheDir
}

// GetTemplateCachePath returns the cache path for a specific template
func (m *Manager) GetTemplateCachePath(archType config.ArchitectureType) string {
	return filepath.Join(m.cacheDir, string(archType))
}

// IsCached checks if a template is cached and still valid
func (m *Manager) IsCached(archType config.ArchitectureType) bool {
	// Load metadata
	if err := m.loadMetadata(); err != nil {
		return false
	}

	info, exists := m.metadata.Templates[string(archType)]
	if !exists {
		return false
	}

	// Check if cache is still valid (within TTL)
	return time.Since(info.CachedAt) < CacheTTL
}

// IsCacheExpired checks if the cache for a template has expired
func (m *Manager) IsCacheExpired(archType config.ArchitectureType) bool {
	if err := m.loadMetadata(); err != nil {
		return true
	}

	info, exists := m.metadata.Templates[string(archType)]
	if !exists {
		return true
	}

	return time.Since(info.CachedAt) >= CacheTTL
}

// UpdateCacheTime updates the cache time for a template
func (m *Manager) UpdateCacheTime(archType config.ArchitectureType) error {
	if err := m.loadMetadata(); err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	m.metadata.Templates[string(archType)] = TemplateCacheInfo{
		CachedAt:    time.Now(),
		LastChecked: time.Now(),
		Path:        m.GetTemplateCachePath(archType),
	}

	if err := m.saveMetadata(); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// GetCacheInfo returns cache information for a template
func (m *Manager) GetCacheInfo(archType config.ArchitectureType) (*TemplateCacheInfo, error) {
	if err := m.loadMetadata(); err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	info, exists := m.metadata.Templates[string(archType)]
	if !exists {
		return nil, fmt.Errorf("template not cached")
	}

	return &info, nil
}

// ClearCache removes all cached templates
func (m *Manager) ClearCache() error {
	// Remove all subdirectories in cache dir
	entries, err := os.ReadDir(m.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(m.cacheDir, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
		}
	}

	// Clear metadata
	m.metadata.Templates = make(map[string]TemplateCacheInfo)
	if err := m.saveMetadata(); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// ClearTemplateCache removes cache for a specific template
func (m *Manager) ClearTemplateCache(archType config.ArchitectureType) error {
	cachePath := m.GetTemplateCachePath(archType)

	// Remove template cache directory
	if _, err := os.Stat(cachePath); err == nil {
		if err := os.RemoveAll(cachePath); err != nil {
			return fmt.Errorf("failed to remove template cache: %w", err)
		}
	}

	// Remove from metadata
	if err := m.loadMetadata(); err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	delete(m.metadata.Templates, string(archType))

	if err := m.saveMetadata(); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// loadMetadata loads cache metadata from disk
func (m *Manager) loadMetadata() error {
	metadataPath := filepath.Join(m.cacheDir, CacheMetadataFile)

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No metadata file yet, that's ok
			m.metadata.Templates = make(map[string]TemplateCacheInfo)
			return nil
		}
		return fmt.Errorf("failed to read metadata file: %w", err)
	}

	if err := json.Unmarshal(data, m.metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	return nil
}

// saveMetadata saves cache metadata to disk
func (m *Manager) saveMetadata() error {
	metadataPath := filepath.Join(m.cacheDir, CacheMetadataFile)

	data, err := json.MarshalIndent(m.metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// GetCacheSize returns the size of the cache directory in bytes
func (m *Manager) GetCacheSize() (int64, error) {
	var size int64

	err := filepath.Walk(m.cacheDir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// GetCacheAge returns the age of the cache for a template
func (m *Manager) GetCacheAge(archType config.ArchitectureType) (time.Duration, error) {
	info, err := m.GetCacheInfo(archType)
	if err != nil {
		return 0, err
	}

	return time.Since(info.CachedAt), nil
}
