package cmd

import (
	"fmt"
	"os"

	"github.com/jhjaggars/gtl/pkg/todolist"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gtl",
	Short: "gtl is a todo list",
	Long:  `a simple todo list application`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

var Filename string
var SortBy string

func main() {
	tl := todolist.Read(Filename)
	tl.Show(SortBy)
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&Filename,
		"filename",
		"f",
		"./todo.json",
		"file containing todo list")
	rootCmd.PersistentFlags().StringVarP(
		&SortBy,
		"sortby",
		"s",
		"description",
		"field to sort by",
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
