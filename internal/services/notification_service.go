package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"community-support-service/internal/config"
	"community-support-service/internal/models"
)

type NotificationService struct {
	config     *config.Config
	httpClient *http.Client
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *NotificationService) SendTicketCreatedNotifications(ctx context.Context, ticket *models.Ticket, studentEmail string) error {
	// Send email to Kemuko admins
	if err := s.sendAdminEmailNotification(ctx, ticket, studentEmail); err != nil {
		return fmt.Errorf("failed to send admin email notification: %w", err)
	}

	// Send email to ticket owner (student)
	if err := s.sendStudentEmailNotification(ctx, ticket, studentEmail); err != nil {
		return fmt.Errorf("failed to send student email notification: %w", err)
	}

	// Send Slack notification to admin channel
	if err := s.sendSlackNotification(ctx, ticket, studentEmail); err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	return nil
}

func (s *NotificationService) SendAdminReplyNotification(ctx context.Context, ticket *models.Ticket, comment *models.TicketComment, studentEmail string) error {
	// Send email to student about admin reply
	emailReq := &models.EmailNotificationRequest{
		To:         []string{studentEmail},
		Subject:    fmt.Sprintf("Kemuko Support - Reply to Ticket #%s", ticket.TicketNumber),
		TemplateID: "admin_reply_notification",
		TemplateData: map[string]interface{}{
			"studentName":    "Student", // You might want to get this from user service
			"ticketNumber":   ticket.TicketNumber,
			"ticketTitle":    ticket.Title,
			"replyMessage":   comment.Content,
			"ticketUrl":      fmt.Sprintf("%s/tickets/%s", s.config.Frontend.BaseURL, ticket.ID),
			"platformName":   "Kemuko",
		},
		ReplyTo: &s.config.Notifications.AdminEmail,
	}

	return s.sendEmailRequest(ctx, emailReq)
}

func (s *NotificationService) sendAdminEmailNotification(ctx context.Context, ticket *models.Ticket, studentEmail string) error {
	adminEmails := s.config.Notifications.AdminEmails
	
	emailReq := &models.EmailNotificationRequest{
		To:         adminEmails,
		Subject:    fmt.Sprintf("Kemuko Support - New Ticket #%s", ticket.TicketNumber),
		TemplateID: "new_ticket_admin_notification",
		TemplateData: map[string]interface{}{
			"ticketNumber":   ticket.TicketNumber,
			"ticketTitle":    ticket.Title,
			"ticketType":     string(ticket.Type),
			"priority":       string(ticket.Priority),
			"studentEmail":   studentEmail,
			"ticketUrl":      fmt.Sprintf("%s/admin/tickets/%s", s.config.Frontend.BaseURL, ticket.ID),
			"courseId":       ticket.CourseID,
			"platformName":   "Kemuko",
			"createdAt":      ticket.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		ReplyTo: &studentEmail,
	}

	return s.sendEmailRequest(ctx, emailReq)
}

func (s *NotificationService) sendStudentEmailNotification(ctx context.Context, ticket *models.Ticket, studentEmail string) error {
	emailReq := &models.EmailNotificationRequest{
		To:         []string{studentEmail},
		Subject:    fmt.Sprintf("Kemuko Support - Ticket Created #%s", ticket.TicketNumber),
		TemplateID: "ticket_created_confirmation",
		TemplateData: map[string]interface{}{
			"ticketNumber":   ticket.TicketNumber,
			"ticketTitle":    ticket.Title,
			"ticketType":     string(ticket.Type),
			"priority":       string(ticket.Priority),
			"ticketUrl":      fmt.Sprintf("%s/tickets/%s", s.config.Frontend.BaseURL, ticket.ID),
			"platformName":   "Kemuko",
			"supportEmail":   s.config.Notifications.AdminEmail,
		},
	}

	return s.sendEmailRequest(ctx, emailReq)
}

func (s *NotificationService) sendSlackNotification(ctx context.Context, ticket *models.Ticket, studentEmail string) error {
	slackReq := &models.SlackNotificationRequest{
		Channel: s.config.Notifications.SlackChannel,
		Message: fmt.Sprintf("🎓 New Kemuko Support Ticket Created"),
		Blocks: []models.SlackBlock{
			{
				Type: "section",
				Text: &models.SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*🎓 New Kemuko Support Ticket*\n*Ticket:* #%s\n*Title:* %s\n*Type:* %s\n*Priority:* %s\n*Student:* %s", 
						ticket.TicketNumber, ticket.Title, ticket.Type, ticket.Priority, studentEmail),
				},
			},
			{
				Type: "section",
				Elements: []models.SlackElement{
					{
						Type:     "button",
						Text:     "View Ticket",
						Value:    ticket.ID.String(),
						ActionID: "view_ticket",
					},
					{
						Type:     "button",
						Text:     "Reply",
						Value:    ticket.ID.String(),
						ActionID: "reply_ticket",
					},
				},
			},
		},
		TemplateData: map[string]interface{}{
			"ticketId":     ticket.ID.String(),
			"ticketNumber": ticket.TicketNumber,
		},
	}

	return s.sendSlackRequest(ctx, slackReq)
}

func (s *NotificationService) sendEmailRequest(ctx context.Context, req *models.EmailNotificationRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.Notifications.EmailServiceURL+"/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.Notifications.EmailServiceToken)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("email service returned status: %d", resp.StatusCode)
	}

	return nil
}

func (s *NotificationService) SendSlackRequest(ctx context.Context, req *models.SlackNotificationRequest) error {
	return s.sendSlackRequest(ctx, req)
}

func (s *NotificationService) sendSlackRequest(ctx context.Context, req *models.SlackNotificationRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal slack request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.Notifications.SlackServiceURL+"/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create slack request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.Notifications.SlackServiceToken)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send slack request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("slack service returned status: %d", resp.StatusCode)
	}

	return nil
}