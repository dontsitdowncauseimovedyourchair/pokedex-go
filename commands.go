package main

import (
	"fmt"
	"os"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokeapi"
	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokecache"
)

type config struct {
	Previous *string
	Next     *string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	var url string
	if cfg.Next != nil {
		url = *cfg.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}

	locations, err := pokeapi.GetResource[pokeapi.LocationResult](url, cfg.cache)
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

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := *cfg.Previous

	locations, err := pokeapi.GetResource[pokeapi.LocationResult](url, cfg.cache)
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
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of the previous 20 previous location areas in the Pokemon world",
			callback:    commandMapb,
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
