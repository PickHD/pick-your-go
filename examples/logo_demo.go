package main

import (
	"fmt"
	"pick-your-go/pkg/ui"
)

func main() {
	// Demo the logo
	ui.ShowLogo()

	// Show some example messages
	fmt.Println()
	ui.ShowInfo("This is a demo of the new simple text logo")
	ui.ShowWarning("Clean, modern, and minimalist design")
	fmt.Println("✓ Logo uses simple text with lipgloss styling")
	fmt.Println("✓ Centered alignment with proper spacing")
	fmt.Println("✓ Works perfectly on all terminals")
}
