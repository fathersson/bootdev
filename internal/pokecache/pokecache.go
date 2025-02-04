package pokecache

import (
	"sync"
	"time"
)

// Структура для хранения данных кэша
type cacheEntry struct {
	createdAt time.Time // Время создания записи
	val       []byte    // Данные, которые мы кэшируем
}

// Структура кэша для хранения записей
type Cache struct {
	data     map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration // Время жизни кэшированных данных
}

func NewCache() *Cache {
	c := &Cache{
		data:     make(map[string]cacheEntry),
		interval: time.Duration(10),
	}
	c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.data[key] = cacheEntry{time.Now(), val}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	if entry, ok := c.data[key]; ok {
		c.mu.Unlock()
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	go func() {
		for range time.Tick(c.interval) {
			c.mu.Lock()
			for key, entry := range c.data {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
