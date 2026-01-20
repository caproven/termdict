package main

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
)

func main() {
	v, err := storage.NewDefaultVocabRepo()
	if err != nil {
		fmt.Println("Failed to instantiate vocab repo")
		os.Exit(1)
	}

	c, err := dictionary.NewFileCache("")
	if err != nil {
		fmt.Println("Failed to instantiate cache")
		os.Exit(1)
	}
	api := dictionary.NewDefaultWebAPI()
	dict := dictionary.NewCachedDefiner(c, api)

	cfg := &cmd.Config{
		Out:   os.Stdout,
		Vocab: v,
		Dict:  dict,
	}
	if err := cmd.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
