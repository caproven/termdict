package dictionary_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/dictionary/dictionarytest"
)

func TestFileCache_Contains(t *testing.T) {
	tests := []struct {
		name    string
		cache   map[string][]dictionary.Definition
		word    string
		want    bool
		wantErr bool
	}{
		{
			name:    "empty cache",
			cache:   map[string][]dictionary.Definition{},
			word:    "kappa",
			want:    false,
			wantErr: false,
		},
		{
			name: "find in single word cache",
			cache: map[string][]dictionary.Definition{
				"kappa": {{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"}},
			},
			word:    "kappa",
			want:    true,
			wantErr: false,
		},
		{
			name: "find in multi-word cache",
			cache: map[string][]dictionary.Definition{
				"kappa":    {{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"}},
				"cucumber": {{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus"}},
			},
			word:    "cucumber",
			want:    true,
			wantErr: false,
		},
		{
			name: "not in cache",
			cache: map[string][]dictionary.Definition{
				"kappa":    {{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"}},
				"cucumber": {{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus"}},
			},
			word:    "dictionary",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testcachecontains")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			fc, err := dictionary.NewFileCache(tempDir)
			if err != nil {
				t.Fatalf("failed to create cache: %v", err)
			}

			if err := writeCache(t.Context(), tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			got, err := fc.Contains(t.Context(), tt.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileCache.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileCache_Save(t *testing.T) {
	tests := []struct {
		name    string
		cache   map[string][]dictionary.Definition
		word    string
		defs    []dictionary.Definition
		want    string
		wantErr bool
	}{
		{
			name:  "single word single definition",
			cache: map[string][]dictionary.Definition{},
			word:  "kappa",
			defs: []dictionary.Definition{
				{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"},
			},
			want:    `[{"PartOfSpeech":"noun","Meaning":"A tortoise-like creature in the Japanese mythology"}]`,
			wantErr: false,
		},
		{
			name:  "single word multiple definitions",
			cache: map[string][]dictionary.Definition{},
			word:  "sponge",
			defs: []dictionary.Definition{
				{PartOfSpeech: "noun", Meaning: "A piece of porous material used for washing"},
				{PartOfSpeech: "verb", Meaning: "To clean, soak up, or dab with a sponge"},
			},
			want:    `[{"PartOfSpeech":"noun","Meaning":"A piece of porous material used for washing"},{"PartOfSpeech":"verb","Meaning":"To clean, soak up, or dab with a sponge"}]`,
			wantErr: false,
		},
		{
			name: "existing word overwritten",
			cache: map[string][]dictionary.Definition{
				"senescence": {{PartOfSpeech: "noun", Meaning: "The state or process of ageing"}},
			},
			word: "senescence",
			defs: []dictionary.Definition{
				{PartOfSpeech: "noun", Meaning: "a new definition"},
			},
			want: `[{"PartOfSpeech":"noun","Meaning":"a new definition"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testcachesave")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			fc, err := dictionary.NewFileCache(tempDir)
			if err != nil {
				t.Fatalf("failed to create cache: %v", err)
			}

			if err := writeCache(t.Context(), tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			if err := fc.Save(t.Context(), tt.word, tt.defs); (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := os.ReadFile(fc.PathFor(tt.word))
			if err != nil {
				t.Fatalf("failed to read expected cache file: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("got %v, expected %v", string(got), tt.want)
			}
		})
	}
}

func TestFileCache_Lookup(t *testing.T) {
	tests := []struct {
		name    string
		cache   map[string][]dictionary.Definition
		word    string
		want    []dictionary.Definition
		wantErr bool
	}{
		{
			name:    "empty cache",
			cache:   map[string][]dictionary.Definition{},
			word:    "kappa",
			want:    nil,
			wantErr: true,
		},
		{
			name: "word not in cache",
			cache: map[string][]dictionary.Definition{
				"kappa":    {{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"}},
				"cucumber": {{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus"}},
			},
			word:    "terminal",
			want:    nil,
			wantErr: true,
		},
		{
			name: "word in cache",
			cache: map[string][]dictionary.Definition{
				"kappa":    {{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology"}},
				"cucumber": {{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus"}},
			},
			word: "cucumber",
			want: []dictionary.Definition{
				{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testcachelookup")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			fc, err := dictionary.NewFileCache(tempDir)
			if err != nil {
				t.Fatalf("failed to create cache: %v", err)
			}

			if err := writeCache(t.Context(), tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			got, err := fc.Lookup(t.Context(), tt.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileCache.Lookup() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("invalid file contents", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "termdict-testcachelookup")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		fc, err := dictionary.NewFileCache(tempDir)
		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		err = os.WriteFile(fc.PathFor("mouse"), []byte("invalid_data"), os.ModePerm)
		if err != nil {
			t.Errorf("failed to write temp file: %v", err)
		}

		_, err = fc.Lookup(t.Context(), "mouse")
		if err == nil {
			t.Errorf("FileCache.Lookup() error = %v, wantErr %v", err, true)
		}
	})
}

func writeCache(ctx context.Context, cache map[string][]dictionary.Definition, c dictionary.Cache) error {
	for word, defs := range cache {
		if err := c.Save(ctx, word, defs); err != nil {
			return err
		}
	}
	return nil
}

func TestCachedDefiner_Define(t *testing.T) {
	type fields struct {
		cache    dictionary.Cache
		fallback dictionary.Definer
	}
	type args struct {
		word string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []dictionary.Definition
		wantErr bool
	}{
		{
			name: "word defined and not in cache",
			fields: fields{
				cache: make(memoryCache),
				fallback: dictionarytest.InMemoryDefiner{
					"splash": {
						{PartOfSpeech: "noun", Meaning: "The sound made by an object hitting a liquid"},
						{PartOfSpeech: "verb", Meaning: "To hit or agitate liquid"},
					},
				},
			},
			args: args{word: "splash"},
			want: []dictionary.Definition{
				{PartOfSpeech: "noun", Meaning: "The sound made by an object hitting a liquid"},
				{PartOfSpeech: "verb", Meaning: "To hit or agitate liquid"},
			},
			wantErr: false,
		},
		{
			name: "word not defined and in cache",
			fields: fields{
				cache: memoryCache{
					"photosynthesis": {{PartOfSpeech: "noun", Meaning: "Any process by which plants and other photoautotrophs convert light energy into chemical energy"}},
				},
				fallback: make(dictionarytest.InMemoryDefiner),
			},
			args:    args{word: "photosynthesis"},
			want:    []dictionary.Definition{{PartOfSpeech: "noun", Meaning: "Any process by which plants and other photoautotrophs convert light energy into chemical energy"}},
			wantErr: false,
		},
		{
			name: "word defined and in cache",
			fields: fields{
				cache: memoryCache{
					"aardvark": {{Meaning: "cached definition"}},
				},
				fallback: dictionarytest.InMemoryDefiner{
					"aardvark": {{Meaning: "fallback definition"}},
				},
			},
			args: args{word: "aardvark"},
			want: []dictionary.Definition{{Meaning: "cached definition"}},
		},
		{
			name: "word not defined and not in cache",
			fields: fields{
				cache:    make(memoryCache),
				fallback: make(dictionarytest.InMemoryDefiner),
			},
			args:    args{word: "platypus"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dictionary.NewCachedDefiner(tt.fields.cache, tt.fields.fallback)
			got, err := d.Define(t.Context(), tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("CachedDefiner.Define() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CachedDefiner.Define() = %v, want %v", got, tt.want)
			}

			if !tt.wantErr {
				// verify word was cached

				found, _ := tt.fields.cache.Contains(t.Context(), tt.args.word)
				if !found {
					t.Errorf("cache did not contain defined word %s", tt.args.word)
				}
				lookup, _ := tt.fields.cache.Lookup(t.Context(), tt.args.word)
				if !reflect.DeepEqual(lookup, tt.want) {
					t.Errorf("cached content = %v, want %v", lookup, tt.want)
				}
			}
		})
	}
}

type memoryCache map[string][]dictionary.Definition

func (mc memoryCache) Contains(_ context.Context, word string) (bool, error) {
	_, ok := mc[word]
	return ok, nil
}

func (mc memoryCache) Save(_ context.Context, word string, defs []dictionary.Definition) error {
	mc[word] = defs
	return nil
}

func (mc memoryCache) Lookup(_ context.Context, word string) ([]dictionary.Definition, error) {
	defs, ok := mc[word]
	if !ok {
		return nil, fmt.Errorf("word %s not found in cache", word)
	}
	return defs, nil
}
