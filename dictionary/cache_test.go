package dictionary_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/dictionary/dictionarytest"
)

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

				found, _ := tt.fields.cache.ContainsWord(t.Context(), tt.args.word)
				if !found {
					t.Errorf("cache did not contain defined word %s", tt.args.word)
				}
				lookup, _ := tt.fields.cache.LookupWord(t.Context(), tt.args.word)
				if !reflect.DeepEqual(lookup, tt.want) {
					t.Errorf("cached content = %v, want %v", lookup, tt.want)
				}
			}
		})
	}
}

type memoryCache map[string][]dictionary.Definition

func (mc memoryCache) ContainsWord(_ context.Context, word string) (bool, error) {
	_, ok := mc[word]
	return ok, nil
}

func (mc memoryCache) SaveWord(_ context.Context, word string, defs []dictionary.Definition) error {
	mc[word] = defs
	return nil
}

func (mc memoryCache) LookupWord(_ context.Context, word string) ([]dictionary.Definition, error) {
	defs, ok := mc[word]
	if !ok {
		return nil, fmt.Errorf("word %s not found in cache", word)
	}
	return defs, nil
}
