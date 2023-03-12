CREATE TABLE IF NOT EXISTS distributors (
	id serial4 NOT NULL,
	distributor_name varchar(256) NULL,
	is_deleted bool DEFAULT false,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	PRIMARY KEY (id)
);