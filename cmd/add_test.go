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
	mockWords := map[string]string{
		"kappa":      `[{"word":"kappa","phonetic":"/ˈkæpə/","phonetics":[{"text":"/ˈkæpə/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/kappa-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=311306"}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A tortoise-like creature in the Japanese mythology.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/kappa"]}]`,
		"cucumber":   `[{"word":"cucumber","phonetic":"/ˈkjuːˌkʌmbə/","phonetics":[{"text":"/ˈkjuːˌkʌmbə/","audio":""},{"text":"/ˈkjuːˌkʌmbɚ/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/cucumber-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=1769363","license":{"name":"BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"}}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A vine in the gourd family, Cucumis sativus.","synonyms":[],"antonyms":[]},{"definition":"The edible fruit of this plant, having a green rind and crisp white flesh.","synonyms":[],"antonyms":[]}],"synonyms":["cuke"],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/cucumber"]}]`,
		"terminal":   `[{"word":"terminal","phonetic":"/ˈtɚmɪnəl/","phonetics":[{"text":"/ˈtɚmɪnəl/","audio":""}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A building in an airport where passengers transfer from ground transportation to the facilities that allow them to board airplanes.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/terminal"]}]`,
		"dictionary": `[{"word":"dictionary","phonetic":"/ˈdɪkʃəˌnɛɹi/","phonetics":[{"text":"/ˈdɪkʃəˌnɛɹi/","audio":""},{"text":"/ˈdɪkʃ(ə)n(ə)ɹi/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/dictionary-uk.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=503422"},{"text":"/ˈdɪkʃəˌnɛɹi/","audio":""}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A reference work with a list of words from one or more languages, normally ordered alphabetically, explaining each word's meaning, and sometimes containing information on its etymology, pronunciation, usage, translations, and other data.","synonyms":["wordbook"],"antonyms":[]}],"synonyms":["wordbook"],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/dictionary"]}]`,
	}

	apiServer := mockDictionaryAPI(mockWords)
	defer apiServer.Close()

	api := dictionary.API{
		URL: apiServer.URL,
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
			cmd:          "add kappa",
			initList:     vocab.List{Words: []string{}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  false,
		},
		{
			name:         "multiple words to empty",
			cmd:          "add terminal dictionary",
			initList:     vocab.List{Words: []string{}},
			expectedList: vocab.List{Words: []string{"terminal", "dictionary"}},
			errExpected:  false,
		},
		{
			name:         "new word to existing",
			cmd:          "add cucumber",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa", "cucumber"}},
			errExpected:  false,
		},
		{
			name:         "multiple new words to existing",
			cmd:          "add dictionary terminal",
			initList:     vocab.List{Words: []string{"kappa", "cucumber"}},
			expectedList: vocab.List{Words: []string{"kappa", "cucumber", "dictionary", "terminal"}},
			errExpected:  false,
		},
		{
			name:         "duplicate",
			cmd:          "add terminal",
			initList:     vocab.List{Words: []string{"cucumber", "terminal"}},
			expectedList: vocab.List{Words: []string{"cucumber", "terminal"}},
			errExpected:  true,
		},
		{
			name:         "multiple words with duplicate",
			cmd:          "add cucumber dictionary",
			initList:     vocab.List{Words: []string{"dictionary", "kappa"}},
			expectedList: vocab.List{Words: []string{"dictionary", "kappa"}},
			errExpected:  true,
		},
		{
			name:         "unknown word",
			cmd:          "add asdf",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  true,
		},
		{
			name:         "multiple words with unknown",
			cmd:          "add cucumber asdf terminal",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa"}},
			errExpected:  true,
		},
		{
			name:         "unknown word with check disabled",
			cmd:          "add asdf --no-check",
			initList:     vocab.List{Words: []string{"kappa"}},
			expectedList: vocab.List{Words: []string{"kappa", "asdf"}},
			errExpected:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testadd")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			v, err := newTempVocab(tempDir, test.initList)
			if err != nil {
				t.Fatalf("failed to create initial vocab storage: %v", err)
			}

			cfg := Config{
				Out:   os.Stdout,
				Vocab: v,
				Dict:  api,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err = cmd.Execute()

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
