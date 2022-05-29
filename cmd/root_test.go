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

var _ storage.Cache = memoryCache{}

type memoryCache struct {
	data map[string][]dictionary.Definition
}

func newMemoryCache(data map[string][]dictionary.Definition) memoryCache {
	if data == nil {
		data = make(map[string][]dictionary.Definition)
	}
	return memoryCache{data: data}
}

func (mc memoryCache) Contains(word string) (bool, error) {
	_, ok := mc.data[word]
	return ok, nil
}

func (mc memoryCache) Save(word string, defs []dictionary.Definition) error {
	mc.data[word] = defs
	return nil
}

func (mc memoryCache) Lookup(word string) ([]dictionary.Definition, error) {
	defs, ok := mc.data[word]
	if !ok {
		return nil, fmt.Errorf("word %s not found in cache", word)
	}
	return defs, nil
}
