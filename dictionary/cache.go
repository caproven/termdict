package dictionary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/caproven/termdict/config"
)

const cacheDir string = "cache"

// defaultCacheDir returns the default subdirectory where
// the dictionary cache may be stored
func defaultCacheDir() string {
	return filepath.Join(config.DefaultConfigDir(), cacheDir)
}

type FileCache struct {
	dir string
}

// NewFileCache returns a new file-based dictionary cache. A default is
// supplied if no dir is passed
func NewFileCache(dir string) (*FileCache, error) {
	if dir == "" {
		dir = defaultCacheDir()
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return &FileCache{dir: dir}, nil
}

func (c *FileCache) Lookup(_ context.Context, word string) ([]Definition, error) {
	path := c.PathFor(word)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var defs []Definition

	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, err
	}

	return defs, nil
}

func (c *FileCache) Contains(_ context.Context, word string) (bool, error) {
	path := c.PathFor(word)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *FileCache) Save(_ context.Context, word string, defs []Definition) error {
	data, err := json.Marshal(defs)
	if err != nil {
		return err
	}

	path := c.PathFor(word)

	return os.WriteFile(path, data, os.ModePerm)
}

func (c *FileCache) PathFor(word string) string {
	return fmt.Sprintf("%s/%s.json", c.dir, word)
}

type Definer interface {
	Define(ctx context.Context, word string) ([]Definition, error)
}

type Cache interface {
	Lookup(ctx context.Context, word string) ([]Definition, error)
	Contains(ctx context.Context, word string) (bool, error)
	Save(ctx context.Context, word string, defs []Definition) error
}

type CachedDefiner struct {
	cache    Cache
	fallback Definer
}

func NewCachedDefiner(c Cache, d Definer) *CachedDefiner {
	return &CachedDefiner{
		cache:    c,
		fallback: d,
	}
}

func (d *CachedDefiner) Define(ctx context.Context, word string) ([]Definition, error) {
	ok, err := d.cache.Contains(ctx, word)
	if err != nil {
		return nil, err
	}
	if ok {
		defs, err := d.cache.Lookup(ctx, word)
		if err != nil {
			return nil, err
		}
		return defs, nil
	}

	defs, err := d.fallback.Define(ctx, word)
	if err != nil {
		return nil, err
	}
	if err = d.cache.Save(ctx, word, defs); err != nil {
		return nil, err
	}
	return defs, nil
}
