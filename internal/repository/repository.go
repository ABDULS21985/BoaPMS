package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Container holds all database connections and domain repositories.
// This is the Go equivalent of the .NET DI-registered DbContexts and repositories.
type Container struct {
	// GORM connection for the primary PostgreSQL database (ORM operations)
	GormDB *gorm.DB

	// sqlx connections for raw SQL queries (Dapper-style)
	CoreSQL     *sqlx.DB // PostgreSQL
	ErpSQL      *sqlx.DB // SQL Server — ERP data (optional)
	StaffIDSQL  *sqlx.DB // SQL Server — Staff ID mask (optional)
	EmailSvcSQL *sqlx.DB // SQL Server — Email service (optional)
	SasSQL      *sqlx.DB // SQL Server — SAS (optional)

	// Audit interceptor (auto-logs entity changes)
	Audit *AuditInterceptor

	// ── Domain-specific repositories (GORM-based, primary PostgreSQL) ──
	Competency  *CompetencyRepository
	Performance *PerformanceRepository
	Project     *ProjectRepository
	Feedback    *FeedbackRepository
	Grievance   *GrievanceRepository

	// ── Multi-database repositories (sqlx-based, SQL Server) ───────────
	Erp   *ErpRepository
	Staff *StaffRepository
	Email *EmailRepository
	Sas   *SasRepository

	log zerolog.Logger
}

// New initializes all database connections, the audit interceptor,
// and all domain-specific repositories.
func New(cfg *config.Config, log zerolog.Logger) (*Container, error) {
	dm, err := NewDatabaseManager(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("initializing database manager: %w", err)
	}

	c := &Container{
		GormDB:      dm.CoreGorm,
		CoreSQL:     dm.CoreSQL,
		ErpSQL:      dm.ErpSQL,
		StaffIDSQL:  dm.StaffIDSQL,
		EmailSvcSQL: dm.EmailSvcSQL,
		SasSQL:      dm.SasSQL,
		log:         log,
	}

	// Register audit log interceptor on the GORM connection
	c.Audit = NewAuditInterceptor(dm.CoreGorm, log)

	// Initialize domain-specific repositories (GORM — primary PostgreSQL)
	c.Competency = NewCompetencyRepository(dm.CoreGorm)
	c.Performance = NewPerformanceRepository(dm.CoreGorm)
	c.Project = NewProjectRepository(dm.CoreGorm)
	c.Feedback = NewFeedbackRepository(dm.CoreGorm)
	c.Grievance = NewGrievanceRepository(dm.CoreGorm)

	// Initialize multi-database repositories (sqlx — SQL Server, optional)
	c.Erp = NewErpRepository(dm.ErpSQL)
	c.Staff = NewStaffRepository(dm.StaffIDSQL)
	c.Email = NewEmailRepository(dm.EmailSvcSQL)
	c.Sas = NewSasRepository(dm.SasSQL)

	log.Info().Msg("Repository container initialized with all domain repositories")
	return c, nil
}

// Close closes all underlying database connections.
func (c *Container) Close() {
	if c.CoreSQL != nil {
		c.CoreSQL.Close()
	}
	for _, db := range []*sqlx.DB{c.ErpSQL, c.StaffIDSQL, c.EmailSvcSQL, c.SasSQL} {
		if db != nil {
			db.Close()
		}
	}
}

// Ping verifies the primary database connection is alive.
func (c *Container) Ping(ctx context.Context) error {
	sqlDB, err := c.GormDB.DB()
	if err != nil {
		return fmt.Errorf("getting underlying DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}
