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

func main() {
	tl := todolist.Read(Filename)
	tl.Show()
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&Filename,
		"filename",
		"f",
		"./todo.json",
		"file containing todo list")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
