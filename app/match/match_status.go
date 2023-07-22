package match

type MatchStatus string

const (
	PENDING  MatchStatus = "PENDING"
	REJECTED MatchStatus = "REJECTED"
	ACCEPTED MatchStatus = "ACCEPTED"
	SKIPPED  MatchStatus = "SKIPPED"
)
