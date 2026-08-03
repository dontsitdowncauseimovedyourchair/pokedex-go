package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, comm := range getCommands() {
		fmt.Printf("%v: %v\n", comm.name, comm.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
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

func commandMapb(cfg *config, args []string) error {
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

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name1> [area_name2 ... area_name_n]")
	}

	for _, location := range args {
		url := "https://pokeapi.co/api/v2/location-area/" + location
		pokemons, err := pokeapi.GetResource[pokeapi.Pokemons](url, cfg.cache)
		if err != nil {
			return fmt.Errorf("%s: %w", location, err)
		}

		fmt.Printf("Pokemons at %s:\n", location)
		for _, encounters := range pokemons.PokemonEncounters {
			fmt.Printf(" - %s\n", encounters.Pokemon.Name)
		}
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}
	name := strings.ToLower(args[0])
	if pokemon, exists := cfg.caughtPokemon[name]; exists {
		fmt.Printf("You have already caught %s\n", pokemon.Name)
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + name
	pokemon, err := pokeapi.GetResource[pokeapi.Pokemon](url, cfg.cache)
	if err != nil {
		return err
	}

	expFactor := 154 - int(float64(pokemon.BaseExperience)/0.9)
	chance := min(95, max(20, expFactor))
	fmt.Printf("Chance of catching %s: %d%%\n", pokemon.Name, chance)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	time.Sleep(time.Second * 2)

	roll := rand.IntN(100)
	if roll < chance {
		//Caught
		cfg.caughtPokemon[name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}
	name := strings.ToLower(args[0])
	pokemon, ok := cfg.caughtPokemon[name]
	if !ok {
		fmt.Printf("You haven't caught %s yet!\n", name)
		return nil
	}

	fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\nStats:\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, typ := range pokemon.Types {
		fmt.Printf(" - %v\n", typ.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: pokedex")
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of the Pokemon you can encounter in a location",
			callback:    commandExplore,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information about your caught pokemon",
			callback:    commandInspect,
		},
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
		"pokedex": {
			name:        "pokedex",
			description: "Display the Pokedex",
			callback:    commandPokedex,
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
