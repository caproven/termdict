package cmd

import (
	"context"

	"github.com/caproven/termdict/dictionary"
	"github.com/stretchr/testify/mock"
)

type mockVocabRepo struct {
	mock.Mock
}

func (m *mockVocabRepo) AddWordsToList(ctx context.Context, words []string) ([]string, error) {
	args := m.Called(ctx, words)
	added, err := args.Get(0), args.Error(1)
	if added == nil {
		return nil, err
	}
	return added.([]string), err
}

func (m *mockVocabRepo) RemoveWordsFromList(ctx context.Context, words []string) ([]string, error) {
	args := m.Called(ctx, words)
	removed, err := args.Get(0), args.Error(1)
	if removed == nil {
		return nil, err
	}
	return removed.([]string), err
}

func (m *mockVocabRepo) GetWordsInList(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	words, err := args.Get(0), args.Error(1)
	if words == nil {
		return nil, err
	}
	return words.([]string), err
}

type mockDefiner struct {
	mock.Mock
}

func (m *mockDefiner) Define(ctx context.Context, word string) ([]dictionary.Definition, error) {
	args := m.Called(ctx, word)
	defs, err := args.Get(0), args.Error(1)
	if defs == nil {
		return nil, err
	}
	return defs.([]dictionary.Definition), err
}
