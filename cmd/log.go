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
	Short: "Manage metric logs",
}

// ct log prompt --timestamp --force

var logPromptCmd = &cobra.Command{
	Use:   "prompt [command]",
	Short: "Prompt a value for all metrics that havent been logged for a timestamp",
	Long: `
Prompt a value for all metrics that havent been logged for a timestamp

EXAMPLES

- Prompt for all metrics with the current timestamp

  $ ct log prompt
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
		_ = viper.BindPFlag("timestamp", cmd.Flags().Lookup("timestamp"))
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

		timestamp, err := parseTimestamp(viper.GetString("timestamp"))
		if err != nil {
			return err
		}

		s := store.NewStore(db)
		ctx := context.Background()

		metrics, err := s.Metric.SelectLimit(ctx, 0)
		if err != nil {
			return err
		}

		for _, metric := range metrics {
			valueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
			if err != nil && err != store.ErrNotFound {
				return err
			}
			if err != nil && err == store.ErrNotFound {
				continue
			}

			_, err = s.Log.SelectOne(ctx, metric.MetricID, timestamp)
			if err != nil && err != store.ErrNotFound {
				return err
			}
			if err != nil && err == store.ErrNotFound {
				// TODO: update getValueFromConsole to log an error type when no value is provided. This will allow ct to skip over the prompt when no value is submitted here.
				value, _ := getValueFromConsole("", valueText)
				if value != "" {
					if err = s.Log.Create(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp}); err != nil {
						return fmt.Errorf("Failed to create log: %s", err)
					}
				}
			}
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

		value, err := getValueFromConsole(viper.GetString("metric-value"), valueText)
		if err != nil {
			return err
		}

		logFunc := s.Log.Create
		if viper.GetBool("update") {
			logFunc = s.Log.Upsert
		}

		err = logFunc(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp})
		if err != nil && !quiet {
			return fmt.Errorf("Failed to create log: %s", err)
		}

		return nil
	},
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

func initLogPromptCmd() {
	c := logPromptCmd
	f := c.Flags()
	f.String("config-file", "", "Filepath of the configuration file")
	f.String("timestamp", time.Now().Format("2006-01-02"), "The timestamp of the metric (format: YYYY-MM-DD)")
}

func init() {
	initLogCreateCmd()
	initLogPromptCmd()
	logCmd.AddCommand(logCreateCmd)
	logCmd.AddCommand(logPromptCmd)
	rootCmd.AddCommand(logCmd)
}
