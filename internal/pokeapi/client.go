package pokeapi

import (
	"encoding/json"
	"fmt"
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

func requestPokeAPILocations(url string) (LocationResult, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationResult{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return LocationResult{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationResult{}, fmt.Errorf("response status code not good, indicates flop: %d", res.StatusCode)
	}

	var locations LocationResult
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locations)
	if err != nil {
		return LocationResult{}, err
	}
	return locations, err
}

func GetLocations(url string, cache *pokecache.Cache) (LocationResult, error) {
	var locations LocationResult

	//Retrieving cached data
	if data, ok := cache.Get(url); ok {
		err := json.Unmarshal(data, &locations)
		if err != nil {
			return LocationResult{}, err
		}
		return locations, nil
	}

	//Requesting locations to PokeAPI
	locations, err := requestPokeAPILocations(url)
	if err != nil {
		return LocationResult{}, nil
	}

	//Saving in cache
	encoded, err := json.Marshal(locations)
	if err != nil {
		return LocationResult{}, err
	}
	cache.Add(url, encoded)

	return locations, nil
}
