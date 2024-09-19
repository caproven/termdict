package cmd

import (
	"github.com/caproven/termdict/vocab"
)

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
