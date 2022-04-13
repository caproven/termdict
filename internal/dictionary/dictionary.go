package dictionary

type Entry struct {
	PartOfSpeech string
	Definition   string
}

type Definer interface {
	Define(word string) ([]Entry, error)
}

type Definition struct {
	PartOfSpeech string
	Meaning      string
}

func Default() Definer {
	return NewDictionaryAPI()
}
