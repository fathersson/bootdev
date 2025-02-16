package main

import (
	"bootdev/internal/commands"
	"bootdev/internal/pokecache"
	"bootdev/internal/types"
	"time"
)

func main() {
	c := &types.Config{
		Next:     "",
		Previous: "",
	}

	interval := 10 * time.Minute
	cache := pokecache.NewCache(interval) // Создаем новый кэш

	commands.ScannerMain(c, cache)
}
