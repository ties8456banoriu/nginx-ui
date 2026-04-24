package main

import (
	"fmt"
	"os"

	"github.com/0xJacky/nginx-ui/cmd"
)

// Version information set by build flags
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Inject version info into the cmd package
	cmd.Version = Version
	cmd.Commit = Commit
	cmd.BuildDate = BuildDate

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		// Use exit code 2 to distinguish application errors from OS-level errors (code 1)
		os.Exit(2)
	}
}
