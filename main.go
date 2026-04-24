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
		// Use exit code 1 for application errors (standard Unix convention).
		// Exit code 2 is typically reserved for misuse of shell builtins.
		// Exit code 126/127 are reserved for command not found / not executable.
		os.Exit(1)
	}
}
