package engine

import "github.com/careermanifest/backend/internal/dto"

// FeatureMapping defines how a question option maps to UserProfile features.
// Each mapping specifies which feature index to contribute to and the weight.
type FeatureMapping struct {
	FeatureIndex int
	Weight       float64
}

// questionFeatureMap defines the feature contributions for each question
// (identified by DisplayOrder 1–30). Each option index maps to a set of
// feature contributions. This replaces the old string-based career scoring.
//
// The values are on a 0–1 normalized scale:
//   - Academic questions → AcademicStrength, TechAffinity
//   - Financial questions → FinancialPressure, IncomeUrgency
//   - Personality questions → RiskTolerance, LeadershipScore, CareerInstability
//   - Career interest questions → GovtInterest, AbroadInterest, etc.
var questionFeatureMap = map[int]map[int][]FeatureMapping{
	// ================================================================
	// SECTION A: ACADEMIC BACKGROUND (Q1–Q8)
	// ================================================================

	// Q1: 10th board marks
	1: {
		0: {{FeatAcademicStrength, 0.10}},
		1: {{FeatAcademicStrength, 0.30}},
		2: {{FeatAcademicStrength, 0.55}},
		3: {{FeatAcademicStrength, 0.80}},
		4: {{FeatAcademicStrength, 1.00}},
	},
	// Q2: Stream (11th/12th)
	2: {
		0: {{FeatTechAffinity, 0.90}, {FeatAcademicStrength, 0.70}},                          // PCM
		1: {{FeatTechAffinity, 0.30}, {FeatAcademicStrength, 0.60}},                           // PCB
		2: {{FeatTechAffinity, 0.20}, {FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.40}}, // Commerce w/ Maths
		3: {{FeatTechAffinity, 0.10}, {FeatAcademicStrength, 0.35}, {FeatLeadershipScore, 0.30}}, // Commerce w/o Maths
		4: {{FeatTechAffinity, 0.05}, {FeatAcademicStrength, 0.30}, {FeatGovtInterest, 0.30}},    // Arts
	},
	// Q3: College major
	3: {
		0: {{FeatTechAffinity, 1.00}, {FeatAcademicStrength, 0.70}},                           // CS/IT
		1: {{FeatTechAffinity, 0.50}, {FeatAcademicStrength, 0.65}},                           // Other Engg
		2: {{FeatTechAffinity, 0.10}, {FeatLeadershipScore, 0.50}, {FeatAcademicStrength, 0.45}}, // Commerce/BBA
		3: {{FeatTechAffinity, 0.20}, {FeatAcademicStrength, 0.55}},                           // B.Sc
		4: {{FeatTechAffinity, 0.05}, {FeatAcademicStrength, 0.30}, {FeatGovtInterest, 0.25}},    // BA/Other
	},
	// Q4: College grades
	4: {
		0: {{FeatAcademicStrength, 0.15}},
		1: {{FeatAcademicStrength, 0.35}},
		2: {{FeatAcademicStrength, 0.55}},
		3: {{FeatAcademicStrength, 0.80}},
		4: {{FeatAcademicStrength, 1.00}},
	},
	// Q5: Competitive exams taken
	5: {
		0: {{FeatAcademicStrength, 0.20}},
		1: {{FeatAcademicStrength, 0.65}, {FeatTechAffinity, 0.50}},             // JEE/NEET
		2: {{FeatAcademicStrength, 0.60}, {FeatLeadershipScore, 0.60}},          // CAT/XAT
		3: {{FeatAcademicStrength, 0.75}, {FeatTechAffinity, 0.40}},             // GATE/NET
		4: {{FeatAcademicStrength, 0.55}, {FeatGovtInterest, 0.80}},             // UPSC/SSC
	},
	// Q6: Coding comfort
	6: {
		0: {{FeatTechAffinity, 0.00}},
		1: {{FeatTechAffinity, 0.25}},
		2: {{FeatTechAffinity, 0.55}},
		3: {{FeatTechAffinity, 0.80}},
		4: {{FeatTechAffinity, 1.00}},
	},
	// Q7: English confidence
	7: {
		0: {{FeatAcademicStrength, 0.20}, {FeatAbroadInterest, 0.05}},
		1: {{FeatAcademicStrength, 0.40}, {FeatAbroadInterest, 0.20}},
		2: {{FeatAcademicStrength, 0.70}, {FeatAbroadInterest, 0.55}},
		3: {{FeatAcademicStrength, 0.85}, {FeatAbroadInterest, 0.80}},
	},
	// Q8: Internships/work experience
	8: {
		0: {{FeatAcademicStrength, 0.20}, {FeatLeadershipScore, 0.10}},
		1: {{FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.30}},
		2: {{FeatAcademicStrength, 0.55}, {FeatLeadershipScore, 0.60}},
		3: {{FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.85}},
	},

	// ================================================================
	// SECTION B: FINANCIAL SITUATION (Q9–Q14)
	// ================================================================

	// Q9: Education loan
	9: {
		0: {{FeatFinancialPressure, 0.10}, {FeatIncomeUrgency, 0.10}},
		1: {{FeatFinancialPressure, 0.35}, {FeatIncomeUrgency, 0.30}},
		2: {{FeatFinancialPressure, 0.65}, {FeatIncomeUrgency, 0.60}},
		3: {{FeatFinancialPressure, 0.90}, {FeatIncomeUrgency, 0.85}},
	},
	// Q10: Family dependents
	10: {
		0: {{FeatFinancialPressure, 0.05}, {FeatIncomeUrgency, 0.05}},
		1: {{FeatFinancialPressure, 0.30}, {FeatIncomeUrgency, 0.25}},
		2: {{FeatFinancialPressure, 0.65}, {FeatIncomeUrgency, 0.60}},
		3: {{FeatFinancialPressure, 0.95}, {FeatIncomeUrgency, 0.90}},
	},
	// Q11: Income urgency timeline
	11: {
		0: {{FeatIncomeUrgency, 0.05}},
		1: {{FeatIncomeUrgency, 0.30}},
		2: {{FeatIncomeUrgency, 0.70}},
		3: {{FeatIncomeUrgency, 1.00}},
	},
	// Q12: MBA-level fees affordability
	12: {
		0: {{FeatFinancialPressure, 0.05}},
		1: {{FeatFinancialPressure, 0.30}},
		2: {{FeatFinancialPressure, 0.55}},
		3: {{FeatFinancialPressure, 0.85}},
	},
	// Q13: Family income
	13: {
		0: {{FeatFinancialPressure, 0.90}, {FeatIncomeUrgency, 0.80}},
		1: {{FeatFinancialPressure, 0.60}, {FeatIncomeUrgency, 0.55}},
		2: {{FeatFinancialPressure, 0.30}, {FeatIncomeUrgency, 0.25}},
		3: {{FeatFinancialPressure, 0.10}, {FeatIncomeUrgency, 0.10}},
		4: {{FeatFinancialPressure, 0.02}, {FeatIncomeUrgency, 0.02}},
	},
	// Q14: Location (where grew up)
	14: {
		0: {{FeatFinancialPressure, 0.40}, {FeatAbroadInterest, 0.05}},
		1: {{FeatFinancialPressure, 0.25}, {FeatAbroadInterest, 0.15}},
		2: {{FeatFinancialPressure, 0.10}, {FeatAbroadInterest, 0.35}},
		3: {{FeatFinancialPressure, 0.05}, {FeatAbroadInterest, 0.50}, {FeatLeadershipScore, 0.30}},
	},

	// ================================================================
	// SECTION C: PERSONALITY & RISK (Q15–Q22)
	// ================================================================

	// Q15: Risk appetite
	15: {
		0: {{FeatRiskTolerance, 0.05}, {FeatGovtInterest, 0.40}},
		1: {{FeatRiskTolerance, 0.20}, {FeatGovtInterest, 0.25}},
		2: {{FeatRiskTolerance, 0.50}},
		3: {{FeatRiskTolerance, 0.75}},
		4: {{FeatRiskTolerance, 1.00}},
	},
	// Q16: Leadership / managing teams
	16: {
		0: {{FeatLeadershipScore, 0.10}},
		1: {{FeatLeadershipScore, 0.35}},
		2: {{FeatLeadershipScore, 0.70}},
		3: {{FeatLeadershipScore, 1.00}},
	},
	// Q17: Security vs growth
	17: {
		0: {{FeatRiskTolerance, 0.05}, {FeatCareerInstability, 0.10}, {FeatGovtInterest, 0.50}},
		1: {{FeatRiskTolerance, 0.25}, {FeatCareerInstability, 0.25}},
		2: {{FeatRiskTolerance, 0.55}, {FeatCareerInstability, 0.50}},
		3: {{FeatRiskTolerance, 0.85}, {FeatCareerInstability, 0.80}},
	},
	// Q18: Handling pressure
	18: {
		0: {{FeatRiskTolerance, 0.10}, {FeatLeadershipScore, 0.05}},
		1: {{FeatRiskTolerance, 0.25}, {FeatLeadershipScore, 0.20}},
		2: {{FeatRiskTolerance, 0.55}, {FeatLeadershipScore, 0.50}},
		3: {{FeatRiskTolerance, 0.80}, {FeatLeadershipScore, 0.70}},
	},
	// Q19: Work-life balance
	19: {
		0: {{FeatGovtInterest, 0.40}, {FeatRiskTolerance, 0.10}},
		1: {{FeatRiskTolerance, 0.25}},
		2: {{FeatRiskTolerance, 0.55}, {FeatLeadershipScore, 0.35}},
		3: {{FeatRiskTolerance, 0.80}, {FeatLeadershipScore, 0.50}},
	},
	// Q20: Problem solving style
	20: {
		0: {{FeatGovtInterest, 0.35}, {FeatTechAffinity, 0.10}},
		1: {{FeatTechAffinity, 0.30}},
		2: {{FeatTechAffinity, 0.65}, {FeatAcademicStrength, 0.30}},
		3: {{FeatTechAffinity, 0.85}, {FeatAcademicStrength, 0.40}},
	},
	// Q21: Learning preference
	21: {
		0: {{FeatGovtInterest, 0.30}},
		1: {{FeatAcademicStrength, 0.20}},
		2: {{FeatAcademicStrength, 0.45}, {FeatTechAffinity, 0.30}},
		3: {{FeatAcademicStrength, 0.65}, {FeatTechAffinity, 0.50}},
	},
	// Q22: Public speaking
	22: {
		0: {{FeatLeadershipScore, 0.10}},
		1: {{FeatLeadershipScore, 0.30}},
		2: {{FeatLeadershipScore, 0.65}},
		3: {{FeatLeadershipScore, 0.90}},
	},

	// ================================================================
	// SECTION D: CAREER INTEREST (Q23–Q30)
	// ================================================================

	// Q23: Govt job appeal
	23: {
		0: {{FeatGovtInterest, 0.00}},
		1: {{FeatGovtInterest, 0.20}},
		2: {{FeatGovtInterest, 0.50}, {FeatCareerInstability, 0.25}},
		3: {{FeatGovtInterest, 0.80}, {FeatCareerInstability, 0.40}},
		4: {{FeatGovtInterest, 1.00}, {FeatCareerInstability, 0.55}},
	},
	// Q24: Entrepreneurship interest
	24: {
		0: {{FeatRiskTolerance, 0.05}},
		1: {{FeatRiskTolerance, 0.20}, {FeatLeadershipScore, 0.15}},
		2: {{FeatRiskTolerance, 0.50}, {FeatLeadershipScore, 0.40}, {FeatCareerInstability, 0.40}},
		3: {{FeatRiskTolerance, 0.75}, {FeatLeadershipScore, 0.65}, {FeatCareerInstability, 0.60}},
		4: {{FeatRiskTolerance, 0.90}, {FeatLeadershipScore, 0.80}, {FeatCareerInstability, 0.75}},
	},
	// Q25: MNC / big company interest
	25: {
		0: {{FeatTechAffinity, 0.05}},
		1: {{FeatTechAffinity, 0.25}},
		2: {{FeatTechAffinity, 0.55}, {FeatLeadershipScore, 0.20}},
		3: {{FeatTechAffinity, 0.75}, {FeatLeadershipScore, 0.35}},
	},
	// Q26: Higher studies in India
	26: {
		0: {{FeatAcademicStrength, 0.10}},
		1: {{FeatAcademicStrength, 0.30}},
		2: {{FeatAcademicStrength, 0.65}},
		3: {{FeatAcademicStrength, 0.90}},
	},
	// Q27: Abroad dream
	27: {
		0: {{FeatAbroadInterest, 0.00}, {FeatGovtInterest, 0.25}},
		1: {{FeatAbroadInterest, 0.30}},
		2: {{FeatAbroadInterest, 0.75}},
		3: {{FeatAbroadInterest, 1.00}},
	},
	// Q28: Salary expectations
	28: {
		0: {{FeatIncomeUrgency, 0.60}, {FeatFinancialPressure, 0.30}},
		1: {{FeatIncomeUrgency, 0.35}},
		2: {{FeatTechAffinity, 0.30}, {FeatLeadershipScore, 0.20}},
		3: {{FeatLeadershipScore, 0.45}, {FeatAbroadInterest, 0.25}},
		4: {{FeatAbroadInterest, 0.55}, {FeatLeadershipScore, 0.40}, {FeatTechAffinity, 0.35}},
	},
	// Q29: Field excitement
	29: {
		0: {{FeatTechAffinity, 0.90}},                                          // Technology
		1: {{FeatLeadershipScore, 0.60}, {FeatTechAffinity, 0.15}},             // Finance
		2: {{FeatGovtInterest, 0.85}},                                          // Government
		3: {{FeatAcademicStrength, 0.40}, {FeatAbroadInterest, 0.25}},          // Healthcare
		4: {{FeatAcademicStrength, 0.70}},                                      // Education/Research
	},
	// Q30: 10-year vision
	30: {
		0: {{FeatTechAffinity, 0.85}, {FeatAcademicStrength, 0.30}},                              // Senior engineer
		1: {{FeatLeadershipScore, 0.85}, {FeatTechAffinity, 0.15}},                               // Business leader
		2: {{FeatGovtInterest, 1.00}},                                                             // IAS/IPS
		3: {{FeatRiskTolerance, 0.80}, {FeatLeadershipScore, 0.75}, {FeatCareerInstability, 0.50}}, // Own company
		4: {{FeatAcademicStrength, 0.85}},                                                          // Professor
		5: {{FeatAbroadInterest, 0.95}, {FeatTechAffinity, 0.30}},                                // Abroad
	},
}

