# Kemuko Support Service

A Go-based microservice for Kemuko's educational technology support and ticketing system with PostgreSQL database.

## Features

- Kemuko student support ticket system with EdTech-specific categories
- Course-related ticket tracking with course ID integration
- Student/Instructor role-based access control for Kemuko platform
- Comment system for student-instructor communication
- File attachment support via URLs (separate file service)
- Ticket history and audit trail
- JSONB metadata for flexible educational data storage
- JWT authentication with user context extraction
- Email and Slack notifications for Kemuko admins
- Slack bot integration for admin replies
- RESTful API with proper authentication middleware

## Project Structure

```
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and helpers
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── migrations/          # Database migrations
├── sql/
│   └── queries/         # SQL queries
├── pkg/
│   ├── middleware/      # HTTP middleware
│   └── utils/           # Utility functions
└── docs/                # Documentation
```

## Database Schema

### Tables:
- **categories**: Educational support categories (technical, course, assignment, grading, etc.)
- **tickets**: Student support tickets with course integration and instructor assignment
- **ticketComments**: Student-instructor communication with internal/external visibility
- **attachments**: File references via URLs (files handled by separate service)
- **ticketHistory**: Complete audit trail of all ticket changes

All tables use camelCase column naming and include JSONB metadata fields for educational context.

## Setup

1. **Database Setup**:
   ```bash
   # Create PostgreSQL database
   createdb community_support
   
   # Run migrations (use a migration tool like golang-migrate)
   migrate -path migrations -database "postgres://user:password@localhost/community_support?sslmode=disable" up
   ```

2. **Environment Configuration**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the Application**:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Public Endpoints
- `GET /health` - Health check for Kemuko support service
- `GET /api/v1/public/categories` - List Kemuko support categories

### Student Endpoints (JWT Required)
- `GET /api/v1/tickets` - List student's tickets
- `POST /api/v1/tickets` - Create new support ticket (triggers Kemuko admin notifications)
- `GET /api/v1/tickets/{id}` - Get ticket details
- `POST /api/v1/tickets/{id}/comments` - Add comment to ticket

### Instructor/Admin Endpoints (Role-Based Access)
- `GET /api/v1/instructor/tickets` - List all tickets for instructor
- `PUT /api/v1/instructor/tickets/{id}` - Update ticket status/assignment
- `POST /api/v1/instructor/tickets/{id}/internal-notes` - Add internal notes

### Slack Integration Endpoints
- `POST /api/v1/slack/reply` - Admin reply via Slack (triggers email to student)
- `POST /api/v1/slack/webhook` - Slack bot webhook for interactive components

## Environment Variables

See `.env.example` for all available configuration options.

## Development

The project follows standard Go project layout and uses:
- **Gorilla Mux** for HTTP routing and middleware
- **sqlx** for database operations
- **PostgreSQL** with JSONB for flexible educational data storage
- **JWT** for authentication and user context extraction
- **UUID** for entity identifiers
- **Role-based access control** (student, instructor, admin)
- **External file service integration** (URLs only, no file uploads)
- **Microservice architecture** with external user management