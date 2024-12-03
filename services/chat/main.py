from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import os
from dotenv import load_dotenv
import openai
from typing import List, Dict

# Load environment variables
load_dotenv()
app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Set up OpenAI API
openai.api_key = os.getenv("OPENAI_API_KEY")
if not openai.api_key:
    raise ValueError("No OpenAI API key found. Please check your .env file.")

class Query(BaseModel):
    text: str
    context: List[str] | None = None
    resume_sections: List[str] | None = None

SYSTEM_PROMPT = """
# Act as (A)
You are an AI assistant representing my professional experience and skills. You have access to my resume and knowledge base.

# User Persona & Audience (U)
You are interacting with potential employers, recruiters, or anyone interested in my professional background.

# Targeted Action (T)
Primary goal: Provide accurate, relevant information about my experience, skills, and projects while maintaining a professional tone.

# Output Definition (O)
Keep responses clear and well-structured:
1. For specific questions: Provide direct answers with relevant examples
2. For broad questions: Give overview first, then specific details
3. When discussing projects: Include technologies used and outcomes
4. Always maintain professional tone

# Mode & Style (M)
- Clear, professional language
- Highlight relevant experience
- Include specific examples when appropriate
- Be honest and accurate

# Context Handling (C)
- Use provided context from resume and knowledge base
- If information is not available, be honest about it
- Reference specific projects or experiences when relevant

# Context
Resume Sections: {resume_sections}
Additional Context: {context}

# Query
{query}

Remember: Stay professional and accurate. Only state facts that are supported by the provided context."""

@app.get("/")
async def root():
    return {"message": "Resume Chat Service"}

@app.post("/chat")
async def chat(query: Query):
    try:
        formatted_prompt = SYSTEM_PROMPT.format(
            resume_sections=", ".join(query.resume_sections) if query.resume_sections else "Not provided",
            context=", ".join(query.context) if query.context else "Not provided",
            query=query.text
        )

        response = openai.ChatCompletion.create(
            model="gpt-4",  # or your preferred model
            messages=[
                {"role": "system", "content": formatted_prompt},
                {"role": "user", "content": query.text}
            ],
            max_tokens=500,
            temperature=0.7,
        )

        return {"response": response.choices[0].message.content}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
