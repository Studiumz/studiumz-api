package match

import (
	"context"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func GetMatchById(matchId string) (match Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get match by id")
		return
	}

	defer tx.Rollback(ctx)

	match, err = findValidMatchById(ctx, tx, ulid.MustParse(matchId))
	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get match by id")
		return
	}
	return match, nil
}

func createNewMatch(matcheeId string, body CreateMatchReq, user auth.User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	m := Match{
		Id:                ulid.Make(),
		MatcherId:         user.Id,
		MatcheeId:         ulid.MustParse(matcheeId),
		Status:            MatchStatus(PENDING),
		InvitationMessage: body.InvitationMessage,
	}
	err = saveNewMatchOrSkip(ctx, tx, m)
	if err != nil {
		log.Err(err).Msg("Failed to create match")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to create match")
		return
	}
	return nil
}

func createNewSkip(matcheeId string, user auth.User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	m := Match{
		Id:                ulid.Make(),
		MatcherId:         user.Id,
		MatcheeId:         ulid.MustParse(matcheeId),
		Status:            MatchStatus(SKIPPED),
		InvitationMessage: null.String{},
	}

	if err = saveNewMatchOrSkip(ctx, tx, m); err != nil {
		log.Err(err).Msg("Failed to create new skipped match")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to create new skipped match")
		return
	}
	return nil
}

func acceptMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	if err = updateMatchStatus(ctx, tx, MatchStatus(ACCEPTED), m.Id); err != nil {
		log.Err(err).Msg("Failed to update match status to accepted")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to update match status to accepted")
		return
	}
	return nil
}

func rejectMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	if err = updateMatchStatus(ctx, tx, MatchStatus(REJECTED), m.Id); err != nil {
		log.Err(err).Msg("Failed to update match status to rejected")
		return
	}
	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to update match status to rejected")
		return
	}
	return nil
}

func withdrawMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	if err = softDeleteMatch(ctx, tx, m.Id); err != nil {
		log.Err(err).Msg("Failed to soft-delete match (withdraw)")
		return
	}
	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to soft-delete match (withdraw)")
		return
	}
	return nil
}

func getUserIncoming(userId ulid.ULID) (matches []UserMatch, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	matches, err = findUserIncomingMatches(ctx, tx, userId)
	if err != nil {
		log.Err(err).Msg("Failed to get user incoming matches")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get user incoming matches")
		return
	}
	return matches, nil
}

func getUserOutgoing(userId ulid.ULID) (matches []Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	matches, err = findUserOutgoingMatches(ctx, tx, userId)
	if err != nil {
		log.Err(err).Msg("Failed to update get user outgoing matches")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to update get user outgoing matches")
		return
	}
	return matches, nil
}
