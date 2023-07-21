package models

import "github.com/oklog/ulid/v2"

type Account struct {
	Id     ulid.ULID `json:"id"`
	UserId ulid.ULID `json:"user_id"`
}
