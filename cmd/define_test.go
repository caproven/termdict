package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/dictionary/dictionarytest"
)

// TODO these tests gonna be a doozy, take the time to think through how they should be restructured. Ideally
// the printer would be dep injected (not resolved) and we can just stub that out too.

func TestDefineCmd(t *testing.T) {
	dict := dictionarytest.InMemoryDefiner{
		"kappa": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A tortoise-like creature in the Japanese mythology."},
		},
		"cucumber": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "A vine in the gourd family, Cucumis sativus."},
			{PartOfSpeech: "noun", Meaning: "The edible fruit of this plant, having a green rind and crisp white flesh."},
		},
		"senescence": []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "definition 1"},
			{PartOfSpeech: "adjective", Meaning: "definition 2"},
			{PartOfSpeech: "verb", Meaning: "definition 3"},
		},
	}

	cases := []struct {
		name    string
		cmd     string
		list    vocab.List
		wantOut string
		wantErr bool
	}{
		{
			name:    "word found",
			cmd:     "define kappa",
			wantOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
		},
		{
			name:    "word not found",
			cmd:     "define asdf",
			wantErr: true,
		},
		{
			name:    "valid definitions limit",
			cmd:     "define cucumber --limit 1",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n",
		},
		{
			name:    "ignored definitions limit",
			cmd:     "define cucumber --limit 0",
			wantOut: "cucumber\n[noun] A vine in the gourd family, Cucumis sativus.\n[noun] The edible fruit of this plant, having a green rind and crisp white flesh.\n",
		},
		{
			name:    "limit higher than def count",
			cmd:     "define kappa --limit 10",
			wantOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
		},
		{
			name:    "random with empty list",
			cmd:     "define --random",
			list:    vocab.List{Words: []string{}},
			wantErr: true,
		},
		{
			name:    "random with single word",
			cmd:     "define --random",
			list:    vocab.List{Words: []string{"kappa"}},
			wantOut: "kappa\n[noun] A tortoise-like creature in the Japanese mythology.\n",
		},
		{
			name:    "random with definitions limit",
			cmd:     "define --random --limit 2",
			list:    vocab.List{Words: []string{"senescence"}},
			wantOut: "senescence\n[noun] definition 1\n[adjective] definition 2\n",
		},
		{
			name:    "random and specific word",
			cmd:     "define --random something",
			list:    vocab.List{Words: []string{"cucumber"}},
			wantErr: true,
		},
		{
			name:    "invalid output format",
			cmd:     "define senescence --output abcd",
			wantErr: true,
		},
		{
			name: "json output",
			cmd:  "define kappa --output json",
			wantOut: `{
	"Word": "kappa",
	"Definitions": [
		{
			"PartOfSpeech": "noun",
			"Meaning": "A tortoise-like creature in the Japanese mythology."
		}
	]
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			v := newMemoryVocabRepo(test.list)

			var b bytes.Buffer

			cfg := Config{
				Out:   &b,
				Vocab: v,
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

func TestTextPrinter(t *testing.T) {
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
			printer := &textPrinter{limit: test.limit}
			if err := printer.Print(&b, test.word, test.definitions); err != nil {
				t.Errorf("failed to print definition: %v", err)
			}

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}

func TestJsonPrinter(t *testing.T) {
	cases := []struct {
		name        string
		word        string
		definitions []dictionary.Definition
		expected    string
	}{
		{
			name: "single entry",
			word: "guava",
			definitions: []dictionary.Definition{
				{
					PartOfSpeech: "noun", Meaning: "A tropical tree or shrub of the myrtle family",
				},
			},
			expected: `{
	"Word": "guava",
	"Definitions": [
		{
			"PartOfSpeech": "noun",
			"Meaning": "A tropical tree or shrub of the myrtle family"
		}
	]
}
`,
		},
		{
			name: "multiple entries",
			word: "super",
			definitions: []dictionary.Definition{
				{PartOfSpeech: "adjective", Meaning: "Of excellent quality"},
				{PartOfSpeech: "adverb", Meaning: "Very; extremely"},
			},
			expected: `{
	"Word": "super",
	"Definitions": [
		{
			"PartOfSpeech": "adjective",
			"Meaning": "Of excellent quality"
		},
		{
			"PartOfSpeech": "adverb",
			"Meaning": "Very; extremely"
		}
	]
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			printer := new(jsonPrinter)
			if err := printer.Print(&b, test.word, test.definitions); err != nil {
				t.Errorf("failed to print definition: %v", err)
			}

			got := b.String()

			if got != test.expected {
				t.Errorf("got %s, expected %s", got, test.expected)
			}
		})
	}
}
