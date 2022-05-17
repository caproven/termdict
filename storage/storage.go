package storage

import (
	"os"
	"path/filepath"
)

const configSubdir string = "termdict"

// CreateConfigDir creates a config directory for the application
func CreateConfigDir() error {
	return os.MkdirAll(defaultConfigDir(), os.ModePerm)
}

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}
