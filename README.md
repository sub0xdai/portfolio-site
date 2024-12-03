# Interactive AI-Powered Resume

A modern, interactive resume that uses AI to engage with visitors, powered by Go, HTMX, and OpenAI. The system integrates with an Obsidian knowledge base to provide detailed insights about experience, projects, and skills.

## Description

This project implements an interactive resume website that allows visitors to have natural conversations about my professional experience. The AI chatbot has access to my Obsidian vault and resume data, enabling it to provide detailed, context-aware responses about my skills, projects, and experience.

## Features

- Interactive AI chat interface with natural language understanding
- Real-time responses using HTMX for smooth interactions
- Integration with Obsidian vault for deep knowledge access
- Modern, responsive UI with TailwindCSS
- Vector-based semantic search for accurate information retrieval
- Resume data visualization and dynamic content updates

## System Architecture

### Backend
- **Go Server**: Fast and efficient HTTP handling using Fiber framework
- **Vector Database**: Stores embeddings of Obsidian notes and resume data
- **OpenAI Integration**: Powers the natural language understanding
- **Knowledge Base**: Processes and indexes Obsidian vault content

### Frontend
- **HTMX**: Enables dynamic content updates without complex JavaScript
- **TailwindCSS**: Provides modern, responsive styling
- **Hyperscript**: Handles simple client-side interactions
- **Static Assets**: Serves optimized images and resources

### Data Flow
1. User sends a message through the chat interface
2. HTMX sends request to Go backend
3. Backend processes query using vector search
4. Relevant context is retrieved from knowledge base
5. OpenAI generates response using retrieved context
6. Response is streamed back to frontend
7. HTMX updates UI smoothly with new content

## Prerequisites

- Go 1.21+
- OpenAI API key
- Obsidian vault with markdown files
- PostgreSQL with pgvector extension (for vector storage)

## Installation

1. Clone the repository:
   ```
   git clone [your-repo-url]
   cd [repo-name]
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Create a `.env` file in the project root:
   ```
   OPENAI_API_KEY=your_api_key_here
   DATABASE_URL=your_postgres_connection_string
   PORT=3000
   ```

4. Build and run:
   ```
   go run main.go
   ```

## Usage

1. Access the website at `http://localhost:3000`
2. Interact with the AI chat interface to learn about my experience
3. Explore the dynamic resume sections
4. Ask detailed questions about projects and skills

## Development

### Project Structure
```
.
├── main.go              # Entry point
├── static/             # Static assets
│   └── index.html      # Main HTML template
├── internal/           # Internal packages
│   ├── api/           # API handlers
│   ├── db/            # Database operations
│   ├── knowledge/     # Knowledge base processing
│   └── models/        # Data models
└── scripts/           # Utility scripts
```

## Contributing

Feel free to submit issues and enhancement requests!


