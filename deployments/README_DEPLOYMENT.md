# Deployment Instructions

## Prerequisites

- Docker and Docker Compose installed
- At least 4GB of free disk space
- Ports 80 and 5432 available

## Quick Start

1. Navigate to the deployments directory:
   ```bash
   cd deployments
   ```

2. (Optional) Create a `.env` file with your JWT secret:
   ```bash
   echo "JWT_SECRET=your-secret-key-here" > .env
   ```

3. Start all services:
   ```bash
   docker-compose up --build
   ```

4. Access the application:
   - Frontend: http://localhost:8080
   - API Gateway: http://localhost:8080/api

## Services

The system consists of the following services:

- **postgres**: PostgreSQL database
- **user_identity_service**: User authentication and registration (port 8080)
- **projects_service**: Project and file management (port 8081)
- **python_analyzer_service**: Python code analysis with flake8 (port 8082)
- **javascript_analyzer_service**: JavaScript code analysis with ESLint (port 8083)
- **java_analyzer_service**: Java code analysis with Checkstyle (port 8084)
- **cpp_analyzer_service**: C/C++ code analysis with cppcheck (port 8085)
- **csharp_analyzer_service**: C# code analysis with .NET CLI (port 8086)
- **json_analyzer_service**: JSON validation (port 8087)
- **gateway**: Nginx API gateway (port 80)
- **frontend**: React web application

## Testing

### Unit Tests

Run unit tests for individual services:

```bash
# User Identity Service
cd ../services/user_identity_service
go test ./internal/service/... -v

# Projects Service
cd ../projects_service
go test ./internal/service/... -v

# Analyzer Services
cd ../analyzers/python_analyzer_service
go test ./internal/service/... -v
```

### Integration Tests

Integration tests require a running database. They will skip if the database is not available:

```bash
cd ../services/user_identity_service
TEST_DATABASE_URL="postgres://user:password@localhost:5432/test_user_identity_db?sslmode=disable" go test -v
```

## Troubleshooting

### Services not starting

1. Check Docker logs:
   ```bash
   docker-compose logs [service_name]
   ```

2. Verify database connection:
   ```bash
   docker-compose exec postgres psql -U user -d user_identity_db
   ```

### Port conflicts

If port 80 is already in use, modify the `docker-compose.yml` file to use a different port for the gateway service.

### Database initialization

The databases are automatically created on first startup via the `init-db.sql` script.

## Stopping Services

```bash
docker-compose down
```

To remove volumes (this will delete all data):
```bash
docker-compose down -v
```

