CREATE TABLE IF NOT EXISTS roles (
	id serial4,
    name varchar(50) NOT NULL,
    label varchar(50) DEFAULT NULL::character varying,
	parent_id int4 DEFAULT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz,
    deleted_at timestamptz,
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES roles (id)
);