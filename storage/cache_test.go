package storage

import (
	"os"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
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

			fc := FileCache{
				DirPath: tempDir,
			}

			if err := writeCache(tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			got, err := fc.Contains(tt.word)
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

			fc := FileCache{
				DirPath: tempDir,
			}

			if err := writeCache(tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			if err := fc.Save(tt.word, tt.defs); (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := os.ReadFile(fc.fileForWord(tt.word))
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

			fc := FileCache{
				DirPath: tempDir,
			}

			if err := writeCache(tt.cache, fc); err != nil {
				t.Fatalf("failed to write initial cache files: %v", err)
			}

			got, err := fc.Lookup(tt.word)
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

		fc := FileCache{
			DirPath: tempDir,
		}

		err = os.WriteFile(fc.fileForWord("mouse"), []byte("invalid_data"), os.ModePerm)
		if err != nil {
			t.Errorf("failed to write temp file: %v", err)
		}

		_, err = fc.Lookup("mouse")
		if err == nil {
			t.Errorf("FileCache.Lookup() error = %v, wantErr %v", err, true)
		}
	})
}

func writeCache(cache map[string][]dictionary.Definition, c Cache) error {
	for word, defs := range cache {
		if err := c.Save(word, defs); err != nil {
			return err
		}
	}
	return nil
}
