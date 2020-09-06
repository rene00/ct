package cmd

import (
	"context"
	"ct/config"
	"ct/internal/store"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use: "dump [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return err
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			return err
		}
		defer db.Close()

		s := store.NewStore(db)
		ctx := context.Background()

		type dumpOutput struct {
			Metrics []store.Metric `json:"metrics"`
			Configs []store.Config `json:"configs"`
			Logs    []store.Log    `json:"logs"`
		}

		d := &dumpOutput{}

		d.Metrics, err = s.Metric.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		d.Configs, err = s.Config.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		d.Logs, err = s.Log.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		dump, err := json.Marshal(d)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", dump)

		return nil
	},
}

func initDumpCmd() {
	c := dumpCmd
	f := c.Flags()
	f.String("config-file", "", "")
}

func init() {
	initDumpCmd()
	rootCmd.AddCommand(dumpCmd)
}
