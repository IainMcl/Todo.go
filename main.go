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
		// Must have a name
		f.PrintDefaults()
		os.Exit(1)
	}

	t := todo{id: 1, name: name, content: content, priority: Priority(priority), completed: 0}

	id, err := insertTodo(dbName, t)
	if err != nil {
		fmt.Println("Error inserting todo: ", err)
		os.Exit(1)
	}
	n := getTodoById(dbName, id)
	fmt.Println("Inserted todo: ", n)
	return t
}

func list(f *flag.FlagSet, dbName string) {
	todos := getAllTodos(dbName)

	for _, t := range todos {
		fmt.Println(t)
	}
}

func delete(f *flag.FlagSet, dbName string) {

	var id int
	f.IntVar(&id, "id", 0, "Id of todo to delete")
	f.Parse(os.Args[2:])
	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}

	todo := getTodoById(dbName, id)
	id, err := deleteTodoById(dbName, id)
	if err != nil {
		fmt.Println("Error deleting todo: ", err)
		os.Exit(1)
	}
	fmt.Println("Deleted todo: ", todo)
}

func complete(f *flag.FlagSet, dbName string) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to complete")
	f.Parse(os.Args[2:])
	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}
	todo := getTodoById(dbName, id)
	todo.completed = 1
	updatedId, err := updateTodoById(dbName, id, todo)
	if err != nil {
		fmt.Println("Error completing todo: ", err)
		os.Exit(1)
	}
	newTodo := getTodoById(dbName, updatedId)
	fmt.Println("Completed todo: ", newTodo)
}

func newDb(f *flag.FlagSet, dbName string, tableName string) {

	if _, err := os.Stat(dbName); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating database: ", dbName)
		_, err = createDB(dbName, tableName)
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB already exists")
	} else {
		fmt.Println("DB already exists: ", err)
	}
}

func main() {
	const dbName = "todo.db"
	const tableName = "todo"

	newCmd := flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	compCmd := flag.NewFlagSet("comp", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'init', 'add', 'del', 'comp' or 'list' subcommands")
		return
	}

	switch os.Args[1] {
	case "init":
		newDb(newCmd, dbName, tableName)
	case "add":
		add(addCmd, dbName)
	case "list":
		list(listCmd, dbName)
	case "del":
		delete(delCmd, dbName)
	case "comp":
		complete(compCmd, dbName)
	}
}
