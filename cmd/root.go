package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// RootCmd is the base command onto which all other commands are added.
var RootCmd = &cobra.Command{
	Use:   "starship-iac",
	Short: "IaC Automation",
}

// Execute starts RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}