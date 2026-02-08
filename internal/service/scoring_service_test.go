package service

import (
	"errors"
	"testing"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// newTestScoringService creates a scoringService with a no-op logger.
func newTestScoringService() *scoringService {
	return newScoringService(zerolog.Nop())
}

// ---------------------------------------------------------------------------
// DetermineGrade
// ---------------------------------------------------------------------------

func TestDetermineGrade(t *testing.T) {
	svc := newTestScoringService()

	tests := []struct {
		name  string
		score float64
		want  enums.PerformanceGrade
	}{
		// Probation: < 30
		{"Probation_29.99", 29.99, enums.PerformanceGradeProbation},
		{"EdgeCase_Zero", 0, enums.PerformanceGradeProbation},
		{"EdgeCase_Negative", -5, enums.PerformanceGradeProbation},

		// Developing: 30 <= score < 50
		{"Developing_30", 30, enums.PerformanceGradeDeveloping},
		{"Developing_49.99", 49.99, enums.PerformanceGradeDeveloping},

		// Progressive: 50 <= score < 66
		{"Progressive_50", 50, enums.PerformanceGradeProgressive},
		{"Progressive_65.99", 65.99, enums.PerformanceGradeProgressive},

		// Competent: 66 <= score < 80
		{"Competent_66", 66, enums.PerformanceGradeCompetent},
		{"Competent_79.99", 79.99, enums.PerformanceGradeCompetent},

		// Accomplished: 80 <= score < 90
		{"Accomplished_80", 80, enums.PerformanceGradeAccomplished},
		{"Accomplished_89.99", 89.99, enums.PerformanceGradeAccomplished},

		// Exemplary: >= 90
		{"Exemplary_90", 90, enums.PerformanceGradeExemplary},
		{"Exemplary_100", 100, enums.PerformanceGradeExemplary},
		{"EdgeCase_Above100", 105, enums.PerformanceGradeExemplary},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.DetermineGrade(decimal.NewFromFloat(tc.score))
			if got != tc.want {
				t.Errorf("DetermineGrade(%v) = %s, want %s",
					tc.score, got.String(), tc.want.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CalculateWorkProductOutcome
// ---------------------------------------------------------------------------

func TestCalculateWorkProductOutcome(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateWorkProductOutcome(
		decimal.NewFromInt(3),
		decimal.NewFromInt(4),
		decimal.NewFromInt(5),
	)

	expected := decimal.NewFromInt(12)
	if !result.Equal(expected) {
		t.Errorf("CalculateWorkProductOutcome(3, 4, 5) = %s, want %s", result, expected)
	}
}

func TestCalculateWorkProductOutcome_Zeros(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateWorkProductOutcome(
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
	)

	if !result.IsZero() {
		t.Errorf("CalculateWorkProductOutcome(0, 0, 0) = %s, want 0", result)
	}
}

// ---------------------------------------------------------------------------
// CalculateWeightedCategoryScore
// ---------------------------------------------------------------------------

func TestCalculateWeightedCategoryScore_Valid(t *testing.T) {
	svc := newTestScoringService()

	scores := []CategoryScore{
		{CategoryID: "CAT-A", Score: decimal.NewFromInt(70), Weight: decimal.NewFromInt(60)},
		{CategoryID: "CAT-B", Score: decimal.NewFromInt(30), Weight: decimal.NewFromInt(40)},
	}

	result, err := svc.CalculateWeightedCategoryScore(scores)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// (70 * 60 / 100) + (30 * 40 / 100) = 42 + 12 = 54
	expected := decimal.NewFromInt(54)
	if !result.Equal(expected) {
		t.Errorf("CalculateWeightedCategoryScore = %s, want %s", result, expected)
	}
}

func TestCalculateWeightedCategoryScore_UnbalancedWeights(t *testing.T) {
	svc := newTestScoringService()

	scores := []CategoryScore{
		{CategoryID: "CAT-A", Score: decimal.NewFromInt(70), Weight: decimal.NewFromInt(50)},
		{CategoryID: "CAT-B", Score: decimal.NewFromInt(30), Weight: decimal.NewFromInt(40)},
	}

	_, err := svc.CalculateWeightedCategoryScore(scores)
	if err == nil {
		t.Fatal("expected error for unbalanced weights (50+40=90)")
	}
	if !errors.Is(err, ErrWeightsNotBalanced) {
		t.Errorf("expected ErrWeightsNotBalanced, got: %v", err)
	}
}

func TestCalculateWeightedCategoryScore_EmptyScores(t *testing.T) {
	svc := newTestScoringService()

	result, err := svc.CalculateWeightedCategoryScore([]CategoryScore{})
	if err == nil {
		t.Fatal("expected error for empty scores")
	}
	if !errors.Is(err, ErrNoScoreData) {
		t.Errorf("expected ErrNoScoreData, got: %v", err)
	}
	if !result.IsZero() {
		t.Errorf("expected zero result for empty scores, got %s", result)
	}
}

// ---------------------------------------------------------------------------
// ValidateCategoryWeights
// ---------------------------------------------------------------------------

func TestValidateCategoryWeights_Valid(t *testing.T) {
	svc := newTestScoringService()

	weights := []decimal.Decimal{
		decimal.NewFromInt(30),
		decimal.NewFromInt(70),
	}

	if err := svc.ValidateCategoryWeights(weights); err != nil {
		t.Errorf("expected valid weights [30, 70], got error: %v", err)
	}
}

func TestValidateCategoryWeights_Invalid(t *testing.T) {
	svc := newTestScoringService()

	weights := []decimal.Decimal{
		decimal.NewFromInt(30),
		decimal.NewFromInt(60),
	}

	err := svc.ValidateCategoryWeights(weights)
	if err == nil {
		t.Fatal("expected error for weights [30, 60] (sum=90)")
	}
	if !errors.Is(err, ErrWeightsNotBalanced) {
		t.Errorf("expected ErrWeightsNotBalanced, got: %v", err)
	}

	// Verify structured error extraction.
	var wvErr *WeightValidationError
	if errors.As(err, &wvErr) {
		if wvErr.ExpectedTotal != 100 {
			t.Errorf("expected ExpectedTotal=100, got %f", wvErr.ExpectedTotal)
		}
		if wvErr.ActualTotal != 90 {
			t.Errorf("expected ActualTotal=90, got %f", wvErr.ActualTotal)
		}
	} else {
		t.Error("expected error to be WeightValidationError")
	}
}

func TestValidateCategoryWeights_WithTolerance(t *testing.T) {
	svc := newTestScoringService()

	// 30.005 + 69.995 = 100.000 (within tolerance of 0.01)
	weights := []decimal.Decimal{
		decimal.NewFromFloat(30.005),
		decimal.NewFromFloat(69.995),
	}

	if err := svc.ValidateCategoryWeights(weights); err != nil {
		t.Errorf("expected weights within tolerance to pass, got error: %v", err)
	}
}

func TestValidateCategoryWeights_JustOutsideTolerance(t *testing.T) {
	svc := newTestScoringService()

	// 30 + 69.98 = 99.98 -> diff = 0.02 > tolerance 0.01
	weights := []decimal.Decimal{
		decimal.NewFromInt(30),
		decimal.NewFromFloat(69.98),
	}

	err := svc.ValidateCategoryWeights(weights)
	if err == nil {
		t.Fatal("expected error for weights just outside tolerance")
	}
}

// ---------------------------------------------------------------------------
// CalculateBehavioralReviewAverage
// ---------------------------------------------------------------------------

func TestCalculateBehavioralReviewAverage_Normal(t *testing.T) {
	svc := newTestScoringService()

	ratings := []decimal.Decimal{
		decimal.NewFromInt(3),
		decimal.NewFromInt(4),
		decimal.NewFromInt(5),
	}

	result := svc.CalculateBehavioralReviewAverage(ratings)

	// (3 + 4 + 5) / 3 = 4
	expected := decimal.NewFromInt(4)
	if !result.Equal(expected) {
		t.Errorf("CalculateBehavioralReviewAverage([3,4,5]) = %s, want %s", result, expected)
	}
}

func TestCalculateBehavioralReviewAverage_WithZeros(t *testing.T) {
	svc := newTestScoringService()

	// Zeros are excluded from the average.
	ratings := []decimal.Decimal{
		decimal.NewFromInt(0),
		decimal.NewFromInt(3),
		decimal.NewFromInt(0),
		decimal.NewFromInt(5),
	}

	result := svc.CalculateBehavioralReviewAverage(ratings)

	// (3 + 5) / 2 = 4
	expected := decimal.NewFromInt(4)
	if !result.Equal(expected) {
		t.Errorf("CalculateBehavioralReviewAverage([0,3,0,5]) = %s, want %s", result, expected)
	}
}

func TestCalculateBehavioralReviewAverage_AllZeros(t *testing.T) {
	svc := newTestScoringService()

	ratings := []decimal.Decimal{
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
	}

	result := svc.CalculateBehavioralReviewAverage(ratings)
	if !result.IsZero() {
		t.Errorf("CalculateBehavioralReviewAverage([0,0,0]) = %s, want 0", result)
	}
}

func TestCalculateBehavioralReviewAverage_Empty(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateBehavioralReviewAverage([]decimal.Decimal{})
	if !result.IsZero() {
		t.Errorf("CalculateBehavioralReviewAverage([]) = %s, want 0", result)
	}
}

// ---------------------------------------------------------------------------
// CalculateTechnicalWeightedScore
// ---------------------------------------------------------------------------

func TestCalculateTechnicalWeightedScore(t *testing.T) {
	svc := newTestScoringService()

	// selfAvg=4, supervisorAvg=3, selfWeight=30, supervisorWeight=70
	// (4*30/100) + (3*70/100) = 1.2 + 2.1 = 3.3
	result := svc.CalculateTechnicalWeightedScore(
		decimal.NewFromInt(4),
		decimal.NewFromInt(3),
		decimal.NewFromInt(30),
		decimal.NewFromInt(70),
	)

	expected := decimal.NewFromFloat(3.3)
	if !result.Equal(expected) {
		t.Errorf("CalculateTechnicalWeightedScore(4,3,30,70) = %s, want %s", result, expected)
	}
}

func TestCalculateTechnicalWeightedScore_DefaultWeights(t *testing.T) {
	svc := newTestScoringService()

	// Using standard 30/70 weighting with different averages.
	// selfAvg=5, supervisorAvg=4, selfWeight=30, supervisorWeight=70
	// (5*30/100) + (4*70/100) = 1.5 + 2.8 = 4.3
	result := svc.CalculateTechnicalWeightedScore(
		decimal.NewFromInt(5),
		decimal.NewFromInt(4),
		decimal.NewFromInt(30),
		decimal.NewFromInt(70),
	)

	expected := decimal.NewFromFloat(4.3)
	if !result.Equal(expected) {
		t.Errorf("CalculateTechnicalWeightedScore(5,4,30,70) = %s, want %s", result, expected)
	}
}

func TestCalculateTechnicalWeightedScore_EqualWeights(t *testing.T) {
	svc := newTestScoringService()

	// selfAvg=3, supervisorAvg=5, selfWeight=50, supervisorWeight=50
	// (3*50/100) + (5*50/100) = 1.5 + 2.5 = 4.0
	result := svc.CalculateTechnicalWeightedScore(
		decimal.NewFromInt(3),
		decimal.NewFromInt(5),
		decimal.NewFromInt(50),
		decimal.NewFromInt(50),
	)

	expected := decimal.NewFromInt(4)
	if !result.Equal(expected) {
		t.Errorf("CalculateTechnicalWeightedScore(3,5,50,50) = %s, want %s", result, expected)
	}
}

// ---------------------------------------------------------------------------
// CalculateCompetencyGap
// ---------------------------------------------------------------------------

func TestCalculateCompetencyGap_HasGap(t *testing.T) {
	svc := newTestScoringService()

	gap, hasGap := svc.CalculateCompetencyGap(
		decimal.NewFromInt(5),
		decimal.NewFromInt(3),
	)

	if !hasGap {
		t.Error("expected hasGap=true when expected > actual")
	}
	expectedGap := decimal.NewFromInt(2)
	if !gap.Equal(expectedGap) {
		t.Errorf("expected gap=2, got %s", gap)
	}
}

func TestCalculateCompetencyGap_NoGap(t *testing.T) {
	svc := newTestScoringService()

	gap, hasGap := svc.CalculateCompetencyGap(
		decimal.NewFromInt(3),
		decimal.NewFromInt(5),
	)

	if hasGap {
		t.Error("expected hasGap=false when actual >= expected")
	}
	if !gap.IsZero() {
		t.Errorf("expected gap=0, got %s", gap)
	}
}

func TestCalculateCompetencyGap_Equal(t *testing.T) {
	svc := newTestScoringService()

	gap, hasGap := svc.CalculateCompetencyGap(
		decimal.NewFromInt(3),
		decimal.NewFromInt(3),
	)

	if hasGap {
		t.Error("expected hasGap=false when expected == actual")
	}
	if !gap.IsZero() {
		t.Errorf("expected gap=0, got %s", gap)
	}
}

// ---------------------------------------------------------------------------
// ApplyHRDDeduction
// ---------------------------------------------------------------------------

func TestApplyHRDDeduction_Normal(t *testing.T) {
	svc := newTestScoringService()

	result := svc.ApplyHRDDeduction(
		decimal.NewFromInt(80),
		decimal.NewFromInt(5),
	)

	expected := decimal.NewFromInt(75)
	if !result.Equal(expected) {
		t.Errorf("ApplyHRDDeduction(80, 5) = %s, want %s", result, expected)
	}
}

func TestApplyHRDDeduction_OverDeduction(t *testing.T) {
	svc := newTestScoringService()

	// Deduction exceeds score -- result should be floored at 0.
	result := svc.ApplyHRDDeduction(
		decimal.NewFromInt(3),
		decimal.NewFromInt(10),
	)

	if !result.IsZero() {
		t.Errorf("ApplyHRDDeduction(3, 10) = %s, want 0 (should not go negative)", result)
	}
}

func TestApplyHRDDeduction_ZeroDeduction(t *testing.T) {
	svc := newTestScoringService()

	result := svc.ApplyHRDDeduction(
		decimal.NewFromInt(80),
		decimal.NewFromInt(0),
	)

	expected := decimal.NewFromInt(80)
	if !result.Equal(expected) {
		t.Errorf("ApplyHRDDeduction(80, 0) = %s, want %s", result, expected)
	}
}

func TestApplyHRDDeduction_ExactDeduction(t *testing.T) {
	svc := newTestScoringService()

	result := svc.ApplyHRDDeduction(
		decimal.NewFromInt(5),
		decimal.NewFromInt(5),
	)

	if !result.IsZero() {
		t.Errorf("ApplyHRDDeduction(5, 5) = %s, want 0", result)
	}
}

// ---------------------------------------------------------------------------
// CalculatePeriodScore
// ---------------------------------------------------------------------------

func TestCalculatePeriodScore_FullCalculation(t *testing.T) {
	svc := newTestScoringService()

	// workProduct=40, objective=30, competency=10 -> finalScore=80
	// hrdDeduction=5 -> adjusted=75
	// maxPoints=100 -> pct=75%
	// 66 <= 75 < 80 -> Competent
	result := svc.CalculatePeriodScore(
		decimal.NewFromInt(40),  // workProductScore
		decimal.NewFromInt(30),  // objectiveScore
		decimal.NewFromInt(10),  // competencyScore
		decimal.NewFromInt(100), // maxPoints
		decimal.NewFromInt(5),   // hrdDeduction
	)

	expectedScore := decimal.NewFromInt(75)
	if !result.FinalScore.Equal(expectedScore) {
		t.Errorf("FinalScore = %s, want %s", result.FinalScore, expectedScore)
	}

	expectedPct := decimal.NewFromInt(75)
	if !result.ScorePercentage.Equal(expectedPct) {
		t.Errorf("ScorePercentage = %s, want %s", result.ScorePercentage, expectedPct)
	}

	if result.Grade != enums.PerformanceGradeCompetent {
		t.Errorf("Grade = %s, want Competent", result.Grade.String())
	}

	if result.IsUnderPerforming {
		t.Error("expected IsUnderPerforming=false for 75%")
	}
}

func TestCalculatePeriodScore_UnderPerforming(t *testing.T) {
	svc := newTestScoringService()

	// workProduct=10, objective=10, competency=5 -> finalScore=25
	// hrdDeduction=0 -> adjusted=25
	// maxPoints=100 -> pct=25%
	// 25 < 30 -> Probation, IsUnderPerforming = true (< 50)
	result := svc.CalculatePeriodScore(
		decimal.NewFromInt(10),
		decimal.NewFromInt(10),
		decimal.NewFromInt(5),
		decimal.NewFromInt(100),
		decimal.NewFromInt(0),
	)

	if result.Grade != enums.PerformanceGradeProbation {
		t.Errorf("Grade = %s, want Probation", result.Grade.String())
	}
	if !result.IsUnderPerforming {
		t.Error("expected IsUnderPerforming=true for 25%")
	}
}

func TestCalculatePeriodScore_ZeroMaxPoints(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculatePeriodScore(
		decimal.NewFromInt(40),
		decimal.NewFromInt(30),
		decimal.NewFromInt(10),
		decimal.NewFromInt(0), // maxPoints = 0
		decimal.NewFromInt(0),
	)

	// ScorePercentage should be 0 (no division by zero).
	if !result.ScorePercentage.IsZero() {
		t.Errorf("ScorePercentage = %s, want 0 when maxPoints=0", result.ScorePercentage)
	}
}

// ---------------------------------------------------------------------------
// CalculateScorePercentage
// ---------------------------------------------------------------------------

func TestCalculateScorePercentage_Normal(t *testing.T) {
	svc := newTestScoringService()

	// 80 / 200 * 100 = 40%
	result := svc.CalculateScorePercentage(
		decimal.NewFromInt(80),
		decimal.NewFromInt(200),
	)

	expected := decimal.NewFromInt(40)
	if !result.Equal(expected) {
		t.Errorf("CalculateScorePercentage(80, 200) = %s, want %s", result, expected)
	}
}

func TestCalculateScorePercentage_ZeroMax(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateScorePercentage(
		decimal.NewFromInt(80),
		decimal.NewFromInt(0),
	)

	if !result.IsZero() {
		t.Errorf("CalculateScorePercentage(80, 0) = %s, want 0 (avoid division by zero)", result)
	}
}

func TestCalculateScorePercentage_PerfectScore(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateScorePercentage(
		decimal.NewFromInt(100),
		decimal.NewFromInt(100),
	)

	expected := decimal.NewFromInt(100)
	if !result.Equal(expected) {
		t.Errorf("CalculateScorePercentage(100, 100) = %s, want %s", result, expected)
	}
}

func TestCalculateScorePercentage_ZeroScore(t *testing.T) {
	svc := newTestScoringService()

	result := svc.CalculateScorePercentage(
		decimal.NewFromInt(0),
		decimal.NewFromInt(100),
	)

	if !result.IsZero() {
		t.Errorf("CalculateScorePercentage(0, 100) = %s, want 0", result)
	}
}

// ---------------------------------------------------------------------------
// Grade boundary edge cases (table-driven)
// ---------------------------------------------------------------------------

func TestDetermineGrade_Boundaries(t *testing.T) {
	svc := newTestScoringService()

	// Test exact boundary values to verify < vs <=.
	tests := []struct {
		score float64
		want  enums.PerformanceGrade
	}{
		{29.99, enums.PerformanceGradeProbation},
		{30.00, enums.PerformanceGradeDeveloping},
		{49.99, enums.PerformanceGradeDeveloping},
		{50.00, enums.PerformanceGradeProgressive},
		{65.99, enums.PerformanceGradeProgressive},
		{66.00, enums.PerformanceGradeCompetent},
		{79.99, enums.PerformanceGradeCompetent},
		{80.00, enums.PerformanceGradeAccomplished},
		{89.99, enums.PerformanceGradeAccomplished},
		{90.00, enums.PerformanceGradeExemplary},
	}

	for _, tc := range tests {
		t.Run(decimal.NewFromFloat(tc.score).String(), func(t *testing.T) {
			got := svc.DetermineGrade(decimal.NewFromFloat(tc.score))
			if got != tc.want {
				t.Errorf("DetermineGrade(%v) = %s, want %s",
					tc.score, got.String(), tc.want.String())
			}
		})
	}
}
