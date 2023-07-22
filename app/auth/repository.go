package auth

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrAccountDoesNotExist = errors.New("Account does not exist")

func findUserById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (user User, err error) {
	q := "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &user, q, id); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return user, ErrUserDoesNotExist
		}

		log.Err(err).Msg("Failed to find user by id")
		return user, err
	}

	return user, nil
}

func findUserByEmail(ctx context.Context, tx pgx.Tx, email string) (user User, err error) {
	q := "SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &user, q, email); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return user, ErrUserDoesNotExist
		}

		log.Err(err).Msg("Failed to find user by email")
		return user, err
	}

	return user, nil
}

func findAccountByProviderAndProviderAccountId(ctx context.Context, tx pgx.Tx, provider AccountProvider, providerAccountId string) (account Account, err error) {
	q := "SELECT * FROM accounts WHERE provider = $1 AND provider_account_id = $2 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &account, q, provider, providerAccountId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return account, ErrAccountDoesNotExist
		}

		log.Err(err).Msg("Failed to find account by provider and provider account id")
		return account, err
	}

	return account, nil
}

func saveUser(ctx context.Context, tx pgx.Tx, user User) (newUser User, err error) {
	q := `
  INSERT INTO users (id, full_name, email, status, created_at) VALUES
  ($1, $2, $3, $4, $5)
  ON CONFLICT ON CONSTRAINT users_email_unique
  DO NOTHING
  RETURNING *
  `

	if err = pgxscan.Get(
		ctx,
		tx,
		&newUser,
		q,
		user.Id,
		user.FullName,
		user.Email,
		user.Status,
		user.CreatedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			newUser, err = findUserByEmail(ctx, tx, user.Email.String)
			if err != nil {
				return
			}

			return newUser, nil
		}

		log.Err(err).Msg("Failed to save user")
		return
	}

	return newUser, nil
}

func saveAccount(ctx context.Context, tx pgx.Tx, account Account) (newAccount Account, err error) {
	if _, err = findUserById(ctx, tx, account.UserId); err != nil {
		return
	}

	q := `
  INSERT INTO accounts (id, user_id, type, password_hash, provider, provider_account_id, created_at) VALUES
  ($1, $2, $3, $4, $5, $6, $7)
  ON CONFLICT ON CONSTRAINT accounts_provider_provider_account_id_unique
  DO NOTHING
  RETURNING *
  `

	if err = pgxscan.Get(
		ctx,
		tx,
		&newAccount,
		q,
		account.Id,
		account.UserId,
		account.Type,
		account.PasswordHash,
		account.Provider,
		account.ProviderAccountId,
		account.CreatedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			newAccount, err = findAccountByProviderAndProviderAccountId(ctx, tx, account.Provider, account.ProviderAccountId)
			if err != nil {
				return
			}

			return newAccount, nil
		}

		log.Err(err).Msg("Failed to save account")
		return
	}

	return newAccount, nil
}

func filterUnmatchedUsers(ctx context.Context, tx pgx.Tx, user User) (filteredUsers []User, err error) {
	q :=
		`SELECT * FROM users where id NOT IN (
			SELECT matcher_id FROM matches WHERE deleted_at IS NULL AND matchee_id = $1
			UNION
			SELECT matchee_id FROM matches WHERE deleted_at IS NULL AND matcher_id = $1
		);
	`
	err = pgxscan.Get(ctx, tx, &filteredUsers, q, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to filter unmatched users")
		return
	}
	return
}

func saveCompleteProfile(ctx context.Context, tx pgx.Tx, user User, body finishOnboardingReq) (err error) {
	q := `
		UPDATE users
		SET full_name = $1, nickname = $2, gender = $3, struggles = $4, birth_date = $5
		WHERE id = $6;
	`
	comm, err := tx.Exec(ctx, q, body.FullName, body.Nickname, body.Gender, body.Struggles, body.BirthDate, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to update user")
		return
	}
	if comm.RowsAffected() == 0 {
		return errors.New("no users affected")
	}
	return
}

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

func CreateUserSubject(ctx context.Context, tx pgx.Tx, subjectNames []string, userId ulid.ULID) (err error) {
	q := "INSERT INTO users_subjects (id, user_id, subject_id) VALUES ($1, $2, $3);"
	// var newSubject Subject
	for _, subjectName := range subjectNames {
		subject, err := findSubjectByName(ctx, tx, subjectName)
		if err != nil {
			log.Err(err).Msg("subject does not exist")
			return err
		}

		// err = pgxscan.Get(ctx, tx, &newSubject, q, ulid.Make(), user.Id, subject.Id)
		comm, err := tx.Exec(ctx, q, ulid.Make(), userId, subject.Id)
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
