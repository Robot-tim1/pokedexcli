package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) ListLocations(pageURL *string) (CurrentLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	return FetchData[CurrentLocations](&c.httpClient, url)
}

func FetchData[T any](pokeClient *http.Client, url string) (T, error) {
	var resultData T

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

	if err = json.NewDecoder(resp.Body).Decode(&resultData); err != nil {
		return resultData, fmt.Errorf("error decoding json: %w", err)
	}

	io.Copy(io.Discard, resp.Body)

	return resultData, nil
}
