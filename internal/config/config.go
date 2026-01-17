// Package config handles application configuration
package config

import (
	"fmt"
	"path/filepath"
)

// ArchitectureType defines the supported architecture patterns
type ArchitectureType string

const (
	// LayeredArchitecture represents traditional layered architecture
	LayeredArchitecture ArchitectureType = "layered"
	// ModularArchitecture represents modular monolithic architecture
	ModularArchitecture ArchitectureType = "modular"
	// HexagonalArchitecture represents ports and adapters architecture
	HexagonalArchitecture ArchitectureType = "hexagonal"
)

// Config holds the application configuration
type Config struct {
	// ProjectName is the name of the Go project to generate
	ProjectName string
	// ModulePath is the Go module path (e.g., github.com/user/project)
	ModulePath string
	// Architecture is the selected architecture pattern
	Architecture ArchitectureType
	// OutputDir is the directory where the project will be created
	OutputDir string
	// Author is the project author name
	Author string
	// Description is the project description
	Description string
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if c.ModulePath == "" {
		return fmt.Errorf("module path is required")
	}
	if c.Architecture == "" {
		return fmt.Errorf("architecture type is required")
	}
	if c.OutputDir == "" {
		c.OutputDir = "." // Default to current directory
	}
	return nil
}

// GetProjectPath returns the full path where the project will be created
// CRITICAL: Always returns absolute path to avoid path resolution issues
func (c *Config) GetProjectPath() string {
	// Get absolute path for output directory
	outputDirAbs, err := filepath.Abs(c.OutputDir)
	if err != nil {
		// Fallback to relative path if absolute path fails
		// This should not happen in normal scenarios
		outputDirAbs = c.OutputDir
	}

	// Join with project name
	projectPath := filepath.Join(outputDirAbs, c.ProjectName)

	// Ensure the result is absolute (double-check)
	projectPathAbs, err := filepath.Abs(projectPath)
	if err != nil {
		return projectPath
	}

	return projectPathAbs
}

// String returns string representation of ArchitectureType
func (a ArchitectureType) String() string {
	return string(a)
}

// DisplayName returns a human-readable name for the architecture
func (a ArchitectureType) DisplayName() string {
	switch a {
	case LayeredArchitecture:
		return "Layered Architecture"
	case ModularArchitecture:
		return "Modular Architecture"
	case HexagonalArchitecture:
		return "Hexagonal Architecture"
	default:
		return "Unknown Architecture"
	}
}

// Description returns a description of the architecture pattern
func (a ArchitectureType) Description() string {
	switch a {
	case LayeredArchitecture:
		return "Traditional layered architecture with clear separation between presentation, business logic, and data layers"
	case ModularArchitecture:
		return "Modular monolith with domain-driven design, organizing code into feature modules"
	case HexagonalArchitecture:
		return "Hexagonal architecture (ports and adapters) with isolation of core logic from external concerns"
	default:
		return "Unknown architecture pattern"
	}
}

// IsValid checks if the architecture type is valid
func (a ArchitectureType) IsValid() bool {
	switch a {
	case LayeredArchitecture, ModularArchitecture, HexagonalArchitecture:
		return true
	default:
		return false
	}
}
