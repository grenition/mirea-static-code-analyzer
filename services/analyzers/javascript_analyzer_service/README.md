# javascript_analyzer_service

Stateless analyzer for JavaScript code.

## Responsibilities
- Expose `POST /api/analyzer/javascript`.
- Accept JSON payload with `files[] { path, content }`.
- Run `ESLint` on provided files and return unified analysis JSON (`comment`, `line_comments`).

## Tech & Architecture
- Language: Golang service invoking `ESLint` via CLI.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence. Aim for low-latency (<3s) responses.
