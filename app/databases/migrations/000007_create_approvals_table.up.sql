CREATE TYPE approval_type AS ENUM ('dpl', 'dpf', 'pssp','entertain');

CREATE TABLE IF NOT EXISTS approvals (
	id bigserial NOT NULL,
	user_id int8 NULL,
	"level" int8 NULL,
	approval_type approval_type DEFAULT NULL,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
	UNIQUE ("level", "approval_type")
);