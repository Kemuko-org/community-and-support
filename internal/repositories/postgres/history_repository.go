package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"community-support-service/internal/database"
	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
)

type historyRepository struct {
	db *database.DB
}

func NewHistoryRepository(db *database.DB) repositories.HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Create(ctx context.Context, history *models.TicketHistory) error {
	query := `
		INSERT INTO ticketHistory (
			id, ticketId, actionType, description, changedBy, changedByType,
			oldValue, newValue, metadata, createdAt
		) VALUES (
			:id, :ticketId, :actionType, :description, :changedBy, :changedByType,
			:oldValue, :newValue, :metadata, :createdAt
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, history)
	return err
}

func (r *historyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketHistory, error) {
	var history models.TicketHistory
	query := `
		SELECT 
			id, ticketId, actionType, description, changedBy, changedByType,
			oldValue, newValue, metadata, createdAt
		FROM ticketHistory 
		WHERE id = $1`
	
	err := r.db.GetContext(ctx, &history, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

func (r *historyRepository) GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*models.TicketHistory, error) {
	query := `
		SELECT 
			id, ticketId, actionType, description, changedBy, changedByType,
			oldValue, newValue, metadata, createdAt
		FROM ticketHistory 
		WHERE ticketId = $1
		ORDER BY createdAt ASC`
	
	var history []*models.TicketHistory
	err := r.db.SelectContext(ctx, &history, query, ticketID)
	return history, err
}

func (r *historyRepository) GetChangeHistory(ctx context.Context, ticketID uuid.UUID, actionType *string) ([]*models.TicketHistory, error) {
	query := `
		SELECT 
			id, ticketId, actionType, description, changedBy, changedByType,
			oldValue, newValue, metadata, createdAt
		FROM ticketHistory 
		WHERE ticketId = $1`
	
	args := []interface{}{ticketID}
	
	if actionType != nil {
		query += " AND actionType = $2"
		args = append(args, *actionType)
	}
	
	query += " ORDER BY createdAt ASC"
	
	var history []*models.TicketHistory
	err := r.db.SelectContext(ctx, &history, query, args...)
	return history, err
}

func (r *historyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ticketHistory WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}