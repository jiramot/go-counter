package counter

import (
	"time"
)

type cacheStore interface {
	Incr() (int, error)
	SetCount(int)
	GetCount() (int, error)
}

type dataStore interface {
	GetCount() int
	SetCount(int) bool
}

type Config struct {
	Key   string
	Ttl   time.Duration
	Limit int
}

type handler struct {
	dataStore  dataStore
	cacheStore cacheStore
	config     *Config
}
