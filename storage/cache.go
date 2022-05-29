package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/caproven/termdict/dictionary"
)

const cacheDir string = "cache"

var _ Cache = FileCache{}

// Cache for word definitions
type Cache interface {
	Contains(word string) (bool, error)
	Save(word string, defs []dictionary.Definition) error
	Lookup(word string) ([]dictionary.Definition, error)
}

// DefaultCacheDir returns the default subdirectory where
// the dictionary cache may be stored
func DefaultCacheDir() string {
	return filepath.Join(defaultConfigDir(), cacheDir)
}

// FileCache stores word definitions on the filesystem
type FileCache struct {
	DirPath string
}

// Contains checks if a word is in the cache
func (fc FileCache) Contains(word string) (bool, error) {
	path := fc.fileForWord(word)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Save a word to the cache
func (fc FileCache) Save(word string, defs []dictionary.Definition) error {
	data, err := json.Marshal(defs)
	if err != nil {
		return err
	}

	path := fc.fileForWord(word)

	return os.WriteFile(path, data, os.ModePerm)
}

// Lookup the definitions for a word in the cache
func (fc FileCache) Lookup(word string) ([]dictionary.Definition, error) {
	path := fc.fileForWord(word)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var defs []dictionary.Definition

	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, err
	}

	return defs, nil
}

func (fc FileCache) fileForWord(word string) string {
	return fmt.Sprintf("%s/%s.json", fc.DirPath, word)
}
