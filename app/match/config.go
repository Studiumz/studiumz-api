package match

import (
	"errors"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	pool *pgxpool.Pool

	ErrMatchNotFound     = errors.New("match not found")
	ErrFailToUpdateMatch = errors.New("could not update existing match")
	ErrNoRowsAffected    = errors.New("no matches affected")
)

func SetPool(newPool *pgxpool.Pool) {
	if newPool == nil {
		log.Fatal().Err(auth.ErrNilPool).Msg("Failed to set connection pool for auth module")
	}

	pool = newPool
}
