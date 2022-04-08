package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the words in your vocab list",
	Run: func(cmd *cobra.Command, args []string) {
		words, err := internal.ListWords()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
