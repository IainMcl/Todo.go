package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func add(d DbTable, f *flag.FlagSet) todo {
	var name string
	var content string
	var priority int
	f.StringVar(&name, "n", "", "Name of todo")
	f.StringVar(&content, "c", "", "Content of todo")
	f.IntVar(&priority, "p", 1, "Priority of todo (1 <low> - 3 <high>)")

	f.Parse(os.Args[2:])
	if len(name) == 0 {
		// Must have a name
		f.PrintDefaults()
		os.Exit(1)
	}

	t := todo{id: 1, name: name, content: content, priority: Priority(priority), completed: 0}

	id, err := d.insertTodo(t)
	if err != nil {
		fmt.Println("Error inserting todo: ", err)
		os.Exit(1)
	}
	n := d.getTodoById(id)
	fmt.Println("Inserted todo: ", n)
	return t
}

func list(d DbTable, f *flag.FlagSet) {
	var status string
	f.StringVar(&status, "s", "incomplete", "Status of todo")
	f.Parse(os.Args[2:])

	var todos []todo
	switch status {
	case "incomplete":
		todos = d.getTodosByStatus(0)
	case "complete":
		todos = d.getTodosByStatus(1)
	case "all":
		todos = d.getAllTodos()
	default:
		fmt.Println("Invalid status")
	}

	for _, t := range todos {
		fmt.Println(t)
	}
}

func delete(d DbTable, f *flag.FlagSet) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to delete")
	f.Parse(os.Args[2:])
	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}

	todo := d.getTodoById(id)
	id, err := d.deleteTodoById(id)
	if err != nil {
		fmt.Println("Error deleting todo: ", err)
		os.Exit(1)
	}
	fmt.Println("Deleted todo: ", todo)
}

func complete(d DbTable, f *flag.FlagSet) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to complete")
	f.Parse(os.Args[2:])

	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}
	todo := d.getTodoById(id)
	todo.completed = 1
	err := d.updateTodoById(id, todo)
	if err != nil {
		fmt.Println("Error completing todo: ", err)
		os.Exit(1)
	}
	newTodo := d.getTodoById(id)
	fmt.Println("Completed todo: ", newTodo)
}

func view(d DbTable, f *flag.FlagSet) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to view")
	f.Parse(os.Args[2:])
	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}
	todo := d.getTodoById(id)
	fmt.Println(todo)
}

func newDb(d DbTable, f *flag.FlagSet) {

	if _, err := os.Stat(d.dbName); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating database: ", d.dbName)
		_, err = d.createDB()
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB already exists")
	} else {
		fmt.Println("DB already exists: ", err)
	}
}

func update(d DbTable, f *flag.FlagSet) {
	var id int
	var name string
	var content string
	var priority int
	f.IntVar(&id, "id", 0, "Id of todo to update")
	f.StringVar(&name, "n", "", "Name of todo")
	f.StringVar(&content, "c", "", "Content of todo")
	f.IntVar(&priority, "p", 1, "Priority of todo (1 <low> - 3 <high>)")
	f.Parse(os.Args[2:])

	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}

	todo := d.getTodoById(id)
	if len(name) > 0 {
		todo.name = name
	}
	if len(content) > 0 {
		todo.content = content
	}
	if priority > 0 {
		todo.priority = Priority(priority)
	}
	err := d.updateTodoById(id, todo)
	if err != nil {
		fmt.Println("Error updating todo: ", err)
		os.Exit(1)
	}
	newTodo := d.getTodoById(id)
	fmt.Println("Updated todo: ", newTodo)
}

func main() {
	const dbName = "todo.db"
	const tableName = "todo"

	d := DbTable{dbName: dbName, tableName: tableName}

	newCmd := flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	compCmd := flag.NewFlagSet("comp", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'init', 'add', 'del', 'comp', 'view' or 'list' subcommands")
		return
	}

	switch os.Args[1] {
	case "init":
		newDb(d, newCmd)
	case "add":
		add(d, addCmd)
	case "list":
		list(d, listCmd)
	case "del":
		delete(d, delCmd)
	case "comp":
		complete(d, compCmd)
	case "view":
		view(d, compCmd)
	case "update":
		update(d, updateCmd)
	}
}
