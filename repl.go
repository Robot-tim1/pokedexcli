package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	UpArrow    = "\x1b[A"
	DownArrow  = "\x1b[B"
	RightArrow = "\x1b[C"
	LeftArrow  = "\x1b[D"
	ClearLine  = "\x1b[2K\r"
)

func startRepl(cfg *config) {
	var history []string
	historyIndex := -1
	var currentInput []byte
	currentInputIndex := -1
	buf := make([]byte, 512)

	prompt := "Pokedex > "
	fmt.Print(prompt)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("%v\r\n", err)
			return
		}

		inputStr := string(buf[:n])
		switch inputStr {
		case UpArrow:
			upInput(history, &historyIndex, &currentInput, &currentInputIndex, prompt)
		case DownArrow:
			downInput(history, &historyIndex, &currentInput, &currentInputIndex, prompt)
		case "\r", "\n":
			if enterInput(&history, &historyIndex, &currentInput, &currentInputIndex, prompt, cfg) {
				return
			}
		case LeftArrow:
			leftInput(&currentInputIndex)
		case RightArrow:
			rightInput(&currentInputIndex, currentInput)
		case "\x7f", "\x08":
			backInput(&currentInput, &currentInputIndex, prompt)
		case "\x03":
			fmt.Print("^C\r\n")
			return
		default:
			if n == 1 && buf[0] >= 32 && buf[0] <= 126 {
				currentInput = append(currentInput, buf[0])
				fmt.Print(string(buf[0]))
				currentInputIndex++
			}
		}
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	cleanText := strings.Fields(lowerText)
	return cleanText
}
