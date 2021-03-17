package cli

import (
	"ct/config"
	"ct/db/migrations"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbCmd = &cobra.Command{
	Use: "db [command]",
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate [command]",
	Short: "run DB migrations",
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, flag := range []string{"run", "config-file"} {
			_ = viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		if viper.IsSet("run") {
			usrCfg := cfg.UserViperConfig
			dbURL := fmt.Sprintf("sqlite3://%s", usrCfg.GetString("ct.db_file"))
			if err = migrations.DoMigrateDb(dbURL); err != nil {
				fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
				os.Exit(1)
			}
		}
		return nil
	},
}

func initDbMigrateCmd() {
	c := dbMigrateCmd
	f := c.Flags()
	f.Bool("run", false, "Run DB migrations")
	f.String("config-file", "", "")
}

func init() {
	initDbMigrateCmd()
	dbCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(dbCmd)
}
