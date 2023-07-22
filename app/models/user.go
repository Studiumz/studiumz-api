package models

import (
	"github.com/oklog/ulid/v2"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
)

type User struct {
	Id        ulid.ULID         `json:"id"`
	FullName  string            `json:"full_name"`
	Nickname  string            `json:"nickname"`
	Email     string            `json:"email"`
	ImageUrl  string            `json:"image_url"`
	Gender    Gender            `json:"gender"`
	Bio       string            `json:"bio"`
	BirthDate date.Date         `json:"birth_date"`
	Status    int               `json:"status"` // iota
	CreatedAt datetime.DateTime `json:"created_at"`
	UpdatedAt datetime.DateTime `json:"updated_at"`
	DeletedAt datetime.DateTime `json:"deleted_at"`
}
