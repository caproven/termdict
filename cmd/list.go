package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the words in your vocab list",
	Run: func(cmd *cobra.Command, args []string) {
		vf := storage.VocabFile{
			Path: storage.DefaultVocabFile(),
		}

		vl, err := vf.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, word := range vl.Words {
			fmt.Println(word)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
