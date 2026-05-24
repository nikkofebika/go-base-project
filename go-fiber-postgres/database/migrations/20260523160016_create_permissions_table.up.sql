CREATE TABLE permissions (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  slug CITEXT NOT NULL UNIQUE,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_id BIGINT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_by_id BIGINT NULL,
  deleted_at TIMESTAMP NULL,
  deleted_by_id BIGINT NULL
);

-- indexes
CREATE INDEX idx_permissions_name ON permissions(name);
-- CREATE INDEX idx_permissions_slug ON permissions(slug);