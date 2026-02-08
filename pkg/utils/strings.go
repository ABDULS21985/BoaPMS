package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TakeFirstLetter returns the first character of a string, or an empty string
// if the input is empty. It is safe for multi-byte (UTF-8) strings.
func TakeFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(s)
	return string(r)
}

// ToLowerTrimmed converts a string to lowercase and trims leading/trailing
// whitespace. Equivalent to the C# extension: value.Trim().ToLower().
func ToLowerTrimmed(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ToUpperTrimmed converts a string to uppercase and trims leading/trailing
// whitespace. Equivalent to the C# extension: value.Trim().ToUpper().
func ToUpperTrimmed(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// ToTitleCase converts a string to title case where each word is capitalized.
// For example, "hello world" becomes "Hello World".
func ToTitleCase(s string) string {
	return strings.Title(strings.ToLower(strings.TrimSpace(s))) //nolint:staticcheck // strings.Title is acceptable here
}

// ToSentenceCase capitalizes the first letter of a string and lowercases the
// rest. For example, "hELLO WORLD" becomes "Hello world".
func ToSentenceCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + strings.ToLower(s[size:])
}

// ToLocalCurrency formats a float64 amount as Nigerian Naira with thousands
// separators and two decimal places. For example, 1234567.89 becomes
// "₦1,234,567.89".
func ToLocalCurrency(amount float64) string {
	// Format with two decimal places.
	raw := fmt.Sprintf("%.2f", amount)

	// Split into integer and decimal parts.
	parts := strings.SplitN(raw, ".", 2)
	intPart := parts[0]
	decPart := parts[1]

	// Handle negative numbers.
	negative := false
	if strings.HasPrefix(intPart, "-") {
		negative = true
		intPart = intPart[1:]
	}

	// Insert thousands separators from right to left.
	var result []byte
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(digit))
	}

	formatted := string(result) + "." + decPart
	if negative {
		formatted = "-" + formatted
	}

	return "\u20A6" + formatted // ₦ prefix
}

// Shuffle randomly reorders a slice in place using the Fisher-Yates algorithm
// and returns it. The original slice is modified.
func Shuffle[T any](slice []T) []T {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
