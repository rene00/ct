package cmd

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //nolint

	"ct/config"
	"ct/internal/store"
	"database/sql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use: "configure [command]",
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, flag := range []string{"config-file", "metric", "data-type", "value-text"} {
			_ = viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
		}
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
		metric, err := s.Metric.SelectOne(ctx, viper.GetString("metric"))
		if err != nil && err != store.ErrNotFound {
			return err
		}
		if err != nil && err == store.ErrNotFound {
			return fmt.Errorf("Metric not found: %s", viper.GetString("metric"))
		}

		valueText := viper.GetString("value-text")
		if valueText != "" {
			config := &store.Config{metric.MetricID, "value_text", valueText}
			if err := s.Config.Upsert(ctx, config); err != nil {
				return err
			}
		}

		dataType := viper.GetString("data-type")
		if dataType != "" {
			config := &store.Config{metric.MetricID, "data_type", dataType}
			if ok := config.IsDataTypeSupported(); !ok {
				return fmt.Errorf("Data type not supported")
			}
			if err := s.Config.Upsert(ctx, config); err != nil {
				return err
			}
		}

		return nil
	},
}

func initConfigureCmd() {
	c := configureCmd
	f := c.Flags()
	f.String("config-file", "", "Config file")
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("data-type", "", "Metric Data Type")
	f.String("value-text", "", "Metric Value Text")
}

func init() {
	initConfigureCmd()
	rootCmd.AddCommand(configureCmd)
}
