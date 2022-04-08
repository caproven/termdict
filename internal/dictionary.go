package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const dictionaryAPIEndpoint string = "https://api.dictionaryapi.dev/api/v2/entries/en"

type Definition struct {
	PartOfSpeech string
	Meaning      string
}

type dictAPIResponse struct {
	Meanings []dictAPIMeanings
}

type dictAPIMeanings struct {
	PartOfSpeech string
	Definitions  []dictAPIDefinitions
}

type dictAPIDefinitions struct {
	Definition string
}

func Define(w string) ([]Definition, error) {
	apiResp, err := queryWord(w)
	if err != nil {
		return nil, err
	}

	definitions := []Definition{}

	for _, meaning := range apiResp.Meanings {
		for _, def := range meaning.Definitions {
			definition := Definition{
				PartOfSpeech: meaning.PartOfSpeech,
				Meaning:      def.Definition,
			}
			definitions = append(definitions, definition)
		}
	}

	return definitions, nil
}

func queryWord(w string) (dictAPIResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", dictionaryAPIEndpoint, w))
	if err != nil {
		return dictAPIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dictAPIResponse{}, err
	}

	var responses []dictAPIResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return dictAPIResponse{}, err
	}
	if len(responses) == 0 {
		return dictAPIResponse{}, errors.New("didn't find any definitions")
	}
	return responses[0], err
}
