package main
import (
	"fmt"
	"os"
	"io"
	"net/http"
  "encoding/json"
)


type cliCommand struct {
	name        string
	description string
	callback    func(*LocationConfig) error
}

type LocationItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationConfig struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous *string        `json:"previous"` // pointer allows null
	Results  []LocationItem `json:"results"`
}



func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}

	return commands
}

func commandHelp(cfg *LocationConfig) error {
	fmt.Println("Welcome to the Pokedex! \nUsage:\n")
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit(cfg *LocationConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *LocationConfig) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var newCfg LocationConfig
	if err := json.Unmarshal(body, &newCfg); err != nil {
		return err
	}

	for _, loc := range newCfg.Results {
		fmt.Println(loc.Name)
	}

	cfg.Count = newCfg.Count
	cfg.Next = newCfg.Next
	cfg.Previous = newCfg.Previous
	cfg.Results = newCfg.Results

	return nil
}

func commandMapb(cfg *LocationConfig) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := *cfg.Previous

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var newCfg LocationConfig
	if err := json.Unmarshal(body, &newCfg); err != nil {
		return err
	}

	for _, loc := range newCfg.Results {
		fmt.Println(loc.Name)
	}

	cfg.Count = newCfg.Count
	cfg.Next = newCfg.Next
	cfg.Previous = newCfg.Previous
	cfg.Results = newCfg.Results

	return nil
}
