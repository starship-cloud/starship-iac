// Copyright 2021 starship studio.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Modified hereafter by contributors to runstarship/starship.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd prints the current version.
type VersionCmd struct {
	starshipVersion string
}

// Init returns the runnable cobra command.
func (v *VersionCmd) Init() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current starship version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("starship %s\n", v.starshipVersion)
		},
	}
}
