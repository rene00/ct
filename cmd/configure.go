package cmd

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"ct/config"
	"ct/internal/storage"
	"database/sql"

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
		if ok := stringInSlice(frequency, supportedFrequencies); !ok {
			return errors.New("Frequency not supported")
		}
		if err := storage.UpsertConfig(db, metricID, "frequency", frequency); err != nil {
			return err
		}
	}

	valueText, _ := flags.GetString("value-text")
	if valueText != "" {
		if err := storage.UpsertConfig(db, metricID, "value_text", valueText); err != nil {
			return err
		}
	}

	supportedDataTypes := []string{
		"int",
		"float",
	}
	dataType, _ := flags.GetString("data-type")
	if dataType != "" {
		if ok := stringInSlice(dataType, supportedDataTypes); !ok {
			return errors.New("Data type not supported")
		}
		if err := storage.UpsertConfig(db, metricID, "data_type", dataType); err != nil {
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
