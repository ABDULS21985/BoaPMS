package helpers

import "testing"

// ---------------------------------------------------------------------------
// PluralizeDay
// ---------------------------------------------------------------------------

func TestPluralizeDay_Zero(t *testing.T) {
	got := PluralizeDay(0)
	if got != "Day" {
		t.Errorf("PluralizeDay(0) = %q; want %q", got, "Day")
	}
}

func TestPluralizeDay_One(t *testing.T) {
	got := PluralizeDay(1)
	if got != "Days" {
		t.Errorf("PluralizeDay(1) = %q; want %q", got, "Days")
	}
}

func TestPluralizeDay_Five(t *testing.T) {
	got := PluralizeDay(5)
	if got != "Days" {
		t.Errorf("PluralizeDay(5) = %q; want %q", got, "Days")
	}
}

func TestPluralizeDay_Negative(t *testing.T) {
	got := PluralizeDay(-1)
	if got != "Day" {
		t.Errorf("PluralizeDay(-1) = %q; want %q", got, "Day")
	}
}

// ---------------------------------------------------------------------------
// IntToMonthName
// ---------------------------------------------------------------------------

func TestIntToMonthName_January(t *testing.T) {
	got := IntToMonthName(1)
	if got != "January" {
		t.Errorf("IntToMonthName(1) = %q; want %q", got, "January")
	}
}

func TestIntToMonthName_December(t *testing.T) {
	got := IntToMonthName(12)
	if got != "December" {
		t.Errorf("IntToMonthName(12) = %q; want %q", got, "December")
	}
}

func TestIntToMonthName_Zero(t *testing.T) {
	got := IntToMonthName(0)
	if got != "" {
		t.Errorf("IntToMonthName(0) = %q; want %q", got, "")
	}
}

func TestIntToMonthName_Thirteen(t *testing.T) {
	got := IntToMonthName(13)
	if got != "" {
		t.Errorf("IntToMonthName(13) = %q; want %q", got, "")
	}
}

// ---------------------------------------------------------------------------
// ConvertToKobo
// ---------------------------------------------------------------------------

func TestConvertToKobo(t *testing.T) {
	got := ConvertToKobo(100)
	if got != 10000 {
		t.Errorf("ConvertToKobo(100) = %d; want %d", got, 10000)
	}
}

func TestConvertToKobo_Zero(t *testing.T) {
	got := ConvertToKobo(0)
	if got != 0 {
		t.Errorf("ConvertToKobo(0) = %d; want %d", got, 0)
	}
}

// ---------------------------------------------------------------------------
// ConvertToNaira
// ---------------------------------------------------------------------------

func TestConvertToNaira(t *testing.T) {
	got := ConvertToNaira(10000)
	if got != 100 {
		t.Errorf("ConvertToNaira(10000) = %d; want %d", got, 100)
	}
}

func TestConvertToNaira_Zero(t *testing.T) {
	got := ConvertToNaira(0)
	if got != 0 {
		t.Errorf("ConvertToNaira(0) = %d; want %d", got, 0)
	}
}

// ---------------------------------------------------------------------------
// PadToSixDigits
// ---------------------------------------------------------------------------

func TestPadToSixDigits_One(t *testing.T) {
	got := PadToSixDigits(1)
	if got != "000001" {
		t.Errorf("PadToSixDigits(1) = %q; want %q", got, "000001")
	}
}

func TestPadToSixDigits_ThreeDigits(t *testing.T) {
	got := PadToSixDigits(123)
	if got != "000123" {
		t.Errorf("PadToSixDigits(123) = %q; want %q", got, "000123")
	}
}

func TestPadToSixDigits_OverflowSevenDigits(t *testing.T) {
	got := PadToSixDigits(1234567)
	if got != "1234567" {
		t.Errorf("PadToSixDigits(1234567) = %q; want %q", got, "1234567")
	}
}

func TestPadToSixDigits_ExactlySix(t *testing.T) {
	got := PadToSixDigits(123456)
	if got != "123456" {
		t.Errorf("PadToSixDigits(123456) = %q; want %q", got, "123456")
	}
}
