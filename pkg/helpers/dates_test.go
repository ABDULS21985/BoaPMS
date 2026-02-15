package helpers

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// ToLocalFormat
// ---------------------------------------------------------------------------

func TestToLocalFormat(t *testing.T) {
	dt := time.Date(2025, time.February, 25, 0, 0, 0, 0, time.UTC)
	got := ToLocalFormat(dt)
	want := "25 Feb 2025"
	if got != want {
		t.Errorf("ToLocalFormat(%v) = %q; want %q", dt, got, want)
	}
}

func TestToLocalFormat_January(t *testing.T) {
	dt := time.Date(2024, time.January, 1, 12, 30, 0, 0, time.UTC)
	got := ToLocalFormat(dt)
	want := "01 Jan 2024"
	if got != want {
		t.Errorf("ToLocalFormat(%v) = %q; want %q", dt, got, want)
	}
}

// ---------------------------------------------------------------------------
// ToShortLocalFormat
// ---------------------------------------------------------------------------

func TestToShortLocalFormat(t *testing.T) {
	dt := time.Date(2025, time.February, 25, 0, 0, 0, 0, time.UTC)
	got := ToShortLocalFormat(dt)
	want := "25 02 2025"
	if got != want {
		t.Errorf("ToShortLocalFormat(%v) = %q; want %q", dt, got, want)
	}
}

func TestToShortLocalFormat_December(t *testing.T) {
	dt := time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC)
	got := ToShortLocalFormat(dt)
	want := "31 12 2023"
	if got != want {
		t.Errorf("ToShortLocalFormat(%v) = %q; want %q", dt, got, want)
	}
}

// ---------------------------------------------------------------------------
// ToLocalFormatPtr
// ---------------------------------------------------------------------------

func TestToLocalFormatPtr_Nil(t *testing.T) {
	got := ToLocalFormatPtr(nil)
	want := "Not Specified"
	if got != want {
		t.Errorf("ToLocalFormatPtr(nil) = %q; want %q", got, want)
	}
}

func TestToLocalFormatPtr_NonNil(t *testing.T) {
	dt := time.Date(2025, time.February, 25, 0, 0, 0, 0, time.UTC)
	got := ToLocalFormatPtr(&dt)
	// ToLocalFormatPtr delegates to ToShortLocalFormat
	want := "25 02 2025"
	if got != want {
		t.Errorf("ToLocalFormatPtr(&%v) = %q; want %q", dt, got, want)
	}
}
