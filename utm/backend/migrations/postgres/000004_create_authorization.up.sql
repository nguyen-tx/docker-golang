CREATE TYPE auth_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled', 'expired', 'revoked');
CREATE TYPE flight_purpose AS ENUM ('commercial', 'recreational', 'survey', 'delivery', 'emergency', 'research');

CREATE TABLE authorization_requests (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_no       VARCHAR(30) NOT NULL UNIQUE,   -- UTM-20241125-0001
    pilot_id         UUID NOT NULL REFERENCES users(id),
    aircraft_id      UUID NOT NULL REFERENCES aircraft(id),
    zone_id          UUID REFERENCES airspace_zones(id),
    purpose          flight_purpose NOT NULL,
    status           auth_status NOT NULL DEFAULT 'pending',
    flight_area_json TEXT NOT NULL,                  -- GeoJSON Polygon khu vực bay
    planned_alt_m    INT NOT NULL,
    planned_start_at TIMESTAMPTZ NOT NULL,
    planned_end_at   TIMESTAMPTZ NOT NULL,
    pilot_notes      TEXT,
    reviewer_id      UUID REFERENCES users(id),
    reviewed_at      TIMESTAMPTZ,
    review_notes     TEXT,
    valid_from       TIMESTAMPTZ,
    valid_until      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_planned_time CHECK (planned_end_at > planned_start_at)
);

CREATE INDEX idx_auth_pilot_id ON authorization_requests(pilot_id);
CREATE INDEX idx_auth_status ON authorization_requests(status);
CREATE INDEX idx_auth_aircraft_id ON authorization_requests(aircraft_id);
CREATE INDEX idx_auth_planned_start ON authorization_requests(planned_start_at);

CREATE TRIGGER auth_updated_at
    BEFORE UPDATE ON authorization_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
