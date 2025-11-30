# Frontend Web App

Implements the UI/UX for the static code analyzer.

## Requirements
- Modern SPA-like experience with shared navigation bar on all main pages.
- Pages: `/registration`, `/authorization`, `/home`, `/projects`, `/projects/{id}` (with file tree), `/projects/{id}/{file}` (editor + analysis, same view as `/projects/{id}`), `/sandbox`.
- Automatic analysis trigger on file changes with debouncing; user can pick analyzer/tool when multiple are available.
- Shared UI components: code review block (line comments), project list CRUD, file tree, auth forms, navigation bar.
- Styling: Tailwind CSS (or similar modern utility-first styling) and modern browser compatibility.

## Behavior
- Communicates with backend via HTTP JSON through the API gateway.
- Renders breadcrumbs for file paths and keeps project/file view as a single client-side experience.
- Surfaces user-friendly errors; avoid exposing internal details.
