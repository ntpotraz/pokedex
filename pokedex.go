package main

import (
	"fmt"
	"sync"
)

type Pokedex struct {
	data map[string]Pokemon
	mu   sync.RWMutex
}

type Stat struct {
	name string
	val  int
}

type Pokemon struct {
	name   string
	height int
	weight int
	stats  []Stat
	types  []string
}

func (p *Pokedex) Add(data PokemonJSON) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	types := []string{}
	for _, pType := range data.Types {
		types = append(types, pType.Type.Name)
	}

	stats := []Stat{}
	for _, stat := range data.Stats {
		stats = append(stats, Stat{name: stat.Stat.Name, val: stat.BaseStat})
	}

	pokemon := Pokemon{
		name:   data.Name,
		height: data.Height,
		weight: data.Weight,
		stats:  stats,
		types:  types,
	}

	p.data[data.Name] = pokemon

	return nil
}

func (p *Pokedex) Get(name string) (Pokemon, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pokemon, ok := p.data[name]
	if !ok {
		return Pokemon{}, fmt.Errorf("error %s not in pokedex", name)
	}

	return pokemon, nil
}

func (p *Pokedex) Check(name string) bool {
	_, err := p.Get(name)
	if err != nil {
		return false
	}
	return true
}

func (p *Pokedex) Inspect(name string) error {
	pokemon, err := p.Get(name)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", pokemon.name)
	fmt.Printf("Height: %d\n", pokemon.height)
	fmt.Printf("Weight: %d\n", pokemon.weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.stats {
		fmt.Printf("  -%s: %d\n", stat.name, stat.val)
	}

	fmt.Println("Types:")
	for _, pType := range pokemon.types {
		fmt.Printf("  - %s\n", pType)
	}

	return nil
}

func (p *Pokedex) Print() {
	fmt.Println("Your Pokedex:")
	if len(p.data) == 0 {
		fmt.Println("~empty~")
		return
	}

	for _, item := range p.data {
		fmt.Println(" ", "-", item.name)
	}
}

func NewPokedex() *Pokedex {
	pokedex := Pokedex{
		data: map[string]Pokemon{},
		mu:   sync.RWMutex{},
	}

	return &pokedex
}
