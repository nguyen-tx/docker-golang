package hapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utm/station-mgmt/internal/service"
)

type OverviewHandler struct {
	overviewSvc *service.OverviewService
}

func NewOverviewHandler(overviewSvc *service.OverviewService) *OverviewHandler {
	return &OverviewHandler{overviewSvc: overviewSvc}
}

// GetOverview godoc
// @Summary     Lấy tổng quan số lượng và trạng thái của tất cả trạm/thiết bị
// @Tags        overview
// @Produce     json
// @Success     200 {object} smpb.OverviewResponse
// @Router      /api/v1/overview [get]
func (h *OverviewHandler) GetOverview(c *gin.Context) {
	resp, err := h.overviewSvc.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
