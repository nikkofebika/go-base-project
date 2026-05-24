CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    last_login_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by_id BIGINT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by_id BIGINT NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by_id BIGINT NULL
);

-- indexes
CREATE INDEX idx_users_name ON users(name);

-- CREATE INDEX idx_users_email ON users(email); // unique already indexed