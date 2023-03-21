package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DbTable struct {
	dbName    string
	tableName string
	tableType todo
	// tableSchema map[string]string
}

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

func (d *DbTable) createDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.dbName)
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
	`, d.tableName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
	fmt.Println("Created table: ", d.tableName)
	db.Close()
	return db, err
}

func (d *DbTable) deleteDb() error {
	return os.Remove(d.dbName)
}

func (d *DbTable) deleteAll() error {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM todo;")
	return err
}

func (d *DbTable) insertTodo(t todo) (int, error) {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("INSERT INTO todo (name, content, priority, completed) VALUES (?, ?, ?, ?);", t.name, t.content, t.priority, t.completed)
	if err != nil {
		if err.Error() == "no such table: todo" {
			fmt.Println("Database not found. Run 'todo init' to create a new database.")
			os.Exit(1)
		}
		panic(err)
	}
	db.Close()
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), err
}

func (d *DbTable) getAllTodos(limit int) []todo {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM todo ORDER BY completed, priority DESC LIMIT ?;", limit)
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

func (d *DbTable) getTodoById(id int) todo {
	db, err := sql.Open("sqlite3", d.dbName)
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

func (d *DbTable) updateTodoById(id int, t todo) error {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("UPDATE todo SET name = ?, content = ?, priority = ?, completed = ? WHERE id = ?;", t.name, t.content, t.priority, t.completed, id)
	if err != nil {
		panic(err)
	}
	db.Close()
	return err
}

func (d *DbTable) deleteTodoById(id int) (int, error) {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec("DELETE FROM todo WHERE id = ?;", id)
	if err != nil {
		panic(err)
	}
	db.Close()
	var newId int64
	if newId, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(newId), err
}

func (d *DbTable) getTodosCount() int {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM todo;").Scan(&count)
	if err != nil {
		panic(err)
	}
	db.Close()
	return count
}

func (d *DbTable) getTodosByStatus(status int, limit int) []todo {
	db, err := sql.Open("sqlite3", d.dbName)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM todo WHERE completed = ? ORDER BY priority DESC LIMIT ?;", status, limit)
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
