package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal/dictionary"
	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

var checkFlag bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add word ...",
	Short: "Add words to your vocab list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vf := storage.VocabFile{
			Path: storage.DefaultVocabFile(),
		}

		vl, err := vf.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dict := dictionary.Default()

		for _, word := range args {
			if checkFlag {
				if _, err := dict.Define(word); err != nil {
					fmt.Printf("failed to add word '%s', couldn't find a definition\n", word)
					os.Exit(1)
				}
			}

			if err := vl.AddWord(word); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if err := vf.Write(vl); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&checkFlag, "check", "c", true, "check that words can be defined before adding")
}
