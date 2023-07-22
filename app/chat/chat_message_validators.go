package chat

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidMessageId         = validation.NewError("chat:invalid_message_id", "Invalid message id")
	ErrInvalidMessageFromUserId = validation.NewError("chat:invalid_message_from_user_id", "Invalid message from user id")
	ErrMessageTextEmpty         = validation.NewError("chat:text_empty", "Text can't be empty")
	ErrMessageTextTooLong       = validation.NewError("chat:text_too_long", "Text can't be longer than 500 characters")
)

func validateMessageId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidMessageId
	}

	return id, nil
}

func validateMessageFromUserId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidMessageFromUserId
	}

	return id, nil
}

func validateMessageText(text string) (err error) {
	text = strings.TrimSpace(text)
	return validation.Validate(
		&text,
		validation.Required.ErrorObject(ErrMessageTextEmpty),
		validation.Length(1, 500).ErrorObject(ErrMessageTextTooLong),
	)

}
