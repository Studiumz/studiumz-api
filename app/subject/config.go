package subject

import (
	"errors"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	pool *pgxpool.Pool

	errSubjectNotFound     = errors.New("match not found")
	errFailToUpdateSubject = errors.New("could not update existing subject")
	errNoRowsAffected      = errors.New("no subjects affected")
)

func SetPool(newPool *pgxpool.Pool) {
	if newPool == nil {
		log.Fatal().Err(auth.ErrNilPool).Msg("Failed to set connection pool for subject module")
	}

	pool = newPool
}
