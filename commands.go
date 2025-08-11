package main

import (
	"fmt"
	"github.com/ntpotraz/pokedex/internal/pokeapi"
	"math/rand"
	"os"
)

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args []string) error {
	commands := getCommands(cfg)
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *Config, args []string) error {
	if cfg.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	pokeMap, err := pokeapi.CallPokeApi[MapJSON](cfg.Next, cfg.Cache)
	if err != nil {
		return fmt.Errorf("error calling poke api: %w", err)
	}

	for _, item := range pokeMap.Results {
		fmt.Println(item.Name)
	}

	cfg.Next = pokeMap.Next
	cfg.Previous = pokeMap.Previous

	return nil
}

func commandMapb(cfg *Config, args []string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	pokeMap, err := pokeapi.CallPokeApi[MapJSON](cfg.Previous, cfg.Cache)
	if err != nil {
		return fmt.Errorf("error calling poke api: %w", err)
	}

	for _, item := range pokeMap.Results {
		fmt.Println(item.Name)
	}

	cfg.Next = pokeMap.Next
	cfg.Previous = pokeMap.Previous

	return nil
}

func commandExplore(cfg *Config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("error invalid number of arguments")
	}

	loc := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + loc
	pokeLoc, err := pokeapi.CallPokeApi[LocationJSON](url, cfg.Cache)
	if err != nil {
		return fmt.Errorf("error retrieving location info: %w", err)
	}
	fmt.Println("Exploring", loc+"...")
	fmt.Println("Found Pokemon:")
	for _, item := range pokeLoc.PokemonEncounters {
		fmt.Printf("- %s\n", item.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("error invalid number of arguments")
	}

	pokemonName := args[0]

	if ok := cfg.Pokedex.Check(pokemonName); ok {
		fmt.Printf("%s has already been added to the Pokedex\n", pokemonName)
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	pokePoke, err := pokeapi.CallPokeApi[PokemonJSON](url, cfg.Cache)
	if err != nil {
		return fmt.Errorf("error retrieving pokemon info: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokePoke.Name)

	minXP := 30.0
	maxXP := 600.0
	baseXP := float64(pokePoke.BaseExperience)
	chance := 1.0 - ((baseXP - minXP) / (maxXP - minXP))
	if chance > 0.95 {
		chance = 0.95
	}
	if chance < 0.05 {
		chance = 0.05
	}

	catch := float64(rand.Intn(100)) / 100

	if catch <= chance {
		cfg.Pokedex.Add(pokePoke)
		fmt.Printf("%s was caught!\n", pokePoke.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokePoke.Name)
	}
	return nil
}

func commandInspect(cfg *Config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("error invalid number of arguments")
	}

	pokemonName := args[0]
	if err := cfg.Pokedex.Inspect(pokemonName); err != nil {
		return err
	}

	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	cfg.Pokedex.Print()
	return nil
}
