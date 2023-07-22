CREATE TABLE IF NOT EXISTS subjects (
  id BYTEA NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id),
  CONSTRAINT name_unique UNIQUE NULLS NOT DISTINCT (name, deleted_at)
);

CREATE TABLE IF NOT EXISTS users_subjects (
  id BYTEA NOT NULL,
  user_id BYTEA NOT NULL,
  subject_id BYTEA NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id, user_id),
  CONSTRAINT user_subject_unique UNIQUE NULLS NOT DISTINCT (user_id, subject_id, deleted_at)
);