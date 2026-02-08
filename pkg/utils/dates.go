package utils

import "time"

// ToLocalFormat formats a time.Time value as "02 Jan 2006" (dd MMM yyyy).
// This is the standard local date display format used throughout the PMS.
func ToLocalFormat(t time.Time) string {
	return t.Format("02 Jan 2006")
}

// ToShortLocalFormat formats a time.Time value as "02 01 2006" (dd MM yyyy).
func ToShortLocalFormat(t time.Time) string {
	return t.Format("02 01 2006")
}

// ToLocalFormatPtr formats a *time.Time using ToLocalFormat. If the pointer is
// nil, it returns "Not Specified".
func ToLocalFormatPtr(t *time.Time) string {
	if t == nil {
		return "Not Specified"
	}
	return ToLocalFormat(*t)
}
