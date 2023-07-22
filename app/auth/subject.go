package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Subject struct {
	Id        ulid.ULID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt null.Time `json:"deleted_at"`
}
