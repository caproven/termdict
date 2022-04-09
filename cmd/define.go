package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/caproven/termdict/internal/dictionary"
	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define word",
	Short: "Lookup the definition of a word",
	Args: func(cmd *cobra.Command, args []string) error {
		random, err := cmd.Flags().GetBool("random")
		if err != nil {
			return errors.New("failed to get random flag")
		}

		if !random && len(args) == 0 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var word string

		random, err := cmd.Flags().GetBool("random")
		if err != nil {
			fmt.Println("failed to get random flag")
			os.Exit(1)
		}

		if random {
			vf := storage.VocabFile{
				Path: storage.DefaultVocabFile(),
			}

			vl, err := vf.Read()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			rand.Seed(time.Now().UnixNano())
			word = vl.Words[rand.Intn(len(vl.Words))]
		} else {
			word = args[0]
		}

		defs, err := dictionary.Define(word)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(word)
		for _, def := range defs {
			fmt.Printf("[%s] %s\n", def.PartOfSpeech, def.Meaning)
		}
	},
}

func init() {
	rootCmd.AddCommand(defineCmd)

	defineCmd.Flags().BoolP("random", "r", false, "Define a random word from your vocab list")
}
