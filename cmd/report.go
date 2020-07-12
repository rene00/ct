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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runReportCmd(cfg, cmd.Flags(), args)
	},
}

func runReportCmd(cfg *config.Config, flags *pflag.FlagSet, args []string) error {

	usrCfg := cfg.UserViperConfig

	db, err := sql.Open("sqlite3", usrCfg.GetString("ct.db_file"))
	if err != nil {
		return err
	}
	defer db.Close()

	reportType, err := flags.GetString("report-type")
	if err != nil {
		return err
	}

	switch rt := reportType; rt {
	case "monthly-average":
		if err = report.MonthlyAverage(db, flags); err != nil {
			return err
		}
	case "all":
		if err = report.All(db, flags); err != nil {
			return err
		}
	default:
		return fmt.Errorf("report type %s not supported", rt)
	}

	return nil
}

func initReportCmd() {
	c := reportCmd
	f := c.Flags()
	f.String("report-type", "", "Report type")
	c.MarkFlagRequired("report-type")
	f.StringSlice("metrics", []string{}, "Metrics")
	f.String("config-file", "", "")
}

func init() {
	initReportCmd()
	rootCmd.AddCommand(reportCmd)
}
