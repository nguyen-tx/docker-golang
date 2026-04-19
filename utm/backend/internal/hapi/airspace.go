package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/middleware"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type AirspaceHandler struct {
	svc *service.AirspaceService
}

func NewAirspaceHandler(svc *service.AirspaceService) *AirspaceHandler {
	return &AirspaceHandler{svc: svc}
}

func (h *AirspaceHandler) RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.GET("/airspace", h.List)
	public.GET("/airspace/:id", h.Get)

	airspace := protected.Group("/airspace")
	airspace.Use(middleware.RequireRole(model.RoleAdmin, model.RoleATC))
	{
		airspace.POST("", h.Create)
		airspace.PUT("/:id", h.Update)
		airspace.DELETE("/:id", middleware.RequireRole(model.RoleAdmin), h.Delete)
	}
}

func (h *AirspaceHandler) List(c *gin.Context) {
	onlyActive := c.Query("active") == "true"
	zones, err := h.svc.List(c.Request.Context(), onlyActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, zones)
}

func (h *AirspaceHandler) Get(c *gin.Context) {
	zone, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "airspace zone not found"})
		return
	}
	c.JSON(http.StatusOK, zone)
}

func (h *AirspaceHandler) Create(c *gin.Context) {
	var input model.AirspaceZone
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	zone, err := h.svc.Create(c.Request.Context(), c.GetString("user_id"), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, zone)
}

func (h *AirspaceHandler) Update(c *gin.Context) {
	var input model.AirspaceZone
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	zone, err := h.svc.Update(c.Request.Context(), c.Param("id"), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, zone)
}

func (h *AirspaceHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
