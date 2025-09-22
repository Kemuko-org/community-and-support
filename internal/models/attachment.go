package models

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID         uuid.UUID  `jsonb:"id" db:"id"`
	TicketID   *uuid.UUID `jsonb:"ticketId" db:"ticketId"`
	CommentID  *uuid.UUID `jsonb:"commentId" db:"commentId"`
	FileName   string     `jsonb:"fileName" db:"fileName"`
	FileUrl    string     `jsonb:"fileUrl" db:"fileUrl"`
	FileType   *string    `jsonb:"fileType" db:"fileType"`
	UploadedBy string     `jsonb:"uploadedBy" db:"uploadedBy"`
	Metadata   JSONB      `jsonb:"metadata" db:"metadata"`
	CreatedAt  time.Time  `jsonb:"createdAt" db:"createdAt"`
}

type CreateAttachmentRequest struct {
	FileName string  `jsonb:"fileName" validate:"required"`
	FileUrl  string  `jsonb:"fileUrl" validate:"required,url"`
	FileType *string `jsonb:"fileType"`
	Metadata JSONB   `jsonb:"metadata"`
}