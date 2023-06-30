package config

import (
	"os"
	"path/filepath"
)

const appName string = "termdict"

func DefaultConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(configDir, appName)
}
