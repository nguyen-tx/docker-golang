package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/utm/backend/internal/model"
	scrpclient "github.com/utm/backend/internal/scrp"
	pgRepo "github.com/utm/backend/internal/repository/postgres"
	scrppb "github.com/utm/backend/pkg/pb/scrp"
)

type FlightNotificationService struct {
	fnRepo       *pgRepo.FlightNotificationRepository
	laneRepo     *pgRepo.LaneRepository
	waypointRepo *pgRepo.WaypointRepository
	vertiportRepo *pgRepo.VertiportRepository
	scrpClient   *scrpclient.Client
}

func NewFlightNotificationService(
	fnRepo *pgRepo.FlightNotificationRepository,
	laneRepo *pgRepo.LaneRepository,
	waypointRepo *pgRepo.WaypointRepository,
	vertiportRepo *pgRepo.VertiportRepository,
	scrpClient *scrpclient.Client,
) *FlightNotificationService {
	return &FlightNotificationService{
		fnRepo:        fnRepo,
		laneRepo:      laneRepo,
		waypointRepo:  waypointRepo,
		vertiportRepo: vertiportRepo,
		scrpClient:    scrpClient,
	}
}

func (s *FlightNotificationService) List(ctx context.Context, operatorID string) ([]*model.FlightNotification, error) {
	return s.fnRepo.List(ctx, operatorID)
}

func (s *FlightNotificationService) GetByID(ctx context.Context, id string) (*model.FlightNotification, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.fnRepo.FindByID(ctx, uid)
}

// Submit tạo thông báo bay mới, gọi SCRP đồng bộ và trả về kết quả
func (s *FlightNotificationService) Submit(ctx context.Context, operatorID string, fn *model.FlightNotification) (*model.FlightNotification, *model.ApprovedPlan, *scrppb.RejectResultProto, error) {
	// 1. parse operator ID
	oid, err := uuid.Parse(operatorID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid operator_id: %w", err)
	}
	fn.OperatorID = oid

	// 2. lưu flight notification với status pending
	fn, err = s.fnRepo.Create(ctx, fn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create flight notification: %w", err)
	}

	// 3. build SCRP request
	req, err := s.buildResolveConflictRequest(ctx, fn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("build scrp request: %w", err)
	}

	// 4. gọi SCRP
	resp, err := s.scrpClient.ResolveConflict(ctx, req)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scrp resolve conflict: %w", err)
	}

	// 5. xử lý kết quả
	switch r := resp.Result.(type) {
	case *scrppb.ResolveConflictResponse_Approve:
		ap, err := s.handleApprove(ctx, fn.ID, r.Approve)
		if err != nil {
			return nil, nil, nil, err
		}
		return fn, ap, nil, nil

	case *scrppb.ResolveConflictResponse_Reject:
		if err := s.fnRepo.UpdateStatus(ctx, fn.ID, model.FlightStatusUTMRejected); err != nil {
			return nil, nil, nil, err
		}
		fn.Status = model.FlightStatusUTMRejected
		return fn, nil, r.Reject, nil
	}

	return nil, nil, nil, fmt.Errorf("unexpected scrp response")
}

// handleApprove lưu ApprovedPlan và cập nhật status → SOFT_RESERVED
func (s *FlightNotificationService) handleApprove(ctx context.Context, notifID uuid.UUID, res *scrppb.ApproveResultProto) (*model.ApprovedPlan, error) {
	wt := make([]model.WaypointTimeEntry, len(res.WaypointTimes))
	for i, w := range res.WaypointTimes {
		wt[i] = model.WaypointTimeEntry{WaypointID: w.WaypointId, Time: w.Time}
	}

	ap := &model.ApprovedPlan{
		NotificationID: notifID,
		TDepStar:       res.TDepStar,
		TLandAssigned:  res.TLandAssigned,
		PadID:          res.PadId,
		SlotIndex:      res.SlotIndex,
		WaypointTimes:  wt,
		ExpiresAt:      res.ExpiresAt,
		DelaySeconds:   res.DelaySeconds,
		DelaySource:    res.DelaySource,
		EnergyEstimate: res.EnergyEstimate,
		SocMeaning:     res.SocMeaning,
	}

	ap, err := s.fnRepo.CreateApprovedPlan(ctx, ap)
	if err != nil {
		return nil, fmt.Errorf("save approved plan: %w", err)
	}

	if err := s.fnRepo.UpdateStatus(ctx, notifID, model.FlightStatusSoftReserved); err != nil {
		return nil, err
	}
	return ap, nil
}

// buildResolveConflictRequest tổng hợp toàn bộ dữ liệu cần thiết cho SCRP
func (s *FlightNotificationService) buildResolveConflictRequest(ctx context.Context, fn *model.FlightNotification) (*scrppb.ResolveConflictRequest, error) {
	// load lane + waypoints của flight notification hiện tại
	laneProto, err := s.loadLaneProto(ctx, fn.LaneID)
	if err != nil {
		return nil, fmt.Errorf("load lane: %w", err)
	}

	// load vertiport hạ cánh để build VertiportStateProto
	vertiportProto, err := s.buildVertiportStateProto(ctx, fn.DestinationVertiport)
	if err != nil {
		return nil, fmt.Errorf("load vertiport: %w", err)
	}

	// load approved plans hiện tại (SOFT_RESERVED + COMMITTED)
	approvedPlans, err := s.buildApprovedPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("load approved plans: %w", err)
	}

	return &scrppb.ResolveConflictRequest{
		FlightIntention: &scrppb.OperatorFlightRequestProto{
			OperatorId:              fn.OperatorID.String(),
			DroneId:                 fn.DroneID,
			UavType:                 fn.UavType,
			Lane:                    laneProto,
			VWaypoints:              fn.VWaypoints,
			DestinationVertiportId:  fn.DestinationVertiport,
			TDes:                    fn.TDes,
			Priority:                fn.Priority,
			Soc_0:                   fn.Soc0,
			CBat:                    fn.CBat,
			PHover:                  fn.PHover,
			TTakeoff:                fn.TTakeoff,
			TLandEstimated:          fn.TLandEstimated,
		},
		ApprovedPlans:  approvedPlans,
		VertiportState: vertiportProto,
		SystemState: &scrppb.SystemStateProto{
			ApprovedPlans:  approvedPlans,
			VertiportState: vertiportProto,
			TNow:           float64(time.Now().Unix()),
		},
	}, nil
}

