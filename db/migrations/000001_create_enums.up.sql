BEGIN;

CREATE TYPE preference AS ENUM ('CARDIO', 'WEIGHT');
CREATE TYPE weight_unit AS ENUM ('KG', 'LBS');
CREATE TYPE height_unit AS ENUM ('CM', 'INCH');
CREATE TYPE activity_type AS ENUM (
    'Walking',
    'Yoga',
    'Stretching',
    'Cycling',
    'Swimming',
    'Dancing',
    'Hiking',
    'Running',
    'HIIT',
    'JumpRope'
);

COMMIT;