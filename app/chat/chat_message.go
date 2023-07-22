package chat

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Message struct {
	Id         ulid.ULID   `json:"id"`
	ChatId     ulid.ULID   `json:"chat_id"`
	FromUserId ulid.ULID   `json:"from_user_id"`
	Text       string      `json:"text"`
	FileUrl    null.String `json:"file_url"`
	CreatedAt  time.Time   `json:"created_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}

func NewMessage(chatIdStr, fromUserIdStr, text string) (Message, map[string]error) {
	errs := make(map[string]error)

	chatId, err := validateChatId(chatIdStr)
	if err != nil {
		errs["chat_id"] = err
	}
	fromUserId, err := validateMessageFromUserId(fromUserIdStr)
	if err != nil {
		errs["from_user_id"] = err
	}
	if err := validateMessageText(text); err != nil {
		errs["text"] = err
	}

	return Message{
		Id:         ulid.Make(),
		ChatId:     chatId,
		FromUserId: fromUserId,
		Text:       text,
		CreatedAt:  time.Now(),
	}, nil
}
