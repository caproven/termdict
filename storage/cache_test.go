package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
)

func TestCacheRepo_Load(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected dictionary.Cache
	}{
		{
			name:     "empty cache",
			input:    `{}`,
			expected: dictionary.Cache{},
		},
		{
			name:  "single word single definition",
			input: `{"kappa":[{"partOfSpeech":"noun","meaning":"A tortoise-like creature in the Japanese mythology."}]}`,
			expected: dictionary.Cache{
				"kappa": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "A tortoise-like creature in the Japanese mythology.",
					},
				},
			},
		},
		{
			name:  "single word multiple definitions",
			input: `{"sponge":[{"partOfSpeech":"noun","meaning":"A piece of porous material used for washing"},{"partOfSpeech":"verb","meaning":"To clean, soak up, or dab with a sponge"}]}`,
			expected: dictionary.Cache{
				"sponge": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "A piece of porous material used for washing",
					},
					{
						PartOfSpeech: "verb",
						Meaning:      "To clean, soak up, or dab with a sponge",
					},
				},
			},
		},
		{
			name:  "multiple words",
			input: `{"aqueduct":[{"partOfSpeech":"noun","meaning":"An artificial channel that is constructed to convey water"}],"senescence":[{"partOfSpeech":"noun","meaning":"The state or process of ageing"}]}`,
			expected: dictionary.Cache{
				"aqueduct": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "An artificial channel that is constructed to convey water",
					},
				},
				"senescence": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "The state or process of ageing",
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fName, err := newFileWithData(test.input)
			if err != nil {
				t.Errorf("failed to create temp storage: %v", err)
			}
			defer os.Remove(fName)

			c := CacheRepo{
				Path: fName,
			}

			got, err := c.Load()
			if err != nil {
				t.Errorf("failed to load storage: %v", err)
			}

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("caches did not match; got %v, expected %v", got, test.expected)
			}
		})
	}

	t.Run("file doesn't exist", func(t *testing.T) {
		c := CacheRepo{
			Path: filepath.Join(os.TempDir(), "thisfileshouldntexist"),
		}

		got, err := c.Load()
		if err != nil {
			t.Errorf("failed to load storage: %v", err)
		}
		expect := dictionary.Cache{}

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("caches did not match; got %v, expected %v", got, expect)
		}
	})
}

func TestCacheRepo_Save(t *testing.T) {
	cases := []struct {
		name     string
		input    dictionary.Cache
		expected string
	}{
		{
			name:     "empty cache",
			input:    dictionary.Cache{},
			expected: `{}`,
		},
		{
			name: "single word single definition",
			input: dictionary.Cache{
				"kappa": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "A tortoise-like creature in the Japanese mythology.",
					},
				},
			},
			expected: `{"kappa":[{"PartOfSpeech":"noun","Meaning":"A tortoise-like creature in the Japanese mythology."}]}`,
		},
		{
			name: "single word multiple definitions",
			input: dictionary.Cache{
				"sponge": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "A piece of porous material used for washing",
					},
					{
						PartOfSpeech: "verb",
						Meaning:      "To clean, soak up, or dab with a sponge",
					},
				},
			},
			expected: `{"sponge":[{"PartOfSpeech":"noun","Meaning":"A piece of porous material used for washing"},{"PartOfSpeech":"verb","Meaning":"To clean, soak up, or dab with a sponge"}]}`,
		},
		{
			name: "multiple words",
			input: dictionary.Cache{
				"aqueduct": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "An artificial channel that is constructed to convey water",
					},
				},
				"senescence": []dictionary.Definition{
					{
						PartOfSpeech: "noun",
						Meaning:      "The state or process of ageing",
					},
				},
			},
			expected: `{"aqueduct":[{"PartOfSpeech":"noun","Meaning":"An artificial channel that is constructed to convey water"}],"senescence":[{"PartOfSpeech":"noun","Meaning":"The state or process of ageing"}]}`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fName, err := newFile()
			if err != nil {
				t.Errorf("failed to create temp file: %v", err)
			}
			defer os.Remove(fName)

			c := CacheRepo{
				Path: fName,
			}

			err = c.Save(test.input)
			if err != nil {
				t.Errorf("failed to save to storage: %v", err)
			}

			got, err := os.ReadFile(c.Path)
			if err != nil {
				t.Errorf("failed to read from storage: %v", err)
			}

			if string(got) != test.expected {
				t.Errorf("got %v, expected %v", string(got), test.expected)
			}
		})
	}
}
