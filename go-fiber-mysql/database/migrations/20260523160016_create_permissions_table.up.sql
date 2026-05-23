CREATE TABLE permissions (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL UNIQUE,
  slug VARCHAR(50) NOT NULL UNIQUE,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_id INT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_by_id INT NULL,
  deleted_at TIMESTAMP NULL,
  deleted_by_id INT NULL,
  
  -- indexes
  INDEX idx_permissions_name (name),
  INDEX idx_permissions_slug (slug)
)