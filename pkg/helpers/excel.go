package helpers

import "strings"

// ConvertToBool converts common Excel boolean strings to bool.
// Returns true for "YES" or "TRUE" (case-insensitive).
func ConvertToBool(s string) bool {
	upper := strings.ToUpper(strings.TrimSpace(s))
	return upper == "YES" || upper == "TRUE"
}

// ValidateYesNo checks if a string represents a boolean-truthy Excel value.
func ValidateYesNo(s string) bool {
	return ConvertToBool(s)
}
