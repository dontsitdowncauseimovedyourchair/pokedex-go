package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
			err := comm.callback()
			if err != nil {
				fmt.Printf("Error calling %v: %w", comm.name, err)
			}
		}
	}
}
