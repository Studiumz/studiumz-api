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
	env                  string
	pool                 *pgxpool.Pool
	firebaseApp          *firebase.App
	firebaseAuth         *auth.Client
	jwtIssuer            string
	jwtAccessTokenSecret []byte

	ErrNilPool                              = errors.New("Connection pool can't be nil")
	ErrFirebaseAdminServiceAccountFileEmpty = errors.New("Firebase admin service account file can't be empty")
	ErrJwtIssuerEmpty                       = errors.New("JWT issuer can't be empty")
	ErrJWTAccessTokenSecretEmpty            = errors.New("JWT access token secret can't be empty")
	errSubjectNotFound                      = errors.New("match not found")
	errFailToUpdateSubject                  = errors.New("could not update existing subject")
	errNoRowsAffected                       = errors.New("no subjects affected")
)

func SetEnv(environment string) {
	env = environment
}

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

func ConfigureJWTProperties(issuer, accessTokenSecret string) {
	issuer = strings.TrimSpace(issuer)
	if len(issuer) == 0 {
		log.Fatal().Err(ErrJwtIssuerEmpty).Msg("Failed to configure JWT properties")
	}

	accessTokenSecret = strings.TrimSpace(accessTokenSecret)
	if len(accessTokenSecret) == 0 {
		log.Fatal().Err(ErrJWTAccessTokenSecretEmpty).Msg("Failed to configure JWT properties")
	}

	jwtIssuer = issuer
	jwtAccessTokenSecret = []byte(accessTokenSecret)
}
