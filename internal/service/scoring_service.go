package service

import (
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// ---------------------------------------------------------------------------
// scoringService consolidates all pure scoring calculations for the PMS.
// All methods are stateless — no database access. They implement the scoring
// rules from the .NET CompetencyReviewService, PeriodScoreService, and
// WorkProductService.
// ---------------------------------------------------------------------------

// CategoryScore holds a single category's score, weight, and max points
// for weighted-score aggregation.
type CategoryScore struct {
	CategoryID string
	Score      decimal.Decimal
	Weight     decimal.Decimal // percentage, e.g. 30 = 30%
	MaxPoints  decimal.Decimal
}

// ScoringResult is the output of a full period-score calculation.
type ScoringResult struct {
	FinalScore        decimal.Decimal
	ScorePercentage   decimal.Decimal
	Grade             enums.PerformanceGrade
	IsUnderPerforming bool
	CategoryBreakdown []CategoryScore
}

// CompetencyScoreResult holds the gap analysis for a single competency.
type CompetencyScoreResult struct {
	CompetencyID   string
	AverageRating  decimal.Decimal
	ExpectedRating decimal.Decimal
	Gap            decimal.Decimal
	HasGap         bool
}

// ---------------------------------------------------------------------------
// Grade thresholds (upper bound is exclusive except for Exemplary).
// ---------------------------------------------------------------------------

var (
	thresholdDeveloping   = decimal.NewFromInt(30)
	thresholdProgressive  = decimal.NewFromInt(50)
	thresholdCompetent    = decimal.NewFromInt(66)
	thresholdAccomplished = decimal.NewFromInt(80)
	thresholdExemplary    = decimal.NewFromInt(90)

	hundred          = decimal.NewFromInt(100)
	zero             = decimal.NewFromInt(0)
	weightTolerance  = decimal.NewFromFloat(0.01)
	underPerfCutoff  = decimal.NewFromInt(50)
)

// ---------------------------------------------------------------------------
// scoringService
// ---------------------------------------------------------------------------

type scoringService struct {
	log zerolog.Logger
}

func newScoringService(log zerolog.Logger) *scoringService {
	return &scoringService{
		log: log.With().Str("sub", "scoring").Logger(),
	}
}

// DetermineGrade maps a score percentage to the appropriate performance grade.
func (s *scoringService) DetermineGrade(scorePercentage decimal.Decimal) enums.PerformanceGrade {
	switch {
	case scorePercentage.LessThan(thresholdDeveloping):
		return enums.PerformanceGradeProbation
	case scorePercentage.LessThan(thresholdProgressive):
		return enums.PerformanceGradeDeveloping
	case scorePercentage.LessThan(thresholdCompetent):
		return enums.PerformanceGradeProgressive
	case scorePercentage.LessThan(thresholdAccomplished):
		return enums.PerformanceGradeCompetent
	case scorePercentage.LessThan(thresholdExemplary):
		return enums.PerformanceGradeAccomplished
	default:
		return enums.PerformanceGradeExemplary
	}
}

// CalculateWorkProductOutcome sums the three evaluation dimensions for a work product.
func (s *scoringService) CalculateWorkProductOutcome(timeliness, quality, output decimal.Decimal) decimal.Decimal {
	return timeliness.Add(quality).Add(output)
}

// CalculateWeightedCategoryScore computes the weighted sum across categories.
// Category weights must sum to 100% (within a tolerance of 0.01).
// Formula: sum(score * weight / 100) for each category.
func (s *scoringService) CalculateWeightedCategoryScore(scores []CategoryScore) (decimal.Decimal, error) {
	if len(scores) == 0 {
		return zero, ErrNoScoreData
	}

	// Extract weights and validate they sum to 100.
	weights := make([]decimal.Decimal, len(scores))
	for i, cs := range scores {
		weights[i] = cs.Weight
	}
	if err := s.ValidateCategoryWeights(weights); err != nil {
		return zero, err
	}

	total := zero
	for _, cs := range scores {
		// weightedScore = score * weight / 100
		weighted := cs.Score.Mul(cs.Weight).Div(hundred)
		total = total.Add(weighted)
	}
	return total, nil
}

// ValidateCategoryWeights checks that the supplied weights sum to 100%.
// Returns a WeightValidationError wrapping ErrWeightsNotBalanced on failure.
func (s *scoringService) ValidateCategoryWeights(weights []decimal.Decimal) error {
	total := zero
	for _, w := range weights {
		total = total.Add(w)
	}

	diff := total.Sub(hundred).Abs()
	if diff.GreaterThan(weightTolerance) {
		totalFloat, _ := total.Float64()
		return &WeightValidationError{
			ExpectedTotal: 100,
			ActualTotal:   totalFloat,
		}
	}
	return nil
}

// CalculateBehavioralReviewAverage computes the mean of non-zero ratings.
// Zero-valued ratings are excluded. Returns zero when no valid ratings remain.
func (s *scoringService) CalculateBehavioralReviewAverage(ratings []decimal.Decimal) decimal.Decimal {
	sum := zero
	count := int64(0)

	for _, r := range ratings {
		if r.IsZero() {
			continue
		}
		sum = sum.Add(r)
		count++
	}

	if count == 0 {
		return zero
	}
	return sum.Div(decimal.NewFromInt(count))
}

// CalculateTechnicalWeightedScore computes the 360-style technical review score.
// Formula: (selfAvg * selfWeight / 100) + (supervisorAvg * supervisorWeight / 100)
func (s *scoringService) CalculateTechnicalWeightedScore(
	selfAvg, supervisorAvg, selfWeight, supervisorWeight decimal.Decimal,
) decimal.Decimal {
	selfPart := selfAvg.Mul(selfWeight).Div(hundred)
	supervisorPart := supervisorAvg.Mul(supervisorWeight).Div(hundred)
	return selfPart.Add(supervisorPart)
}

// CalculateCompetencyGap determines the gap between expected and actual ratings.
// Gap = max(0, expected - actual). HasGap is true when expected != actual and gap > 0.
func (s *scoringService) CalculateCompetencyGap(expected, actual decimal.Decimal) (gap decimal.Decimal, hasGap bool) {
	diff := expected.Sub(actual)
	if diff.LessThanOrEqual(zero) {
		return zero, false
	}
	// expected != actual is implied because diff > 0.
	return diff, true
}

// ApplyHRDDeduction subtracts HRD-deducted points from a score, floored at zero.
func (s *scoringService) ApplyHRDDeduction(score, deductedPoints decimal.Decimal) decimal.Decimal {
	result := score.Sub(deductedPoints)
	if result.LessThan(zero) {
		return zero
	}
	return result
}

// CalculatePeriodScore assembles the final period score from all component scores.
//
//	FinalScore      = workProductScore + objectiveScore + competencyScore
//	AdjustedScore   = FinalScore - hrdDeduction (floored at 0)
//	ScorePercentage = (adjustedScore / maxPoints) * 100
//	Grade           = DetermineGrade(ScorePercentage)
//	IsUnderPerforming = ScorePercentage < 50
func (s *scoringService) CalculatePeriodScore(
	workProductScore, objectiveScore, competencyScore,
	maxPoints, hrdDeduction decimal.Decimal,
) ScoringResult {
	finalScore := workProductScore.Add(objectiveScore).Add(competencyScore)
	adjusted := s.ApplyHRDDeduction(finalScore, hrdDeduction)
	pct := s.CalculateScorePercentage(adjusted, maxPoints)
	grade := s.DetermineGrade(pct)

	return ScoringResult{
		FinalScore:        adjusted,
		ScorePercentage:   pct,
		Grade:             grade,
		IsUnderPerforming: pct.LessThan(underPerfCutoff),
	}
}

// CalculateScorePercentage converts a raw score to a percentage of maxPoints.
// Returns zero when maxPoints is zero to avoid division by zero.
func (s *scoringService) CalculateScorePercentage(score, maxPoints decimal.Decimal) decimal.Decimal {
	if maxPoints.IsZero() {
		return zero
	}
	return score.Div(maxPoints).Mul(hundred)
}
