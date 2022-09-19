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
	defer db.Close()

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
		panic(err)
	}
	fmt.Println("Created database: ", tableName)
	return db, err
}

func deleteDb(name string) error {
	return os.Remove(name)
}

func deleteAll(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todo")
	return err
}
