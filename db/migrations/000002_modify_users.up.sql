ALTER TABLE users
RENAME COLUMN bio TO struggles;

ALTER TABLE users
RENAME COLUMN image_url TO avatar;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS gender INTEGER;
