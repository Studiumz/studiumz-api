package match

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findValidMatchById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (match Match, err error) {
	q := "SELECT * FROM matches WHERE id = $1 AND deleted_at IS NULL AND status != 'SKIPPED';"

	if err = pgxscan.Get(ctx, tx, &match, q, id); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Match{}, nil
		}

		log.Err(err).Msg("Failed to find match by id")
		return match, err
	}

	return match, nil
}

func saveNewMatchOrSkip(ctx context.Context, tx pgx.Tx, m Match) (err error) {
	q := "INSERT INTO matches(id, matcher_id, matchee_id, status, invitation_message) VALUES($1, $2, $3, $4, $5);"

	n, err := tx.Exec(ctx, q, m.Id, m.MatcherId, m.MatcheeId, m.Status, m.InvitationMessage)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return err
		}

		log.Err(err).Msg("could not save new match")
		return err
	}
	if n.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func updateMatchStatus(ctx context.Context, tx pgx.Tx, newStatus MatchStatus, id ulid.ULID) (err error) {
	q := "UPDATE matches SET status = $1 where id = $2 AND deleted_at IS NULL AND status != 'SKIPPED';"

	n, err := tx.Exec(ctx, q, newStatus, id)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return err
		}

		log.Err(err).Msg("Failed to update match")
		return err
	}
	if n.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func softDeleteMatch(ctx context.Context, tx pgx.Tx, id ulid.ULID) (err error) {
	q := "UPDATE matches SET deleted_at = now() where id = $1 AND status != 'SKIPPED';"

	n, err := tx.Exec(ctx, q, id)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return err
		}

		log.Err(err).Msg("Failed to soft delete match")
		return err
	}
	if n.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func findUserIncomingMatches(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (matches []UserMatch, err error) {
	q := `SELECT * 
	FROM matches m JOIN user u ON m.matcher_id = u.id
	WHERE m.matchee_id = $1 
		AND m.deleted_at IS NULL 
		AND m.status != 'SKIPPED';`

	if err = pgxscan.Select(ctx, tx, &matches, q, userId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return []UserMatch{}, nil
		}

		log.Err(err).Msg("Failed to find incoming matches")
		return matches, err
	}

	return matches, nil
}

func findUserOutgoingMatches(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (matches []Match, err error) {
	q := "SELECT * FROM matches WHERE matchee_id = $1 AND deleted_at IS NULL AND status != 'SKIPPED';"

	if err = pgxscan.Select(ctx, tx, &matches, q, userId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return []Match{}, nil
		}

		log.Err(err).Msg("Failed to find outgoing matches")
		return matches, err
	}

	return matches, nil
}
