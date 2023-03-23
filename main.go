package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Read internal config
	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	d := &DbTable{dbName: config.GetDbName(), tableName: config.GetTableName()}

	newCmd := flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	compCmd := flag.NewFlagSet("comp", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	expectedInput := "Expected 'init', 'add', 'del', 'comp', 'view', 'update', 'list', or 'config' subcommands"

	inputHelp :=
		`Usage of todo:
  todo init
      Initialize a new todo database and config file
  todo add 
	  Add a new todo item
  todo del
	  Delete a todo item
  todo comp
	  Mark a todo item as complete
  todo view
	  View an individual todo item
  todo update
	  Update a todo item
  todo list
	  List multiple todo items
  todo config
	  View or update config values - Not yet implemented
`

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
	case "help":
		fallthrough
	case "-h":
		fallthrough
	case "--help":
		fmt.Println(inputHelp)
	default:
		fmt.Println(expectedInput)
	}
}
