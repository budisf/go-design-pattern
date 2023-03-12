CREATE TYPE sales_zone_type AS ENUM ('group_territories', 'areas', 'regions');

CREATE TABLE IF NOT EXISTS user_zones (
	id bigserial NOT NULL,
	user_id int8 NULL,
	sales_zone_id int8 NULL,
	"sales_zone_type" sales_zone_type DEFAULT NULL,
	assigned_date timestamptz NULL,
	finished_date timestamptz NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);



