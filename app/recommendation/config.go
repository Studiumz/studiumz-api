package recommendation

import (
	"errors"

	"github.com/cohere-ai/cohere-go"
	"github.com/rs/zerolog/log"
)

var (
	co              *cohere.Client
	errInjectCohere = errors.New("Cohere could not be injected to recommendation module")
)

func InjectCohereClientAdapter(client *cohere.Client) {
	if client == nil {
		log.Fatal().Err(errInjectCohere).Send()
	}
	co = client
}
