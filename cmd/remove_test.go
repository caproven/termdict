package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary/dictionarytest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRemoveCmd(t *testing.T) {
	t.Run("no words specified", func(t *testing.T) {
		cfg := Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list", "remove"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("failure removing from list", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("RemoveWordsFromList", mock.Anything, mock.Anything).Return(errors.New("failure")).Once()

		cfg := Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list", "remove", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("remove single word", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("RemoveWordsFromList", mock.Anything, mock.MatchedBy(func(words []string) bool {
			return reflect.DeepEqual(words, []string{"cucumber"})
		})).Return(nil).Once()

		cfg := Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list", "remove", "cucumber"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("remove multiple words", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("RemoveWordsFromList", mock.Anything, mock.MatchedBy(func(words []string) bool {
			return reflect.DeepEqual(words, []string{"kappa", "cucumber"})
		})).Return(nil).Once()

		cfg := Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list", "remove", "kappa", "cucumber"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
}
