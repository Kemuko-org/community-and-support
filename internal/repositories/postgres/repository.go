package postgres

import (
	"community-support-service/internal/database"
	"community-support-service/internal/repositories"
)

func NewRepository(db *database.DB) *repositories.Repository {
	return &repositories.Repository{
		Ticket:     NewTicketRepository(db),
		Comment:    NewCommentRepository(db),
		Category:   NewCategoryRepository(db),
		Attachment: NewAttachmentRepository(db),
		History:    NewHistoryRepository(db),
	}
}