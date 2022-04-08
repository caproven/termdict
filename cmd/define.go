package cmd

import (
	"fmt"
	"os"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define word",
	Short: "Lookup the definition of a word",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defs, err := internal.Define(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, def := range defs {
			fmt.Printf("[%s] %s\n", def.PartOfSpeech, def.Meaning)
		}
	},
}

func init() {
	rootCmd.AddCommand(defineCmd)
}
