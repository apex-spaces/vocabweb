# VocabWeb Backend

Go backend for VocabWeb vocabulary learning platform.

## Tech Stack

- **Go**: 1.22+
- **Router**: chi (github.com/go-chi/chi/v5)
- **Database**: PostgreSQL with pgx driver
- **Auth**: Google Identity Platform (Firebase Admin SDK)
- **Config**: Environment variables

## Project Structure

```
backend/
├── main.go                 # Entry point
├── internal/
│   ├── config/            # Configuration loading
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware (auth, CORS)
│   ├── model/             # Data models
│   ├── repository/        # Database layer
│   └── router/            # Route registration
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL database (optional for initial setup)
- Firebase project with Identity Platform enabled

### Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Edit `.env` with your configuration

3. Install dependencies:
```bash
go mod download
```

4. Run the server:
```bash
go run main.go
```

Server will start on `http://localhost:8080`

## API Endpoints

### Public
- `GET /health` - Health check

### Protected (requires Firebase JWT)
- `GET /api/v1/auth/profile` - Get user profile
- `GET /api/v1/words` - List words
- `GET /api/v1/words/{id}` - Get word by ID

## Docker

Build image:
```bash
docker build -t vocabweb-backend .
```

Run container:
```bash
docker run -p 8080:8080 --env-file .env vocabweb-backend
```

## Deployment to Cloud Run

```bash
gcloud run deploy vocabweb-backend \
  --source . \
  --region asia-east2 \
  --platform managed \
  --allow-unauthenticated
```
