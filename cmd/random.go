package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/caproven/termdict/internal/dictionary"
	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Define a random word from your vocab list",
	Run: func(cmd *cobra.Command, args []string) {
		vf := storage.VocabFile{
			Path: storage.DefaultVocabFile(),
		}

		vl, err := vf.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(vl.Words) == 0 {
			fmt.Println("no words in vocab list")
			os.Exit(1)
		}

		word := vl.Words[rand.Intn(len(vl.Words))]

		defs, err := dictionary.Define(word)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printDefinition(os.Stdout, word, defs)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.AddCommand(randomCmd)
}
