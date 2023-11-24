package pokedexcli

import (
	"fmt"
	"os"
	"time"

	"github.com/this-is-richard/pokedex/pokeapi"
)

type Command struct {
	Key         string
	Description string
	Callback    func()
}

func getCommands(pokedex *Pokedex, args ...string) map[string]Command {
	return map[string]Command{
		"exit": {
			Key:         "exit",
			Description: "Exit the program gracefully.",
			Callback: func() {
				fmt.Println("Bye.")
				os.Exit(0)
			},
		},
		"help": {
			Key:         "help",
			Description: "See available commands.",
			Callback: func() {
				fmt.Printf("\nWelcome to pokedex. Available commands:\n\n")
				allCommands := getCommands(pokedex)
				for key := range allCommands {
					command := allCommands[key]
					fmt.Printf("%v: %v\n", command.Key, command.Description)
				}
			},
		},
		"map": {
			Key:         "map",
			Description: "See the next 20 locations",
			Callback: func() {
				pokedex.offset += pokedex.limit

				locations, err := pokedex.pokeapiClient.GetLocations(pokedex.offset, pokedex.limit)
				if err != nil {
					fmt.Printf("failed to list locations, %v\n", err.Error())
					return
				}

				for _, location := range locations.Results {
					fmt.Printf("%v: %v\n", location.Name, location.URL)
				}
			},
		},
		"mapb": {
			Key:         "mapb",
			Description: "See the previous 20 locations",
			Callback: func() {
				if pokedex.offset-pokedex.limit > 0 {
					pokedex.offset -= pokedex.limit
				}

				locations, err := pokedex.pokeapiClient.GetLocations(pokedex.offset, pokedex.limit)
				if err != nil {
					fmt.Printf("failed to list locations, %v\n", err.Error())
					return
				}

				for _, location := range locations.Results {
					fmt.Printf("%v: %v\n", location.Name, location.URL)
				}
			},
		},
		"explore": {
			Key:         "explore <area_code>",
			Description: "See pokemons in a certain area",
			Callback: func() {
				if len(args) == 0 {
					fmt.Println("Please enter area code, for example, `explore sunyshore-city-area`")
					return
				}

				areaCode := args[0]
				locationArea, err := pokedex.pokeapiClient.GetLocationArea(areaCode)
				if err != nil {
					fmt.Printf("failed to explore area %v: %v\n", areaCode, err)
					return
				}

				for _, pe := range locationArea.PokemonEncounters {
					fmt.Printf("%v: %v\n", pe.Pokemon.Name, pe.Pokemon.URL)
				}
			},
		},
		"catch": {
			Key:         "catch <pokemon_name>",
			Description: "Catch Pokemon",
			Callback: func() {
				if len(args) == 0 {
					fmt.Println("Please enter pokemon_name, for example, `catch pikachu`")
					return
				}

				name := args[0]
				pokemon, err := pokedex.pokeapiClient.GetPokemon(name)
				if err != nil {
					fmt.Printf("failed to find pokemeon %v: %v\n", name, err)
					return
				}

				newPokemon := NewPokemon(*pokemon)
				fmt.Printf("catching %v...\n", newPokemon.Name)
				time.Sleep(1200 * time.Millisecond)
				caught := newPokemon.Catch()
				if caught {
					pokedex.pokemons[pokemon.Name] = newPokemon
					fmt.Printf("%v caught!\n", newPokemon.Name)
					fmt.Printf("you now have %v pokemons!\n", len(pokedex.pokemons))
				} else {
					fmt.Printf("%v escaped! try next time!\n", newPokemon.Name)
				}
			},
		},
		"inspect": {
			Key:         "inspect <pokemon_name>",
			Description: "Inspect a Pokemon you've caught",
			Callback: func() {
				if len(args) == 0 {
					fmt.Println("Please enter pokemon_name, for example, `catch pikachu`")
					return
				}

				name := args[0]
				pokemon, ok := pokedex.pokemons[name]
				if !ok {
					fmt.Printf("%v not caught yet, cannot inspect it\n", name)
					return
				}

				fmt.Printf("Name: %v\n", pokemon.Name)
				fmt.Printf("Height: %v\n", pokemon.Height)
				fmt.Printf("Weight: %v\n", pokemon.Weight)
				fmt.Println("Stats:")
				for _, stat := range pokemon.Stats {
					fmt.Printf("- %v: %v\n", stat.Stat.Name, stat.BaseStat)
				}
				fmt.Println("Types:")
				for _, t := range pokemon.Types {
					fmt.Printf("- %v\n", t.Type.Name)
				}
				fmt.Println()
			},
		},
		"pokedex": {
			Key:         "pokedex",
			Description: "List all Pokemons in your Pokedex",
			Callback: func() {
				if len(pokedex.pokemons) == 0 {
					fmt.Println("Your Pokedex is empty. Catch a Pokemon by `catch <pokemon_name>`")
					return
				}
				fmt.Println("In your Pokedex:")
				for name := range pokedex.pokemons {
					fmt.Printf("- %v\n", name)
				}
			},
		},
	}
}

type Pokedex struct {
	offset        int
	limit         int
	pokeapiClient *pokeapi.PokeapiClient
	pokemons      map[string]Pokemon
}

func NewPokedex() Pokedex {
	return Pokedex{limit: 20, pokeapiClient: pokeapi.NewPokeapiClient(), pokemons: make(map[string]Pokemon)}
}

func (p *Pokedex) Run(key string, args ...string) error {
	command, ok := getCommands(p, args...)[key]
	if !ok {
		return fmt.Errorf("unknown command `%v`, type `help` to see available commands", key)
	}

	command.Callback()
	return nil
}
