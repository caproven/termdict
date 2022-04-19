package vocab

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected List
	}{
		{
			name:     "empty list",
			input:    `[]`,
			expected: List{Words: []string{}},
		},
		{
			name:     "single word",
			input:    `["kappa"]`,
			expected: List{Words: []string{"kappa"}},
		},
		{
			name:     "multiple words",
			input:    `["kappa", "cucumber", "terminal", "dictionary"]`,
			expected: List{Words: []string{"kappa", "cucumber", "terminal", "dictionary"}},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fName, err := newFileWithData(test.input)
			if err != nil {
				t.Errorf("failed to create temp storage: %v", err)
			}
			defer os.Remove(fName)

			s := Storage{
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
		s := Storage{
			Path: filepath.Join(os.TempDir(), "thisfileshouldntexist"),
		}

		got, err := s.Load()
		if err != nil {
			t.Errorf("failed to load storage: %v", err)
		}
		expect := List{Words: []string{}}

		assertLists(t, got, expect)
	})
}

func TestSave(t *testing.T) {
	cases := []struct {
		name     string
		input    List
		expected string
	}{
		{
			name:     "empty list",
			input:    List{Words: []string{}},
			expected: `[]`,
		},
		{
			name:     "single word",
			input:    List{Words: []string{"kappa"}},
			expected: `["kappa"]`,
		},
		{
			name:     "multiple words",
			input:    List{Words: []string{"kappa", "cucumber", "terminal", "dictionary"}},
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

			s := Storage{
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
