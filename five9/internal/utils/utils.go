package utils

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type CacheResponse[Key comparable, T any] struct {
	LastUpdated *time.Time `json:"lastUpdated"`
	Items       map[Key]T  `json:"items"`
}

type MemoryCacheInstance[Key comparable, T any] struct {
	mutex       *sync.Mutex
	lastUpdated *time.Time
	maxAge      *time.Duration
	items       map[Key]T
}

func NewMemoryCacheInstance[Key comparable, T any](maxAllowedAge *time.Duration) *MemoryCacheInstance[Key, T] {
	return &MemoryCacheInstance[Key, T]{
		mutex:       &sync.Mutex{},
		items:       map[Key]T{},
		maxAge:      maxAllowedAge,
		lastUpdated: nil,
	}
}

func (cache *MemoryCacheInstance[Key, T]) Replace(freshData map[Key]T) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	replaceTime := time.Now()

	cache.items = freshData
	cache.lastUpdated = &replaceTime
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

func (cache *MemoryCacheInstance[Key, T]) Delete(key Key) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.items, key)
}

func (cache *MemoryCacheInstance[Key, T]) Get(key Key) (T, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	target := *new(T)

	{ // Check cache age
		if cache.lastUpdated == nil {
			return *new(T), false
		}

		if cache.maxAge != nil {
			if time.Since(*cache.lastUpdated) > *cache.maxAge {
				return *new(T), false
			}
		}
	}

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

func (cache *MemoryCacheInstance[Key, T]) GetAll() (CacheResponse[Key, T], error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	{ // Check cache age
		if cache.lastUpdated == nil {
			return CacheResponse[Key, T]{}, ErrWebSocketCacheNotReady
		}

		if cache.maxAge != nil {
			if time.Since(*cache.lastUpdated) > *cache.maxAge {
				return CacheResponse[Key, T]{}, ErrWebSocketCacheStale
			}
		}
	}

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
	}, nil
}

func (cache *MemoryCacheInstance[Key, T]) GetLastUpdated() *time.Time {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	return cache.lastUpdated
}

func (cache *MemoryCacheInstance[Key, T]) GetCacheAge() *time.Duration {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.lastUpdated != nil {
		age := time.Since(*cache.lastUpdated)

		return &age
	}

	return nil
}

var (
	ErrWebSocketCacheNotReady error = errors.New("webSocket cache is not ready")
	ErrWebSocketCacheStale    error = errors.New("webSocket cache is stale")
)
