# csharp_analyzer_service

Stateless analyzer for C# code.

## Responsibilities
- Expose `POST /api/analyzer/csharp`.
- Accept JSON payload with `files[] { path, content }`.
- Use .NET SDK/Roslyn analyzers (e.g., `dotnet format analyze` or similar) to produce unified analysis JSON.

## Tech & Architecture
- Language: Golang service invoking `.NET` CLI tooling.
- Pattern: MVC internally (controllers → services → models).
- Stateless; no persistence. Keep analysis latency around 3 seconds for typical files.
