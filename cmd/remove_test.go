package cmd

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/caproven/termdict/vocab"
)

func TestRemoveCmd(t *testing.T) {
	cases := []struct {
		name         string
		cmd          string
		initList     vocab.List
		expectedList vocab.List
		errExpected  bool
	}{
		{
			name:         "from empty list",
			cmd:          "list remove cucumber",
			initList:     vocab.List{Words: []string{}},
			expectedList: vocab.List{Words: []string{}},
			errExpected:  true,
		},
		{
			name:         "from beginning",
			cmd:          "list remove dictionary",
			initList:     vocab.List{Words: []string{"dictionary", "cucumber", "kappa"}},
			expectedList: vocab.List{Words: []string{"cucumber", "kappa"}},
			errExpected:  false,
		},
		{
			name:         "from middle",
			cmd:          "list remove cucumber",
			initList:     vocab.List{Words: []string{"dictionary", "cucumber", "kappa"}},
			expectedList: vocab.List{Words: []string{"dictionary", "kappa"}},
			errExpected:  false,
		},
		{
			name:         "from end",
			cmd:          "list remove kappa",
			initList:     vocab.List{Words: []string{"dictionary", "cucumber", "kappa"}},
			expectedList: vocab.List{Words: []string{"dictionary", "cucumber"}},
			errExpected:  false,
		},
		{
			name:         "multiple words",
			cmd:          "list remove kappa terminal",
			initList:     vocab.List{Words: []string{"kappa", "terminal", "dictionary", "cucumber"}},
			expectedList: vocab.List{Words: []string{"dictionary", "cucumber"}},
			errExpected:  false,
		},
		{
			name:         "word that doesn't exist",
			cmd:          "list remove asdf",
			initList:     vocab.List{Words: []string{"dictionary", "cucumber", "kappa"}},
			expectedList: vocab.List{Words: []string{"dictionary", "cucumber", "kappa"}},
			errExpected:  true,
		},
		{
			name:         "multiple words not all exist",
			cmd:          "list remove cucumber terminal",
			initList:     vocab.List{Words: []string{"kappa", "cucumber", "dictionary"}},
			expectedList: vocab.List{Words: []string{"kappa", "cucumber", "dictionary"}},
			errExpected:  true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			v := newMemoryVocabRepo(test.initList)

			cfg := Config{
				Out:   os.Stdout,
				Vocab: v,
				Dict:  memoryDefiner{},
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
