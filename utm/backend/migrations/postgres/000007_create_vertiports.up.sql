CREATE TABLE vertiports (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name              VARCHAR(100) NOT NULL,
    slot_duration_sec DOUBLE PRECISION NOT NULL DEFAULT 60,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pads (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vertiport_id     UUID NOT NULL REFERENCES vertiports(id) ON DELETE CASCADE,
    name             VARCHAR(50) NOT NULL,
    compatible_types TEXT[] NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pads_vertiport_id ON pads(vertiport_id);

CREATE TRIGGER vertiports_updated_at
    BEFORE UPDATE ON vertiports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER pads_updated_at
    BEFORE UPDATE ON pads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
