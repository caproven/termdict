package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
)

func TestDefineCmd(t *testing.T) {
	dict := memoryDefiner{
		"kappa": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology."},
		},
		"cucumber": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus."},
			{PartOfSpeech: "noun", Meaning: "The edible fruit of this plant, having a green rind and crisp white flesh."},
		},
	}

	cases := []struct {
		name    string
		cmd     string
		word    string
		wantOut string
		wantErr bool
	}{
		{
			name:    "word found",
			cmd:     "define kappa",
			word:    "kappa",
			wantOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
		},
		{
			name:    "word not found",
			cmd:     "define asdf",
			word:    "asdf",
			wantErr: true,
		},
		{
			name:    "valid definitions limit",
			cmd:     "define cucumber --limit 1",
			word:    "cucumber",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n",
		},
		{
			name:    "ignored definitions limit",
			cmd:     "define cucumber --limit 0",
			word:    "cucumber",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n[noun] The edible fruit of this plant, having a green rind and crisp white flesh.\n",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer

			cfg := Config{
				Out:   &b,
				Vocab: nil, // shouldn't be used by this cmd
				Dict:  dict,
			}

			cmd := NewRootCmd(&cfg)
			cmd.SetArgs(strings.Split(test.cmd, " "))

			err := cmd.Execute()
			gotOut := b.String()

			if err != nil {
				if (err != nil) != test.wantErr {
					t.Errorf("define cmd error = %v, wantErr %v", err, test.wantErr)
				}
				return
			}

			if gotOut != test.wantOut {
				t.Errorf("got %v, expected %v", gotOut, test.wantOut)
			}
		})
	}
}

func TestPrintDefinition(t *testing.T) {
	cases := []struct {
		name        string
		word        string
		limit       int
		definitions []dictionary.Definition
		expected    string
	}{
		{
			name: "single entry",
			word: "sponge",
			definitions: []dictionary.Definition{
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
			name: "multiple entries no limit",
			word: "sponge",
			definitions: []dictionary.Definition{
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
		{
			name:  "negative limit ignored",
			word:  "sponge",
			limit: -5,
			definitions: []dictionary.Definition{
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
		{
			name:  "positive limit obeyed",
			word:  "sponge",
			limit: 1,
			definitions: []dictionary.Definition{
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
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			printDefinition(&b, test.word, test.definitions, test.limit)

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}
