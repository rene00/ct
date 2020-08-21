package cmd

import (
	"ct/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ct [command]",
}

// Execute the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error")
		os.Exit(1)
	}
}

func init() {
	BinaryName := os.Args[0]
	config.SetDefaultDirName(BinaryName)
}
