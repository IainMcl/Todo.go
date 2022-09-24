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
	if size.Width < 80 || size.Height < 10 {
		fmt.Println("Terminal size too small for nice formatting")
		correctSize = false
	}
	return size.Width, size.Height, correctSize
}

func NewConsolePrint() *ConsolePrint {
	w, h, correctSize := consoleSize()
	wId := 5
	wPriority := 9
	wCompleted := 10

	remainingWidth := w - wId - wPriority - wCompleted
	wName := remainingWidth / 4
	if wName < 20 {
		wName = 20
	}

	totalWidth := wId + wName + wPriority + wCompleted + 5 // 5 is for the 4 spaces and the border
	wContent := w - totalWidth

	return &ConsolePrint{
		width:  w,
		height: h,
		colour: map[string]string{
			"border":    "\033[34m",
			"error":     "\033[31m",
			"warning":   "\033[33m",
			"success":   "\033[32m",
			"null":      "\033[37m",
			"white":     "\033[37m",
			"bold":      "\033[1m",
			"italic":    "\033[3m",
			"underline": "\033[4m",
			"normal":    "\033[0m",
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
	c.printHeaderDivider()
	// Id | Priority | Completed are fixed width elements. Name and content are variable width
	// Name and content should wrap to next line if they are too long
	// Print the column names for the table
	fmt.Print("| ",
		c.colour["white"],
		c.colour["bold"],
		"ID", strings.Repeat(" ", c.idWidth),
		"Name", strings.Repeat(" ", c.nameWidth-4),
		"Content", strings.Repeat(" ", c.contentWidth-7),
		"Priority", strings.Repeat(" ", c.priorityWidth-8),
		"Completed",
		c.colour["normal"],
		c.colour["border"],
		" |\n")
	c.printHeaderDivider()
}

func (c ConsolePrint) printHeaderDivider() {
	fmt.Println(c.colour["border"], strings.Repeat("=", c.width-2))
}

func (c ConsolePrint) printFooter() {
	fmt.Println(c.colour["border"], strings.Repeat("=", c.width-2))
}

func (c ConsolePrint) printDivider() {
	fmt.Println(c.colour["border"], strings.Repeat("-", c.width-2))
}

func getLineContent(content string, width int, lineNumber int) string {
	var nameString string
	if lineNumber <= width && len(content) > width {
		if width > len(content[width*lineNumber:]) {
			nameString = content[lineNumber*width:]
		} else {
			nameString = content[lineNumber*width : (lineNumber+1)*width]
		}
	} else if lineNumber == 0 {
		nameString = content
	} else {
		nameString = ""
	}
	if nameString != "" && nameString[0] == ' ' {
		nameString = nameString[1:]
	}
	return nameString
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

	for i := 0; i <= nWraps; i++ {

		nameString := getLineContent(t.name, c.nameWidth, i)
		contentString := getLineContent(t.content, c.contentWidth, i)

		if i == 0 {
			fmt.Print("| ",
				c.colour["white"],
				t.id, strings.Repeat(" ", c.idWidth-len(strconv.FormatInt(int64(t.id), 10))+2),
				nameString, strings.Repeat(" ", c.nameWidth-len(nameString)),
				contentString, strings.Repeat(" ", c.contentWidth-len(contentString)),
				strings.Repeat(" ", len("priority")),
				t.priority, strings.Repeat(" ", c.priorityWidth-len(strconv.FormatInt(int64(t.priority), 10))),
				check,
				c.colour["border"],
				" |")
		} else {
			fmt.Print("| ",
				c.colour["white"],
				strings.Repeat(" ", c.idWidth+2),
				nameString, strings.Repeat(" ", c.nameWidth-len(nameString)),
				contentString, strings.Repeat(" ", c.contentWidth-len(contentString)),
				strings.Repeat(" ", c.priorityWidth),
				strings.Repeat(" ", c.completedWidth-1),
				c.colour["border"],
				" |")
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
