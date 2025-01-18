BEGIN;

CREATE TABLE activities (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_type       activity_type NOT NULL,
    done_at             timestamp NOT NULL,
    duration_in_minutes integer NOT NULL,
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;