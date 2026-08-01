package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type config struct {
	Previous *string
	Next     *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type locationResult struct {
	Count    int
	Next     *string
	Previous *string
	Results  []struct {
		Name string
		Url  string
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, comm := range getCommands() {
		fmt.Printf("%v: %v\n", comm.name, comm.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
	}

	var url string
	if cfg.Next != nil {
		url = *cfg.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("response status code not good, indicates flop: %d", res.StatusCode)
	}

	var locations locationResult
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locations)
	if err != nil {
		return err
	}

	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
