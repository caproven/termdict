package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/caproven/termdict/dictionary"
)

const cacheFile string = "cache.json"

// CacheRepo lets you fetch or save a dictionary cache
type CacheRepo struct {
	Path string
}

// Load a dictionary cache from storage
func (c CacheRepo) Load() (dictionary.Cache, error) {
	if _, err := os.Stat(c.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty map
			return make(dictionary.Cache), nil
		}
		return nil, err
	}

	data, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}

	var cache dictionary.Cache

	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return cache, nil
}

// Save a dictionary cache to storage
func (c CacheRepo) Save(cache dictionary.Cache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(c.Path, data, os.ModePerm)
}

// DefaultCacheFilepath returns the default filepath for where
// the dictionary cache may be stored on the filesystem
func DefaultCacheFilepath() string {
	return filepath.Join(defaultConfigDir(), cacheFile)
}
