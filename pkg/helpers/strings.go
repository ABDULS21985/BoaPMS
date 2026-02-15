package helpers

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TakeFirstLetter returns the first character of s as a string.
// Returns empty string if s is empty. Unicode-safe.
func TakeFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(s)
	return string(r)
}

// ToUpperTrimmed converts s to uppercase and trims whitespace.
func ToUpperTrimmed(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// ToLowerTrimmed converts s to lowercase and trims whitespace.
func ToLowerTrimmed(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ToSentenceCase converts s to sentence case (first letter upper, rest lower).
func ToSentenceCase(s string) string {
	if s == "" {
		return ""
	}
	lower := strings.ToLower(s)
	r, size := utf8.DecodeRuneInString(lower)
	if r == utf8.RuneError {
		return lower
	}
	return string(unicode.ToUpper(r)) + lower[size:]
}

// ToTitleCase converts s to title case using English locale rules.
func ToTitleCase(s string) string {
	return cases.Title(language.English).String(s)
}

// FormatNaira formats amount as Nigerian Naira with comma grouping.
// Example: 1234.56 → "₦1,234.56"
func FormatNaira(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}

	// Format with 2 decimal places
	raw := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(raw, ".")

	// Add comma grouping to integer part
	intPart := parts[0]
	var result strings.Builder
	for i, ch := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(ch)
	}

	prefix := "₦"
	if negative {
		prefix = "-₦"
	}
	return prefix + result.String() + "." + parts[1]
}

// Shuffle returns a new slice with elements in random order.
// Uses Fisher-Yates shuffle. The original slice is not modified.
func Shuffle[T any](list []T) []T {
	if len(list) == 0 {
		return nil
	}
	result := make([]T, len(list))
	copy(result, list)
	for i := len(result) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}
