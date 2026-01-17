// Package cmd provides the CLI commands implementation
package cmd

import (
	"fmt"

	"pick-your-go/internal/config"
	"pick-your-go/internal/generator"
	"pick-your-go/pkg/ui"
	"github.com/spf13/cobra"
)

// InitCommand represents the init command
type InitCommand struct {
	cmd      *cobra.Command
	archType string
	yes      bool // Skip confirmation
}

// NewInitCommand creates a new init command
func NewInitCommand() *cobra.Command {
	initCmd := &InitCommand{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project with chosen architecture",
		Long: `Initialize a new Go project with your chosen architecture pattern.
This command will guide you through an interactive process to:
1. Choose your architecture pattern (Layered, Modular, or Hexagonal)
2. Provide project details (name, module path, author, description)
3. Generate a complete project structure based on your selection`,
		RunE: initCmd.Run,
	}

	// Add flags
	cmd.Flags().StringVarP(&initCmd.archType, "architecture", "a", "", "Architecture type: layered, modular, or hexagonal")
	cmd.Flags().StringP("name", "n", "", "Project name")
	cmd.Flags().StringP("module", "m", "", "Go module path (e.g., github.com/user/project)")
	cmd.Flags().StringP("output", "o", ".", "Output directory for the project")
	cmd.Flags().StringP("author", "u", "", "Author name")
	cmd.Flags().StringP("description", "d", "", "Project description")
	cmd.Flags().BoolVarP(&initCmd.yes, "yes", "y", false, "Skip confirmation prompt")

	initCmd.cmd = cmd
	return cmd
}

// Run executes the init command
func (c *InitCommand) Run(cmd *cobra.Command, args []string) error {
	// Get flag values
	name, _ := cmd.Flags().GetString("name")
	module, _ := cmd.Flags().GetString("module")
	output, _ := cmd.Flags().GetString("output")
	author, _ := cmd.Flags().GetString("author")
	description, _ := cmd.Flags().GetString("description")

	// Check if running in interactive mode
	interactiveMode := name == "" || module == "" || c.archType == ""

	var cfg *config.Config
	var err error

	if interactiveMode {
		// Run interactive form
		cfg, err = ui.RunInitForm(c.archType, name, module, output, author, description)
		if err != nil {
			return fmt.Errorf("interactive form failed: %w", err)
		}
	} else {
		// Create config from flags
		cfg = &config.Config{
			ProjectName:  name,
			ModulePath:   module,
			OutputDir:    output,
			Author:       author,
			Description:  description,
			Architecture: config.ArchitectureType(c.archType),
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
	}

	// Show summary
	ui.ShowSummary(cfg)

	// Confirm before proceeding (skip if --yes flag is set)
	var confirmed bool
	if c.yes {
		confirmed = true
		fmt.Println("\nSkipping confirmation (--yes flag set)")
	} else {
		confirmed, err = ui.ConfirmGeneration()
		if err != nil {
			return fmt.Errorf("confirmation failed: %w", err)
		}
	}

	if !confirmed {
		fmt.Println("\nOperation cancelled.")
		return nil
	}

	// Generate the project
	fmt.Printf("\nGenerating %s project...\n\n", cfg.Architecture.DisplayName())

	factory := generator.NewGeneratorFactory()
	gen, err := factory.CreateGenerator(cfg.Architecture)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	if err := gen.Generate(cfg); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Show success message
	ui.ShowSuccess(cfg)

	return nil
}
