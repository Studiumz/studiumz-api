package auth

import (
	"context"

	"github.com/rs/zerolog/log"
)

func signInWithGoogle(ctx context.Context, idToken string) (signIn UserSignIn, errs map[string]error, err error) {
	authToken, err := validateUserFirebaseIdToken(ctx, idToken)
	if err != nil {
		return
	}

	accountId, name, email := extractProfileFromFirebaseIdToken(authToken)

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to sign in with Google")
		return
	}

	return UserSignIn{}, nil, nil
}
