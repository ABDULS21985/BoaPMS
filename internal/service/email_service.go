package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/email"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// emailService implements the EmailService interface.
// Converted from the .NET NotificationService class which saves outbound
// email records to the EmailService SQL Server database (EmailLogs table).
// A separate mail-sender process picks up rows with Status = 'New' and
// delivers them via SMTP.
// ---------------------------------------------------------------------------

// defaultSenderEmail is the fallback sender address when the global setting
// SENDER_EMAIL is not configured. Mirrors the .NET default.
const defaultSenderEmail = "pms@cbn.gov.ng"

// emailService persists email records to the EmailLogs table via the
// EmailSvcSQL (sqlx) connection, matching the .NET EmailServiceDBContext.
type emailService struct {
	emailDB *sqlx.DB
	gs      GlobalSettingService
	cfg     *config.Config
	log     zerolog.Logger
}

func newEmailService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	gs GlobalSettingService,
) EmailService {
	return &emailService{
		emailDB: repos.EmailSvcSQL,
		gs:      gs,
		cfg:     cfg,
		log:     log.With().Str("service", "email").Logger(),
	}
}

// SendEmail saves a single email record to the database for asynchronous
// delivery. It checks the ENABLE_EMAIL_NOTIFICATION global setting before
// persisting. Mirrors the .NET NotificationService.SaveEmailToDb method.
func (s *emailService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	return s.saveEmailToDb(ctx, to, subject, body, "SEND_EMAIL")
}

