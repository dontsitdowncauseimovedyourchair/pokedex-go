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
			var args []string = words[1:]
			err := comm.callback(&globalConfig, args)
			if err != nil {
				fmt.Printf("Error calling %v: %v\n", comm.name, err)
			}
		}
	}
}
