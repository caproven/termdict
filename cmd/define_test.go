package cmd

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
)

func TestDefineCmd(t *testing.T) {
	mockWords := map[string]string{
		"kappa":    `[{"word":"kappa","phonetic":"/ˈkæpə/","phonetics":[{"text":"/ˈkæpə/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/kappa-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=311306"}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A tortoise-like creature in the Japanese mythology.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/kappa"]}]`,
		"cucumber": `[{"word":"cucumber","phonetic":"/ˈkjuːˌkʌmbə/","phonetics":[{"text":"/ˈkjuːˌkʌmbə/","audio":""},{"text":"/ˈkjuːˌkʌmbɚ/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/cucumber-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=1769363","license":{"name":"BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"}}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A vine in the gourd family, Cucumis sativus.","synonyms":[],"antonyms":[]},{"definition":"The edible fruit of this plant, having a green rind and crisp white flesh.","synonyms":[],"antonyms":[]}],"synonyms":["cuke"],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/cucumber"]}]`,
	}

	apiServer := mockDictionaryAPI(mockWords)
	defer apiServer.Close()

	api := dictionary.API{
		URL: apiServer.URL,
	}

	cases := []struct {
		name          string
		cmd           string
		initCache     dictionary.Cache
		expectedCache dictionary.Cache
		expectedOut   string
		errExpected   bool
	}{
		{
			name:      "word not in cache",
			cmd:       "define kappa",
			initCache: dictionary.Cache{},
			expectedCache: dictionary.Cache{
				"kappa": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "A tortoise-like creature in the Japanese mythology.",
					},
				},
			},
			expectedOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
			errExpected: false,
		},
		{
			name:          "word not found",
			cmd:           "define asdf",
			initCache:     dictionary.Cache{},
			expectedCache: dictionary.Cache{},
			errExpected:   true,
		},
		{
			name: "word in cache",
			cmd:  "define cucumber",
			initCache: dictionary.Cache{
				"cucumber": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						// different from mock word definition(s), test that cache value
						// is used and not modified
						Meaning: "Custom meaning",
					},
				},
			},
			expectedCache: dictionary.Cache{
				"cucumber": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "Custom meaning",
					},
				},
			},
			expectedOut: "cucumber\n[noun] Custom meaning\n",
			errExpected: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer

			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testdefine")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			c, err := newTempCache(tempDir, test.initCache)
			if err != nil {
				t.Fatalf("failed to create initial vocab storage: %v", err)
			}

			cfg := Config{
				Out:   &b,
				Vocab: storage.VocabRepo{}, // shouldn't be used by this cmd
				Cache: c,
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
				return
			}

			if out != test.expectedOut {
				t.Errorf("got %v, expected %v", out, test.expectedOut)
			}

			if err != nil {
				t.Errorf("didn't expect err but got: %v", err)
			}

			gotCache, err := c.Load()
			if err != nil {
				t.Errorf("failed to load storage after executing command: %v", err)
			}
			if !reflect.DeepEqual(gotCache, test.expectedCache) {
				t.Errorf("got %v, expected %v", gotCache, test.expectedCache)
			}
		})
	}
}
