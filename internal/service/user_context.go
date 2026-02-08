package service

import (
	"context"

	"github.com/enterprise-pms/pms-api/internal/middleware"
	"github.com/rs/zerolog"
)

// userContextService extracts current user information from the request context.
// Claims are injected by the JWT middleware.
type userContextService struct {
	log zerolog.Logger
}

func newUserContextService(log zerolog.Logger) UserContextService {
	return &userContextService{log: log}
}

func (s *userContextService) GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		return v
	}
	return ""
}

func (s *userContextService) GetEmail(ctx context.Context) string {
	if v, ok := ctx.Value(middleware.EmailKey).(string); ok {
		return v
	}
	return ""
}

func (s *userContextService) GetRoles(ctx context.Context) []string {
	if v, ok := ctx.Value(middleware.RolesKey).([]string); ok {
		return v
	}
	return nil
}

func (s *userContextService) IsAuthenticated(ctx context.Context) bool {
	return s.GetUserID(ctx) != ""
}

func (s *userContextService) IsInRole(ctx context.Context, role string) bool {
	for _, r := range s.GetRoles(ctx) {
		if r == role {
			return true
		}
	}
	return false
}
