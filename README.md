# 🚀 CareerManifest — ML-Powered Career Decision Platform

> **Helping Indian students choose the right career path** using a machine-learning-trained scoring engine, risk analysis, 5-year salary projections, and AI-powered explanations.

[![Engine](https://img.shields.io/badge/Engine-v3.0.0--ml-purple)](https://github.com/Chiragj2003/Carriermanifest)
[![ML Model](https://img.shields.io/badge/Model-Random%20Forest-blue)](ml/)
[![Accuracy](https://img.shields.io/badge/Accuracy-88.45%25-brightgreen)](ml/ml_weights.json)
[![Next.js](https://img.shields.io/badge/Next.js-15.5-black)](frontend/)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8)](backend/)
[![License](https://img.shields.io/badge/License-MIT-yellow)](#license)

**🌐 Live:** [carriermanifest.vercel.app](https://carriermanifest.vercel.app) &nbsp;|&nbsp; **API:** [careermanifest-api.onrender.com](https://careermanifest-api.onrender.com)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **30-Question Assessment** | Conversational-tone questions covering academics, finances, personality, and goals |
| **ML-Trained Scoring** | Weight matrix derived from a Random Forest model trained on 10K student profiles |
| **6 Career Paths** | IT/Software, MBA, Government Exams, Startup, Higher Studies (India), MS Abroad |
| **9-Dimensional Profiling** | Academic Strength, Financial Pressure, Risk Tolerance, Leadership, Tech Affinity, Govt Interest, Abroad Interest, Income Urgency, Career Instability |
| **Risk Analysis** | Weighted risk scoring with career-specific penalties |
| **5-Year Salary Projections** | Realistic salary curves for each career path |
| **AI Career Explanations** | Optional Groq Llama3-70B powered detailed explanations |
| **AI Chatbot** | Ask follow-up questions about your results |
| **Google OAuth** | One-click sign-in with Google |
| **Dark/Light Mode** | Toggle between themes |
| **Mobile Responsive** | Works on all screen sizes |
| **Admin Panel** | Manage questions, view stats |

---

## 🏗️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.22 · Gin v1.10 |
| **Database** | PostgreSQL (Neon) |
| **Auth** | JWT (HS256) · bcrypt · Google OAuth |
| **Frontend** | Next.js 15.5 (App Router) · React 19 · TypeScript · TailwindCSS |
| **ML Pipeline** | Python 3.13 · scikit-learn · Random Forest |
| **AI (Optional)** | Groq (Llama3-70B) |
| **Deployment** | Render (backend) · Vercel (frontend) |

---

## 📁 Project Structure

```
CareerManifest/
├── backend/
│   ├── cmd/server/main.go              # Entry point
│   └── internal/
│       ├── config/                     # Environment configuration
│       ├── database/                   # PostgreSQL connection & migration
│       ├── models/                     # Database models
│       ├── dto/                        # Request/Response DTOs
│       ├── repository/                 # Data access layer
│       ├── service/                    # Business logic
│       ├── engine/                     # 🧠 ML-trained scoring engine
│       │   ├── career.go              #   Career enum (6 paths)
│       │   ├── profile.go             #   9-dimensional user profile
│       │   ├── aggregator.go          #   Question → feature mapping
│       │   ├── matrix.go              #   ML-derived weight matrix
│       │   ├── scorer.go              #   Dot-product scoring
│       │   ├── risk.go                #   Risk calculation & penalties
│       │   ├── normalize.go           #   Score normalization
│       │   ├── explain.go             #   Result explanation generator
│       │   ├── enrichment.go          #   Salary, skills, colleges data
│       │   └── version.go             #   Engine + ML version tracking
│       ├── handler/                    # HTTP handlers
│       ├── middleware/                 # JWT, CORS, Admin middleware
│       ├── router/                     # Route registration
│       └── seed/                       # Question seed data
│
├── frontend/
│   └── src/
│       ├── app/                        # Next.js App Router pages
│       │   ├── page.tsx               #   Landing page
│       │   ├── login/                 #   Login (with Google OAuth)
│       │   ├── register/              #   Register (with Google OAuth)
│       │   ├── dashboard/             #   User dashboard
│       │   ├── assessment/            #   Take assessment
│       │   ├── result/[id]/           #   View result + AI chatbot
│       │   ├── admin/                 #   Admin panel
│       │   └── icon.svg               #   Favicon
│       ├── components/                 # Navbar, UI components
│       └── lib/                        # API client, auth context, types
│
├── ml/                                 # 🤖 ML Training Pipeline
│   ├── generate_dataset.py            #   Synthetic dataset generator (10K samples)
│   ├── train_model.py                 #   Train & compare 5 ML models
│   ├── career_training_data.csv       #   Generated training data
│   ├── ml_weights.json                #   Exported model weights & metadata
│   ├── matrix_ml.go.txt               #   Ready-to-paste Go weight matrix
│   ├── best_model_Random_Forest.joblib #  Serialized best model
│   └── scaler.joblib                  #   Feature scaler
│
└── render.yaml                         # Render deployment config
```

---

## 🤖 ML-Trained Scoring Engine (v3.0.0-ml)

The scoring engine uses a **weight matrix derived from machine learning** rather than hand-tuned values.

### Training Pipeline

1. **Synthetic Dataset Generation** — 10,000 Indian student profiles with 9 features across 6 career archetypes, including feature correlations, noise injection, and boundary cases
2. **Model Training** — 5 classifiers trained and compared:

| Model | Accuracy | F1 Score | CV F1 Mean |
|-------|----------|----------|------------|
| **🏆 Random Forest** | **88.45%** | **0.8844** | **0.8869** |
| Neural Network (MLP) | 88.20% | 0.8816 | 0.8853 |
| SVM (RBF Kernel) | 87.75% | 0.8779 | 0.8789 |
| Gradient Boosting | 87.00% | 0.8697 | 0.8814 |
| Logistic Regression | 85.75% | 0.8584 | 0.8706 |

3. **Weight Extraction** — Feature importances and class-conditional weights exported to JSON + Go format
4. **Domain Adjustment** — Blend of 60% tree importance + 40% logistic regression coefficients, with domain amplification

### Feature Importance Ranking

```
GovtInterest      ████████████████████ 0.188  (#1)
AbroadInterest    ██████████████████   0.172  (#2)
TechAffinity      ███████████████      0.142  (#3)
LeadershipScore   ██████████████       0.136  (#4)
CareerInstability █████████████        0.125  (#5)
RiskTolerance     ████████████         0.115  (#6)
AcademicStrength  ████████             0.076  (#7)
IncomeUrgency     ███                  0.026  (#8)
FinancialPressure ██                   0.018  (#9)
```

### Scoring Formula

The engine computes a **dot product** of the user's 9-dimensional profile vector against each career's weight vector:

$$\text{Score}_c = \sum_{i=0}^{8} \text{Profile}[i] \times W_c[i]$$

Scores are normalized, risk-adjusted, and ranked to produce the final career recommendation.

### Risk Calculation

$$\text{RiskScore} = 0.35 \times \text{IncomeUrgency} + 0.25 \times \text{FinancialPressure} + 0.20 \times \text{RiskTolerance} + 0.20 \times \text{CareerInstability}$$

Career-specific penalties apply (e.g., high financial pressure → Startup score reduced by 20%).

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.22+
- **Node.js** 18+
- **PostgreSQL** 14+ (or a [Neon](https://neon.tech) cloud database)
- **Python** 3.10+ (only for ML training)

### 1. Backend

```bash
cd backend

# Install dependencies
go mod tidy

# Set environment variables (see table below)
export DATABASE_URL="postgresql://user:pass@host/dbname?sslmode=require"
export JWT_SECRET="your-strong-secret-key"

# Run the server
go run cmd/server/main.go
```

The server starts on `http://localhost:8080`. On first run:
- Tables are auto-created/migrated
- 30 assessment questions are seeded
- An admin user is created

### 2. Frontend

```bash
cd frontend

# Install dependencies
npm install

# Set environment variables
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
echo "NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id" >> .env.local

# Start dev server
npm run dev
```

The frontend runs on `http://localhost:3000`.

### 3. ML Training (Optional)

```bash
cd ml

# Install Python dependencies
pip install scikit-learn pandas numpy matplotlib seaborn joblib

# Generate synthetic dataset
python generate_dataset.py

# Train models & export weights
python train_model.py
```

---

## ⚙️ Environment Variables

### Backend

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `DATABASE_URL` | — | **Yes** | PostgreSQL connection string |
| `JWT_SECRET` | `default-secret-change-me` | **Yes** | JWT signing key |
| `PORT` | `8080` | No | Server port |
| `GIN_MODE` | `debug` | No | `debug` or `release` |
| `JWT_EXPIRY_HOURS` | `72` | No | Token expiry |
| `ALLOWED_ORIGINS` | `http://localhost:3000` | No | CORS origins (comma-separated) |
| `LLM_PROVIDER` | — | No | `groq` to enable AI explanations |
| `LLM_API_KEY` | — | No | API key for LLM provider |
| `LLM_MODEL` | — | No | e.g. `llama3-70b-8192` |
| `GOOGLE_CLIENT_ID` | — | No | Google OAuth client ID |
| `ADMIN_EMAIL` | `admin@careermanifest.in` | No | Default admin email |
| `ADMIN_PASSWORD` | `Admin@123` | No | Default admin password |

### Frontend (`.env.local`)

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_API_URL` | Backend API URL (e.g. `http://localhost:8080`) |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | Google OAuth client ID |

---

## 📡 API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/register` | Register new user |
| `POST` | `/api/auth/login` | Login |
| `POST` | `/api/auth/google` | Google OAuth login |

### Protected (JWT Required)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/auth/profile` | Get current user profile |
| `GET` | `/api/questions` | Get active assessment questions |
| `POST` | `/api/assessment` | Submit assessment answers |
| `GET` | `/api/assessment` | List user's past assessments |
| `GET` | `/api/assessment/:id` | Get specific assessment result |
| `POST` | `/api/chat` | AI chatbot for result follow-ups |

### Admin (JWT + Admin Role)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/admin/stats` | Dashboard statistics |
| `GET` | `/api/admin/questions` | List all questions |
| `POST` | `/api/admin/questions` | Create a question |
| `PUT` | `/api/admin/questions/:id` | Update a question |

---

## 🌐 Deployment

### Current Production

- **Backend**: [Render](https://render.com) — `careermanifest-api.onrender.com`
- **Frontend**: [Vercel](https://vercel.com) — `carriermanifest.vercel.app`
- **Database**: [Neon](https://neon.tech) — Serverless PostgreSQL

### Deploy Your Own

#### Backend (Render)

1. Fork the repo
2. Create a new Web Service on Render
3. Set root directory to `backend/`
4. Build command: `go build -o server ./cmd/server/`
5. Start command: `./server`
6. Add environment variables (see table above)

#### Frontend (Vercel)

1. Import the repo on Vercel
2. Set root directory to `frontend/`
3. Framework preset: Next.js
4. Add `NEXT_PUBLIC_API_URL` and `NEXT_PUBLIC_GOOGLE_CLIENT_ID` env vars

---

## 🖥️ Screenshots

| Page | Description |
|------|-------------|
| **Landing** | Animated hero with CTA |
| **Assessment** | Smart conditional questions with progress bar |
| **Results** | Career rankings, risk analysis, salary chart, AI explanation |
| **Dashboard** | Past assessments history |
| **Admin** | Question management & statistics |

---

## 🧪 How It Works

```
User answers 30 questions
        │
        ▼
  ┌─────────────────┐
  │   Aggregator     │  Maps answers → 9 feature scores
  │  (aggregator.go) │  Uses ML feature importance for
  │                   │  cross-feature signal extraction
  └────────┬──────────┘
           │
           ▼
  ┌─────────────────┐
  │   Scorer         │  Dot-product: Profile × Weight Matrix
  │  (scorer.go)     │  Weight matrix trained by Random Forest
  └────────┬──────────┘
           │
           ▼
  ┌─────────────────┐
  │   Risk Engine    │  Calculates risk score (0-10)
  │  (risk.go)       │  Applies career-specific penalties
  └────────┬──────────┘
           │
           ▼
  ┌─────────────────┐
  │   Normalizer     │  Softmax + percentile normalization
  │  (normalize.go)  │  Final ranking with confidence %
  └────────┬──────────┘
           │
           ▼
  ┌─────────────────┐
  │   Enrichment     │  Salary projections, skills,
  │  (enrichment.go) │  colleges, preparation roadmap
  └────────┬──────────┘
           │
           ▼
  ┌─────────────────┐
  │   AI Explain     │  Optional Groq/Claude analysis
  │  (explain.go)    │  Detailed career reasoning
  └────────┬──────────┘
           │
           ▼
    Final Result Page
    with ML badge, rankings,
    risk analysis & chatbot
```

---

## 📄 License

MIT

---

<p align="center">
  Built for Indian students, by India 🇮🇳
  <br>
  <strong>CareerManifest</strong> — Your career, data-driven.
</p>
