// Package generator provides the core template generation logic
package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PickHD/pick-your-go/internal/config"
)

// Generator defines the interface for architecture-specific generators
// Each architecture pattern implements this interface with its own structure
type Generator interface {
	// Generate creates the project structure based on the architecture
	Generate(cfg *config.Config) error
	// Validate checks if the configuration is valid for this architecture
	Validate(cfg *config.Config) error
	// GetStructure returns the directory structure that will be created
	GetStructure() []string
}

// GeneratorFactory creates generators based on architecture type
type GeneratorFactory struct{}

// NewGeneratorFactory creates a new generator factory
func NewGeneratorFactory() *GeneratorFactory {
	return &GeneratorFactory{}
}

// CreateGenerator returns a generator for the specified architecture type
func (f *GeneratorFactory) CreateGenerator(archType config.ArchitectureType) (Generator, error) {
	switch archType {
	case config.LayeredArchitecture:
		return NewLayeredGenerator(), nil
	case config.ModularArchitecture:
		return NewModularGenerator(), nil
	case config.HexagonalArchitecture:
		return NewHexagonalGenerator(), nil
	default:
		return nil, fmt.Errorf("unsupported architecture type: %s", archType)
	}
}

// BaseGenerator provides common functionality for all generators
type BaseGenerator struct {
	createFile func(path string, content string) error
	createDir  func(path string) error
}

// NewBaseGenerator creates a new base generator with real filesystem operations
func NewBaseGenerator() *BaseGenerator {
	return &BaseGenerator{
		createFile: writeFile,
		createDir:  createDirectory,
	}
}

// writeFile writes content to a file
func writeFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// createDirectory creates a directory with all parent directories
func createDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// CreateDirectory creates a directory
func (b *BaseGenerator) CreateDirectory(path string) error {
	return b.createDir(path)
}

// CreateFile creates a file with content
func (b *BaseGenerator) CreateFile(path string, content string) error {
	return b.createFile(path, content)
}

// ValidateConfig performs common validation
func (b *BaseGenerator) ValidateConfig(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

// GetProjectPath returns the full project path
func (b *BaseGenerator) GetProjectPath(cfg *config.Config) string {
	return cfg.GetProjectPath()
}
