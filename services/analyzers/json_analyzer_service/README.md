# json_analyzer_service

Stateless analyzer for JSON documents.

## Responsibilities
- Expose `POST /api/analyzer/json`.
- Accept JSON payload with `files[] { path, content }`.
- Validate using `github.com/xeipuuv/gojsonschema` (or equivalent) and return unified analysis JSON.

## Tech & Architecture
- Language: Golang.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence. Maintain quick response (~3s) for typical payloads.
