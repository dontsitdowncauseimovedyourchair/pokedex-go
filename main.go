package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var globalConfig config = config{
		Previous: nil,
		Next:     nil,
		cache:    pokecache.NewCache(time.Minute * 5),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		inputCommand := words[0]
		if comm, ok := getCommands()[inputCommand]; !ok {
			fmt.Println("Unknown command")
			continue
		} else {
			err := comm.callback(&globalConfig)
			if err != nil {
				fmt.Printf("Error calling %v: %v", comm.name, err)
			}
		}
	}
}
