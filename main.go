package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

type ChatRequest struct {
	Text           string   `json:"text"`
	ResumeSection  []string `json:"resume_sections"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Resume sections as static content
var resumeSections = []string{
	`PERSONAL INFORMATION
	Daniel Palazzolo
	www.youtube.com/@m0xu44
	danielpalazzolo4@gmail.com
	www.linkedin.com/in/dpalazzolo
	GitHub: sub0xdai`,

	`SUMMARY
	Enthusiastic and adaptable professional transitioning from a successful sales career to software engineering. 
	With a growth mindset, everyday is an opportunity to be improve on the day before. 
	Holistic approach with a focus on systems and process.`,

	`TECHNICAL SKILLS
	Programming Languages: Python, JavaScript, Go, PHP, C#
	Database: SQL, MongoDB
	Tools: Git, GitHub, Docker, VSCode, Vim
	Operating Systems: Linux, Windows`,

	`EDUCATION
	Diploma of Advanced Programming - Holmesglen Tafe (2024)
	- Coursework: Data Structures, Algorithms, Object-Oriented Programming, Data Modelling, System Design
	- Projects: CSIRO Research Management App, Library Management System, Dynamic Full-Stack App
	
	Linux Administration Course - ProLUG
	Web Developer Bootcamp - Generation Australia (2022)
	Graduate Certificate Commerce - RMIT (2020-2021)`,

	`PROJECTS
	Cryptocurrency Trading TUI
	- Developing Terminal User Interface for trading perpetual swaps in Go
	- Using Bubbletea library for intuitive interfaces
	
	Library Management System
	- Built web-based system using PHP, Bootstrap, Express.js
	- Full-stack project with database management and RESTful API design`,

	`EXPERIENCE
	Tech Development Analyst at Accenture (2023-2024)
	- Collaborated on cloud-based solutions
	- Monitored system performance for major energy retailer
	- Developed tools to streamline operational tasks`,
}

const indexHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Daniel Palazzolo | Portfolio</title>
    <link rel="icon" type="image/x-icon" href="data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAA/4QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEREQAAAAAAEAAAEAAAAAEAAAABAAAAEAAAAAAQAAAQAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAEAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD8HwAA++8AAPf3AAD3+wAA9/sAAPf7AAD3+wAA9/sAAP//AAD//wAA//8AAP//AAD//wAA" />
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary-bg: #0a0d16;
            --nav-bg: #12151f;
            --chat-bg: #161a26;
            --text-color: #e4e8f7;
            --text-muted: #8b92a8;
            --accent-color: #7c6aff;
            --accent-hover: #6a57ff;
            --accent-glow: rgba(124, 106, 255, 0.15);
            --border-color: #2a2e3d;
            --shadow-color: rgba(0, 0, 0, 0.3);
            --input-bg: #1e2332;
            --message-user-bg: #2a2440;
            --message-ai-bg: #1e2332;
            --highlight-color: #ff6b6b;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            transition: all 0.3s ease;
        }

        body {
            font-family: 'JetBrains Mono', monospace;
            background-color: var(--primary-bg);
            background-image: url('/static/images/evening-sky.png');
            background-size: cover;
            background-position: center;
            background-attachment: fixed;
            color: var(--text-color);
            line-height: 1.6;
            min-height: 100vh;
        }

        nav {
            background-color: rgba(18, 21, 31, 0.85);
            backdrop-filter: blur(10px);
            -webkit-backdrop-filter: blur(10px);
            padding: 1.2rem 2rem;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            z-index: 1000;
            box-shadow: 0 2px 20px var(--shadow-color);
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }

        .nav-content {
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .logo {
            font-size: 1.4rem;
            font-weight: 600;
            color: var(--text-color);
            text-decoration: none;
            letter-spacing: -0.5px;
        }

        .nav-links {
            display: flex;
            gap: 2.5rem;
            align-items: center;
        }

        .nav-links a {
            color: var(--text-muted);
            text-decoration: none;
            font-size: 0.9rem;
            font-weight: 500;
            letter-spacing: 0.5px;
            transition: all 0.2s ease;
            position: relative;
        }

        .nav-links a:hover {
            color: var(--text-color);
        }

        .nav-links a::after {
            content: '';
            position: absolute;
            bottom: -4px;
            left: 0;
            width: 0;
            height: 2px;
            background: var(--accent-color);
            transition: width 0.2s ease;
        }

        .nav-links a:hover::after {
            width: 100%;
        }

        .chat-button {
            background-color: var(--accent-color);
            color: white;
            padding: 0.6rem 1.2rem;
            border-radius: 8px;
            font-weight: 500;
            letter-spacing: 0.5px;
            box-shadow: 0 2px 10px var(--accent-glow);
            border: 1px solid rgba(255, 255, 255, 0.1);
            cursor: pointer;
        }

        .chat-button:hover {
            background-color: var(--accent-hover);
            transform: translateY(-1px);
            box-shadow: 0 4px 15px var(--accent-glow);
        }

        .chat-button.active {
            background-color: var(--highlight-color);
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0% {
                box-shadow: 0 0 0 0 rgba(255, 107, 107, 0.4);
            }
            70% {
                box-shadow: 0 0 0 10px rgba(255, 107, 107, 0);
            }
            100% {
                box-shadow: 0 0 0 0 rgba(255, 107, 107, 0);
            }
        }

        .main-content {
            padding-top: 5rem;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .welcome-content {
            text-align: center;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
        }

        .welcome-content h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            background: linear-gradient(45deg, var(--accent-color), var(--highlight-color));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            letter-spacing: -1px;
            line-height: 1.2;
        }

        .welcome-content p {
            font-size: 1.1rem;
            color: var(--text-muted);
            margin-bottom: 2rem;
            letter-spacing: 0.5px;
        }

        .welcome-content .highlight {
            color: var(--highlight-color);
            font-weight: 500;
        }

        /* Modal styles */
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background-color: rgba(0, 0, 0, 0.5);
            z-index: 2000;
            align-items: center;
            justify-content: center;
            backdrop-filter: blur(5px);
        }

        .modal.active {
            display: flex;
        }

        #chat-container {
            background-color: rgba(26, 29, 40, 0.95);
            backdrop-filter: blur(10px);
            -webkit-backdrop-filter: blur(10px);
            width: 100%;
            max-width: 800px;
            border-radius: 20px;
            box-shadow: 0 8px 30px var(--shadow-color);
            padding: 30px;
            position: relative;
            overflow: hidden;
            margin: 2rem;
            border: 1px solid rgba(255, 255, 255, 0.1);
            animation: modalFadeIn 0.3s ease;
        }

        @keyframes modalFadeIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .close-modal {
            position: absolute;
            top: 20px;
            right: 20px;
            color: var(--text-muted);
            cursor: pointer;
            font-size: 1.5rem;
            line-height: 1;
            padding: 5px;
            transition: color 0.2s ease;
        }

        .close-modal:hover {
            color: var(--text-color);
        }

        /* Welcome content styles */
        #chat-output {
            margin: 20px 0;
            padding: 15px;
            border-radius: 15px;
            background-color: var(--input-bg);
            min-height: 300px;
            max-height: 500px;
            overflow-y: auto;
            scrollbar-width: thin;
            scrollbar-color: var(--accent-color) var(--input-bg);
        }

        #chat-output::-webkit-scrollbar {
            width: 8px;
        }

        #chat-output::-webkit-scrollbar-track {
            background: var(--input-bg);
        }

        #chat-output::-webkit-scrollbar-thumb {
            background-color: var(--accent-color);
            border-radius: 4px;
        }

        .message {
            padding: 12px 16px;
            border-radius: 12px;
            margin: 8px 0;
            max-width: 85%;
            animation: fadeIn 0.3s ease;
        }

        .message.user {
            background-color: var(--message-user-bg);
            margin-left: auto;
            border-bottom-right-radius: 4px;
        }

        .message.ai {
            background-color: var(--message-ai-bg);
            margin-right: auto;
            border-bottom-left-radius: 4px;
        }

        #input-container {
            position: relative;
            margin-top: 20px;
        }

        #chat-input {
            width: 100%;
            padding: 16px;
            border: 2px solid var(--border-color);
            border-radius: 12px;
            background-color: var(--input-bg);
            font-size: 16px;
            color: var(--text-color);
            resize: none;
            transition: border-color 0.3s ease;
        }

        #chat-input::placeholder {
            color: var(--text-muted);
        }

        #chat-input:focus {
            outline: none;
            border-color: var(--accent-color);
        }

        button {
            background-color: var(--accent-color);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 12px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 500;
            margin-top: 12px;
            width: 100%;
            transition: transform 0.2s ease, background-color 0.2s ease;
        }

        button:hover {
            background-color: var(--accent-hover);
            transform: translateY(-2px);
        }

        button:active {
            transform: translateY(0);
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @media (max-width: 768px) {
            nav {
                padding: 1rem;
            }

            .nav-links {
                gap: 1rem;
            }

            #chat-container {
                margin: 1rem;
                padding: 20px;
                border-radius: 15px;
            }

            .message {
                max-width: 90%;
            }
        }
    </style>
</head>
<body>
    <nav>
        <div class="nav-content">
            <a href="/" class="logo">Daniel Palazzolo</a>
            <div class="nav-links">
                <a href="#projects">Projects</a>
                <a href="#contact">Contact</a>
                <div class="chat-button" onclick="toggleChat()">AI Chat</div>
            </div>
        </div>
    </nav>

    <div class="main-content">
        <div class="welcome-content">
            <h1>Building Digital<br>Experiences</h1>
            <p>Software Engineer with a passion for <span class="highlight">problem-solving</span> and <span class="highlight">innovation</span></p>
        </div>
    </div>

    <div class="modal" id="chat-modal">
        <div id="chat-container">
            <span class="close-modal" onclick="closeChat()">&times;</span>
            <h1>AI Resume Assistant</h1>
            <div id="chat-output"></div>
            <div id="input-container">
                <textarea id="chat-input" placeholder="Ask me anything about Daniel's experience, skills, or background..." rows="3"></textarea>
                <button onclick="sendMessage()">Send Message</button>
            </div>
        </div>
    </div>

    <script>
        let chatActive = false;

        function toggleChat() {
            const chatButton = document.querySelector('.chat-button');
            chatActive = !chatActive;
            
            if (chatActive) {
                chatButton.classList.add('active');
                chatButton.textContent = 'Chat Active';
                document.getElementById('chat-modal').classList.add('active');
                document.body.style.overflow = 'hidden';
            } else {
                chatButton.classList.remove('active');
                chatButton.textContent = 'AI Chat';
                document.getElementById('chat-modal').classList.remove('active');
                document.body.style.overflow = 'auto';
            }
        }

        function closeChat() {
            toggleChat();
        }

        // Close modal when clicking outside the chat container
        document.getElementById('chat-modal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeChat();
            }
        });

        function sendMessage() {
            const input = document.getElementById("chat-input");
            const output = document.getElementById("chat-output");
            const text = input.value.trim();
            
            if (!text) return;
            
            // Create and append user message
            const userMessage = document.createElement("div");
            userMessage.className = "message user";
            userMessage.textContent = text;
            output.appendChild(userMessage);
            
            input.value = "";
            output.scrollTop = output.scrollHeight;
            
            fetch("/api/chat", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ text: text })
            })
            .then(response => response.json())
            .then(data => {
                // Create and append AI message
                const aiMessage = document.createElement("div");
                aiMessage.className = "message ai";
                aiMessage.textContent = data.response;
                output.appendChild(aiMessage);
                output.scrollTop = output.scrollHeight;
            })
            .catch(error => {
                const errorMessage = document.createElement("div");
                errorMessage.className = "message ai";
                errorMessage.style.color = "#dc3545";
                errorMessage.textContent = "Sorry, I encountered an error. Please try again.";
                output.appendChild(errorMessage);
            });
        }
        
        // Allow Enter key to send message, Shift+Enter for new line
        document.getElementById("chat-input").addEventListener("keydown", function(e) {
            if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        });

        // Add initial greeting when chat is opened
        document.getElementById('chat-modal').addEventListener('transitionend', function() {
            if (this.classList.contains('active') && !this.hasGreeted) {
                const output = document.getElementById("chat-output");
                const greeting = document.createElement("div");
                greeting.className = "message ai";
                greeting.textContent = "Hello! I am here to help you learn more about Daniel's experience and skills. What would you like to know?";
                output.appendChild(greeting);
                this.hasGreeted = true;
            }
        });
    </script>
</body>
</html>`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable not set")
	}

	app := fiber.New(fiber.Config{
		AppName: "Personal AI Resume",
	})

	app.Use(cors.New())

	// Serve static files
	app.Static("/static", "./static")

	// Serve the HTML page at root
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(indexHTML)
	})

	api := app.Group("/api")
	api.Post("/chat", handleChat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

func handleChat(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Add resume sections as context
	req.ResumeSection = resumeSections

	// Forward request to Python service
	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to process request",
		})
	}

	resp, err := http.Post("http://localhost:8000/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to reach chat service",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.JSON(chatResp)
}
