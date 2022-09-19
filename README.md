# Todo.go

Todo app built in Go using an sqlite db

## Usage

1. Create a new database
  `todo init`
2. List all todos
  `todo list`
3. Add a new todo
  `todo add -n "<Name of todo>" -c "<Content of todo>" -p "<Priority {1=Low,2,3=High}>"`
4. Complete a todo
  `todo comp -id <Id of todo>`
5. Delete a todo
  `todo del -id <Id of todo>`
  
  
## Install

1. Install go
2. Clone into a folder

  2.1 `mkdir todo`
  
  2.2 `cd todo`
  
  2.3 `make build` <- Creates an executable in `/todo/bin/todo.exe`. Add this to path to use anywhere
  
