package main

import (
	"bufio"
	"fmt"
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

	pokemon, err := internal.GetPokemons(url, name, c.cache)
	if err != nil {
		return fmt.Errorf("Invalid area name %v", err)
	}

	fmt.Printf("Found Pokemon:\n")
	for _, i := range pokemon.PokemonEncounters {
		fmt.Printf("- %v\n", i.Pokemon.Name)
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
