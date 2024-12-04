package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sub0x/resume-ai/internal/api"
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
	log.Printf("PORT: %s", os.Getenv("PORT"))

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize API handlers with simplified knowledge service
	apiHandler := api.NewHandler(nil)  // We'll modify the api package to work without knowledge service

	// Create Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New())

	// Serve static files
	app.Static("/static", "./static")

	// API routes
	app.Post("/api/query", apiHandler.HandleQuery)
	app.Get("/api/tags", apiHandler.HandleGetTags)
	app.Get("/api/notes/tag/:tag", apiHandler.HandleGetNotesByTag)

	// Web routes
	app.Get("/", func(c *fiber.Ctx) error {
		component := templates.Home(nil)
		return component.Render(c.Context(), c.Response().BodyWriter())
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Starting server on %s...", addr)
	log.Fatal(app.Listen(addr))
}
