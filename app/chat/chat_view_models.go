package chat

import (
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type ChatViewModel struct {
	Chat
	LastMsgTime    null.Time `json:"last_msg_time"`
	RecipientEmail string    `json:"recipient_email"`
	RecipientName  string    `json:"recipient_name"`
}

type ChatMetadata struct {
	Id             ulid.ULID `json:"id"`
	RecipientName  string    `json:"recipient_name"`
	RecipientEmail string    `json:"recipient_email"`
}

type ChatDetail struct {
	ChatMetadata
	Messages []*Message `json:"messages"`
}
