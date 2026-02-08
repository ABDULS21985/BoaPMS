package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/erp"
	"github.com/jmoiron/sqlx"
)

// EmailRepository provides data access for the Email Service SQL Server database.
type EmailRepository struct {
	db *sqlx.DB // EmailSvcSQL connection
}

// NewEmailRepository creates a new Email repository.
func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	if db == nil {
		return nil
	}
	return &EmailRepository{db: db}
}

// InsertEmail inserts a new email record into the email service database.
func (r *EmailRepository) InsertEmail(ctx context.Context, email *erp.EmailObject) error {
	if r.db == nil {
		return fmt.Errorf("emailRepo.InsertEmail: Email database not configured")
	}
	query := `INSERT INTO dbo.EmailObjects
		([From], [To], CC, BCC, Subject, Body, Status, NoOfRetry,
		 ExpectedSendDate, Action, AppSource, CreatedBy, DateCreated, EmailGuid)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14)`
	_, err := r.db.ExecContext(ctx, query,
		email.From, email.To, email.CC, email.BCC,
		email.Subject, email.Body, email.Status, email.NoOfRetry,
		email.ExpectedSendDate, email.Action, email.AppSource,
		email.CreatedBy, time.Now().UTC(), email.EmailGUID)
	if err != nil {
		return fmt.Errorf("emailRepo.InsertEmail: %w", err)
	}
	return nil
}

// GetPendingEmails retrieves emails with "Pending" status.
func (r *EmailRepository) GetPendingEmails(ctx context.Context) ([]erp.EmailObject, error) {
	if r.db == nil {
		return nil, fmt.Errorf("emailRepo.GetPendingEmails: Email database not configured")
	}
	var results []erp.EmailObject
	err := r.db.SelectContext(ctx, &results,
		`SELECT * FROM dbo.EmailObjects WHERE Status = 'Pending' ORDER BY DateCreated DESC`)
	if err != nil {
		return nil, fmt.Errorf("emailRepo.GetPendingEmails: %w", err)
	}
	return results, nil
}

// GetNewEmails retrieves emails with Status='New' that are ready to send.
// An email is ready when its ExpectedSendDate is null (immediate) or in the past.
// Limit controls the batch size to prevent overloading the SMTP server.
func (r *EmailRepository) GetNewEmails(ctx context.Context, limit int) ([]erp.EmailObject, error) {
	if r.db == nil {
		return nil, fmt.Errorf("emailRepo.GetNewEmails: Email database not configured")
	}
	var results []erp.EmailObject
	err := r.db.SelectContext(ctx, &results,
		`SELECT TOP(@p1) * FROM dbo.EmailObjects
		 WHERE Status = 'New'
		   AND (ExpectedSendDate IS NULL OR ExpectedSendDate <= GETDATE())
		 ORDER BY DateCreated ASC`, limit)
	if err != nil {
		return nil, fmt.Errorf("emailRepo.GetNewEmails: %w", err)
	}
	return results, nil
}

// MarkEmailProcessing atomically transitions an email from 'New' to 'Processing'
// to prevent concurrent pickup by multiple mail sender instances.
func (r *EmailRepository) MarkEmailProcessing(ctx context.Context, id int) error {
	if r.db == nil {
		return fmt.Errorf("emailRepo.MarkEmailProcessing: Email database not configured")
	}
	result, err := r.db.ExecContext(ctx,
		`UPDATE dbo.EmailObjects SET Status = 'Processing', LastUpdatedDate = @p1
		 WHERE Id = @p2 AND Status = 'New'`,
		time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("emailRepo.MarkEmailProcessing: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("emailRepo.MarkEmailProcessing: email %d already picked up", id)
	}
	return nil
}

// UpdateEmailStatus updates the status and retry count of an email.
func (r *EmailRepository) UpdateEmailStatus(ctx context.Context, id int, status string, actualSendDate *time.Time) error {
	if r.db == nil {
		return fmt.Errorf("emailRepo.UpdateEmailStatus: Email database not configured")
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE dbo.EmailObjects SET Status = @p1, ActualSendDate = @p2, NoOfRetry = NoOfRetry + 1,
		 LastUpdatedDate = @p3 WHERE Id = @p4`,
		status, actualSendDate, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("emailRepo.UpdateEmailStatus: %w", err)
	}
	return nil
}
