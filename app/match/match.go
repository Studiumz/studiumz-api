package match

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Match struct {
	Id                ulid.ULID   `json:"id"`
	MatcherId         ulid.ULID   `json:"matcher_id"`
	MatcheeId         ulid.ULID   `json:"matchee_id"`
	Status            MatchStatus `json:"match_status"`
	InvitationMessage null.String `json:"invitation_message"`
	CreatedAt         time.Time   `json:"created_at"`
	DeletedAt         time.Time   `json:"deleted_at"`
}
