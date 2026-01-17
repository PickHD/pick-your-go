// Package ui provides interactive UI components using huh and lipgloss
package ui

import (
	"fmt"
	"strings"

	"github.com/PickHD/pick-your-go/internal/config"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the UI
var (
	// LogoStyle is the style for the "PICK YOUR GO" logo
	LogoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			MarginTop(1).
			MarginBottom(0)

	// LogoBorderStyle is the style for the logo border/separator
	LogoBorderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(2).
			MarginTop(0)

	// TitleStyle is the style for titles
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	// SubtitleStyle is the style for subtitles
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	// SuccessStyle is the style for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	// ErrorStyle is the style for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	// InfoStyle is the style for info messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			MarginBottom(1)

	// WarningStyle is the style for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			MarginBottom(1)

	// SummaryLabelStyle is the style for summary labels (right-aligned)
	SummaryLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true).
				Width(20).
				MarginRight(1)

	// SummaryValueStyle is the style for summary values
	SummaryValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	// SummarySeparatorStyle is the style for summary separator
	SummarySeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))
)

// InitFormData holds the form data
type InitFormData struct {
	Architecture string
	ProjectName  string
	ModulePath   string
	Author       string
	Description  string
	OutputDir    string
}

// RunInitForm runs the interactive initialization form
func RunInitForm(archType, name, module, output, author, description string) (*config.Config, error) {
	// Show logo at the beginning
	ShowLogo()

	formData := InitFormData{
		Architecture: archType,
		ProjectName:  name,
		ModulePath:   module,
		OutputDir:    output,
		Author:       author,
		Description:  description,
	}

	// Architecture selection options
	archOptions := []huh.Option[string]{
		huh.NewOption("Layered Architecture - Traditional layered architecture", config.LayeredArchitecture.String()),
		huh.NewOption("Modular Architecture - Modular monolith with DDD", config.ModularArchitecture.String()),
		huh.NewOption("Hexagonal Architecture - Ports and adapters pattern", config.HexagonalArchitecture.String()),
	}

	// Create form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your architecture pattern").
				Description("Select the architectural pattern you want to use for your project").
				Options(archOptions...).
				Value(&formData.Architecture),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("The name of your project (e.g., my-awesome-app)").
				Prompt("> ").
				Validate(validateNotEmpty).
				Value(&formData.ProjectName),

			huh.NewInput().
				Title("Go module path").
				Description("Your Go module path (e.g., github.com/username/project)").
				Prompt("> ").
				Validate(validateModulePath).
				Value(&formData.ModulePath),

			huh.NewInput().
				Title("Output directory").
				Description("Directory where the project will be created (default: current directory)").
				Prompt("> ").
				Value(&formData.OutputDir).
				Placeholder("."),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Author name").
				Description("Your name or organization name").
				Prompt("> ").
				Value(&formData.Author).
				Placeholder("Your Name"),

			huh.NewText().
				Title("Project description").
				Description("A brief description of your project").
				Placeholder("A brief description...").
				Lines(3).
				Value(&formData.Description),
		),
	)

	// Set default values if provided
	if output == "" {
		formData.OutputDir = "."
	}

	// Run the form
	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("form error: %w", err)
	}

	// Convert to config
	cfg := &config.Config{
		Architecture: config.ArchitectureType(formData.Architecture),
		ProjectName:  formData.ProjectName,
		ModulePath:   formData.ModulePath,
		OutputDir:    formData.OutputDir,
		Author:       formData.Author,
		Description:  formData.Description,
	}

	return cfg, nil
}

// ShowLogo displays the "PICK YOUR GO" logo at the top of the form
func ShowLogo() {
	// Calculate terminal width for centering
	width := 80
	logoText := "PICK YOUR GO"
	logoWidth := len(logoText)

	// Create styled logo components with proper centering
	mainLogo := LogoStyle.
			Width(width).
			Align(lipgloss.Center).
			Render(logoText)

	// Create separator line - adjust to match text width
	separatorWidth := logoWidth + 4
	separator := LogoBorderStyle.
			Width(width).
			Align(lipgloss.Center).
			Render(strings.Repeat("=", separatorWidth))

	// Add a subtle tagline
	taglineStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(width).
			Align(lipgloss.Center).
			MarginTop(0).
			MarginBottom(2)

	tagline := taglineStyle.Render("Interactive Go Project Scaffolder")

	// Render everything with proper spacing
	fmt.Println()
	fmt.Println(mainLogo)
	fmt.Println(separator)
	fmt.Println(tagline)
}

