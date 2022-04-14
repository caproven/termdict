package main

import (
	"os"

	"github.com/caproven/termdict/cmd"
	"github.com/caproven/termdict/vocab"
)

func main() {
	s := vocab.File{
		Path: vocab.DefaultFilepath(),
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
