package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove word ...",
	Short: "Remove words from your vocab list",
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

		for _, word := range args {
			if err := vl.RemoveWord(word); err != nil {
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
	rootCmd.AddCommand(removeCmd)
}
