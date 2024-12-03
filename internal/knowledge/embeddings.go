package knowledge

import (
	"github.com/sub0x/resume-ai/internal/types"
)

// Manager provides an interface for embedding operations
type Manager interface {
	UpdateEmbeddings(notes []types.Note) error
	FindSimilar(query string, limit int) ([]types.Chunk, error)
}

// embeddingManager implements the Manager interface
type embeddingManager struct {
	manager *types.EmbeddingManager
}

// NewEmbeddingManager creates a new embedding manager
func NewEmbeddingManager(openAIKey string) Manager {
	return &embeddingManager{
		manager: types.NewEmbeddingManager(openAIKey),
	}
}

// UpdateEmbeddings processes notes and updates their embeddings
func (e *embeddingManager) UpdateEmbeddings(notes []types.Note) error {
	return e.manager.UpdateEmbeddings(notes)
}

// FindSimilar finds chunks similar to the query
func (e *embeddingManager) FindSimilar(query string, limit int) ([]types.Chunk, error) {
	return e.manager.FindSimilar(query, limit)
}
