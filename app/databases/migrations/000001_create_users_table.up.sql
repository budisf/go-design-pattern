CREATE TABLE IF NOT EXISTS users (
	id serial4,
    auth_server_id int4 DEFAULT NULL,
    nip varchar(18) NOT NULL,
	name varchar(255) DEFAULT NULL::character varying,
	status varchar(255) DEFAULT NULL::character varying,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz,
    deleted_at timestamptz,
    PRIMARY KEY (id)
);

