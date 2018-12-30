package cmd

import (
	"fmt"

	"github.com/jhjaggars/gtl/pkg/todolist"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "toggle an item's done-ness",
	Run: func(cmd *cobra.Command, args []string) {
		tl := todolist.Read(Filename)
		for _, ch := range args {
			err := tl.ToggleDone(ch)
			if err != nil {
				fmt.Println(err)
			}
		}
		tl.Show(SortBy)
		tl.Save(Filename)
	},
}
