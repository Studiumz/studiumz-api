package match

import (
	"context"

	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/oklog/ulid/v2"
)

func GetMatchById(matchId string) (match Match, err error) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	return findMatchById(ctx, tx, ulid.MustParse(matchId))
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
	return saveNewMatch(ctx, tx, m)
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