// loadLaneProto load Lane từ DB và build LaneProto đầy đủ với waypoint data
func (s *FlightNotificationService) loadLaneProto(ctx context.Context, laneID uuid.UUID) (*scrppb.LaneProto, error) {
	lane, err := s.laneRepo.FindByID(ctx, laneID)
	if err != nil {
		return nil, err
	}

	waypointMap, err := s.waypointRepo.FindByIDs(ctx, lane.WaypointIDs)
	if err != nil {
		return nil, err
	}

	wpProtos := make([]*scrppb.WaypointProto, 0, len(lane.WaypointIDs))
	for _, wid := range lane.WaypointIDs {
		w, ok := waypointMap[wid]
		if !ok {
			return nil, fmt.Errorf("waypoint %s not found", wid)
		}
		wpProtos = append(wpProtos, &scrppb.WaypointProto{
			Id: w.ID.String(), X: w.X, Y: w.Y, Z: w.Z,
		})
	}

	segProtos := make([]*scrppb.SegmentProto, 0, len(lane.Segments))
	for _, seg := range lane.Segments {
		fromW, ok1 := waypointMap[seg.FromWaypointID]
		toW, ok2 := waypointMap[seg.ToWaypointID]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("segment waypoint not found")
		}
		segProtos = append(segProtos, &scrppb.SegmentProto{
			PStart: &scrppb.WaypointProto{Id: fromW.ID.String(), X: fromW.X, Y: fromW.Y, Z: fromW.Z},
			PEnd:   &scrppb.WaypointProto{Id: toW.ID.String(), X: toW.X, Y: toW.Y, Z: toW.Z},
			Length: seg.Length,
			VMin:   seg.VMin,
			VMax:   seg.VMax,
		})
	}

	return &scrppb.LaneProto{
		Id:       lane.ID.String(),
		Waypoints: wpProtos,
		Segments:  segProtos,
	}, nil
}

// buildVertiportStateProto load Vertiport từ DB
func (s *FlightNotificationService) buildVertiportStateProto(ctx context.Context, vertiportID string) (*scrppb.VertiportStateProto, error) {
	vid, err := uuid.Parse(vertiportID)
	if err != nil {
		return nil, fmt.Errorf("invalid vertiport_id: %w", err)
	}

	v, err := s.vertiportRepo.FindByID(ctx, vid)
	if err != nil {
		return nil, err
	}

	padProtos := make([]*scrppb.PadProto, 0, len(v.Pads))
	for _, p := range v.Pads {
		padProtos = append(padProtos, &scrppb.PadProto{
			Id:              p.ID.String(),
			CompatibleTypes: p.CompatibleTypes,
		})
	}

	return &scrppb.VertiportStateProto{
		VertiportId:  v.ID.String(),
		Pads:         padProtos,
		SlotDuration: v.SlotDurationSec,
		Slots:        []*scrppb.SlotEntryProto{}, // slot booking sẽ bổ sung sau
	}, nil
}

// buildApprovedPlans lấy các notification SOFT_RESERVED/COMMITTED và build ApprovedPlanProto
func (s *FlightNotificationService) buildApprovedPlans(ctx context.Context) ([]*scrppb.ApprovedPlanProto, error) {
	notifications, err := s.fnRepo.ListActiveWithApprovedPlan(ctx)
	if err != nil {
		return nil, err
	}

	plans := make([]*scrppb.ApprovedPlanProto, 0, len(notifications))
	for _, fn := range notifications {
		ap, err := s.fnRepo.FindApprovedPlanByNotification(ctx, fn.ID)
		if err != nil {
			continue // bỏ qua nếu chưa có approved plan
		}

		laneProto, err := s.loadLaneProto(ctx, fn.LaneID)
		if err != nil {
			continue
		}

		wtProtos := make([]*scrppb.WaypointTimeEntryProto, 0, len(ap.WaypointTimes))
		for _, wt := range ap.WaypointTimes {
			wtProtos = append(wtProtos, &scrppb.WaypointTimeEntryProto{
				Waypoint: &scrppb.WaypointProto{Id: wt.WaypointID},
				Time:     wt.Time,
			})
		}

		plans = append(plans, &scrppb.ApprovedPlanProto{
			DroneId:       fn.DroneID,
			UavType:       fn.UavType,
			Lane:          laneProto,
			VWaypoints:    fn.VWaypoints,
			TLand:         ap.TLandAssigned,
			PadId:         ap.PadID,
			SlotIndex:     ap.SlotIndex,
			Status:        string(fn.Status),
			WaypointTimes: wtProtos,
			ExpiresAt:     ap.ExpiresAt,
		})
	}

	return plans, nil
}
