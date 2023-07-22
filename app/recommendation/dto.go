package recommendation

import "github.com/Studiumz/studiumz-api/app/auth"

type ScoredUser struct {
	// ranging from 0 to 1
	Similarity float64 `json:"similarity"`

	// extend User
	auth.User
}
