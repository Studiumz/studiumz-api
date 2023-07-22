package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type UserStatus uint

const (
	ONBOARDING UserStatus = iota
	ACTIVE
)

type User struct {
	Id        ulid.ULID   `json:"id"`
	FullName  null.String `json:"full_name"`
	Nickname  null.String `json:"nickname"`
	Email     null.String `json:"email"`
	Avatar    null.String `json:"avatar"`
	Gender    null.Int    `json:"gender"`
	Struggles null.String `json:"struggles"`
	BirthDate null.Time   `json:"birth_date"`
	Status    UserStatus  `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}

func NewUserWithOAuth(name, email null.String) (User, map[string]error) {
	id := ulid.Make()

	return User{
		Id:        id,
		FullName:  name,
		Email:     email,
		Status:    ONBOARDING,
		CreatedAt: time.Now(),
	}, nil
}
