package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caproven/termdict/vocab"
)

func TestListCmd(t *testing.T) {
	cases := []struct {
		name        string
		cmd         string
		initList    vocab.List
		expectedOut string
		errExpected bool
	}{
		{
			name:        "empty list",
			cmd:         "list",
			initList:    vocab.List{Words: []string{}},
			expectedOut: "no words in vocab list\n",
			errExpected: false,
		},
		{
			name:        "single word",
			cmd:         "list",
			initList:    vocab.List{Words: []string{"kappa"}},
			expectedOut: "kappa\n",
			errExpected: false,
		},
		{
			name:        "multiple words",
			cmd:         "list",
			initList:    vocab.List{Words: []string{"kappa", "cucumber", "terminal", "dictionary"}},
			expectedOut: "kappa\ncucumber\nterminal\ndictionary\n",
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
				Dict:  memoryDefiner{},
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
