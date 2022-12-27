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

type memoryCache struct {
	data map[string][]dictionary.Definition
}

func newMemoryCache(data map[string][]dictionary.Definition) memoryCache {
	if data == nil {
		data = make(map[string][]dictionary.Definition)
	}
	return memoryCache{data: data}
}

func (mc memoryCache) Contains(word string) (bool, error) {
	_, ok := mc.data[word]
	return ok, nil
}

func (mc memoryCache) Save(word string, defs []dictionary.Definition) error {
	mc.data[word] = defs
	return nil
}

func (mc memoryCache) Lookup(word string) ([]dictionary.Definition, error) {
	defs, ok := mc.data[word]
	if !ok {
		return nil, fmt.Errorf("word %s not found in cache", word)
	}
	return defs, nil
}
