package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func dbType(goType string) string {
	switch goType {
	case "int":
		return "INTEGER"
	case "string":
		return "TEXT"
	case "float64":
		return "REAL"
	case "bool":
		return "INTEGER"
	default:
		return "BLOB"
	}
}

func createDB(tableName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", tableName)
	if err != nil {
		panic(err)
	}

	// Create table
	sqlStmt := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			content TEXT,
			priority INTEGER
		);
	`, tableName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
	fmt.Println("Created database: ", tableName)
	db.Close()
	return db, err
}

func deleteDb(name string) error {
	return os.Remove(name)
}

func deleteAll(dbName string) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM todo")
	return err
}

func insertTodo(dbName string, t todo) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO todo (name, content, priority) VALUES (?, ?, ?)", t.name, t.content, t.priority)
	db.Close()
	return err
}
