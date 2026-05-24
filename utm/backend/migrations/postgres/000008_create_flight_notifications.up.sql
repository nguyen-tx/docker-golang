CREATE TYPE flight_notification_status AS ENUM (
    'pending',
    'soft_reserved',
    'committed',
    'utm_rejected',
    'dcs_rejected'
);

CREATE TABLE flight_notifications (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    operator_id           UUID NOT NULL REFERENCES users(id),
    drone_id              VARCHAR(100) NOT NULL,
    uav_type              VARCHAR(50) NOT NULL,
    lane_id               UUID NOT NULL REFERENCES lanes(id),
    v_waypoints           DOUBLE PRECISION[] NOT NULL DEFAULT '{}',
    destination_vertiport VARCHAR(100) NOT NULL,
    t_des                 DOUBLE PRECISION NOT NULL,
    t_takeoff             DOUBLE PRECISION NOT NULL,
    t_land_estimated      DOUBLE PRECISION NOT NULL,
    priority              INT NOT NULL DEFAULT 0,
    soc_0                 DOUBLE PRECISION NOT NULL DEFAULT 100,
    c_bat                 DOUBLE PRECISION NOT NULL DEFAULT 0,
    p_hover               DOUBLE PRECISION NOT NULL DEFAULT 0,
    status                flight_notification_status NOT NULL DEFAULT 'pending',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fn_operator_id ON flight_notifications(operator_id);
CREATE INDEX idx_fn_status      ON flight_notifications(status);
CREATE INDEX idx_fn_t_takeoff   ON flight_notifications(t_takeoff);

CREATE TRIGGER flight_notifications_updated_at
    BEFORE UPDATE ON flight_notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- approved_plans: lưu kết quả SCRP approve, độc lập với flight_notifications
CREATE TABLE approved_plans (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_id UUID NOT NULL REFERENCES flight_notifications(id) ON DELETE CASCADE,
    t_dep_star      DOUBLE PRECISION NOT NULL,
    t_land_assigned DOUBLE PRECISION NOT NULL,
    pad_id          VARCHAR(100) NOT NULL,
    slot_index      INT NOT NULL,
    waypoint_times  JSONB NOT NULL DEFAULT '[]',
    expires_at      DOUBLE PRECISION NOT NULL,
    delay_seconds   DOUBLE PRECISION NOT NULL DEFAULT 0,
    delay_source    VARCHAR(100) NOT NULL DEFAULT '',
    energy_estimate DOUBLE PRECISION NOT NULL DEFAULT 0,
    soc_meaning     DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_approved_plan_notification UNIQUE (notification_id)
);

CREATE INDEX idx_ap_notification_id ON approved_plans(notification_id);
