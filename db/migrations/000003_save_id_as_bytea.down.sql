ALTER TABLE IF EXISTS users
ALTER COLUMN id TYPE CHAR(26);

ALTER TABLE IF EXISTS accounts
ALTER COLUMN id TYPE CHAR(26),
ALTER COLUMN user_id TYPE CHAR(26);
