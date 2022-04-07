package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var word string

	flag.StringVar(&word, "w", "", "word to define")
	flag.Parse()

	if word == "" {
		fmt.Println("Must provide word")
		os.Exit(1)
	}

	defs, err := define(word)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Definitions for %s:\n", word)
	for _, def := range defs {
		fmt.Printf("\t%s - %s\n", def.PartOfSpeech, def.Meaning)
	}
}
