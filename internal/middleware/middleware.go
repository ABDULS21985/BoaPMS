package middleware

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	UserIDKey             contextKey = "user_id"
	EmailKey              contextKey = "email"
	RolesKey              contextKey = "roles"
	PermissionsKey        contextKey = "permissions"
	OrganizationalUnitKey contextKey = "organizational_unit"
)

// Stack holds all middleware instances.
type Stack struct {
	cfg *config.Config
	log zerolog.Logger
}

// New creates a middleware stack.
func New(cfg *config.Config, log zerolog.Logger) *Stack {
	return &Stack{cfg: cfg, log: log}
}

// RequestLogger logs every HTTP request with method, path, status, and duration.
func (s *Stack) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		s.log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", sw.status).
			Dur("duration", time.Since(start)).
			Str("remote", r.RemoteAddr).
			Msg("HTTP request")
	})
}

// Recover catches panics and returns a 500 error.
func (s *Stack) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.Error().Interface("panic", err).
					Str("path", r.URL.Path).
					Msg("Recovered from panic")
				response.Error(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS handles Cross-Origin Resource Sharing headers.
func (s *Stack) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.CORS.AllowAll {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if len(s.cfg.CORS.AllowedOrigins) > 0 {
			origin := r.Header.Get("Origin")
			for _, allowed := range s.cfg.CORS.AllowedOrigins {
				if allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(s.cfg.CORS.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(s.cfg.CORS.AllowedHeaders, ", "))

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// APIKeyAuth validates the api-key header.
// Bypasses /swagger and /health endpoints.
func (s *Stack) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip API key check for certain paths
		path := r.URL.Path
		if strings.HasPrefix(path, "/swagger") ||
			strings.HasPrefix(path, "/health") ||
			strings.HasPrefix(path, "/hangfire") {
			next.ServeHTTP(w, r)
			return
		}

		if s.cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("api-key")
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(s.cfg.APIKey)) != 1 {
			response.Error(w, http.StatusUnauthorized, "Invalid or missing API key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JWTAuth validates JWT bearer tokens and injects user claims into context.
func (s *Stack) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "Missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.cfg.JWT.Secret), nil
		},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer(s.cfg.JWT.Issuer),
			jwt.WithAudience(s.cfg.JWT.Audience),
		)

		if err != nil || !token.Valid {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Inject claims into context
		ctx := r.Context()
		if userID, ok := claims["sub"].(string); ok {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}
		if email, ok := claims["email"].(string); ok {
			ctx = context.WithValue(ctx, EmailKey, email)
		}
		if rolesRaw, ok := claims["roles"]; ok {
			if roles, ok := rolesRaw.([]interface{}); ok {
				roleStrings := make([]string, 0, len(roles))
				for _, r := range roles {
					if s, ok := r.(string); ok {
						roleStrings = append(roleStrings, s)
					}
				}
				ctx = context.WithValue(ctx, RolesKey, roleStrings)
			}
		}
		if permsRaw, ok := claims["permissions"]; ok {
			if perms, ok := permsRaw.([]interface{}); ok {
				permStrings := make([]string, 0, len(perms))
				for _, p := range perms {
					if ps, ok := p.(string); ok {
						permStrings = append(permStrings, ps)
					}
				}
				ctx = context.WithValue(ctx, PermissionsKey, permStrings)
			}
		}
		if ou, ok := claims["organizational_unit"].(string); ok {
			ctx = context.WithValue(ctx, OrganizationalUnitKey, ou)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns middleware that restricts access to users with at least
// one of the specified roles. Must be used after JWTAuth in the middleware chain.
// This mirrors the .NET [Authorize(Roles = "...")] attribute.
func (s *Stack) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value(RolesKey).([]string)
			if !ok || len(userRoles) == 0 {
				response.Error(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			for _, required := range roles {
				for _, userRole := range userRoles {
					if userRole == required {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			response.Error(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}

// RequirePermission returns middleware that checks if the user has the specified
// permission in their JWT claims. Permissions are loaded from the RolePermission
// junction table at login time and embedded in the token.
// This mirrors the .NET Permission/RolePermission fine-grained authorization.
func (s *Stack) RequirePermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userPerms, ok := r.Context().Value(PermissionsKey).([]string)
			if !ok || len(userPerms) == 0 {
				response.Error(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			for _, required := range permissions {
				for _, userPerm := range userPerms {
					if userPerm == required {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			response.Error(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}

// SecurityHeaders adds common security headers.
func (s *Stack) SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// statusWriter wraps http.ResponseWriter to capture the status code.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
