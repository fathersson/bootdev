package pokecache

import (
	"sync"
	"time"
)

// Структура для хранения данных кэша
type cacheEntry struct {
	CreatedAt time.Time // Время создания записи
	Val       []byte    // Данные, которые мы кэшируем
}

// Структура кэша для хранения записей
type Cache struct {
	data     map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration // Время жизни кэшированных данных
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
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
		return entry.Val, true
	} else {
		c.mu.Unlock()
		return nil, false
	}

}

func (c *Cache) GetAll() map[string]cacheEntry {
	c.mu.Lock()
	result := c.data
	c.mu.Unlock()
	return result
}

func (c *Cache) reapLoop() {
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			c.mu.Lock()
			for key, entry := range c.data {
				if time.Since(entry.CreatedAt) > c.interval {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
