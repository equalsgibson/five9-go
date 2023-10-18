package utils

import (
	"encoding/json"
	"sync"
	"time"
)

type CacheResponse[Key comparable, T any] struct {
	LastUpdated *time.Time `json:"lastUpdated"`
	Items       map[Key]T  `json:"items"`
}

func NewMemoryCacheInstance[Key comparable, T any]() *MemoryCacheInstance[Key, T] {
	return &MemoryCacheInstance[Key, T]{
		mutex:       &sync.Mutex{},
		items:       map[Key]T{},
		lastUpdated: nil,
	}
}

type MemoryCacheInstance[Key comparable, T any] struct {
	mutex       *sync.Mutex
	lastUpdated *time.Time
	items       map[Key]T
}

func (cache *MemoryCacheInstance[Key, T]) Reset() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.lastUpdated = nil

	cache.items = map[Key]T{}
}

func (cache *MemoryCacheInstance[Key, T]) Update(key Key, item T) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	cache.lastUpdated = &now

	cache.items[key] = item
}

func (cache *MemoryCacheInstance[Key, T]) Get(key Key) (T, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	target := *new(T)

	item, ok := cache.items[key]
	if !ok {
		return *new(T), false
	}

	itemBytes, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(itemBytes, &target); err != nil {
		panic(err)
	}

	return item, true
}

func (cache *MemoryCacheInstance[Key, T]) GetAll() CacheResponse[Key, T] {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	target := map[Key]T{}

	itemBytes, err := json.Marshal(cache.items)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(itemBytes, &target); err != nil {
		panic(err)
	}

	return CacheResponse[Key, T]{
		LastUpdated: cache.lastUpdated,
		Items:       target,
	}
}
