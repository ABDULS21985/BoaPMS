package service

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
)

// activeDirectoryService handles LDAP/AD authentication and user lookups.
// Mirrors .NET's ActiveDirectoryService with Bind + Search pattern.
type activeDirectoryService struct {
	cfg config.ActiveDirectoryConfig
	log zerolog.Logger
}

func newActiveDirectoryService(cfg config.ActiveDirectoryConfig, log zerolog.Logger) ActiveDirectoryService {
	return &activeDirectoryService{cfg: cfg, log: log}
}

// Authenticate verifies credentials against Active Directory using LDAP bind.
func (s *activeDirectoryService) Authenticate(username, password string) (bool, error) {
	conn, err := s.connect()
	if err != nil {
		return false, fmt.Errorf("connecting to LDAP: %w", err)
	}
	defer conn.Close()

	// Bind with the user's credentials (DOMAIN\username format)
	bindDN := fmt.Sprintf("%s\\%s", s.cfg.Domain, username)
	if err := conn.Bind(bindDN, password); err != nil {
		s.log.Debug().Str("user", username).Msg("AD authentication failed")
		return false, nil
	}

	s.log.Debug().Str("user", username).Msg("AD authentication succeeded")
	return true, nil
}

// GetUser retrieves user details from Active Directory via LDAP search.
func (s *activeDirectoryService) GetUser(username string) (interface{}, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to LDAP: %w", err)
	}
	defer conn.Close()

	// Bind with service account for search
	if s.cfg.Username != "" && s.cfg.Password != "" {
		bindDN := fmt.Sprintf("%s\\%s", s.cfg.Domain, s.cfg.Username)
		if err := conn.Bind(bindDN, s.cfg.Password); err != nil {
			return nil, fmt.Errorf("service account bind failed: %w", err)
		}
	}

	// Build search base from domain (e.g., CENBANK → DC=CENBANK,DC=local)
	baseDN := domainToBaseDN(s.cfg.Domain)
	filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", ldap.EscapeFilter(username))

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
		filter,
		[]string{"sAMAccountName", "mail", "givenName", "sn", "displayName"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %w", err)
	}
	if len(result.Entries) == 0 {
		return nil, nil
	}

	entry := result.Entries[0]
	adUser := &auth.ADUser{
		Username:    entry.GetAttributeValue("sAMAccountName"),
		Email:       entry.GetAttributeValue("mail"),
		FirstName:   entry.GetAttributeValue("givenName"),
		LastName:    entry.GetAttributeValue("sn"),
		DisplayName: entry.GetAttributeValue("displayName"),
	}

	return adUser, nil
}

func (s *activeDirectoryService) connect() (*ldap.Conn, error) {
	if strings.HasPrefix(s.cfg.LdapURL, "ldaps://") {
		return ldap.DialURL(s.cfg.LdapURL, ldap.DialWithTLSConfig(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // internal AD server
		}))
	}
	return ldap.DialURL(s.cfg.LdapURL)
}

// domainToBaseDN converts a flat domain name to an LDAP base DN.
// e.g., "CENBANK" → "DC=CENBANK,DC=local"
func domainToBaseDN(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) == 1 {
		// Simple domain like "CENBANK" — append DC=local
		return fmt.Sprintf("DC=%s,DC=local", domain)
	}
	var dcParts []string
	for _, p := range parts {
		dcParts = append(dcParts, "DC="+p)
	}
	return strings.Join(dcParts, ",")
}
