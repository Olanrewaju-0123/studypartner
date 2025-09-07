# 📚 AI Student Study Partner

An **AI-powered study assistant** that helps students learn smarter by generating **summaries, flashcards, and quizzes** from uploaded notes.  
Built with **Go (backend)** + **Next.js + Tailwind (frontend)** + **Postgres (DB)** + **local/free AI models**.

---

## 🚀 Features

- Upload notes (PDF, DOCX, TXT)
- AI-generated **summaries**
- AI-generated **flashcards (Q&A)**
- AI-generated **quizzes (MCQs)**
- Semantic search over notes (using pgvector)
- User authentication (JWT)

---

## 🛠️ Tech Stack

### Frontend

- **Next.js** (React framework)
- **TailwindCSS** (styling)
- Fetch APIs → Connects to Go backend

### Backend

- **Go** (Gin/Fiber for REST API)
- **PostgreSQL** (main DB)
- **pgvector** (store embeddings for semantic search)
- **AI integration**:
  - [Ollama](https://ollama.ai/) → Run free local models (Llama 2/3, Mistral)
  - Or HuggingFace free inference API

### AI Models

- Summarization / Q&A: `llama2`, `mistral`, or any HuggingFace free model
- Embeddings: `all-MiniLM-L6-v2` (free on HuggingFace)

---

## ⚙️ Architecture

```mermaid
flowchart LR
  A[Frontend (Next.js + Tailwind)] -->|API Calls| B[Backend (Go)]
  B -->|Store Notes| C[(PostgreSQL + pgvector)]
  B -->|Send Text| D[AI Models via Ollama/HuggingFace]
  D -->|Return Summary/Flashcards/Quiz| B
  B -->|JSON Response| A
```

## 📂 Project Structure

### Frontend

```
/frontend
  /pages
    index.tsx        # Landing page
    upload.tsx       # Upload notes
    study.tsx        # Summaries, flashcards, quizzes
  /components
    UploadForm.tsx
    Flashcard.tsx
    Quiz.tsx
  /utils
    api.ts           # Fetch wrapper
```

### Backend

```
/backend
  /cmd
    main.go          # Entry point
  /routes
    notes.go         # Upload & process notes
    study.go         # Get summary, flashcards, quiz
  /services
    ai.go            # AI model calls
    embedding.go     # Embedding + pgvector
  /db
    db.go            # DB connection
    models.go        # User, Notes, Flashcards, Quiz
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+ with pgvector extension
- Ollama (for local AI models) or HuggingFace API key

### Backend Setup

```bash
cd backend
go mod init studypartner
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/pgvector/pgvector-go
go run cmd/main.go
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

### Database Setup

```sql
-- Create database
CREATE DATABASE studypartner;

-- Enable pgvector extension
CREATE EXTENSION vector;

-- Create tables (see backend/db/migrations.sql)
```

## 📡 API Endpoints

### Upload Notes

```
POST /api/notes/upload
{
  "userId": "123",
  "file": "base64/pdf"
}
```

### Get Summary

```
GET /api/notes/:id/summary
```

### Get Flashcards

```
GET /api/notes/:id/flashcards
```

### Get Quiz

```
GET /api/notes/:id/quiz
```

## 🗺️ Roadmap

### MVP

- [x] Project setup
- [ ] Upload notes
- [ ] Generate summaries
- [ ] Flashcards
- [ ] Quizzes
- [ ] JWT auth

### Advanced

- [ ] Semantic search
- [ ] Study progress tracking
- [ ] Collaborative features
- [ ] Mobile app
