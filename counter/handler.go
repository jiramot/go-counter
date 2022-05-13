package counter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHandler(svc *service) *handler {
	return &handler{svc: svc}
}

func (h *handler) Reset(c *gin.Context) {
	count := h.svc.Reset()
	c.JSON(http.StatusOK, countResponse{Count: count})
}

func (h *handler) Info(c *gin.Context) {
	count, err := h.svc.Info()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, countResponse{Count: count})
}

func (h *handler) Increase(c *gin.Context) {
	count, err := h.svc.Increase()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, countResponse{Count: count})

}

type countResponse struct {
	Count int `json:"count"`
}
