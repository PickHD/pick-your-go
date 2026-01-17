// Package generator provides architecture-specific generators
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pick-your-go/internal/config"
	"pick-your-go/internal/template"
)

// LayeredGenerator generates projects with layered architecture
type LayeredGenerator struct {
	*BaseGenerator
	templateManager *template.Manager
}

// NewLayeredGenerator creates a new layered architecture generator
func NewLayeredGenerator() *LayeredGenerator {
	return &LayeredGenerator{
		BaseGenerator:   NewBaseGenerator(),
		templateManager: template.NewManager(),
	}
}

// Generate creates a layered architecture project
func (g *LayeredGenerator) Generate(cfg *config.Config) error {
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
	if err := g.templateManager.EnsureTemplateCached(config.LayeredArchitecture, token); err != nil {
		return fmt.Errorf("failed to ensure template is cached: %w", err)
	}

	// Copy template to destination
	fmt.Println("Copying template to destination...")
	if err := g.templateManager.CopyTemplateToDestination(config.LayeredArchitecture, projectPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Customize project-specific files
	fmt.Println("Customizing project files...")
	if err := g.customizeProject(cfg, projectPath); err != nil {
		return fmt.Errorf("failed to customize project: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid for layered architecture
func (g *LayeredGenerator) Validate(cfg *config.Config) error {
	return g.ValidateConfig(cfg)
}

// GetStructure returns the directory structure for layered architecture
func (g *LayeredGenerator) GetStructure() []string {
	return []string{
		"cmd/",
		"internal/domain/",
		"internal/presentation/http/",
		"internal/infrastructure/database/",
		"internal/infrastructure/cache/",
		"pkg/",
		"configs/",
		"docs/",
	}
}

// customizeProject customizes the project with user-specific information
func (g *LayeredGenerator) customizeProject(cfg *config.Config, projectPath string) error {

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

// updateGoModule updates the module path in go.mod
func updateGoModule(goModPath, modulePath string) error {

	// Verify goModPath is absolute
	if !filepath.IsAbs(goModPath) {
		return fmt.Errorf("BUG: goModPath is not absolute: %s", goModPath)
	}

	// Verify file exists before trying to read
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod file does not exist at path: %s", goModPath)
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file %s: %w", goModPath, err)
	}

	contentStr := string(content)
	lines := splitLines(contentStr)

	for i, line := range lines {
		if strings.HasPrefix(line, "module ") {
			lines[i] = fmt.Sprintf("module %s", modulePath)
			break
		}
	}

	if err := os.WriteFile(goModPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod file %s: %w", goModPath, err)
	}

	return nil
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

// extractOriginalModulePath extracts the module path from go.mod before updating
func extractOriginalModulePath(goModPath string) (string, error) {
	// Verify goModPath is absolute
	if !filepath.IsAbs(goModPath) {
		return "", fmt.Errorf("BUG: goModPath is not absolute: %s", goModPath)
	}

	// Verify file exists before trying to read
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "", fmt.Errorf("go.mod file does not exist at path: %s", goModPath)
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod file %s: %w", goModPath, err)
	}

	contentStr := string(content)
	lines := splitLines(contentStr)

	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			// Extract module name (remove "module " prefix)
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			return moduleName, nil
		}
	}

	return "", fmt.Errorf("no module declaration found in go.mod")
}

// updateImportPaths updates all import paths in .go files from oldModule to newModule
func updateImportPaths(projectPath, oldModule, newModule string) error {
	// CRITICAL SAFETY CHECK: Ensure oldModule and newModule are different
	if oldModule == newModule {
		return fmt.Errorf("oldModule and newModule are the same: %s", oldModule)
	}

	// CRITICAL SAFETY CHECK: Ensure both are provided
	if oldModule == "" || newModule == "" {
		return fmt.Errorf("oldModule and newModule must not be empty (old: '%s', new: '%s')", oldModule, newModule)
	}


	// Walk through all files in projectPath
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip vendor directory and hidden directories
			baseName := filepath.Base(path)
			if baseName == "vendor" || baseName == ".git" || strings.HasPrefix(baseName, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip go.mod file (already handled)
		if filepath.Base(path) == "go.mod" {
			return nil
		}

		// Update import paths in this file
		if err := updateImportPathsInFile(path, oldModule, newModule); err != nil {
			// Log error but continue processing other files
			fmt.Printf("Warning: failed to update import paths in %s: %v\n", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking project directory: %w", err)
	}

	return nil
}

// updateImportPathsInFile updates import paths in a single file
func updateImportPathsInFile(filePath, oldModule, newModule string) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	originalContent := contentStr

	// Replace import paths
	// Pattern: "oldModule/path" -> "newModule/path"
	// We need to handle various import formats:
	// 1. import "oldModule/path"
	// 2. import oldModule "oldModule/path" (aliased imports)
	// 3. Multi-line import blocks

	// Replace quoted import paths first (most common case)
	contentStr = replaceImportPaths(contentStr, oldModule, newModule)

	// Only write if content changed
	if contentStr != originalContent {
		if err := os.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

// replaceImportPaths replaces module paths in import statements
func replaceImportPaths(content, oldModule, newModule string) string {
	lines := splitLines(content)
	inImportBlock := false

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect import block start
		if trimmedLine == "import (" {
			inImportBlock = true
			continue
		}

		// Detect import block end
		if inImportBlock && trimmedLine == ")" {
			inImportBlock = false
			continue
		}

		// Check if this line contains an import
		if strings.HasPrefix(trimmedLine, "import ") || inImportBlock {
			// Replace old module path with new module path
			// Handle both: import "oldModule/path" and oldModule "oldModule/path"
			newLine := replaceModulePathInLine(line, oldModule, newModule)
			if newLine != line {
				lines[i] = newLine
			}
		}
	}

	return strings.Join(lines, "\n")
}

// replaceModulePathInLine replaces module path in a single line
func replaceModulePathInLine(line, oldModule, newModule string) string {
	// Extract quoted string from line
	// This handles both:
	// - import "oldModule/path"
	// - alias "oldModule/path"
	// - "oldModule/path" (in import block)

	// Find all quoted strings in the line
	start := 0
	for {
		// Find opening quote
		quoteStart := strings.Index(line[start:], `"`)
		if quoteStart == -1 {
			break
		}
		quoteStart += start

		// Find closing quote
		quoteEnd := strings.Index(line[quoteStart+1:], `"`)
		if quoteEnd == -1 {
			break
		}
		quoteEnd += quoteStart + 1

		// Extract the quoted string
		quotedPath := line[quoteStart+1 : quoteEnd]

		// Check if it starts with oldModule
		if quotedPath == oldModule || strings.HasPrefix(quotedPath, oldModule+"/") {
			// Replace with newModule
			newPath := strings.Replace(quotedPath, oldModule, newModule, 1)
			line = line[:quoteStart+1] + newPath + line[quoteEnd:]
			// Adjust quoteEnd since we modified the line
			quoteEnd = quoteStart + 1 + len(newPath)
		}

		start = quoteEnd + 1
	}

	return line
}
