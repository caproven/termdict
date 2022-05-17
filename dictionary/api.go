package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultURL is the default URL for the dictionary API
const DefaultURL string = "https://api.dictionaryapi.dev"

// EndpointPath is the endpoint of the dictionary API used
// for defining words
const EndpointPath string = "/api/v2/entries/en/"

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

// NewDefaultAPI creates a new instance for connecting to a dictionary API
func NewDefaultAPI() API {
	return API{
		URL: DefaultURL,
	}
}

// Define defines a word
func (d API) Define(word string) ([]Definition, error) {
	apiResp, err := d.query(word)
	if err != nil {
		return nil, fmt.Errorf("failed to define word '%s'", word)
	}

	defs := []Definition{}

	for _, respMeaning := range apiResp.Meanings {
		for _, respDef := range respMeaning.Definitions {
			def := Definition{
				PartOfSpeech: respMeaning.PartOfSpeech,
				Meaning:      respDef.Definition,
			}
			defs = append(defs, def)
		}
	}

	return defs, nil
}

func (d API) query(w string) (APIResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", d.URL, EndpointPath, w))
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
