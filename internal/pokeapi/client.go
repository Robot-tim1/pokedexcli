package pokeapi

import (
	"net/http"
	"time"

	"github.com/Robot-tim1/pokedexcli/internal/pokecache"
)

type Client struct {
	Pokedex    map[string]Pokemon
	cache      *pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		Pokedex: make(map[string]Pokemon),
		cache:   pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
