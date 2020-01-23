package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// DefaultDirName is the default name used for the config directory.
	DefaultDirName string
	viperConfig    *viper.Viper
)

type Config struct {
	// The users Operating System.
	OS string
	// The name of the directory where the configuration file will be
	// stored.
	Dir string
	// The name of the directory which matches the name of the ct binary
	// used to append to the config directory name.
	DefaultDirName  string
	UserViperConfig *viper.Viper
	Persister       Persister
}

func NewConfig(flags *pflag.FlagSet) (*Config, error) {

	dir := Dir()
	configName := "ct.json"
	configType := "json"

	if configFile, _ := flags.GetString("config-file"); configFile != "" {
		abs, err := filepath.Abs(configFile)
		if err != nil {
			return nil, err
		}
		if fmt.Sprintf("%s", filepath.Ext(abs)) != fmt.Sprintf(".%s", configType) {
			err = errors.New(fmt.Sprintf("File extension must be %s, %s", configType, configFile))
			return nil, err
		}
		dir = filepath.Dir(abs)
		configName = filepath.Base(abs)
	}

	configFilePath := filepath.Join(dir, configName)
	viperConfig := viper.New()
	viperConfig.AddConfigPath(dir)
	viperConfig.SetConfigName(configName)
	viperConfig.SetConfigFile(configFilePath)
	viperConfig.SetConfigType(configType)
	viperConfig.ReadInConfig()

	return &Config{
		OS:              runtime.GOOS,
		Dir:             dir,
		DefaultDirName:  DefaultDirName,
		Persister:       FilePersister{Dir: dir, ConfigFilePath: configFilePath},
		UserViperConfig: viperConfig,
	}, nil
}

func SetDefaultDirName(binaryName string) {
	binaryNameBase := filepath.Base(binaryName)
	// Rename binaryNameBase to ct when user runs cli via
	// `go run main.go`.
	if binaryNameBase == "main" {
		binaryNameBase = "ct"
	}
	DefaultDirName = strings.Replace(binaryNameBase, ".exe", "", 1)
}

func Dir() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir != "" {
			return filepath.Join(dir, DefaultDirName)
		}
	}

	dir = os.Getenv("CT_CONFIG_HOME")
	if dir != "" {
		return dir
	}

	dir = os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if dir != "" {
		return filepath.Join(dir, DefaultDirName)
	}

	dir, _ = os.Getwd()
	return dir
}

func (c Config) Save(basename string) error {
	return c.Persister.Save(c.UserViperConfig, basename)
}

// Validates the users configuration file.
func ValidateConfig(cfg *viper.Viper, flags *pflag.FlagSet) error {
	return nil
}
