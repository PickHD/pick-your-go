// Package cli provides the CLI command structure using cobra
package cli

import (
	"github.com/PickHD/pick-your-go/internal/cli/cmd"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pick-your-go",
	Short: "A CLI tool to generate Go projects with various architecture patterns",
	Long: `pick-your-go is a CLI tool that helps you quickly scaffold Go projects
with different architecture patterns like Layered, Modular, and Hexagonal.

It uses interactive prompts to gather project information and generates
a complete, production-ready project structure based on your chosen architecture.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(cmd.NewInitCommand())
	rootCmd.AddCommand(cmd.NewTemplatesCommand())
}

// GetRootCommand returns the root command for testing purposes
func GetRootCommand() *cobra.Command {
	return rootCmd
}
