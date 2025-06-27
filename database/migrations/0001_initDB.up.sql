BEGIN;

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP  DEFAULT NULL
);

-- Roles Table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name TEXT UNIQUE NOT NULL
);

-- User_Roles Mapping Table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID REFERENCES users(id),
    role_id UUID REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- Addresses Table
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    label TEXT,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    archived_at TIMESTAMP DEFAULT NULL
);

-- Restaurants Table
CREATE TABLE IF NOT EXISTS restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurantname TEXT NOT NULL,
    created_by UUID REFERENCES users(id),
    archived_at TIMESTAMP DEFAULT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL
);

-- Dishes Table
CREATE TABLE IF NOT EXISTS dishes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dishname TEXT NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id),
    price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    created_by UUID REFERENCES users(id),
    archived_at TIMESTAMP DEFAULT NULL
    );

-- Sessions Table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    refresh_token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pre-populate roles
INSERT INTO roles (role_name)
SELECT 'admin'
    WHERE NOT EXISTS (SELECT 1 FROM roles WHERE role_name = 'admin');

INSERT INTO roles (role_name)
SELECT 'subadmin'
    WHERE NOT EXISTS (SELECT 1 FROM roles WHERE role_name = 'subadmin');

INSERT INTO roles (role_name)
SELECT 'user'
    WHERE NOT EXISTS (SELECT 1 FROM roles WHERE role_name = 'user');

-- Insert hardcoded admin user only if not exists
INSERT INTO users (id, username, email, password, created_by)
SELECT 'd4e5f6c2-1234-4abc-9def-111122223333',
       'AdminUser',
       'admin@example.com',
       '$2a$10$TV0sEg/w1/GPMeIPSO2t0.pnFsxEEhNbu528/UWGPoj1SGeSuwodC',  -- hashed "admin123"
       NULL
    WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@example.com'

);

-- Assign 'admin' role to the hardcoded admin user (only if not already assigned)
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'admin@example.com'
  AND r.role_name = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = r.id
);


COMMIT;
