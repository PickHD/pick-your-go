// Package generator provides architecture-specific generators
package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PickHD/pick-your-go/internal/config"
	"github.com/PickHD/pick-your-go/internal/template"
)

// HexagonalGenerator generates projects with hexagonal architecture
type HexagonalGenerator struct {
	*BaseGenerator
	templateManager *template.Manager
}

// NewHexagonalGenerator creates a new hexagonal architecture generator
func NewHexagonalGenerator() *HexagonalGenerator {
	return &HexagonalGenerator{
		BaseGenerator:   NewBaseGenerator(),
		templateManager: template.NewManager(),
	}
}

// Generate creates a hexagonal architecture project
func (g *HexagonalGenerator) Generate(cfg *config.Config) error {
	if err := g.ValidateConfig(cfg); err != nil {
		return err
	}

	projectPath := g.GetProjectPath(cfg)

	// Check if directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory already exists: %s", projectPath)
	}

	// Get GitHub token from environment
	token := os.Getenv("PICK_YOUR_GO_GITHUB_TOKEN")

	// Ensure template is cached
	fmt.Println("Ensuring template is cached...")
	if err := g.templateManager.EnsureTemplateCached(config.HexagonalArchitecture, token); err != nil {
		return fmt.Errorf("failed to ensure template is cached: %w", err)
	}

	// Copy template to destination
	fmt.Println("Copying template to destination...")
	if err := g.templateManager.CopyTemplateToDestination(config.HexagonalArchitecture, projectPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Customize project-specific files
	fmt.Println("Customizing project files...")
	if err := g.customizeProject(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to customize project: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid for hexagonal architecture
func (g *HexagonalGenerator) Validate(cfg *config.Config) error {
	return g.ValidateConfig(cfg)
}

// GetStructure returns the directory structure for hexagonal architecture
func (g *HexagonalGenerator) GetStructure() []string {
	return []string{
		"cmd/",
		"internal/domain/",
		"internal/ports/",
		"internal/ports/in/",
		"internal/ports/out/",
		"internal/adapters/",
		"internal/adapters/in/",
		"internal/adapters/out/",
		"internal/app/",
		"pkg/",
		"configs/",
		"docs/",
	}
}

// customizeProject customizes the project with user-specific information
func (g *HexagonalGenerator) customizeProject(cfg *config.Config, projectPath string) error {

	// Verify projectPath is absolute
	if !filepath.IsAbs(projectPath) {
		return fmt.Errorf("BUG: projectPath is not absolute: %s", projectPath)
	}

	// Update go.mod with correct module path
	goModPath := filepath.Join(projectPath, "go.mod")

	// CRITICAL: Extract original module path BEFORE updating go.mod
	oldModule, err := extractOriginalModulePath(goModPath)
	if err != nil {
		return fmt.Errorf("failed to extract original module path: %w", err)
	}

	if err := updateGoModule(goModPath, cfg.ModulePath); err != nil {
		fmt.Printf("Warning: failed to update go.mod: %v\n", err)
		// Don't return error here, just warn
	}

	// CRITICAL: Update all import paths in .go files
	// This is necessary because the template uses its own module name in imports
	if oldModule != cfg.ModulePath {
		fmt.Println("Updating import paths in Go files...")
		if err := updateImportPaths(projectPath, oldModule, cfg.ModulePath); err != nil {
			return fmt.Errorf("failed to update import paths: %w", err)
		}
		fmt.Printf("Successfully updated import paths from '%s' to '%s'\n", oldModule, cfg.ModulePath)
	}

	return nil
}
