package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var globalConfig config = config{
		Previous: nil,
		Next:     nil,
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
				fmt.Printf("Error calling %v: %w", comm.name, err)
			}
		}
	}
}
