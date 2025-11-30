# java_analyzer_service

Stateless analyzer for Java code.

## Responsibilities
- Expose `POST /api/analyzer/java`.
- Accept JSON payload with `files[] { path, content }`.
- Run `Checkstyle` CLI on provided files and return unified analysis JSON.

## Tech & Architecture
- Language: Golang service invoking `Checkstyle` via CLI.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence. Target sub-3s response for typical files.
