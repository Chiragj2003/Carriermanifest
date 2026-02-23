// Package engine implements the vector-based career recommendation system.
// This is the core AI system of CareerManifest — a production-grade,
// mathematically robust, explainable, and ML-ready scoring engine.
//
// Architecture (3-layer):
//
//	Assessment Engine
//	  ├── Feature Aggregation Layer  (aggregator.go)
//	  ├── Career Scoring Engine      (scorer.go + matrix.go)
//	  ├── Risk Engine                (risk.go)
//	  ├── Normalization Layer        (normalize.go)
//	  ├── Ranking + Confidence       (normalize.go)
//	  └── Explanation Generator      (explain.go)
package engine

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/careermanifest/backend/internal/dto"
)

// ScoringEngine evaluates assessment answers and produces career recommendations
// using vector-based scoring with feature aggregation.
type ScoringEngine struct{}

// NewScoringEngine creates a new ScoringEngine.
func NewScoringEngine() *ScoringEngine {
	return &ScoringEngine{}
}

// QuestionData is a simplified question structure for the engine.
type QuestionData struct {
	ID           uint64               `json:"id"`
	Category     string               `json:"category"`
	Weights      []dto.QuestionWeight `json:"weights"`
	DisplayOrder int                  `json:"display_order"`
}

// ParseQuestionWeights parses JSON weight data from the database.
func ParseQuestionWeights(weightsJSON string) ([]dto.QuestionWeight, error) {
	var weights []dto.QuestionWeight
	if err := json.Unmarshal([]byte(weightsJSON), &weights); err != nil {
		return nil, fmt.Errorf("failed to parse weights: %w", err)
	}
	return weights, nil
}

// ComputeResult processes all answers against question data and produces a full result.
//
// Flow:
//  1. Answers → Feature Aggregator → UserProfile
//  2. UserProfile → Vector Engine → Raw CareerScores
//  3. Risk Engine computes independent risk + applies career penalties
//  4. Scores normalized (min-max 0–1)
//  5. Careers ranked, confidence computed
//  6. Explanations generated (deterministic, data-driven)
//  7. Enrichment data attached (salary, roadmap, skills, exams, colleges)
//  8. Top 3 returned with full breakdown
func (e *ScoringEngine) ComputeResult(answers []dto.AnswerItem, questionsJSON []QuestionData) (*dto.AssessmentResult, error) {
	// ── Step 1: Feature Aggregation ──────────────────────────
	profile := AggregateProfile(answers, questionsJSON)

	// ── Step 2: Vector-Based Scoring ─────────────────────────
	rawScores := ComputeRawScores(profile)

	// ── Step 3: Risk Engine (independent) ────────────────────
	riskResult := ComputeRisk(profile)

	// ── Step 4: Apply Risk Penalties ─────────────────────────
	adjustedScores, appliedPenalties := ApplyRiskPenalties(rawScores, profile)

	// ── Step 5: Normalize & Rank ─────────────────────────────
	ranked := NormalizeAndRank(adjustedScores)

	if len(ranked.Rankings) == 0 {
		return nil, fmt.Errorf("no career scores computed")
	}

	bestCareer := ranked.Rankings[0].Career

	// ── Step 6: Build DTO career scores (backward compatible) ──
	careerScores := make([]dto.CareerScore, len(ranked.Rankings))
	for i, r := range ranked.Rankings {
		careerScores[i] = dto.CareerScore{
			Category:   r.Career.String(),
			Score:      math.Round(r.RawScore*100) / 100,
			MaxScore:   math.Round(ranked.Rankings[0].RawScore*100) / 100, // best raw as reference
			Percentage: r.Percentage,
		}
	}

	// ── Step 7: Generate Explanations ────────────────────────
	explanations := GenerateTop3Explanations(ranked.Rankings, profile, appliedPenalties)

	explanationDTOs := make([]dto.CareerExplanationDTO, len(explanations))
	for i, exp := range explanations {
		factors := make([]dto.FeatureContributionDTO, len(exp.TopFactors))
		for j, f := range exp.TopFactors {
			factors[j] = dto.FeatureContributionDTO{
				Feature:      f.Feature,
				UserValue:    f.UserValue,
				CareerWeight: f.CareerWeight,
				Contribution: f.Contribution,
				Percentage:   f.Percentage,
			}
		}
		var penalties []string
		for _, p := range exp.RiskPenalties {
			penalties = append(penalties, fmt.Sprintf("-%0.f%%: %s", p.Penalty*100, p.Reason))
		}
		explanationDTOs[i] = dto.CareerExplanationDTO{
			Career:     exp.Career.String(),
			TopFactors: factors,
			Summary:    exp.Summary,
			Penalties:  penalties,
		}
	}

	// ── Step 8: Risk DTO ─────────────────────────────────────
	risk := dto.RiskAssessment{
		Score:   riskResult.Score,
		Level:   riskResult.Level,
		Factors: riskResult.Factors,
	}

	// ── Step 9: Profile DTO ──────────────────────────────────
	profileDTO := dto.UserProfileDTO{
		AcademicStrength:  math.Round(profile.AcademicStrength()*1000) / 1000,
		FinancialPressure: math.Round(profile.FinancialPressure()*1000) / 1000,
		RiskTolerance:     math.Round(profile.RiskTolerance()*1000) / 1000,
		LeadershipScore:   math.Round(profile.LeadershipScore()*1000) / 1000,
		TechAffinity:      math.Round(profile.TechAffinity()*1000) / 1000,
		GovtInterest:      math.Round(profile.GovtInterest()*1000) / 1000,
		AbroadInterest:    math.Round(profile.AbroadInterest()*1000) / 1000,
		IncomeUrgency:     math.Round(profile.IncomeUrgency()*1000) / 1000,
		CareerInstability: math.Round(profile.CareerInstability()*1000) / 1000,
	}

	// ── Step 10: Enrichment (salary, roadmap, skills, exams, colleges) ──
	salaryProjection := GetSalaryProjection(bestCareer)
	roadmap := GetRoadmap(bestCareer)
	skills := GetRequiredSkills(bestCareer)
	exams := GetSuggestedExams(bestCareer)
	colleges := GetSuggestedColleges(bestCareer)

	return &dto.AssessmentResult{
		Scores:            careerScores,
		BestCareerPath:    bestCareer.String(),
		Confidence:        ranked.Confidence,
		IsMultiFit:        ranked.IsMultiFit,
		Risk:              risk,
		Profile:           profileDTO,
		Explanations:      explanationDTOs,
		SalaryProjection:  salaryProjection,
		Roadmap:           roadmap,
		RequiredSkills:    skills,
		SuggestedExams:    exams,
		SuggestedColleges: colleges,
		Version: dto.VersionInfo{
			Assessment:    AssessmentVersion,
			WeightMatrix:  WeightMatrixVersion,
			FeatureMap:    FeatureMapVersion,
			ModelType:     ModelType,
			ModelAccuracy: ModelAccuracy,
			ModelF1Score:  ModelF1Score,
		},
	}, nil
}
