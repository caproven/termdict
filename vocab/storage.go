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
	Read() (List, error)
	Write(List) error
}

// File represents a vocab list on the filesystem
type File struct {
	Path string
}

// Read reads a vocab list from the filesystem
func (vf File) Read() (List, error) {
	var vl List

	if _, err := os.Stat(vf.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty list
			return vl, nil
		}
		return List{}, err
	}

	data, err := os.ReadFile(vf.Path)
	if err != nil {
		return List{}, err
	}

	if err := json.Unmarshal(data, &vl.Words); err != nil {
		return List{}, err
	}

	return vl, nil
}

// Write writes a vocab list to the filesystem
func (vf File) Write(vl List) error {
	data, err := json.Marshal(vl.Words)
	if err != nil {
		return err
	}

	if err := createConfigDir(); err != nil {
		return err
	}

	return os.WriteFile(vf.Path, data, os.ModePerm)
}

// DefaultFilepath returns the default filepath for where
// the vocab list may be stored on the filesystem
func DefaultFilepath() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}

func createConfigDir() error {
	return os.MkdirAll(defaultConfigDir(), os.ModePerm)
}
