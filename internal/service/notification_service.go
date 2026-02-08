package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// notificationService implements the NotificationService interface.
// Mirrors the .NET NotificationService + NotificationTemplates pattern.
// It sends email notifications using template-based content.
// ---------------------------------------------------------------------------

// Notification template constants — direct conversion from
// CompetencyApp.Utilities/NotificationMessages.cs (NotificationTemplates class).
// .NET placeholders (%NAME%, %REQUEST_NAME%, etc.) are mapped to Go template
// syntax ({{.Name}}, {{.RequestName}}, etc.).
const (
	tplSlaGenericNewRequest = `<p>Dear {{.Name}}</p>` +
		`<p>A request on {{.RequestName}} has been assigned to you on {{.AssignedDate}}.</p>` +
		`<p>Kindly logon to your account on Performance Management System to treat this request ` +
		`as soon as possible to avoid breaching SLA of {{.SLAHours}} hours after initiation.</p>` +
		`<p>Thank you, <br/>CBN PMS</p>`

	tplGenericNewRequest = `<p>Dear {{.Name}}</p>` +
		`<p>A request on {{.RequestName}} has been assigned to you on {{.AssignedDate}}</p>` +
		`<p>Kindly logon to your account on Performance Management System to treat this request.</p>` +
		`<p>Thank you, <br/>CBN PMS</p>`

	tplAssignerGenericNewRequest = `<p>Dear {{.Name}}</p>` +
		`<p>Your request on {{.RequestName}} has been initated on {{.AssignedDate}}</p>` +
		`<p>Thank you, <br/>CBN PMS</p>`

	tplUpdateRequest = `<p>Dear {{.Name}}</p>` +
		`<p>You have treated the request on {{.RequestName}} on {{.TreatedDate}}</p>` +
		`<p>Thank you, <br/>CBN PMS</p>`

	tplAssignerUpdateRequest = `<p>Dear {{.Name}}</p>` +
		`<p>Your request on {{.RequestName}} has been treated on {{.TreatedDate}}</p>` +
		`<p>Thank you, <br/>CBN PMS</p>`
)

// notificationTemplateData holds values for notification template rendering.
type notificationTemplateData struct {
	Name         string
	RequestName  string
	AssignedDate string
	TreatedDate  string
	SLAHours     int
}

// Pre-parsed templates for notification messages.
var (
	notifTplSlaNew     = template.Must(template.New("sla_new").Parse(tplSlaGenericNewRequest))
	notifTplNew        = template.Must(template.New("new").Parse(tplGenericNewRequest))
	notifTplAssignerNew = template.Must(template.New("assigner_new").Parse(tplAssignerGenericNewRequest))
	notifTplUpdate     = template.Must(template.New("update").Parse(tplUpdateRequest))
	notifTplAssignerUp = template.Must(template.New("assigner_update").Parse(tplAssignerUpdateRequest))
)

type notificationService struct {
	emailSvc EmailService
	cfg      *config.Config
	log      zerolog.Logger
}

func newNotificationService(
	emailSvc EmailService,
	cfg *config.Config,
	log zerolog.Logger,
) NotificationService {
	return &notificationService{
		emailSvc: emailSvc,
		cfg:      cfg,
		log:      log.With().Str("service", "notification").Logger(),
	}
}

// Send sends a plain-text/HTML notification email to the specified user.
func (s *notificationService) Send(ctx context.Context, userID string, message string) error {
	return s.emailSvc.SendEmail(ctx, userID, "PMS Notification", message)
}

// SendNewRequestNotification sends a notification for a new feedback request.
// Uses the SLA template if hasSLA is true, otherwise the generic template.
// Mirrors .NET NotificationTemplates.SlaGenericNewRequest / GenericNewRequest.
func (s *notificationService) SendNewRequestNotification(
	ctx context.Context,
	recipientEmail, recipientName, requestName string,
	assignedDate time.Time,
	hasSLA bool,
	slaHours int,
) error {
	data := notificationTemplateData{
		Name:         recipientName,
		RequestName:  requestName,
		AssignedDate: assignedDate.Format("Monday, 02 January 2006 03:04 PM"),
		SLAHours:     slaHours,
	}

	var tpl *template.Template
	if hasSLA {
		tpl = notifTplSlaNew
	} else {
		tpl = notifTplNew
	}

	body, err := renderNotifTemplate(tpl, data)
	if err != nil {
		return fmt.Errorf("notification: rendering new request template: %w", err)
	}

	subject := fmt.Sprintf("PMS | New Request: %s", requestName)
	return s.emailSvc.SendEmail(ctx, recipientEmail, subject, body)
}

// SendAssignerNewRequestNotification notifies the assigner that their request was initiated.
// Mirrors .NET NotificationTemplates.AssignerGenericNewRequest.
func (s *notificationService) SendAssignerNewRequestNotification(
	ctx context.Context,
	recipientEmail, recipientName, requestName string,
	assignedDate time.Time,
) error {
	data := notificationTemplateData{
		Name:         recipientName,
		RequestName:  requestName,
		AssignedDate: assignedDate.Format("Monday, 02 January 2006 03:04 PM"),
	}

	body, err := renderNotifTemplate(notifTplAssignerNew, data)
	if err != nil {
		return fmt.Errorf("notification: rendering assigner new request template: %w", err)
	}

	subject := fmt.Sprintf("PMS | Request Initiated: %s", requestName)
	return s.emailSvc.SendEmail(ctx, recipientEmail, subject, body)
}

// SendRequestTreatedNotification notifies when a request is treated.
// Mirrors .NET NotificationTemplates.UpdateRequest.
func (s *notificationService) SendRequestTreatedNotification(
	ctx context.Context,
	recipientEmail, recipientName, requestName string,
	treatedDate time.Time,
) error {
	data := notificationTemplateData{
		Name:        recipientName,
		RequestName: requestName,
		TreatedDate: treatedDate.Format("Monday, 02 January 2006 03:04 PM"),
	}

	body, err := renderNotifTemplate(notifTplUpdate, data)
	if err != nil {
		return fmt.Errorf("notification: rendering treated request template: %w", err)
	}

	subject := fmt.Sprintf("PMS | Request Treated: %s", requestName)
	return s.emailSvc.SendEmail(ctx, recipientEmail, subject, body)
}

// SendAssignerRequestTreatedNotification notifies the assigner when their request is treated.
// Mirrors .NET NotificationTemplates.AssignerUpdateRequest.
func (s *notificationService) SendAssignerRequestTreatedNotification(
	ctx context.Context,
	recipientEmail, recipientName, requestName string,
	treatedDate time.Time,
) error {
	data := notificationTemplateData{
		Name:        recipientName,
		RequestName: requestName,
		TreatedDate: treatedDate.Format("Monday, 02 January 2006 03:04 PM"),
	}

	body, err := renderNotifTemplate(notifTplAssignerUp, data)
	if err != nil {
		return fmt.Errorf("notification: rendering assigner treated template: %w", err)
	}

	subject := fmt.Sprintf("PMS | Your Request Treated: %s", requestName)
	return s.emailSvc.SendEmail(ctx, recipientEmail, subject, body)
}

func renderNotifTemplate(tpl *template.Template, data notificationTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func init() {
	// Compile-time interface compliance check.
	var _ NotificationService = (*notificationService)(nil)
}
