BEGIN;

-- user id
-- 9bfc3585-e92d-4506-917d-ed9eb0bfb13b
-- user password 12345678
-- $2a$10$ThYUBp8mOhpXWNaKMSWnZ.mZKBUq82l8/KbWcTsBYyjz4qHXCVuSe

INSERT INTO users (
    email, 
    password, 
    preference, 
    weight_unit, 
    height_unit, 
    weight, 
    height, 
    name, 
    image_uri
) VALUES (
    'a@a.a',
    '$2a$10$ThYUBp8mOhpXWNaKMSWnZ.mZKBUq82l8/KbWcTsBYyjz4qHXCVuSe',
    'CARDIO',
    'KG',
    'CM',
    75,
    180,
    'John Doe',
    'https://example.com/john.jpg'
);

COMMIT;