package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server           ServerConfig           `mapstructure:"server"`
	Database         DatabaseConfig         `mapstructure:"database"`
	JWT              JWTConfig              `mapstructure:"jwt"`
	ActiveDirectory  ActiveDirectoryConfig  `mapstructure:"active_directory"`
	Email            EmailConfig            `mapstructure:"email"`
	Bitly            BitlyConfig            `mapstructure:"bitly"`
	CORS             CORSConfig             `mapstructure:"cors"`
	Logging          LoggingConfig          `mapstructure:"logging"`
	Redis            RedisConfig            `mapstructure:"redis"`
	Jobs             JobsConfig             `mapstructure:"jobs"`
	APIKey           string                 `mapstructure:"api_key"`
	HangfireSchema   string                 `mapstructure:"hangfire_schema"`
	ReCaptcha        GoogleReCaptchaConfig  `mapstructure:"recaptcha"`
	PasswordGen      PasswordGenConfig      `mapstructure:"password_gen"`
	General          GeneralConfig          `mapstructure:"general"`
	RSA              RSAConfig              `mapstructure:"rsa"`
	Storage          StorageConfig          `mapstructure:"storage"`
}

// JobsConfig holds background job processing settings.
type JobsConfig struct {
	WorkerPoolSize     int           `mapstructure:"worker_pool_size"`
	WorkerQueueSize    int           `mapstructure:"worker_queue_size"`
	MailSenderInterval time.Duration `mapstructure:"mail_sender_interval"`
	CronSchedule       string       `mapstructure:"cron_schedule"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds all database connection settings.
type DatabaseConfig struct {
	Core        PostgresConfig   `mapstructure:"core"`
	Hangfire    PostgresConfig   `mapstructure:"hangfire"`
	ErpData     SQLServerConfig  `mapstructure:"erp_data"`
	StaffIDMask SQLServerConfig  `mapstructure:"staff_id_mask"`
	EmailSvc    SQLServerConfig  `mapstructure:"email_service"`
	Sas         SQLServerConfig  `mapstructure:"sas"`
}

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN returns the PostgreSQL connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// SQLServerConfig holds SQL Server connection settings.
type SQLServerConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// DSN returns the SQL Server connection string.
func (c SQLServerConfig) DSN() string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	)
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	Issuer             string        `mapstructure:"issuer"`
	Audience           string        `mapstructure:"audience"`
	TokenExpiryMinutes time.Duration `mapstructure:"token_expiry_minutes"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}

