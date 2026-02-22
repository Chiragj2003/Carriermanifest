package engine

import (
	"math"
	"sort"
)

// NormalizedScore holds a career's score after min-max normalization.
type NormalizedScore struct {
	Career     Career
	RawScore   float64
	Normalized float64 // 0–1 after min-max scaling
	Percentage float64 // 0–100 for display
}

// RankedResult holds the final ranked list with confidence metric.
type RankedResult struct {
	Rankings   []NormalizedScore
	Confidence float64 // 0–1: how clearly the top career dominates
	IsMultiFit bool    // true if confidence < 0.1
}

// NormalizeAndRank applies min-max normalization, sorts descending, and computes confidence.
func NormalizeAndRank(scores []RawCareerScore) RankedResult {
	if len(scores) == 0 {
		return RankedResult{}
	}

	// Find min and max raw scores
	minScore := scores[0].Score
	maxScore := scores[0].Score
	for _, s := range scores[1:] {
		if s.Score < minScore {
			minScore = s.Score
		}
		if s.Score > maxScore {
			maxScore = s.Score
		}
	}

	// Min-max normalize
	scoreRange := maxScore - minScore
	if scoreRange == 0 {
		scoreRange = 1 // prevent division by zero
	}

	normalized := make([]NormalizedScore, len(scores))
	for i, s := range scores {
		norm := (s.Score - minScore) / scoreRange
		normalized[i] = NormalizedScore{
			Career:     s.Career,
			RawScore:   math.Round(s.Score*1000) / 1000,
			Normalized: math.Round(norm*1000) / 1000,
			Percentage: math.Round(norm*10000) / 100, // 2 decimal places
		}
	}

	// Sort descending by normalized score
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].Normalized > normalized[j].Normalized
	})

	// Compute confidence: (TopScore - SecondScore) / TopScore
	confidence := 0.0
	if len(normalized) >= 2 && normalized[0].Normalized > 0 {
		confidence = (normalized[0].Normalized - normalized[1].Normalized) / normalized[0].Normalized
	}
	confidence = math.Round(confidence*1000) / 1000

	return RankedResult{
		Rankings:   normalized,
		Confidence: confidence,
		IsMultiFit: confidence < 0.1,
	}
}
