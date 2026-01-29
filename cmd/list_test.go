package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/caproven/termdict/dictionary/dictionarytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCmd(t *testing.T) {
	t.Run("failure fetching list", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return(nil, errors.New("failure")).Once()

		cfg := Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("empty list", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{}, nil).Once()

		var b bytes.Buffer
		cfg := Config{
			Out:   &b,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list"})

		err := cmd.Execute()
		require.NoError(t, err)

		assert.Equal(t, "no words in vocab list\n", b.String())
	})

	t.Run("single word", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{"kappa"}, nil).Once()

		var b bytes.Buffer
		cfg := Config{
			Out:   &b,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list"})

		err := cmd.Execute()
		require.NoError(t, err)

		assert.Equal(t, "kappa\n", b.String())
	})

	t.Run("multiple words", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{"kappa", "cucumber", "terminal", "dictionary"}, nil).Once()

		var b bytes.Buffer
		cfg := Config{
			Out:   &b,
			Vocab: vocabRepo,
			Dict:  dictionarytest.InMemoryDefiner{},
		}

		cmd := NewRootCmd(&cfg)
		cmd.SetArgs([]string{"list"})

		err := cmd.Execute()
		require.NoError(t, err)

		assert.Equal(t, "kappa\ncucumber\nterminal\ndictionary\n", b.String())
	})
}
