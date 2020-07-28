package cmd

import (
	"ct/config"
	"ct/internal/storage"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var dbCmd = &cobra.Command{
	Use: "db [command]",
}

var migrateDbCmd = &cobra.Command{
	Use:   "migrate [command]",
	Short: "run DB migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runMigrateDbCmd(cfg, cmd.Flags(), args)
	},
}

func runMigrateDbCmd(cfg *config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	dbUrl := fmt.Sprintf("sqlite3://%s", usrCfg.GetString("ct.db_file"))
	return storage.DoMigrateDb(dbUrl)
}

func initMigrateDbCmd() {
	c := migrateDbCmd
	f := c.Flags()
	f.Bool("run", false, "Run DB migrations")
	f.String("config-file", "", "")
}

func init() {
	initMigrateDbCmd()
	dbCmd.AddCommand(migrateDbCmd)
	rootCmd.AddCommand(dbCmd)
}
