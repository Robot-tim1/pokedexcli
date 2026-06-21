package pokesave

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Robot-tim1/pokedexcli/internal/pokeapi"
)

func SavePokedex(pokedex map[string]pokeapi.Pokemon) error {
	filename, err := getFileNamePath()
	if err != nil {
		return fmt.Errorf("error getting filepath: %w", err)
	}

	data, err := json.Marshal(pokedex)
	if err != nil {
		return fmt.Errorf("error marshalling savedata: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadPokedex() (map[string]pokeapi.Pokemon, error) {
	filename, err := getFileNamePath()
	if err != nil {
		return nil, fmt.Errorf("error getting filepath: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading savedata: %w", err)
	}
	var pokedex map[string]pokeapi.Pokemon
	err = json.Unmarshal(data, &pokedex)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling savedata: %w", err)
	}
	return pokedex, err
}

func DeletePokedex() error {
	filename, err := getFileNamePath()
	if err != nil {
		return fmt.Errorf("error getting filepath: %w", err)
	}

	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting savedata: %w", err)
	}
	return nil
}

func getFileNamePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	var filename string
	switch os := runtime.GOOS; os {
	case "linux":
		filename = filepath.Join(home, ".local", "share", "pokedexcli", "pokedex.json")
	case "windows", "darwin":
		fallthrough
	default:
		filename = filepath.Join(home, "Documents", "pokedexcli", "pokedex.json")
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("error creating directory savedata: %w", err)
	}
	return filename, nil
}
