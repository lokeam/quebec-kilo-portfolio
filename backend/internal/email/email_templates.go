package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// TemplateEngine handles email template rendering
type TemplateEngine struct {
	templates map[string]*template.Template
	templateDir string
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(templateDir string) (*TemplateEngine, error) {
	engine := &TemplateEngine{
		templates:   make(map[string]*template.Template),
		templateDir: templateDir,
	}

	// Load all templates
	if err := engine.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return engine, nil
}

// loadTemplates loads all HTML templates from the template directory
func (te *TemplateEngine) loadTemplates() error {
	templateFiles := []string{
		"deletion_request.html",
		"grace_period_reminder.html",
		"deletion_confirmation.html",
		"data_export.html",
		"welcome_back.html",
	}

	for _, filename := range templateFiles {
		filepath := filepath.Join(te.templateDir, filename)

		// Read template file
		content, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", filename, err)
		}

		// Parse template
		tmpl, err := template.New(filename).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", filename, err)
		}

		// Store template
		te.templates[filename] = tmpl
	}

	return nil
}

// TemplateData contains data for email templates
type TemplateData struct {
	UserID         string
	Email          string
	GracePeriodEnd time.Time
	DaysLeft       int
	ExportURL      string
	DeletionDate   time.Time
}

// renderTemplate renders a template with the given data
func (te *TemplateEngine) renderTemplate(templateName string, data TemplateData) (string, error) {
	tmpl, exists := te.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// RenderDeletionRequest renders the deletion request email template
func (te *TemplateEngine) RenderDeletionRequest(userID, email string, gracePeriodEnd time.Time) (string, error) {
	data := TemplateData{
		UserID:         userID,
		Email:          email,
		GracePeriodEnd: gracePeriodEnd,
	}
	return te.renderTemplate("deletion_request.html", data)
}

// RenderGracePeriodReminder renders the grace period reminder email template
func (te *TemplateEngine) RenderGracePeriodReminder(userID, email string, daysLeft int) (string, error) {
	data := TemplateData{
		UserID:   userID,
		Email:    email,
		DaysLeft: daysLeft,
	}
	return te.renderTemplate("grace_period_reminder.html", data)
}

// RenderDeletionConfirmation renders the deletion confirmation email template
func (te *TemplateEngine) RenderDeletionConfirmation(userID, email string) (string, error) {
	data := TemplateData{
		UserID:       userID,
		Email:        email,
		DeletionDate: time.Now(),
	}
	return te.renderTemplate("deletion_confirmation.html", data)
}

// RenderDataExport renders the data export email template
func (te *TemplateEngine) RenderDataExport(userID, email, exportURL string) (string, error) {
	data := TemplateData{
		UserID:    userID,
		Email:     email,
		ExportURL: exportURL,
	}
	return te.renderTemplate("data_export.html", data)
}

// RenderWelcomeBack renders the welcome back email template
func (te *TemplateEngine) RenderWelcomeBack(userID, email string) (string, error) {
	data := TemplateData{
		UserID: userID,
		Email:  email,
	}
	return te.renderTemplate("welcome_back.html", data)
}