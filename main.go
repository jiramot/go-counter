package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jiramot/go-counter/counter"
	"github.com/jiramot/go-counter/database"
	"time"
)

func main() {
	config := &counter.Config{
		Limit: 10,
		Key:   "counter",
		Ttl:   0 * time.Second,
	}
	db := database.NewInMemory()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := gin.Default()
	cache := database.NewCacheStore(rdb, config)
	handler := counter.NewHandle(db, cache, config)
	r.GET("/info", handler.Info)
	r.GET("/increase", handler.Increase)
	r.GET("/setup/:count", handler.Setup)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
