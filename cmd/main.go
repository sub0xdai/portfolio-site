package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sub0x/resume-ai/internal/api"
	"github.com/sub0x/resume-ai/internal/knowledge"
	"github.com/sub0x/resume-ai/internal/templates"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Print current working directory and environment variables for debugging
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)
	log.Printf("VAULT_PATH: %s", os.Getenv("VAULT_PATH"))
	log.Printf("PORT: %s", os.Getenv("PORT"))

	// Initialize knowledge service
	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		log.Fatal("VAULT_PATH environment variable is required")
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	service, err := knowledge.NewService(vaultPath, openAIKey)
	if err != nil {
		log.Fatalf("Failed to create knowledge service: %v", err)
	}

	// Load notes but don't wait for embeddings
	if err := service.LoadNotes(); err != nil {
		log.Fatalf("Failed to load notes: %v", err)
	}

	// Initialize API handlers
	apiHandler := api.NewHandler(service)

	// Create Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New())

	// API routes
	app.Post("/api/query", apiHandler.HandleQuery)
	app.Get("/api/tags", apiHandler.HandleGetTags)
	app.Get("/api/notes/tag/:tag", apiHandler.HandleGetNotesByTag)

	// Web routes
	app.Get("/", func(c *fiber.Ctx) error {
		notes := service.GetAllNotes()
		component := templates.Home(notes)
		return component.Render(c.Context(), c.Response().BodyWriter())
	})

	// Start refreshing knowledge base in the background
	go func() {
		log.Println("Starting knowledge base refresh in background...")
		if err := service.RefreshKnowledge(); err != nil {
			log.Printf("Error refreshing knowledge base: %v", err)
		}
		log.Println("Knowledge base refresh complete")
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Starting server on %s...", addr)
	log.Fatal(app.Listen(addr))
}
