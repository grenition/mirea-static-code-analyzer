# cpp_analyzer_service

Stateless analyzer for C/C++ code.

## Responsibilities
- Expose `POST /api/analyzer/cpp`.
- Accept JSON payload with `files[] { path, content }`.
- Run `cppcheck` on provided files and return unified analysis JSON (`comment`, `line_comments`).

## Tech & Architecture
- Language: Golang service invoking `cppcheck` via CLI.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence. Aim for sub-3-second responses on typical files.
