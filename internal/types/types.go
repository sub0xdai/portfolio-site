package types

// Query represents the chat query structure
type Query struct {
	Text           string   `json:"text"`
	Context        []string `json:"context,omitempty"`
	ResumeSections []string `json:"resume_sections,omitempty"`
}

// ChatResponse represents the response structure
type ChatResponse struct {
	Response string `json:"response"`
}

// ScoredChunk represents a scored piece of text
type ScoredChunk struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// Config represents the application configuration
type Config struct {
	ResumePath      string
	ObsidianPath    string
	CurrentRole     string
	ExperienceYears string
	KeySkills       []string
}
