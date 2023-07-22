package subject

import (
	"context"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/rs/zerolog/log"
)

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

func getSubjectsOfUser(user auth.User) (subjects []Subject, err error) {
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

func UserAddSubjects(subjectNames []string, user auth.User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	err = CreateUserSubject(ctx, tx, subjectNames, user)
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
