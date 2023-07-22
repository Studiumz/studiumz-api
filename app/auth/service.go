package auth

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func createTestFirebaseIdToken(ctx context.Context) (token string, err error) {
	token, err = firebaseAuth.CustomTokenWithClaims(ctx, "secButWXDCewaqsP3YRqzGEjg2l1", map[string]interface{}{"name": "Nayyara Airlangga Raharjo"})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create test firebase id token")
		return token, err
	}

	return token, nil
}

func FilterUnmatched(targetUser User) (filteredUsers []User, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("Failed to get unmatched users against %s", targetUser.Id.String()))
		return
	}
	return filterUnmatchedUsers(ctx, tx, targetUser)
}

func signInWithGoogle(ctx context.Context, idToken string) (signIn UserSignIn, errs map[string]error, err error) {
	authToken, err := validateUserFirebaseIdToken(ctx, idToken)
	if err != nil {
		return
	}

	accountId, fullName, email := extractProfileFromFirebaseIdToken(authToken)

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to sign in with Google")
		return
	}

	defer tx.Rollback(ctx)

	var user User
	var account Account

	user, err = findUserByEmail(ctx, tx, email.String)
	if err != nil {
		if err != ErrUserDoesNotExist {
			return
		}

		user, errs = NewUserWithOAuth(fullName, email)
		if errs != nil {
			return
		}

		user, err = saveUser(ctx, tx, user)
		if err != nil {
			return
		}

		account, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
		if err != nil {
			if err != ErrAccountDoesNotExist {
				return
			}

			account, errs = NewAccount(
				user.Id.String(),
				string(OAUTH),
				null.NewString("", false),
				string(GOOGLE),
				accountId,
			)
			if errs != nil {
				return
			}

			account, err = saveAccount(ctx, tx, account)
			if err != nil {
				return
			}
		}
	} else {
		account, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
		if err != nil {
			if err != ErrAccountDoesNotExist {
				return
			}

			account, errs = NewAccount(
				user.Id.String(),
				string(OAUTH),
				null.NewString("", false),
				string(GOOGLE),
				accountId,
			)
			if errs != nil {
				return
			}

			account, err = saveAccount(ctx, tx, account)
			if err != nil {
				return
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to sign in with Google")
		return
	}

	accessToken, err := generateUserAccessToken(user.Id.String())
	if err != nil {
		return
	}

	return UserSignIn{
		UserId:      user.Id,
		AccessToken: accessToken,
	}, nil, nil
}

func finishOnboarding(user User, body finishOnboardingReq) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	err = saveCompleteProfile(ctx, tx, user, body)
	if err != nil {
		log.Err(err).Msg("Failed to update profile of user")
		return
	}

	err = CreateUserSubject(ctx, tx, body.SubjectNames, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to add subjects to user")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to complete onboarding for user")
		return
	}
	return nil
}

func getAllSubjects() (subjects []Subject, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	subjects, err = findAllSubjects(ctx, tx)
	if err != nil {
		log.Err(err).Msg("Failed to fetch subjects")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to fetch subjects")
		return
	}
	return subjects, nil
}

func getSubjectsOfUser(user User) (subjects []Subject, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	subjects, err = findSubjectsOfUser(ctx, tx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to fetch subjects")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to fetch subjects")
		return
	}
	return subjects, nil
}

func UserAddSubjects(subjectNames []string, user User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	err = CreateUserSubject(ctx, tx, subjectNames, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to add subjects")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to add subjects")
		return
	}
	return nil
}

func AddSubjects(subjectNames []string) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	err = CreateSubjects(ctx, tx, subjectNames)
	if err != nil {
		log.Err(err).Msg("Failed to create subjects")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to create subjects")
		return
	}
	return nil
}
