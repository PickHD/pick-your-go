// Package cmd provides the CLI commands implementation
package cmd

import (
	"fmt"
	"os"

	"github.com/PickHD/pick-your-go/internal/cache"
	"github.com/PickHD/pick-your-go/internal/template"

	"github.com/spf13/cobra"
)

// TemplatesCommand represents the templates command
type TemplatesCommand struct {
	cmd *cobra.Command
}

// NewTemplatesCommand creates a new templates command with subcommands
func NewTemplatesCommand() *cobra.Command {
	templatesCmd := &TemplatesCommand{}

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage project templates",
		Long:  `Manage and interact with project templates. List available templates or update local cache.`,
	}

	// Add subcommands
	cmd.AddCommand(templatesCmd.NewListCommand())
	cmd.AddCommand(templatesCmd.NewUpdateCommand())

	templatesCmd.cmd = cmd
	return cmd
}

// ListCommand represents the templates list command
type ListCommand struct {
	cmd *cobra.Command
}

// NewListCommand creates a new list command
func (c *TemplatesCommand) NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available templates",
		Long:  `List all available architecture templates with their descriptions and status.`,
		RunE:  c.runList,
	}
}

// runList executes the list command
func (c *TemplatesCommand) runList(cmd *cobra.Command, args []string) error {
	manager := template.NewManager()

	templates, err := manager.GetTemplates()
	if err != nil {
		return fmt.Errorf("failed to get templates: %w", err)
	}

	fmt.Println("\nAvailable Templates:")
	fmt.Println("===================")

	for _, tmpl := range templates {
		status := "  (cached)"
		if !manager.IsCached(tmpl.Type) {
			status = "  (not cached)"
		}
		fmt.Printf("\n%s - %s%s\n", tmpl.Type.DisplayName(), tmpl.Name, status)
		fmt.Printf("  %s\n", tmpl.Description)
	}

	fmt.Println()

	return nil
}

// UpdateCommand represents the templates update command
type UpdateCommand struct {
	cmd *cobra.Command
}

// NewUpdateCommand creates a new update command
func (c *TemplatesCommand) NewUpdateCommand() *cobra.Command {
	updateCmd := &UpdateCommand{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update template cache from remote repositories",
		Long: `Update the local template cache by pulling the latest changes from
remote GitHub repositories. This requires PICK_YOUR_GO_GITHUB_TOKEN
environment variable to be set for private repositories.`,
		RunE: updateCmd.Run,
	}

	updateCmd.cmd = cmd
	return cmd
}

// Run executes the update command
func (c *UpdateCommand) Run(cmd *cobra.Command, args []string) error {
	// Check for GitHub token
	token := os.Getenv("PICK_YOUR_GO_GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("PICK_YOUR_GO_GITHUB_TOKEN environment variable is required for accessing private repositories")
	}

	fmt.Println("Updating template cache...")

	manager := template.NewManager()
	cacheMgr := cache.NewManager()

	templates, err := manager.GetTemplates()
	if err != nil {
		return fmt.Errorf("failed to get templates: %w", err)
	}

	for _, tmpl := range templates {
		fmt.Printf("\nUpdating %s template...\n", tmpl.Type.DisplayName())

		if err := manager.UpdateTemplate(tmpl.Type, token); err != nil {
			fmt.Printf("  Warning: Failed to update %s: %v\n", tmpl.Type.DisplayName(), err)
			continue
		}

		// Update cache metadata
		if err := cacheMgr.UpdateCacheTime(tmpl.Type); err != nil {
			fmt.Printf("  Warning: Failed to update cache metadata: %v\n", err)
		}

		fmt.Printf("  %s template updated successfully\n", tmpl.Type.DisplayName())
	}

	fmt.Println("\nTemplate cache update completed!")

	return nil
}
