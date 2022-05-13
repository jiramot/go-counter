package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jiramot/go-counter/counter"
	"strconv"
)

type cache struct {
	rdb    *redis.Client
	config *counter.Config
}

func NewCacheStore(rdb *redis.Client, config *counter.Config) *cache {
	return &cache{rdb: rdb, config: config}
}

func (c *cache) Incr() (int, error) {
	ctx := context.Background()
	val, err := c.rdb.Incr(ctx, c.config.Key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (c *cache) GetCount() (int, error) {
	ctx := context.Background()
	val, err := c.rdb.Get(ctx, c.config.Key).Result()
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *cache) SetCount(count int) {
	ctx := context.Background()
	_, _ = c.rdb.Set(ctx, c.config.Key, count, c.config.Ttl).Result()
}

//func (c *cache) IsOverLimit() bool {
//	ctx := context.Background()
//	val, err := c.rdb.Get(ctx, c.config.Key).Result()
//	if err != nil {
//		return true
//	}
//	count, err := strconv.Atoi(val)
//	if err != nil {
//		return true
//	}
//	if count >= c.config.Limit {
//		return true
//	}
//	return false
//}
