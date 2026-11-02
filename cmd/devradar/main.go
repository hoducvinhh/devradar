package main

import (
	"fmt"
	"os"

	"github.com/devradar/devradar/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.RootCmd.Version = fmt.Sprintf("%s (commit: %s, built at: %s)", version, commit, date)
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
