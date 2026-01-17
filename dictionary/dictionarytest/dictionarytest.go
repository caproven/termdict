package dictionarytest

import (
	"context"
	"fmt"

	"github.com/caproven/termdict/dictionary"
)

type InMemoryDefiner map[string][]dictionary.Definition

func (m InMemoryDefiner) Define(_ context.Context, word string) ([]dictionary.Definition, error) {
	defs, ok := m[word]
	if !ok {
		return nil, fmt.Errorf("word '%s' not found", word)
	}
	return defs, nil
}
