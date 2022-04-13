package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const endpoint string = "https://api.dictionaryapi.dev/api/v2/entries/en"

type APIResponse struct {
	Meanings []APIMeanings
}

type APIMeanings struct {
	PartOfSpeech string
	Definitions  []APIDefinition
}

type APIDefinition struct {
	Definition string
}

type DictionaryAPI struct {
	URL string
}

func NewDictionaryAPI() DictionaryAPI {
	return DictionaryAPI{
		URL: endpoint,
	}
}

func (d DictionaryAPI) Define(word string) ([]Entry, error) {
	apiResp, err := d.query(word)
	if err != nil {
		return nil, fmt.Errorf("failed to define word '%s': %w", word, err)
	}

	entries := []Entry{}

	for _, meaning := range apiResp.Meanings {
		for _, def := range meaning.Definitions {
			entry := Entry{
				PartOfSpeech: meaning.PartOfSpeech,
				Definition:   def.Definition,
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (d DictionaryAPI) query(w string) (APIResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", d.URL, w))
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var responses []APIResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return APIResponse{}, err
	}
	if len(responses) == 0 {
		return APIResponse{}, errors.New("didn't find any definitions")
	}
	return responses[0], err
}
