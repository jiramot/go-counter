package counter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
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

func NewHandle(db dataStore, counter cacheStore, config *Config) *handler {
	return &handler{dataStore: db, cacheStore: counter, config: config}
}

func (h *handler) Reset(c *gin.Context) {
	h.dataStore.SetCount(0)
	h.cacheStore.SetCount(0)
	c.JSON(http.StatusOK, countResponse{Count: 0})
}

func (h *handler) Info(c *gin.Context) {
	count, err := recover(h.cacheStore, h.dataStore)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, countResponse{Count: count})
}

func (h *handler) Increase(c *gin.Context) {
	count, _ := h.cacheStore.GetCount()
	if count < h.config.Limit {
		count, err := recover(h.cacheStore, h.dataStore)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		count, err = h.cacheStore.Incr()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(200, countResponse{Count: count})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "over limit",
		})
	}
}

type countResponse struct {
	Count int `json:"count"`
}

func recover(cache cacheStore, db dataStore) (int, error) {
	count, err := cache.GetCount()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			count = db.GetCount()
		} else {
			return 0, err
		}
	}
	return count, nil
}
