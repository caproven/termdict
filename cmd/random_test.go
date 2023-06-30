package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

func TestRandomCmd(t *testing.T) {
	dict := memoryDefiner{
		"kappa": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology."},
		},
		"senescence": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "definition 1"},
			{PartOfSpeech: "adjective", Meaning: "definition 2"},
			{PartOfSpeech: "verb", Meaning: "definition 3"},
		},
	}

	cases := []struct {
		name        string
		cmd         string
		initList    vocab.List
		expectedOut string
		errExpected bool
	}{
		{
			name:        "empty list",
			cmd:         "random",
			initList:    vocab.List{Words: []string{}},
			expectedOut: "no words in vocab list\n",
			errExpected: false,
		},
		{
			name:        "single word",
			cmd:         "random",
			initList:    vocab.List{Words: []string{"kappa"}},
			expectedOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
			errExpected: false,
		},
		{
			name:        "valid definitions limit",
			cmd:         "random --limit 2",
			initList:    vocab.List{Words: []string{"senescence"}},
			expectedOut: "senescence\n[noun] definition 1\n[adjective] definition 2\n",
			errExpected: false,
		},
		{
			name:        "ignored definitions limit",
			cmd:         "random",
			initList:    vocab.List{Words: []string{"senescence"}},
			expectedOut: "senescence\n[noun] definition 1\n[adjective] definition 2\n[verb] definition 3\n",
			errExpected: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			v := newMemoryVocabRepo(test.initList)

			var b bytes.Buffer

			cfg := Config{
				Out:   &b,
				Vocab: v,
				Dict:  dict,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err := cmd.Execute()
			out := b.String()

			if test.errExpected {
				if err == nil {
					t.Error("expected err but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("didn't expect err but got: %v", err)
				}
			}

			if out != test.expectedOut {
				t.Errorf("got %v, expected %v", out, test.expectedOut)
			}
		})
	}
}
