CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('admin', 'pilot', 'operator', 'atc', 'inspector');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending', 'banned');

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255) NOT NULL,
    phone         VARCHAR(20),
    role          user_role NOT NULL DEFAULT 'pilot',
    status        user_status NOT NULL DEFAULT 'pending',
    license_no    VARCHAR(100),
    organization  VARCHAR(255),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);

-- Trigger tự động cập nhật updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Admin mặc định (password: Admin@123)
INSERT INTO users (id, email, password_hash, full_name, role, status)
VALUES (
    uuid_generate_v4(),
    'admin@utm.vn',
    '$2a$10$rHzHN0X3aG6vIh2XhJv8heADcaQiKx5pYgGqNjMg6i7QbZ7ZQDK9m',
    'System Administrator',
    'admin',
    'active'
);
