package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type TelemetryHandler struct {
	svc *service.TelemetryService
}

func NewTelemetryHandler(svc *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{svc: svc}
}

// RegisterRoutes chỉ còn REST — WebSocket đã chuyển sang utm-alert.
func (h *TelemetryHandler) RegisterRoutes(protected *gin.RouterGroup) {
	telemetry := protected.Group("/telemetry")
	{
		telemetry.POST("/push", h.Push)
		telemetry.GET("/sessions", h.LiveSessions)
		telemetry.GET("/sessions/:session_id/history", h.History)
	}
}

func (h *TelemetryHandler) Push(c *gin.Context) {
	var data model.Telemetry
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.svc.Push(c.Request.Context(), &data); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.MessageResponse{Message: "telemetry received"})
}

func (h *TelemetryHandler) History(c *gin.Context) {
	data, err := h.svc.GetSessionHistory(c.Request.Context(), c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TelemetryHandler) LiveSessions(c *gin.Context) {
	sessions, err := h.svc.GetActiveSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}
