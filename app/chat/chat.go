package chat

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Chat struct {
	Id           ulid.ULID `json:"id"`
	FirstUserId  ulid.ULID `json:"first_user_id"`
	SecondUserId ulid.ULID `json:"second_user_id"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    null.Time `json:"deleted_at"`
}

func NewChat(firstUserIdStr, secondUserIdStr string) (Chat, map[string]error) {
	errs := make(map[string]error)

	firstUserId, err := validateChatUserId(firstUserIdStr)
	if err != nil {
		errs["first_user_id"] = err
	}
	secondUserId, err := validateChatUserId(secondUserIdStr)
	if err != nil {
		errs["second_user_id"] = err
	}
	if len(errs) != 0 {
		return Chat{}, errs
	}

	id := ulid.Make()

	return Chat{Id: id, FirstUserId: firstUserId, SecondUserId: secondUserId, CreatedAt: time.Now()}, nil
}