// ActiveDirectoryConfig holds LDAP/AD settings.
type ActiveDirectoryConfig struct {
	LdapURL  string `mapstructure:"ldap_url"`
	Domain   string `mapstructure:"domain"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// EmailConfig holds SMTP email settings.
type EmailConfig struct {
	SMTPServer     string `mapstructure:"smtp_server"`
	SMTPPort       int    `mapstructure:"smtp_port"`
	SenderEmail    string `mapstructure:"sender_email"`
	SenderPassword string `mapstructure:"sender_password"`
	SenderDisplay  string `mapstructure:"sender_display"`
	ApplicationURL string `mapstructure:"application_url"`
	ToAllStaff     string `mapstructure:"to_all_staff"`
}

// BitlyConfig holds Bitly API settings.
type BitlyConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowAll       bool     `mapstructure:"allow_all"`
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level          string `mapstructure:"level"`
	Format         string `mapstructure:"format"`
	FilePath       string `mapstructure:"file_path"`
	MaxSizeMB      int    `mapstructure:"max_size_mb"`
	MaxBackups     int    `mapstructure:"max_backups"`
	MaxAgeDays     int    `mapstructure:"max_age_days"`
	Compress       bool   `mapstructure:"compress"`
	ConsoleEnabled bool   `mapstructure:"console_enabled"`
	FileEnabled    bool   `mapstructure:"file_enabled"`
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// GoogleReCaptchaConfig holds Google reCAPTCHA verification settings.
type GoogleReCaptchaConfig struct {
	SiteKey   string `mapstructure:"site_key"`
	SecretKey string `mapstructure:"secret_key"`
	VerifyURL string `mapstructure:"verify_url"`
}

// PasswordGenConfig holds password generation policy settings.
type PasswordGenConfig struct {
	MaxIdenticalConsecutiveChars int    `mapstructure:"max_identical_consecutive_chars"`
	LowercaseChars              string `mapstructure:"lowercase_chars"`
	UppercaseChars              string `mapstructure:"uppercase_chars"`
	NumericChars                string `mapstructure:"numeric_chars"`
	SpecialChars                string `mapstructure:"special_chars"`
	SpaceChar                   string `mapstructure:"space_char"`
	PasswordLengthMin           int    `mapstructure:"password_length_min"`
	PasswordLengthMax           int    `mapstructure:"password_length_max"`
}

// GeneralConfig holds general application settings.
type GeneralConfig struct {
	Domain       string `mapstructure:"domain"`
	Name         string `mapstructure:"name"`
	Logo         string `mapstructure:"logo"`
	NoReplyEmail string `mapstructure:"no_reply_email"`
}

// RSAConfig holds RSA integration settings.
type RSAConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

// StorageConfig holds file storage settings.
type StorageConfig struct {
	BasePath string `mapstructure:"base_path"`
}

// Load reads the configuration from files and environment variables.
func Load() (*Config, error) {
	v := viper.New()

	// Set config file name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/pms-api")

	// Set defaults
	setDefaults(v)

	// Read base config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		// Config file not found — rely on defaults and env vars
	}

	// Load environment-specific overrides (like appsettings.Development.json)
	env := v.GetString("APP_ENV")
	if env != "" {
		v.SetConfigName(fmt.Sprintf("config.%s", strings.ToLower(env)))
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("reading env config file: %w", err)
			}
		}
	}

	// Environment variables override (PMS_ prefix)
	v.SetEnvPrefix("PMS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")

	// Database - Core (PostgreSQL)
	v.SetDefault("database.core.host", "localhost")
	v.SetDefault("database.core.port", 5432)
	v.SetDefault("database.core.database", "pms_db")
	v.SetDefault("database.core.username", "postgres")
	v.SetDefault("database.core.password", "postgres")
	v.SetDefault("database.core.ssl_mode", "disable")
	v.SetDefault("database.core.max_open_conns", 25)
	v.SetDefault("database.core.max_idle_conns", 10)
	v.SetDefault("database.core.conn_max_lifetime", "5m")

	// Database - Hangfire (PostgreSQL)
	v.SetDefault("database.hangfire.host", "localhost")
	v.SetDefault("database.hangfire.port", 5432)
	v.SetDefault("database.hangfire.database", "pms_db")
	v.SetDefault("database.hangfire.username", "postgres")
	v.SetDefault("database.hangfire.password", "postgres")
	v.SetDefault("database.hangfire.ssl_mode", "disable")

	// Database - ERP (SQL Server)
	v.SetDefault("database.erp_data.host", "localhost")
	v.SetDefault("database.erp_data.port", 1433)
	v.SetDefault("database.erp_data.database", "ERP_DATA")

	// Database - StaffIDMask (SQL Server)
	v.SetDefault("database.staff_id_mask.host", "localhost")
	v.SetDefault("database.staff_id_mask.port", 1433)
	v.SetDefault("database.staff_id_mask.database", "CBN_LAPTOP_REG")

	// Database - EmailService (SQL Server)
	v.SetDefault("database.email_service.host", "localhost")
	v.SetDefault("database.email_service.port", 1433)
	v.SetDefault("database.email_service.database", "XXCBN_EMAIL_SERVICE")

	// Database - SAS (SQL Server)
	v.SetDefault("database.sas.host", "localhost")
	v.SetDefault("database.sas.port", 1433)
	v.SetDefault("database.sas.database", "SAS")

	// JWT
	v.SetDefault("jwt.issuer", "https://pms-api.local")
	v.SetDefault("jwt.audience", "pms-api")
	v.SetDefault("jwt.token_expiry_minutes", "20m")
	v.SetDefault("jwt.refresh_token_expiry", "7d")

	// Active Directory
	v.SetDefault("active_directory.ldap_url", "ldap://localhost:389")
	v.SetDefault("active_directory.domain", "CENBANK")

	// Email
	v.SetDefault("email.smtp_port", 25)

	// CORS
	v.SetDefault("cors.allow_all", true)
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"})
	v.SetDefault("cors.allowed_headers", []string{"*"})

	// Logging
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.file_path", "logs/pms-api.log")
	v.SetDefault("logging.max_size_mb", 100)
	v.SetDefault("logging.max_backups", 10)
	v.SetDefault("logging.max_age_days", 10)
	v.SetDefault("logging.compress", true)
	v.SetDefault("logging.console_enabled", true)
	v.SetDefault("logging.file_enabled", true)

	// Redis
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Jobs
	v.SetDefault("jobs.worker_pool_size", 5)
	v.SetDefault("jobs.worker_queue_size", 100)
	v.SetDefault("jobs.mail_sender_interval", "30s")
	v.SetDefault("jobs.cron_schedule", "@every 10m")

	// Hangfire
	v.SetDefault("hangfire_schema", "WebAPiHangfire")

	// ReCaptcha
	v.SetDefault("recaptcha.site_key", "")
	v.SetDefault("recaptcha.secret_key", "")
	v.SetDefault("recaptcha.verify_url", "https://www.google.com/recaptcha/api/siteverify")

	// Password Generation
	v.SetDefault("password_gen.max_identical_consecutive_chars", 2)
	v.SetDefault("password_gen.lowercase_chars", "abcdefghijklmnopqrstuvwxyz")
	v.SetDefault("password_gen.uppercase_chars", "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	v.SetDefault("password_gen.numeric_chars", "0123456789")
	v.SetDefault("password_gen.special_chars", "!@#$%^&*()-_=+[]{}|;:,.<>?")
	v.SetDefault("password_gen.space_char", " ")
	v.SetDefault("password_gen.password_length_min", 8)
	v.SetDefault("password_gen.password_length_max", 128)

	// General
	v.SetDefault("general.domain", "")
	v.SetDefault("general.name", "PMS")
	v.SetDefault("general.logo", "")
	v.SetDefault("general.no_reply_email", "")

	// RSA
	v.SetDefault("rsa.base_url", "")
	v.SetDefault("rsa.api_key", "")

	// Storage
	v.SetDefault("storage.base_path", "./uploads")
}
