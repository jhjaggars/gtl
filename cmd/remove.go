package cmd

import (
	"fmt"

	"github.com/jhjaggars/gtl/pkg/todolist"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove item(s) from the list",
	Run: func(cmd *cobra.Command, args []string) {
		tl := todolist.Read(Filename)
		defer tl.Save(Filename)
		for _, ch := range args {
			_, prs := tl.Items[ch]
			if prs {
				delete(tl.Items, ch)
			} else {
				fmt.Printf("Invalid index '%s'.\n", ch)
			}
		}
	},
}
