CREATE TYPE aircraft_type AS ENUM ('fixed_wing', 'multi_rotor', 'helicopter', 'hybrid');
CREATE TYPE aircraft_status AS ENUM ('active', 'inactive', 'maintenance', 'grounded');

CREATE TABLE aircraft (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registration_no     VARCHAR(50) NOT NULL UNIQUE,
    model               VARCHAR(100) NOT NULL,
    manufacturer        VARCHAR(100) NOT NULL,
    serial_no           VARCHAR(100) NOT NULL UNIQUE,
    type                aircraft_type NOT NULL,
    status              aircraft_status NOT NULL DEFAULT 'active',
    weight_kg           DECIMAL(8,3) NOT NULL,
    max_altitude_m      INT NOT NULL DEFAULT 120,
    max_speed_kmh       INT NOT NULL DEFAULT 50,
    max_flight_time_min INT NOT NULL DEFAULT 30,
    has_camera          BOOLEAN NOT NULL DEFAULT FALSE,
    insurance_no        VARCHAR(100),
    insurance_expiry    DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aircraft_owner_id ON aircraft(owner_id);
CREATE INDEX idx_aircraft_status ON aircraft(status);
CREATE INDEX idx_aircraft_registration ON aircraft(registration_no);

CREATE TRIGGER aircraft_updated_at
    BEFORE UPDATE ON aircraft
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
