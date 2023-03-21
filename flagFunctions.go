package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	var limit int
	f.StringVar(&status, "s", "incomplete", "Status of todo")
	f.IntVar(&limit, "l", 10, "Limit number of todos to return")
	f.Parse(os.Args[2:])

	var todos []todo
	switch status {
	case "incomplete":
		todos = d.getTodosByStatus(0, limit)
	case "complete":
		todos = d.getTodosByStatus(1, limit)
	case "all":
		todos = d.getAllTodos(limit)
	default:
		fmt.Println("Invalid status")
	}

	NewConsolePrint().printTodos(todos)
	countTodos := d.getTodosCount()
	returnedTodos := len(todos)
	if returnedTodos < countTodos {
		fmt.Printf("Showing %d of %d todos", returnedTodos, countTodos)
	}
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

func readConfig() (*Config, error) {
	// Read config file from default location ~/.config/todo/config.json
	// If the config file does not exist then create it with default values

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// join home directory with config file path
	configPath := filepath.Join(home, ".todo", "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// Config file does not exist
		// Create config file with default values
		config := NewConfig()
		// save config file
		err = config.WriteConfig()
		return config, err
	}

	// Config file exists
	// Read config file
	config, err := ReadConfig(configPath)

	return config, err
}

func newDb(d *DbTable, f *flag.FlagSet, config *Config) {
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

func configCmd(args []string, config *Config) {
	configOptions := "Invalid config command. Valid commands are: view, set, del"
	if len(args) < 1 {
		fmt.Print("Too few arguments\n", configOptions)
		os.Exit(1)
	}
	switch args[0] {
	case "view":
		err := config.View()
		if err != nil {
			fmt.Println("Error viewing config: ", err)
			os.Exit(1)
		}
	case "set":
		if len(args) < 3 {
			fmt.Println("Too few arguments for config set")
			fmt.Println("Usage: todo config set <key> <value>")
			os.Exit(1)
		}
		err := config.Set(args[1], args[2])
		if err != nil {
			fmt.Println("Error setting config: <", args[1], "> with value <", args[2], ">: ", err)
			os.Exit(1)
		}
	case "del":
		err := config.Delete(args[1])
		if err != nil {
			fmt.Println("Error deleting config <", args[1], ">: ", err)
			os.Exit(1)
		}
	default:
		fmt.Println(configOptions)
	}
}
