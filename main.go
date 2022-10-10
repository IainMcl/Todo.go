package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func add(d *DbTable, f *flag.FlagSet) todo {
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
	fmt.Println("Inserted todo: ")
	NewConsolePrint().printTodos([]todo{n})
	return t
}

func list(d *DbTable, f *flag.FlagSet) {
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

	NewConsolePrint().printTodos(todos)
}

func delete(d *DbTable, f *flag.FlagSet) {
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

func complete(d *DbTable, f *flag.FlagSet) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to complete")
	f.Parse(os.Args[2:])

	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}
	todoFetched := d.getTodoById(id)
	todoFetched.completed = 1
	err := d.updateTodoById(id, todoFetched)
	if err != nil {
		fmt.Println("Error completing todo: ", err)
		os.Exit(1)
	}
	newTodo := d.getTodoById(id)
	fmt.Println("Completed todo: ")
	NewConsolePrint().printTodos([]todo{newTodo})
}

func view(d *DbTable, f *flag.FlagSet) {
	var id int
	f.IntVar(&id, "id", 0, "Id of todo to view")
	f.Parse(os.Args[2:])
	if id == 0 {
		// Must have an id
		f.PrintDefaults()
		os.Exit(1)
	}
	todoView := d.getTodoById(id)
	NewConsolePrint().printTodos([]todo{todoView})
}

func newDb(d *DbTable, f *flag.FlagSet, config IConfig) {
	// Create config file if one doesn't exist
	if _, err := os.Stat(config.GetConfigPath()); os.IsNotExist(err) {
		fmt.Println("Creating config file: ", config.GetConfigPath())
		err := config.CreateDefaultConfig()
		if err != nil {
			fmt.Println("Error creating config file: ", err)
			os.Exit(1)
		}
	}
	// Create database file if one doesn't exist
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

func update(d *DbTable, f *flag.FlagSet) {
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

	todoUpdate := d.getTodoById(id)
	if len(name) > 0 {
		todoUpdate.name = name
	}
	if len(content) > 0 {
		todoUpdate.content = content
	}
	if priority > 0 {
		todoUpdate.priority = Priority(priority)
	}
	err := d.updateTodoById(id, todoUpdate)
	if err != nil {
		fmt.Println("Error updating todo: ", err)
		os.Exit(1)
	}
	newTodo := d.getTodoById(id)
	fmt.Println("Updated todo: ")
	NewConsolePrint().printTodos([]todo{newTodo})
}

func configCmd(args []string, config IConfig) {
	switch args[1] {
	case "view":
		err := config.View()
		if err != nil {
			fmt.Println("Error viewing config: ", err)
			os.Exit(1)
		}
	case "set":
		err := config.Set(args[2], args[3])
		if err != nil {
			fmt.Println("Error setting config: <", args[2], "> with value <", args[3], ">: ", err)
			os.Exit(1)
		}
	case "del":
		err := config.Delete(args[2])
		if err != nil {
			fmt.Println("Error deleting config <", args[2], ">: ", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid config command. Valid commands are: view, set, del")
	}
}

func main() {
	// Read internal config
	var config IConfig = &Config{
		ConfigPath: "~/.todo/config.json",
		DbName:     "~/.todo/todo.db",
		TableName:  "todo",
	}

	// err := config.ReadConfig() // TODO: Read user config and update config with its contents. This should probably be done in the Config struct
	// if err != nil {
	// 	panic(err)
	// }

	d := &DbTable{dbName: config.GetDbName(), tableName: config.GetTableName()}

	newCmd := flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	compCmd := flag.NewFlagSet("comp", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	expectedInput := "Expected 'init', 'add', 'del', 'comp', 'view', 'update', 'list', or 'config' subcommands"
	if len(os.Args) < 2 {
		fmt.Println(expectedInput)
		return
	}

	switch os.Args[1] {
	case "init":
		newDb(d, newCmd, config)
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
	case "config":
		configCmd(os.Args[2:], config)
	default:
		fmt.Println(expectedInput)
	}
}
