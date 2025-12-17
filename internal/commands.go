package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
	ExploreArgs []string
	PokemonArgs []string
	Pokedex     map[string]Pokemon
	Cache       *Cache
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order int `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
	} `json:"past_abilities"`
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
		"catch": {
			name:        "catch",
			description: "Try and catch a pokemon",
			Callback:    CommandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View a caught pokemon",
			Callback:    CommandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all caught pokemon",
			Callback:    CommandPokedex,
		}, // Display pokemon information
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
		return data, nil
	}

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

func CommandCatch(cfg *LocationConfig) error {
	if len(cfg.PokemonArgs) == 0 || cfg.PokemonArgs[0] == "" {
		return errors.New("if you aim at nothing you will catch nothing -Prof Oak")
	}

	if cfg.Pokedex == nil {
		cfg.Pokedex = make(map[string]Pokemon)
	}

	fmt.Printf("\nThrowing a Pokeball at %s...\n", cfg.PokemonArgs[0])

	pokemon := cfg.PokemonArgs[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"

	body, err := fetchWithCache(url, cfg.Cache)
	if err != nil {
		return err
	}

	var pokemonData Pokemon

	if err := json.Unmarshal(body, &pokemonData); err != nil {
		return err
	}

	caughtPokemon := rand.Intn(pokemonData.BaseExperience*5) <= pokemonData.BaseExperience

	if caughtPokemon {
		cfg.Pokedex[pokemonData.Name] = pokemonData
		fmt.Printf("You caught %s! \n", pokemonData.Name)
	} else {
		fmt.Printf("%s ran away :( \n", pokemonData.Name)
	}

	return nil
}

func CommandInspect(cfg *LocationConfig) error {
	if len(cfg.PokemonArgs) == 0 {
		return errors.New("if you search for nothing that is what you will find -Brok Lesnar")
	}

	pokemonName := cfg.PokemonArgs[0]

	pokemon, found := cfg.Pokedex[pokemonName]
	if !found {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func CommandPokedex(cfg *LocationConfig) error {
	if len(cfg.Pokedex) <= 0 || cfg.Pokedex == nil {
		fmt.Print("nothing caught yet!\n")
		return nil
	}
	fmt.Print("your pokedex:\n")
	for pokemon := range cfg.Pokedex {
		fmt.Printf("-%s \n", pokemon)
	}

	return nil
}
