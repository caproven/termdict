package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add word ...",
	Short: "Add words to your vocab list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, word := range args {
			if err := internal.AddWord(word); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Successfully added word", word)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
