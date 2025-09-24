package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationChannel string

const (
	NotificationTypeTicketCreated     NotificationType = "ticket_created"
	NotificationTypeTicketUpdated     NotificationType = "ticket_updated"
	NotificationTypeTicketAssigned    NotificationType = "ticket_assigned"
	NotificationTypeAdminReply        NotificationType = "admin_reply"
	NotificationTypeStudentReply      NotificationType = "student_reply"
)

const (
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSlack NotificationChannel = "slack"
)

type EmailNotificationRequest struct {
	To          []string               `jsonb:"to" validate:"required"`
	Subject     string                 `jsonb:"subject" validate:"required"`
	TemplateID  string                 `jsonb:"templateId" validate:"required"`
	TemplateData map[string]interface{} `jsonb:"templateData"`
	ReplyTo     *string                `jsonb:"replyTo"`
}

type SlackNotificationRequest struct {
	Channel     string                 `jsonb:"channel" validate:"required"`
	Message     string                 `jsonb:"message" validate:"required"`
	Blocks      []SlackBlock           `jsonb:"blocks,omitempty"`
	ThreadTS    *string                `jsonb:"threadTs,omitempty"`
	TemplateData map[string]interface{} `jsonb:"templateData"`
}

type SlackBlock struct {
	Type string      `jsonb:"type"`
	Text *SlackText  `jsonb:"text,omitempty"`
	Elements []SlackElement `jsonb:"elements,omitempty"`
}

type SlackText struct {
	Type string `jsonb:"type"`
	Text string `jsonb:"text"`
}

type SlackElement struct {
	Type     string `jsonb:"type"`
	Text     string `jsonb:"text,omitempty"`
	Value    string `jsonb:"value,omitempty"`
	ActionID string `jsonb:"actionId,omitempty"`
}

type SlackReplyRequest struct {
	TicketID    uuid.UUID `jsonb:"ticketId" validate:"required"`
	AdminUserID string    `jsonb:"adminUserId" validate:"required"`
	Message     string    `jsonb:"message" validate:"required"`
	IsInternal  bool      `jsonb:"isInternal"`
	Metadata    JSONB     `jsonb:"metadata"`
}

type NotificationEvent struct {
	Type         NotificationType `jsonb:"type"`
	TicketID     uuid.UUID        `jsonb:"ticketId"`
	TriggerUserID string          `jsonb:"triggerUserId"`
	Recipients   []string         `jsonb:"recipients"`
	Data         JSONB            `jsonb:"data"`
	CreatedAt    time.Time        `jsonb:"createdAt"`
}