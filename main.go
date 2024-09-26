package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/cheezecakee/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

var cliCommandMap map[string]cliCommand

func createCliCommand() map[string]cliCommand {
	config := newConfig()
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location in the Pokemon world",
			callback:    config.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location in the Pokemon world",
			callback:    config.commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Displays pokemons in the area",
			callback:    config.commmandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Trys to catch a pokemon",
			callback:    config.commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Displays pokemon stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all caught pokemons",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(args []string) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n\n")

	for _, cmd := range cliCommandMap {
		// Adjust the width in the format string to suit your needs
		fmt.Printf("%-20s : %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(args []string) error {
	fmt.Printf("Exiting...")
	os.Exit(0)
	return nil
}

type config struct {
	Next     string
	Previous string
	cache    *internal.Cache
}

func newConfig() *config {
	cache := internal.NewCache(1 * time.Minute)
	return &config{cache: cache}
}

func (c *config) commandMap(args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if c.Next != "" {
		url = c.Next
	}
	locations, err := internal.GetLocations(url, c.cache)
	c.Next = locations.Next
	c.Previous = locations.Previous
	if err != nil {
		return err
	}
	for _, i := range locations.Results {
		fmt.Printf("%v\n", i.Name)
	}
	return nil
}

func (c *config) commandMapb(args []string) error {
	if c.Previous == "" {
		fmt.Printf("No previous locations to display, try a different typing [map]\n")
		return nil
	}
	locations, err := internal.GetLocations(c.Previous, c.cache)
	c.Next = locations.Next
	c.Previous = locations.Previous
	if err != nil {
		return err
	}
	for _, i := range locations.Results {
		fmt.Printf("%v\n", i.Name)
	}
	return nil
}

func (c *config) commmandExplore(args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	var name string

	if len(args) == 0 {
		return fmt.Errorf("Insert area name [explore <area_name>]")
	}
	name = args[0]

	pokemon, err := internal.GetPokemonsInArea(url, name, c.cache)
	if err != nil {
		return fmt.Errorf("Invalid area name %v", err)
	}

	fmt.Printf("Exploring %v...\n", name)
	fmt.Printf("Found Pokemon:\n")
	for _, i := range pokemon.PokemonEncounters {
		fmt.Printf("- %v\n", i.Pokemon.Name)
	}

	return nil
}

var PokeDex = make(map[string]internal.PokemonDetails)

func (c *config) commandCatch(args []string) error {
	url := "https://pokeapi.co/api/v2/pokemon/"
	var name string

	if len(args) == 0 {
		return fmt.Errorf("Insert pokemon name [catch <pokemon>]")
	}
	name = args[0]

	pokemon, err := internal.GetPokemon(url, name, c.cache)
	if err != nil {
		fmt.Printf("Invalid pokemon name %v\n", err)
		return err
	}
	fmt.Printf("Throwing a pokeball at %+v...\n", pokemon.Name)

	catchThreshold := 100 - pokemon.BaseExperience/10 // Smaller value -> harder catch

	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	min, max := 0, 100
	randomInRange := rng.Intn(max-min+1) + min

	if randomInRange <= catchThreshold {
		fmt.Printf("Oh no! %v got away!\n", pokemon.Name)
		fmt.Printf("Better luck next time!\n")
		return nil
	}

	fmt.Printf("Congratulations! You caught %v!\n", pokemon.Name)
	fmt.Printf("Pokedex: [%v]\n", pokemon.Name)

	PokeDex[pokemon.Name] = pokemon

	// fmt.Printf("Pokedex: %v\n", PokeDex)
	return nil
}

// Constants for the layout
const (
	asciiHeight  = 10 // Height of the ASCII art (blank space for now)
	asciiWidth   = 30 // Width of the ASCII art (blank space for now)
	consoleWidth = 80 // Total width of the output
)

// Function to print blank ASCII space (or placeholder)
func printAsciiPlaceholder() {
	for i := 0; i < asciiHeight; i++ {
		fmt.Printf("%s", strings.Repeat(" ", asciiWidth))
		fmt.Println() // move to the next line
	}
}

func commandInspect(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Insert area name [explore <area_name>]")
	}
	name := args[0]

	if name, exists := PokeDex[name]; !exists {
		fmt.Printf("Haven't caught %v yet! No data available.\n", name)
		return nil
	}

	// Print the ASCII placeholder on the left
	printAsciiPlaceholder()

	// Print the box with the same formatting as before, pushed to the right
	fmt.Println(strings.Repeat("-", consoleWidth))
	fmt.Printf("%50sPokedex%v\n", "", PokeDex[name].Name)
	fmt.Println(strings.Repeat("-", consoleWidth))

	// Pokémon details
	fmt.Printf("%30sName: %-10v\n", "", PokeDex[name].Name)
	fmt.Printf("%30sHeight: %-10v\n", "", PokeDex[name].Height)
	fmt.Printf("%30sWeight: %-10v\n", "", PokeDex[name].Weight)

	// Stats
	fmt.Printf("%30sStats:\n", "")
	for _, i := range PokeDex[name].Stats {
		fmt.Printf("%30s  - %v: %-10v\n", "", i.Stat.Name, i.BaseStat)
	}

	// Types
	fmt.Printf("%30sType:\n", "")
	for _, i := range PokeDex[name].Types {
		fmt.Printf("%30s  - %-10v\n", "", i.Type.Name)
	}

	fmt.Println(strings.Repeat("-", consoleWidth))
	return nil
}

func commandPokedex(args []string) error {
	fmt.Printf("Your Pokedex: \n")
	for pokemon := range PokeDex {
		fmt.Printf("  - %v\n", PokeDex[pokemon].Name)
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommandMap = createCliCommand()
	for {
		fmt.Printf("Pokedex > ")
		for scanner.Scan() {
			input := scanner.Text()

			parts := strings.Fields(input)
			if len(parts) == 0 {
				continue
			}

			command := parts[0]
			args := parts[1:]

			if cmd, ok := cliCommandMap[command]; ok {
				cmd.callback(args)
			} else {
				fmt.Printf("Unknown command: %s\n", command)
				break
			}
			break
		}
	}
}
