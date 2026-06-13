package pokeapi

import (
	"net/http"
	"sync"
	"time"

	"github.com/Robot-tim1/pokedexcli/internal/pokecache"
)

type Client struct {
	dexMu      *sync.Mutex
	Pokedex    map[string]Pokemon
	cache      *pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		dexMu:   &sync.Mutex{},
		Pokedex: make(map[string]Pokemon),
		cache:   pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetPokedex(pokemon string) (Pokemon, bool) {
	c.dexMu.Lock()
	defer c.dexMu.Unlock()
	poke, ok := c.Pokedex[pokemon]
	return poke, ok
}

func (c *Client) SetPokedex(key string, pokemon Pokemon) {
	c.dexMu.Lock()
	defer c.dexMu.Unlock()
	c.Pokedex[key] = pokemon
}
