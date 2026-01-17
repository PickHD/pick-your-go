// Package generator provides architecture-specific generators
package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExtractOriginalModulePath tests extracting module path from go.mod
func TestExtractOriginalModulePath(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "go-mod-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test go.mod file
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := `module github.com/test/example

go 1.21

require (
	github.com/stretchr/testify v1.8.0
)
`
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	// Test extraction
	modulePath, err := extractOriginalModulePath(goModPath)
	if err != nil {
		t.Fatalf("extractOriginalModulePath failed: %v", err)
	}

	expectedModule := "github.com/test/example"
	if modulePath != expectedModule {
		t.Errorf("expected module path '%s', got '%s'", expectedModule, modulePath)
	}
}

// TestReplaceModulePathInLine tests replacing module path in a single line
func TestReplaceModulePathInLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		oldModule  string
		newModule  string
		expected   string
		shouldChange bool
	}{
		{
			name:       "Simple import",
			line:       `import "github.com/old/module/internal/config"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `import "github.com/new/module/internal/config"`,
			shouldChange: true,
		},
		{
			name:       "Aliased import",
			line:       `oldconfig "github.com/old/module/internal/config"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `oldconfig "github.com/new/module/internal/config"`,
			shouldChange: true,
		},
		{
			name:       "Import block entry",
			line:       `\t"github.com/old/module/internal/domain"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `\t"github.com/new/module/internal/domain"`,
			shouldChange: true,
		},
		{
			name:       "External package (should not change)",
			line:       `import "github.com/gin-gonic/gin"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `import "github.com/gin-gonic/gin"`,
			shouldChange: false,
		},
		{
			name:       "Standard library (should not change)",
			line:       `import "fmt"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `import "fmt"`,
			shouldChange: false,
		},
		{
			name:       "Root module path",
			line:       `import "github.com/old/module"`,
			oldModule:  "github.com/old/module",
			newModule:  "github.com/new/module",
			expected:   `import "github.com/new/module"`,
			shouldChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceModulePathInLine(tt.line, tt.oldModule, tt.newModule)
			if result != tt.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, result)
			}
			if tt.shouldChange && result == tt.line {
				t.Errorf("expected line to change, but it remained the same")
			}
			if !tt.shouldChange && result != tt.line {
				t.Errorf("expected line to remain the same, but it changed")
			}
		})
	}
}

// TestReplaceImportPaths tests replacing import paths in a full file content
func TestReplaceImportPaths(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		oldModule string
		newModule string
		// Strings that should be present in the result
		mustContain []string
		// Strings that should NOT be present in the result
		mustNotContain []string
	}{
		{
			name: "Single import statement",
			content: `package main

import "github.com/old/module/internal/config"

func main() {
	config.Load()
}`,
			oldModule: "github.com/old/module",
			newModule: "github.com/new/module",
			mustContain: []string{`import "github.com/new/module/internal/config"`},
			mustNotContain: []string{`github.com/old/module`},
		},
		{
			name: "Import block",
			content: `package main

import (
	"github.com/old/module/internal/config"
	"github.com/old/module/internal/domain"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
}`,
			oldModule: "github.com/old/module",
			newModule: "github.com/new/module",
			mustContain: []string{
				`"github.com/new/module/internal/config"`,
				`"github.com/new/module/internal/domain"`,
				`"github.com/gin-gonic/gin"`,
			},
			mustNotContain: []string{
				`github.com/old/module`,
			},
		},
		{
			name: "Mixed imports",
			content: `package main

import (
	"fmt"
	oldconfig "github.com/old/module/internal/config"
	"github.com/old/module/internal/domain"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(oldconfig.Load())
}`,
			oldModule: "github.com/old/module",
			newModule: "github.com/user/project",
			mustContain: []string{
				`"fmt"`,
				`oldconfig "github.com/user/project/internal/config"`,
				`"github.com/user/project/internal/domain"`,
				`"github.com/gin-gonic/gin"`,
			},
			mustNotContain: []string{
				`github.com/old/module`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceImportPaths(tt.content, tt.oldModule, tt.newModule)

			// Check mustContain
			for _, s := range tt.mustContain {
				if !strings.Contains(result, s) {
					t.Errorf("expected result to contain: %s\nbut got:\n%s", s, result)
				}
			}

			// Check mustNotContain
			for _, s := range tt.mustNotContain {
				if strings.Contains(result, s) {
					t.Errorf("expected result NOT to contain: %s\nbut got:\n%s", s, result)
				}
			}
		})
	}
}

// TestUpdateImportPathsInFile tests updating import paths in a file
func TestUpdateImportPathsInFile(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "import-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

import "github.com/old/module/internal/config"

func main() {
	config.Load()
}`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Update import paths
	oldModule := "github.com/old/module"
	newModule := "github.com/new/module"
	if err := updateImportPathsInFile(testFile, oldModule, newModule); err != nil {
		t.Fatalf("updateImportPathsInFile failed: %v", err)
	}

	// Read updated content
	updatedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	contentStr := string(updatedContent)

	// Verify changes
	if !strings.Contains(contentStr, newModule) {
		t.Errorf("expected content to contain new module %s, got:\n%s", newModule, contentStr)
	}

	if strings.Contains(contentStr, oldModule) {
		t.Errorf("expected content NOT to contain old module %s, got:\n%s", oldModule, contentStr)
	}

	// Verify the specific import line
	expectedImport := `import "github.com/new/module/internal/config"`
	if !strings.Contains(contentStr, expectedImport) {
		t.Errorf("expected import line '%s', got:\n%s", expectedImport, contentStr)
	}
}
