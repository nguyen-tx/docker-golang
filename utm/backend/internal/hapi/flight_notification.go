package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type FlightNotificationHandler struct {
	svc *service.FlightNotificationService
}

func NewFlightNotificationHandler(svc *service.FlightNotificationService) *FlightNotificationHandler {
	return &FlightNotificationHandler{svc: svc}
}

func (h *FlightNotificationHandler) RegisterRoutes(protected *gin.RouterGroup) {
	g := protected.Group("/flight-notifications")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Submit)
}

func (h *FlightNotificationHandler) List(c *gin.Context) {
	operatorID := c.GetString("user_id")
	list, err := h.svc.List(c.Request.Context(), operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *FlightNotificationHandler) Get(c *gin.Context) {
	fn, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "flight notification not found"})
		return
	}
	c.JSON(http.StatusOK, fn)
}

// Submit nhận thông báo bay, gọi SCRP đồng bộ và trả về kết quả ngay
func (h *FlightNotificationHandler) Submit(c *gin.Context) {
	var input model.FlightNotification
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	fn, approvedPlan, rejected, err := h.svc.Submit(c.Request.Context(), c.GetString("user_id"), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if rejected != nil {
		c.JSON(http.StatusOK, gin.H{
			"flight_notification": fn,
			"status":              "utm_rejected",
			"reason":              rejected.Reason,
			"detail":              rejected.Detail,
			"earliest_possible":   rejected.EarliestPossible,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"flight_notification": fn,
		"status":              "soft_reserved",
		"approved_plan":       approvedPlan,
	})
}
