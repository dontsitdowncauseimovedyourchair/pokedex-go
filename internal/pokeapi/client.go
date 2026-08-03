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

type Pokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height  int    `json:"height"`
	Name    string `json:"name"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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
		if res.StatusCode == 404 {
			return zero, fmt.Errorf("404 not found")
		}
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

	// Cache the data, I first Marshal the data so that the unneeded fields are not stored
	toCache, err := json.Marshal(resource)
	if err != nil {
		return zero, err
	}
	cache.Add(url, toCache)

	return resource, nil
}
