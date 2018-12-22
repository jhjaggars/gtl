package cmd

import (
	"github.com/jhjaggars/gtl/pkg/todolist"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add an item to the list",
	Run: func(cmd *cobra.Command, args []string) {
		tl := todolist.Read(Filename)
		defer tl.Save(Filename)
		for _, desc := range args {
			tl.Add(desc)
		}
	},
}
