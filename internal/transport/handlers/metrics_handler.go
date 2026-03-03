package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gostream/internal/domain"
)

type MetricsHandler struct {
	store domain.MetricsStore
}

func NewMetricsHandler(store domain.MetricsStore) *MetricsHandler {
	return &MetricsHandler{store: store}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	snapshot := h.store.GetSnapshot()
	c.JSON(http.StatusOK, snapshot)
}