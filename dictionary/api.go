package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// defaultURL is the default API's URL
const defaultURL string = "https://api.dictionaryapi.dev"

// defaultEndpoint is the API's endpoint for defining words
const defaultEndpoint string = "/api/v2/entries/en/"

// apiResponse is the dictionary API response
type apiResponse struct {
	Meanings []apiMeanings
}

// apiMeanings is a series of definitions broken up
// by a specific part of speech
type apiMeanings struct {
	PartOfSpeech string
	Definitions  []apiDefinition
}

// apiDefinition is a single definition for a word
type apiDefinition struct {
	Definition string
}

// WebAPI lets you interact with a dictionary API
type WebAPI struct {
	url string
	// Word to define should be appended to the end
	endpoint string
}

// NewDefaultWebAPI creates a new instance for connecting to a dictionary API
func NewDefaultWebAPI() WebAPI {
	return WebAPI{
		url:      defaultURL,
		endpoint: defaultEndpoint,
	}
}

func (api WebAPI) Define(word string) ([]Definition, error) {
	apiResp, err := api.query(word)
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

func (api WebAPI) query(w string) (apiResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", api.url, api.endpoint, w))
	if err != nil {
		return apiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse{}, err
	}

	var responses []apiResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return apiResponse{}, err
	}
	if len(responses) == 0 {
		return apiResponse{}, errors.New("didn't find any definitions")
	}
	return responses[0], err
}
