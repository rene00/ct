package cli

import (
	"bytes"
	"ct/db/migrations"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initTestConfig(t *testing.T) *cli {
	t.Helper()
	dir, err := ioutil.TempDir("", "ct")
	if err != nil {
		t.Fatal(err)
	}
	config := config{DBFile: path.Join(dir, "ct.db")}
	cli := &cli{configFile: path.Join(dir, "config.json"), config: config}
	if err := cli.persistConfig(); err != nil {
		t.Fatal(err)
	}
	if err := migrations.DoMigrateDb(fmt.Sprintf("sqlite3://%s", config.DBFile)); err != nil {
		t.Fatal(err)
	}
	return cli
}

func initTest(t *testing.T, createMetricArgs [][]string, createLogArgs [][]string) *cli {
	t.Helper()
	cli := initTestConfig(t)

	for _, i := range createMetricArgs {
		cmd := createMetricCmd(cli)
		cmd.SetArgs(i)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}

	for _, i := range createLogArgs {
		cmd := createLogCmd(cli)
		cmd.SetArgs(i)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		configDir := path.Dir(cli.configFile)
		for _, i := range []string{cli.configFile, path.Join(configDir, "ct.db"), configDir} {
			if err := os.Remove(i); err != nil {
				t.Fatal(err)
			}
		}
	})

	return cli
}

func dumpTest(t *testing.T, cli *cli) dumpOutput {
	t.Helper()
	cmd := dumpCmd(cli)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_, err := cmd.ExecuteC()
	if err != nil {
		t.Fatal(err)
	}
	dumpOutput := dumpOutput{}
	if err := json.Unmarshal([]byte(buf.String()), &dumpOutput); err != nil {
		t.Fatal(err)
	}
	return dumpOutput
}
