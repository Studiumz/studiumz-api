package chat

import "gopkg.in/guregu/null.v4"

type ChatViewModel struct {
	Chat
	LastMsgTime    null.Time `json:"last_msg_time"`
	RecipientEmail string    `json:"recipient_email"`
	RecipientName  string    `json:"recipient_name"`
}
