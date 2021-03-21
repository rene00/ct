package cli

import (
	"context"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func metricCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "metric",
	}

	cmd.AddCommand(configMetricCmd(cli))
	cmd.AddCommand(deleteMetricCmd(cli))
	cmd.AddCommand(listMetricCmd(cli))
	cmd.AddCommand(createMetricCmd(cli))

	return cmd
}

func configMetricCmd(cli *cli) *cobra.Command {
	var flags struct {
		DataType   string
		ValueText  string
		MetricType string
	}

	var cmd = &cobra.Command{
		Use: "configure",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("data-type", cmd.Flags().Lookup("data-type"))
			_ = viper.BindPFlag("value-text", cmd.Flags().Lookup("value-text"))
			_ = viper.BindPFlag("metric-type", cmd.Flags().Lookup("metric-type"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Missing metric")
			}
			m := args[0]

			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			s := store.NewStore(db)

			ctx := context.Background()
			metric, err := s.Metric.SelectOne(ctx, m)
			if err != nil && err != store.ErrNotFound {
				return err
			}
			if err != nil && err == store.ErrNotFound {
				return fmt.Errorf("Metric not found: %s", m)
			}

			if flags.ValueText != "" {
				config := store.NewConfig(metric.MetricID, "value_text", flags.ValueText)
				if err := s.Config.Upsert(ctx, config); err != nil {
					return err
				}
			}

			if flags.DataType != "" {
				config := store.NewConfig(metric.MetricID, "data_type", flags.DataType)
				if ok := config.IsDataTypeSupported(); !ok {
					return fmt.Errorf("Data type not supported")
				}
				if err := s.Config.Upsert(ctx, config); err != nil {
					return err
				}
			}

			if flags.MetricType != "" {
				config := store.NewConfig(metric.MetricID, "metric_type", flags.MetricType)
				if ok := config.IsMetricTypeSupported(); !ok {
					return fmt.Errorf("Metric type not supported")
				}
				if err := s.Config.Upsert(ctx, config); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.DataType, "data-type", "", "Metric data type")
	cmd.Flags().StringVar(&flags.ValueText, "value-text", "", "Metric value text")
	cmd.Flags().StringVar(&flags.MetricType, "metric-type", "", "Metric type")

	return cmd
}

func deleteMetricCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "delete",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("metric-name", cmd.Flags().Lookup("metric-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Missing metric")
			}
			m := args[0]

			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			s := store.NewStore(db)

			ctx := context.Background()
			metric, err := s.Metric.SelectOne(ctx, m)
			if err != nil {
				return err
			}

			return s.Metric.Delete(ctx, metric.MetricID)
		},
	}
	return cmd
}

func listMetricCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("sqlite3", cli.config.DBFile)
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
	return cmd
}

func createMetricCmd(cli *cli) *cobra.Command {
	var flags struct {
		DataType   string
		ValueText  string
		MetricType string
	}

	var cmd = &cobra.Command{
		Use: "create",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("data-type", cmd.Flags().Lookup("data-type"))
			_ = viper.BindPFlag("value-text", cmd.Flags().Lookup("value-text"))
			_ = viper.BindPFlag("metric-type", cmd.Flags().Lookup("metric-type"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Missing metric")
			}
			m := args[0]

			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			s := store.NewStore(db)
			ctx := context.Background()

			metric, err := s.Metric.Create(ctx, m)
			if err != nil {
				return err
			}

			config := store.NewConfig(metric.MetricID, "data_type", flags.DataType)
			if ok := config.IsDataTypeSupported(); !ok {
				return fmt.Errorf("Data type not supported: %s", config.Val)
			}
			if err := s.Config.Upsert(ctx, config); err != nil {
				return err
			}

			config = store.NewConfig(metric.MetricID, "metric_type", flags.MetricType)
			if ok := config.IsMetricTypeSupported(); !ok {
				return fmt.Errorf("Metric type not supported: %s", config.Val)
			}
			if err := s.Config.Upsert(ctx, config); err != nil {
				return err
			}

			if flags.ValueText != "" {
				config = store.NewConfig(metric.MetricID, "value_text", flags.ValueText)
				if err := s.Config.Upsert(ctx, config); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.DataType, "data-type", "float", "Metric data type")
	cmd.Flags().StringVar(&flags.MetricType, "metric-type", "gauge", "Metric type")
	cmd.Flags().StringVar(&flags.ValueText, "value-text", "", "Metric value text")

	return cmd
}
