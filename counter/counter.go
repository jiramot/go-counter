package counter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"time"
)

type cacheStore interface {
	Incr() (int, error)
	SetCount(int)
	GetCount() (int, error)
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
	db     dataStore
	cache  cacheStore
	config *Config
}

func NewHandle(db dataStore, counter cacheStore, config *Config) *handler {
	return &handler{db: db, cache: counter, config: config}
}

func (h *handler) Setup(c *gin.Context) {
	query := c.Param("count")
	count, _ := strconv.Atoi(query)
	h.db.SetCount(count)
	h.cache.SetCount(count)
	c.JSON(http.StatusOK, countResponse{Count: count})
}

func (h *handler) Info(c *gin.Context) {
	count, err := recover(h.cache, h.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, countResponse{Count: count})
}

func (h *handler) Increase(c *gin.Context) {
	count, _ := h.cache.GetCount()
	if count < h.config.Limit {
		count, err := recover(h.cache, h.db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		count, err = h.cache.Incr()
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
