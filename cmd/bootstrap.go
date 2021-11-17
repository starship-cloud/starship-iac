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
	"os"

	"github.com/runstarship/starship/testdrive"
	"github.com/spf13/cobra"
)

// TestdriveCmd starts the testdrive process for testing out starship.
type TestdriveCmd struct{}

// Init returns the runnable cobra command.
func (b *TestdriveCmd) Init() *cobra.Command {
	return &cobra.Command{
		Use:   "testdrive",
		Short: "Start a guided tour of starship",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := testdrive.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError: %s\033[39m\n\n", err.Error())
			}
			return err
		},
		SilenceErrors: true,
	}
}
