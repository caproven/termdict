package vocab

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

const configSubdir string = "termdict"
const vocabFile string = "vocab.json"

// Storage is an interface for the storage of a vocab list
type Storage interface {
	// Load a vocab list from storage
	Load() (List, error)
	// Save a vocab list to storage
	Save(List) error
}

// File represents vocab list storage on the filesystem. Implements
// the Storage interface
type File struct {
	Path string
}

// Load a vocab list from a file
func (f File) Load() (List, error) {
	var vl List

	if _, err := os.Stat(f.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty list
			return vl, nil
		}
		return List{}, err
	}

	data, err := os.ReadFile(f.Path)
	if err != nil {
		return List{}, err
	}

	if err := json.Unmarshal(data, &vl.Words); err != nil {
		return List{}, err
	}

	return vl, nil
}

// Save a vocab list to a file
func (f File) Save(vl List) error {
	data, err := json.Marshal(vl.Words)
	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, data, os.ModePerm)
}

// DefaultFilepath returns the default filepath for where
// the vocab list may be stored on the filesystem
func DefaultFilepath() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}

// CreateConfigDir creates a config directory for the application
func CreateConfigDir() error {
	return os.MkdirAll(defaultConfigDir(), os.ModePerm)
}

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}
