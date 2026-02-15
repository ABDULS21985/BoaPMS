package helpers

import "time"

// ToLocalFormat formats t as "02 Jan 2006" (e.g., "25 Feb 2025").
func ToLocalFormat(t time.Time) string {
	return t.Format("02 Jan 2006")
}

// ToShortLocalFormat formats t as "02 01 2006" (e.g., "25 02 2025").
func ToShortLocalFormat(t time.Time) string {
	return t.Format("02 01 2006")
}

// ToLocalFormatPtr formats a time pointer. Returns "Not Specified" if t is nil.
func ToLocalFormatPtr(t *time.Time) string {
	if t == nil {
		return "Not Specified"
	}
	return ToShortLocalFormat(*t)
}
