# python_analyzer_service

Stateless analyzer for Python code.

## Responsibilities
- Expose `POST /api/analyzer/python` endpoint.
- Accept JSON payload with `files[] { path, content }`.
- Run `flake8` against provided files and return unified analysis JSON (`comment`, `line_comments`).

## Tech & Architecture
- Language: Golang service shelling out to `flake8`.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence.
- HTTP JSON via API gateway; respond within ~3 seconds for typical files.
