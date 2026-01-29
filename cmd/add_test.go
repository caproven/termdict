package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddCmd(t *testing.T) {
	sampleDefs := []dictionary.Definition{
		{PartOfSpeech: "verb", Meaning: "to foo"},
		{PartOfSpeech: "verb", Meaning: "to bar"},
	}
	sampleErr := errors.New("failure")

	t.Run("no words specified", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  &mockDefiner{},
		})
		cmd.SetArgs([]string{"list", "add"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("failure adding words", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, mock.Anything).Return(sampleErr).Once()

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "foo").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("add word that cannot be defined", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "foo").Return(nil, sampleErr).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("add word that cannot be defined with no check", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, []string{"foo"}).Return(nil).Once()

		// Shouldn't be called
		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "foo", "--no-check"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("add word that can be defined", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, []string{"fortitude"}).Return(nil).Once()

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "fortitude").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "fortitude"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("add multiple words", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, []string{"porter", "placate"}).Return(nil).Once()

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "porter").Return(sampleDefs, nil).Once()
		definer.On("Define", mock.Anything, "placate").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "porter", "placate"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("add multiple words with one that cannot be defined", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "erudite").Return(sampleDefs, nil).Once()
		definer.On("Define", mock.Anything, "sanguine").Return(nil, sampleErr).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"list", "add", "erudite", "sanguine"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
