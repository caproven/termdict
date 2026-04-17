package cmd

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDefineCmd(t *testing.T) {
	sampleDefs := []dictionary.Definition{
		{PartOfSpeech: "noun", Meaning: "something"},
	}
	sampleErr := errors.New("failure")

	t.Run("no word specified", func(t *testing.T) {
		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  &mockDefiner{},
		})
		cmd.SetArgs([]string{"define"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("word not found", func(t *testing.T) {
		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "foo").Return(nil, sampleErr).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  definer,
		})
		cmd.SetArgs([]string{"define", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("word found", func(t *testing.T) {
		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "bar").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  definer,
		})
		cmd.SetArgs([]string{"define", "bar"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("random with empty list", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{}, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  &mockDefiner{},
		})
		cmd.SetArgs([]string{"define", "--random"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	// TODO create test for random from multiple entries (inject rand source)
	t.Run("random with single word in list", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{"a"}, nil).Once()

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "a").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"define", "--random"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("random flag cannot be given alongside a positional arg", func(t *testing.T) {
		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  &mockDefiner{},
		})
		cmd.SetArgs([]string{"define", "--random", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("invalid output format", func(t *testing.T) {
		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  &mockDefiner{},
		})
		cmd.SetArgs([]string{"define", "--output", "invalid", "foo"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("json output", func(t *testing.T) {
		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "b").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: &mockVocabRepo{},
			Dict:  definer,
		})
		cmd.SetArgs([]string{"define", "--output", "json", "b"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("random with output flag", func(t *testing.T) {
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("GetWordsInList", mock.Anything).Return([]string{"c"}, nil).Once()

		definer := &mockDefiner{}
		defer definer.AssertExpectations(t)
		definer.On("Define", mock.Anything, "c").Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  definer,
		})
		cmd.SetArgs([]string{"define", "--output", "json", "--random"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("random and save flags are mutually exclusive", func(t *testing.T) {
		cmd := NewRootCmd(&Config{
			Out: os.Stdout,
		})
		cmd.SetArgs([]string{"define", "--save", "--random"})
		require.ErrorContains(t, cmd.Execute(), "flags")
	})

	t.Run("save new word", func(t *testing.T) {
		word := "prism"
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, mock.MatchedBy(func(words []string) bool {
			return reflect.DeepEqual(words, []string{word})
		})).Return([]string{word}, nil).Once()

		dict := &mockDefiner{}
		defer dict.AssertExpectations(t)
		dict.On("Define", mock.Anything, word).Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dict,
		})
		cmd.SetArgs([]string{"define", word, "--save"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("unknown word is not saved", func(t *testing.T) {
		word := "growth"
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)

		dict := &mockDefiner{}
		defer dict.AssertExpectations(t)
		dict.On("Define", mock.Anything, word).Return(nil, sampleErr).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dict,
		})
		cmd.SetArgs([]string{"define", word, "--save"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("saving word fails", func(t *testing.T) {
		word := "cumulonimbus"
		vocabRepo := &mockVocabRepo{}
		defer vocabRepo.AssertExpectations(t)
		vocabRepo.On("AddWordsToList", mock.Anything, mock.Anything).Return(nil, sampleErr).Once()

		dict := &mockDefiner{}
		defer dict.AssertExpectations(t)
		dict.On("Define", mock.Anything, word).Return(sampleDefs, nil).Once()

		cmd := NewRootCmd(&Config{
			Out:   os.Stdout,
			Vocab: vocabRepo,
			Dict:  dict,
		})
		cmd.SetArgs([]string{"define", word, "--save"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}

func TestTextPrinter(t *testing.T) {
	cases := []struct {
		name        string
		word        string
		definitions []dictionary.Definition
		expected    string
	}{
		{
			name: "single entry",
			word: "sponge",
			definitions: []dictionary.Definition{
				{
					PartOfSpeech: "noun",
					Meaning:      "A piece of porous material used for washing",
				},
			},
			expected: `sponge
[noun] A piece of porous material used for washing
`,
		},
		{
			name: "multiple entries",
			word: "sponge",
			definitions: []dictionary.Definition{
				{
					PartOfSpeech: "noun",
					Meaning:      "A piece of porous material used for washing",
				},
				{
					PartOfSpeech: "verb",
					Meaning:      "To clean, soak up, or dab with a sponge",
				},
			},
			expected: `sponge
[noun] A piece of porous material used for washing
[verb] To clean, soak up, or dab with a sponge
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			printer := &textPrinter{}
			if err := printer.Print(&b, test.word, test.definitions); err != nil {
				t.Errorf("failed to print definition: %v", err)
			}

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestJsonPrinter(t *testing.T) {
	cases := []struct {
		name        string
		word        string
		definitions []dictionary.Definition
		expected    string
	}{
		{
			name: "single entry",
			word: "guava",
			definitions: []dictionary.Definition{
				{
					PartOfSpeech: "noun", Meaning: "A tropical tree or shrub of the myrtle family",
				},
			},
			expected: `{
	"Word": "guava",
	"Definitions": [
		{
			"PartOfSpeech": "noun",
			"Meaning": "A tropical tree or shrub of the myrtle family"
		}
	]
}
`,
		},
		{
			name: "multiple entries",
			word: "super",
			definitions: []dictionary.Definition{
				{PartOfSpeech: "adjective", Meaning: "Of excellent quality"},
				{PartOfSpeech: "adverb", Meaning: "Very; extremely"},
			},
			expected: `{
	"Word": "super",
	"Definitions": [
		{
			"PartOfSpeech": "adjective",
			"Meaning": "Of excellent quality"
		},
		{
			"PartOfSpeech": "adverb",
			"Meaning": "Very; extremely"
		}
	]
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			printer := new(jsonPrinter)
			if err := printer.Print(&b, test.word, test.definitions); err != nil {
				t.Errorf("failed to print definition: %v", err)
			}

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}
