package main

import (
	"fmt"
	"slices"
)

func upInput(history []string, historyIndex *int, currentInput *[]byte, currentInputIndex *int, prompt string) {
	if len(history) > 0 {
		if *historyIndex == -1 {
			*historyIndex = len(history) - 1
		} else if *historyIndex != 0 {
			*historyIndex--
		}
		lastCmd := history[*historyIndex]
		*currentInput = []byte(lastCmd)
		*currentInputIndex = len(*currentInput) - 1

		fmt.Print(ClearLine + prompt + lastCmd)
	}
}

func downInput(history []string, historyIndex *int, currentInput *[]byte, currentInputIndex *int, prompt string) {
	if *historyIndex != -1 {
		if *historyIndex != len(history)-1 {
			*historyIndex++
		} else {
			*historyIndex = -1
		}

		var lastCmd string
		if *historyIndex != -1 {
			lastCmd = history[*historyIndex]
		} else {
			lastCmd = ""
		}
		*currentInput = []byte(lastCmd)
		*currentInputIndex = len(*currentInput) - 1

		fmt.Print(ClearLine + prompt + lastCmd)
	}
}

func leftInput(currentInputIndex *int) {
	if *currentInputIndex >= 0 {
		fmt.Print(LeftArrow)
		*currentInputIndex--
	}
}

func rightInput(currentInputIndex *int, currentInput []byte) {
	if *currentInputIndex < len(currentInput)-1 {
		fmt.Print(RightArrow)
		*currentInputIndex++
	}
}

func enterInput(history *[]string, historyIndex *int, currentInput *[]byte, currentInputIndex *int, prompt string, cfg *config) bool {
	dirtyInput := string(*currentInput)
	if dirtyInput != "" {
		fmt.Print("\r\n")
		input := cleanInput(dirtyInput)
		if len(input) == 0 {
			fmt.Print(prompt)
			return false
		}

		command := input[0]
		var args []string
		if len(input) > 1 {
			args = input[1:]
		}

		if commandStruct, ok := commandRegistry[command]; ok {
			err := commandStruct.callback(cfg, args...)
			if err != nil {
				fmt.Printf("%v\r\n", err)
			}
		} else {
			fmt.Print("Unknown command\r\n")
		}

		if len(*history) == 0 || (*history)[len(*history)-1] != dirtyInput {
			*history = append(*history, dirtyInput)
		}

		if command == "exit" {
			return true
		}

		*currentInput = nil
		*currentInputIndex = -1
		*historyIndex = -1
		fmt.Print(prompt)
	}
	return false
}

func backInput(currentInput *[]byte, currentInputIndex *int, prompt string) {
	if len(*currentInput) > 0 && *currentInputIndex != -1 {

		*currentInput = slices.Delete(*currentInput, *currentInputIndex, *currentInputIndex+1)

		fmt.Print(ClearLine + prompt + string(*currentInput))
		*currentInputIndex--
		if *currentInputIndex != len(*currentInput)-1 {
			fmt.Printf("\x1b[%dD", len(*currentInput)-*currentInputIndex-1)
		}
	}
}
