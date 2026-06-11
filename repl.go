package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	getCommand := commandRegistry()
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		reader.Scan()
		input := cleanInput(reader.Text())

		if len(input) == 0 {
			continue
		}

		command := input[0]

		if commandStruct, ok := getCommand[command]; ok {
			commandStruct.callback()
		} else {
			fmt.Println("Unknown command")
		}

		if err := reader.Err(); err != nil {
			fmt.Printf("Error encountered: %v", err)
			break
		}
	}

}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	cleanText := strings.Fields(lowerText)
	return cleanText
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	allCommands := commandRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, command := range allCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
