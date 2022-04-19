package dictionary

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDefine(t *testing.T) {
	mockWords := map[string]string{
		"prickly":        `[{"word":"prickly","phonetics":[{"audio":"https://api.dictionaryapi.dev/media/pronunciations/en/prickly.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=77605874","license":{"name":"BY-SA 4.0","url":"https://creativecommons.org/licenses/by-sa/4.0"}}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"Something that gives a pricking sensation; a sharp object.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]},{"partOfSpeech":"adjective","definitions":[{"definition":"Covered with sharp points.","synonyms":[],"antonyms":[],"example":"The prickly pear is a cactus; you have to peel it before eating it to remove the spines and the tough skin."},{"definition":"Easily irritated.","synonyms":[],"antonyms":[],"example":"He has a prickly personality. He doesn't get along with people because he is easily set off."}],"synonyms":["spiny","thorny"],"antonyms":[]},{"partOfSpeech":"adverb","definitions":[{"definition":"In a prickly manner.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/prickly"]}]`,
		"snow":           `[{"word":"snow","phonetic":"/snəʊ/","phonetics":[{"text":"/snəʊ/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/snow-1-uk.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=9027438","license":{"name":"BY 3.0 US","url":"https://creativecommons.org/licenses/by/3.0/us"}},{"text":"/snoʊ/","audio":"https://api.dictionaryapi.dev/media/pronunciations/en/snow-1-us.mp3","sourceUrl":"https://commons.wikimedia.org/w/index.php?curid=1157887","license":{"name":"BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"}}],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"The frozen, crystalline state of water that falls as precipitation.","synonyms":[],"antonyms":[]},{"definition":"A snowfall; a blanket of frozen, crystalline water.","synonyms":[],"antonyms":[],"example":"We have had several heavy snows this year."},{"definition":"A shade of the color white.","synonyms":[],"antonyms":[]}],"synonyms":["blow","shash"],"antonyms":[]},{"partOfSpeech":"verb","definitions":[{"definition":"To have snow fall from the sky.","synonyms":[],"antonyms":[],"example":"It is snowing."}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/snow"]},{"word":"snow","phonetics":[],"meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A square-rigged vessel, differing from a brig only in that she has a trysail mast close abaft the mainmast, on which a large trysail is hoisted.","synonyms":[],"antonyms":[]}],"synonyms":[],"antonyms":[]}],"license":{"name":"CC BY-SA 3.0","url":"https://creativecommons.org/licenses/by-sa/3.0"},"sourceUrls":["https://en.wiktionary.org/wiki/snow"]}]`,
		"no_definitions": `{"title":"No Definitions Found","message":"Sorry pal, we couldn't find definitions for the word you were looking for.","resolution":"You can try the search again at later time or head to the web instead."}`,
		"empty_response": `[]`,
	}

	apiServer := mockDictionaryAPI(mockWords)
	defer apiServer.Close()

	api := API{
		URL: apiServer.URL,
	}

	cases := []struct {
		name        string
		word        string
		entries     []Entry
		errExpected bool
	}{
		{
			name: "standard response",
			word: "prickly",
			entries: []Entry{
				{PartOfSpeech: "noun", Definition: "Something that gives a pricking sensation; a sharp object."},
				{PartOfSpeech: "adjective", Definition: "Covered with sharp points."},
				{PartOfSpeech: "adjective", Definition: "Easily irritated."},
				{PartOfSpeech: "adverb", Definition: "In a prickly manner."},
			},
			errExpected: false,
		},
		{
			name: "multiple definitions in array",
			word: "snow",
			entries: []Entry{
				{PartOfSpeech: "noun", Definition: "The frozen, crystalline state of water that falls as precipitation."},
				{PartOfSpeech: "noun", Definition: "A snowfall; a blanket of frozen, crystalline water."},
				{PartOfSpeech: "noun", Definition: "A shade of the color white."},
				{PartOfSpeech: "verb", Definition: "To have snow fall from the sky."},
			},
			errExpected: false,
		},
		{
			name:        "word with no definition from api",
			word:        "no_definitions",
			entries:     nil,
			errExpected: true,
		},
		{
			name:        "word with empty response from api",
			word:        "empty_response",
			entries:     nil,
			errExpected: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := api.Define(test.word)

			if test.errExpected {
				if err == nil {
					t.Error("expected err but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("didn't expect err but got: %v", err)
				}
			}

			if !reflect.DeepEqual(got, test.entries) {
				t.Errorf("got entries %v, expected %v", got, test.entries)
			}
		})
	}
}

func mockDictionaryAPI(data map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(EndpointPath, func(w http.ResponseWriter, r *http.Request) {
		word := strings.TrimPrefix(r.URL.Path, EndpointPath)
		fmt.Fprintf(w, data[word])
	})
	return httptest.NewServer(mux)
}
