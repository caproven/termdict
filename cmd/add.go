package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add word ...",
	Short: "Add words to your vocab list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.AddWords(args...); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) == 1 {
			fmt.Println("Successfully added word", args[0])
		} else {
			fmt.Println("Successfully added words", strings.Join(args, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
