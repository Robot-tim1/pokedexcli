package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
	"github.com/Robot-tim1/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	pokeapiClient pokeapi.Client
	cache         *pokecache.Cache
	Next          *string
	Previous      *string
}

var commandRegistry map[string]cliCommand

func init() {
	commandRegistry = map[string]cliCommand{
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	locationStruct, err := cfg.pokeapiClient.ListLocations(cfg.Next, cfg.cache)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationStruct.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = locationStruct.Next
	cfg.Previous = locationStruct.Previous

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		return errors.New("you're on the first page")
	}

	locationStruct, err := cfg.pokeapiClient.ListLocations(cfg.Previous, cfg.cache)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationStruct.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = locationStruct.Next
	cfg.Previous = locationStruct.Previous

	return nil
}
