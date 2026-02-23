package engine

import "math"

// RiskResult holds the computed risk assessment.
type RiskResult struct {
	Score   float64            // Overall risk score (0–10)
	Level   string             // "Low", "Medium", "High"
	Factors map[string]float64 // Individual factor scores (0–10 scale)
}

// ComputeRisk calculates the overall risk score from the UserProfile.
// Risk is computed independently from career scoring.
//
// Formula:
//
//	RiskScore = (IncomeUrgency × 0.35) + (FinancialPressure × 0.25) +
//	            (RiskTolerance × 0.20) + (CareerInstability × 0.20)
//
// All profile features are on 0–1 scale; risk is projected to 0–10.
func ComputeRisk(profile *UserProfile) RiskResult {
	// Scale 0–1 features to 0–10 for risk factors
	incomeUrgency := profile.IncomeUrgency() * 10
	financialPressure := profile.FinancialPressure() * 10
	riskTolerance := profile.RiskTolerance() * 10
	careerInstability := profile.CareerInstability() * 10

	riskScore := (incomeUrgency * 0.35) +
		(financialPressure * 0.25) +
		(riskTolerance * 0.20) +
		(careerInstability * 0.20)

	riskScore = math.Round(riskScore*100) / 100

	level := "Low"
	if riskScore > 6 {
		level = "High"
	} else if riskScore > 3 {
		level = "Medium"
	}

	factors := map[string]float64{
		"income_urgency":     math.Round(incomeUrgency*100) / 100,
		"financial_pressure": math.Round(financialPressure*100) / 100,
		"risk_tolerance":     math.Round(riskTolerance*100) / 100,
		"career_instability": math.Round(careerInstability*100) / 100,
	}

	return RiskResult{
		Score:   riskScore,
		Level:   level,
		Factors: factors,
	}
}

// careerRiskPenalty defines per-career penalty rules based on risk profile.
// Returns a multiplier (0–1) where 1.0 means no penalty.
type careerPenaltyRule struct {
	Career    Career
	Condition func(profile *UserProfile) bool
	Penalty   float64 // fractional reduction, e.g. 0.20 = -20%
	Reason    string
}

// riskPenaltyRules defines career-specific risk penalties.
var riskPenaltyRules = []careerPenaltyRule{
	{
		Career:    CareerStartup,
		Condition: func(p *UserProfile) bool { return p.FinancialPressure() > 0.6 },
		Penalty:   0.20,
		Reason:    "High financial pressure makes startup risky",
	},
	{
		Career:    CareerStartup,
		Condition: func(p *UserProfile) bool { return p.IncomeUrgency() > 0.7 },
		Penalty:   0.15,
		Reason:    "Urgent income need conflicts with startup timeline",
	},
	{
		Career:    CareerMSAbroad,
		Condition: func(p *UserProfile) bool { return p.FinancialPressure() > 0.65 },
		Penalty:   0.15,
		Reason:    "High financial pressure makes MS Abroad difficult",
	},
	{
		Career:    CareerMSAbroad,
		Condition: func(p *UserProfile) bool { return p.IncomeUrgency() > 0.6 },
		Penalty:   0.10,
		Reason:    "Income urgency conflicts with 2-year study abroad",
	},
	{
		Career:    CareerHigherStudies,
		Condition: func(p *UserProfile) bool { return p.IncomeUrgency() > 0.65 },
		Penalty:   0.25,
		Reason:    "Urgent income need conflicts with extended studies",
	},
	{
		Career:    CareerMBA,
		Condition: func(p *UserProfile) bool { return p.FinancialPressure() > 0.7 },
		Penalty:   0.15,
		Reason:    "High financial pressure makes MBA fees challenging",
	},
	{
		Career:    CareerGovt,
		Condition: func(p *UserProfile) bool { return p.RiskTolerance() > 0.8 },
		Penalty:   0.10,
		Reason:    "High risk appetite may lead to dissatisfaction with govt job security",
	},
	{
		Career:    CareerCreative,
		Condition: func(p *UserProfile) bool { return p.FinancialPressure() > 0.7 },
		Penalty:   0.15,
		Reason:    "High financial pressure makes freelance/creative career risky",
	},
	{
		Career:    CareerCreative,
		Condition: func(p *UserProfile) bool { return p.IncomeUrgency() > 0.7 },
		Penalty:   0.10,
		Reason:    "Urgent income need conflicts with creative career ramp-up",
	},
	{
		Career:    CareerHealthcare,
		Condition: func(p *UserProfile) bool { return p.IncomeUrgency() > 0.7 },
		Penalty:   0.20,
		Reason:    "Urgent income need conflicts with long medical training (5.5+ years)",
	},
	{
		Career:    CareerHealthcare,
		Condition: func(p *UserProfile) bool { return p.RiskTolerance() > 0.8 },
		Penalty:   0.08,
		Reason:    "High risk appetite may lead to dissatisfaction with structured medical career",
	},
}

// RiskPenalty holds the penalty information applied to a career.
type RiskPenalty struct {
	Penalty float64
	Reason  string
}

// ApplyRiskPenalties adjusts raw career scores based on risk profile.
// Returns the adjusted scores and a map of applied penalties per career.
func ApplyRiskPenalties(scores []RawCareerScore, profile *UserProfile) ([]RawCareerScore, map[Career][]RiskPenalty) {
	adjusted := make([]RawCareerScore, len(scores))
	copy(adjusted, scores)
	appliedPenalties := make(map[Career][]RiskPenalty)

	for _, rule := range riskPenaltyRules {
		if rule.Condition(profile) {
			idx := int(rule.Career)
			if idx < len(adjusted) {
				adjusted[idx].Score *= (1.0 - rule.Penalty)
				appliedPenalties[rule.Career] = append(appliedPenalties[rule.Career], RiskPenalty{
					Penalty: rule.Penalty,
					Reason:  rule.Reason,
				})
			}
		}
	}

	return adjusted, appliedPenalties
}
