package cmd

import (
	"ct/config"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var dumpCmd = &cobra.Command{
	Use: "dump [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runDumpCmd(cfg, cmd.Flags(), args)
	},
}

func runDumpCmd(cfg *config.Config, flags *pflag.FlagSet, args []string) error {
	var sqlStmt string

	usrCfg := cfg.UserViperConfig

	db, err := sql.Open("sqlite3", usrCfg.GetString("ct.db_file"))
	if err != nil {
		return err
	}
	defer db.Close()

	metrics := []*Metric{}

	sqlStmt = `SELECT id, name from metric`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		metric := &Metric{ID: id, Name: name}
		metrics = append(metrics, metric)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	for _, m := range metrics {
		sqlStmt = `SELECT opt, val FROM config WHERE metric_id = ?`
		rows, err := db.Query(sqlStmt, m.ID)
		if err != nil {
			return err
		}
		defer rows.Close()

		metricConfig := MetricConfig{}

		for rows.Next() {
			var opt string
			var val string

			if err := rows.Scan(&opt, &val); err != nil {
				return err
			}

			switch opt {
			case "frequency":
				metricConfig.Frequency = val
			case "value_text":
				metricConfig.ValueText = val
			default:
				return errors.New("Unsupported config option")
			}

		}

		err = rows.Err()
		if err != nil {
			return err
		}

		m.Config = metricConfig

		sqlStmt = `SELECT timestamp, value FROM ct WHERE metric_id = ?`
		rows, err = db.Query(sqlStmt, m.ID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var timestamp time.Time
			var value float64

			metricData := MetricData{}

			if err := rows.Scan(&timestamp, &value); err != nil {
				return err
			}

			metricData.Timestamp = timestamp
			metricData.Value = value
			m.Data = append(m.Data, metricData)
		}

		err = rows.Err()
		if err != nil {
			return err
		}
	}

	dump, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", dump)

	return nil
}

func initDumpCmd() {
	c := dumpCmd
	f := c.Flags()
	f.String("config-file", "", "")
}

func init() {
	initDumpCmd()
	rootCmd.AddCommand(dumpCmd)
}
