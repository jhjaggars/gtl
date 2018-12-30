package todolist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ttacon/chalk"
)

const checkmark string = "\u2713"
const circle string = "\u25cf"
const week = 7 * 24 * time.Hour
const month = 30 * 24 * time.Hour

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
	}
	elapsed := time.Since(item.Updated)
	switch {
	case elapsed > month:
		return chalk.Red.Color(circle)
	case elapsed > week:
		return chalk.Yellow.Color(circle)
	default:
		return circle
	}
}

func (item *TodoItem) Show() {
	fmt.Printf("%s %s %s %s\n",
		item.Index,
		item.getDone(),
		item.Description,
		chalk.White.NewStyle().
			WithBackground(chalk.Black).
			WithTextStyle(chalk.Dim).
			Style(humanize.Time(item.Updated)),
	)
}

func (item *TodoItem) Toggle() {
	item.Done = !item.Done
}

type By func(item1, item2 *TodoItem) bool

type itemSorter struct {
	items []TodoItem
	by    func(item1, item2 *TodoItem) bool
}

func (s *itemSorter) Len() int {
	return len(s.items)
}

func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *itemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j])
}

func (by By) Sort(items []TodoItem) {
	is := &itemSorter{
		items: items,
		by:    by,
	}
	sort.Sort(is)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

var SortMappings = map[string]By{
	"index": func(item1, item2 *TodoItem) bool {
		return item1.Index < item2.Index
	},
	"description": func(item1, item2 *TodoItem) bool {
		return item1.Description < item2.Description
	},
	"done": func(item1, item2 *TodoItem) bool {

		return boolToInt(item1.Done) < boolToInt(item2.Done)
	},
	"created": func(item1, item2 *TodoItem) bool {
		return item1.Created.Unix() < item2.Created.Unix()
	},
	"updated": func(item1, item2 *TodoItem) bool {
		return item1.Updated.Unix() < item2.Updated.Unix()
	},
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
		if bytes.IndexByte(used, value) == -1 {
			return string(value), nil
		}
	}

	return "", errors.New("no more available indices")
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

func (todolist *TodoList) Show(sortcolumn string) {
	by, ok := SortMappings[sortcolumn]
	if !ok {
		by = SortMappings["updated"]
	}
	var items []TodoItem
	for k := range todolist.Items {
		items = append(items, todolist.Items[k])
	}
	By(by).Sort(items)
	for _, v := range items {
		v.Show()
	}
}

func (todolist *TodoList) ToggleDone(index string) error {
	item, prs := todolist.Items[index]
	if !prs {
		return errors.New("index did not exist")
	}
	item.Toggle()
	item.Updated = time.Now()
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
