package utils

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

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

func (cache *MemoryCacheInstance[Key, T]) Replace(freshData map[Key]T) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	replaceTime := time.Now()

	cache.items = freshData
	cache.lastUpdated = &replaceTime

	return nil
}

func (cache *MemoryCacheInstance[Key, T]) Reset() error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.lastUpdated = nil

	cache.items = map[Key]T{}

	return nil
}

func (cache *MemoryCacheInstance[Key, T]) Update(key Key, item T) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	cache.lastUpdated = &now

	cache.items[key] = item

	return nil
}

func (cache *MemoryCacheInstance[Key, T]) Delete(key Key) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.items, key)

	return nil
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

func (cache *MemoryCacheInstance[Key, T]) GetAll() (five9types.CacheResponse[Key, T], error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	{ // Check cache age
		if cache.lastUpdated == nil {
			return five9types.CacheResponse[Key, T]{}, ErrWebSocketCacheNotReady
		}

		if cache.maxAge != nil {
			if time.Since(*cache.lastUpdated) > *cache.maxAge {
				return five9types.CacheResponse[Key, T]{}, ErrWebSocketCacheStale
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

	return five9types.CacheResponse[Key, T]{
		LastUpdated: cache.lastUpdated,
		Items:       target,
	}, nil
}

func (cache *MemoryCacheInstance[Key, T]) GetLastUpdated() (*time.Time, error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	return cache.lastUpdated, nil
}

func (cache *MemoryCacheInstance[Key, T]) GetCacheAge() (*time.Duration, error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.lastUpdated != nil {
		age := time.Since(*cache.lastUpdated)

		return &age, nil
	}

	return nil, ErrWebSocketCacheNotReady
}

var (
	ErrWebSocketCacheNotReady error = errors.New("webSocket cache is not ready")
	ErrWebSocketCacheStale    error = errors.New("webSocket cache is stale")
)
