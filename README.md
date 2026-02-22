# CareerManifest — AI-Powered Career Decision Platform for Indian Students

A full-stack career assessment platform that helps Indian students choose the right career path using a weighted scoring engine, risk analysis, 5-year salary projections, and optional AI-powered explanations.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.22 · Gin v1.9.1 · MySQL |
| **Auth** | JWT (HS256) · bcrypt |
| **Frontend** | Next.js 14 (App Router) · TypeScript · TailwindCSS |
| **AI (Optional)** | Groq (Llama3) / Claude API |

## Architecture

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                  # Environment configuration
│   ├── database/                # MySQL connection & migration
│   ├── models/                  # Database models
│   ├── dto/                     # Request/Response DTOs
│   ├── repository/              # Data access layer
│   ├── service/                 # Business logic
│   ├── engine/                  # Scoring engine
│   ├── handler/                 # HTTP handlers
│   ├── middleware/              # JWT, CORS, Admin middleware
│   ├── router/                  # Route registration
│   └── seed/                    # Question seed data
├── schema.sql
├── .env.example
└── go.mod

frontend/
├── src/
│   ├── app/                     # Next.js App Router pages
│   │   ├── page.tsx             #   Landing page
│   │   ├── login/               #   Login page
│   │   ├── register/            #   Register page
│   │   ├── dashboard/           #   User dashboard
│   │   ├── assessment/          #   Take assessment
│   │   ├── result/[id]/         #   View result
│   │   └── admin/               #   Admin panel
│   ├── components/ui/           # ShadCN-style components
│   └── lib/                     # Utils, API client, auth context, types
├── package.json
├── tailwind.config.ts
└── next.config.js
```

## Prerequisites

- **Go** 1.22+
- **Node.js** 18+ (npm/yarn)
- **MySQL** 8.0+

## Quick Start

### 1. Database Setup

```sql
CREATE DATABASE IF NOT EXISTS career_manifest;
CREATE USER IF NOT EXISTS 'career_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON career_manifest.* TO 'career_user'@'localhost';
FLUSH PRIVILEGES;
```

Or run the full schema:

```bash
mysql -u root -p < backend/schema.sql
```

### 2. Backend

```bash
cd backend

# Copy env and configure
cp .env.example .env
# Edit .env with your MySQL credentials and JWT secret

# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

The server starts on `http://localhost:8080` by default.

**On first run:**
- Tables are auto-created/migrated
- 30 assessment questions are seeded
- An admin user is created (see `ADMIN_EMAIL` / `ADMIN_PASSWORD` in `.env`)

### 3. Frontend

```bash
cd frontend

# Install dependencies
npm install

# Copy env
cp .env.local.example .env.local

# Start dev server
npm run dev
```

The frontend runs on `http://localhost:3000` and proxies API requests to the backend.

## Environment Variables

### Backend (`.env`)

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | Server port |
| `DB_HOST` | Yes | `localhost` | MySQL host |
| `DB_PORT` | No | `3306` | MySQL port |
| `DB_USER` | Yes | `root` | MySQL user |
| `DB_PASSWORD` | Yes | — | MySQL password |
| `DB_NAME` | Yes | `career_manifest` | Database name |
| `JWT_SECRET` | Yes | — | JWT signing key (use a strong random string) |
| `JWT_EXPIRY_HOURS` | No | `72` | Token expiry in hours |
| `CORS_ORIGIN` | No | `http://localhost:3000` | Allowed CORS origin |
| `LLM_PROVIDER` | No | — | `groq` or `claude` (leave empty to disable) |
| `LLM_API_KEY` | No | — | API key for the LLM provider |
| `ADMIN_EMAIL` | No | `admin@careermanifest.com` | Default admin email |
| `ADMIN_PASSWORD` | No | `admin123` | Default admin password |

### Frontend (`.env.local`)

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | No | `http://localhost:8080` | Backend API base URL |

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login |

### Protected (requires JWT)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/auth/profile` | Get current user profile |
| GET | `/api/questions` | Get active questions |
| POST | `/api/assessment/submit` | Submit assessment answers |
| GET | `/api/assessment/history` | Get user's assessment history |
| GET | `/api/assessment/:id` | Get specific assessment result |

### Admin (requires JWT + admin role)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/admin/stats` | Dashboard statistics |
| GET | `/api/admin/questions` | List all questions |
| POST | `/api/admin/questions` | Create a question |
| PUT | `/api/admin/questions/:id` | Update a question |

## Scoring Engine

The platform evaluates users across **6 career categories**:

1. IT / Software Jobs
2. MBA (India)
3. Government Exams
4. Startup / Entrepreneurship
5. Higher Studies (India)
6. MS Abroad

Each question option carries weighted scores for each category. The engine:
- Accumulates weighted scores from 30 questions
- Ranks careers by total score
- Computes a **Risk Score** using the formula:

```
RiskScore = (IncomeUrgency × 0.35) + (FamilyDependency × 0.25) +
            (RiskTolerance × 0.20) + (CareerInstabilityIndex × 0.20)
```

- Generates salary projections, preparation roadmaps, required skills, suggested exams, and colleges for the best career path.

## Mobile App Support

The REST API is fully compatible with a React Native (Expo) mobile app. The same endpoints work — just point your mobile API client to the backend URL.

## Optional LLM Integration

To enable AI-powered explanations in results:

1. Set `LLM_PROVIDER=groq` or `LLM_PROVIDER=claude` in `.env`
2. Set `LLM_API_KEY` to your API key
3. Restart the backend

If no LLM is configured, a template-based explanation is generated instead.

## Production Deployment

### Backend
```bash
# Build binary
cd backend
CGO_ENABLED=0 GOOS=linux go build -o careermanifest-server cmd/server/main.go

# Run with production env
./careermanifest-server
```

### Frontend
```bash
cd frontend
npm run build
npm start
```

### Docker (Example)

```dockerfile
# Backend Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /usr/local/bin/
CMD ["server"]
```

```dockerfile
# Frontend Dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./
RUN npm ci --production
CMD ["npm", "start"]
```

## License

MIT
