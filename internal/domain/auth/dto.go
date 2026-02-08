package auth

import "time"

// AuthenticateRequest is the login request payload.
type AuthenticateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthenticateResponse is returned after successful authentication.
type AuthenticateResponse struct {
	UserID             string   `json:"user_id"`
	Username           string   `json:"username"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Email              string   `json:"email"`
	Roles              []string `json:"roles"`
	Permissions        []string `json:"permissions"`
	OrganizationalUnit string   `json:"organizational_unit,omitempty"`
	AccessToken        string   `json:"access_token"`
	RefreshToken       string   `json:"refresh_token"`
	ExpiresAt          int64    `json:"expires_at"`
}

// TokenResponse is returned when refreshing an access token.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// RefreshTokenRequest is the request payload for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ADUser holds the user information retrieved from Active Directory.
type ADUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

// CurrentUserData holds extracted claims from the JWT for request-scoped use.
type CurrentUserData struct {
	UserID             string   `json:"user_id"`
	Username           string   `json:"username"`
	Email              string   `json:"email"`
	Name               string   `json:"name"`
	Roles              []string `json:"roles"`
	Permissions        []string `json:"permissions"`
	OrganizationalUnit string   `json:"organizational_unit,omitempty"`
}

// RegisterUserRequest is the payload for creating a new local user.
type RegisterUserRequest struct {
	Username  string `json:"username"   validate:"required,min=3,max=256"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
}

// UpdateUserRequest is the payload for updating a user.
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"omitempty,email"`
	IsActive  bool   `json:"is_active"`
}

// CookieData holds session cookie information for the authenticated user.
type CookieData struct {
	ID                string    `json:"id"`
	JambRegNo         string    `json:"jambRegNo"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	Name              string    `json:"name"`
	Token             string    `json:"token"`
	Phone             string    `json:"phone"`
	NeedPasswordReset bool      `json:"needPasswordReset"`
	Expiry            time.Time `json:"expiry"`
	ReturnUrl         string    `json:"returnUrl"`
}

// ForgotPasswordRequest is the payload for initiating a password reset.
type ForgotPasswordRequest struct {
	Email      string `json:"email"      validate:"required"`
	ClientHost string `json:"clientHost"`
}

// ResetPasswordRequest is the payload for resetting a user's password.
type ResetPasswordRequest struct {
	Email           string `json:"email"           validate:"required"`
	Password        string `json:"password"        validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword"`
	Code            string `json:"code"`
}

// ConfirmEmailRequest is the payload for confirming a user's email address.
type ConfirmEmailRequest struct {
	UserId string `json:"userId" validate:"required"`
	Code   string `json:"code"   validate:"required"`
}

// ResendConfirmationEmailRequest is the payload for resending a confirmation email.
type ResendConfirmationEmailRequest struct {
	Email      string `json:"email"      validate:"required"`
	ClientHost string `json:"clientHost"`
}

// ActiveDirectoryLoginResponseVm is returned after an Active Directory login attempt.
type ActiveDirectoryLoginResponseVm struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
	ADUser    *ADUser `json:"adUser"`
}
