package dictionarytest

import (
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
)

func TestInMemoryDefiner_Define(t *testing.T) {
	tests := []struct {
		name    string
		m       InMemoryDefiner
		word    string
		want    []dictionary.Definition
		wantErr bool
	}{
		{
			name: "gives definition",
			m: InMemoryDefiner{
				"exacerbate": []dictionary.Definition{
					{
						PartOfSpeech: "verb",
						Meaning:      "To make worse",
					},
				},
			},
			word: "exacerbate",
			want: []dictionary.Definition{
				{
					PartOfSpeech: "verb",
					Meaning:      "To make worse",
				},
			},
			wantErr: false,
		},
		{
			name:    "fails for unknown word",
			m:       InMemoryDefiner{},
			word:    "nonchalant",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Define(t.Context(), tt.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Define() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Define() got = %v, want %v", got, tt.want)
			}
		})
	}
}
