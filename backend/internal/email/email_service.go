package email

import (
	"context"
	"time"
)

// EmailService defines the interface for email operations
type EmailService interface {
	// User deletion related emails
	SendDeletionRequestEmail(ctx context.Context, userID, email string, userName string, gracePeriodEnd time.Time) error
	SendGracePeriodReminderEmail(ctx context.Context, userID, email string, userName string, daysLeft int, deletionDate time.Time) error
	SendDeletionConfirmationEmail(ctx context.Context, userID, email string, userName string, deletionDate time.Time) error
	SendDataExportEmail(ctx context.Context, userID, email, userName string, exportURL string) error
	SendWelcomeBackEmail(ctx context.Context, userID, email string, userName string) error

	// Utility methods
	SendEmail(ctx context.Context, to, subject, htmlContent string) error
	Close() error
}

// EmailData contains data for email templates
type EmailData struct {
	UserID         string
	Email          string
	Name           string
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