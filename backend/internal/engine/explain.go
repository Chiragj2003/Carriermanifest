package engine

import (
	"fmt"
	"math"
	"sort"
)

// FeatureContribution describes how a single feature contributed to a career's score.
type FeatureContribution struct {
	Feature      string  // Human-readable feature name
	UserValue    float64 // The user's profile value for this feature (0–1)
	CareerWeight float64 // The career's weight for this feature
	Contribution float64 // UserValue × CareerWeight (signed)
	Percentage   float64 // Percentage of total positive contribution
}

// Explanation holds the deterministic, data-driven explanation for a career recommendation.
type Explanation struct {
	Career        Career
	TopFactors    []FeatureContribution // Sorted by |contribution| descending
	Summary       string                // Human-readable summary
	RiskPenalties []RiskPenalty          // Any risk adjustments applied
}

// GenerateExplanation builds a deterministic explanation for why a career was recommended.
// It computes per-feature contributions and ranks them by impact.
func GenerateExplanation(career Career, profile *UserProfile, penalties []RiskPenalty) Explanation {
	weights := GetCareerWeights(career)
	userVec := profile.Vector()

	var contributions []FeatureContribution
	totalPositive := 0.0

	for i := 0; i < NumFeatures; i++ {
		contrib := userVec[i] * weights[i]
		if contrib > 0 {
			totalPositive += contrib
		}
		contributions = append(contributions, FeatureContribution{
			Feature:      FeatureNames[i],
			UserValue:    math.Round(userVec[i]*1000) / 1000,
			CareerWeight: weights[i],
			Contribution: math.Round(contrib*1000) / 1000,
		})
	}

	// Calculate percentage of positive contribution
	if totalPositive > 0 {
		for i := range contributions {
			if contributions[i].Contribution > 0 {
				contributions[i].Percentage = math.Round((contributions[i].Contribution/totalPositive)*10000) / 100
			}
		}
	}

	// Sort by absolute contribution descending
	sort.Slice(contributions, func(i, j int) bool {
		return math.Abs(contributions[i].Contribution) > math.Abs(contributions[j].Contribution)
	})

	// Build human-readable summary
	summary := fmt.Sprintf("%s recommended because:", career.String())
	count := 0
	for _, c := range contributions {
		if c.Contribution > 0 && count < 4 {
			summary += fmt.Sprintf("\n• %s (+%.0f%%)", c.Feature, c.Percentage)
			count++
		}
	}

	// Note any negative factors
	for _, c := range contributions {
		if c.Contribution < 0 && math.Abs(c.Contribution) > 0.05 {
			summary += fmt.Sprintf("\n• %s (caution: %.0f%% drag)", c.Feature,
				math.Abs(c.Contribution/totalPositive)*100)
		}
	}

	// Note risk penalties
	if len(penalties) > 0 {
		summary += "\n\nRisk adjustments:"
		for _, p := range penalties {
			summary += fmt.Sprintf("\n• -%0.f%%: %s", p.Penalty*100, p.Reason)
		}
	}

	return Explanation{
		Career:        career,
		TopFactors:    contributions,
		Summary:       summary,
		RiskPenalties: penalties,
	}
}

// GenerateTop3Explanations generates explanations for the top 3 ranked careers.
func GenerateTop3Explanations(rankings []NormalizedScore, profile *UserProfile, penalties map[Career][]RiskPenalty) []Explanation {
	n := 3
	if len(rankings) < n {
		n = len(rankings)
	}

	explanations := make([]Explanation, n)
	for i := 0; i < n; i++ {
		career := rankings[i].Career
		var careerPenalties []RiskPenalty
		if p, ok := penalties[career]; ok {
			careerPenalties = p
		}
		explanations[i] = GenerateExplanation(career, profile, careerPenalties)
	}
	return explanations
}
