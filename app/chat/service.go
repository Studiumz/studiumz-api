package chat

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func getChats(ctx context.Context, userId ulid.ULID) (chats []*ChatViewModel, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get chats")
		return
	}

	defer tx.Rollback(ctx)

	chats, err = findChatsByUserId(ctx, tx, userId)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get chats")
		return
	}

	return chats, nil
}

func createChat(ctx context.Context, firstUserId ulid.ULID, body createChatReq) (chat Chat, errs map[string]error, err error) {
	secondUserId, err := validateChatUserId(body.SecondUserId)
	if err != nil {
		return
	}

	chat, errs = NewChat(firstUserId.String(), body.SecondUserId)
	if errs != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create chat")
		return
	}

	defer tx.Rollback(ctx)

	chat, err = saveChat(ctx, tx, chat)
	if err != nil && err != ErrChatAlreadyExists {
		return
	} else if err == ErrChatAlreadyExists {
		tx.Rollback(ctx)

		tx, err = pool.Begin(ctx)
		if err != nil {
			log.Err(err).Msg("Failed to create chat")
			return
		}

		chat, err = findChatByUserIds(ctx, tx, firstUserId, secondUserId)
		if err != nil {
			return
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to create chat")
		return
	}

	return chat, nil, nil
}
