# user_identity_service

Authenticates users and manages credentials.

## Responsibilities
- Register users and issue login tokens (e.g., JWT or session id).
- Hash passwords securely (bcrypt or equivalent).
- Enforce ownership checks: only project owners access their data.

## Interfaces
- `POST /api/users/register` – create user.
- `POST /api/users/login` – authenticate user and return token/session.
- Communicates with other services via shared auth tokens through the API gateway.

## Tech & Architecture
- Language: Golang.
- Persistence: PostgreSQL.
- Pattern: MVC (controllers → services → repositories → models).
- Errors must be user-friendly; no sensitive details in responses.
