package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/this-is-richard/pokedex/cache"
)

const LOCATION_URL = "https://pokeapi.co/api/v2/location"
const LOCATION_AREA_URL = "https://pokeapi.co/api/v2/location-area"
const POKEMON_URL = "https://pokeapi.co/api/v2/pokemon"

type PokeapiClient struct {
	cache *cache.Cache[[]byte]
}

func NewPokeapiClient() *PokeapiClient {
	return &PokeapiClient{
		cache: cache.NewCache[[]byte](10 * time.Second),
	}
}

func (c *PokeapiClient) GetLocations(offset int, limit int) (*Locations, error) {
	url := fmt.Sprintf("%v?offset=%v&limit=%v", LOCATION_URL, offset, limit)
	var locations Locations

	bytes, found := c.cache.Get(url)
	if err := json.Unmarshal(bytes, &locations); found && err == nil {
		return &locations, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %v", err.Error())
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("fetched locations with non-successful status code %v", resp.StatusCode)

	}

	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp.Body: %v", err.Error())
	}

	c.cache.Add(url, bytes)

	err = json.Unmarshal(bytes, &locations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resp.Body: %v", err.Error())
	}

	return &locations, nil
}

func (c *PokeapiClient) GetLocationArea(areaCode string) (*LocationArea, error) {
	url := fmt.Sprintf("%v/%v", LOCATION_AREA_URL, areaCode)
	var locationArea LocationArea

	bytes, found := c.cache.Get(url)
	if err := json.Unmarshal(bytes, &locationArea); found && err == nil {
		return &locationArea, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location-area: %v", err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("location-area %v not found", areaCode)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("fetched location-area with non-successful status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp.Body: %v", err.Error())
	}

	c.cache.Add(url, bytes)

	err = json.Unmarshal(bytes, &locationArea)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resp.Body: %v", err.Error())
	}

	return &locationArea, nil
}

func (c *PokeapiClient) GetPokemon(name string) (*Pokemon, error) {
	url := fmt.Sprintf("%v/%v", POKEMON_URL, name)
	var pokemon Pokemon

	bytes, found := c.cache.Get(url)
	if err := json.Unmarshal(bytes, &pokemon); found && err == nil {
		return &pokemon, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon: %v", err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon %v not found", name)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("fetched pokemon with non-successful status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp.Body: %v", err.Error())
	}

	c.cache.Add(url, bytes)

	err = json.Unmarshal(bytes, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resp.Body: %v", err.Error())
	}

	return &pokemon, nil
}
