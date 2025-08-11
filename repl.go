package main

import (
	"bufio"
	"fmt"
	"github.com/ntpotraz/pokedex/internal/pokecache"
	"os"
	"strings"
	"time"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  *Pokedex
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
	config      *Config
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(time.Minute * 5)

	cfg := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
		Cache:    cache,
		Pokedex:  NewPokedex(),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		commandName, args := cleanInput(userInput)
		if commandName == "" {
			continue
		}

		command, exists := getCommands(cfg)[commandName]

		if exists {
			if err := command.callback(cfg, args); err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Printf("Unknown command: %s\n", commandName)
			continue
		}
	}
}

func cleanInput(text string) (string, []string) {
	if text == "" {
		return "", []string{}
	}

	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	if len(words) == 1 {
		return words[0], []string{}
	}
	return words[0], words[1:]
}

func getCommands(cfg *Config) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      cfg,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      cfg,
		},
		"map": {
			name:        "map",
			description: "Gets a list of locations",
			callback:    commandMap,
			config:      cfg,
		},
		"mapb": {
			name:        "mapb",
			description: "Goes back to previous list of locations",
			callback:    commandMapb,
			config:      cfg,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location on the map",
			callback:    commandExplore,
			config:      cfg,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified pokemon",
			callback:    commandCatch,
			config:      cfg,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon in your Pokedex",
			callback:    commandInspect,
			config:      cfg,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon in your pokedex",
			callback:    commandPokedex,
			config:      cfg,
		},
	}
}
