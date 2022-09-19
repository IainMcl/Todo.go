package main

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

type todo struct {
	id        int
	name      string
	content   string
	priority  Priority
	completed int
}

// type db struct {
// 	db     string
// 	fields []todo
// }
