package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CurrentMap struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var pokeClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchMap(url string) (CurrentMap, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CurrentMap{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("user-agent", "pokedexcli for boot.dev course")

	resp, err := pokeClient.Do(req)
	if err != nil {
		return CurrentMap{}, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return CurrentMap{}, fmt.Errorf("returned status code: %d", resp.StatusCode)
	}

	var mapData CurrentMap
	if err = json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
		return CurrentMap{}, fmt.Errorf("error decoding json: %w", err)
	}

	return mapData, nil
}
