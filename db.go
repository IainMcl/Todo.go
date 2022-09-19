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

func createDB(dbName string, tableName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	// Create table
	sqlStmt := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			content TEXT,
			priority INTEGER,
			completed INTEGER
		);
	`, tableName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
	fmt.Println("Created table: ", tableName)
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
	_, err = db.Exec("DELETE FROM todo;")
	return err
}

func insertTodo(dbName string, t todo) (int, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("INSERT INTO todo (name, content, priority, completed) VALUES (?, ?, ?, ?);", t.name, t.content, t.priority, t.completed)
	db.Close()
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), err
}

func getAllTodos(dbName string) []todo {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM todo ORDER BY priority DESC;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var todos []todo
	for rows.Next() {
		var t todo
		err = rows.Scan(&t.id, &t.name, &t.content, &t.priority, &t.completed)
		if err != nil {
			panic(err)
		}
		todos = append(todos, t)
	}
	db.Close()
	return todos
}

func getTodoById(dbName string, id int) todo {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	var t todo
	err = db.QueryRow("SELECT * FROM todo WHERE id = ?;", id).Scan(&t.id, &t.name, &t.content, &t.priority, &t.completed)
	if err != nil {
		panic(err)
	}
	db.Close()
	return t
}

func updateTodoById(dbName string, id int, t todo) (int, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("UPDATE todo SET name = ?, content = ?, priority = ?, completed = ? WHERE id = ?;", t.name, t.content, t.priority, t.completed, id)
	db.Close()
	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(newId), err
}

func deleteTodoById(dbName string, id int) (int, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("DELETE FROM todo WHERE id = ?;", id)
	db.Close()
	var newId int64
	if newId, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(newId), err
}
