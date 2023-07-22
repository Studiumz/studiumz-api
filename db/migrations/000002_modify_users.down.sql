ALTER TABLE users
DROP COLUMN IF EXISTS gender;

ALTER TABLE users
RENAME COLUMN avatar TO image_url;

ALTER TABLE users
RENAME COLUMN struggles TO bio;
