package auth

import (
	"context"
	"errors"

	"firebase.google.com/go/auth"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

var ErrInvalidFirebaseIdToken = errors.New("Invalid firebase id token")

func validateUserFirebaseIdToken(ctx context.Context, idToken string) (authToken *auth.Token, err error) {
	authToken, err = firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to validate user firebase id token")
		return authToken, ErrInvalidFirebaseIdToken
	}

	return authToken, nil
}

func extractProfileFromFirebaseIdToken(authToken *auth.Token) (uid string, name, email null.String) {
	claims := authToken.Claims

	if nameStr, ok := claims["name"].(string); !ok {
		name = null.NewString("", false)
	} else {
		name = null.StringFrom(nameStr)
	}

	if emailStr, ok := claims["email"].(string); !ok {
		email = null.NewString("", false)
	} else {
		email = null.StringFrom(emailStr)
	}

	return authToken.UID, name, email
}
