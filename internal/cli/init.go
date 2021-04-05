package cli

import (
	"ct/db/migrations"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCmd(cli *cli) *cobra.Command {
	var flags struct {
		DBFile string
	}

	var cmd = &cobra.Command{
		Use: "init",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("db-file", cmd.Flags().Lookup("db-file"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(cli.configFile); err == nil {
				return fmt.Errorf("not clobbering config file %s", cli.configFile)
			}

			cli.config.DBFile = flags.DBFile
			if err := cli.persistConfig(); err != nil {
				return err
			}

			if err := migrations.DoMigrateDb(fmt.Sprintf("sqlite3://%s", flags.DBFile)); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&flags.DBFile, "db-file", path.Join(os.Getenv("HOME"), ".config", "ct", "ct.db"), "Database file path")

	return cmd
}