// SendBulkEmail saves one email record per recipient. Each record is
// independently persisted so partial failures do not block other recipients.
func (s *emailService) SendBulkEmail(ctx context.Context, to []string, subject string, body string) error {
	var errs []string
	for _, recipient := range to {
		if err := s.saveEmailToDb(ctx, recipient, subject, body, "SEND_BULK_EMAIL"); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", recipient, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to save emails for %d recipient(s): %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}

// saveEmailToDb is the core persistence method mirroring the .NET
// NotificationService.SaveEmailToDb. It:
//  1. Reads SENDER_EMAIL and ENABLE_EMAIL_NOTIFICATION from global settings.
//  2. If notifications are disabled, logs a warning and returns early.
//  3. Inserts a row into the EmailLogs table with Status = 'New'.
func (s *emailService) saveEmailToDb(ctx context.Context, emailTo, subject, body, action string) error {
	if s.emailDB == nil {
		s.log.Warn().Msg("EmailSvcSQL database is not configured, skipping email save")
		return nil
	}

	// Resolve global settings with safe defaults (mirrors the .NET try/catch pattern)
	senderEmail := defaultSenderEmail
	enableNotification := false

	if s.gs != nil {
		if val, err := s.gs.GetStringValue(ctx, "SENDER_EMAIL"); err == nil && val != "" {
			senderEmail = val
		}
		if val, err := s.gs.GetBoolValue(ctx, "ENABLE_EMAIL_NOTIFICATION"); err == nil {
			enableNotification = val
		}
	}

	if !enableNotification {
		s.log.Warn().
			Str("to", emailTo).
			Str("action", action).
			Msg("email notification is not enabled")
		return nil
	}

	now := time.Now()
	const insertQuery = `
		INSERT INTO EmailLogs (
			[From], [To], [Subject], [Body], [Action], [Status],
			[NoOfRetry], [AppSource], [CreatedBy], [DateCreated]
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6,
			@p7, @p8, @p9, @p10
		)`

	_, err := s.emailDB.ExecContext(ctx, insertQuery,
		senderEmail,  // From
		emailTo,      // To
		subject,      // Subject
		body,         // Body
		action,       // Action
		"New",        // Status
		0,            // NoOfRetry
		"PMS",        // AppSource
		"SYSTEM",     // CreatedBy
		now,          // DateCreated
	)
	if err != nil {
		s.log.Error().Err(err).
			Str("to", emailTo).
			Str("action", action).
			Msg("failed to save email to database")
		return fmt.Errorf("saving email to database: %w", err)
	}

	s.log.Info().
		Str("to", emailTo).
		Str("action", action).
		Msg("email saved successfully")
	return nil
}

// SendEmailWithCC saves an email record with CC recipients for asynchronous delivery.
// Mirrors the .NET pattern where emails are sent with CC addresses (e.g. supervisor CCs).
func (s *emailService) SendEmailWithCC(ctx context.Context, to string, cc []string, subject string, body string) error {
	if s.emailDB == nil {
		s.log.Warn().Msg("EmailSvcSQL database is not configured, skipping email save")
		return nil
	}

	// Resolve global settings with safe defaults
	senderEmail := defaultSenderEmail
	enableNotification := false

	if s.gs != nil {
		if val, err := s.gs.GetStringValue(ctx, "SENDER_EMAIL"); err == nil && val != "" {
			senderEmail = val
		}
		if val, err := s.gs.GetBoolValue(ctx, "ENABLE_EMAIL_NOTIFICATION"); err == nil {
			enableNotification = val
		}
	}

	if !enableNotification {
		s.log.Warn().
			Str("to", to).
			Msg("email notification is not enabled")
		return nil
	}

	now := time.Now()
	ccStr := strings.Join(cc, ";")

	const insertQuery = `
		INSERT INTO EmailLogs (
			[From], [To], [Cc], [Subject], [Body], [Action], [Status],
			[NoOfRetry], [AppSource], [CreatedBy], [DateCreated]
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7,
			@p8, @p9, @p10, @p11
		)`

	_, err := s.emailDB.ExecContext(ctx, insertQuery,
		senderEmail,       // From
		to,                // To
		ccStr,             // Cc
		subject,           // Subject
		body,              // Body
		"SEND_EMAIL_CC",   // Action
		"New",             // Status
		0,                 // NoOfRetry
		"PMS",             // AppSource
		"SYSTEM",          // CreatedBy
		now,               // DateCreated
	)
	if err != nil {
		s.log.Error().Err(err).
			Str("to", to).
			Str("cc", ccStr).
			Msg("failed to save email with CC to database")
		return fmt.Errorf("saving email with CC to database: %w", err)
	}

	s.log.Info().
		Str("to", to).
		Str("cc", ccStr).
		Msg("email with CC saved successfully")
	return nil
}

// ProcessEmail dispatches an email notification based on the request title.
// Mirrors .NET MailServices.ProcessEmail which uses a switch on EmailRequest.Title
// to determine which notification method to invoke. In the Go implementation,
// the actual email body construction and sending is delegated to saveEmailToDb,
// as we use a database-driven email queue rather than Hangfire background jobs.
//
// The supported notification titles mirror the .NET switch cases:
//   - NewCompetency, ApprovedCompetency, NewReviewPeriod, ApprovedReviewPeriod
//   - Self, Peers, Subordinates, Superior
//   - DevTaskAssigned, DevTaskCompleted, DevTaskApproved
//   - UpdateJobRole, ApprovedJobRole, RejectJobRole, NotifyCmsTeamApprovedJobRole
//   - ReminderEmail
func (s *emailService) ProcessEmail(ctx context.Context, req interface{}) (interface{}, error) {
	emailReq, ok := req.(*performance.EmailRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *performance.EmailRequest")
	}

	s.log.Info().
		Str("title", emailReq.Title).
		Str("userId", emailReq.UserID).
		Msg("processing email notification")

	// Determine subject based on notification title (mirrors .NET switch cases)
	var subject string
	switch emailReq.Title {
	case "NewCompetency":
		subject = "Review Request for New Competency Entry"
	case "ApprovedCompetency":
		subject = "Competency Management System | Approved Competency Notification"
	case "NewReviewPeriod":
		subject = "Competency Management System | Request to Approve New Competency Review Period Notification"
	case "ApprovedReviewPeriod":
		subject = "Competency Management System | New Competency Review Period Notification"
	case "Self":
		subject = emailReq.EmailDescription + " Competency Review Completed by Self Notification"
	case "Peers":
		subject = emailReq.EmailDescription + " Competency Review Completed by Peer Notification"
	case "Subordinates":
		subject = emailReq.EmailDescription + " Competency Review Completed by Subordinate Notification"
	case "Superior":
		subject = emailReq.EmailDescription + " Reviews Completed by " + emailReq.Title + " Notification"
	case "DevTaskAssigned":
		subject = "Competency Management System | Assigned Development Task Notification"
	case "DevTaskCompleted":
		subject = "Competency Management System | Development Task Completed Notification"
	case "DevTaskApproved":
		subject = "Competency Management System | Development Task Approved Notification"
	case "UpdateJobRole":
		subject = "Competency Management System | Job Role Update Request Notification"
	case "ApprovedJobRole":
		subject = "Competency Management System | Job Role Approval Notification"
	case "RejectJobRole":
		subject = "Competency Management System | Job Role Rejection Notification"
	case "NotifyCmsTeamApprovedJobRole":
		subject = "Competency Management System | Job Role Approval Notification"
	case "ReminderEmail":
		subject = "Competency Review Reminder Notification"
	default:
		s.log.Warn().Str("title", emailReq.Title).Msg("unknown email notification title, skipping")
		return &performance.GenericResponseVm{
			IsSuccess: true,
			Message:   "The request was Sent Successfully!",
		}, nil
	}

	// Render the HTML body using the appropriate template from the email package.
	// This replaces the .NET HtmlMessageTemplate string-interpolation pattern.
	// Template data is built from the EmailRequest fields; recipient resolution
	// via ErpEmployeeService will be added when that service is fully implemented.
	appURL := s.cfg.Email.ApplicationURL
	data := email.TemplateData{
		Recipient:      emailReq.UserID,
		Requester:      emailReq.UserID,
		RequestorName:  emailReq.UserID,
		AppURL:         appURL,
		CompetencyInfo: emailReq.EmailDescription,
		ReviewPeriod:   emailReq.EmailDescription,
		DeadlineDate:   emailReq.EmailDescription,
		CompetencyName: emailReq.EmailDescription,
		Description:    emailReq.EmailDescription,
	}

	// Map EmailRequest.Title to the template key used in the email package.
	var templateKey string
	switch emailReq.Title {
	case "NewCompetency":
		templateKey = email.TplNewCompetencyToApprover
	case "ApprovedCompetency":
		templateKey = email.TplApprovedCompetencyToRequestor
	case "NewReviewPeriod":
		templateKey = email.TplNewReviewPeriodToApprover
	case "ApprovedReviewPeriod":
		templateKey = email.TplApprovedReviewPeriodToAll
	case "Self":
		templateKey = email.TplSelfReviewCompleted
	case "Peers":
		templateKey = email.TplPeersReviewCompleted
	case "Subordinates":
		templateKey = email.TplSubordinateReviewCompleted
	case "Superior":
		templateKey = email.TplSuperiorReviewCompleted
	case "DevTaskAssigned":
		templateKey = email.TplDevTaskAssigned
	case "DevTaskCompleted":
		templateKey = email.TplDevTaskCompleted
	case "DevTaskApproved":
		templateKey = email.TplDevTaskApproved
	case "UpdateJobRole":
		templateKey = email.TplJobRoleUpdateRequest
	case "ApprovedJobRole":
		templateKey = email.TplApprovedJobRole
	case "RejectJobRole":
		templateKey = email.TplRejectJobRole
	case "NotifyCmsTeamApprovedJobRole":
		templateKey = email.TplCMSTeamApprovedJobRole
	case "ReminderEmail":
		templateKey = email.TplReminderMessage
	}

	body, err := email.Render(templateKey, data)
	if err != nil {
		// Fallback to a generic body if template rendering fails.
		s.log.Warn().Err(err).Str("template", templateKey).Msg("template rendering failed, using fallback")
		body = fmt.Sprintf(
			"<p>Notification: <strong>%s</strong></p><p>%s</p><p>User: %s</p>",
			emailReq.Title, emailReq.EmailDescription, emailReq.UserID,
		)
	}

	// Save the email for async delivery via the mail sender worker.
	if err := s.saveEmailToDb(ctx, "", subject, body, emailReq.Title); err != nil {
		s.log.Error().Err(err).Str("title", emailReq.Title).Msg("failed to process email notification")
		return nil, fmt.Errorf("processing email notification: %w", err)
	}

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "The request was Sent Successfully!",
	}, nil
}

func init() {
	// Compile-time interface compliance check.
	var _ EmailService = (*emailService)(nil)
}
