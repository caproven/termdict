package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal/storage"
	"github.com/spf13/cobra"
)

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

		for _, word := range args {
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
}
