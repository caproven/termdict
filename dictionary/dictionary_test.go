package dictionary

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEntryStringer(t *testing.T) {
	var _ fmt.Stringer = Definition{}
}

func TestPrintDefinition(t *testing.T) {
	cases := []struct {
		name        string
		word        string
		definitions []Definition
		expected    string
	}{
		{
			name: "single entry",
			word: "sponge",
			definitions: []Definition{
				{
					PartOfSpeech: "noun",
					Meaning:      "A piece of porous material used for washing",
				},
			},
			expected: `sponge
[noun] A piece of porous material used for washing
`,
		},
		{
			name: "multiple entries",
			word: "sponge",
			definitions: []Definition{
				{
					PartOfSpeech: "noun",
					Meaning:      "A piece of porous material used for washing",
				},
				{
					PartOfSpeech: "verb",
					Meaning:      "To clean, soak up, or dab with a sponge",
				},
			},
			expected: `sponge
[noun] A piece of porous material used for washing
[verb] To clean, soak up, or dab with a sponge
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			PrintDefinition(&b, test.word, test.definitions)

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}
