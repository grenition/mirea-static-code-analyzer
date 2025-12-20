You are an expert Golang web developer. Design and implement the **most primitive but valid** web service for **static code analysis** that satisfies ALL requirements below. Keep the solution simple, deterministic, and easy to defend in a university course project.

## Mandatory constraints (must be present)
- Tech stack: **HTML5, CSS3, JavaScript, Golang, PostgreSQL**
- Server architecture pattern: **MVC** (Controllers / Models / Views must be explicit)
- UI requirement: **modern-looking pages** (Bootstrap 5 via CDN is allowed)
- Must have **multi-page navigation** (header menu visible on all pages)
- Must include: server logic layer, database layer, client presentation layer
- Must be runnable locally from VS Code

## Product idea (keep it primitive)
Build a **language-agnostic** static analyzer: it scans source code as plain text (no AST).  
User can either:
1) upload a `.zip` archive with multiple files, OR  
2) **paste a code snippet directly in the browser** (single “virtual file”).

The service analyzes the input with simple rules, stores results in PostgreSQL, and renders reports.

---

## Pages (must be real separate routes + templates)
At minimum implement these pages:
1) `GET /` — Home (service description + buttons)
2) `GET /upload` — Upload page with TWO modes (Zip + Paste)
3) `POST /upload/zip` — Accept `.zip` and run analysis
4) `POST /upload/snippet` — Accept pasted code and run analysis
5) `GET /analyses` — Analyses history list
6) `GET /analyses/{id}` — Analysis details (summary + issues table)
7) `GET /about` — About page (stack + limitations)

### Navigation (required)
Global header menu (in `layout.html`):
- Home
- Upload
- History
- About

---

## Upload page UI (must support BOTH inputs)
On `/upload`, render:
### A) Zip Upload Form
Fields:
- Project name (text)
- Zip file input (`accept=".zip"`)
- Submit button “Upload ZIP & Analyze”

### B) Paste Code Form
Fields:
- Project name (text)
- “Language / file extension” selector (optional but recommended):
  - dropdown: `.go .cs .js .ts .py .java .cpp .php .kt .sql .html .css .txt`
  - default: `.txt`
- Code textarea (required)
- Submit button “Analyze Snippet”

Both forms should create a new analysis run and redirect to `/analyses/{id}`.

---

## Minimal analysis pipeline
### Input options
#### 1) ZIP mode
- Accept only `.zip`
- Limit size: **20 MB**
- Safely unzip into `./storage/{analysis_id}/`
  - Prevent path traversal (`../`, absolute paths)
  - Ignore symlinks
- Analyze only allowed extensions (whitelist):  
  `.go, .cs, .js, .ts, .py, .java, .cpp, .h, .hpp, .php, .rb, .kt, .swift, .sql, .html, .css, .md`

#### 2) Snippet mode (pasted code)
- Accept textarea input up to **50,000 chars** (or 200 KB) to keep it safe
- Treat snippet as a single virtual file:
  - file path stored as: `snippet{ext}` (e.g. `snippet.go` or `snippet.txt`)
- Do NOT write snippet to disk unless you want to (optional).  
  For the minimal solution, you can analyze it in memory.

### Rules (implement at least 8)
Produce `Issue` entries with: `severity`, `rule_code`, `message`, `file_path`, `line`.
Use simple string checks / regex. Suggested rules:
1) `LINE_TOO_LONG` — line length > 120 (warn)
2) `TRAILING_WHITESPACE` — trailing spaces/tabs (info)
3) `TAB_INDENT` — contains `\t` (info)
4) `TODO_FOUND` — contains `TODO` or `FIXME` (warn)
5) `DEBUG_PRINT` — contains `console.log(` OR `fmt.Println(` OR `print(` (info)
6) `HARDCODED_SECRET` — regex for `password|apikey|secret` near `=` and quotes (error)
7) `INSECURE_HTTP_URL` — contains `http://` (warn)
8) `FILE_TOO_LARGE` — file lines > 800 (warn)  
   - In snippet mode this can still apply (if pasted code is huge)
