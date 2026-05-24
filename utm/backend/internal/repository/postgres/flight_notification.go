package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type FlightNotificationRepository struct {
	db *pgxpool.Pool
}

func NewFlightNotificationRepository(db *pgxpool.Pool) *FlightNotificationRepository {
	return &FlightNotificationRepository{db: db}
}

const fnCols = `id, operator_id, drone_id, uav_type, lane_id, v_waypoints,
	destination_vertiport, t_des, t_takeoff, t_land_estimated,
	priority, soc_0, c_bat, p_hover, status, created_at, updated_at`

func scanFlightNotification(row interface{ Scan(...any) error }) (*model.FlightNotification, error) {
	fn := &model.FlightNotification{}
	return fn, row.Scan(
		&fn.ID, &fn.OperatorID, &fn.DroneID, &fn.UavType, &fn.LaneID, &fn.VWaypoints,
		&fn.DestinationVertiport, &fn.TDes, &fn.TTakeoff, &fn.TLandEstimated,
		&fn.Priority, &fn.Soc0, &fn.CBat, &fn.PHover, &fn.Status,
		&fn.CreatedAt, &fn.UpdatedAt,
	)
}

func (r *FlightNotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.FlightNotification, error) {
	row := r.db.QueryRow(ctx, `SELECT `+fnCols+` FROM flight_notifications WHERE id = $1`, id)
	return scanFlightNotification(row)
}

func (r *FlightNotificationRepository) List(ctx context.Context, operatorID string) ([]*model.FlightNotification, error) {
	query := `SELECT ` + fnCols + ` FROM flight_notifications`
	args := []any{}
	if operatorID != "" {
		query += ` WHERE operator_id = $1`
		args = append(args, operatorID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.FlightNotification
	for rows.Next() {
		fn, err := scanFlightNotification(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, fn)
	}
	return list, nil
}

// ListActiveWithApprovedPlan trả về các notification SOFT_RESERVED/COMMITTED
// kèm approved_plan — dùng để build SystemStateProto.approved_plans cho SCRP
func (r *FlightNotificationRepository) ListActiveWithApprovedPlan(ctx context.Context) ([]*model.FlightNotification, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+fnCols+`
		FROM flight_notifications
		WHERE status IN ('soft_reserved', 'committed')
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.FlightNotification
	for rows.Next() {
		fn, err := scanFlightNotification(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, fn)
	}
	return list, nil
}

func (r *FlightNotificationRepository) Create(ctx context.Context, fn *model.FlightNotification) (*model.FlightNotification, error) {
	fn.ID = uuid.New()
	fn.Status = model.FlightStatusPending
	_, err := r.db.Exec(ctx, `
		INSERT INTO flight_notifications
		(id, operator_id, drone_id, uav_type, lane_id, v_waypoints,
		 destination_vertiport, t_des, t_takeoff, t_land_estimated,
		 priority, soc_0, c_bat, p_hover, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`,
		fn.ID, fn.OperatorID, fn.DroneID, fn.UavType, fn.LaneID, fn.VWaypoints,
		fn.DestinationVertiport, fn.TDes, fn.TTakeoff, fn.TLandEstimated,
		fn.Priority, fn.Soc0, fn.CBat, fn.PHover, fn.Status,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, fn.ID)
}

func (r *FlightNotificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.FlightNotificationStatus) error {
	_, err := r.db.Exec(ctx,
		`UPDATE flight_notifications SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, id,
	)
	return err
}

// --- ApprovedPlan ---

func (r *FlightNotificationRepository) CreateApprovedPlan(ctx context.Context, ap *model.ApprovedPlan) (*model.ApprovedPlan, error) {
	ap.ID = uuid.New()
	wtJSON, err := json.Marshal(ap.WaypointTimes)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO approved_plans
		(id, notification_id, t_dep_star, t_land_assigned, pad_id, slot_index,
		 waypoint_times, expires_at, delay_seconds, delay_source, energy_estimate, soc_meaning)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`,
		ap.ID, ap.NotificationID, ap.TDepStar, ap.TLandAssigned, ap.PadID, ap.SlotIndex,
		wtJSON, ap.ExpiresAt, ap.DelaySeconds, ap.DelaySource, ap.EnergyEstimate, ap.SocMeaning,
	)
	return ap, err
}

func (r *FlightNotificationRepository) FindApprovedPlanByNotification(ctx context.Context, notificationID uuid.UUID) (*model.ApprovedPlan, error) {
	ap := &model.ApprovedPlan{}
	var wtJSON []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, notification_id, t_dep_star, t_land_assigned, pad_id, slot_index,
		       waypoint_times, expires_at, delay_seconds, delay_source, energy_estimate, soc_meaning, created_at
		FROM approved_plans WHERE notification_id = $1
	`, notificationID).Scan(
		&ap.ID, &ap.NotificationID, &ap.TDepStar, &ap.TLandAssigned, &ap.PadID, &ap.SlotIndex,
		&wtJSON, &ap.ExpiresAt, &ap.DelaySeconds, &ap.DelaySource, &ap.EnergyEstimate, &ap.SocMeaning, &ap.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ap, json.Unmarshal(wtJSON, &ap.WaypointTimes)
}
