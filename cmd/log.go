package cmd

import (
	"errors"
	"ct/config"
	"database/sql"
	"strconv"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var logCmd = &cobra.Command{
	Use: "log [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runLogCmd(cfg, cmd.Flags(), args)
	},
}

func runLogCmd(cfg *config.Config, flags *pflag.FlagSet, args []string) error {
	var sqlStmt string

	usrCfg := cfg.UserViperConfig
	metric := Metric{}

	db, err := sql.Open("sqlite3", usrCfg.GetString("ct.db_file"))
	if err != nil {
		return err
	}
	defer db.Close()

	if metric.Name, err = flags.GetString("metric"); err != nil {
		return err
	}

	value, err := flags.GetString("value")
	if err != nil {
		return err
	}

	quiet, err := flags.GetBool("quiet")
	if err != nil {
		return err
	}

	timestamp, err := flags.GetString("timestamp")
	if err != nil {
		return err
	}

	ts := time.Now()
	if timestamp != "" {
		ts, err = parseTimestamp(timestamp)
		if err != nil {
			return err
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	metric.ID, err = getMetricID(db, metric)
	if err != nil {
		return err
	}

	metric.Config, err = getMetricConfig(db, metric)
	if err != nil {
		return err
	}

	value, err = getValueFromConsole(value, metric.Config.ValueText)
	if err != nil {
		return err
	}

	if metric.Config.Frequency == "daily" {
		var count int


		sqlStmt = `
		SELECT COUNT(1)
			FROM ct
			WHERE metric_id = ?
			AND
			timestamp >= datetime(?, '-1 day')
		`
		stmt, err := db.Prepare(sqlStmt)
		if err != nil {
			return err
		}
		if err = stmt.QueryRow(metric.ID, ts.Format("2006-01-02")).Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			if !quiet {
				return errors.New("Already logged metric within frequency")
			}
			return nil
		}
	}


	switch metric.Config.DataType {
	case "int":
		_, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return errors.New("Value not an int")
		}
	case "float":
		_, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return errors.New("Value not a float")
		}
	default:
	}

	sqlStmt = `
	INSERT INTO ct
		(
			id,
			timestamp,
			metric_id,
			value
		)
		VALUES
		(
			NULL,
			?,
			?,
			?
		)
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(ts.Format("2006-01-02"), metric.ID, value); err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func initLogCmd() {
	c := logCmd
	f := c.Flags()
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("value", "", "Value")
	f.String("config-file", "", "")
	f.Bool("quiet", false, "")
	f.String("timestamp", "", "")
}

func init() {
	initLogCmd()
	rootCmd.AddCommand(logCmd)
}
