package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const endpoint string = "https://api.dictionaryapi.dev/api/v2/entries/en"

// APIResponse is the dictionary API response
type APIResponse struct {
	Meanings []APIMeanings
}

// APIMeanings is a series of definitions broken up
// by a specific part of speech
type APIMeanings struct {
	PartOfSpeech string
	Definitions  []APIDefinition
}

// APIDefinition is a single definition for a word
type APIDefinition struct {
	Definition string
}

// API lets you interact with a dictionary API
type API struct {
	URL string
}

// NewAPI creates a new instance for connecting to a dictionary API
func NewAPI() API {
	return API{
		URL: endpoint,
	}
}

// Define defines a word
func (d API) Define(word string) ([]Entry, error) {
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

func (d API) query(w string) (APIResponse, error) {
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
