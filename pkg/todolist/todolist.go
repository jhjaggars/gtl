package todolist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ttacon/chalk"
)

const checkmark string = "\u2713"
const circle string = "\u25cf"

type TodoItem struct {
	Index       string
	Description string
	Done        bool
	Created     time.Time
	Updated     time.Time
}

func (item *TodoItem) getDone() string {
	if item.Done {
		return chalk.Green.Color(checkmark)
	} else {
		return circle
	}
}

func (item *TodoItem) Show() {
	fmt.Printf("%s %s %s %s\n",
		item.Index,
		item.getDone(),
		item.Description,
		chalk.Black.NewStyle().
			WithTextStyle(chalk.Bold).
			Style(humanize.Time(item.Updated)),
	)
}

func (item *TodoItem) Toggle() {
	item.Done = !item.Done
}

func getbytes(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}

func Read(filename string) TodoList {
	var items map[string]TodoItem

	err := json.Unmarshal(getbytes(filename), &items)
	if err != nil {
		panic(err)
	}
	return TodoList{
		Items: items,
		File:  filename,
	}
}

type TodoList struct {
	Items map[string]TodoItem
	File  string
}

func (todolist *TodoList) getNextIndex() (string, error) {
	options := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	used := make([]byte, len(todolist.Items))

	i := 0
	for key := range todolist.Items {
		used[i] = []byte(key)[0]
		i++
	}

	for _, value := range options {
		if bytes.IndexByte(used, value) != -1 {
			return string(value), nil
		}
	}

	return "", errors.New("No more available indices.")
}

func (todolist *TodoList) Add(description string) error {
	i, err := todolist.getNextIndex()
	if err != nil {
		return err
	}
	todolist.Items[i] = TodoItem{
		Index:       i,
		Done:        false,
		Description: description,
		Created:     time.Now(),
		Updated:     time.Now(),
	}
	return nil
}

func (todolist *TodoList) Show() {
	for _, v := range todolist.Items {
		v.Show()
	}
}

func (todolist *TodoList) ToggleDone(index string) error {
	item, prs := todolist.Items[index]
	if !prs {
		return errors.New("Index did not exist.")
	}
	item.Toggle()
	todolist.Items[index] = item
	return nil
}

func (todolist *TodoList) Save(filename string) error {
	bytes, err := json.Marshal(todolist.Items)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bytes, 0600)
	return err
}
