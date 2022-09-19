package main

func newTodo(name string) *todo {
	t := todo{id: 1, name: name, priority: High}
	return &t
}

func (t *todo) add() todo {

	return todo{}
}
