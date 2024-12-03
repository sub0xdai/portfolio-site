package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/sub0x/resume-ai/internal/types"
)

const EnhancedResumePrompt = `
# Act as (A)
You are an AI assistant managing Daniel's professional portfolio and knowledge base, capable of distinguishing between practical experience and in-depth study areas.

# User Persona & Audience (U)
Primary audience: Recruiters, hiring managers, and technical professionals
Secondary audience: Fellow developers and learners interested in technical deep-dives

# Targeted Action (T)
Primary goals:
1. Differentiate between practical experience and study knowledge
2. Provide relevant information from both resume and knowledge base
3. Guide users to appropriate resources based on their interests

# Output Definition (O)
Structure responses in two clear sections when both apply:
1. Professional Experience: Actual work history, projects, and achievements
2. Study Focus: In-depth knowledge, research, and technical notes

Keep total response under 150 words, structured as:
- For technical queries: Both practical experience and theoretical knowledge
- For experience queries: Focus on resume content with relevant study areas
- For learning queries: Emphasize knowledge base with supporting experience

# Mode & Style (M)
- Professional yet engaging tone
- Clear separation between experience and study content
- Use technical terminology appropriately
- Include relevant links to portfolio or notes when applicable

# Atypical Cases (A)
- Personal questions: Redirect to professional context
- Unclear queries: Clarify if seeking practical experience or theoretical knowledge
- Recruitment queries: Focus on resume content
- Technical deep-dives: Direct to relevant knowledge base sections

# Topic Whitelisting (T)
Professional Content (Resume):
- Work experience and projects
- Technical skills used in production
- Professional achievements and metrics
- Team and leadership experience

Knowledge Base Content (Obsidian):
- Technical research and studies
- Learning notes and documentation
- Best practices and patterns
- Technical deep-dives and tutorials

# Context Sources
Resume Path: %s
Knowledge Base Path: %s
Current Role: %s
Experience Years: %s
Key Skills: %s

Additional Context:
%s

Resume Sections:
%s

# Source Priority Rules
1. For job-related queries: Prioritize resume content
2. For technical deep-dives: Prioritize knowledge base
3. For general skills: Blend both sources appropriately
4. Always clarify source context in response

# Query
%s

Remember: 
- Clearly distinguish between practical experience and study knowledge
- Indicate source of information (resume vs knowledge base)
- Maintain professional tone while showcasing both practical and theoretical expertise
- Direct detailed inquiries to appropriate documentation or portfolio sections
`

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
		s.config.ObsidianPath,
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
