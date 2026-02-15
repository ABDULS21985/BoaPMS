package auth

import "time"

// RefreshToken represents a hashed refresh token stored in the database.
// The plaintext token is never persisted; only a SHA-256 hash is stored.
type RefreshToken struct {
	ID        string    `json:"id"         gorm:"column:id;primaryKey;size:450"`
	UserID    string    `json:"user_id"    gorm:"column:user_id;not null;size:450;index"`
	Token     string    `json:"-"          gorm:"column:token;not null;size:64;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Revoked   bool      `json:"revoked"    gorm:"column:revoked;default:false;index"`
}

// TableName returns the fully-qualified PostgreSQL table name including the schema prefix.
func (RefreshToken) TableName() string { return "CoreSchema.refresh_tokens" }

// IsExpired reports whether the refresh token has passed its expiry time.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(rt.ExpiresAt)
}

// IsValid reports whether the refresh token is usable (not revoked, not expired).
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && !rt.IsExpired()
}
