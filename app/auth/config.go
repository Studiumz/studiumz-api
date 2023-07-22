package auth

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var (
	pool         *pgxpool.Pool
	firebaseApp  *firebase.App
	firebaseAuth *auth.Client

	ErrNilPool                              = errors.New("Connection pool can't be nil")
	ErrFirebaseAdminServiceAccountFileEmpty = errors.New("Firebase admin service account file can't be empty")
)

func SetPool(newPool *pgxpool.Pool) {
	if newPool == nil {
		log.Fatal().Err(ErrNilPool).Msg("Failed to set connection pool for auth module")
	}

	pool = newPool
}

func ConfigureFirebaseAdminSdk(serviceAccountFile string) {
	serviceAccountFile = strings.TrimSpace(serviceAccountFile)
	if len(serviceAccountFile) == 0 {
		log.Fatal().Err(ErrFirebaseAdminServiceAccountFileEmpty).Msg("Failed to configure firebase admin sdk")
	}

	ctx := context.Background()

	opt := option.WithCredentialsFile(serviceAccountFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure firebase admin sdk")
	}

	firebaseApp = app

	auth, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure firebase admin sdk")
	}

	firebaseAuth = auth
}
