package dictionary

import (
	"context"
)

type Definer interface {
	Define(ctx context.Context, word string) ([]Definition, error)
}

type Cache interface {
	LookupWord(ctx context.Context, word string) ([]Definition, error)
	ContainsWord(ctx context.Context, word string) (bool, error)
	SaveWord(ctx context.Context, word string, defs []Definition) error
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
	ok, err := d.cache.ContainsWord(ctx, word)
	if err != nil {
		return nil, err
	}
	if ok {
		defs, err := d.cache.LookupWord(ctx, word)
		if err != nil {
			return nil, err
		}
		return defs, nil
	}

	defs, err := d.fallback.Define(ctx, word)
	if err != nil {
		return nil, err
	}
	if err = d.cache.SaveWord(ctx, word, defs); err != nil {
		return nil, err
	}
	return defs, nil
}
