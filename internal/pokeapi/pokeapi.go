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

type CurrentLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocations(pageURL *string, cache *pokecache.Cache) (CurrentLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	return FetchData[CurrentLocations](&c.httpClient, url, cache)
}

func FetchData[T any](pokeClient *http.Client, url string, cache *pokecache.Cache) (T, error) {
	var resultData T
	var value []byte

	cache.Mu.Lock()
	if entry, ok := cache.Entries[url]; ok {
		value = entry.Val
	}
	cache.Mu.Unlock()

	if value != nil {
		if err := json.Unmarshal(value, &resultData); err != nil {
			return resultData, fmt.Errorf("error decoding json: %w", err)
		}
		return resultData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resultData, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("user-agent", "pokedexcli for boot.dev course")

	resp, err := pokeClient.Do(req)
	if err != nil {
		return resultData, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return resultData, fmt.Errorf("returned status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resultData, fmt.Errorf("error reading response body: %w", err)
	}

	if err = json.Unmarshal(bodyBytes, &resultData); err != nil {
		return resultData, fmt.Errorf("error decoding json: %w", err)
	}

	cache.Add(url, bodyBytes)

	return resultData, nil
}
