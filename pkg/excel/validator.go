package excel

import "strings"

// CheckTrueFalse validates string values that represent boolean TRUE/YES.
// Returns true for "TRUE", "YES", "true", "yes", and any mixed-case variants.
func CheckTrueFalse(value string) bool {
	v := strings.TrimSpace(strings.ToUpper(value))
	return v == "TRUE" || v == "YES"
}

// ValidateYesNo validates a single YES/NO value.
// Returns true only when the trimmed, case-insensitive value is "YES" or "NO".
func ValidateYesNo(value string) bool {
	v := strings.TrimSpace(strings.ToUpper(value))
	return v == "YES" || v == "NO"
}
