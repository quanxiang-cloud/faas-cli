/*
Copyright 2020 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package command

import (
	"io"
	"os"

	"github.com/quanxiang-cloud/faas-cli/internal/command/subcommand"
	"github.com/spf13/cobra"
)

// TODO pre and post scaffold

// NewDefaultCommand creates the default command with default arguments
func NewDefaultCommand() *cobra.Command {
	return NewDefaultCommandWithArgs(os.Args, os.Stdin, os.Stdout, os.Stderr)
}

// NewDefaultCommandWithArgs creates the default command with arguments
func NewDefaultCommandWithArgs(args []string, in io.Reader, out, errout io.Writer) *cobra.Command {
	return NewCommand(in, out, errout)
}

// NewCommand creates the command
func NewCommand(in io.Reader, out, errout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qlf",
		Short: "QUANXIANG lowcode faas command line",
		Long: `
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(subcommand.NewRun())
	return cmd
}
