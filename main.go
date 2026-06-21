package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
	"github.com/Robot-tim1/pokedexcli/internal/pokesave"
	"golang.org/x/term"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	pokedexData, err := pokesave.LoadPokedex()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error loading savedata: %v\n", err)
	} else if err == nil {
		cfg.pokeapiClient.Pokedex = pokedexData
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	startRepl(cfg)
}
