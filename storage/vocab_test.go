package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caproven/termdict/vocab"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected vocab.List
	}{
		{
			name:     "empty list",
			input:    `[]`,
			expected: vocab.List{Words: []string{}},
		},
		{
			name:     "single word",
			input:    `["kappa"]`,
			expected: vocab.List{Words: []string{"kappa"}},
		},
		{
			name:     "multiple words",
			input:    `["kappa", "cucumber", "terminal", "dictionary"]`,
			expected: vocab.List{Words: []string{"kappa", "cucumber", "terminal", "dictionary"}},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fName, err := newFileWithData(test.input)
			if err != nil {
				t.Errorf("failed to create temp storage: %v", err)
			}
			defer os.Remove(fName)

			s := VocabRepo{
				Path: fName,
			}

			got, err := s.Load()
			if err != nil {
				t.Errorf("failed to load storage: %v", err)
			}

			assertLists(t, got, test.expected)
		})
	}

	t.Run("file doesn't exist", func(t *testing.T) {
		s := VocabRepo{
			Path: filepath.Join(os.TempDir(), "thisfileshouldntexist"),
		}

		got, err := s.Load()
		if err != nil {
			t.Errorf("failed to load storage: %v", err)
		}
		expect := vocab.List{Words: []string{}}

		assertLists(t, got, expect)
	})
}

func TestSave(t *testing.T) {
	cases := []struct {
		name     string
		input    vocab.List
		expected string
	}{
		{
			name:     "empty list",
			input:    vocab.List{Words: []string{}},
			expected: `[]`,
		},
		{
			name:     "single word",
			input:    vocab.List{Words: []string{"kappa"}},
			expected: `["kappa"]`,
		},
		{
			name:     "multiple words",
			input:    vocab.List{Words: []string{"kappa", "cucumber", "terminal", "dictionary"}},
			expected: `["kappa","cucumber","terminal","dictionary"]`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fName, err := newFile()
			if err != nil {
				t.Errorf("failed to create temp file: %v", err)
			}
			defer os.Remove(fName)

			s := VocabRepo{
				Path: fName,
			}

			err = s.Save(test.input)
			if err != nil {
				t.Errorf("failed to save to storage: %v", err)
			}

			got, err := os.ReadFile(s.Path)
			if err != nil {
				t.Errorf("failed to read from storage: %v", err)
			}

			if string(got) != test.expected {
				t.Errorf("got %v, expected %v", string(got), test.expected)
			}
		})
	}
}

func newFileWithData(data string) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(data)
	f.Close()
	return f.Name(), err
}

func newFile() (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), err
}

func assertLists(t testing.TB, got, expected vocab.List) {
	t.Helper()

	if len(got.Words) != len(expected.Words) {
		t.Fatalf("lists not the same length; got %v, expected %v", got, expected)
	}

	for i, v := range got.Words {
		if v != expected.Words[i] {
			t.Errorf("lists did not match; got %v, expected %v", got, expected)
		}
	}
}
