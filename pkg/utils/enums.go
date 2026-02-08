package utils

import (
	"strings"
	"unicode"
)

// HumanizeEnum converts an enum-style name into a human-readable string.
//
// It handles two common naming conventions:
//   - CamelCase / PascalCase: "AppraisalInProgress" becomes "Appraisal in progress"
//   - SCREAMING_SNAKE_CASE: "APPRAISAL_IN_PROGRESS" becomes "Appraisal in progress"
//
// The result always has the first letter capitalized and the remainder in
// lowercase, with spaces inserted at word boundaries.
func HumanizeEnum(name string) string {
	if name == "" {
		return ""
	}

	// Handle SCREAMING_SNAKE_CASE: replace underscores with spaces and sentence-case.
	if strings.Contains(name, "_") {
		words := strings.Split(name, "_")
		for i, w := range words {
			if w == "" {
				continue
			}
			if i == 0 {
				words[i] = capitalizeFirst(strings.ToLower(w))
			} else {
				words[i] = strings.ToLower(w)
			}
		}
		return strings.Join(words, " ")
	}

	// Handle CamelCase / PascalCase: insert spaces before uppercase letters.
	var builder strings.Builder
	runes := []rune(name)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			// Check if this is the start of a new word:
			// - Previous rune is lowercase, OR
			// - Next rune is lowercase (handles sequences like "XMLParser" -> "XML Parser")
			prevLower := unicode.IsLower(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || nextLower {
				builder.WriteRune(' ')
			}
		}

		if i == 0 {
			builder.WriteRune(unicode.ToUpper(r))
		} else {
			builder.WriteRune(unicode.ToLower(r))
		}
	}

	return builder.String()
}

// capitalizeFirst uppercases the first rune of a string.
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
