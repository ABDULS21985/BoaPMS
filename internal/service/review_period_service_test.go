package service

import (
	"testing"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// getStartOrEndDate — Quarterly
// ---------------------------------------------------------------------------

func TestGetStartOrEndDate_QuarterlyQ1Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeQuarterly, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q1 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ1End(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeQuarterly, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// month=1, +2=3 (March), last day=31
	want := time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q1 end = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ2Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 2, enums.ReviewPeriodRangeQuarterly, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q2 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ2End(t *testing.T) {
	got, err := getStartOrEndDate(2024, 2, enums.ReviewPeriodRangeQuarterly, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// month=4, +2=6 (June), last day=30
	want := time.Date(2024, time.June, 30, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q2 end = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ3Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 3, enums.ReviewPeriodRangeQuarterly, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q3 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ4Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 4, enums.ReviewPeriodRangeQuarterly, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.October, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q4 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_QuarterlyQ4End(t *testing.T) {
	got, err := getStartOrEndDate(2024, 4, enums.ReviewPeriodRangeQuarterly, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// month=10, +2=12 (December), last day=31
	want := time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Q4 end = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// getStartOrEndDate — BiAnnual
// ---------------------------------------------------------------------------

func TestGetStartOrEndDate_BiAnnualH1Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeBiAnnual, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("H1 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_BiAnnualH1End(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeBiAnnual, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// BiAnnual H1: month = (1-1)*6+1 = 1, +2 = 3 (March), last day = 31
	want := time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("H1 end = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_BiAnnualH2Start(t *testing.T) {
	got, err := getStartOrEndDate(2024, 2, enums.ReviewPeriodRangeBiAnnual, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// BiAnnual H2: month = (2-1)*6+1 = 7
	want := time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("H2 start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_BiAnnualH2End(t *testing.T) {
	got, err := getStartOrEndDate(2024, 2, enums.ReviewPeriodRangeBiAnnual, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// BiAnnual H2: month = 7, +2 = 9 (September), last day = 30
	want := time.Date(2024, time.September, 30, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("H2 end = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// getStartOrEndDate — Annual
// ---------------------------------------------------------------------------

func TestGetStartOrEndDate_AnnualStart(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeAnnual, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Annual start = %v, want %v", got, want)
	}
}

func TestGetStartOrEndDate_AnnualEnd(t *testing.T) {
	got, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRangeAnnual, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Annual: month = 1, +2 = 3 (March), last day = 31
	want := time.Date(2024, time.March, 31, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Annual end = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// getStartOrEndDate — Invalid values
// ---------------------------------------------------------------------------

func TestGetStartOrEndDate_InvalidQuarterValue(t *testing.T) {
	_, err := getStartOrEndDate(2024, 5, enums.ReviewPeriodRangeQuarterly, true)
	if err == nil {
		t.Fatal("expected error for quarter value=5")
	}
}

func TestGetStartOrEndDate_InvalidBiAnnualValue(t *testing.T) {
	_, err := getStartOrEndDate(2024, 3, enums.ReviewPeriodRangeBiAnnual, true)
	if err == nil {
		t.Fatal("expected error for bi-annual value=3")
	}
}

func TestGetStartOrEndDate_InvalidAnnualValue(t *testing.T) {
	_, err := getStartOrEndDate(2024, 2, enums.ReviewPeriodRangeAnnual, true)
	if err == nil {
		t.Fatal("expected error for annual value=2")
	}
}

func TestGetStartOrEndDate_InvalidRangeType(t *testing.T) {
	_, err := getStartOrEndDate(2024, 1, enums.ReviewPeriodRange(99), true)
	if err == nil {
		t.Fatal("expected error for invalid range type")
	}
}

func TestGetStartOrEndDate_ZeroValueDefaultsToOne(t *testing.T) {
	// When value=0, it defaults to 1.
	got, err := getStartOrEndDate(2024, 0, enums.ReviewPeriodRangeQuarterly, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("zero-value default = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// daysInMonth
// ---------------------------------------------------------------------------

func TestDaysInMonth_February_LeapYear(t *testing.T) {
	got := daysInMonth(2024, 2)
	if got != 29 {
		t.Errorf("daysInMonth(2024, 2) = %d, want 29", got)
	}
}

func TestDaysInMonth_February_NonLeapYear(t *testing.T) {
	got := daysInMonth(2023, 2)
	if got != 28 {
		t.Errorf("daysInMonth(2023, 2) = %d, want 28", got)
	}
}

func TestDaysInMonth_January(t *testing.T) {
	got := daysInMonth(2024, 1)
	if got != 31 {
		t.Errorf("daysInMonth(2024, 1) = %d, want 31", got)
	}
}

func TestDaysInMonth_April(t *testing.T) {
	got := daysInMonth(2024, 4)
	if got != 30 {
		t.Errorf("daysInMonth(2024, 4) = %d, want 30", got)
	}
}

func TestDaysInMonth_AllMonths(t *testing.T) {
	// Table-driven test for every month of 2024.
	expected := map[int]int{
		1: 31, 2: 29, 3: 31, 4: 30, 5: 31, 6: 30,
		7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31,
	}

	for month, want := range expected {
		t.Run(time.Month(month).String(), func(t *testing.T) {
			got := daysInMonth(2024, month)
			if got != want {
				t.Errorf("daysInMonth(2024, %d) = %d, want %d", month, got, want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// padWithZeros
// ---------------------------------------------------------------------------

func TestPadWithZeros_Normal(t *testing.T) {
	got, err := padWithZeros("42", 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "000042" {
		t.Errorf("padWithZeros(\"42\", 6) = %q, want %q", got, "000042")
	}
}

func TestPadWithZeros_ExactLength(t *testing.T) {
	got, err := padWithZeros("123456", 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "123456" {
		t.Errorf("padWithZeros(\"123456\", 6) = %q, want %q", got, "123456")
	}
}

func TestPadWithZeros_ExceedsLength(t *testing.T) {
	_, err := padWithZeros("1234567", 6)
	if err == nil {
		t.Fatal("expected error when input exceeds max length")
	}
}

func TestPadWithZeros_SingleChar(t *testing.T) {
	got, err := padWithZeros("1", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0001" {
		t.Errorf("padWithZeros(\"1\", 4) = %q, want %q", got, "0001")
	}
}

func TestPadWithZeros_EmptyString(t *testing.T) {
	got, err := padWithZeros("", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "000" {
		t.Errorf("padWithZeros(\"\", 3) = %q, want %q", got, "000")
	}
}

// ---------------------------------------------------------------------------
// isValidSettingType
// ---------------------------------------------------------------------------

func TestIsValidSettingType_Valid(t *testing.T) {
	validTypes := []string{"Bool", "DateTime", "Decimal", "Double", "Float", "Int", "Long", "String"}

	for _, st := range validTypes {
		t.Run(st, func(t *testing.T) {
			if !isValidSettingType(st) {
				t.Errorf("isValidSettingType(%q) = false, want true", st)
			}
		})
	}
}

func TestIsValidSettingType_Invalid(t *testing.T) {
	invalidTypes := []string{"Map", "Array", "", "bool", "string", "int", "LIST", "Object"}

	for _, st := range invalidTypes {
		name := st
		if name == "" {
			name = "empty_string"
		}
		t.Run(name, func(t *testing.T) {
			if isValidSettingType(st) {
				t.Errorf("isValidSettingType(%q) = true, want false", st)
			}
		})
	}
}

func TestIsValidSettingType_CaseSensitive(t *testing.T) {
	// The valid types are capitalized. Lowercase variants should fail.
	if isValidSettingType("bool") {
		t.Error("expected 'bool' (lowercase) to be invalid")
	}
	if isValidSettingType("STRING") {
		t.Error("expected 'STRING' (uppercase) to be invalid")
	}
	if isValidSettingType("int") {
		t.Error("expected 'int' (lowercase) to be invalid")
	}
}

// ---------------------------------------------------------------------------
// getStartOrEndDate — table-driven comprehensive test
// ---------------------------------------------------------------------------

func TestGetStartOrEndDate_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		year    int
		value   int
		rng     enums.ReviewPeriodRange
		isStart bool
		want    time.Time
		wantErr bool
	}{
		// Quarterly starts
		{"Q1 start", 2024, 1, enums.ReviewPeriodRangeQuarterly, true,
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"Q2 start", 2024, 2, enums.ReviewPeriodRangeQuarterly, true,
			time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), false},
		{"Q3 start", 2024, 3, enums.ReviewPeriodRangeQuarterly, true,
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), false},
		{"Q4 start", 2024, 4, enums.ReviewPeriodRangeQuarterly, true,
			time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), false},

		// Quarterly ends
		{"Q1 end", 2024, 1, enums.ReviewPeriodRangeQuarterly, false,
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), false},
		{"Q2 end", 2024, 2, enums.ReviewPeriodRangeQuarterly, false,
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC), false},
		{"Q3 end", 2024, 3, enums.ReviewPeriodRangeQuarterly, false,
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), false},
		{"Q4 end", 2024, 4, enums.ReviewPeriodRangeQuarterly, false,
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), false},

		// Invalid quarter
		{"Q5 invalid", 2024, 5, enums.ReviewPeriodRangeQuarterly, true,
			time.Time{}, true},

		// BiAnnual
		{"H1 start", 2024, 1, enums.ReviewPeriodRangeBiAnnual, true,
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"H2 start", 2024, 2, enums.ReviewPeriodRangeBiAnnual, true,
			time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), false},
		{"H1 end", 2024, 1, enums.ReviewPeriodRangeBiAnnual, false,
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), false},
		{"H2 end", 2024, 2, enums.ReviewPeriodRangeBiAnnual, false,
			time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), false},

		// Invalid bi-annual
		{"H3 invalid", 2024, 3, enums.ReviewPeriodRangeBiAnnual, true,
			time.Time{}, true},

		// Annual
		{"Annual start", 2024, 1, enums.ReviewPeriodRangeAnnual, true,
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"Annual end", 2024, 1, enums.ReviewPeriodRangeAnnual, false,
			time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), false},

		// Invalid annual
		{"Annual v=2 invalid", 2024, 2, enums.ReviewPeriodRangeAnnual, true,
			time.Time{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getStartOrEndDate(tc.year, tc.value, tc.rng, tc.isStart)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
