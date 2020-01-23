package cmd

import (
	"ct/config"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"path/filepath"
	"os"
)

var initCmd = &cobra.Command{
	Use: "init [command]",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(cmd.Flags())
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%v\n", err))
			return err
		}
		return runInitCmd(cfg, cmd.Flags(), args)
	},
}

// The init command:
// - creates the ct config file
// - creates the sqlite db
func runInitCmd(cfg *config.Config, flags *pflag.FlagSet, args []string) error {

	var sqlStmt string

	dbFile, _ := flags.GetString("db-file")
	if dbFile == "" {
		dbFile = filepath.Join(cfg.Dir, "ct.db")
	}

	usrCfg := cfg.UserViperConfig
	usrCfg.Set("ct", UserConfig{DbFile: dbFile})

	if err := cfg.Save("ct"); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt = `
	CREATE TABLE
		IF NOT EXISTS metric
		(
			id INTEGER NOT NULL PRIMARY KEY,
			name text
		)
	`
	if _, err := db.Exec(sqlStmt); err != nil {
		return err
	}

	sqlStmt = `
	CREATE TABLE
		IF NOT EXISTS ct
		(
			id INTEGER NOT NULL PRIMARY KEY,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			metric_id INTEGER NOT NULL,
			value int NOT NULL,
			FOREIGN KEY(metric_id) REFERENCES metric(id)
		)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	CREATE TABLE
		IF NOT EXISTS config 
		(
			metric_id INTEGER NOT NULL,
			opt text NOT NULL,
			val text NOT NULL,
			UNIQUE(metric_id, opt)
				ON CONFLICT REPLACE,
			FOREIGN KEY(metric_id) REFERENCES metric(id)
		)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func initInitCmd() {
	c := initCmd
	f := c.Flags()
	f.String("db-file", "", "DB file")
	f.String("config-file", "", "Config file")
}

func init() {
	initInitCmd()
	rootCmd.AddCommand(initCmd)
}
