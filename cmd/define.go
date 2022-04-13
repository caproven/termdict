package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal/dictionary"
	"github.com/spf13/cobra"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define WORD",
	Short: "Define a word",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]
		dict := dictionary.Default()

		defs, err := dict.Define(word)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printDefinition(os.Stdout, word, defs)
	},
}

func init() {
	rootCmd.AddCommand(defineCmd)
}
