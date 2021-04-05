package cli

import (
	"context"
	"ct/internal/store"
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
)

type dumpOutput struct {
	Metrics []store.Metric `json:"metrics"`
	Configs []store.Config `json:"configs"`
	Logs    []store.Log    `json:"logs"`
}

func dumpCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			s := store.NewStore(db)
			ctx := context.Background()

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

			cmd.Print(string(dump))

			return nil
		},
	}

	return cmd
}
