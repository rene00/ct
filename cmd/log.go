package cmd

import (
	"context"
	"ct/config"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logCmd = &cobra.Command{
	Use:   "log [command]",
	Short: "Log a metric",
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, flag := range []string{"config-file", "metric", "value", "timestamp", "quiet", "edit"} {
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

		metricName := viper.GetString("metric")
		quiet := viper.GetBool("quiet")
		value := viper.GetString("value")

		timestamp, err := parseTimestamp(viper.GetString("timestamp"))
		if err != nil {
			return err
		}

		ctx := context.Background()
		metric, err := s.Metric.SelectOne(ctx, metricName)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		if err != nil && err == store.ErrNotFound {
			metric, err = s.Metric.Create(ctx, metricName)
			if err != nil {
				return err
			}
			if err = s.Config.Create(ctx, metric.MetricID); err != nil {
				return err
			}
		}

		log, err := s.Log.SelectOne(ctx, metric.MetricID, timestamp)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		if log != nil && !viper.GetBool("edit") {
			if quiet {
				return nil
			}
			return fmt.Errorf("log for %s with timestamp of %s already exists", metric.Name, timestamp.Format("2006-01-02"))
		}

		valueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
		if err != nil && err != store.ErrNotFound {
			return err
		}

		value, err = getValueFromConsole(value, valueText)
		if err != nil {
			return err
		}

		logFunc := s.Log.Create
		if viper.GetBool("edit") {
			logFunc = s.Log.Upsert
		}

		err = logFunc(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp})
		if err != nil && !quiet {
			return err
		}

		return nil
	},
}

var logCreateCmd = &cobra.Command{
	Use:   "create [command]",
	Short: "Create a new log entry",
	Long: `
Create a new log entry

EXAMPLES

- Create a new log entry for the weight metric and prompt for the value:

  $ ct log create --metric-name weight

- Same as above but specify the value on the command line:

  $ ct log create --metric-name weight --metric-value 100

- Same as above but with the timestamp of 2020-01-01:

  $ ct log create --metric-name weight --metric-value 95 --timestamp 2020-01-01

- Same as above but update an existing log entry for the same timestamp:

  $ ct log create --metric-name weight --metric-value 96 --timestamp 2020-01-01 --update
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
		_ = viper.BindPFlag("metric-name", cmd.Flags().Lookup("metric-name"))
		_ = viper.BindPFlag("metric-value", cmd.Flags().Lookup("metric-value"))
		_ = viper.BindPFlag("timestamp", cmd.Flags().Lookup("timestamp"))
		_ = viper.BindPFlag("quiet", cmd.Flags().Lookup("quiet"))
		_ = viper.BindPFlag("update", cmd.Flags().Lookup("update"))
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

		metricName := viper.GetString("metric-name")

		metric, err := s.Metric.SelectOne(ctx, metricName)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		if err != nil && err == store.ErrNotFound {
			return fmt.Errorf("Metric not found: %s", metricName)
		}

		timestamp, err := parseTimestamp(viper.GetString("timestamp"))
		if err != nil {
			return err
		}

		log, err := s.Log.SelectOne(ctx, metric.MetricID, timestamp)
		if err != nil && err != store.ErrNotFound {
			return err
		}

		quiet := viper.GetBool("quiet")
		if log != nil && !viper.GetBool("update") {
			if quiet {
				return nil
			}
			return fmt.Errorf("log for %s with timestamp of %s already exists", metric.Name, timestamp.Format("2006-01-02"))
		}

		valueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
		if err != nil && err != store.ErrNotFound {
			return err
		}

		value, err := getValueFromConsole(viper.GetString("value"), valueText)
		if err != nil {
			return err
		}

		logFunc := s.Log.Create
		if viper.GetBool("update") {
			logFunc = s.Log.Upsert
		}

		err = logFunc(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp})
		if err != nil && !quiet {
			return err
		}

		return nil
	},
}

func initLogCmd() {
	c := logCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("value", "", "Value")
	f.String("config-file", "", "")
	f.Bool("quiet", false, "")
	f.String("timestamp", time.Now().Format("2006-01-02"), "")
	f.Bool("edit", false, "")
}

func initLogCreateCmd() {
	c := logCreateCmd
	f := c.Flags()
	f.String("metric-name", "", "The name of the metric")
	c.MarkFlagRequired("metric-name")
	f.String("metric-value", "", "The value of the metric")
	f.String("config-file", "", "Filepath of the configuration file")
	f.Bool("quiet", false, "Dont print warnings")
	f.String("timestamp", time.Now().Format("2006-01-02"), "The timestamp of the metric (format: YYYY-MM-DD)")
	f.Bool("update", false, "Update an existing metric value logged for the same timestamp")
}

func init() {
	initLogCmd()
	initLogCreateCmd()
	logCmd.AddCommand(logCreateCmd)
	rootCmd.AddCommand(logCmd)
}
