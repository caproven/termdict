package dictionary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	endpoint   string
	httpClient *http.Client
}

// NewDefaultWebAPI creates a new instance for connecting to a dictionary API
func NewDefaultWebAPI() WebAPI {
	return WebAPI{
		url:        defaultURL,
		endpoint:   defaultEndpoint,
		httpClient: http.DefaultClient,
	}
}

func (api WebAPI) Define(ctx context.Context, word string) ([]Definition, error) {
	apiResp, err := api.query(ctx, word)
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

func (api WebAPI) query(ctx context.Context, w string) (apiResponse, error) {
	reqURL := fmt.Sprintf("%s%s%s", api.url, api.endpoint, w)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return apiResponse{}, fmt.Errorf("build request: %w", err)
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return apiResponse{}, err
	}
	defer func(body io.Closer) {
		if err := body.Close(); err != nil {
			slog.Warn("Failed to close http response body", "error", err)
		}
	}(resp.Body)

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
