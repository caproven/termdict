package cmd

import (
	"fmt"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

type memoryDefiner map[string][]dictionary.Definition

func (m memoryDefiner) Define(word string) ([]dictionary.Definition, error) {
	defs, ok := m[word]
	if !ok {
		return nil, fmt.Errorf("word '%s' not found", word)
	}
	return defs, nil
}

type memoryVocabRepo struct {
	list vocab.List
}

func newMemoryVocabRepo(init vocab.List) *memoryVocabRepo {
	return &memoryVocabRepo{list: init}
}

func (mvr *memoryVocabRepo) Load() (vocab.List, error) {
	list := vocab.List{Words: []string{}}
	list.Words = append(list.Words, mvr.list.Words...)
	return list, nil
}

func (mvr *memoryVocabRepo) Save(vl vocab.List) error {
	mvr.list = vl
	return nil
}
