package main

import (
	"fmt"
	"os"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	Next     string
	Previous string
}

func commandRegistry() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows previous 20 locations",
			callback:    commandMapb,
		},
	}
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

func commandMap() error {
	replaceLater := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	if replaceLater.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	locationStruct, err := pokeapi.FetchMap(replaceLater.Next)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationStruct.Results {
		fmt.Println(location.Name)
	}

	if locationStruct.Next != nil {
		replaceLater.Next = *locationStruct.Next
	} else {
		replaceLater.Next = ""
	}

	if locationStruct.Previous != nil {
		replaceLater.Previous = *locationStruct.Previous
	} else {
		replaceLater.Previous = ""
	}

	return nil
}

func commandMapb() error {
	replaceLater := config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
		Previous: "",
	}

	if replaceLater.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationStruct, err := pokeapi.FetchMap(replaceLater.Previous)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationStruct.Results {
		fmt.Println(location.Name)
	}

	if locationStruct.Next != nil {
		replaceLater.Next = *locationStruct.Next
	} else {
		replaceLater.Next = ""
	}

	if locationStruct.Previous != nil {
		replaceLater.Previous = *locationStruct.Previous
	} else {
		replaceLater.Previous = ""
	}

	return nil
}
