# Static Code Analyzer

A primitive but functional web service for static code analysis built with Golang, PostgreSQL, HTML5, CSS3, and JavaScript.

## Features

- **Dual Input Modes**: Upload ZIP archives or paste code snippets directly
- **Multi-Language Linters**: Uses popular CLI linters (ESLint, Pylint, golangci-lint, cppcheck, Rubocop, PHP linters) via Docker
- **Language-Agnostic Analysis**: Works on plain text (no AST required)
- **8+ Analysis Rules**: Detects common code issues
- **Docker-Based Analysis**: Platform-independent linter execution in isolated containers
- **Modern UI**: Bootstrap 5 with responsive design
- **MVC Architecture**: Clear separation of concerns
- **PostgreSQL Storage**: Persistent storage of analysis results

## Technology Stack

- **Frontend**: HTML5, CSS3, JavaScript, Bootstrap 5 (CDN)
- **Backend**: Golang
- **Database**: PostgreSQL
- **Architecture**: MVC (Model-View-Controller)

## Project Structure

```
/cmd/webapp/main.go              # Application entry point
/internal/app/
  /controllers/                  # HTTP handlers (Controllers)
  /models/                       # Data models
  /repositories/                 # Database access layer
  /services/                     # Business logic
  /views/templates/              # HTML templates (Views)
/internal/db/
  /migrations/                   # SQL migrations
/web/static/                     # CSS and JavaScript files
/storage/                        # Uploaded ZIP files (gitignored)
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Docker and Docker Compose (for linter analysis)
- VS Code (or any Go-compatible IDE)

## Setup Instructions

### 1. Install Dependencies

```bash
go mod download
```

### 2. Set Up Docker Services

The project includes Docker Compose configuration for both PostgreSQL and the linter container.

Start all services:

```bash
docker-compose up -d
```

This will start:
- **PostgreSQL**: Database for storing analysis results
- **Linters Container**: Container with multiple language linters (ESLint, Pylint, golangci-lint, cppcheck, Rubocop, PHP linters)

The linter container must be running for Docker-based analysis to work. If it's not running, the analyzer will fall back to simple pattern-based rules only.

To rebuild the linter container after changes:

```bash
docker-compose build linters
docker-compose up -d linters
```

#### Option B: Local PostgreSQL (without Docker)

1. Install PostgreSQL on your system
2. Create a database:
   ```sql
   CREATE DATABASE static_analyzer;
   ```
3. Note: You'll still need Docker for the linter container, or analysis will use only basic pattern-based rules

### 3. Configure Database Connection

The application uses the following default connection string:
```
postgres://postgres:postgres@localhost:5432/static_analyzer?sslmode=disable
```

To use a custom connection, set the `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/static_analyzer?sslmode=disable"
```

### 4. Run Migrations

Migrations run automatically when the server starts. The application will create all necessary tables on first run.

### 5. Run the Server

From the project root directory:

```bash
go run cmd/webapp/main.go
```

Or build and run:

```bash
go build -o bin/webapp cmd/webapp/main.go
./bin/webapp
```

The server will start on `http://localhost:8080` by default.

To use a custom port, set the `PORT` environment variable:

```bash
export PORT=3000
go run cmd/webapp/main.go
```

## Usage

1. **Home Page** (`/`): Overview and navigation
2. **Upload Page** (`/upload`): 
   - Upload a ZIP file (max 20MB)
   - Paste a code snippet (max 50,000 characters)
3. **History** (`/analyses`): View all past analyses
4. **Analysis Details** (`/analyses/{id}`): View detailed results with issues
5. **About** (`/about`): Information about the service

## Analysis Rules

The analyzer uses a two-tier approach:

### Docker-Based Linters (when container is running)

The analyzer automatically detects file types and runs appropriate linters:

- **JavaScript/TypeScript**: ESLint
- **Python**: Pylint
- **Go**: golangci-lint (errcheck, govet, staticcheck)
- **C/C++**: cppcheck
- **Ruby**: Rubocop
- **PHP**: PHP syntax checker

### Pattern-Based Rules (always active)

These rules work on all file types:

1. **LINE_TOO_LONG** (warn) - Lines exceeding 120 characters
2. **TRAILING_WHITESPACE** (info) - Trailing spaces or tabs
3. **TAB_INDENT** (info) - Tab characters in code
4. **TODO_FOUND** (warn) - TODO or FIXME comments
5. **DEBUG_PRINT** (info) - Debug print statements (console.log, fmt.Println, print)
6. **HARDCODED_SECRET** (error) - Potential hardcoded secrets (password, apikey, secret, token)
7. **INSECURE_HTTP_URL** (warn) - Insecure HTTP URLs
8. **FILE_TOO_LARGE** (warn) - Files exceeding 800 lines

Both linter results and pattern-based rule results are combined in the final analysis report.

## Supported File Extensions

The analyzer processes files with these extensions:
`.go`, `.cs`, `.js`, `.ts`, `.py`, `.java`, `.cpp`, `.h`, `.hpp`, `.php`, `.rb`, `.kt`, `.swift`, `.sql`, `.html`, `.css`, `.md`

## Limitations

- ZIP uploads limited to 20 MB
- Code snippets limited to 50,000 characters
- Synchronous analysis (no background processing)
- Docker container must be running for linter-based analysis
- Some linters may require configuration files for optimal results (ESLint uses --no-eslintrc by default)

## Safety Features

- Path traversal protection for ZIP files
- File size limits enforced
- Extension whitelist for security
- Symlink detection and prevention

## Development

### Running from VS Code

1. Open the project in VS Code
2. Install the Go extension
3. Set up your launch configuration (optional):
   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Launch",
         "type": "go",
         "request": "launch",
         "mode": "auto",
         "program": "${workspaceFolder}/cmd/webapp",
         "env": {
           "DATABASE_URL": "postgres://postgres:postgres@localhost:5432/static_analyzer?sslmode=disable"
         }
       }
     ]
   }
   ```
4. Press F5 to run or use the Run button

## Database Schema

### Tables

- **projects**: Stores project information
- **analysis_runs**: Stores analysis execution records
- **issues**: Stores detected issues from analyses

See `internal/db/migrations/001_init.sql` for the complete schema.

## License

This is a university course project.

