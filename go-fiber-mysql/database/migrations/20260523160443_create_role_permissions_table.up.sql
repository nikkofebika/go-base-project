CREATE TABLE role_permissions (
  role_id INT UNSIGNED NOT NULL,
  permission_id INT UNSIGNED NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_id INT NULL,

  PRIMARY KEY (role_id, permission_id),

  CONSTRAINT fk_role_permissions_role
    FOREIGN KEY (role_id)
    REFERENCES roles(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_role_permissions_permission
    FOREIGN KEY (permission_id)
    REFERENCES permissions(id)
    ON DELETE CASCADE,

  INDEX idx_role_permissions_permission_id (permission_id)
)