// AggregateProfile converts raw assessment answers into a structured UserProfile.
// It maps each answer to feature contributions, accumulates them, and normalizes
// each feature to 0–1 by dividing by the number of contributing questions.
func AggregateProfile(answers []dto.AnswerItem, questions []QuestionData) *UserProfile {
	profile := &UserProfile{}
	featureCounts := [NumFeatures]int{}

	for _, answer := range answers {
		// Find matching question by ID, get its DisplayOrder
		var displayOrder int
		found := false
		for _, q := range questions {
			if q.ID == answer.QuestionID {
				displayOrder = q.DisplayOrder
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Lookup feature mappings for this question + selected option
		optionMap, qExists := questionFeatureMap[displayOrder]
		if !qExists {
			continue
		}
		mappings, optExists := optionMap[answer.Selected]
		if !optExists {
			continue
		}

		// Accumulate feature contributions
		for _, m := range mappings {
			profile.Features[m.FeatureIndex] += m.Weight
			featureCounts[m.FeatureIndex]++
		}
	}

	// Normalize each feature by the number of contributing questions
	for i := 0; i < NumFeatures; i++ {
		if featureCounts[i] > 0 {
			profile.Features[i] /= float64(featureCounts[i])
		}
	}

	// Clamp all features to [0, 1]
	for i := 0; i < NumFeatures; i++ {
		if profile.Features[i] > 1.0 {
			profile.Features[i] = 1.0
		}
		if profile.Features[i] < 0.0 {
			profile.Features[i] = 0.0
		}
	}

	return profile
}
