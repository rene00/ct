package cli

import (
	"ct/db/migrations"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func dbCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "db",
	}

	cmd.AddCommand(migrateDbCmd(cli))

	return cmd
}

func migrateDbCmd(cli *cli) *cobra.Command {
	var flags struct {
		Run bool
	}

	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "run DB migrations",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("run", cmd.Flags().Lookup("run"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.Run {
				if err := migrations.DoMigrateDb(fmt.Sprintf("sqlite3://%s", cli.config.DBFile)); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&flags.Run, "run", false, "Run DB migrations")
	return cmd
}
