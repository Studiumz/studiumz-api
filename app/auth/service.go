package auth

import "context"

func signInWithGoogle(ctx context.Context, idToken string) (signIn UserSignIn, errs map[string]error, err error) {
	validateUserFirebaseIdToken(ctx, idToken)

	return UserSignIn{}, nil, nil
}
