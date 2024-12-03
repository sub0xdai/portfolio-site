package types

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Note represents a single document in the knowledge vault
type Note struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	FilePath string   `json:"file_path"`
}

// Chunk represents a piece of text with its embedding and metadata
type Chunk struct {
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	Score     float32   `json:"score"`
	Embedding []float32 `json:"embedding"`
	FilePath  string    `json:"file_path"`
}

// EmbeddingManager handles the generation and management of embeddings
type EmbeddingManager struct {
	Client     *openai.Client
	Embeddings map[string][]Chunk
	OpenAIKey  string
}

// NewEmbeddingManager creates a new EmbeddingManager instance
func NewEmbeddingManager(openAIKey string) *EmbeddingManager {
	return &EmbeddingManager{
		Client:     openai.NewClient(openAIKey),
		Embeddings: make(map[string][]Chunk),
		OpenAIKey:  openAIKey,
	}
}

// UpdateEmbeddings processes notes and updates their embeddings
func (e *EmbeddingManager) UpdateEmbeddings(notes []Note) error {
	newEmbeddings := make(map[string][]Chunk)
	
	for _, note := range notes {
		chunks, err := e.processNote(note)
		if err != nil {
			return fmt.Errorf("failed to process note %s: %v", note.FilePath, err)
		}
		newEmbeddings[note.FilePath] = chunks
	}
	
	e.Embeddings = newEmbeddings
	return nil
}

// processNote splits a note into chunks and generates embeddings
func (e *EmbeddingManager) processNote(note Note) ([]Chunk, error) {
	// Split content into chunks (simplified for now)
	chunks := strings.Split(note.Content, "\n\n")
	var result []Chunk
	
	for i, content := range chunks {
		if len(strings.TrimSpace(content)) == 0 {
			continue
		}
		
		embedding, err := e.generateEmbedding(content)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for chunk %d: %v", i, err)
		}
		
		result = append(result, Chunk{
			Content:   content,
			Embedding: embedding,
			FilePath:  note.FilePath,
		})
	}
	
	return result, nil
}

// generateEmbedding creates an embedding for a piece of text
func (e *EmbeddingManager) generateEmbedding(text string) ([]float32, error) {
	resp, err := e.Client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.AdaEmbeddingV2,
		},
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %v", err)
	}
	
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data received")
	}
	
	return resp.Data[0].Embedding, nil
}

// FindSimilar finds chunks similar to the query
func (e *EmbeddingManager) FindSimilar(query string, limit int) ([]Chunk, error) {
	queryEmbedding, err := e.generateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %v", err)
	}

	var allChunks []Chunk
	for _, chunks := range e.Embeddings {
		for _, chunk := range chunks {
			score := cosineSimilarity(queryEmbedding, chunk.Embedding)
			chunk.Score = float32(score)
			allChunks = append(allChunks, chunk)
		}
	}

	// Sort chunks by similarity score (highest first)
	sort.Slice(allChunks, func(i, j int) bool {
		return allChunks[i].Score > allChunks[j].Score
	})

	// Return top k results
	if len(allChunks) > limit {
		allChunks = allChunks[:limit]
	}
	return allChunks, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Vault manages the collection of notes
type Vault struct {
	Notes    []Note  `json:"notes"`
	RootPath string  `json:"root_path"`
}

// NewVault creates a new Vault instance
func NewVault(rootPath string) (*Vault, error) {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault root path does not exist: %s", rootPath)
	}
	return &Vault{
		RootPath: rootPath,
		Notes:    make([]Note, 0),
	}, nil
}

// LoadNotes reads all notes from the vault directory
func (v *Vault) LoadNotes() error {
	v.Notes = make([]Note, 0)
	
	err := filepath.Walk(v.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}
		
		var note Note
		if err := json.Unmarshal(data, &note); err != nil {
			return fmt.Errorf("failed to parse note %s: %v", path, err)
		}
		
		note.FilePath = path
		v.Notes = append(v.Notes, note)
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to load notes: %v", err)
	}
	
	return nil
}

// GetNoteByPath retrieves a note by its file path
func (v *Vault) GetNoteByPath(path string) (*Note, error) {
	for i := range v.Notes {
		if v.Notes[i].FilePath == path {
			return &v.Notes[i], nil
		}
	}
	return nil, fmt.Errorf("note not found: %s", path)
}

// GetNotesByTag returns all notes with a specific tag
func (v *Vault) GetNotesByTag(tag string) []Note {
	var result []Note
	for _, note := range v.Notes {
		for _, t := range note.Tags {
			if t == tag {
				result = append(result, note)
				break
			}
		}
	}
	return result
}

// GetAllTags returns all unique tags from the vault
func (v *Vault) GetAllTags() []string {
	tagMap := make(map[string]bool)
	for _, note := range v.Notes {
		for _, tag := range note.Tags {
			tagMap[tag] = true
		}
	}
	
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}
