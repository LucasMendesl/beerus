package main

import (
	"log/slog"
	"os"

	"github.com/lucasmendesl/beerus/cmd"
)

// main is the entry point for the beerus application.
//
// It will run the root command and exit with a non-zero status code if
// any error occurs.
func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error on executing application", "err", err)
		os.Exit(1)
	}
}
