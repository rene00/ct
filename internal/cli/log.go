package cli

import (
	"context"
	"ct/internal/report"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func logCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "log",
	}

	cmd.AddCommand(promptLogCmd(cli))
	cmd.AddCommand(createLogCmd(cli))

	return cmd
}

func promptLogCmd(cli *cli) *cobra.Command {
	var flags struct {
		Timestamp string
	}

	var cmd = &cobra.Command{
		Use:   "prompt",
		Short: "Prompt a value for all metrics that havent been logged for a timestamp",
		Long: `
Prompt a value for all metrics that havent been logged for a timestamp

EXAMPLES

- Prompt for all metrics with the current timestamp

  $ ct log prompt

- Prompt for all metrics with the timestamp of 2020-01-01

  $ ct log prompt --timestamp 2020-01-01
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("timestamp", cmd.Flags().Lookup("timestamp"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("sqlite3", cli.config.DBFile)
			if err != nil {
				return err
			}
			defer db.Close()

			timestamp, err := parseTimestamp(flags.Timestamp)
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
						_, err = s.Log.Create(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp})
						if err != nil {
							return fmt.Errorf("Failed to create log: %s", err)
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.Timestamp, "timestamp", time.Now().Format("2006-01-02"), "The timestamp of the metric (format: YYYY-MM-DD)")
	return cmd
}

func createLogCmd(cli *cli) *cobra.Command {
	var flags struct {
		Timestamp string
		Comment   string
		Update    bool
		Quiet     bool
		Feedback  bool
	}
	var cmd = &cobra.Command{
		Use:   "create [command]",
		Short: "Create a new log entry",
		Long: `
Create a new log entry

EXAMPLES

- Create a new log entry for the weight metric and prompt for the value:

  $ ct log create weight

- Same as above but specify the value on the command line:

  $ ct log create weight 100

- Same as above but with the timestamp of 2020-01-01:

  $ ct log create weight 95 --timestamp 2020-01-01

- Same as above but update an existing log entry for the same timestamp:

  $ ct log create weight 96 --timestamp 2020-01-01 --update

- Create a new log entry with a comment:

  $ ct log create walk 3.0 --comment "Walked to the beach and back"
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("timestamp", cmd.Flags().Lookup("timestamp"))
			_ = viper.BindPFlag("quiet", cmd.Flags().Lookup("quiet"))
			_ = viper.BindPFlag("update", cmd.Flags().Lookup("update"))
			_ = viper.BindPFlag("comment", cmd.Flags().Lookup("comment"))
			_ = viper.BindPFlag("feedback", cmd.Flags().Lookup("feedback"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("Missing metric name and/or value")
			}
			m := args[0]
			v := args[1]

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

			timestamp, err := parseTimestamp(flags.Timestamp)
			if err != nil {
				return err
			}

			log, err := s.Log.SelectOne(ctx, metric.MetricID, timestamp)
			if err != nil && err != store.ErrNotFound {
				return err
			}

			if log != nil && !flags.Update {
				if flags.Quiet {
					return nil
				}
				return fmt.Errorf("log for %s with timestamp of %s already exists", metric.Name, timestamp.Format("2006-01-02"))
			}

			valueText, err := s.Config.SelectOne(ctx, metric.MetricID, "value_text")
			if err != nil && err != store.ErrNotFound {
				return err
			}

			value, err := getValueFromConsole(v, valueText)
			if err != nil {
				return err
			}

			logFunc := s.Log.Create
			if flags.Update {
				logFunc = s.Log.Upsert
			}

			log, err = logFunc(ctx, &store.Log{MetricID: metric.MetricID, Value: value, Timestamp: timestamp})
			if err != nil && !flags.Quiet {
				return fmt.Errorf("Failed to create log: %s", err)
			}

			if flags.Comment != "" {
				if err = s.LogComment.Upsert(ctx, log, flags.Comment); err != nil {
					return fmt.Errorf("Failed to insert/update log comment: %s", err)
				}
			}

			if flags.Feedback {
				configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
				if err != nil && err != store.ErrNotFound {
					return err
				}
				if err != nil && err == store.ErrNotFound {
					return fmt.Errorf("Missing config option metric_type: %s", metric.Name)
				}

				switch configMetricType {
				case "gauge":
					r := report.NewReport(db, metric)
					output, err := r.MonthlyGuage(ctx, report.WithStartTimestamp(time.Now().AddDate(0, -1, 0)))
					if err != nil {
						return err
					}
					cmd.Print(output)
				}
			}

			return nil
		},
	}
	cmd.Flags().BoolVar(&flags.Update, "update", false, "Update an existing log entry")
	cmd.Flags().BoolVar(&flags.Quiet, "quiet", false, "Dont print warnings")
	cmd.Flags().BoolVar(&flags.Feedback, "feedback", false, "Provide feedback when log created")
	cmd.Flags().StringVar(&flags.Timestamp, "timestamp", time.Now().Format("2006-01-02"), "The timestamp of the metric (format: YYYY-MM-DD)")
	cmd.Flags().StringVar(&flags.Comment, "comment", "", "A log comment")
	return cmd
}
