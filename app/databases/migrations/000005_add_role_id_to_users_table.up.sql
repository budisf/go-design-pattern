ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id INTEGER;

ALTER TABLE users ADD FOREIGN KEY (role_id) REFERENCES roles (id);