package cmd

import (
	"ct/config"
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
	var sqlStmt string

	usrCfg := cfg.UserViperConfig

	db, err := sql.Open("sqlite3", usrCfg.GetString("ct.db_file"))
	if err != nil {
		return err
	}
	defer db.Close()

	metricName, err := flags.GetString("metric")
	if err != nil {
		return err
	}

	metric := Metric{Name: metricName}

	metricID, err := getMetricID(db, metric)
	if err != nil {
		return err
	}

	metric.Config, err = getMetricConfig(db, metric)
	if err != nil {
		return err
	}

	sqlStmt = `
	SELECT strftime("%Y-%m-%d", timestamp), value
		FROM ct
		WHERE metric_id = ?
		ORDER BY timestamp
		DESC
	`

	rows, err := db.Query(sqlStmt, metricID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp string
		var value string
		if err := rows.Scan(&timestamp, &value); err != nil {
			return err
		}
		fmt.Println(timestamp, metric.Name, value)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func initReportCmd() {
	c := reportCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("config-file", "", "")
}

func init() {
	initReportCmd()
	rootCmd.AddCommand(reportCmd)
}
