package vocab

import "testing"

func TestAddWord(t *testing.T) {
	cases := []struct {
		name      string
		list      List
		word      string
		expected  List
		expectErr bool
	}{
		{
			name:      "empty list",
			list:      List{},
			word:      "aardvark",
			expected:  List{Words: []string{"aardvark"}},
			expectErr: false,
		},
		{
			name:      "duplicate",
			list:      List{Words: []string{"apple", "aardvark"}},
			word:      "aardvark",
			expected:  List{Words: []string{"apple", "aardvark"}},
			expectErr: true,
		},
		{
			name:      "not duplicate",
			list:      List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			word:      "aardvark",
			expected:  List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy", "aardvark"}},
			expectErr: false,
		},
		{
			name:      "duplicate capitalized",
			list:      List{Words: []string{"apple", "aardvark"}},
			word:      "AARDVARK",
			expected:  List{Words: []string{"apple", "aardvark"}},
			expectErr: true,
		},
		{
			name:      "not duplicate capitalized",
			list:      List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			word:      "AARDVARK",
			expected:  List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy", "aardvark"}},
			expectErr: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := test.list.AddWord(test.word)

			assertLists(t, test.list, test.expected)
			gotErr := err != nil
			if gotErr != test.expectErr {
				t.Errorf("got error: %t, expected error: %t", err == nil, test.expectErr)
			}
		})
	}
}

func TestRemoveWord(t *testing.T) {
	cases := []struct {
		name      string
		list      List
		word      string
		expected  List
		expectErr bool
	}{
		{
			name:      "empty list",
			list:      List{},
			word:      "aardvark",
			expected:  List{},
			expectErr: true,
		},
		{
			name:      "word exists",
			list:      List{Words: []string{"apple", "aardvark"}},
			word:      "aardvark",
			expected:  List{Words: []string{"apple"}},
			expectErr: false,
		},
		{
			name:      "word doesn't exist",
			list:      List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			word:      "aardvark",
			expected:  List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			expectErr: true,
		},
		{
			name:      "capitalized word exists",
			list:      List{Words: []string{"apple", "aardvark"}},
			word:      "AARDVARK",
			expected:  List{Words: []string{"apple"}},
			expectErr: false,
		},
		{
			name:      "capitalized word doesn't exist",
			list:      List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			word:      "AARDVARK",
			expected:  List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			expectErr: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := test.list.RemoveWord(test.word)

			assertLists(t, test.list, test.expected)
			gotErr := err != nil
			if gotErr != test.expectErr {
				t.Errorf("got error: %t, expected error: %t", err == nil, test.expectErr)
			}
		})
	}
}

func TestWordExists(t *testing.T) {
	cases := []struct {
		name     string
		list     List
		word     string
		expected bool
	}{
		{
			name:     "empty list",
			list:     List{},
			word:     "aardvark",
			expected: false,
		},
		{
			name:     "word doesn't exist",
			list:     List{Words: []string{"apple", "banana", "capricorn", "delta", "entropy"}},
			word:     "aardvark",
			expected: false,
		},
		{
			name:     "word exists at first idx",
			list:     List{Words: []string{"aardvark", "banana", "capricorn", "delta", "entropy"}},
			word:     "aardvark",
			expected: true,
		},
		{
			name:     "word exists at middle idx",
			list:     List{Words: []string{"apple", "banana", "aardvark", "delta", "entropy"}},
			word:     "aardvark",
			expected: true,
		},
		{
			name:     "word exists at last idx",
			list:     List{Words: []string{"apple", "banana", "capricorn", "delta", "aardvark"}},
			word:     "aardvark",
			expected: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := test.list.WordExists(test.word)

			if got != test.expected {
				t.Errorf("got %t, expected %t", got, test.expected)
			}
		})
	}
}

func assertLists(t testing.TB, got, expected List) {
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
