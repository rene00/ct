package cli

import (
	"context"
	"ct/config"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"os"
	"time"

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
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
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
		metrics, err := s.Metric.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Config", "Count", "Last"})

		for _, metric := range metrics {
			configDataType, err := s.Config.SelectOne(ctx, metric.MetricID, "data_type")
			if err != nil && err != store.ErrNotFound {
				return fmt.Errorf("Failed finding data_type config: %s", err)
			}
			configValueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
			if err != nil && err != store.ErrNotFound {
				return fmt.Errorf("Failed finding value_text config: %s", err)
			}
			configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
			if err != nil && err != store.ErrNotFound {
				return fmt.Errorf("Failed finding metric_type config: %s", err)
			}

			last30Days, err := s.Log.SelectWithTimestamp(ctx, metric.MetricID, time.Now().AddDate(0, 0, -30))
			if err != nil {
				return fmt.Errorf("Failed finding last 30 days of logs: %s", err)
			}

			lastLog, err := s.Log.SelectLast(ctx, metric.MetricID)
			if err != nil && err != store.ErrNotFound {
				return fmt.Errorf("Failed finding last log entry: %s", err)
			}
			lastLogFriendlyTimestamp := "None"
			if lastLog != nil {
				lastLogFriendlyTimestamp = lastLog.Timestamp.Format("2006-01-02")
			}

			configText := fmt.Sprintf("%s;%s;%s", configDataType, configMetricType, configValueText)
			table.Append([]string{metric.Name, configText, fmt.Sprintf("%d", len(last30Days)), lastLogFriendlyTimestamp})
		}

		table.Render()

		return nil
	},
}

var metricCreateCmd = &cobra.Command{
	Use: "create [command]",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
		_ = viper.BindPFlag("metric-name", cmd.Flags().Lookup("metric-name"))
		_ = viper.BindPFlag("data-type", cmd.Flags().Lookup("data-type"))
		_ = viper.BindPFlag("value-text", cmd.Flags().Lookup("value-text"))
		_ = viper.BindPFlag("metric-type", cmd.Flags().Lookup("metric-type"))
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

		metric, err := s.Metric.Create(ctx, viper.GetString("metric-name"))
		if err != nil {
			return err
		}

		config := &store.Config{metric.MetricID, "data_type", viper.GetString("data-type")}
		if ok := config.IsDataTypeSupported(); !ok {
			return fmt.Errorf("Data type not supported: %s", config.Val)
		}
		if err := s.Config.Upsert(ctx, config); err != nil {
			return err
		}

		config = &store.Config{metric.MetricID, "metric_type", viper.GetString("metric-type")}
		if ok := config.IsMetricTypeSupported(); !ok {
			return fmt.Errorf("Metric type not supported: %s", config.Val)
		}
		if err := s.Config.Upsert(ctx, config); err != nil {
			return err
		}

		valueText := viper.GetString("value-text")
		if valueText != "" {
			if err := s.Config.Upsert(ctx, &store.Config{metric.MetricID, "value_text", valueText}); err != nil {
				return err
			}
		}

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

func initMetricCreateCmd() {
	c := metricCreateCmd
	f := c.Flags()
	f.String("metric-name", "", "name of metric to create")
	c.MarkFlagRequired("metric-name")
	f.String("config-file", "", "")
	f.String("data-type", "float", "metric data type (bool, float or int)")
	f.String("value-text", "", "metric value text")
	f.String("metric-type", "gauge", "metric type (counter or gauge)")
}

func initMetricListCmd() {
	c := metricListCmd
	f := c.Flags()
	f.String("config-file", "", "")
}

func init() {
	initMetricDeleteCmd()
	initMetricCreateCmd()
	initMetricListCmd()
	metricCmd.AddCommand(metricDeleteCmd)
	metricCmd.AddCommand(metricCreateCmd)
	metricCmd.AddCommand(metricListCmd)
	rootCmd.AddCommand(metricCmd)
}
