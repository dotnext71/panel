// SPDX-License-Identifier: MIT

// Package panelcmd implements the dotNext HTTP/API daemon (`dotnext-api`)
// and administrative CLI (`dotnext-cli`) commands.
package panelcmd

import (
	"github.com/dotnext-panel/dotnext/pkg/paths"
)

const defaultConfigPath = paths.ConfigFile

// ExecuteCLI runs the dotnext-cli Cobra command tree.
func ExecuteCLI() error {
	return newRootCmd().Execute()
}

// ExecuteAPI runs the dotnext-api server entrypoint, defaulting to `serve`
// when invoked without an explicit subcommand.
func ExecuteAPI(args []string) error {
	cmd := newRootCmd()
	if len(args) == 0 {
		cmd.SetArgs([]string{"serve"})
	} else {
		first := args[0]
		switch first {
		case "serve", "help", "--help", "-h", "--version", "-v":
			cmd.SetArgs(args)
		default:
			if len(first) > 0 && first[0] == '-' {
				cmd.SetArgs(append([]string{"serve"}, args...))
			} else {
				cmd.SetArgs(args)
			}
		}
	}
	return cmd.Execute()
}
