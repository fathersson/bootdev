package main

import (
	"bootdev/internal/pokecache"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
			description: "map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "mapb",
			callback:    commandMapb,
		},
	}

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
		if str == commands["map"].name {
			commandMap(c)
			str = ""
			continue
		}
		if str == commands["mapb"].name {
			commandMapb(c)
			str = ""
			continue
		}
		if str == commands["help"].name {
			commandHelp(c)
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			str = ""
			continue
		}
		if str == commands["exit"].name {
			commandExit(c)
			os.Exit(0)
		} else {
			fmt.Println("Unknown command")
			str = ""
		}
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return errors.New("exit")
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	return errors.New("help")
}

func commandMap(c *config) error {
	if c.Next == "" {
		c.Next = "https://pokeapi.co/api/v2/location-area/"
	}

	var res *http.Response
	var err error
	if data, ok := pokecache.Cache.Get(c.Next); ok {
		res = data
	} else {
		res, err = http.Get(c.Next)
		if err != nil {
			log.Fatal(err)
		}
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	dat := []byte(body)

	var area struct {
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Results  []struct {
			Name string `json:"name"`
		}
	}

	err = json.Unmarshal(dat, &area)
	if err != nil {

		fmt.Println(err)

	}

	for _, result := range area.Results {
		fmt.Printf("%s\n", result.Name)
	}
	c.Next = area.Next
	c.Previous = area.Previous
	return errors.New("map")
}

func commandMapb(c *config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		c.Next = c.Previous
		commandMap(c)
	}
	return errors.New("mapb")
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
	fmt.Println(result)
	return result
}
