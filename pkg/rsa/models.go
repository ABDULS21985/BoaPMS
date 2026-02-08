package rsa

// ---------------------------------------------------------------------------
// Initialize
// ---------------------------------------------------------------------------

// InitializeRequest is the payload sent to the RSA SecurID "initialize"
// endpoint to start an authentication attempt.
type InitializeRequest struct {
	SubjectName         string             `json:"subjectName"`
	AuthnAttemptTimeout int                `json:"authnAttemptTimeout,omitempty"`
	Lang                string             `json:"lang,omitempty"`
	AuthMethodId        string             `json:"authMethodId,omitempty"`
	Context             *InitializeContext `json:"context,omitempty"`
}

// InitializeContext carries the message correlation ID when initializing.
type InitializeContext struct {
	MessageId string `json:"messageId,omitempty"`
}

// NewInitializeRequest returns an InitializeRequest pre-filled with the
// standard defaults used by the PMS application.
func NewInitializeRequest(subjectName, messageID string) *InitializeRequest {
	return &InitializeRequest{
		SubjectName:         subjectName,
		AuthnAttemptTimeout: 300,
		Lang:                "us_EN",
		AuthMethodId:        "TOKEN",
		Context: &InitializeContext{
			MessageId: messageID,
		},
	}
}

// ---------------------------------------------------------------------------
// Verify
// ---------------------------------------------------------------------------

// VerifyRequest is the payload sent to the RSA SecurID "verify" endpoint to
// validate a token code that the user entered.
type VerifyRequest struct {
	SubjectCredentials []SubjectCredential `json:"subjectCredentials"`
	Context            *VerifyContext       `json:"context"`
}

// SubjectCredential describes one authentication credential being submitted.
type SubjectCredential struct {
	MethodId        string           `json:"methodId"`
	VersionId       string           `json:"versionId,omitempty"`
	CollectedInputs []CollectedInput `json:"collectedInputs"`
}

// CollectedInput is a single name/value pair representing user input (e.g. a
// TOKEN code).
type CollectedInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// VerifyContext carries the correlation IDs that tie a verify call back to the
// original initialize attempt.
type VerifyContext struct {
	AuthnAttemptId string `json:"authnAttemptId"`
	MessageId      string `json:"messageId"`
	InResponseTo   string `json:"inResponseTo"`
}

// NewVerifyRequest builds a VerifyRequest for a TOKEN credential.
func NewVerifyRequest(tokenCode, authnAttemptID, messageID, inResponseTo string) *VerifyRequest {
	return &VerifyRequest{
		SubjectCredentials: []SubjectCredential{
			{
				MethodId:  "TOKEN",
				VersionId: "string",
				CollectedInputs: []CollectedInput{
					{Name: "TOKEN", Value: tokenCode},
				},
			},
		},
		Context: &VerifyContext{
			AuthnAttemptId: authnAttemptID,
			MessageId:      messageID,
			InResponseTo:   inResponseTo,
		},
	}
}

// ---------------------------------------------------------------------------
// Shared response types (Initialize and Verify share the same shape)
// ---------------------------------------------------------------------------

// Response is the common response envelope returned by both the initialize
// and verify endpoints.
type Response struct {
	AttemptResponseCode        string                      `json:"attemptResponseCode"`
	AttemptReasonCode          string                      `json:"attemptReasonCode"`
	Context                    *ResponseContext             `json:"context,omitempty"`
	CredentialValidationResults []CredentialValidationResult `json:"credentialValidationResults,omitempty"`
	ChallengeMethods           *ChallengeMethods            `json:"challengeMethods,omitempty"`
}

// ResponseContext contains the IDs that correlate requests and responses
// across the authentication conversation.
type ResponseContext struct {
	AuthnAttemptId string `json:"authnAttemptId"`
	MessageId      string `json:"messageId"`
	InResponseTo   string `json:"inResponseTo"`
}

// CredentialValidationResult holds the per-method outcome of a verify call.
type CredentialValidationResult struct {
	MethodId           string           `json:"methodId"`
	MethodResponseCode string           `json:"methodResponseCode"`
	MethodReasonCode   string           `json:"methodReasonCode"`
	AuthnAttributes    []AuthnAttribute `json:"authnAttributes,omitempty"`
}

// AuthnAttribute is a key/value pair returned inside validation results.
type AuthnAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ---------------------------------------------------------------------------
// Challenge method types
// ---------------------------------------------------------------------------

// ChallengeMethods wraps the list of challenges that the server may present
// after an initialize call.
type ChallengeMethods struct {
	Challenges []Challenge `json:"challenges,omitempty"`
}

// Challenge represents one authentication challenge.
type Challenge struct {
	RequiredMethods []RequiredMethod `json:"requiredMethods,omitempty"`
}

// RequiredMethod describes a method the user must complete.
type RequiredMethod struct {
	MethodId string    `json:"methodId"`
	Versions []Version `json:"versions,omitempty"`
}

// Version describes a specific version of a required method, including the
// prompts the UI should present.
type Version struct {
	VersionId       string            `json:"versionId"`
	Prompts         []Prompt          `json:"prompts,omitempty"`
	MethodAttributes []MethodAttribute `json:"methodAttributes,omitempty"`
}

// Prompt describes one input that the user should be asked to provide.
type Prompt struct {
	PromptInfoType string `json:"promptInfoType"`
	PromptValue    string `json:"promptValue"`
}

// MethodAttribute is a key/value pair attached to a method version.
type MethodAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
