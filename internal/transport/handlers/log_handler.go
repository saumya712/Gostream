package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gostream/internal/application"
	"gostream/internal/domain"
)

type LogHandler struct {
	service *application.LogService
}

func NewLogHandler(service *application.LogService) *LogHandler {
	return &LogHandler{service: service}
}

type logRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Level       string `json:"level" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

func (h *LogHandler) Ingest(c *gin.Context) {
	var req logRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log := &domain.Log{
		ID:          "",
		SERVICENAME: req.ServiceName,
		LEVEL	:       domain.LogLevel(req.Level),
		MESSAGE:     req.Message,
		TIMESTAMP:   time.Now(),
	}

	if err := h.service.Ingest(c.Request.Context(), log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "log accepted"})
}