package config

import (
	"os"

	"github.com/spf13/viper"
)

// Persister saves viper configs.
type Persister interface {
	Save(*viper.Viper, string) error
}

// FilePersister saves viper configs to the file system.
type FilePersister struct {
	Dir            string
	ConfigFilePath string
}

// Save writes the viper config to the target location on the filesystem.
func (p FilePersister) Save(v *viper.Viper, basename string) error {
	if _, err := os.Stat(p.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(p.Dir, os.FileMode(0755)); err != nil {
			return err
		}
	}
	return v.WriteConfigAs(p.ConfigFilePath)
}

// InMemoryPersister is a noop persister for use in unit tests.
type InMemoryPersister struct{}

func (p InMemoryPersister) Save(*viper.Viper, string) error {
	return nil
}
