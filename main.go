package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func add(f *flag.FlagSet) todo {
	var name string
	var content string
	var priority int
	f.StringVar(&name, "n", "", "Name of todo")
	f.StringVar(&content, "c", "Default content", "Content of todo")
	f.IntVar(&priority, "p", 1, "Content of todo")

	if len(name) == 0 {
		f.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Hello %s\n", name)
	fmt.Printf("Content: %s\n", content)
	fmt.Printf("Priority: %d\n", priority)

	t := todo{id: 1, name: name, content: content, priority: Priority(priority)}
	return t
}

func list(f *flag.FlagSet) {
	fmt.Println("List")
}

func delete() {
	fmt.Println("delete")
}

func new() {
	const dbName = "todo.db"
	if _, err := os.Stat(dbName); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating database: ", dbName)
		_, err := createDB(dbName)
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB already exists")
	}
}

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add' or 'list' subcommands")
		return
	}
	switch os.Args[1] {
	case "add":
		add(addCmd)
	case "list":
		list(listCmd)
	}
}
