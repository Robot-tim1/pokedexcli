package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	pokeapiClient pokeapi.Client
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
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	if args != nil {
		fmt.Printf("You can't use an argument on exit idiot!\n")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationStruct, err := cfg.pokeapiClient.ListLocations(cfg.Next)
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

func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == nil {
		return errors.New("you're on the first page")
	}

	locationList, err := cfg.pokeapiClient.ListLocations(cfg.Previous)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationList.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = locationList.Next
	cfg.Previous = locationList.Previous

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if args == nil {
		return fmt.Errorf("you must provide a location name")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	location, err := cfg.pokeapiClient.ListEncounters(args[0])
	if err != nil {
		return fmt.Errorf("error getting encounters data: %w", err)
	}

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon: ")

	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if args == nil {
		return fmt.Errorf("you must provide a pokemon name")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	if poke, ok := cfg.pokeapiClient.Pokedex[args[0]]; ok {
		return fmt.Errorf("you already have %s in your pokedex", poke.Name)
	}

	pokemon, err := cfg.pokeapiClient.GetPokemon(args[0])
	if err != nil {
		return fmt.Errorf("error getting pokemon data: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	randomNum := rand.Intn(pokemon.BaseExperience)
	if randomNum <= 30 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.pokeapiClient.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}
