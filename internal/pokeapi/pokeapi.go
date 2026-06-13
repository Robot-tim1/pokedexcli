package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Robot-tim1/pokedexcli/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

func (c *Client) ListLocations(pageURL *string) (CurrentLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	return FetchData[CurrentLocations](&c.httpClient, url, c.cache)
}

func (c *Client) ListEncounters(areaName string) (PokemonEncounters, error) {
	url := baseURL + "/location-area/" + areaName

	encounters, err := FetchData[PokemonEncounters](&c.httpClient, url, c.cache)
	if err != nil {
		return PokemonEncounters{}, fmt.Errorf("error getting area data: %w", err)
	}

	return encounters, nil
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	pokemon, err := FetchData[Pokemon](&c.httpClient, url, c.cache)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error getting pokemon data: %w", err)
	}

	return pokemon, nil
}

func FetchData[T any](pokeClient *http.Client, url string, cache *pokecache.Cache) (T, error) {
	var resultData T
	var zero T

	value, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(value, &resultData); err != nil {
			return zero, fmt.Errorf("error decoding json from cache: %w", err)
		}
		return resultData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "pokedexcli for boot.dev course")

	resp, err := pokeClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return zero, fmt.Errorf("returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return zero, fmt.Errorf("error reading response body: %w", err)
	}

	if err = json.Unmarshal(bodyBytes, &resultData); err != nil {
		return zero, fmt.Errorf("error decoding json: %w", err)
	}

	cache.Add(url, bodyBytes)

	return resultData, nil
}
