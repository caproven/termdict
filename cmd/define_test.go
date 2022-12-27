package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
)

func TestDefineCmd(t *testing.T) {
	dict := memoryDefiner{
		"kappa": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology."},
		},
		"cucumber": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus."},
			{PartOfSpeech: "noun", Meaning: "The edible fruit of this plant, having a green rind and crisp white flesh."},
		},
	}

	cases := []struct {
		name    string
		cmd     string
		cache   Cache
		word    string
		wantOut string
		wantErr bool
	}{
		{
			name:    "word found but not in cache",
			cmd:     "define kappa",
			cache:   newMemoryCache(nil),
			word:    "kappa",
			wantOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
		},
		{
			name: "word found and in cache",
			cmd:  "define kappa",
			cache: newMemoryCache(map[string][]dictionary.Definition{
				"kappa": {{PartOfSpeech: "noun", Meaning: "cache-specific meaning"}},
			}),
			word:    "kappa",
			wantOut: "kappa\n[noun] cache-specific meaning\n",
		},
		{
			name:    "word not found and not in cache",
			cmd:     "define asdf",
			cache:   newMemoryCache(nil),
			word:    "asdf",
			wantErr: true,
		},
		{
			name: "word not found but in cache",
			cmd:  "define sponge",
			cache: newMemoryCache(map[string][]dictionary.Definition{
				"sponge": {{PartOfSpeech: "noun", Meaning: "A piece of porous material used for washing"}},
			}),
			word:    "sponge",
			wantOut: "sponge\n[noun] A piece of porous material used for washing\n",
		},
		{
			name:    "valid definitions limit",
			cmd:     "define cucumber --limit 1",
			cache:   newMemoryCache(nil),
			word:    "cucumber",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n",
		},
		{
			name:    "ignored definitions limit",
			cmd:     "define cucumber --limit 0",
			cache:   newMemoryCache(nil),
			word:    "cucumber",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n[noun] The edible fruit of this plant, having a green rind and crisp white flesh.\n",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer

			cfg := Config{
				Out:   &b,
				Vocab: nil, // shouldn't be used by this cmd
				Cache: test.cache,
				Dict:  dict,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err := cmd.Execute()
			gotOut := b.String()

			if err != nil {
				if (err != nil) != test.wantErr {
					t.Errorf("define cmd error = %v, wantErr %v", err, test.wantErr)
				}
				return
			}

			if gotOut != test.wantOut {
				t.Errorf("got %v, expected %v", gotOut, test.wantOut)
			}

			found, _ := test.cache.Contains(test.word)
			if !found {
				t.Errorf("cache did not contain defined word %s", test.word)
			}
		})
	}
}
