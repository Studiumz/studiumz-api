package match

import (
	"github.com/Studiumz/studiumz-api/app/auth"
	"gopkg.in/guregu/null.v4"
)

type CreateMatchReq struct {
	// Id                ulid.ULID   `json:"id"`
	// MatcherId         ulid.ULID   `json:"matcher_id"`
	// Status            MatchStatus `json:"match_status"`
	// CreatedAt         time.Time   `json:"created_at"`
	// DeletedAt         time.Time   `json:"deleted_at"`
	// MatcheeId         ulid.ULID   `json:"matchee_id"`
	InvitationMessage null.String `json:"invitation_message"`
}

type UserMatch struct {
	auth.User
	Match
}
