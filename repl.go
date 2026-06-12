package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) {
	getCommand := commandRegistry
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !reader.Scan() {
			break
		}

		input := cleanInput(reader.Text())

		if len(input) == 0 {
			continue
		}

		command := input[0]

		if commandStruct, ok := getCommand[command]; ok {
			err := commandStruct.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
	if err := reader.Err(); err != nil {
		fmt.Printf("Error encountered: %v", err)
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	cleanText := strings.Fields(lowerText)
	return cleanText
}
