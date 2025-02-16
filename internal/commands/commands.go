package commands

import (
	"bootdev/internal/myhttp"
	"bootdev/internal/myjson"
	"bootdev/internal/pokecache"
	"bootdev/internal/types"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func ScannerMain(c *types.Config, cache *pokecache.Cache) {
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
		commands := Commands()
		commandsName := words[0]
		if command, ok := commands[commandsName]; ok {
			command.Callback(c, cache, params...)
			str = ""
			continue
		}
		fmt.Println("Unknown command")
		str = ""
	}
}

func Commands() map[string]types.CliCommand {
	commands := map[string]types.CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "map now page",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "map back page",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "explore now page",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "catch pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "inspect pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "pokedex all pokemons",
			Callback:    commandPokedex,
		},
	}
	return commands
}

func commandPokedex(c *types.Config, cache *pokecache.Cache, params ...string) error {
	pokedex := cache.GetAll()
	substring := "pokedex"
	for key, value := range pokedex {
		if strings.Contains(key, substring) {
			fmt.Println("    -" + string(value.Val))
		}
	}
	return nil
}

func commandInspect(c *types.Config, cache *pokecache.Cache, params ...string) error {
	cacheBody := myhttp.CacheGet(cache, params)
	var inspect types.Inspect
	myjson.Unmarshal(cacheBody, &inspect)

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
	return nil
}

func commandCatch(c *types.Config, cache *pokecache.Cache, params ...string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", params[0])

	urlPokemon := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", params[0])
	body := myhttp.Get(urlPokemon)
	var Exp types.PokemonExp
	myjson.Unmarshal(body, &Exp)

	rand := rand.Intn(11)
	catchPokemon := make(map[string]types.Pokemon)

	if rand == 0 {
		fmt.Printf("%s escaped!\n", params[0])
		return nil
	}
	resultCatch := Exp.BaseExperience / rand
	fmt.Println(resultCatch)
	if resultCatch <= 20 {
		fmt.Printf("%s was caught!\n", params[0])
		catchPokemon[params[0]] = types.Pokemon{
			Val: []byte(params[0]),
		}
		cache.Add(params[0]+"pokedex", catchPokemon[params[0]].Val)
		cache.Add(params[0], body)
	} else {
		fmt.Printf("%s escaped!\n", params[0])
	}
	return nil
}

func commandExit(c *types.Config, cache *pokecache.Cache, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *types.Config, cache *pokecache.Cache, params ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	commandsHelp := Commands()
	for _, command := range commandsHelp {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMap(c *types.Config, cache *pokecache.Cache, params ...string) error {
	if c.Next == "" {
		c.Next = "https://pokeapi.co/api/v2/location-area/"
	}

	dat := myhttp.CacheGetMap(c, cache)
	var area types.Area
	myjson.Unmarshal(dat, &area)

	for _, result := range area.Results {
		fmt.Printf("%s\n", result.Name)
	}

	c.Next = area.Next
	c.Previous = area.Previous
	return nil
}

func commandMapb(c *types.Config, cache *pokecache.Cache, params ...string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		c.Next = c.Previous
		commandMap(c, cache)
	}
	return nil
}

func commandExplore(c *types.Config, cache *pokecache.Cache, params ...string) error {
	areaParams := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", params[0])
	body := myhttp.Get(areaParams)

	var name types.PokeName
	myjson.Unmarshal(body, &name)
	for _, result := range name.Pokemon_encounters {
		fmt.Printf("%s\n", result.Pokemon.Name)
	}

	cache.Add(areaParams, body)
	return nil
}

func CleanInput(text string) []string {
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
