CREATE TABLE waypoints (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    x          DOUBLE PRECISION NOT NULL,
    y          DOUBLE PRECISION NOT NULL,
    z          DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER waypoints_updated_at
    BEFORE UPDATE ON waypoints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
