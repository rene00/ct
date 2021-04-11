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
)

func reportCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "report",
	}
	cmd.AddCommand(streakReportCmd(cli))
	cmd.AddCommand(dailyReportCmd(cli))
	cmd.AddCommand(weeklyReportCmd(cli))
	cmd.AddCommand(monthlyReportCmd(cli))
	return cmd
}

func dailyReportCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daily [command]",
		Short: "run the daily report",
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

			s1, err := report.Daily(ctx, db, metric)
			if err != nil {
				return err
			}
			cmd.Print(s1)

			return nil
		},
	}
	return cmd
}

func streakReportCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "streak",
		Short: "run the streak report",
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

			if err = report.Streak(db, m); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func monthlyReportCmd(cli *cli) *cobra.Command {
	var flags struct {
		StartDate string
		EndDate   string
	}
	var cmd = &cobra.Command{
		Use:   "monthly",
		Short: "run the monthly report",
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

			configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
			if err != nil && err != store.ErrNotFound {
				return err
			}
			if err != nil && err == store.ErrNotFound {
				return fmt.Errorf("Missing config option metric_type: %s", metric.Name)
			}

			output := ""

			start, err := time.Parse("2006-01-02", flags.StartDate)
			if err != nil {
				return err
			}

			end, err := time.Parse("2006-01-02", flags.EndDate)
			if err != nil {
				return err
			}

			r := report.NewReport(db, metric)
			switch configMetricType {
			case "counter":
				output, err = r.MonthlyCounter(ctx, report.WithStartTimestamp(start), report.WithEndTimestamp(end))
				if err != nil {
					return err
				}
			case "gauge":
				output, err = r.MonthlyGuage(ctx, report.WithStartTimestamp(start), report.WithEndTimestamp(end))
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unsupported reporting for metric type: %s", configMetricType)
			}

			cmd.Print(output)
			return nil
		},
	}
	cmd.Flags().StringVar(&flags.StartDate, "start", time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), "Start date to report on")
	cmd.Flags().StringVar(&flags.EndDate, "end", time.Now().Format("2006-01-02"), "End date to report on")

	return cmd
}

func weeklyReportCmd(cli *cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "weekly",
		Short: "run the weekly report",
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

			configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
			if err != nil && err != store.ErrNotFound {
				return err
			}
			if err != nil && err == store.ErrNotFound {
				return err
			}

			switch configMetricType {
			case "counter":
				if err = report.WeeklyCounter(ctx, db, metric); err != nil {
					return err
				}
			case "gauge":
				if err = report.WeeklyGauge(ctx, db, metric); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unsupported reporting for metric type: %s", configMetricType)
			}

			return nil
		},
	}
	return cmd
}
