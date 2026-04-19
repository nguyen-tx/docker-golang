package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/backend/internal/service"
	"github.com/utm/backend/internal/types"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	auth := public.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
	protected.GET("/auth/me", h.Me)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"token_type":   "Bearer",
		"user":         result.User,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authSvc.Register(c.Request.Context(), service.RegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Organization: req.Organization,
		LicenseNo:    req.LicenseNo,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.authSvc.GetProfile(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
