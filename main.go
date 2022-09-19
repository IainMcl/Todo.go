package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func add(f *flag.FlagSet, dbName string) todo {
	var name string
	var content string
	var priority int
	f.StringVar(&name, "n", "", "Name of todo")
	f.StringVar(&content, "c", "", "Content of todo")
	f.IntVar(&priority, "p", 1, "Content of todo")

	f.Parse(os.Args[2:])
	if len(name) == 0 {
		f.PrintDefaults()
		os.Exit(1)
	}

	t := todo{id: 1, name: name, content: content, priority: Priority(priority)}

	insertTodo(dbName, t)
	return t
}

func list(f *flag.FlagSet, dbName string) {
	fmt.Println("List")
}

func delete() {
	fmt.Println("delete")
}

func newDb(f *flag.FlagSet, dbName string) {
	if _, err := os.Stat(dbName); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating database: ", dbName)
		_, err = createDB(dbName)
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB already exists")
	} else {
		panic(err)
	}
}

func main() {
	const dbName = "todo.db"

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'new', 'add' or 'list' subcommands")
		return
	}

	switch os.Args[1] {
	case "new":
		newDb(newCmd, dbName)
	case "add":
		add(addCmd, dbName)
	case "list":
		list(listCmd, dbName)
	}
}
