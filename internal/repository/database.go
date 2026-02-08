package repository

import (
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseManager manages multiple GORM and sqlx database connections.
// This replaces the .NET pattern of multiple DbContext registrations.
type DatabaseManager struct {
	// GORM connections (for ORM operations)
	CoreGorm *gorm.DB // PostgreSQL — primary PMS database

	// sqlx connections (for raw SQL / Dapper-style queries)
	CoreSQL      *sqlx.DB // PostgreSQL — primary
	ErpSQL       *sqlx.DB // SQL Server — ERP data (optional)
	StaffIDSQL   *sqlx.DB // SQL Server — Staff ID mask (optional)
	EmailSvcSQL  *sqlx.DB // SQL Server — Email service (optional)
	SasSQL       *sqlx.DB // SQL Server — SAS data (optional)

	log zerolog.Logger
}

// NewDatabaseManager creates a DatabaseManager with all configured connections.
func NewDatabaseManager(cfg *config.Config, log zerolog.Logger) (*DatabaseManager, error) {
	dm := &DatabaseManager{log: log}

	// Primary PostgreSQL via GORM
	gormDB, err := connectGormPostgres(cfg.Database.Core, log)
	if err != nil {
		return nil, fmt.Errorf("connecting GORM to core PostgreSQL: %w", err)
	}
	dm.CoreGorm = gormDB

	// Primary PostgreSQL via sqlx (for raw queries)
	coreSQL, err := connectSqlxPostgres(cfg.Database.Core)
	if err != nil {
		return nil, fmt.Errorf("connecting sqlx to core PostgreSQL: %w", err)
	}
	dm.CoreSQL = coreSQL

	// Optional SQL Server connections via sqlx
	dm.ErpSQL = connectSqlxSQLServerOptional(cfg.Database.ErpData, "ERP", log)
	dm.StaffIDSQL = connectSqlxSQLServerOptional(cfg.Database.StaffIDMask, "StaffIDMask", log)
	dm.EmailSvcSQL = connectSqlxSQLServerOptional(cfg.Database.EmailSvc, "EmailService", log)
	dm.SasSQL = connectSqlxSQLServerOptional(cfg.Database.Sas, "SAS", log)

	return dm, nil
}

// Close closes all database connections.
func (dm *DatabaseManager) Close() {
	if dm.CoreSQL != nil {
		dm.CoreSQL.Close()
	}
	for _, db := range []*sqlx.DB{dm.ErpSQL, dm.StaffIDSQL, dm.EmailSvcSQL, dm.SasSQL} {
		if db != nil {
			db.Close()
		}
	}
}

func connectGormPostgres(cfg config.PostgresConfig, log zerolog.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()

	logLevel := gormlogger.Warn
	if log.GetLevel() <= zerolog.DebugLevel {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Info().Str("db", cfg.Database).Msg("Connected GORM to PostgreSQL")
	return db, nil
}

func connectSqlxPostgres(cfg config.PostgresConfig) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", cfg.DSN())
}

func connectSqlxSQLServerOptional(cfg config.SQLServerConfig, name string, log zerolog.Logger) *sqlx.DB {
	if cfg.Host == "" {
		log.Info().Str("db", name).Msg("SQL Server not configured, skipping")
		return nil
	}
	db, err := sqlx.Connect("sqlserver", cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Str("db", name).Msg("Failed to connect to SQL Server, skipping")
		return nil
	}
	log.Info().Str("db", name).Str("host", cfg.Host).Msg("Connected sqlx to SQL Server")
	return db
}
