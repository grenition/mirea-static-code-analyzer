# Static code analyzer

## Requirements

### Functional

#### 1. Static analysis of program code in different languages

Supported languages:

- Python
- Java
- JavaScript
- C#
- C++
- Json

User can analyze file within project or without them. 

#### 2. Users identity

Authorization by login and password.

#### 3. Storing user projects

User can create/read/delete/update any count of projects.

Project is:

- Name of project
- Files with paths (need to save origin project structure)

Project can be created with:

- Uploading archive file: .zip, max size: 25 mb
- Manually file by file

Only project owner can access project data

#### 3. UI/UX

Web site with modern design.

Pages:

- /registration
- /authorization
- /home - tutorial page with navigation to projects page or to single file page
- /projects - list of projects with availabillity to create/delete/edit projects (preferred the same component that in home page)
- /projects/{id} - list of project files
- /projects/{id}/{file} - file edit and analysis with breadcrumbs (note that /projects/{id} and /projects/{id}/{file} should be a single page that renders on client side to provide fluent user experience)
- /sandbox - single file edit and analysis

Page navigation bar is required

Analysis should be started automatically after file appearance or changing.

User can select analyzer manually.

### Non functional requirements

- Low latency of analysis, less than 3 seconds
- Low latency on crud operations
- Error messages must be user-friendly and not expose internal details
- System must follow MVC architecture

## Architecture

### Backend

#### API Gateway (nginx)

#### Analysers microservices group:

5 stateless microservices that only analyze code for specified language:

##### python_analyzer_service 
use library: flake8
api: api/analyzer/python
language: golang
achitecture: mvc

##### java_analyzer_service 
use library: Checkstyle
api: api/analyzer/java
language: golang
achitecture: mvc

##### javascript_analyzer_service 
use library: ESLint
api: api/analyzer/javascript
language: golang
achitecture: mvc

##### csharp_analyzer_service 
use dotnet sdk
api: api/analyzer/csharp
language: golang
achitecture: mvc

##### cpp_analyzer_service 
use library: cppcheck
api: api/analyzer/cpp
language: golang
achitecture: mvc

##### json_analyzer_service 
use library: github.com/xeipuuv/gojsonschema
api: api/analyzer/json
language: golang
achitecture: mvc

common api contract:

GET:
```json
{
    "files" : [
        {
            "path" : "path/to/file",
            "content" : "content"
        },
        {
            "path" : "path/to/file/2",
            "content" : "content"
        }
    ]
}
```
returns
```json
{
    "files" : [
        {
            "path" : "path/to/file",
            "comment" : "verdict",
            "line_comments" : [
                { "0" : "bad syntax" },
                { "13" : "no all paths return value" }
            ]
        },
        {
            "path" : "path/to/file/2",
            "comment" : "OK"
        }
    ]
}
```

#### user_identity_service

Stateful service for authenticating users

Database: PostgreSQL

Contracts:

- POST: api/users/register
- POST: api/users/login

language: golang
achitecture: mvc

#### projects_service

Stateful service for storing user projects

Database: PostgreSQL

Contracts:

- api/projects/..

language: golang
achitecture: mvc

#### Communication between services:

- Synchronous, HTTP

### Frontend

Pages should use shared components:

- code review block (block with code and line comments if provided by analyzers)
- project list (with control elements)
- file structure (in projects) 
- auth form
- register form
etc..

Should be used modern styles like Tailwind CSS.

