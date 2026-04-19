package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/middleware"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type AircraftHandler struct {
	svc *service.AircraftService
}

func NewAircraftHandler(svc *service.AircraftService) *AircraftHandler {
	return &AircraftHandler{svc: svc}
}

func (h *AircraftHandler) RegisterRoutes(protected *gin.RouterGroup) {
	aircraft := protected.Group("/aircraft")
	{
		aircraft.GET("", h.List)
		aircraft.GET("/:id", h.Get)
		aircraft.POST("", h.Create)
		aircraft.PUT("/:id", h.Update)
		aircraft.DELETE("/:id", middleware.RequireRole(model.RoleAdmin), h.Delete)
	}
}

func (h *AircraftHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.MustGet("user_role").(model.UserRole)

	ownerID := ""
	if role == model.RolePilot || role == model.RoleOperator {
		ownerID = userID
	}

	list, err := h.svc.List(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *AircraftHandler) Get(c *gin.Context) {
	aircraft, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "aircraft not found"})
		return
	}
	c.JSON(http.StatusOK, aircraft)
}

func (h *AircraftHandler) Create(c *gin.Context) {
	var input model.Aircraft
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	aircraft, err := h.svc.Create(c.Request.Context(), c.GetString("user_id"), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, aircraft)
}

func (h *AircraftHandler) Update(c *gin.Context) {
	var input model.Aircraft
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	aircraft, err := h.svc.Update(c.Request.Context(), c.Param("id"), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, aircraft)
}

func (h *AircraftHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