9) (optional) `DUPLICATE_FILE` — identical files by sha256 (info)

### AnalysisRun statuses
`Created -> Running -> Done / Failed`

Even if analysis is synchronous:
- create run `Running`
- on success set `Done`
- on error set `Failed` (optional: store error_text)

---

## Database (PostgreSQL) — must be real SQL + migrations
Implement SQL migrations in `/internal/db/migrations`.

### Tables (minimum)
**projects**
- id UUID PK
- name TEXT
- created_at TIMESTAMP

**analysis_runs**
- id UUID PK
- project_id UUID FK -> projects.id
- status TEXT
- created_at TIMESTAMP
- started_at TIMESTAMP NULL
- finished_at TIMESTAMP NULL
- input_type TEXT NOT NULL  -- 'zip' | 'snippet'
- input_meta JSONB NULL     -- optional: ext, file_count, size_bytes
- summary_json JSONB NULL   -- counts by severity (recommended)

**issues**
- id UUID PK
- analysis_id UUID FK -> analysis_runs.id
- severity TEXT
- rule_code TEXT
- message TEXT
- file_path TEXT
- line INT NULL

### Indexes (minimum)
- `issues(analysis_id)`
- `analysis_runs(project_id, created_at desc)`
- `analysis_runs(input_type)` (optional)

---

## MVC structure (must be obvious)
Provide a clear project tree and follow it.

Recommended structure:
```
/cmd/webapp/main.go

/internal/app/controllers/
home_controller.go
upload_controller.go // GET /upload + POST handlers for zip/snippet
analyses_controller.go
about_controller.go

/internal/app/models/
project.go
analysis_run.go
issue.go

/internal/app/repositories/
project_repo.go
analysis_repo.go
issue_repo.go

/internal/app/services/
analyzer_service.go // core rule engine, works on []FileInput
storage_service.go // zip unpacking + safe file enumeration

/internal/app/views/templates/
layout.html
home.html
upload.html // contains BOTH forms
analyses.html
analysis_details.html
about.html
errors/404.html
errors/500.html

/internal/db/migrations/
001_init.sql

/web/static/
css/styles.css
js/app.js

/storage/ (gitignored)
```


### Service interface suggestion (to keep both modes unified)
Represent inputs as:
- `type FileInput struct { Path string; Content []byte }`
Zip mode -> build multiple `FileInput` by reading files.  
Snippet mode -> single `FileInput{ Path: "snippet"+ext, Content: []byte(snippet) }`.

AnalyzerService returns:
- list of `Issue` + summary counts.

---

## UI requirements (keep it simple but “modern”)
- Use Bootstrap 5 via CDN OR decent custom CSS
- All pages share the same layout (header/footer)
- History page: table of analyses (project, date, input_type, status, link)
- Details page:
  - summary cards: counts by severity
  - issues table with filter by severity using small JS (client-side filtering)

---

## Output requirements (what you must produce)
1) Minimal runnable code (Go web server + templates + migrations)
2) README with:
   - how to start Postgres (docker-compose optional)
   - how to run migrations
   - how to run server
3) Demonstrate multi-page navigation + MVC in code organization
4) Deterministic analysis results (same input -> same issues)
5) Upload supports BOTH ZIP and pasted code snippet

---

## Acceptance checklist (must pass)
- [ ] Multi-page navigation present and visible
- [ ] MVC pattern clearly implemented in folders and usage
- [ ] Postgres schema + migrations exist and work
- [ ] Upload zip -> analysis -> save -> view history + details works
- [ ] Paste snippet -> analysis -> save -> view history + details works
- [ ] Modern UI (bootstrap or equivalent)
- [ ] Safety: zip path traversal blocked + size limited + extension whitelist
- [ ] Snippet input limited (length cap) + validated (non-empty)

Now implement it.
