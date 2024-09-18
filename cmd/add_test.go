package cmd

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

func TestAddCmd(t *testing.T) {
	dict := memoryDefiner{
		"kappa": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology."},
		},
		"cucumber": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus."},
			{PartOfSpeech: "noun", Meaning: "The edible fruit of this plant, having a green rind and crisp white flesh."},
		},
		"terminal": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A building in an airport where passengers transfer from ground transportation to the facilities that allow them to board airplanes."},
		},
		"dictionary": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A reference work with a list of words from one or more languages, normally ordered alphabetically, explaining each word's meaning."},
		},
	}

	cases := []struct {
		name         string
		cmd          string
		initList     vocab.List
		expectedList vocab.List
		errExpected  bool
	}{
		{
			name:         "word to empty",
			cmd:          "list add kappa",
			initList:     vocab.List{Words: []string{}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  false,
		},
		{
			name:         "multiple words to empty",
			cmd:          "list add terminal dictionary",
			initList:     vocab.List{Words: []string{}},
			expectedList: vocab.List{Words: []string{"terminal", "dictionary"}},
			errExpected:  false,
		},
		{
			name:         "new word to existing",
			cmd:          "list add cucumber",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa", "cucumber"}},
			errExpected:  false,
		},
		{
			name:         "multiple new words to existing",
			cmd:          "list add dictionary terminal",
			initList:     vocab.List{Words: []string{"kappa", "cucumber"}},
			expectedList: vocab.List{Words: []string{"kappa", "cucumber", "dictionary", "terminal"}},
			errExpected:  false,
		},
		{
			name:         "duplicate",
			cmd:          "list add terminal",
			initList:     vocab.List{Words: []string{"cucumber", "terminal"}},
			expectedList: vocab.List{Words: []string{"cucumber", "terminal"}},
			errExpected:  true,
		},
		{
			name:         "multiple words with duplicate",
			cmd:          "list add cucumber dictionary",
			initList:     vocab.List{Words: []string{"dictionary", "kappa"}},
			expectedList: vocab.List{Words: []string{"dictionary", "kappa"}},
			errExpected:  true,
		},
		{
			name:         "unknown word",
			cmd:          "list add asdf",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  true,
		},
		{
			name:         "multiple words with unknown",
			cmd:          "list add cucumber asdf terminal",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  true,
		},
		{
			name:         "unknown word with check disabled",
			cmd:          "list add asdf --no-check",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa", "asdf"}},
			errExpected:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			v := newMemoryVocabRepo(test.initList)

			cfg := Config{
				Out:   os.Stdout,
				Vocab: v,
				Dict:  dict,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err := cmd.Execute()

			if test.errExpected {
				if err == nil {
					t.Error("expected err but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("didn't expect err but got: %v", err)
				}
			}

			got, err := v.Load()
			if err != nil {
				t.Errorf("failed to load storage after executing command: %v", err)
			}
			if !reflect.DeepEqual(got, test.expectedList) {
				t.Errorf("got %v, expected %v", got, test.expectedList)
			}
		})
	}
}
