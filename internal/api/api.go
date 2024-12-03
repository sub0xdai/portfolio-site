package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/sub0x/resume-ai/internal/types"
)

const EnhancedResumePrompt = `You are a concise AI assistant representing Daniel's professional portfolio and knowledge base.

Guidelines:
1. Keep responses under 3 sentences when possible
2. Use bullet points for lists
3. Focus on key achievements and metrics
4. Distinguish between practical experience and study knowledge
5. Use active voice

Context:
Resume Path: %s
Knowledge Base Path: %s
Current Role: %s
Experience Years: %s
Key Skills: %s

Additional Context:
%s

Resume Sections:
%s

Query: %s`

type Server struct {
	config *types.Config
	client *openai.Client
}

func NewServer(config *types.Config, client *openai.Client) *Server {
	return &Server{
		config: config,
		client: client,
	}
}

func (s *Server) SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./index.html")

	// Root endpoint
	r.GET("/api", s.handleRoot)

	// Chat endpoint
	r.POST("/chat", s.handleChat)

	return r
}

func (s *Server) handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the abyss",
	})
}

func (s *Server) handleChat(c *gin.Context) {
	var query types.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format context and resume sections
	contextStr := "No additional context provided"
	if len(query.Context) > 0 {
		contextStr = strings.Join(query.Context, "\n")
	}

	resumeSectionsStr := "No specific resume sections provided"
	if len(query.ResumeSections) > 0 {
		resumeSectionsStr = strings.Join(query.ResumeSections, "\n")
	}

	// Format the prompt
	formattedPrompt := fmt.Sprintf(
		EnhancedResumePrompt,
		s.config.ResumePath,
		s.config.KnowledgeBasePath,
		s.config.CurrentRole,
		s.config.ExperienceYears,
		strings.Join(s.config.KeySkills, ", "),
		contextStr,
		resumeSectionsStr,
		query.Text,
	)

	// Create OpenAI chat completion request
	resp, err := s.client.CreateChatCompletion(
		c.Request.Context(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: formattedPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query.Text,
				},
			},
			MaxTokens:   300,
			Temperature: 0.7,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resp.Choices) > 0 {
		c.JSON(http.StatusOK, types.ChatResponse{
			Response: resp.Choices[0].Message.Content,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No response generated",
		})
	}
}

// HandleChat processes a chat request and returns a response
func (s *Server) HandleChat(text string, resumeSections []string) (string, error) {
	prompt := fmt.Sprintf(EnhancedResumePrompt,
		s.config.ResumePath,
		s.config.KnowledgeBasePath,
		s.config.CurrentRole,
		s.config.ExperienceYears,
		strings.Join(s.config.KeySkills, ", "),
		"", // Additional context
		strings.Join(resumeSections, "\n"),
		text,
	)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
			MaxTokens:   250,
			Temperature: 0.5,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
