package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"community-support-service/internal/database"
	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
)

type ticketRepository struct {
	db *database.DB
}

func NewTicketRepository(db *database.DB) repositories.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		) VALUES (
			:id, :ticketNumber, :title, :description, :status, :priority, :type,
			:studentId, :courseId, :instructorId, :categoryId, :metadata,
			:createdAt, :updatedAt
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, ticket)
	return err
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets 
		WHERE id = $1`
	
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets 
		WHERE ticketNumber = $1`
	
	err := r.db.GetContext(ctx, &ticket, query, ticketNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByStudentID(ctx context.Context, studentID string, filters repositories.TicketFilters) ([]*models.Ticket, error) {
	baseQuery := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets 
		WHERE studentId = $1`
	
	args := []interface{}{studentID}
	query, args := r.applyFilters(baseQuery, filters, args)
	
	var tickets []*models.Ticket
	err := r.db.SelectContext(ctx, &tickets, query, args...)
	return tickets, err
}

func (r *ticketRepository) GetByCourseID(ctx context.Context, courseID string, filters repositories.TicketFilters) ([]*models.Ticket, error) {
	baseQuery := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets 
		WHERE courseId = $1`
	
	args := []interface{}{courseID}
	query, args := r.applyFilters(baseQuery, filters, args)
	
	var tickets []*models.Ticket
	err := r.db.SelectContext(ctx, &tickets, query, args...)
	return tickets, err
}

func (r *ticketRepository) GetByInstructorID(ctx context.Context, instructorID string, filters repositories.TicketFilters) ([]*models.Ticket, error) {
	baseQuery := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets 
		WHERE instructorId = $1`
	
	args := []interface{}{instructorID}
	query, args := r.applyFilters(baseQuery, filters, args)
	
	var tickets []*models.Ticket
	err := r.db.SelectContext(ctx, &tickets, query, args...)
	return tickets, err
}

func (r *ticketRepository) List(ctx context.Context, filters repositories.TicketFilters, pagination repositories.Pagination) ([]*models.Ticket, int64, error) {
	baseQuery := `
		SELECT 
			id, ticketNumber, title, description, status, priority, type,
			studentId, courseId, instructorId, categoryId, metadata,
			createdAt, updatedAt
		FROM tickets`
	
	var args []interface{}
	query, args := r.applyFilters(baseQuery, filters, args)
	
	// Get total count
	countQuery := strings.Replace(query, "SELECT id, ticketNumber, title, description, status, priority, type, studentId, courseId, instructorId, categoryId, metadata, createdAt, updatedAt", "SELECT COUNT(*)", 1)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	orderBy := "createdAt"
	orderDir := "DESC"
	if pagination.OrderBy != "" {
		orderBy = pagination.OrderBy
	}
	if pagination.OrderDir != "" && (pagination.OrderDir == "ASC" || pagination.OrderDir == "DESC") {
		orderDir = pagination.OrderDir
	}
	
	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	
	if pagination.PageSize > 0 {
		offset := (pagination.Page - 1) * pagination.PageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.PageSize, offset)
	}
	
	var tickets []*models.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	return tickets, total, err
}

func (r *ticketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	query := `
		UPDATE tickets SET 
			title = :title,
			description = :description,
			status = :status,
			priority = :priority,
			type = :type,
			instructorId = :instructorId,
			categoryId = :categoryId,
			metadata = :metadata,
			updatedAt = :updatedAt
		WHERE id = :id`
	
	_, err := r.db.NamedExecContext(ctx, query, ticket)
	return err
}

func (r *ticketRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, updatedBy string) error {
	query := `
		UPDATE tickets SET 
			status = $1,
			updatedAt = $2
		WHERE id = $3`
	
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *ticketRepository) AssignInstructor(ctx context.Context, id uuid.UUID, instructorID string, updatedBy string) error {
	query := `
		UPDATE tickets SET 
			instructorId = $1,
			updatedAt = $2
		WHERE id = $3`
	
	_, err := r.db.ExecContext(ctx, query, instructorID, time.Now(), id)
	return err
}

func (r *ticketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tickets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ticketRepository) applyFilters(baseQuery string, filters repositories.TicketFilters, args []interface{}) (string, []interface{}) {
	var conditions []string
	argIndex := len(args) + 1
	
	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}
	
	if filters.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *filters.Priority)
		argIndex++
	}
	
	if filters.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *filters.Type)
		argIndex++
	}
	
	if filters.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("categoryId = $%d", argIndex))
		args = append(args, *filters.CategoryID)
		argIndex++
	}
	
	if filters.CourseID != nil {
		conditions = append(conditions, fmt.Sprintf("courseId = $%d", argIndex))
		args = append(args, *filters.CourseID)
		argIndex++
	}
	
	if filters.InstructorID != nil {
		conditions = append(conditions, fmt.Sprintf("instructorId = $%d", argIndex))
		args = append(args, *filters.InstructorID)
		argIndex++
	}
	
	if filters.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("createdAt >= $%d", argIndex))
		args = append(args, *filters.FromDate)
		argIndex++
	}
	
	if filters.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("createdAt <= $%d", argIndex))
		args = append(args, *filters.ToDate)
		argIndex++
	}
	
	if filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + *filters.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR ticketNumber ILIKE $%d)", argIndex, argIndex+1, argIndex+2))
		args = append(args, searchTerm, searchTerm, searchTerm)
		argIndex += 3
	}
	
	if len(conditions) > 0 {
		connector := " WHERE "
		if strings.Contains(baseQuery, "WHERE") {
			connector = " AND "
		}
		baseQuery += connector + strings.Join(conditions, " AND ")
	}
	
	return baseQuery, args
}