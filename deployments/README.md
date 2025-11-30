# Deployments (docker-compose & API Gateway)

Deployment layer using `docker-compose` plus Nginx gateway.

## Composition
- `nginx` gateway container routing external HTTP to internal services.
- `frontend` web app container.
- `user_identity_service` and `projects_service` containers (Golang + PostgreSQL dependencies).
- Analyzer microservices containers: python, java, javascript, csharp, cpp, json.
- `postgres` container (shared by stateful services).
- Shared internal network for service-to-service JSON/HTTP.

## Gateway Responsibilities
- Route external traffic to frontend and backend services.
- Enforce HTTPS termination if configured.
- Forward requests to analyzers, `user_identity_service`, and `projects_service`.
- Keep routing stateless; no business logic inside gateway.

## docker-compose Notes
- Define services, networks, and volumes in `docker-compose.yml` (to be added alongside this file).
- Map `/api/*` paths to backend services; serve frontend assets via Nginx.
- Use environment variables for DB credentials, JWT secrets, and analyzer tool paths; never bake secrets into images.
- Keep logs free of internal paths or stack traces; prefer structured logging.

## Bring-up (example)
1. Ensure required env files exist (e.g., `deployments/.env` with DB creds and service secrets).
2. From repo root: `docker-compose -f deployments/docker-compose.yml up --build`.
3. Access via gateway (default `http://localhost:80` or configured port).
