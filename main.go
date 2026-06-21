package main

import (
	"os"
	"time"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
	"golang.org/x/term"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	startRepl(cfg)
}
