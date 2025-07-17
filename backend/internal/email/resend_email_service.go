package email

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/resendlabs/resend-go"
)

// ResendEmailService implements EmailService using Resend API
type ResendEmailService struct {
	appCtx        *appcontext.AppContext
	client        *resend.Client
	templateEngine *TemplateEngine
	config        *config.EmailConfig
}

// NewResendEmailService creates a new Resend email service
func NewResendEmailService(appCtx *appcontext.AppContext) (EmailService, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context cannot be nil")
	}

	if appCtx.Config.Email == nil {
		return nil, fmt.Errorf("email configuration is required")
	}

	if appCtx.Config.Email.ResendAPIKey == "" {
		return nil, fmt.Errorf("resend API key is required")
	}

	// Create Resend client
	client := resend.NewClient(appCtx.Config.Email.ResendAPIKey)

	// Create template engine
	templateEngine, err := NewTemplateEngine(appCtx.Config.Email.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create template engine: %w", err)
	}

	return &ResendEmailService{
		appCtx:        appCtx,
		client:        client,
		templateEngine: templateEngine,
		config:        appCtx.Config.Email,
	}, nil
}

// SendEmail sends a generic email using Resend
func (res *ResendEmailService) SendEmail(
	ctx context.Context,
	to,
	subject,
	htmlContent string,
) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", res.config.FromName, res.config.FromAddress),
		To:      []string{to},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := res.client.Emails.Send(params)
	if err != nil {
		res.appCtx.Logger.Error("Failed to send email via Resend", map[string]any{
			"to":      to,
			"subject": subject,
			"error":   err.Error(),
		})
		return fmt.Errorf("failed to send email: %w", err)
	}

	res.appCtx.Logger.Info("Email sent successfully", map[string]any{
		"to":      to,
		"subject": subject,
	})

	return nil
}

// SendDeletionRequestEmail sends email notification when deletion is requested
func (res *ResendEmailService) SendDeletionRequestEmail(
	ctx context.Context,
	userID,
	email,
	userName string,
	gracePeriodEnd time.Time,
) error {
	// Render email template
	htmlContent, err := res.templateEngine.RenderDeletionRequest(
		userID,
		email,
		userName,
		gracePeriodEnd,
	)
	if err != nil {
		return fmt.Errorf("failed to render deletion request template: %w", err)
	}

	subject := "Account Deletion Requested - QKO"

	return res.SendEmail(ctx, email, subject, htmlContent)
}

// SendGracePeriodReminderEmail sends reminder email before grace period ends
func (res *ResendEmailService) SendGracePeriodReminderEmail(
	ctx context.Context,
	userID,
	email,
	userName string,
	daysLeft int,
	deletionDate time.Time,
) error {
	// Render email template
	htmlContent, err := res.templateEngine.RenderGracePeriodReminder(
		userID,
		email,
		userName,
		daysLeft,
		deletionDate,
	)
	if err != nil {
		return fmt.Errorf("failed to render grace period reminder template: %w", err)
	}

	subject := fmt.Sprintf("Final Reminder: Account Deletion in %d days - QKO", daysLeft)

	return res.SendEmail(ctx, email, subject, htmlContent)
}

// SendDeletionConfirmationEmail sends confirmation email after permanent deletion
func (res *ResendEmailService) SendDeletionConfirmationEmail(
	ctx context.Context,
	userID,
	email,
	userName string,
	deletionDate time.Time,
) error {
	// Render email template
	htmlContent, err := res.templateEngine.RenderDeletionConfirmation(
		userID,
		email,
		userName,
		deletionDate,
	)
	if err != nil {
		return fmt.Errorf("failed to render deletion confirmation template: %w", err)
	}

	subject := "Account Deletion Confirmed - QKO"

	return res.SendEmail(ctx, email, subject, htmlContent)
}

// SendDataExportEmail sends email with data export download link
func (res *ResendEmailService) SendDataExportEmail(
	ctx context.Context,
	userID,
	email,
	userName,
	exportURL string,
) error {
	// Render email template
	htmlContent, err := res.templateEngine.RenderDataExport(
		userID,
		email,
		userName,
		exportURL,
	)
	if err != nil {
		return fmt.Errorf("failed to render data export template: %w", err)
	}

	subject := "Your Data Export is Ready - QKO"

	return res.SendEmail(ctx, email, subject, htmlContent)
}

// SendWelcomeBackEmail sends welcome back email when deletion is cancelled
func (res *ResendEmailService) SendWelcomeBackEmail(
	ctx context.Context,
	userID,
	email,
	userName string,
) error {
	// Render email template
	htmlContent, err := res.templateEngine.RenderWelcomeBack(
		userID,
		email,
		userName,
	)
	if err != nil {
		return fmt.Errorf("failed to render welcome back template: %w", err)
	}

	subject := "Welcome Back to QKO!"

	return res.SendEmail(ctx, email, subject, htmlContent)
}

// Close closes the email service
func (res *ResendEmailService) Close() error {
	// Resend client doesn't need explicit closing
	return nil
}