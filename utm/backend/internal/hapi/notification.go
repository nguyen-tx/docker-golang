package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) RegisterRoutes(protected *gin.RouterGroup) {
	notif := protected.Group("/notifications")
	{
		notif.GET("", h.List)
		notif.PATCH("/:id/read", h.MarkRead)
	}
}

func (h *NotificationHandler) List(c *gin.Context) {
	notifs, err := h.svc.ListForUser(c.Request.Context(), c.GetString("user_id"), c.Query("unread") == "true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	if err := h.svc.MarkRead(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.MessageResponse{Message: "marked as read"})
}
