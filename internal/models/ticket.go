package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string
type TicketPriority string
type TicketType string

const (
	TicketStatusOpen              TicketStatus = "open"
	TicketStatusInProgress        TicketStatus = "inProgress"
	TicketStatusWaitingForCustomer TicketStatus = "waitingForCustomer"
	TicketStatusResolved          TicketStatus = "resolved"
	TicketStatusClosed            TicketStatus = "closed"
)

const (
	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityUrgent TicketPriority = "urgent"
)

const (
	TicketTypeGeneral    TicketType = "general"
	TicketTypeTechnical  TicketType = "technical"
	TicketTypeCourse     TicketType = "course"
	TicketTypeAssignment TicketType = "assignment"
	TicketTypeGrading    TicketType = "grading"
	TicketTypePlatform   TicketType = "platform"
	TicketTypeContent    TicketType = "content"
)

type Ticket struct {
	ID               uuid.UUID       `jsonb:"id" db:"id"`
	TicketNumber     string          `jsonb:"ticketNumber" db:"ticketNumber"`
	Title            string          `jsonb:"title" db:"title"`
	Description      string          `jsonb:"description" db:"description"`
	Status           TicketStatus    `jsonb:"status" db:"status"`
	Priority         TicketPriority  `jsonb:"priority" db:"priority"`
	Type             TicketType      `jsonb:"type" db:"type"`
	StudentID        string          `jsonb:"studentId" db:"studentId"`
	InstructorID     *string         `jsonb:"instructorId" db:"instructorId"`
	CourseID         *string         `jsonb:"courseId" db:"courseId"`
	CategoryID       *uuid.UUID      `jsonb:"categoryId" db:"categoryId"`
	Metadata         JSONB           `jsonb:"metadata" db:"metadata"`
	CreatedAt        time.Time       `jsonb:"createdAt" db:"createdAt"`
	UpdatedAt        time.Time       `jsonb:"updatedAt" db:"updatedAt"`
	ResolvedAt       *time.Time      `jsonb:"resolvedAt" db:"resolvedAt"`
	ClosedAt         *time.Time      `jsonb:"closedAt" db:"closedAt"`
	
	// Related entities (populated via joins)
	Category       *Category `jsonb:"category,omitempty"`
}

type CreateTicketRequest struct {
	Title       string         `jsonb:"title" validate:"required,max=255"`
	Description string         `jsonb:"description" validate:"required"`
	Priority    TicketPriority `jsonb:"priority" validate:"required"`
	Type        TicketType     `jsonb:"type" validate:"required"`
	CourseID    *string        `jsonb:"courseId"`
	CategoryID  *uuid.UUID     `jsonb:"categoryId"`
	Metadata    JSONB          `jsonb:"metadata"`
}

type UpdateTicketRequest struct {
	Title           *string         `jsonb:"title" validate:"omitempty,max=255"`
	Description     *string         `jsonb:"description"`
	Status          *TicketStatus   `jsonb:"status"`
	Priority        *TicketPriority `jsonb:"priority"`
	Type            *TicketType     `jsonb:"type"`
	InstructorID    *string         `jsonb:"instructorId"`
	CourseID        *string         `jsonb:"courseId"`
	CategoryID      *uuid.UUID      `jsonb:"categoryId"`
	Metadata        JSONB           `jsonb:"metadata"`
}

type TicketComment struct {
	ID         uuid.UUID `jsonb:"id" db:"id"`
	TicketID   uuid.UUID `jsonb:"ticketId" db:"ticketId"`
	UserID     string    `jsonb:"userId" db:"userId"`
	Content    string    `jsonb:"content" db:"content"`
	IsInternal bool      `jsonb:"isInternal" db:"isInternal"`
	Metadata   JSONB     `jsonb:"metadata" db:"metadata"`
	CreatedAt  time.Time `jsonb:"createdAt" db:"createdAt"`
	UpdatedAt  time.Time `jsonb:"updatedAt" db:"updatedAt"`
}

type CreateCommentRequest struct {
	Content    string `jsonb:"content" validate:"required"`
	IsInternal bool   `jsonb:"isInternal"`
	Metadata   JSONB  `jsonb:"metadata"`
}

type TicketHistory struct {
	ID          uuid.UUID `jsonb:"id" db:"id"`
	TicketID    uuid.UUID `jsonb:"ticketId" db:"ticketId"`
	UserID      string    `jsonb:"userId" db:"userId"`
	Action      string    `jsonb:"action" db:"action"`
	OldValue    *string   `jsonb:"oldValue" db:"oldValue"`
	NewValue    *string   `jsonb:"newValue" db:"newValue"`
	Description *string   `jsonb:"description" db:"description"`
	Metadata    JSONB     `jsonb:"metadata" db:"metadata"`
	CreatedAt   time.Time `jsonb:"createdAt" db:"createdAt"`
}