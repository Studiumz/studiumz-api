package auth

import (
	"context"
	"errors"
)

var ErrInvalidFirebaseIdToken = errors.New("Invalid firebase id token")

func validateUserFirebaseIdToken(ctx context.Context, idToken string) {

}
