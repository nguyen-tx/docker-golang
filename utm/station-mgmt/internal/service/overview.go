package service

import (
	"context"
	"time"

	"github.com/utm/station-mgmt/internal/client"
	mipb "github.com/utm/station-mgmt/pkg/pb/manage_info/v1"
	nspb "github.com/utm/station-mgmt/pkg/pb/network_station/v1"
	rtkpb "github.com/utm/station-mgmt/pkg/pb/rtk_station/v1"
	smpb "github.com/utm/station-mgmt/pkg/pb/station_mgmt/v1"
	"golang.org/x/sync/errgroup"
)

type OverviewService struct {
	ns  *client.NetworkStationClient
	rtk *client.RTKStationClient
	mi  *client.ManageInfoClient
}

func NewOverviewService(
	ns *client.NetworkStationClient,
	rtk *client.RTKStationClient,
	mi *client.ManageInfoClient,
) *OverviewService {
	return &OverviewService{ns: ns, rtk: rtk, mi: mi}
}

func (s *OverviewService) GetOverview(ctx context.Context) (*smpb.OverviewResponse, error) {
	var (
		nsSummary  *nspb.StationSummary
		rtkSummary *rtkpb.StationSummary
		miSummary  *mipb.ManageInfoSummary
	)

	// gọi 3 service song song, fail nhanh nếu bất kỳ service nào lỗi
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		nsSummary, err = s.ns.GetSummary(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		rtkSummary, err = s.rtk.GetSummary(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		miSummary, err = s.mi.GetSummary(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &smpb.OverviewResponse{
		FiveG: &smpb.GroupSummary{
			Total:       nsSummary.Total,
			Online:      nsSummary.Online,
			Offline:     nsSummary.Offline,
			Maintenance: nsSummary.Maintenance,
		},
		Rtk: &smpb.GroupSummary{
			Total:       rtkSummary.Total,
			Online:      rtkSummary.Online,
			Offline:     rtkSummary.Offline,
			Maintenance: rtkSummary.Maintenance,
		},
		Radar: &smpb.GroupSummary{
			Total:       miSummary.Radar.Total,
			Online:      miSummary.Radar.Online,
			Offline:     miSummary.Radar.Offline,
			Maintenance: miSummary.Radar.Maintenance,
		},
		Tdoa: &smpb.GroupSummary{
			Total:       miSummary.Tdoa.Total,
			Online:      miSummary.Tdoa.Online,
			Offline:     miSummary.Tdoa.Offline,
			Maintenance: miSummary.Tdoa.Maintenance,
		},
		Hub: &smpb.GroupSummary{
			Total:       miSummary.Hub.Total,
			Online:      miSummary.Hub.Online,
			Offline:     miSummary.Hub.Offline,
			Maintenance: miSummary.Hub.Maintenance,
		},
		FetchedAt: time.Now().UnixMilli(),
	}, nil
}
