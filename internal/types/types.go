package types

import (
	"bootdev/internal/pokecache"
)

type Pokemon struct {
	Val []byte
}

type PokeName struct {
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		}
	}
}

type PokemonExp struct {
	BaseExperience int `json:"base_experience"`
}

type Area struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}

type Inspect struct {
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

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, *pokecache.Cache, ...string) error
}

type Config struct {
	Next     string
	Previous string
}
