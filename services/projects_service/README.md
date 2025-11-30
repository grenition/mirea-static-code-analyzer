# projects_service

Stores and manages user projects and files.

## Responsibilities
- CRUD for projects and files while preserving original structure.
- Accept project creation via archive upload (`.zip`, max 25 MB) or manual file creation with arbitrary nested paths.
- Delete all files and related analysis results when a project is removed.
- Trigger analyzer requests and enforce ownership checks.

## Interfaces (examples)
- `GET /api/projects` – list authenticated user projects.
- `POST /api/projects` – create project.
- `GET /api/projects/{id}` – fetch project metadata.
- `PUT /api/projects/{id}` – update project.
- `DELETE /api/projects/{id}` – delete project and related data.
- `GET /api/projects/{id}/files` – list files.
- CRUD `/api/projects/{id}/files/{fileId}` – manage files and their analysis metadata.

## Tech & Architecture
- Language: Golang.
- Persistence: PostgreSQL.
- Pattern: MVC (controllers → services → repositories → models).
- JSON over HTTP through the API gateway; low-latency operations expected.
