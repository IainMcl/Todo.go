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
