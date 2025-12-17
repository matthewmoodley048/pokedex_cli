package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/matthewmoodley048/pokedex_cli/internal"
)

func main() {
	cache := internal.NewCache(5 * time.Minute)

	cfg := &internal.LocationConfig{
		Cache: cache,
	}

	commands := internal.GetCommands()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		commandName := strings.ToLower(parts[0])
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command. Type 'help' for available commands.")
			continue
		}

		cfg.ExploreArgs = args
		cfg.PokemonArgs = args

		err := command.Callback(cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
