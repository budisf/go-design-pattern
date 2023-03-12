CREATE TABLE IF NOT EXISTS regions (
	id serial4,
    name varchar(50) NOT NULL,
	is_deleted bool DEFAULT FALSE,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz,
    deleted_at timestamptz,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS areas (
	id serial4,
    name varchar(50) NOT NULL,
    region_id int4 not NULL,
	is_deleted bool DEFAULT FALSE,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz,
    deleted_at timestamptz,
    PRIMARY KEY (id),
    FOREIGN KEY (region_id) REFERENCES regions (id)
);

CREATE TABLE IF NOT EXISTS group_territories (
	id serial4,
    name varchar(50) NOT NULL,
    area_id int4 not NULL,
	is_deleted bool DEFAULT FALSE,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz,
    deleted_at timestamptz,
    PRIMARY KEY (id),
    FOREIGN KEY (area_id) REFERENCES areas (id)
);