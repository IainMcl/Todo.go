package main

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

type todo struct {
	id       int
	name     string
	content  string
	priority Priority
}

// type db struct {
// 	db     string
// 	fields []todo
// }
