package app

import (
	"github.com/cohere-ai/cohere-go"
	"github.com/rs/zerolog/log"
)

func ConfigureCohere(c Config) *cohere.Client {
	client, err := cohere.CreateClient(c.CohereApiKey)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not configure Cohere Client")
	}

	c.CohereClient = client
	return client
}
