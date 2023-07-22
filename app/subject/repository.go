package subject

import (
	"context"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findAllSubjects(ctx context.Context, tx pgx.Tx) (subjects []Subject, err error) {
	q := "SELECT * FROM subjects WHERE deleted_at IS NULL;"

	if err = pgxscan.Select(ctx, tx, &subjects, q); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return []Subject{}, nil
		}

		log.Err(err).Msg("Failed to find subjects")
		return subjects, err
	}

	return subjects, nil
}

func findSubjectsOfUser(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (subjects []Subject, err error) {
	q := `
	SELECT s.id, s.name, s.created_at, s.deleted_at 
	FROM subjects s 
	JOIN users_subjects us ON s.id = us.subject_id 
	WHERE us.user_id = $1 AND s.deleted_at IS NULL;
	`

	if err = pgxscan.Select(ctx, tx, &subjects, q, userId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return []Subject{}, nil
		}

		log.Err(err).Msg("Failed to find subjects")
		return subjects, err
	}

	return subjects, nil
}

func findSubjectByName(ctx context.Context, tx pgx.Tx, name string) (subject Subject, err error) {
	var subjects []Subject
	q := "SELECT * FROM subjects WHERE name = $1 AND deleted_at IS NULL;"

	if err = pgxscan.Select(ctx, tx, &subjects, q, name); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Subject{}, nil
		}

		log.Err(err).Msg("Failed to find subjects")
		return Subject{}, err
	}

	return subjects[0], nil
}

func CreateUserSubject(ctx context.Context, tx pgx.Tx, subjectNames []string, user auth.User) (err error) {
	q := "INSERT INTO users_subjects (id, user_id, subject_id) VALUES ($1, $2, $3);"
	// var newSubject Subject
	for _, subjectName := range subjectNames {
		subject, err := findSubjectByName(ctx, tx, subjectName)
		if err != nil {
			log.Err(err).Msg("subject does not exist")
			return err
		}

		// err = pgxscan.Get(ctx, tx, &newSubject, q, ulid.Make(), user.Id, subject.Id)
		comm, err := tx.Exec(ctx, q, ulid.Make(), user.Id, subject.Id)
		if err != nil {
			if err.Error() == "scanning one: no rows in result set" {
				return err
			}

			log.Err(err).Msg("could not save subject")
			return err
		}
		if comm.RowsAffected() == 0 {
			return errNoRowsAffected
		}
	}

	log.Info().Msg("all subjects saved to user with no problems")
	return nil
}

func CreateSubjects(ctx context.Context, tx pgx.Tx, subjectNames []string) (err error) {
	q := "INSERT INTO subjects(id, name) VALUES($1, $2)"
	for _, subjectName := range subjectNames {
		comm, err := tx.Exec(ctx, q, ulid.Make(), subjectName)
		if err != nil {
			if err.Error() == "scanning one: no rows in result set" {
				return err
			}

			log.Err(err).Msg("could not create subject")
			return err
		}
		if comm.RowsAffected() == 0 {
			return errNoRowsAffected
		}
	}

	log.Info().Msg("all subjects created with no problems")
	return nil
}
