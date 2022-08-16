package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

func TestRandomCmd(t *testing.T) {
	mockWords := map[string]string{
		"kappa":      `[{"word":"kappa","phonetic":"/ˈkæpə/","phonetics":[{"text":"/ˈkæpə/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/kappa-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=311306"}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A tortoise-like creature in the Japanese mythology.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/kappa"]}]`,
		"senescence": `[{"word":"senescence","phonetic":"/sɨnˈɛsəns/","meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"definition 1","synonyms":[],"antonyms":[]},{"definition":"definition 2","synonyms":[],"antonyms":[]},{"definition":"definition 3","synonyms":[],"antonyms":[]}]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/senescence"]}]`,
	}

	apiServer := mockDictionaryAPI(mockWords)
	defer apiServer.Close()

	api := dictionary.API{
		URL: apiServer.URL,
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
			expectedOut: "senescence\n[noun] definition 1\n[noun] definition 2\n",
			errExpected: false,
		},
		{
			name:        "ignored definitions limit",
			cmd:         "random",
			initList:    vocab.List{Words: []string{"senescence"}},
			expectedOut: "senescence\n[noun] definition 1\n[noun] definition 2\n[noun] definition 3\n",
			errExpected: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testrandom")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			v, err := newTempVocab(tempDir, test.initList)
			if err != nil {
				t.Fatalf("failed to create initial vocab storage: %v", err)
			}

			var b bytes.Buffer

			cfg := Config{
				Out:   &b,
				Vocab: v,
				Cache: newMemoryCache(nil),
				Dict:  api,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err = cmd.Execute()
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
