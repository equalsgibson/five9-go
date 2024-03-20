package five9

import (
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type Cache[Key comparable, T any] interface {
	Replace(freshData map[Key]T) error
	Reset() error
	Update(key Key, item T) error
	Delete(key Key) error
	Get(key Key) (T, bool)
	GetAll() (five9types.CacheResponse[Key, T], error)
	GetLastUpdated() (*time.Time, error)
	GetCacheAge() (*time.Duration, error)
}
