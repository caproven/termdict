package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove word ...",
	Short: "Remove words from your vocab list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, word := range args {
			if err := internal.RemoveWord(word); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Successfully removed word", word)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
