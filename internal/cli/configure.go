package cli

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //nolint

	"ct/internal/store"
	"database/sql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configureCmd(cli *cli) *cobra.Command {
	var flags struct {
		Metric     string
		DataType   string
		ValueText  string
		MetricType string
	}
	var cmd = &cobra.Command{
		Use:  "configure",
		Args: cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("data-type", cmd.Flags().Lookup("data-type"))
			_ = viper.BindPFlag("value-text", cmd.Flags().Lookup("value-text"))
			_ = viper.BindPFlag("metric-type", cmd.Flags().Lookup("metric-type"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			s := store.NewStore(db)

			ctx := context.Background()
			m := args[0]
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
