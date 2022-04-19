package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

func mockDictionaryAPI(data map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(dictionary.EndpointPath, func(w http.ResponseWriter, r *http.Request) {
		word := strings.TrimPrefix(r.URL.Path, dictionary.EndpointPath)
		fmt.Fprintf(w, data[word])
	})
	return httptest.NewServer(mux)
}

func newTempStorage(init vocab.List) (vocab.Storage, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return vocab.Storage{}, err
	}
	if err := f.Close(); err != nil {
		return vocab.Storage{}, err
	}

	s := vocab.Storage{
		Path: f.Name(),
	}

	if err := s.Save(init); err != nil {
		return vocab.Storage{}, err
	}

	return s, nil
}
