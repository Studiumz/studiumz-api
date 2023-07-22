package auth

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/jellydator/validation/is"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidUserId       = validation.NewError("auth:invalid_user_id", "Invalid user id")
	ErrUserFullNameTooLong = validation.NewError("auth:user_full_name_too_long", "Full name can't be longer than 255 characters")
	ErrInvalidUserEmail    = validation.NewError("auth:invalid_user_email", "Invalid user email")
	ErrInvalidUserStatus   = validation.NewError("auth:invalid_user_status", "Invalid user status")
)

func validateUserId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidUserId
	}

	return id, nil
}

func validateUserFullName(fullName string) (err error) {
	fullName = strings.TrimSpace(fullName)
	return validation.Validate(
		&fullName,
		validation.When(
			!validation.IsEmpty(fullName),
			validation.Length(1, 255).ErrorObject(ErrUserFullNameTooLong),
		),
	)
}

func validateUserEmail(email string) (err error) {
	email = strings.TrimSpace(email)
	return validation.Validate(
		&email,
		validation.When(
			!validation.IsEmpty(email),
			is.Email.ErrorObject(ErrInvalidUserEmail),
		),
	)
}

func validateUserStatus(status int) (err error) {
	if status < int(ONBOARDING) || status > int(ACTIVE) {
		return ErrInvalidUserStatus
	}

	return nil
}
