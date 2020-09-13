package cmd

import (
	"context"
	"ct/config"
	"ct/internal/store"
	"database/sql"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var metricCmd = &cobra.Command{
	Use: "metric [command]",
}

var metricDeleteCmd = &cobra.Command{
	Use: "delete [command]",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metric-name", cmd.Flags().Lookup("metric-name"))
	},
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
		metric, err := s.Metric.SelectOne(ctx, viper.GetString("metric-name"))
		if err != nil {
			return err
		}

		return s.Metric.Delete(ctx, metric.MetricID)
	},
}

var metricListCmd = &cobra.Command{
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
		metrics, err := s.Metric.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Value Text"})

		for _, metric := range metrics {
			configValueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
			if err != nil && err != store.ErrNotFound {
				return err
			}
			table.Append([]string{metric.Name, configValueText})
		}

		table.Render()

		return nil
	},
}

func initMetricDeleteCmd() {
	c := metricDeleteCmd
	f := c.Flags()
	f.String("metric-name", "", "Name of metric to delete")
	c.MarkFlagRequired("metric-name")
	f.String("config-file", "", "")
}

func init() {
	initMetricDeleteCmd()
	metricCmd.AddCommand(metricDeleteCmd)
	metricCmd.AddCommand(metricListCmd)
	rootCmd.AddCommand(metricCmd)
}
