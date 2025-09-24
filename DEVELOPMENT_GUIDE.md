# EdTech Support Service - Development Guide

## 📋 Project Overview

This is a Go-based microservice for educational technology support and ticketing system. It's designed to work within a larger EdTech ecosystem with separate user management and file upload services.

### Key Characteristics
- **Microservice Architecture**: No user management, uses external authentication
- **EdTech Context**: Student/instructor support tickets with course integration
- **JWT Authentication**: Token-based auth with role-based access control
- **File URLs Only**: No file uploads, only references to external file service
- **PostgreSQL + JSONB**: Flexible data storage with educational metadata

## 🏗️ Current Project State

### ✅ Completed Components

1. **Database Schema** (camelCase + JSONB)
   - `categories` - Support categories
   - `tickets` - Student tickets with course integration
   - `ticketComments` - Student-instructor communication
   - `attachments` - File URL references
   - `ticketHistory` - Audit trail

2. **Go Models** (with JSONB support)
   - Category, Ticket, TicketComment, Attachment models
   - EdTech-specific types (TicketType, Priority, Status)
   - Request/Response DTOs

3. **Authentication & Middleware**
   - JWT token extraction from headers
   - User context management
   - Role-based access control (student, instructor, admin)

4. **Basic Server Setup**
   - Configuration management
   - Database connection
   - Route structure with authentication

### 🚧 Still Needed

1. **Database Layer (Repositories)**
   - CRUD operations for all models
   - Query builders with filtering
   - Transaction support

2. **Business Logic (Services)**
   - Ticket management logic
   - Comment handling
   - File attachment logic
   - History tracking

3. **HTTP Handlers**
   - REST API endpoints
   - Request validation
   - Response formatting

4. **Advanced Features**
   - Pagination
   - Search and filtering
   - Email notifications
   - Metrics and logging

## 📁 Project Structure

```
community-and-support/
├── cmd/
│   └── server/main.go              # Application entry point
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── database/db.go              # Database connection
│   ├── models/                     # Data models
│   │   ├── category.go
│   │   ├── ticket.go
│   │   └── attachment.go
│   ├── handlers/                   # HTTP handlers (TODO)
│   ├── services/                   # Business logic (TODO)
│   └── repositories/               # Database layer (TODO)
├── migrations/                     # Database migrations
│   ├── 001_create_categories_table.up/down.sql
│   ├── 002_create_tickets_table.up/down.sql
│   ├── 003_create_ticket_comments_table.up/down.sql
│   ├── 004_create_attachments_table.up/down.sql
│   └── 005_create_ticket_history_table.up/down.sql
├── pkg/
│   ├── middleware/auth.go          # JWT authentication
│   └── utils/response.go           # API response helpers
├── .env.example                    # Environment variables template
└── go.mod                          # Go dependencies
```

## 🔧 Environment Setup

### Prerequisites
- Go 1.22.2+
- PostgreSQL 12+
- Migration tool (golang-migrate recommended)

### Environment Variables
```bash
# Copy and configure
cp .env.example .env

# Required variables:
SERVER_HOST=localhost
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=community_support
JWT_SECRET=your-secret-key
```

### Database Setup
```bash
# Create database
createdb community_support

# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://user:password@localhost/community_support?sslmode=disable" up
```

### Run Application
```bash
go mod tidy
go run cmd/server/main.go
```

## 🎯 Next Development Steps

### Phase 1: Core Database Layer (Priority: HIGH)

1. **Create Repository Interfaces**
   ```go
   // internal/repositories/interfaces.go
   type TicketRepository interface {
       Create(ctx context.Context, ticket *models.Ticket) error
       GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
       GetByStudentID(ctx context.Context, studentID string, filters TicketFilters) ([]*models.Ticket, error)
       Update(ctx context.Context, ticket *models.Ticket) error
       // ... more methods
   }
   ```

2. **Implement PostgreSQL Repositories**
   ```go
   // internal/repositories/postgres/ticket_repository.go
   type ticketRepository struct {
       db *database.DB
   }
   ```

3. **Add Repository Tests**
   - Unit tests for each repository method
   - Database integration tests

### Phase 2: Business Logic Services (Priority: HIGH)

1. **Create Service Layer**
   ```go
   // internal/services/ticket_service.go
   type TicketService struct {
       ticketRepo    repositories.TicketRepository
       commentRepo   repositories.CommentRepository
       historyRepo   repositories.HistoryRepository
   }
   ```

2. **Implement Core Features**
   - Ticket creation with automatic numbering
   - Status transitions with validation
   - Comment handling with visibility rules
   - History tracking for all changes

### Phase 3: HTTP API Layer (Priority: HIGH)

1. **Create HTTP Handlers**
   ```go
   // internal/handlers/ticket_handler.go
   type TicketHandler struct {
       ticketService *services.TicketService
   }
   ```

2. **Implement API Endpoints**
   - Student endpoints (CRUD tickets, comments)
   - Instructor endpoints (assignment, status updates)
   - Admin endpoints (full management)

### Phase 4: Advanced Features (Priority: MEDIUM)

