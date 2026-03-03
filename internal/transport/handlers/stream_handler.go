package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gostream/internal/domain"

	"github.com/gin-gonic/gin"
)

type StreamHandler struct {
	store domain.MetricsStore
}

func NewStreamHandler(store domain.MetricsStore) *StreamHandler {
	return &StreamHandler{store: store}
}

func (h *StreamHandler) Stream(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Writer.WriteHeader(http.StatusOK)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return

		case <-ticker.C:
			snapshot := h.store.GetSnapshot()

			data, _ := json.Marshal(snapshot)

			c.Writer.Write([]byte("data: "))
			c.Writer.Write(data)
			c.Writer.Write([]byte("\n\n"))

			c.Writer.Flush()
		}
	}
}
