package service

// ---------------------------------------------------------------------------
// RSA Adaptive Authentication API models.
// These mirror the .NET RSAVms/ request/response types for the
// RSA Adaptive Authentication REST API.
// ---------------------------------------------------------------------------

// RSAContext is used in responses to correlate authentication attempts.
type RSAContext struct {
	AuthnAttemptID string `json:"authnAttemptId,omitempty"`
	MessageID      string `json:"messageId,omitempty"`
	InResponseTo   string `json:"inResponseTo,omitempty"`
}

// RSAInitializeContext is used in initialize requests.
type RSAInitializeContext struct {
	MessageID string `json:"messageId"`
}

// RSAInitializeRequest is sent to the /initialize endpoint.
type RSAInitializeRequest struct {
	AuthnAttemptTimeout int                  `json:"authnAttemptTimeout"`
	SubjectName         string               `json:"subjectName"`
	Lang                string               `json:"lang"`
	AuthMethodID        string               `json:"authMethodId"`
	Context             RSAInitializeContext  `json:"context"`
}

// RSAInitializeResponse is returned from the /initialize endpoint.
type RSAInitializeResponse struct {
	Context                    RSAContext                  `json:"context,omitempty"`
	CredentialValidationResults []interface{}              `json:"credentialValidationResults,omitempty"`
	AttemptResponseCode        string                     `json:"attemptResponseCode,omitempty"`
	AttemptReasonCode          string                     `json:"attemptReasonCode,omitempty"`
	ChallengeMethods           *ChallengeMethods          `json:"challengeMethods,omitempty"`
}

// RSAVerifyRequest is sent to the /verify endpoint.
type RSAVerifyRequest struct {
	SubjectCredentials []SubjectCredential `json:"subjectCredentials"`
	Context            RSAContext          `json:"context"`
}

// RSAVerifyResponse is returned from the /verify endpoint.
type RSAVerifyResponse struct {
	Context                    RSAContext                     `json:"context,omitempty"`
	CredentialValidationResults []CredentialValidationResult  `json:"credentialValidationResults,omitempty"`
	AttemptResponseCode        string                        `json:"attemptResponseCode,omitempty"`
	AttemptReasonCode          string                        `json:"attemptReasonCode,omitempty"`
	ChallengeMethods           *ChallengeMethods             `json:"challengeMethods,omitempty"`
}

// SubjectCredential represents a credential submitted for verification.
type SubjectCredential struct {
	MethodID        string           `json:"methodId"`
	VersionID       string           `json:"versionId,omitempty"`
	CollectedInputs []CollectedInput `json:"collectedInputs,omitempty"`
}

// CollectedInput is a name-value pair of user-provided credential data.
type CollectedInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ChallengeMethods describes available authentication challenges.
type ChallengeMethods struct {
	Challenges []Challenge `json:"challenges,omitempty"`
}

// Challenge is a set of required authentication methods.
type Challenge struct {
	MethodSetID     string           `json:"methodSetId,omitempty"`
	RequiredMethods []RequiredMethod `json:"requiredMethods,omitempty"`
}

// RequiredMethod describes a single authentication method.
type RequiredMethod struct {
	MethodID    string          `json:"methodId"`
	DisplayName string          `json:"displayName,omitempty"`
	Priority    int             `json:"priority,omitempty"`
	Versions    []MethodVersion `json:"versions,omitempty"`
}

// MethodVersion describes a version of an authentication method.
type MethodVersion struct {
	VersionID        string            `json:"versionId"`
	MethodAttributes []MethodAttribute `json:"methodAttributes,omitempty"`
	ValueRequired    bool              `json:"valueRequired,omitempty"`
	ReferenceID      string            `json:"referenceId,omitempty"`
	Prompt           *MethodPrompt     `json:"prompt,omitempty"`
}

// MethodAttribute is a key-value attribute of a method version.
type MethodAttribute struct {
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
	DataType string `json:"dataType,omitempty"`
}

// MethodPrompt describes UI prompt configuration for credential collection.
type MethodPrompt struct {
	PromptResourceID    string        `json:"promptResourceId,omitempty"`
	DefaultText         string        `json:"defaultText,omitempty"`
	FormatRegex         string        `json:"formatRegex,omitempty"`
	DefaultValue        string        `json:"defaultValue,omitempty"`
	ValueBeingDefined   bool          `json:"valueBeingDefined,omitempty"`
	Sensitive           bool          `json:"sensitive,omitempty"`
	MinLength           int           `json:"minLength,omitempty"`
	MaxLength           int           `json:"maxLength,omitempty"`
	PromptArgs          []interface{} `json:"promptArgs,omitempty"`
	SubjectNameRequired bool          `json:"subjectNameRequired,omitempty"`
}

// CredentialValidationResult is returned after credential verification.
type CredentialValidationResult struct {
	MethodID           string            `json:"methodId"`
	MethodResponseCode string            `json:"methodResponseCode,omitempty"`
	MethodReasonCode   string            `json:"methodReasonCode,omitempty"`
	AuthnAttributes    []MethodAttribute `json:"authnAttributes,omitempty"`
}
