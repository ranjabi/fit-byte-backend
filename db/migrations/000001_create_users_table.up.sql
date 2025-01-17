BEGIN;

CREATE TABLE users (
	id 					uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email				text NOT NULL UNIQUE,
	password			text NOT NULL
);

COMMIT;