package helpers

import (
	"sort"
	"testing"
)

// ---------------------------------------------------------------------------
// TakeFirstLetter
// ---------------------------------------------------------------------------

func TestTakeFirstLetter_EmptyString(t *testing.T) {
	got := TakeFirstLetter("")
	if got != "" {
		t.Errorf("TakeFirstLetter(%q) = %q; want %q", "", got, "")
	}
}

func TestTakeFirstLetter_SingleChar(t *testing.T) {
	got := TakeFirstLetter("A")
	if got != "A" {
		t.Errorf("TakeFirstLetter(%q) = %q; want %q", "A", got, "A")
	}
}

func TestTakeFirstLetter_MultiChar(t *testing.T) {
	got := TakeFirstLetter("Hello")
	if got != "H" {
		t.Errorf("TakeFirstLetter(%q) = %q; want %q", "Hello", got, "H")
	}
}

func TestTakeFirstLetter_Unicode(t *testing.T) {
	// U+00FC = 'u with diaeresis', a multi-byte UTF-8 character
	got := TakeFirstLetter("\u00fcber")
	want := "\u00fc"
	if got != want {
		t.Errorf("TakeFirstLetter(%q) = %q; want %q", "\u00fcber", got, want)
	}
}

// ---------------------------------------------------------------------------
// ToUpperTrimmed
// ---------------------------------------------------------------------------

func TestToUpperTrimmed_Normal(t *testing.T) {
	got := ToUpperTrimmed("hello")
	if got != "HELLO" {
		t.Errorf("ToUpperTrimmed(%q) = %q; want %q", "hello", got, "HELLO")
	}
}

func TestToUpperTrimmed_WithSpaces(t *testing.T) {
	got := ToUpperTrimmed("  hello world  ")
	if got != "HELLO WORLD" {
		t.Errorf("ToUpperTrimmed(%q) = %q; want %q", "  hello world  ", got, "HELLO WORLD")
	}
}

func TestToUpperTrimmed_Empty(t *testing.T) {
	got := ToUpperTrimmed("")
	if got != "" {
		t.Errorf("ToUpperTrimmed(%q) = %q; want %q", "", got, "")
	}
}

// ---------------------------------------------------------------------------
// ToLowerTrimmed
// ---------------------------------------------------------------------------

func TestToLowerTrimmed_Normal(t *testing.T) {
	got := ToLowerTrimmed("HELLO")
	if got != "hello" {
		t.Errorf("ToLowerTrimmed(%q) = %q; want %q", "HELLO", got, "hello")
	}
}

func TestToLowerTrimmed_WithSpaces(t *testing.T) {
	got := ToLowerTrimmed("  HELLO WORLD  ")
	if got != "hello world" {
		t.Errorf("ToLowerTrimmed(%q) = %q; want %q", "  HELLO WORLD  ", got, "hello world")
	}
}

func TestToLowerTrimmed_Empty(t *testing.T) {
	got := ToLowerTrimmed("")
	if got != "" {
		t.Errorf("ToLowerTrimmed(%q) = %q; want %q", "", got, "")
	}
}

// ---------------------------------------------------------------------------
// ToSentenceCase
// ---------------------------------------------------------------------------

func TestToSentenceCase_Normal(t *testing.T) {
	got := ToSentenceCase("hello world")
	if got != "Hello world" {
		t.Errorf("ToSentenceCase(%q) = %q; want %q", "hello world", got, "Hello world")
	}
}

func TestToSentenceCase_Empty(t *testing.T) {
	got := ToSentenceCase("")
	if got != "" {
		t.Errorf("ToSentenceCase(%q) = %q; want %q", "", got, "")
	}
}

func TestToSentenceCase_SingleChar(t *testing.T) {
	got := ToSentenceCase("a")
	if got != "A" {
		t.Errorf("ToSentenceCase(%q) = %q; want %q", "a", got, "A")
	}
}

func TestToSentenceCase_AllUpper(t *testing.T) {
	got := ToSentenceCase("HELLO")
	if got != "Hello" {
		t.Errorf("ToSentenceCase(%q) = %q; want %q", "HELLO", got, "Hello")
	}
}

// ---------------------------------------------------------------------------
// ToTitleCase
// ---------------------------------------------------------------------------

func TestToTitleCase_Normal(t *testing.T) {
	got := ToTitleCase("hello world")
	if got != "Hello World" {
		t.Errorf("ToTitleCase(%q) = %q; want %q", "hello world", got, "Hello World")
	}
}

func TestToTitleCase_Empty(t *testing.T) {
	got := ToTitleCase("")
	if got != "" {
		t.Errorf("ToTitleCase(%q) = %q; want %q", "", got, "")
	}
}

// ---------------------------------------------------------------------------
// FormatNaira
// ---------------------------------------------------------------------------

func TestFormatNaira_Zero(t *testing.T) {
	got := FormatNaira(0)
	want := "\u20a60.00"
	if got != want {
		t.Errorf("FormatNaira(0) = %q; want %q", got, want)
	}
}

func TestFormatNaira_WithDecimals(t *testing.T) {
	got := FormatNaira(1234.56)
	want := "\u20a61,234.56"
	if got != want {
		t.Errorf("FormatNaira(1234.56) = %q; want %q", got, want)
	}
}

func TestFormatNaira_Negative(t *testing.T) {
	got := FormatNaira(-500)
	want := "-\u20a6500.00"
	if got != want {
		t.Errorf("FormatNaira(-500) = %q; want %q", got, want)
	}
}

func TestFormatNaira_Million(t *testing.T) {
	got := FormatNaira(1000000)
	want := "\u20a61,000,000.00"
	if got != want {
		t.Errorf("FormatNaira(1000000) = %q; want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Shuffle
// ---------------------------------------------------------------------------

func TestShuffle_NonNilResult(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := Shuffle(input)
	if got == nil {
		t.Fatal("Shuffle returned nil for non-empty input")
	}
}

func TestShuffle_SameLength(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := Shuffle(input)
	if len(got) != len(input) {
		t.Errorf("Shuffle length = %d; want %d", len(got), len(input))
	}
}

func TestShuffle_PreservesElements(t *testing.T) {
	input := []int{5, 3, 1, 4, 2}
	got := Shuffle(input)

	sortedInput := make([]int, len(input))
	copy(sortedInput, input)
	sort.Ints(sortedInput)

	sortedGot := make([]int, len(got))
	copy(sortedGot, got)
	sort.Ints(sortedGot)

	for i := range sortedInput {
		if sortedInput[i] != sortedGot[i] {
			t.Errorf("Shuffle changed elements: sorted input %v != sorted output %v", sortedInput, sortedGot)
			break
		}
	}
}

func TestShuffle_EmptySlice(t *testing.T) {
	got := Shuffle([]int{})
	if got != nil {
		t.Errorf("Shuffle(empty) = %v; want nil", got)
	}
}

func TestShuffle_DoesNotModifyOriginal(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	original := make([]int, len(input))
	copy(original, input)

	_ = Shuffle(input)

	for i := range input {
		if input[i] != original[i] {
			t.Errorf("Shuffle modified original slice: got %v; want %v", input, original)
			break
		}
	}
}
