CREATE TABLE IF NOT EXISTS matches (
  id CHAR(26) NOT NULL,
  matcher_id CHAR(26) NOT NULL,
  matchee_id CHAR(26) NOT NULL,
  status VARCHAR(10),
  invitation_message TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id),
  CONSTRAINT match_user_combination_unique UNIQUE NULLS NOT DISTINCT (matcher_id, matchee_id, deleted_at)
);