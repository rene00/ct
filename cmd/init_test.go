package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"ct/config"
	"testing"
	"io/ioutil"
	"path/filepath"
	"database/sql"
	"os"
)

func TestRunInitCmd(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "run-init-cmd")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cfg := config.Config{
		Persister: config.InMemoryPersister{},
		UserViperConfig: viper.New(),
	}

	dbFile := filepath.Join(tmpDir, "ct.db")
	flags := pflag.NewFlagSet("test", pflag.PanicOnError)
	flags.String("db-file", dbFile, "")

	err = runInitCmd(&cfg, flags, []string{})
	if err != nil {
		t.Error(err.Error())
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	if _, err := db.Exec("DROP TABLE config; DROP TABLE ct; DROP TABLE metric"); err != nil {
		t.Error(err.Error())
	}
}

