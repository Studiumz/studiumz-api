package chat

import (
	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidChatId     = validation.NewError("chat:invalid_chat_id", "Invalid chat id")
	ErrInvalidChatUserId = validation.NewError("chat:invalid_chat_user_id", "Invalid chat user id")
)

func validateChatId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidChatId
	}

	return id, nil
}

func validateChatUserId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidChatUserId
	}

	return id, nil
}
