package counter

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

func NewService(db dataStore, counter cacheStore, config *Config) *service {
	return &service{dataStore: db, cacheStore: counter, config: config}
}

func (s *service) Reset() int {
	s.dataStore.SetCount(0)
	s.cacheStore.SetCount(0)
	return 0
}

func (s *service) Info() (int, error) {
	return recover(s.cacheStore, s.dataStore)
}

func (s *service) Increase() (int, error) {
	count, err := recover(s.cacheStore, s.dataStore)
	count, err = s.cacheStore.Incr()
	if count > s.config.Limit {
		return 0, errors.New("over limit")
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func recover(cache cacheStore, db dataStore) (int, error) {
	count, err := cache.GetCount()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			count = db.GetCount()
			cache.SetCount(count)
		} else {
			return 0, err
		}
	}
	return count, nil
}
