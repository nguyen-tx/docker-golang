package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/middleware"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type AuthorizationHandler struct {
	svc *service.AuthorizationService
}

func NewAuthorizationHandler(svc *service.AuthorizationService) *AuthorizationHandler {
	return &AuthorizationHandler{svc: svc}
}

func (h *AuthorizationHandler) RegisterRoutes(protected *gin.RouterGroup) {
	authz := protected.Group("/authorizations")
	{
		authz.GET("", h.List)
		authz.GET("/:id", h.Get)
		authz.POST("", middleware.RequireRole(model.RolePilot, model.RoleOperator), h.Create)
		authz.POST("/:id/review", middleware.RequireRole(model.RoleATC, model.RoleAdmin), h.Review)
		authz.POST("/:id/cancel", h.Cancel)
	}
}

func (h *AuthorizationHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.MustGet("user_role").(model.UserRole)

	filter := service.AuthFilter{Status: c.Query("status")}
	if role == model.RolePilot || role == model.RoleOperator {
		filter.PilotID = userID
	}

	list, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *AuthorizationHandler) Get(c *gin.Context) {
	req, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "authorization request not found"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *AuthorizationHandler) Create(c *gin.Context) {
	var req types.CreateAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.svc.Create(c.Request.Context(), c.GetString("user_id"), service.CreateAuthInput{
		AircraftID:     req.AircraftID,
		ZoneID:         req.ZoneID,
		Purpose:        req.Purpose,
		FlightAreaJSON: req.FlightAreaJSON,
		PlannedAltM:    req.PlannedAltM,
		PlannedStartAt: req.PlannedStartAt,
		PlannedEndAt:   req.PlannedEndAt,
		PilotNotes:     req.PilotNotes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *AuthorizationHandler) Review(c *gin.Context) {
	var req types.ReviewAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.svc.Review(c.Request.Context(), c.Param("id"), c.GetString("user_id"), req.Action, req.ReviewNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AuthorizationHandler) Cancel(c *gin.Context) {
	if err := h.svc.Cancel(c.Request.Context(), c.Param("id"), c.GetString("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.MessageResponse{Message: "authorization request cancelled"})
}
