package main

import (
	"fmt"
	"strconv"
	"strings"

	tsize "github.com/kopoli/go-terminal-size"
)

type ConsolePrint struct {
	width          int
	height         int
	idWidth        int
	nameWidth      int
	contentWidth   int
	priorityWidth  int
	completedWidth int
	colour         map[string]string
	prettyPrint    bool
}

func consoleSize() (int, int, bool) {
	size, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}
	// fmt.Println("Width: ", size.Width, "Height: ", size.Height)
	correctSize := true
	if size.Width < 80 || size.Height < 20 {
		fmt.Println("Terminal size too small for nice formatting")
		correctSize = false
	}
	return size.Width, size.Height, correctSize
}

func NewConsolePrint() *ConsolePrint {
	w, h, correctSize := consoleSize()
	wId := 4
	wPriority := 9
	wCompleted := 10

	remainingWidth := w - wId - wPriority - wCompleted - 13
	wName := remainingWidth / 3
	wContent := remainingWidth * 2 / 3

	totalWidth := wId + wName + wContent + wPriority + wCompleted + 13
	// If the division leave a remainder, add one to the content width
	if totalWidth != w {
		// fmt.Println("Width calculation error: ", totalWidth, w)
		wContent += w - totalWidth
	}

	return &ConsolePrint{
		width:  w,
		height: h,
		colour: map[string]string{
			"border":  "\033[34m",
			"error":   "\033[31m",
			"warning": "\033[33m",
			"success": "\033[32m",
			"null":    "\033[30m",
			"white":   "\033[37m",
		},
		idWidth:        wId,
		priorityWidth:  wPriority,
		completedWidth: wCompleted,
		nameWidth:      wName,
		contentWidth:   wContent,
		prettyPrint:    correctSize,
	}
}

func (c ConsolePrint) printHeader() {
	fmt.Println(c.colour["border"], strings.Repeat("=", c.width-2))
	// Id | Priority | Completed are fixed width elements. Name and content are variable width
	// Name and content should wrap to next line if they are too long
	// Print the column names for the table
	fmt.Println("|",
		c.colour["white"],
		"ID", strings.Repeat(" ", c.idWidth-2),
		"Name", strings.Repeat(" ", c.nameWidth-4),
		"Content", strings.Repeat(" ", c.contentWidth-7),
		"Priority", strings.Repeat(" ", c.priorityWidth-8),
		"Completed",
		c.colour["border"],
		"|")
	fmt.Println(c.colour["border"], strings.Repeat("=", c.width-2))
}

func (c ConsolePrint) printFooter() {
	fmt.Println(c.colour["border"], strings.Repeat("=", c.width-2))
}

func (c ConsolePrint) printDivider() {
	fmt.Println(c.colour["border"], strings.Repeat("-", c.width-2))
}

func (c ConsolePrint) printTodo(t todo) {
	var check string
	if t.completed == 1 {
		check = "\u2713"
	} else {
		check = "\u2717"
	}

	contentWraps := len(t.content) / c.contentWidth
	nameWraps := len(t.name) / c.nameWidth

	var nWraps int
	if nameWraps > contentWraps {
		nWraps = nameWraps
	} else {
		nWraps = contentWraps
	}

	// fmt.Println("Number of wraps: ", nWraps)

	for i := 0; i <= nWraps; i++ {
		var nameString string
		if i <= nameWraps && len(t.name) > c.nameWidth {
			if c.nameWidth > len(t.name[c.nameWidth*i:]) {
				nameString = t.name[i*c.nameWidth:]
			} else {
				nameString = t.name[i*c.nameWidth : (i+1)*c.nameWidth]
			}
		} else if i == 0 {
			nameString = t.name
		} else {
			nameString = strings.Repeat(" ", c.nameWidth)
		}
		var contentString string
		if i <= contentWraps && len(t.content) > c.contentWidth {
			if c.contentWidth > len(t.content[c.contentWidth*i:]) {
				contentString = t.content[i*c.contentWidth:]
			} else {
				contentString = t.content[i*c.contentWidth : (i+1)*c.contentWidth]
			}
		} else if i == 0 {
			contentString = t.content
		} else {
			contentString = strings.Repeat(" ", c.contentWidth)
		}
		if i == 0 {
			fmt.Println("|",
				c.colour["white"],
				t.id, strings.Repeat(" ", c.idWidth-len(strconv.FormatInt(int64(t.id), 10))),
				nameString, strings.Repeat(" ", c.nameWidth-len(nameString)),
				contentString, strings.Repeat(" ", c.contentWidth-len(contentString)+7),
				t.priority, strings.Repeat(" ", c.priorityWidth-len(strconv.FormatInt(int64(t.priority), 10))+1),
				check,
				c.colour["border"],
				"|")
		} else {
			fmt.Println("|",
				c.colour["white"],
				strings.Repeat(" ", c.idWidth+len(strconv.FormatInt(int64(t.id), 10))),
				nameString, strings.Repeat(" ", c.nameWidth-len(nameString)),
				contentString, strings.Repeat(" ", c.contentWidth-len(contentString)),
				strings.Repeat(" ", c.priorityWidth),
				strings.Repeat(" ", c.completedWidth),
				c.colour["border"],
				"|")
		}
	}
}

func (c ConsolePrint) resetColour() {
	fmt.Println(c.colour["null"])
}

func (c ConsolePrint) printTodos(todos []todo) {
	if !c.prettyPrint {
		for _, t := range todos {
			fmt.Println(t)
		}
		return
	}
	c.printHeader()
	for _, t := range todos {
		c.printTodo(t)
		c.printDivider()
	}
	c.resetColour()
}
