package hapi

import (
	"github.com/gin-gonic/gin"
	"github.com/utm/station-mgmt/internal/config"
	"github.com/utm/station-mgmt/internal/service"
)

func NewRouter(cfg *config.Config, overviewSvc *service.OverviewService) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "station-mgmt"})
	})

	h := NewOverviewHandler(overviewSvc)

	v1 := r.Group("/api/v1")
	v1.GET("/overview", h.GetOverview)

	return r
}
