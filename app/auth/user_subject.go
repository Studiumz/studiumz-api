package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type UserSubject struct {
	Id        ulid.ULID `json:"id"`
	UserId    ulid.ULID `json:"user_id"`
	SubjectId ulid.ULID `json:"subject_id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
