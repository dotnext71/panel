// SPDX-License-Identifier: MIT

// Command dotnext-api runs the dotNext Control Panel HTTP(S)/UDS API server,
// embedded React SPA, reconcilers, and background workers.
package main

import (
	"fmt"
	"os"

	"github.com/dotnext-panel/dotnext/internal/panelcmd"
)

func main() {
	if err := panelcmd.ExecuteAPI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
