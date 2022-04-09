package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/caproven/termdict/internal"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove word ...",
	Short: "Remove words from your vocab list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.RemoveWords(args...); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) == 1 {
			fmt.Println("Successfully removed word", args[0])
		} else {
			fmt.Println("Successfully removed words", strings.Join(args, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
