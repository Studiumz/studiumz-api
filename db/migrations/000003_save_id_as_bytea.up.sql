ALTER TABLE IF EXISTS users
ALTER COLUMN id TYPE BYTEA USING id::BYTEA;

ALTER TABLE IF EXISTS accounts
ALTER COLUMN id TYPE BYTEA USING id::BYTEA,
ALTER COLUMN user_id TYPE BYTEA USING user_id::BYTEA;
