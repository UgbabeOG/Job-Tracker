# JobTrack

A web application for job seekers to organize applications, track deadlines, and compare CVs against job descriptions.

## Tech stack

| Layer | Technology |
|---|---|
| Frontend | React, TypeScript, Vite |
| Backend | Go, `net/http` |
| Database | SQLite |
| API style | REST / JSON |

## Quick start

### Backend

```bash
cd backend
go run .
# API available at http://localhost:8080
```

See [`backend/README.md`](backend/README.md) for full setup instructions, environment variables, and how to run tests.

### Frontend

```bash
cd frontend
npm install
npm run dev
# UI available at http://localhost:5173
```

See [`frontend/README.md`](frontend/README.md) for setup details.

## Documentation

| Document | Description |
|---|---|
| [`docs/requirements.md`](docs/requirements.md) | Product requirements (PRD) |
| [`docs/api-contract.md`](docs/api-contract.md) | REST API contract |
| [`docs/test-plan.md`](docs/test-plan.md) | Test plan |

## Recent changes

- Added a `job_description` field to job applications (backend model, database, API, and frontend form).
- Keyword matching now uses `job_description` instead of notes or job title.
- Keyword Matcher page shows scores for all saved jobs and allows matching an individual job.
- Jobs without a description show a "No job description available" message.
