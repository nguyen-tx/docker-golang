package gapi

import (
	"context"

	"github.com/utm/station-mgmt/internal/service"
	smpb "github.com/utm/station-mgmt/pkg/pb/station_mgmt/v1"
)

type StationMgmtHandler struct {
	smpb.UnimplementedStationManagementServiceServer
	overviewSvc *service.OverviewService
}

func NewStationMgmtHandler(overviewSvc *service.OverviewService) *StationMgmtHandler {
	return &StationMgmtHandler{overviewSvc: overviewSvc}
}

func (h *StationMgmtHandler) GetOverview(ctx context.Context, _ *smpb.GetOverviewRequest) (*smpb.OverviewResponse, error) {
	return h.overviewSvc.GetOverview(ctx)
}
