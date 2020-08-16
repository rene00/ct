package cmd

import (
	"ct/config"
	"ct/internal/report"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
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

var reportAllCmd = &cobra.Command{
	Use:   "all [command]",
	Short: "run the all report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metrics", cmd.Flags().Lookup("metrics"))
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

		if err = report.All(db, cmd.Flags()); err != nil {
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

var reportMonthlyAverageCmd = &cobra.Command{
	Use:   "monthly-average [command]",
	Short: "run the monthly average report",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("metrics", cmd.Flags().Lookup("metrics"))
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

		if err = report.MonthlyAverage(db, viper.GetStringSlice("metrics")); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}

		return nil
	},
}

func initReportAllCmd() {
	c := reportAllCmd
	f := c.Flags()
	f.StringSlice("metrics", []string{}, "Metrics")
	c.MarkFlagRequired("metrics")
	f.String("config-file", "", "")
}

func initReportMonthlyAverageCmd() {
	c := reportMonthlyAverageCmd
	f := c.Flags()
	f.StringSlice("metrics", []string{}, "Metrics")
	c.MarkFlagRequired("metrics")
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
	initReportAllCmd()
	initReportMonthlyAverageCmd()
	initReportStreakCmd()
	reportCmd.AddCommand(reportAllCmd)
	reportCmd.AddCommand(reportMonthlyAverageCmd)
	reportCmd.AddCommand(reportStreakCmd)
	rootCmd.AddCommand(reportCmd)
}
