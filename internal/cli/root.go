package cli

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// Execute the root command.
func Execute() {
	cli := &cli{}

	rootCmd := &cobra.Command{
		Use: "ct",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "init" {
				return nil
			}
			return cli.setup(cmd.Context())
		},
	}

	rootCmd.PersistentFlags().BoolVar(&cli.debug, "debug", false, "Enable debug mode.")
	rootCmd.PersistentFlags().StringVar(&cli.configFile, "config-file", path.Join(os.Getenv("HOME"), ".config", "ct", "config.json"), "Config file path")

	rootCmd.AddCommand(metricCmd(cli))
	rootCmd.AddCommand(logCmd(cli))
	rootCmd.AddCommand(initCmd(cli))
	rootCmd.AddCommand(dumpCmd(cli))
	rootCmd.AddCommand(dbCmd(cli))
	rootCmd.AddCommand(reportCmd(cli))
	rootCmd.AddCommand(configureCmd(cli))

	if err := rootCmd.ExecuteContext(context.TODO()); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
