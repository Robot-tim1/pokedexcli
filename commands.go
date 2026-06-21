package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
	"github.com/Robot-tim1/pokedexcli/internal/pokesave"
)

type config struct {
	pokeapiClient pokeapi.Client
	Next          *string
	Previous      *string
}

func commandExit(cfg *config, args ...string) error {
	if len(cfg.pokeapiClient.Pokedex) != 0 {
		err := pokesave.SavePokedex(cfg.pokeapiClient.Pokedex)
		if err != nil {
			return fmt.Errorf("error saving data: %w", err)
		}
	}
	fmt.Print("Closing the Pokedex... Goodbye!\r\n")
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\r\n")
	fmt.Print("Usage:\r\n\n")
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\r\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationList, err := cfg.pokeapiClient.ListLocations(cfg.Next)
	if err != nil {
		return fmt.Errorf("error getting map data: %w", err)
	}

	for _, location := range locationList.Results {
		fmt.Printf("%s\r\n", location.Name)
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
		fmt.Printf("%s\r\n", location.Name)
	}

	cfg.Next = locationList.Next
	cfg.Previous = locationList.Previous

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	location, err := cfg.pokeapiClient.ListEncounters(args[0])
	if err != nil {
		return fmt.Errorf("error getting encounters data: %w", err)
	}

	fmt.Printf("Exploring %s...\r\n", location.Name)
	fmt.Print("Found Pokemon: \r\n")

	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" - %s\r\n", enc.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	if poke, ok := cfg.pokeapiClient.GetPokedex(args[0]); ok {
		return fmt.Errorf("you already have %s in your pokedex", poke.Name)
	}

	pokemon, err := cfg.pokeapiClient.GetPokemon(args[0])
	if err != nil {
		return fmt.Errorf("error getting pokemon data: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\r\n", pokemon.Name)
	randomNum := rand.Intn(pokemon.BaseExperience)
	if randomNum <= 50 {
		fmt.Printf("%s was caught!\r\n", pokemon.Name)
		fmt.Print("You may now inspect it with the inspect command.\r\n")
		cfg.pokeapiClient.SetPokedex(pokemon.Name, pokemon)
		cfg.pokeapiClient.Cache.Delete("https://pokeapi.co/api/v2/pokemon/" + pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\r\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	p, ok := cfg.pokeapiClient.GetPokedex(args[0])
	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\r\nHeight: %d\r\nWeight: %d\r\n", p.Name, p.Height, p.Weight)

	fmt.Printf("Stats:\r\n")
	for i := range p.Stats {
		fmt.Printf("  -%s: %d\r\n", p.Stats[i].Stat.Name, p.Stats[i].BaseStat)
	}

	fmt.Printf("Types:\r\n")
	for i := range p.Types {
		fmt.Printf("  - %s\r\n", p.Types[i].Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	cfg.pokeapiClient.DexMu.Lock()
	defer cfg.pokeapiClient.DexMu.Unlock()
	for name := range cfg.pokeapiClient.Pokedex {
		fmt.Printf(" - %s\r\n", name)
	}
	return nil
}

func commandDelete(cfg *config, args ...string) error {
	cfg.pokeapiClient.Pokedex = make(map[string]pokeapi.Pokemon)
	err := pokesave.DeletePokedex()
	if err != nil {
		return fmt.Errorf("error deleting savedata: %w", err)
	}
	return nil
}
