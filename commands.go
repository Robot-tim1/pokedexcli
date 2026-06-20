package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient pokeapi.Client
	Next          *string
	Previous      *string
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
	locationList, err := cfg.pokeapiClient.ListLocations(cfg.Next)
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
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location name")
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
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	if poke, ok := cfg.pokeapiClient.GetPokedex(args[0]); ok {
		return fmt.Errorf("you already have %s in your pokedex", poke.Name)
	}

	pokemon, err := cfg.pokeapiClient.GetPokemon(args[0])
	if err != nil {
		return fmt.Errorf("error getting pokemon data: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	randomNum := rand.Intn(pokemon.BaseExperience)
	if randomNum <= 50 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Printf("You may now inspect it with the inspect command.\n")
		cfg.pokeapiClient.SetPokedex(pokemon.Name, pokemon)
		cfg.pokeapiClient.Cache.Delete("https://pokeapi.co/api/v2/pokemon/" + pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	p, ok := cfg.pokeapiClient.GetPokedex(args[0])
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", p.Name, p.Height, p.Weight)

	fmt.Printf("Stats:\n")
	for i := range p.Stats {
		fmt.Printf("  -%s: %d\n", p.Stats[i].Stat.Name, p.Stats[i].BaseStat)
	}

	fmt.Printf("Types:\n")
	for i := range p.Types {
		fmt.Printf("  - %s\n", p.Types[i].Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	cfg.pokeapiClient.DexMu.Lock()
	defer cfg.pokeapiClient.DexMu.Unlock()
	for name := range cfg.pokeapiClient.Pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
