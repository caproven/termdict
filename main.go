package main

import (
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
)

func main() {
	if err := storage.CreateConfigDir(); err != nil {
		panic(err)
	}

	v := storage.NewDefaultVocabRepo()

	c := storage.NewDefaultFSCache()

	api := dictionary.NewDefaultWebAPI()

	cfg := &cmd.Config{
		Out:   os.Stdout,
		Vocab: v,
		Cache: c,
		Dict:  api,
	}
	if err := cmd.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
