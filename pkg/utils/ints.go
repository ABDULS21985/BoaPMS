package utils

import (
	"fmt"
	"time"
)

// PluralizeDay returns "Days" if the value is greater than 1, or "Day"
// otherwise. Useful for display strings like "3 Days" vs "1 Day".
func PluralizeDay(value int) string {
	if value > 1 {
		return "Days"
	}
	return "Day"
}

// ConvertToMonth converts an integer (1-12) to the full English month name.
// For example, 1 returns "January", 12 returns "December".
// Returns an empty string for values outside the 1-12 range.
func ConvertToMonth(value int) string {
	if value < 1 || value > 12 {
		return ""
	}
	return time.Month(value).String()
}

// ConvertToKobo multiplies a Naira value by 100 to convert it to Kobo.
func ConvertToKobo(value int) int {
	return value * 100
}

// ConvertToNaira divides a Kobo value by 100 to convert it to Naira.
// Integer division; any remainder is truncated.
func ConvertToNaira(value int) int {
	return value / 100
}

// PadToSixDigits formats a number as a zero-padded 6-digit string.
// For example, 42 becomes "000042".
func PadToSixDigits(number int) string {
	return fmt.Sprintf("%06d", number)
}
