package myhttp

import (
	"bootdev/internal/pokecache"
	"bootdev/internal/types"
	"io"
	"log"
	"net/http"
)

func Get(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при выполнении GET-запроса: %v", err)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Ошибка при чтении тела ответа: %v", err)
		return nil
	}
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	defer res.Body.Close()
	return body
}

func CacheGet(cache *pokecache.Cache, params []string) []byte {
	body, ok := cache.Get(params[0])
	if !ok {
		log.Println("в кэше нет такого ключа")
	}
	return body
}

func CacheGetMap(c *types.Config, cache *pokecache.Cache) []byte {
	var dat []byte
	body2, ok := cache.Get(c.Next) // Попытка получить данные из кэша
	dat = []byte(body2)
	if !ok {
		body := Get(c.Next)
		dat = []byte(body)
		cache.Add(c.Next, dat)
	}
	return dat
}
