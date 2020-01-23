package cmd

import (
	"fmt"
	"os"
	"errors"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"ct/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type UserConfig struct {
	DbFile string `json:"db_file"`
}

var configureCmd = &cobra.Command{
	Use: "configure [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runConfigure(cfg, cmd.Flags())
	},
}

func runConfigure(cfg *config.Config, flags *pflag.FlagSet) error {
	var sqlStmt string

	usrCfg := cfg.UserViperConfig
	dbFile := usrCfg.GetString("ct.db_file")
	if dbFile == "" {
		return errors.New("db_file not set")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	metricName, _ := flags.GetString("metric")

	// No more configuration needed.
	if metricName == "" {
		return nil
	}

	metric := Metric{Name: metricName}
	metricID, err := getMetricID(db, metric)
	if err != nil {
		return err
	}

	supportedFrequencies := []string{
		"daily",
	}

	frequency, _ := flags.GetString("frequency")
	if frequency != "" {
		supportedFrequency := stringInSlice(frequency, supportedFrequencies)
		if !supportedFrequency {
			return errors.New("Frequency not supported")
		}
		sqlStmt = `
		INSERT INTO config
			(
				metric_id,
				opt,
				val
			)
			VALUES
			(
				?,
				"frequency",
				?
			)
			ON CONFLICT(metric_id, opt)
			DO UPDATE SET val=?
		`
		stmt, err := db.Prepare(sqlStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err = stmt.Exec(metricID, frequency, frequency); err != nil {
			return err
		}

	}

	valueText, _ := flags.GetString("value-text")
	if valueText != "" {
		sqlStmt = `
		INSERT INTO config
			(
				metric_id,
				opt,
				val
			)
			VALUES
			(
				?,
				"value_text",
				?
			)
			ON CONFLICT(metric_id, opt)
			DO UPDATE SET val=?
		`
		stmt, err := db.Prepare(sqlStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err = stmt.Exec(metricID, valueText, valueText); err != nil {
			return err
		}
	}

	supportedDataTypes := []string{
		"int",
		"float",
	}
	dataType, _ := flags.GetString("data-type")
	if dataType != "" {
		supportedDataType := stringInSlice(dataType, supportedDataTypes)
		if !supportedDataType {
			return errors.New("Data type not supported")
		}
		sqlStmt = `
		INSERT INTO config
			(
				metric_id,
				opt,
				val
			)
			VALUES
			(
				?,
				"data_type",
				?
			)
			ON CONFLICT(metric_id, opt)
			DO UPDATE SET val=?
		`
		stmt, err := db.Prepare(sqlStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err = stmt.Exec(metricID, dataType, dataType); err != nil {
			return err
		}
	}

	return nil
}

func initConfigureCmd() {
	c := configureCmd
	f := c.Flags()
	f.String("config-file", "", "Config file")
	f.String("metric", "", "Metric")
	c.MarkFlagRequired("metric")
	f.String("frequency", "", "Metric Frequency")
	f.String("data-type", "", "Metric Data Type")
	f.String("value-text", "", "Metric Value Text")
}

func init() {
	initConfigureCmd()
	rootCmd.AddCommand(configureCmd)

}