// ShowSummary displays a summary of the configuration
func ShowSummary(cfg *config.Config) {
	// Calculate separator width based on terminal width
	separatorWidth := 80
	separator := strings.Repeat("─", separatorWidth)

	fmt.Println()
	fmt.Println(TitleStyle.Render("Project Configuration Summary"))
	fmt.Println(SummarySeparatorStyle.Render(separator))
	fmt.Println()

	// Print configuration with proper alignment
	// Using lipgloss Width to ensure fixed-width label column
	printSummaryRow("Project Name:", cfg.ProjectName)
	printSummaryRow("Module Path:", cfg.ModulePath)
	printSummaryRow("Architecture:", cfg.Architecture.DisplayName())
	printSummaryRow("Output Directory:", cfg.OutputDir)

	if cfg.Author != "" {
		printSummaryRow("Author:", cfg.Author)
	}

	if cfg.Description != "" {
		printSummaryRow("Description:", cfg.Description)
	}

	printSummaryRow("Project Path:", cfg.GetProjectPath())
	fmt.Println()
}

// printSummaryRow prints a single row of the summary with proper alignment
func printSummaryRow(label, value string) {
	// Render label with fixed width (right-aligned by default with Width())
	labelRendered := SummaryLabelStyle.Render(label)

	// Render value
	valueRendered := SummaryValueStyle.Render(value)

	// Print the row
	fmt.Println(labelRendered + valueRendered)
}

// ConfirmGeneration prompts user to confirm project generation
func ConfirmGeneration() (bool, error) {
	var confirm bool

	confirmForm := huh.NewConfirm().
		Title("Generate project?").
		Description("This will create a new project directory with the selected architecture pattern.").
		Value(&confirm)

	if err := confirmForm.Run(); err != nil {
		return false, err
	}

	return confirm, nil
}

// ShowSuccess displays success message
func ShowSuccess(cfg *config.Config) {
	fmt.Println()
	fmt.Println(SuccessStyle.Render("✓ Project generated successfully!"))
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", cfg.GetProjectPath())
	fmt.Println("  2. Review the generated structure")
	fmt.Println("  3. Start building your application!")
	fmt.Println()

	switch cfg.Architecture {
	case config.LayeredArchitecture:
		fmt.Println(InfoStyle.Render("Layered Architecture Notes:"))
		fmt.Println("  - Presentation layer is in /internal/presentation")
		fmt.Println("  - Business logic is in /internal/domain")
		fmt.Println("  - Data access is in /internal/infrastructure")
	case config.ModularArchitecture:
		fmt.Println(InfoStyle.Render("Modular Architecture Notes:"))
		fmt.Println("  - Each module is self-contained in /internal/modules")
		fmt.Println("  - Shared code is in /internal/shared")
		fmt.Println("  - Follow DDD principles for module boundaries")
	case config.HexagonalArchitecture:
		fmt.Println(InfoStyle.Render("Hexagonal Architecture Notes:"))
		fmt.Println("  - Domain logic is in /internal/domain")
		fmt.Println("  - Ports are in /internal/ports")
		fmt.Println("  - Adapters are in /internal/adapters")
	}

	fmt.Println()
}

// ShowError displays an error message
func ShowError(message string) {
	fmt.Println()
	fmt.Println(ErrorStyle.Render("✗ Error: " + message))
	fmt.Println()
}

// ShowWarning displays a warning message
func ShowWarning(message string) {
	fmt.Println()
	fmt.Println(WarningStyle.Render("⚠ Warning: " + message))
	fmt.Println()
}

// ShowInfo displays an info message
func ShowInfo(message string) {
	fmt.Println()
	fmt.Println(InfoStyle.Render("ℹ " + message))
	fmt.Println()
}

// Validation functions

func validateNotEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("this field cannot be empty")
	}
	return nil
}

func validateModulePath(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if !strings.Contains(s, "/") {
		return fmt.Errorf("module path should be a valid path (e.g., github.com/username/project)")
	}
	return nil
}