1. **Search and Filtering**
   - Full-text search on tickets
   - Advanced filtering (status, priority, course, etc.)
   - Sorting options

2. **Pagination and Performance**
   - Cursor-based pagination
   - Database query optimization
   - Caching layer (Redis)

### Phase 5: Production Features (Priority: LOW)

1. **Observability**
   - Structured logging
   - Metrics collection
   - Health checks
   - Tracing

2. **Integration Features**
   - Webhook notifications
   - Email notifications
   - External service integrations

## 💡 Implementation Guidelines

### Database Patterns
```go
// Always use transactions for multi-table operations
err := repo.db.WithTx(func(tx *sqlx.Tx) error {
    // Create ticket
    if err := createTicket(tx, ticket); err != nil {
        return err
    }
    // Create history entry
    if err := createHistory(tx, history); err != nil {
        return err
    }
    return nil
})
```

### Authentication Patterns
```go
// Extract user from context in handlers
user, ok := middleware.GetUserFromContext(r.Context())
if !ok {
    utils.WriteError(w, http.StatusUnauthorized, "User not found")
    return
}

// Use user ID for data filtering
tickets, err := h.ticketService.GetByStudentID(ctx, user.UserID, filters)
```

### Error Handling
```go
// Use structured errors
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// Return consistent API responses
utils.WriteError(w, http.StatusBadRequest, "Invalid request data")
utils.WriteSuccess(w, "Ticket created successfully", ticket)
```

### JSONB Metadata Usage
```go
// Store flexible educational data
metadata := models.JSONB{
    "courseTitle": "Introduction to Computer Science",
    "assignmentName": "Final Project",
    "dueDate": "2024-12-15",
    "moduleNumber": 3,
}
ticket.Metadata = metadata
```

## 🧪 Testing Strategy

### Unit Tests
- Repository layer with database mocks
- Service layer with repository mocks
- Handler layer with service mocks

### Integration Tests
- Database operations with test database
- API endpoints with test server
- Authentication flow testing

### Test Structure
```go
func TestTicketService_CreateTicket(t *testing.T) {
    // Setup
    mockRepo := &mocks.TicketRepository{}
    service := services.NewTicketService(mockRepo)
    
    // Test cases
    tests := []struct {
        name    string
        input   *models.CreateTicketRequest
        wantErr bool
    }{
        // ... test cases
    }
    
    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

## 🚀 Deployment Considerations

### Docker Setup
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment-Specific Configs
- Development: Local PostgreSQL, debug logging
- Staging: Managed database, structured logging
- Production: High availability setup, monitoring

## 📚 Key Resources

### Documentation
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [sqlx](https://github.com/jmoiron/sqlx) - SQL extensions
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [jwt-go](https://github.com/golang-jwt/jwt) - JWT handling

### EdTech Domain Knowledge
- Student support ticket types and workflows
- Instructor assignment and escalation patterns
- Course integration and metadata requirements
- Educational data privacy considerations

## 🔄 Git Workflow

### Branch Strategy
```bash
# Feature development
git checkout -b feature/ticket-repository
git commit -m "feat: implement ticket repository with CRUD operations"

# Bug fixes
git checkout -b fix/authentication-middleware
git commit -m "fix: handle missing JWT claims gracefully"
```

### Commit Message Format
```
type(scope): description

feat: new feature
fix: bug fix
docs: documentation
style: formatting
refactor: code restructuring
test: adding tests
chore: maintenance
```

## 📝 Progress Tracking (MANDATORY)

### Critical Pattern
**After completing ANY work, you MUST update `.claude/progress.md`**

This is the permanent record that persists across all Claude sessions and enables seamless project continuation.

### Required Progress Entry Format
```markdown
## YYYY-MM-DD HH:MM - [Feature/Component Name]

### Completed:
- Specific files created/modified
- Functions/features implemented
- Tests added

### Key Decisions:
- Technical choices made
- Patterns established
- Trade-offs considered

### Next Steps:
- Immediate follow-up tasks
- Dependencies to address
- Areas requiring attention

### Files Modified:
- path/to/file1.go (new/modified)
- path/to/file2.go (new/modified)
```

### Example Progress Update
```bash
# After implementing ticket repository
echo "## $(date '+%Y-%m-%d %H:%M') - Ticket Repository Implementation

### Completed:
- Created TicketRepository interface with full CRUD operations
- Implemented PostgreSQL ticket repository with transaction support
- Added filtering capabilities for student/instructor views
- Implemented JSONB metadata handling

### Key Decisions:
- Used sqlx named parameters for SQL injection prevention
- Implemented dynamic WHERE clause building for flexible filtering
- Added proper error wrapping with context

### Next Steps:
- Implement CommentRepository interface and PostgreSQL implementation
- Add comprehensive repository tests
- Create service layer for business logic

### Files Modified:
- internal/repositories/interfaces.go (new)
- internal/repositories/postgres/ticket_repository.go (new)
- internal/repositories/postgres/ticket_repository_test.go (new)
" >> .claude/progress.md
```

This guide provides everything needed to continue development efficiently. The project is well-structured and ready for the next phase of implementation.