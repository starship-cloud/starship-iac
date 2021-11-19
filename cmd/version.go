
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd prints the current version.
type VersionCmd struct {
	StarshipVersion string
}

// Init returns the runnable cobra command.
func (v *VersionCmd) Init() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current Starship-IaC version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("starship-iac %s\n", v.	StarshipVersion)
		},
	}
}
