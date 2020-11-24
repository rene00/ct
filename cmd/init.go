package cmd

import (
	"ct/config"
	"ct/db/migrations"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use: "init [command]",
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, flag := range []string{"db-file", "config-file"} {
			_ = viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}

		dbFile := viper.GetString("db-file")
		if dbFile == "" {
			dbFile = filepath.Join(cfg.Dir, "ct.db")
		}

		usrCfg := cfg.UserViperConfig
		usrCfg.Set("ct", struct {
			DbFile string `json:"db_file"`
		}{DbFile: dbFile})

		configFilePath := usrCfg.ConfigFileUsed()
		if _, err := os.Stat(configFilePath); err == nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("Not clobbering config file %s.\n", configFilePath))
			return err
		}
		if err := cfg.Save("ct"); err != nil {
			return err
		}

		dbURL := fmt.Sprintf("sqlite3://%s", dbFile)
		if err := migrations.DoMigrateDb(dbURL); err != nil {
			return err
		}

		return nil
	},
}

func initInitCmd() {
	c := initCmd
	f := c.Flags()
	f.String("db-file", "", "DB file")
	f.String("config-file", "", "Config file")
}

func init() {
	initInitCmd()
	rootCmd.AddCommand(initCmd)
}
