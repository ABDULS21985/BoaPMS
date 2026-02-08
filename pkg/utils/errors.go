package utils

import (
	"errors"
	"strings"
)

// FullErrorMessage recursively builds a complete error message from a chain of
// wrapped errors. Each error's message is joined with " -> " to show the full
// causal chain. For example, an error wrapping two inner errors might produce:
//
//	"outer error -> middle error -> root cause"
//
// Returns an empty string if err is nil.
func FullErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var messages []string
	for current := err; current != nil; current = errors.Unwrap(current) {
		msg := current.Error()
		// Avoid duplicating the inner message that fmt.Errorf("%w") appends.
		// When wrapping with %w, Go's fmt.Errorf produces messages like
		// "outer: inner", so we extract only the prefix before the inner
		// error's text to avoid repetition.
		inner := errors.Unwrap(current)
		if inner != nil {
			innerMsg := inner.Error()
			if strings.HasSuffix(msg, ": "+innerMsg) {
				msg = strings.TrimSuffix(msg, ": "+innerMsg)
			}
		}
		messages = append(messages, msg)
	}

	return strings.Join(messages, " -> ")
}
