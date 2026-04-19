package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type TelemetryHandler struct {
	svc *service.TelemetryService
}

func NewTelemetryHandler(svc *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{svc: svc}
}

func (h *TelemetryHandler) RegisterRoutes(root *gin.Engine, protected *gin.RouterGroup) {
	root.GET("/ws/monitor", h.Stream)

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

func (h *TelemetryHandler) Stream(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("websocket upgrade failed")
		return
	}
	defer conn.Close()
	h.svc.StreamToClient(c.Request.Context(), conn)
}
