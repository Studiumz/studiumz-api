package recommendation

import (
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/cohere-ai/cohere-go"
	"gonum.org/v1/gonum/mat"
)

const (
	MAX_RECOMMENDATION_BATCH_SIZE = 10
)

func getEmbeddings(targets []string) (*cohere.EmbedResponse, error) {
	return co.Embed(cohere.EmbedOptions{
		Model:    "embed-english-light-v2.0",
		Texts:    targets,
		Truncate: "END",
	})
}

func getCosineSimilarity(embeddingA []float64, embeddingB []float64) float64 {
	vecA := mat.NewVecDense(len(embeddingA), embeddingA)
	vecB := mat.NewVecDense(len(embeddingB), embeddingB)

	return mat.Dot(vecA, vecB) / (mat.Norm(vecA, 2) * mat.Norm(vecB, 2))
}

func CreateRecommendation(currentUser auth.User) (recommendedUsers []ScoredUser, err error) {
	recommendedUsers = []ScoredUser{}

	filteredUsers, err := auth.FilterUnmatched(currentUser)
	if err != nil {
		return
	}
	// get tags (subjects)
	tags := []string{"astronomy, thermodynamics, gravity"}                                                                       // TODO
	otherTags := []string{"human anatomy, genetics, carbon cycle", "anthropology, world history, sociology", "physics, science"} // TODO
	tags = append(tags, otherTags...)

	// get embeddings from tags, index 0 = currentUser's emb
	res, err := getEmbeddings(tags)
	if err != nil {
		return
	}

	// get similarity and append to recommendedUsers
	for i, embedding := range res.Embeddings[1 : len(res.Embeddings)-1] {
		// limit recommendations to max batch size
		if i == MAX_RECOMMENDATION_BATCH_SIZE || i == len(filteredUsers) {
			return
		}
		score := getCosineSimilarity(res.Embeddings[0], embedding)
		recommendedUsers = append(recommendedUsers, ScoredUser{score, filteredUsers[i]})
	}

	return
}
