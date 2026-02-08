package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// JWTService handles JWT token generation and validation.
// Mirrors the .NET JwtMiddleware with GenerateToken, ValidateToken,
// GenerateAccessToken, and GenerateRefreshToken.
type JWTService struct {
	cfg config.JWTConfig
	log zerolog.Logger
}

// NewJWTService creates a new JWT service.
func NewJWTService(cfg config.JWTConfig, log zerolog.Logger) *JWTService {
	return &JWTService{cfg: cfg, log: log}
}

// TokenClaims contains the claims embedded in the JWT.
// Includes: UserId, Email, FullName, Roles, Permissions, OrganizationalUnit.
type TokenClaims struct {
	UserID             string   `json:"user_id"`
	Email              string   `json:"email"`
	Name               string   `json:"name"`
	Roles              []string `json:"roles"`
	Permissions        []string `json:"permissions"`
	OrganizationalUnit string   `json:"organizational_unit"`
}

// GenerateAccessToken creates a signed JWT with user claims.
// Maps to .NET's GenerateAccessToken / GenerateToken methods.
func (s *JWTService) GenerateAccessToken(claims TokenClaims, expiryMinutes int) (string, int64, error) {
	if expiryMinutes <= 0 {
		expiryMinutes = int(s.cfg.TokenExpiryMinutes.Minutes())
	}
	if expiryMinutes <= 0 {
		expiryMinutes = 20 // fallback default
	}

	expiresAt := time.Now().UTC().Add(time.Duration(expiryMinutes) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":                 claims.UserID,
		"email":               claims.Email,
		"name":                claims.Name,
		"roles":               claims.Roles,
		"permissions":         claims.Permissions,
		"organizational_unit": claims.OrganizationalUnit,
		"iss":                 s.cfg.Issuer,
		"aud":                 s.cfg.Audience,
		"iat":                 time.Now().UTC().Unix(),
		"exp":                 expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", 0, fmt.Errorf("signing token: %w", err)
	}

	return tokenString, expiresAt.Unix(), nil
}

// GenerateRefreshToken creates a cryptographically random refresh token.
// Mirrors .NET's GenerateRefreshToken (32 random bytes, base64 encoded).
func (s *JWTService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateToken parses and validates a JWT, returning the extracted claims.
// Mirrors .NET's ValidateToken method.
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(s.cfg.Issuer),
		jwt.WithAudience(s.cfg.Audience),
	)
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	claims := &TokenClaims{}
	if sub, ok := mapClaims["sub"].(string); ok {
		claims.UserID = sub
	}
	if email, ok := mapClaims["email"].(string); ok {
		claims.Email = email
	}
	if name, ok := mapClaims["name"].(string); ok {
		claims.Name = name
	}
	if rolesRaw, ok := mapClaims["roles"].([]interface{}); ok {
		for _, r := range rolesRaw {
			if rs, ok := r.(string); ok {
				claims.Roles = append(claims.Roles, rs)
			}
		}
	}
	if permsRaw, ok := mapClaims["permissions"].([]interface{}); ok {
		for _, p := range permsRaw {
			if ps, ok := p.(string); ok {
				claims.Permissions = append(claims.Permissions, ps)
			}
		}
	}
	if ou, ok := mapClaims["organizational_unit"].(string); ok {
		claims.OrganizationalUnit = ou
	}

	return claims, nil
}
