package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cheezecakee/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n\n")
	fmt.Printf("Help: Displays a help message\n")
	fmt.Printf("Exit: Exit the Pokedex\n")
	fmt.Printf("Map: Display the next 20 location in the Pokemon world\n")
	fmt.Printf("Mapb: Display the previous 20 location in the Pokemon world\n\n")
	return nil
}

func commandExit() error {
	fmt.Printf("Exiting...")
	os.Exit(0)
	return nil
}

type config struct {
	Next     string
	Previous string
}

func newConfig() *config {
	return &config{}
}

func (c *config) commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if c.Next != "" {
		url = c.Next
	}
	locations, err := internal.GetLocations(url)
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

func (c *config) commandMapb() error {
	if c.Previous == "" {
		fmt.Printf("No previous locations to display, try a different typing [map]\n")
		return nil
	}
	locations, err := internal.GetLocations(c.Previous)
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommandMap := createCliCommand()
	for {
		fmt.Printf("Pokedex > ")
		for scanner.Scan() {
			command := scanner.Text()
			switch command {
			case "help":
				cliCommandMap[command].callback()
			case "exit":
				cliCommandMap[command].callback()
			case "map":
				cliCommandMap[command].callback()
			case "mapb":
				cliCommandMap[command].callback()
			}
			break
		}
	}
}
