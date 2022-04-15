package main

import (
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/vocab"
)

func main() {
	if err := vocab.CreateConfigDir(); err != nil {
		panic(err)
	}

	s := vocab.File{
		Path: vocab.DefaultFilepath(),
	}

	cfg := &cmd.Config{
		Out:     os.Stdout,
		Storage: s,
	}
	if err := cmd.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
