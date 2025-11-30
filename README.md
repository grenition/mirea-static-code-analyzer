# Static code analyzer

## Requirements

### Functional

#### 1. Static analysis of program code in different languages

Supported languages:

* Python
* Java
* JavaScript
* C#
* C++
* JSON

User can analyze a file within a project or as a standalone file (sandbox).

#### 2. User identity

* Authorization by login and password.
* Passwords must be stored as a secure cryptographic hash (e.g. bcrypt).
* Only the project owner can access and manage project data.

#### 3. Storing user projects

User can create, read, update and delete any number of projects.

Project is:

* Project name
* Files with paths (original project structure must be preserved)

Project can be created with:

* Uploading an archive file: `.zip`, max size: **25 MB**
* Creating files manually (file by file, with arbitrary nested paths)

When a project is deleted, all its files and related analysis results must be deleted as well.

#### 4. UI/UX

Web site with a modern design.

Pages:

* `/registration` – user registration
* `/authorization` – user login
* `/home` – tutorial/landing page with navigation to projects page and sandbox
* `/projects` – list of projects with ability to create / edit / delete projects (preferably using the same UI component as on `/home`)
* `/projects/{id}` – project view with file tree
* `/projects/{id}/{file}` – file edit and analysis with breadcrumbs

  * `/projects/{id}` and `/projects/{id}/{file}` must be implemented as a single client-side page (SPA-like) for fluent user experience
* `/sandbox` – single file edit and analysis (without project)

Additional UI requirements:

* Global navigation bar is required on all main pages.
* Analysis should be triggered automatically after a file appears or changes (with reasonable debouncing on the client).
* User can select the analyzer (language / tool) manually when multiple choices are available.

### Non functional requirements

* Low latency of analysis: less than **3 seconds** for typical files.
* Low latency on CRUD operations with projects and files.
* Error messages must be user-friendly and must not expose internal details (no stack traces, internal paths, etc.).
* System must follow MVC architecture inside backend services (clear separation of controllers, business logic and data access).
* Web UI must use modern styles (e.g. Tailwind CSS) and render correctly in modern browsers (Chrome, Firefox, Edge).

---

## Architecture

### Backend

#### API Gateway (Nginx)

See [deployments/README.md](./deployments/README.md) for gateway and docker-compose architecture.

#### Analyzer microservices group

Stateless microservices that only analyze code for a specified language.

Each analyzer:

* Is stateless.
* Implements MVC internally (controllers → services → models).
* Exposes a single HTTP API endpoint.
* Receives file contents in the request body and returns analysis results in a unified JSON format.

##### python_analyzer_service

See [services/analyzers/python_analyzer_service/README.md](./services/analyzers/python_analyzer_service/README.md).

##### java_analyzer_service

See [services/analyzers/java_analyzer_service/README.md](./services/analyzers/java_analyzer_service/README.md).

##### javascript_analyzer_service

See [services/analyzers/javascript_analyzer_service/README.md](./services/analyzers/javascript_analyzer_service/README.md).

##### csharp_analyzer_service

See [services/analyzers/csharp_analyzer_service/README.md](./services/analyzers/csharp_analyzer_service/README.md).

##### cpp_analyzer_service

See [services/analyzers/cpp_analyzer_service/README.md](./services/analyzers/cpp_analyzer_service/README.md).

##### json_analyzer_service

See [services/analyzers/json_analyzer_service/README.md](./services/analyzers/json_analyzer_service/README.md).

**Common analyzer API contract**

Request (POST):

```json
{
  "files": [
    {
      "path": "path/to/file",
      "content": "content"
    },
    {
      "path": "path/to/file/2",
      "content": "content"
    }
  ]
}
```

Response:

```json
{
  "files": [
    {
      "path": "path/to/file",
      "comment": "verdict",
      "line_comments": [
        {
          "line": 0,
          "comment": "bad syntax"
        },
        {
          "line": 13,
          "comment": "not all code paths return a value"
        }
      ]
    },
    {
      "path": "path/to/file/2",
      "comment": "OK",
      "line_comments": []
    }
  ]
}
```

* `comment` – overall summary/verdict for the file.
* `line_comments` – optional list of per-line issues.

#### user_identity_service

See [services/user_identity_service/README.md](./services/user_identity_service/README.md).

#### projects_service

See [services/projects_service/README.md](./services/projects_service/README.md).

#### Communication between services

* Synchronous HTTP communication over internal network.
* Data format: JSON.
* API gateway routes external requests to appropriate internal services.
* Services authenticate requests using shared tokens / session data from `user_identity_service`.

---

### Frontend
See [frontend/README.md](./frontend/README.md).
