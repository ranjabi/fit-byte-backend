BEGIN;

CREATE TABLE users (
	id 				uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email			text NOT NULL UNIQUE,
	password		text NOT NULL,
	preference		preference,
	weight_unit		weight_unit,
	height_unit		height_unit,
	weight			integer,
	height			integer,
	name			text,
	image_uri		text
);

COMMIT;