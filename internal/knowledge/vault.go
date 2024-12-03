package knowledge

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sub0x/resume-ai/internal/types"
)

// Note represents a single document in the knowledge vault
type Note struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	FilePath string   `json:"file_path"`
}

// VaultManager provides an interface for vault operations
type VaultManager interface {
	LoadNotes() error
	GetNoteByPath(path string) (*types.Note, error)
	GetNotesByTag(tag string) []types.Note
	GetAllTags() []string
	GetAllNotes() []types.Note
}

// vaultManager implements the VaultManager interface
type vaultManager struct {
	rootPath string
	notes    []types.Note
}

// NewVault creates a new vault instance
func NewVault(rootPath string) (VaultManager, error) {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault root path does not exist: %s", rootPath)
	}

	return &vaultManager{
		rootPath: rootPath,
		notes:    make([]types.Note, 0),
	}, nil
}

// LoadNotes reads all notes from the vault directory
func (v *vaultManager) LoadNotes() error {
	v.notes = make([]types.Note, 0)
	
	err := filepath.Walk(v.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}

		note := types.Note{
			Title:    strings.TrimSuffix(filepath.Base(path), ".md"),
			Content:  string(content),
			FilePath: path,
			Tags:     []string{}, // TODO: Parse tags from content
		}

		v.notes = append(v.notes, note)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk vault directory: %v", err)
	}

	return nil
}

// GetNoteByPath retrieves a note by its file path
func (v *vaultManager) GetNoteByPath(path string) (*types.Note, error) {
	for _, note := range v.notes {
		if note.FilePath == path {
			return &note, nil
		}
	}
	return nil, fmt.Errorf("note not found: %s", path)
}

// GetNotesByTag returns all notes with a specific tag
func (v *vaultManager) GetNotesByTag(tag string) []types.Note {
	var result []types.Note
	for _, note := range v.notes {
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
func (v *vaultManager) GetAllTags() []string {
	tagMap := make(map[string]bool)
	for _, note := range v.notes {
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

// GetAllNotes returns all notes from the vault
func (v *vaultManager) GetAllNotes() []types.Note {
	return v.notes
}
