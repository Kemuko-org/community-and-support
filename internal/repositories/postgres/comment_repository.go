package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"community-support-service/internal/database"
	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
)

type commentRepository struct {
	db *database.DB
}

func NewCommentRepository(db *database.DB) repositories.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.TicketComment) error {
	query := `
		INSERT INTO ticketComments (
			id, ticketId, commentText, isInternal, authorId, authorType,
			metadata, createdAt, updatedAt
		) VALUES (
			:id, :ticketId, :commentText, :isInternal, :authorId, :authorType,
			:metadata, :createdAt, :updatedAt
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, comment)
	return err
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketComment, error) {
	var comment models.TicketComment
	query := `
		SELECT 
			id, ticketId, commentText, isInternal, authorId, authorType,
			metadata, createdAt, updatedAt
		FROM ticketComments 
		WHERE id = $1`
	
	err := r.db.GetContext(ctx, &comment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) GetByTicketID(ctx context.Context, ticketID uuid.UUID, includeInternal bool) ([]*models.TicketComment, error) {
	query := `
		SELECT 
			id, ticketId, commentText, isInternal, authorId, authorType,
			metadata, createdAt, updatedAt
		FROM ticketComments 
		WHERE ticketId = $1`
	
	args := []interface{}{ticketID}
	
	if !includeInternal {
		query += " AND isInternal = false"
	}
	
	query += " ORDER BY createdAt ASC"
	
	var comments []*models.TicketComment
	err := r.db.SelectContext(ctx, &comments, query, args...)
	return comments, err
}

func (r *commentRepository) GetCommentCount(ctx context.Context, ticketID uuid.UUID, includeInternal bool) (int64, error) {
	query := `SELECT COUNT(*) FROM ticketComments WHERE ticketId = $1`
	args := []interface{}{ticketID}
	
	if !includeInternal {
		query += " AND isInternal = false"
	}
	
	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *commentRepository) Update(ctx context.Context, comment *models.TicketComment) error {
	query := `
		UPDATE ticketComments SET 
			commentText = :commentText,
			isInternal = :isInternal,
			metadata = :metadata,
			updatedAt = :updatedAt
		WHERE id = :id`
	
	_, err := r.db.NamedExecContext(ctx, query, comment)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ticketComments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}