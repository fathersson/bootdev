package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"bootdev/internal/pokecache"

	"math/rand"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, ...string) error
}

type pokemon struct {
	Val []byte
}

type config struct {
	Next     string
	Previous string
}

func main() {
	c := &config{
		Next:     "",
		Previous: "",
	}

	// Создаем новый кэш
	interval := 10 * time.Minute
	cache := pokecache.NewCache(interval)

	scanner := bufio.NewScanner(os.Stdin)
	var str string
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for i, word := range words {
			for _, s := range word {
				if strings.ToUpper(string(s)) == string(s) {
					symbol := strings.ToLower(string(s))
					str = str + symbol
					continue
				}
				str = str + string(s)
			}
			if i < len(words)-1 {
				str = str + " "
			}
		}
		params := words[1:]
		commands := commands()
		commandsName := words[0]
		if command, ok := commands[commandsName]; ok {
			command.callback(c, cache, params...)
			str = ""
			continue
		} else {
			fmt.Println("Unknown command")
			str = ""
		}
	}
}

func commands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "map now page",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "map back page",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore now page",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "pokedex all pokemons",
			callback:    commandPokedex,
		},
	}
	return commands
}

func commandPokedex(c *config, cache *pokecache.Cache, params ...string) error {
	pokedex := cache.GetAll()
	substring := "pokedex"
	for key, value := range pokedex {
		if strings.Contains(key, substring) {
			fmt.Println("    -" + string(value.Val))
		}
	}
	return errors.New("pokedex")
}

func commandInspect(c *config, cache *pokecache.Cache, params ...string) error {
	var inspect struct {
		Name   string `json:"name"`
		Height int    `json:"height"`
		Weight int    `json:"weight"`
		Stats  []struct {
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
	}
	body1, ok := cache.Get(params[0])
	if !ok {
		fmt.Println("в кэше нет такого ключа")
		return errors.New("qwe")
	}

	err := json.Unmarshal(body1, &inspect)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name: %s\n", inspect.Name)
	fmt.Printf("Weight: %d\n", inspect.Weight)
	fmt.Printf("Height: %d\n", inspect.Height)
	fmt.Println("Stats:")
	for _, field := range inspect.Stats {
		fmt.Printf("	-%s: %d\n", field.Stat.Name, field.BaseStat)
	}
	fmt.Println("Types:")
	for _, field := range inspect.Types {
		fmt.Printf("	-%s\n", field.Type.Name)
	}

	return errors.New("inspect")
}

func commandCatch(c *config, cache *pokecache.Cache, params ...string) error {

	fmt.Printf("Throwing a Pokeball at %s...\n", params[0])

	urlPokemon := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", params[0])

	res, err := http.Get(urlPokemon)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var pokemonExp struct {
		Base_experience int `json:"base_experience"`
	}

	err = json.Unmarshal(body, &pokemonExp)
	if err != nil {
		log.Printf("Ошибка при разборе JSON commandCatch: %v", err)
		return err
	}

	rand := rand.Intn(11)
	catchPokemon := make(map[string]pokemon)

	if rand != 0 {
		resultCatch := pokemonExp.Base_experience / rand
		if resultCatch <= 20 {
			fmt.Printf("%s was caught!\n", params[0])
			catchPokemon[params[0]] = pokemon{
				Val: []byte(params[0]),
			}
			cache.Add(params[0]+"pokedex", catchPokemon[params[0]].Val)
			cache.Add(params[0], body)
		} else {
			fmt.Printf("%s escaped!\n", params[0])
		}
	} else {
		fmt.Printf("%s escaped!\n", params[0])
	}

	return errors.New("catch")
}

func commandExit(c *config, cache *pokecache.Cache, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("exit")
}

func commandHelp(c *config, cache *pokecache.Cache, params ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	commandsHelp := commands()
	for _, command := range commandsHelp {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return errors.New("help")
}

func commandMap(c *config, cache *pokecache.Cache, params ...string) error {
	if c.Next == "" {
		c.Next = "https://pokeapi.co/api/v2/location-area/"
	}
	var dat []byte

	// Попытка получить данные из кэша
	body2, ok := cache.Get(c.Next)

	dat = []byte(body2)
	if !ok {
		res, err := http.Get(c.Next)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}

		dat = []byte(body)
		cache.Add(c.Next, dat)
	}

	var area struct {
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Results  []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		}
	}

	err := json.Unmarshal(dat, &area)
	if err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		return err
	}

	for _, result := range area.Results {
		fmt.Printf("%s\n", result.Name)
	}
	c.Next = area.Next
	c.Previous = area.Previous
	return errors.New("map")
}

func commandMapb(c *config, cache *pokecache.Cache, params ...string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		c.Next = c.Previous
		commandMap(c, cache)
	}
	return errors.New("mapb")
}

func commandExplore(c *config, cache *pokecache.Cache, params ...string) error {

	areaParams := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", params[0])

	res, err := http.Get(areaParams)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var pokeName struct {
		Pokemon_encounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			}
		}
	}

	err = json.Unmarshal(body, &pokeName)
	if err != nil {
		log.Printf("Ошибка при разборе JSON commandExplore: %v", err)
		return err
	}

	for _, result := range pokeName.Pokemon_encounters {
		fmt.Printf("%s\n", result.Pokemon.Name)
	}

	cache.Add(areaParams, body)

	return errors.New("explore")
}

func cleanInput(text string) []string {
	var c int
	var result []string
	for i := 0; i < len(text); i++ {
		if text[i] != ' ' {
			var str string
			for _, s := range text[i:] {
				if s == ' ' {
					break
				}
				if strings.ToUpper(string(s)) == string(s) {
					symbol := strings.ToLower(string(s))
					str = str + symbol
					c++
					continue
				}
				c++
				str = str + string(s)
			}
			result = append(result, str)
			str = ""
			i = i + c
			c = 0

		}

	}
	return result
}
