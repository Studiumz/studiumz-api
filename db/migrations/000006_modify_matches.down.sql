ALTER TABLE IF EXISTS matches
ALTER COLUMN id TYPE CHAR(26);
ALTER TABLE IF EXISTS matches
ALTER COLUMN matcher_id TYPE CHAR(26);
ALTER TABLE IF EXISTS matches
ALTER COLUMN matchee_id TYPE CHAR(26);