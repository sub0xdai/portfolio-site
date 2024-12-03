package knowledge

import (
	"fmt"
	"log"
	"sync"

	"github.com/sub0x/resume-ai/internal/types"
)

// Service coordinates vault reading and embedding generation
type Service struct {
	vault      VaultManager
	embeddings Manager
	mu         sync.RWMutex
	vaultPath  string
}

// NewService creates a new knowledge service
func NewService(vaultPath, openAIKey string) (*Service, error) {
	vault, err := NewVault(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault: %v", err)
	}

	if err := vault.LoadNotes(); err != nil {
		return nil, fmt.Errorf("failed to load notes: %v", err)
	}

	embeddings := NewEmbeddingManager(openAIKey)

	service := &Service{
		vault:      vault,
		embeddings: embeddings,
		mu:         sync.RWMutex{},
		vaultPath:  vaultPath,
	}

	return service, nil
}

// RefreshKnowledge updates the knowledge base with current vault contents
func (s *Service) RefreshKnowledge() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Refreshing knowledge base...")
	if err := s.vault.LoadNotes(); err != nil {
		return fmt.Errorf("failed to reload notes: %v", err)
	}

	if err := s.embeddings.UpdateEmbeddings(s.vault.GetAllNotes()); err != nil {
		return fmt.Errorf("failed to update embeddings: %v", err)
	}

	log.Println("Knowledge base refreshed successfully")
	return nil
}

// LoadAndProcess initializes and processes the knowledge base
func (s *Service) LoadAndProcess() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Starting LoadAndProcess...")
	if err := s.vault.LoadNotes(); err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	notes := s.vault.GetAllNotes()
	log.Printf("Processed %d notes total", len(notes))
	
	if err := s.RefreshKnowledge(); err != nil {
		return fmt.Errorf("failed to process knowledge base: %v", err)
	}

	return nil
}

// LoadNotes loads notes from the vault without generating embeddings
func (s *Service) LoadNotes() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Loading notes...")
	if err := s.vault.LoadNotes(); err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	notes := s.vault.GetAllNotes()
	log.Printf("Loaded %d notes", len(notes))
	return nil
}

// Query finds relevant content for a query
func (s *Service) Query(query string, limit int) ([]types.Chunk, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chunks, err := s.embeddings.FindSimilar(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar content: %v", err)
	}

	return chunks, nil
}

// GetAllTags returns all unique tags from the vault
func (s *Service) GetAllTags() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vault.GetAllTags()
}

// GetNotesByTag returns all notes with a specific tag
func (s *Service) GetNotesByTag(tag string) []types.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vault.GetNotesByTag(tag)
}

// GetAllNotes returns all notes from the vault
func (s *Service) GetAllNotes() []types.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.vault.GetAllNotes()
}
