package main

import (
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
)

func main() {
	if err := vocab.CreateConfigDir(); err != nil {
		panic(err)
	}

	s := vocab.Storage{
		Path: vocab.DefaultFilepath(),
	}

	api := dictionary.NewDefaultAPI()

	cfg := &cmd.Config{
		Out:     os.Stdout,
		Storage: s,
		DictAPI: api,
	}
	if err := cmd.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
