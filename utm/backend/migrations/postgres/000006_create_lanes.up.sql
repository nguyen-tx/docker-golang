CREATE TABLE lanes (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(100) NOT NULL,
    -- mảng UUID theo thứ tự các waypoint trong lane
    waypoint_ids UUID[] NOT NULL DEFAULT '{}',
    -- [{from_waypoint_id, to_waypoint_id, length, v_min, v_max}]
    segments     JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER lanes_updated_at
    BEFORE UPDATE ON lanes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
