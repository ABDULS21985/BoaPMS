package helpers

import (
	"fmt"
	"time"
)

// PluralizeDay returns "Days" if n > 0, otherwise "Day".
func PluralizeDay(n int) string {
	if n > 0 {
		return "Days"
	}
	return "Day"
}

// IntToMonthName converts an integer (1-12) to a month name.
// Returns empty string for invalid values.
func IntToMonthName(n int) string {
	if n < 1 || n > 12 {
		return ""
	}
	return time.Month(n).String()
}

// ConvertToKobo converts Naira to Kobo (multiply by 100).
func ConvertToKobo(naira int) int {
	return naira * 100
}

// ConvertToNaira converts Kobo to Naira (integer division by 100).
func ConvertToNaira(kobo int) int {
	return kobo / 100
}

// PadToSixDigits zero-pads an integer to 6 digits.
// Examples: 1 → "000001", 123 → "000123", 1234567 → "1234567"
func PadToSixDigits(n int) string {
	return fmt.Sprintf("%06d", n)
}
