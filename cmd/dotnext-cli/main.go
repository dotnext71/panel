// SPDX-License-Identifier: MIT

// Command dotnext-cli provides the administrative CLI for dotNext.
package main

import (
	"fmt"
	"os"

	"github.com/dotnext-panel/dotnext/internal/panelcmd"
)

func main() {
	if err := panelcmd.ExecuteCLI(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
