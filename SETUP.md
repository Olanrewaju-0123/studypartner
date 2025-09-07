# 🚀 AI Study Partner - Setup Guide

This guide will help you set up the AI Study Partner application on your local machine.

## 📋 Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **Node.js 18+** - [Download here](https://nodejs.org/)
- **PostgreSQL 14+** with pgvector extension
- **Docker** (optional, for easy database setup) - [Download here](https://www.docker.com/)

## 🗄️ Database Setup

### Option 1: Using Docker (Recommended)

1. Start the PostgreSQL database with pgvector:

```bash
docker-compose up postgres -d
```

2. Wait for the database to be ready (check with `docker-compose ps`)

### Option 2: Manual PostgreSQL Setup

1. Install PostgreSQL 14+ on your system
2. Install the pgvector extension:

```bash
# On Ubuntu/Debian
sudo apt install postgresql-14-pgvector

# On macOS with Homebrew
brew install pgvector

# On Windows, download from: https://github.com/pgvector/pgvector/releases
```

3. Create the database:

```sql
CREATE DATABASE studypartner;
\c studypartner;
CREATE EXTENSION vector;
```

## 🤖 AI Model Setup

### Option 1: Using Ollama (Local AI Models)

1. Install Ollama:

```bash
# On macOS/Linux
curl -fsSL https://ollama.ai/install.sh | sh

# On Windows, download from: https://ollama.ai/download
```

2. Start Ollama:

```bash
ollama serve
```

3. Pull a model (in a new terminal):

```bash
ollama pull llama2
# or
ollama pull mistral
```

### Option 2: Using HuggingFace API

1. Get a free API key from [HuggingFace](https://huggingface.co/settings/tokens)
2. Add it to your environment variables

## ⚙️ Backend Setup

1. Navigate to the backend directory:

```bash
cd backend
```

2. Copy the environment file:

```bash
cp env.example .env
```

3. Edit `.env` with your configuration:

```env
DATABASE_URL=postgres://studypartner:studypartner123@localhost/studypartner?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here
OLLAMA_URL=http://localhost:11434
HUGGINGFACE_API_KEY=your-huggingface-api-key-here
PORT=8080
```

4. Install dependencies:

```bash
go mod tidy
```

5. Run the backend:

```bash
go run cmd/main.go
```

The backend will be available at `http://localhost:8080`

## 🎨 Frontend Setup

1. Navigate to the frontend directory:

```bash
cd frontend
```

2. Copy the environment file:

```bash
cp env.example .env.local
```

3. Edit `.env.local` with your configuration:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

4. Install dependencies:

```bash
npm install
```

5. Start the development server:

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 🧪 Testing the Setup

1. **Backend Health Check**: Visit `http://localhost:8080/health`
2. **Frontend**: Visit `http://localhost:3000`
3. **Upload a Test File**: Try uploading a PDF, DOCX, or TXT file
4. **Generate Study Materials**: Create summaries, flashcards, and quizzes

## 📁 Project Structure

```
studypartner/
├── backend/                 # Go backend
│   ├── cmd/                # Application entry point
│   ├── routes/             # API routes
│   ├── services/           # Business logic
│   ├── db/                 # Database models and migrations
│   ├── middleware/         # HTTP middleware
│   └── config/             # Configuration
├── frontend/               # Next.js frontend
│   ├── src/
│   │   ├── app/           # Next.js app directory
│   │   ├── components/    # React components
│   │   ├── utils/         # Utility functions
│   │   └── types/         # TypeScript types
│   └── public/            # Static assets
├── docker-compose.yml     # Database setup
└── README.md             # Project documentation
```

## 🔧 API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/me` - Get current user

### Notes

- `POST /api/notes/upload` - Upload a note
- `GET /api/notes` - Get user's notes
- `GET /api/notes/:id` - Get specific note
- `DELETE /api/notes/:id` - Delete note
- `POST /api/notes/search` - Search notes

### Study Materials

- `GET /api/study/notes/:id/summary` - Get summary
- `POST /api/study/notes/:id/summary` - Generate summary
- `GET /api/study/notes/:id/flashcards` - Get flashcards
- `POST /api/study/notes/:id/flashcards` - Generate flashcards
- `GET /api/study/notes/:id/quiz` - Get quiz
- `POST /api/study/notes/:id/quiz` - Generate quiz

## 🐛 Troubleshooting

### Database Connection Issues

- Ensure PostgreSQL is running
- Check the DATABASE_URL in your `.env` file
- Verify pgvector extension is installed

### AI Model Issues

- For Ollama: Ensure the service is running and models are pulled
- For HuggingFace: Check your API key and rate limits

### Frontend Issues

- Clear browser cache
- Check browser console for errors
- Ensure backend is running on the correct port

## 🚀 Deployment

For production deployment:

1. **Backend**: Build with `go build -o main cmd/main.go`
2. **Frontend**: Build with `npm run build`
3. **Database**: Use a managed PostgreSQL service with pgvector support
4. **AI Models**: Consider using cloud AI services for better performance

## 📝 Environment Variables

### Backend (.env)

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT tokens
- `OLLAMA_URL`: Ollama service URL (if using local models)
- `HUGGINGFACE_API_KEY`: HuggingFace API key (if using cloud models)
- `PORT`: Server port (default: 8080)

### Frontend (.env.local)

- `NEXT_PUBLIC_API_URL`: Backend API URL

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.
