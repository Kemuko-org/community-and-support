package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"community-support-service/internal/database"
	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
)

type attachmentRepository struct {
	db *database.DB
}

func NewAttachmentRepository(db *database.DB) repositories.AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	query := `
		INSERT INTO attachments (
			id, ticketId, commentId, fileName, fileUrl, fileType, fileSize,
			metadata, uploadedBy, createdAt
		) VALUES (
			:id, :ticketId, :commentId, :fileName, :fileUrl, :fileType, :fileSize,
			:metadata, :uploadedBy, :createdAt
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, attachment)
	return err
}

func (r *attachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Attachment, error) {
	var attachment models.Attachment
	query := `
		SELECT 
			id, ticketId, commentId, fileName, fileUrl, fileType, fileSize,
			metadata, uploadedBy, createdAt
		FROM attachments 
		WHERE id = $1`
	
	err := r.db.GetContext(ctx, &attachment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attachment, nil
}

func (r *attachmentRepository) GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*models.Attachment, error) {
	query := `
		SELECT 
			id, ticketId, commentId, fileName, fileUrl, fileType, fileSize,
			metadata, uploadedBy, createdAt
		FROM attachments 
		WHERE ticketId = $1
		ORDER BY createdAt ASC`
	
	var attachments []*models.Attachment
	err := r.db.SelectContext(ctx, &attachments, query, ticketID)
	return attachments, err
}

func (r *attachmentRepository) GetByCommentID(ctx context.Context, commentID uuid.UUID) ([]*models.Attachment, error) {
	query := `
		SELECT 
			id, ticketId, commentId, fileName, fileUrl, fileType, fileSize,
			metadata, uploadedBy, createdAt
		FROM attachments 
		WHERE commentId = $1
		ORDER BY createdAt ASC`
	
	var attachments []*models.Attachment
	err := r.db.SelectContext(ctx, &attachments, query, commentID)
	return attachments, err
}

func (r *attachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM attachments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *attachmentRepository) DeleteByTicketID(ctx context.Context, ticketID uuid.UUID) error {
	query := `DELETE FROM attachments WHERE ticketId = $1`
	_, err := r.db.ExecContext(ctx, query, ticketID)
	return err
}

func (r *attachmentRepository) DeleteByCommentID(ctx context.Context, commentID uuid.UUID) error {
	query := `DELETE FROM attachments WHERE commentId = $1`
	_, err := r.db.ExecContext(ctx, query, commentID)
	return err
}