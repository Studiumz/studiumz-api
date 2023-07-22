package chat

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrChatDoesNotExist  = errors.New("Chat does not exist")
	ErrChatAlreadyExists = errors.New("Chat already exists")
)

func findChatsByUserId(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (chats []*ChatViewModel, err error) {
	q := `
  SELECT cs.*, (CASE WHEN cm.created_at IS NULL THEN cs.created_at ELSE cm.created_at END) "last_msg_time"
  FROM (
    SELECT cs1.*, u.full_name "recipient_name", u.email "recipient_email"
    FROM (
      SELECT c.*
      FROM chats c
      INNER JOIN users u
      ON u.id = c.first_user_id
      WHERE c.first_user_id = $1
    ) cs1
    INNER JOIN users u
    ON u.id = cs1.second_user_id
    UNION
    SELECT cs2.*, u.full_name "recipient_name", u.email "recipient_email"
    FROM (
      SELECT c.*
      FROM chats c
      INNER JOIN users u
      ON u.id = c.second_user_id
      WHERE c.second_user_id = $1
    ) cs2
    INNER JOIN users u
    ON u.id = cs2.first_user_id
  ) cs
  LEFT JOIN chat_messages cm
  ON cm.chat_id = cs.id AND
  cm.created_at = (
    SELECT MAX(cm.created_at)
    FROM chat_messages cm
    INNER JOIN chats c
    ON c.id = cm.chat_id
  )
  ORDER BY last_msg_time DESC
  `

	chats = []*ChatViewModel{}
	if err = pgxscan.Select(ctx, tx, &chats, q, userId); err != nil {
		log.Err(err).Msg("Failed to find chats")
		return
	}

	return chats, nil
}

func findChatByUserIds(ctx context.Context, tx pgx.Tx, firstUserId, secondUserId ulid.ULID) (chat Chat, err error) {
	q := `
  SELECT *
  FROM chats
  WHERE ((first_user_id = $1 AND second_user_id = $2) OR (first_user_id = $2 AND second_user_id = $1)) AND deleted_at IS NULL
  `

	if err = pgxscan.Get(ctx, tx, &chat, q, firstUserId, secondUserId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return chat, ErrChatDoesNotExist
		}

		log.Err(err).Msg("Failed to find chat by user ids")
		return chat, err
	}

	return chat, nil
}

func saveChat(ctx context.Context, tx pgx.Tx, chat Chat) (newChat Chat, err error) {
	q := `
  INSERT INTO chats (id, first_user_id, second_user_id, created_at) VALUES
  ($1, $2, $3, $4)
  RETURNING *
  `
	if err = pgxscan.Get(
		ctx,
		tx,
		&newChat,
		q,
		chat.Id,
		chat.FirstUserId,
		chat.SecondUserId,
		chat.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return newChat, ErrChatAlreadyExists
			}

			log.Err(err).Msg("Failed to save chat")
			return newChat, err
		}

		log.Err(err).Msg("Failed to save chat")
		return
	}

	return newChat, nil
}
