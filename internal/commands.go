package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(*LocationConfig) error
}

type LocationItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationConfig struct {
	Count       int            `json:"count"`
	Next        string         `json:"next"`
	Previous    *string        `json:"previous"`
	Results     []LocationItem `json:"results"`
	Cache       *Cache
	ExploreArgs []string
}

func GetCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			Callback:    CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			Callback:    CommandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area for pokemon",
			Callback:    CommandExplore,
		},
	}
	return commands
}

func CommandHelp(cfg *LocationConfig) error {
	fmt.Println("Welcome to the Pokedex! \nUsage:")
	commands := GetCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func CommandExit(cfg *LocationConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func fetchWithCache(url string, cache *Cache) ([]byte, error) {
	if data, found := cache.Get(url); found {
		fmt.Println("Cache hit! Using cached data.")
		return data, nil
	}

	fmt.Println("Cache miss! Fetching from API.")
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	return body, nil
}

func CommandMap(cfg *LocationConfig) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	body, err := fetchWithCache(url, cfg.Cache)
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

func CommandMapb(cfg *LocationConfig) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := *cfg.Previous

	body, err := fetchWithCache(url, cfg.Cache)
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

func CommandExplore(cfg *LocationConfig) error {
	if len(cfg.ExploreArgs) == 0 || cfg.ExploreArgs[0] == "" {
		return errors.New("please enter a location")
	}

	location := cfg.ExploreArgs[0]
	url := "https://pokeapi.co/api/v2/location-area/" + location + "/"

	body, err := fetchWithCache(url, cfg.Cache)
	if err != nil {
		return err
	}

	var exploreData struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(body, &exploreData); err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreData.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
