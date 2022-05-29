package cmd

import (
	"bytes"
	"os"
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
		name    string
		cmd     string
		cache   storage.Cache
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
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer

			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testdefine")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			cfg := Config{
				Out:   &b,
				Vocab: storage.VocabRepo{}, // shouldn't be used by this cmd
				Cache: test.cache,
				Dict:  api,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err = cmd.Execute()
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
