package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"community-support-service/internal/models"
)

type TicketFilters struct {
	Status       *string
	Priority     *string
	Type         *string
	CategoryID   *uuid.UUID
	CourseID     *string
	InstructorID *string
	Search       *string
	FromDate     *time.Time
	ToDate       *time.Time
}

type Pagination struct {
	Page     int
	PageSize int
	OrderBy  string
	OrderDir string // ASC or DESC
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *models.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	GetByStudentID(ctx context.Context, studentID string, filters TicketFilters) ([]*models.Ticket, error)
	GetByCourseID(ctx context.Context, courseID string, filters TicketFilters) ([]*models.Ticket, error)
	GetByInstructorID(ctx context.Context, instructorID string, filters TicketFilters) ([]*models.Ticket, error)
	Update(ctx context.Context, ticket *models.Ticket) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters TicketFilters, pagination Pagination) ([]*models.Ticket, int64, error)
	GetByTicketNumber(ctx context.Context, ticketNumber string) (*models.Ticket, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, updatedBy string) error
	AssignInstructor(ctx context.Context, id uuid.UUID, instructorID string, updatedBy string) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.TicketComment) error
	GetByTicketID(ctx context.Context, ticketID uuid.UUID, includeInternal bool) ([]*models.TicketComment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TicketComment, error)
	Update(ctx context.Context, comment *models.TicketComment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetCommentCount(ctx context.Context, ticketID uuid.UUID, includeInternal bool) (int64, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetAll(ctx context.Context, activeOnly bool) ([]*models.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetByType(ctx context.Context, categoryType string, activeOnly bool) ([]*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error
}

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*models.Attachment, error)
	GetByCommentID(ctx context.Context, commentID uuid.UUID) ([]*models.Attachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByTicketID(ctx context.Context, ticketID uuid.UUID) error
	DeleteByCommentID(ctx context.Context, commentID uuid.UUID) error
}

type HistoryRepository interface {
	Create(ctx context.Context, history *models.TicketHistory) error
	GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*models.TicketHistory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TicketHistory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetChangeHistory(ctx context.Context, ticketID uuid.UUID, actionType *string) ([]*models.TicketHistory, error)
}

type Repository struct {
	Ticket     TicketRepository
	Comment    CommentRepository
	Category   CategoryRepository
	Attachment AttachmentRepository
	History    HistoryRepository
}