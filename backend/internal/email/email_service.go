package email

import (
	"context"
	"time"
)

// EmailService defines the interface for email operations
type EmailService interface {
	// User deletion related emails
	SendDeletionRequestEmail(ctx context.Context, userID, email string, gracePeriodEnd time.Time) error
	SendGracePeriodReminderEmail(ctx context.Context, userID, email string, daysLeft int) error
	SendDeletionConfirmationEmail(ctx context.Context, userID, email string) error
	SendDataExportEmail(ctx context.Context, userID, email, exportURL string) error
	SendWelcomeBackEmail(ctx context.Context, userID, email string) error

	// Utility methods
	SendEmail(ctx context.Context, to, subject, htmlContent string) error
	Close() error
}

// EmailData contains data for email templates
type EmailData struct {
	UserID         string
	Email          string
	GracePeriodEnd time.Time
	DaysLeft       int
	ExportURL      string
	Subject        string
	TemplateName   string
}

// EmailTemplate defines email template structure
type EmailTemplate struct {
	Subject string
	HTML    string
	Text    string
}