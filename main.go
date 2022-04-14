package main

import (
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/storage"
)

func main() {
	s := storage.VocabFile{
		Path: storage.DefaultVocabFile(),
	}

	cfg := &cmd.Config{
		Out:     os.Stdout,
		Storage: s,
	}
	err := cmd.NewRootCmd(cfg).Execute()
	if err != nil {
		os.Exit(1)
	}
}
