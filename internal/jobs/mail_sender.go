package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	gomail "github.com/wneessen/go-mail"
)

// MailSenderWorker polls the EmailObjects table for emails with Status='New'
// and delivers them via SMTP. This replaces the .NET MailSender.SendEmailAsync
// and the separate mail-sender process that picks up queued emails.
type MailSenderWorker struct {
	emailRepo *repository.EmailRepository
	cfg       config.EmailConfig
	log       zerolog.Logger
	interval  time.Duration
	batchSize int
}

// NewMailSenderWorker creates a new SMTP mail sender worker.
func NewMailSenderWorker(
	emailRepo *repository.EmailRepository,
	cfg config.EmailConfig,
	interval time.Duration,
	log zerolog.Logger,
) *MailSenderWorker {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &MailSenderWorker{
		emailRepo: emailRepo,
		cfg:       cfg,
		log:       log.With().Str("component", "mail_sender").Logger(),
		interval:  interval,
		batchSize: 50,
	}
}

// Run starts the polling loop. It blocks until the context is cancelled.
func (w *MailSenderWorker) Run(ctx context.Context) {
	w.log.Info().
		Dur("interval", w.interval).
		Str("smtp_server", w.cfg.SMTPServer).
		Int("smtp_port", w.cfg.SMTPPort).
		Msg("mail sender worker started")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("mail sender worker stopping")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *MailSenderWorker) processBatch(ctx context.Context) {
	if w.emailRepo == nil {
		return
	}

	emails, err := w.emailRepo.GetNewEmails(ctx, w.batchSize)
	if err != nil {
		w.log.Error().Err(err).Msg("failed to fetch new emails")
		return
	}

	if len(emails) == 0 {
		return
	}

	w.log.Info().Int("count", len(emails)).Msg("processing email batch")

	for _, email := range emails {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Atomically mark as processing to prevent concurrent pickup.
		if err := w.emailRepo.MarkEmailProcessing(ctx, email.ID); err != nil {
			w.log.Debug().Err(err).Int("id", email.ID).Msg("email already picked up, skipping")
			continue
		}

		if err := w.sendEmail(email.To, email.CC, email.BCC, email.Subject, email.Body); err != nil {
			w.log.Error().Err(err).
				Int("id", email.ID).
				Str("to", email.To).
				Msg("failed to send email via SMTP")
			_ = w.emailRepo.UpdateEmailStatus(ctx, email.ID, "Failed", nil)
			continue
		}

		now := time.Now().UTC()
		if err := w.emailRepo.UpdateEmailStatus(ctx, email.ID, "Sent", &now); err != nil {
			w.log.Error().Err(err).Int("id", email.ID).Msg("failed to update email status to Sent")
		} else {
			w.log.Info().Int("id", email.ID).Str("to", email.To).Msg("email sent successfully")
		}
	}
}

// sendEmail delivers a single email via SMTP using go-mail.
// Mirrors the .NET MailSender.SendEmailAsync method.
func (w *MailSenderWorker) sendEmail(to, cc, bcc, subject, body string) error {
	msg := gomail.NewMsg()

	// From
	from := w.cfg.SenderEmail
	if from == "" {
		from = "pms@cbn.gov.ng"
	}
	if err := msg.From(from); err != nil {
		return err
	}
	if w.cfg.SenderDisplay != "" {
		if err := msg.FromFormat(w.cfg.SenderDisplay, from); err != nil {
			return err
		}
	}

	// To (may be semicolon or comma separated)
	toAddrs := splitAddresses(to)
	if len(toAddrs) == 0 {
		w.log.Warn().Msg("no recipient address, skipping email")
		return nil
	}
	if err := msg.To(toAddrs...); err != nil {
		return err
	}

	// CC
	if cc != "" {
		ccAddrs := splitAddresses(cc)
		if len(ccAddrs) > 0 {
			if err := msg.Cc(ccAddrs...); err != nil {
				return err
			}
		}
	}

	// BCC
	if bcc != "" {
		bccAddrs := splitAddresses(bcc)
		if len(bccAddrs) > 0 {
			if err := msg.Bcc(bccAddrs...); err != nil {
				return err
			}
		}
	}

	msg.Subject(subject)
	msg.SetBodyString(gomail.TypeTextHTML, body)

	// Create SMTP client
	port := w.cfg.SMTPPort
	if port == 0 {
		port = 25
	}

	client, err := gomail.NewClient(w.cfg.SMTPServer,
		gomail.WithPort(port),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(w.cfg.SenderEmail),
		gomail.WithPassword(w.cfg.SenderPassword),
		gomail.WithTLSPolicy(gomail.TLSOpportunistic),
	)
	if err != nil {
		return err
	}

	return client.DialAndSend(msg)
}

// splitAddresses splits a string of email addresses separated by semicolons
// or commas, trimming whitespace from each address.
func splitAddresses(s string) []string {
	s = strings.ReplaceAll(s, ",", ";")
	parts := strings.Split(s, ";")
	var result []string
	for _, p := range parts {
		addr := strings.TrimSpace(p)
		if addr != "" {
			result = append(result, addr)
		}
	}
	return result
}
