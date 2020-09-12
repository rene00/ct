package cmd

import (
	"context"
	"ct/config"
	"ct/internal/store"
	"database/sql"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			return err
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			return err
		}
		defer db.Close()

		s := store.NewStore(db)

		ctx := context.Background()
		metrics, err := s.Metric.SelectAll(ctx)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name"})

		for _, metric := range metrics {
			table.Append([]string{m.Name})
		}

		table.Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
