# Task Management API

A production-ready RESTful API built with Go for managing tasks and projects with user authentication, role-based access control, and team collaboration features.

## Features

- **User Authentication**: JWT-based authentication with access and refresh tokens
- **Project Management**: Create, update, and delete projects with full CRUD operations
- **Task Management**: Complete task lifecycle management with status tracking
- **Collaboration**: Add team members to projects with role-based permissions
- **Role-Based Access Control**: Owner, Editor, and Viewer roles with permission enforcement
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Production-Ready**: Docker containerization, proper error handling, and database migrations
- **Security**: Password hashing with bcrypt, SQL injection prevention, and secure session management

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (lightweight HTTP framework)
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## Project Structure

```
task-management-api/
├── cmd/main.go                 # Application entry point
├── pkg/
│   ├── config/                 # Configuration management
│   ├── models/                 # Data models and DTOs
│   ├── controllers/            # HTTP handlers
│   ├── services/               # Business logic
│   ├── repositories/           # Data access layer
│   ├── middleware/             # HTTP middleware
│   ├── utils/                  # Utilities (JWT, errors, validators)
│   └── database/               # Database connection and migrations
├── Dockerfile                  # Production container image
├── docker-compose.yml          # Local development setup
├── go.mod & go.sum            # Go module files
└── README.md                   # This file
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd task-management-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file from template:
```bash
cp .env.example .env
```

4. Update environment variables in `.env` file

### Running Locally

#### Option 1: With Docker Compose (Recommended)

```bash
docker-compose up
```

The API will be available at `http://localhost:8080`

#### Option 2: Manual Setup

1. Start PostgreSQL:
```bash
# Ensure PostgreSQL is running on localhost:5432
```

2. Run the application:
```bash
go run cmd/main.go
```

## API Endpoints

### Authentication

```
POST   /api/v1/auth/register          # Register new user
POST   /api/v1/auth/login             # Login and get tokens
POST   /api/v1/auth/refresh           # Refresh access token
GET    /api/v1/auth/me                # Get current user profile
PUT    /api/v1/auth/me                # Update user profile
```

### Projects

```
POST   /api/v1/projects               # Create project
GET    /api/v1/projects               # List user's projects
GET    /api/v1/projects/{id}          # Get project details
PUT    /api/v1/projects/{id}          # Update project
DELETE /api/v1/projects/{id}          # Delete project
GET    /api/v1/projects/{id}/members  # List project members
POST   /api/v1/projects/{id}/members  # Add member to project
DELETE /api/v1/projects/{id}/members/{userId}  # Remove project member
```

### Tasks

```
POST   /api/v1/projects/{projectId}/tasks     # Create task
GET    /api/v1/projects/{projectId}/tasks     # List project tasks
GET    /api/v1/tasks/{taskId}                 # Get task details
PUT    /api/v1/tasks/{taskId}                 # Update task
PATCH  /api/v1/tasks/{taskId}/status          # Update task status
DELETE /api/v1/tasks/{taskId}                 # Delete task
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Example Login Flow

1. Register user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

2. Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

3. Use access token:
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

## API Documentation

Swagger documentation is available at: `http://localhost:8080/swagger/index.html`

To regenerate Swagger docs after code changes:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

## Database Schema

### Users Table
- id (UUID, primary key)
- email (unique, indexed)
- password (hashed with bcrypt)
- name
- role (admin, user, viewer)
- created_at, updated_at

### Projects Table
- id (UUID, primary key)
- owner_id (foreign key → users)
- name (indexed)
- description
- status (active, archived)
- created_at, updated_at

### Tasks Table
- id (UUID, primary key)
- project_id (foreign key → projects)
- assigned_to (nullable, foreign key → users)
- title (indexed)
- description
- status (todo, in_progress, done)
- priority (low, medium, high)
- due_date (nullable)
- created_at, updated_at

### Project Members Table
- id (UUID, primary key)
- project_id (foreign key → projects)
- user_id (foreign key → users)
- role (owner, editor, viewer)
- created_at

## Error Handling

All error responses follow a consistent format:

```json
{
  "error": "error_type",
  "message": "Human-readable error message",
  "code": 400
}
```

Common HTTP Status Codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input or validation error
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Security Considerations

1. **Password Security**: All passwords are hashed with bcrypt before storage
2. **JWT Tokens**: 
   - Access tokens expire after 1 hour
   - Refresh tokens expire after 7 days
   - Tokens are stateless and validated on each request
3. **SQL Injection Prevention**: All database queries use parameterized statements
4. **Input Validation**: All inputs are validated before processing
5. **CORS**: Configured to allow only specified origins
6. **Role-Based Access**: Permissions checked at service layer

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| APP_ENV | development | Application environment |
| APP_PORT | 8080 | Server port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | PostgreSQL user |
| DB_PASSWORD | postgres | PostgreSQL password |
| DB_NAME | task_management_db | Database name |
| DB_SSL_MODE | disable | PostgreSQL SSL mode |
| JWT_SECRET | change-me | Secret key for JWT signing |
| JWT_EXPIRATION | 3600 | Access token expiration (seconds) |
| JWT_REFRESH_EXPIRATION | 604800 | Refresh token expiration (seconds) |
| CORS_ALLOWED_ORIGINS | http://localhost:3000 | Comma-separated allowed origins |

## Development

### Code Structure Guidelines

1. **Models** (`pkg/models/`): Data structures and DTOs
2. **Repositories** (`pkg/repositories/`): Database access layer
3. **Services** (`pkg/services/`): Business logic layer
4. **Controllers** (`pkg/controllers/`): HTTP request handlers
5. **Middleware** (`pkg/middleware/`): Request/response processing
6. **Utils** (`pkg/utils/`): Helper functions

### Running Tests

```bash
go test ./...
```

### Linting

```bash
go vet ./...
golint ./...
```

## Deployment

### Docker Deployment

1. Build the image:
```bash
docker build -t task-management-api:latest .
```

2. Run the container:
```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  task-management-api:latest
```

### Production Checklist

- [ ] Change JWT_SECRET in environment variables
- [ ] Set DB_SSL_MODE=require for production PostgreSQL
- [ ] Configure CORS_ALLOWED_ORIGINS for your domain
- [ ] Set APP_ENV=production
- [ ] Use strong database password
- [ ] Enable HTTPS
- [ ] Configure database backups
- [ ] Set up monitoring and logging

## Performance Optimization

1. **Connection Pooling**: Configured with max connections and idle timeouts
2. **Database Indexes**: Indexes on frequently queried columns (email, project_id, title)
3. **Pagination**: List endpoints support offset/limit pagination
4. **Filtering**: Task endpoint supports filtering by status and priority

## Troubleshooting

### Database Connection Error
```
Check that PostgreSQL is running and credentials in .env are correct
```

### JWT Token Expired
```
Use the refresh endpoint to get a new access token
POST /api/v1/auth/refresh with your refresh token
```

### Permission Denied
```
Ensure you have the required role in the project (owner, editor)
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or suggestions, please open an issue in the repository or contact the development team.
