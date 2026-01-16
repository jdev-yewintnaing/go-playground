package generic

import "sync"

type Cache[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		items: make(map[string]T),
	}
}

func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value

}

func (c *Cache[T]) Get(key string) (value T, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok = c.items[key]
	return value, ok
}

func Filter[T any](ts []T, f func(T) bool) []T {
	result := make([]T, 0, len(ts))
	for _, t := range ts {
		if f(t) {
			result = append(result, t)
		}
	}

	return result
}
