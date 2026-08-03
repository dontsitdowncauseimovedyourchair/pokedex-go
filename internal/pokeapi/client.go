package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokecache"
)

type LocationResult struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Location struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetResource[T any](url string, cache *pokecache.Cache) (T, error) {
	var resource T
	var zero T
	//Retrieving cached data
	if data, ok := cache.Get(url); ok {
		err := json.Unmarshal(data, &resource)
		if err != nil {
			return zero, err
		}
		return resource, nil
	}

	//Requesting resource to PokeAPI
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, err
	}

	res, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return zero, fmt.Errorf("response status code not good, indicates flop: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return zero, err
	}

	// I will only cache if the data was actually unmarshal-able
	cache.Add(url, data)
	return resource, nil
}
