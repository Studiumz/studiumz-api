package match

import (
	"context"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func GetMatchById(matchId string) (match Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return findValidMatchById(ctx, tx, ulid.MustParse(matchId))
}

func createNewMatch(matcheeId string, body CreateMatchReq, user auth.User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	m := Match{
		Id:                ulid.Make(),
		MatcherId:         user.Id,
		MatcheeId:         ulid.MustParse(matcheeId),
		Status:            MatchStatus(PENDING),
		InvitationMessage: body.InvitationMessage,
	}
	return saveNewMatchOrSkip(ctx, tx, m)
}

func createNewSkip(matcheeId string, user auth.User) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	m := Match{
		Id:                ulid.Make(),
		MatcherId:         user.Id,
		MatcheeId:         ulid.MustParse(matcheeId),
		Status:            MatchStatus(SKIPPED),
		InvitationMessage: null.String{},
	}
	return saveNewMatchOrSkip(ctx, tx, m)
}

func acceptMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return updateMatchStatus(ctx, tx, MatchStatus(ACCEPTED), m.Id)
}

func rejectMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return updateMatchStatus(ctx, tx, MatchStatus(REJECTED), m.Id)
}

func withdrawMatch(m Match) (err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return softDeleteMatch(ctx, tx, m.Id)
}

func getUserIncoming(userId ulid.ULID) (matches []Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return findUserIncomingMatches(ctx, tx, userId)
}

func getUserOutgoing(userId ulid.ULID) (matches []Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return findUserOutgoingMatches(ctx, tx, userId)
}
