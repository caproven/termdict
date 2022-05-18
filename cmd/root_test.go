package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
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

func newTempVocab(dir string, init vocab.List) (storage.VocabRepo, error) {
	vocabFile, err := os.CreateTemp(dir, "vocab")
	if err != nil {
		return storage.VocabRepo{}, err
	}

	v := storage.VocabRepo{
		Path: vocabFile.Name(),
	}
	if err := v.Save(init); err != nil {
		return storage.VocabRepo{}, err
	}

	return v, nil
}

func newTempCache(dir string, init dictionary.Cache) (storage.CacheRepo, error) {
	cacheFile, err := os.CreateTemp(dir, "cache")
	if err != nil {
		return storage.CacheRepo{}, err
	}

	s := storage.CacheRepo{
		Path: cacheFile.Name(),
	}
	if err := s.Save(init); err != nil {
		return storage.CacheRepo{}, err
	}

	return s, nil
}
