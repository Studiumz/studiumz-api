package chat

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrChatDoesNotExist = errors.New("Chat does not exist")
)

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
  ON CONFLICT ON CONSTRAINT chats_user_ids_unique
  DO NOTHING
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
		if err.Error() == "scanning one: no rows in result set" {
			newChat, err = findChatByUserIds(ctx, tx, chat.FirstUserId, chat.SecondUserId)
			if err != nil {
				return
			}

			return newChat, nil
		}
	}

	return newChat, nil
}
