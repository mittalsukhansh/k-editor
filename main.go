package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	//"golang.org/x/text/width"
)

func main() {

	var startingBuffer []rune
	var filename string

	if len(os.Args) > 1 {
		filename = os.Args[1]
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		} else if err == nil {
			startingBuffer = []rune(string(data))
		}
	}

	m := model{
		buffer:   startingBuffer,
		cursor:   0,
		filename: os.Args[1],
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// model
type model struct {
	cursor   int // which to-do list item our cursor is pointing at
	buffer   []rune
	height   int
	width    int
	filename string
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// update ->  The update function is called when “things happen.”
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		case "left":
			if m.cursor > 0 {
				m.cursor -= 1
			}

		case "right":
			if m.cursor < len(m.buffer) {
				m.cursor += 1
			}
		case "backspace":
			if m.cursor > 0 {
				m.buffer = append(m.buffer[:m.cursor-1], m.buffer[m.cursor:]...)
				m.cursor -= 1
			}
		case "space":
			var newBuffer []rune
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, ' ')
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "tab":
			var newBuffer []rune
			spaces := []rune("    ")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, spaces...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 4

		case "enter":
			var newBuffer []rune
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, '\n')
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "ctrl+s":
			data := []byte(string(m.buffer))
			err := os.WriteFile(m.filename, data, 0644)
			if err != nil {
				panic(err)
			}

		case "[":
			var newBuffer []rune
			bracket := []rune("[]")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, bracket...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "{":
			var newBuffer []rune
			bracket := []rune("{}")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, bracket...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "(":
			var newBuffer []rune
			bracket := []rune("()")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, bracket...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "'":
			var newBuffer []rune
			quotes := []rune("''")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, quotes...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "\"":
			var newBuffer []rune
			quotes := []rune("\"\"")
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, quotes...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)
			m.buffer = newBuffer
			m.cursor += 1

		case "up":
			currentLineStart := m.find_start_of_current_row(m.cursor)
			xoffset := m.cursor - currentLineStart
			prevLineEnd := currentLineStart - 1
			prevLineStart := m.find_start_of_current_row(prevLineEnd)

			if prevLineStart+xoffset <= prevLineEnd {
				m.cursor = prevLineStart + xoffset
			} else {
				m.cursor = prevLineEnd
			}

		default:
			runes := []rune(msg.String())

			var newBuffer []rune
			newBuffer = append(newBuffer, m.buffer[:m.cursor]...)
			newBuffer = append(newBuffer, runes...)
			newBuffer = append(newBuffer, m.buffer[m.cursor:]...)

			m.buffer = newBuffer
			m.cursor += len(runes)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil
	}

	return m, nil
}

//view -> to render our UI

func (m model) View() tea.View {

	// //s := "\n WELCOME TO THE FUCKING TERMINAL \n\n Press 'q' or 'ctrl + c' to quit"
	// leftSide := string(m.buffer[:m.cursor])
	// rightSide := string(m.buffer[m.cursor:])

	// v := tea.NewView(leftSide + "█" + rightSide)
	// v.AltScreen = true

	var lines []string
	var currentLine string
	for i, char := range m.buffer {

		if i == m.cursor {
			currentLine += "|"
		}

		if char != '\n' {
			currentLine += string(char)
		}

		if len(currentLine) >= m.width || char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		}
	}

	if m.cursor >= len(m.buffer) {
		currentLine += "|"
	}

	lines = append(lines, currentLine)

	v := tea.NewView(strings.Join(lines, "\n"))
	v.AltScreen = true
	return v
}

func (m model) find_start_of_current_row(pos int) int {

	if pos <= 0 {
		return 0
	}
	//finding start of currentline
	for j := pos - 1; j >= 0; j-- {
		if m.buffer[j] == '\n' {
			// scanner1 = j+1
			return j + 1
		}
	}
	return 0
}
