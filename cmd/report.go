package cmd

import (
	"context"
	"ct/config"
	"ct/internal/report"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type reportCmdContext struct {
	flags  *pflag.FlagSet
	usrCfg *viper.Viper
	DB     *sql.DB
}

func newReportCmdContext(usrCfg *viper.Viper, flags *pflag.FlagSet) *reportCmdContext {
	return &reportCmdContext{
		flags:  flags,
		usrCfg: usrCfg,
	}
}

var reportCmd = &cobra.Command{
	Use: "report [command]",
}

var reportDailyCmd = &cobra.Command{
	Use:   "daily [command]",
	Short: "run the daily report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metric", cmd.Flags().Lookup("metric"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		defer db.Close()

		s := store.NewStore(db)
		ctx := context.Background()
		metric, err := s.Metric.SelectOne(ctx, viper.GetString("metric"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		if err = report.Daily(ctx, db, metric); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		return nil
	},
}

var reportStreakCmd = &cobra.Command{
	Use:   "streak [command]",
	Short: "run the streak report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metric", cmd.Flags().Lookup("metric"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		defer db.Close()

		if err = report.Streak(db, viper.GetString("metric")); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		return nil
	},
}

var reportMonthlyCmd = &cobra.Command{
	Use:   "monthly [command]",
	Short: "run the monthly report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metric", cmd.Flags().Lookup("metric"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		defer db.Close()

		s := store.NewStore(db)
		ctx := context.Background()
		metric, err := s.Metric.SelectOne(ctx, viper.GetString("metric"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
		if err != nil && err != store.ErrNotFound {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		if err != nil && err == store.ErrNotFound {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Missing config option metric_type: %s", metric.Name))
			os.Exit(1)
		}

		switch configMetricType {
		case "counter":
			if err = report.MonthlyCounter(ctx, db, metric); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
				os.Exit(1)
			}
		case "gauge":
			if err = report.MonthlyGauge(ctx, db, metric); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Unsupported reporting for metric type: %s", configMetricType))
			os.Exit(1)
		}

		return nil
	},
}

var reportWeeklyCmd = &cobra.Command{
	Use:   "weekly [command]",
	Short: "run the weekly report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metric", cmd.Flags().Lookup("metric"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		db, err := sql.Open("sqlite3", cfg.UserViperConfig.GetString("ct.db_file"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		defer db.Close()

		s := store.NewStore(db)
		ctx := context.Background()
		metric, err := s.Metric.SelectOne(ctx, viper.GetString("metric"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		configMetricType, err := s.Config.SelectOne(ctx, metric.MetricID, "metric_type")
		if err != nil && err != store.ErrNotFound {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
		if err != nil && err == store.ErrNotFound {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Missing config option metric_type: %s", metric.Name))
			os.Exit(1)
		}

		switch configMetricType {
		case "counter":
			if err = report.WeeklyCounter(ctx, db, metric); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
				os.Exit(1)
			}
		case "gauge":
			if err = report.WeeklyGauge(ctx, db, metric); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Unsupported reporting for metric type: %s", configMetricType))
			os.Exit(1)
		}

		return nil
	},
}

func initReportDailyCmd() {
	c := reportDailyCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("config-file", "", "")
}

func initReportMonthlyCmd() {
	c := reportMonthlyCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("config-file", "", "")
}

func initReportWeeklyCmd() {
	c := reportWeeklyCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("config-file", "", "")
}

func initReportStreakCmd() {
	c := reportStreakCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("config-file", "", "")
}

func init() {
	initReportDailyCmd()
	initReportWeeklyCmd()
	initReportMonthlyCmd()
	initReportStreakCmd()
	reportCmd.AddCommand(reportDailyCmd)
	reportCmd.AddCommand(reportWeeklyCmd)
	reportCmd.AddCommand(reportMonthlyCmd)
	reportCmd.AddCommand(reportStreakCmd)
	rootCmd.AddCommand(reportCmd)
}
