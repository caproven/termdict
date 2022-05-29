package storage

import (
	"os"
	"path/filepath"
)

const configSubdir string = "termdict"

// CreateConfigDir creates a config directory for the application
func CreateConfigDir() error {
	if err := os.MkdirAll(defaultConfigDir(), os.ModePerm); err != nil {
		return err
	}
	return os.MkdirAll(DefaultCacheDir(), os.ModePerm)
}

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}
