CREATE TABLE IF NOT EXISTS chats (
  id BYTEA NOT NULL,
  first_user_id BYTEA NOT NULL,
  second_user_id BYTEA NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS chats_user_ids_unique ON chats(GREATEST(first_user_id, second_user_id), LEAST(first_user_id, second_user_id), deleted_at) NULLS NOT DISTINCT 
WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS chat_messages (
  id BYTEA NOT NULL,
  chat_id BYTEA NOT NULL,
  from_user_id BYTEA NOT NULL,
  text TEXT,
  file_url TEXT, 
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMPTZ,
  
  PRIMARY KEY(id)
);
