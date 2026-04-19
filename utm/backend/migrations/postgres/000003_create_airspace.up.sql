CREATE TYPE zone_type AS ENUM ('restricted', 'controlled', 'uncontrolled', 'tfr', 'nfz');

CREATE TABLE airspace_zones (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    type             zone_type NOT NULL,
    min_altitude_m   INT NOT NULL DEFAULT 0,
    max_altitude_m   INT NOT NULL DEFAULT 500,
    boundary_geojson TEXT NOT NULL,       -- GeoJSON Polygon
    active_from      TIMESTAMPTZ,
    active_until     TIMESTAMPTZ,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_by       UUID NOT NULL REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_airspace_type ON airspace_zones(type);
CREATE INDEX idx_airspace_is_active ON airspace_zones(is_active);

CREATE TRIGGER airspace_updated_at
    BEFORE UPDATE ON airspace_zones
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
