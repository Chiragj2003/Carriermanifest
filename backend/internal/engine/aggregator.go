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
// ML-OPTIMIZED: Feature mappings refined based on Random Forest feature
// importance analysis (10,000 training samples, 88.45% accuracy).
//
// Feature importance ranking:
//   #1 GovtInterest (0.188) — MOST CRITICAL differentiator
//   #2 AbroadInterest (0.172) — CRITICAL for MS Abroad vs India paths
//   #3 TechAffinity (0.142) — key IT/Software signal
//   #4 LeadershipScore (0.136) — MBA/Startup differentiator
//   #5 CareerInstability (0.125) — Startup vs stable career signal
//   #6 RiskTolerance (0.115) — Risk profile
//   #7 AcademicStrength (0.076) — Less differentiating than expected
//   #8 IncomeUrgency (0.026) — Low differentiation power
//   #9 FinancialPressure (0.018) — Weakest differentiator
//
// Cross-feature signal extraction: Many questions now contribute to multiple
// high-importance features to improve classification accuracy.
var questionFeatureMap = map[int]map[int][]FeatureMapping{
	// ================================================================
	// SECTION A: ACADEMIC BACKGROUND (Q1–Q8)
	// ML insight: AcademicStrength rank #7, but TechAffinity #3 and
	// AbroadInterest #2 can be extracted from academic questions.
	// ================================================================

	// Q1: 10th board marks — strong academic signal + mild abroad (high scores → MS ready)
	1: {
		0: {{FeatAcademicStrength, 0.10}},
		1: {{FeatAcademicStrength, 0.30}},
		2: {{FeatAcademicStrength, 0.55}},
		3: {{FeatAcademicStrength, 0.80}, {FeatAbroadInterest, 0.15}},
		4: {{FeatAcademicStrength, 1.00}, {FeatAbroadInterest, 0.25}},
	},
	// Q2: Stream — PCM strongly signals TechAffinity + AbroadInterest (STEM for MS)
	2: {
		0: {{FeatTechAffinity, 0.90}, {FeatAcademicStrength, 0.70}, {FeatAbroadInterest, 0.30}},  // PCM → STEM abroad
		1: {{FeatTechAffinity, 0.30}, {FeatAcademicStrength, 0.60}, {FeatAbroadInterest, 0.20}},  // PCB
		2: {{FeatTechAffinity, 0.20}, {FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.40}}, // Commerce w/ Maths
		3: {{FeatTechAffinity, 0.10}, {FeatAcademicStrength, 0.35}, {FeatLeadershipScore, 0.30}}, // Commerce w/o Maths
		4: {{FeatTechAffinity, 0.05}, {FeatAcademicStrength, 0.30}, {FeatGovtInterest, 0.40}},    // Arts → stronger govt signal
	},
	// Q3: College major — CS/IT boosts abroad (most common MS Abroad major)
	3: {
		0: {{FeatTechAffinity, 1.00}, {FeatAcademicStrength, 0.70}, {FeatAbroadInterest, 0.35}},  // CS/IT → strong abroad signal
		1: {{FeatTechAffinity, 0.50}, {FeatAcademicStrength, 0.65}, {FeatAbroadInterest, 0.25}},  // Other Engg → mild abroad
		2: {{FeatTechAffinity, 0.10}, {FeatLeadershipScore, 0.50}, {FeatAcademicStrength, 0.45}}, // Commerce/BBA
		3: {{FeatTechAffinity, 0.20}, {FeatAcademicStrength, 0.55}, {FeatAbroadInterest, 0.15}},  // B.Sc → some abroad
		4: {{FeatTechAffinity, 0.05}, {FeatAcademicStrength, 0.30}, {FeatGovtInterest, 0.35}},    // BA/Other → govt signal
	},
	// Q4: College grades — high grades boost abroad readiness
	4: {
		0: {{FeatAcademicStrength, 0.15}},
		1: {{FeatAcademicStrength, 0.35}},
		2: {{FeatAcademicStrength, 0.55}},
		3: {{FeatAcademicStrength, 0.80}, {FeatAbroadInterest, 0.15}},
		4: {{FeatAcademicStrength, 1.00}, {FeatAbroadInterest, 0.25}},
	},
	// Q5: Competitive exams — UPSC/SSC strongly boosts GovtInterest (#1 feature)
	5: {
		0: {{FeatAcademicStrength, 0.20}},
		1: {{FeatAcademicStrength, 0.65}, {FeatTechAffinity, 0.50}},                             // JEE/NEET
		2: {{FeatAcademicStrength, 0.60}, {FeatLeadershipScore, 0.60}},                          // CAT/XAT
		3: {{FeatAcademicStrength, 0.75}, {FeatTechAffinity, 0.40}, {FeatAbroadInterest, 0.20}}, // GATE/NET → abroad research
		4: {{FeatAcademicStrength, 0.55}, {FeatGovtInterest, 0.90}},                             // UPSC/SSC → stronger govt signal
	},
	// Q6: Coding comfort — strong tech signal + mild career instability (more skills = more options)
	6: {
		0: {{FeatTechAffinity, 0.00}},
		1: {{FeatTechAffinity, 0.25}},
		2: {{FeatTechAffinity, 0.55}, {FeatCareerInstability, 0.10}},
		3: {{FeatTechAffinity, 0.80}, {FeatCareerInstability, 0.15}},
		4: {{FeatTechAffinity, 1.00}, {FeatCareerInstability, 0.20}, {FeatAbroadInterest, 0.15}}, // top coders → abroad option
	},
	// Q7: English confidence — critical for AbroadInterest (#2 feature)
	7: {
		0: {{FeatAcademicStrength, 0.20}, {FeatAbroadInterest, 0.05}},
		1: {{FeatAcademicStrength, 0.40}, {FeatAbroadInterest, 0.20}},
		2: {{FeatAcademicStrength, 0.70}, {FeatAbroadInterest, 0.60}}, // good English → strong abroad signal
		3: {{FeatAcademicStrength, 0.85}, {FeatAbroadInterest, 0.85}}, // excellent English → very strong abroad
	},
	// Q8: Internships — signals leadership + tech affinity
	8: {
		0: {{FeatAcademicStrength, 0.20}, {FeatLeadershipScore, 0.10}},
		1: {{FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.30}, {FeatTechAffinity, 0.20}},
		2: {{FeatAcademicStrength, 0.55}, {FeatLeadershipScore, 0.60}, {FeatTechAffinity, 0.30}},
		3: {{FeatAcademicStrength, 0.50}, {FeatLeadershipScore, 0.85}, {FeatTechAffinity, 0.35}},
	},

	// ================================================================
	// SECTION B: FINANCIAL SITUATION (Q9–Q14)
	// ML insight: FinancialPressure (#9) and IncomeUrgency (#8) are
	// WEAKEST differentiators. Extract cross-signals for GovtInterest
	// (stability-seeking) and inverse AbroadInterest.
	// ================================================================

	// Q9: Education loan — loan pressure → GovtInterest (stability), inverse AbroadInterest
	9: {
		0: {{FeatFinancialPressure, 0.10}, {FeatIncomeUrgency, 0.10}},
		1: {{FeatFinancialPressure, 0.35}, {FeatIncomeUrgency, 0.30}},
		2: {{FeatFinancialPressure, 0.65}, {FeatIncomeUrgency, 0.60}, {FeatGovtInterest, 0.15}},
		3: {{FeatFinancialPressure, 0.90}, {FeatIncomeUrgency, 0.85}, {FeatGovtInterest, 0.25}},
	},
	// Q10: Family dependents — high dependents → GovtInterest (stability)
	10: {
		0: {{FeatFinancialPressure, 0.05}, {FeatIncomeUrgency, 0.05}},
		1: {{FeatFinancialPressure, 0.30}, {FeatIncomeUrgency, 0.25}},
		2: {{FeatFinancialPressure, 0.65}, {FeatIncomeUrgency, 0.60}, {FeatGovtInterest, 0.20}},
		3: {{FeatFinancialPressure, 0.95}, {FeatIncomeUrgency, 0.90}, {FeatGovtInterest, 0.30}},
	},
	// Q11: Income urgency timeline — urgent income inversely signals AbroadInterest
	11: {
		0: {{FeatIncomeUrgency, 0.05}, {FeatAbroadInterest, 0.20}}, // no rush → open to abroad
		1: {{FeatIncomeUrgency, 0.30}},
		2: {{FeatIncomeUrgency, 0.70}, {FeatGovtInterest, 0.15}},
		3: {{FeatIncomeUrgency, 1.00}, {FeatGovtInterest, 0.20}}, // very urgent → stability
	},
	// Q12: MBA-level fees affordability
	12: {
		0: {{FeatFinancialPressure, 0.05}},
		1: {{FeatFinancialPressure, 0.30}},
		2: {{FeatFinancialPressure, 0.55}},
		3: {{FeatFinancialPressure, 0.85}},
	},
	// Q13: Family income — low income boosts GovtInterest (stability)
	13: {
		0: {{FeatFinancialPressure, 0.90}, {FeatIncomeUrgency, 0.80}, {FeatGovtInterest, 0.25}},
		1: {{FeatFinancialPressure, 0.60}, {FeatIncomeUrgency, 0.55}, {FeatGovtInterest, 0.15}},
		2: {{FeatFinancialPressure, 0.30}, {FeatIncomeUrgency, 0.25}},
		3: {{FeatFinancialPressure, 0.10}, {FeatIncomeUrgency, 0.10}, {FeatAbroadInterest, 0.15}},
		4: {{FeatFinancialPressure, 0.02}, {FeatIncomeUrgency, 0.02}, {FeatAbroadInterest, 0.25}}, // wealthy → abroad viable
	},
	// Q14: Location — metro signals TechAffinity + AbroadInterest
	14: {
		0: {{FeatFinancialPressure, 0.40}, {FeatAbroadInterest, 0.05}, {FeatGovtInterest, 0.20}},                              // village → govt
		1: {{FeatFinancialPressure, 0.25}, {FeatAbroadInterest, 0.15}, {FeatGovtInterest, 0.10}},                              // small town
		2: {{FeatFinancialPressure, 0.10}, {FeatAbroadInterest, 0.35}, {FeatTechAffinity, 0.20}},                              // city → tech + abroad
		3: {{FeatFinancialPressure, 0.05}, {FeatAbroadInterest, 0.55}, {FeatLeadershipScore, 0.30}, {FeatTechAffinity, 0.30}}, // metro → strong tech + abroad
	},

	// ================================================================
	// SECTION C: PERSONALITY & RISK (Q15–Q22)
	// ML insight: Extract GovtInterest (#1) from safety-oriented answers,
	// AbroadInterest (#2) from openness/creativity answers.
	// ================================================================

	// Q15: Risk appetite — safety-first strongly boosts GovtInterest
	15: {
		0: {{FeatRiskTolerance, 0.05}, {FeatGovtInterest, 0.55}}, // safety first → strong govt signal
		1: {{FeatRiskTolerance, 0.20}, {FeatGovtInterest, 0.35}},
		2: {{FeatRiskTolerance, 0.50}},
		3: {{FeatRiskTolerance, 0.75}, {FeatCareerInstability, 0.20}},
		4: {{FeatRiskTolerance, 1.00}, {FeatCareerInstability, 0.35}}, // risk-lover → instability ok
	},
	// Q16: Leadership — strong leadership signals MBA/Startup, add TechAffinity for tech leadership
	16: {
		0: {{FeatLeadershipScore, 0.10}},
		1: {{FeatLeadershipScore, 0.35}},
		2: {{FeatLeadershipScore, 0.70}, {FeatCareerInstability, 0.15}},
		3: {{FeatLeadershipScore, 1.00}, {FeatCareerInstability, 0.25}}, // strong leaders tolerate instability
	},
	// Q17: Security vs growth — CRITICAL for GovtInterest (#1 feature)
	17: {
		0: {{FeatRiskTolerance, 0.05}, {FeatCareerInstability, 0.05}, {FeatGovtInterest, 0.60}}, // security → very strong govt
		1: {{FeatRiskTolerance, 0.25}, {FeatCareerInstability, 0.25}, {FeatGovtInterest, 0.25}},
		2: {{FeatRiskTolerance, 0.55}, {FeatCareerInstability, 0.50}, {FeatAbroadInterest, 0.15}},
		3: {{FeatRiskTolerance, 0.85}, {FeatCareerInstability, 0.80}, {FeatAbroadInterest, 0.20}}, // growth → abroad mindset
	},
	// Q18: Handling pressure — high pressure tolerance → leadership + risk
	18: {
		0: {{FeatRiskTolerance, 0.10}, {FeatLeadershipScore, 0.05}},
		1: {{FeatRiskTolerance, 0.25}, {FeatLeadershipScore, 0.20}},
		2: {{FeatRiskTolerance, 0.55}, {FeatLeadershipScore, 0.50}},
		3: {{FeatRiskTolerance, 0.80}, {FeatLeadershipScore, 0.70}, {FeatCareerInstability, 0.20}},
	},
	// Q19: Work-life balance — strong WLB preference → GovtInterest
	19: {
		0: {{FeatGovtInterest, 0.50}, {FeatRiskTolerance, 0.05}}, // WLB priority → strong govt
		1: {{FeatRiskTolerance, 0.25}, {FeatGovtInterest, 0.20}},
		2: {{FeatRiskTolerance, 0.55}, {FeatLeadershipScore, 0.35}},
		3: {{FeatRiskTolerance, 0.80}, {FeatLeadershipScore, 0.50}, {FeatCareerInstability, 0.30}},
	},
	// Q20: Problem solving — innovation preference signals TechAffinity + AbroadInterest
	20: {
		0: {{FeatGovtInterest, 0.40}, {FeatTechAffinity, 0.05}}, // proven methods → govt mindset
		1: {{FeatTechAffinity, 0.30}},
		2: {{FeatTechAffinity, 0.65}, {FeatAcademicStrength, 0.30}, {FeatAbroadInterest, 0.15}},
		3: {{FeatTechAffinity, 0.85}, {FeatAcademicStrength, 0.40}, {FeatAbroadInterest, 0.25}}, // innovators → abroad
	},
	// Q21: Learning preference — love of learning → academic + abroad
	21: {
		0: {{FeatGovtInterest, 0.30}},
		1: {{FeatAcademicStrength, 0.20}},
		2: {{FeatAcademicStrength, 0.45}, {FeatTechAffinity, 0.30}},
		3: {{FeatAcademicStrength, 0.65}, {FeatTechAffinity, 0.50}, {FeatAbroadInterest, 0.20}}, // lifelong learner → abroad
	},
	// Q22: Public speaking — strong ability signals leadership + MBA path
	22: {
		0: {{FeatLeadershipScore, 0.10}},
		1: {{FeatLeadershipScore, 0.30}},
		2: {{FeatLeadershipScore, 0.65}, {FeatGovtInterest, 0.10}}, // good speaker → govt/admin
		3: {{FeatLeadershipScore, 0.90}, {FeatGovtInterest, 0.15}},
	},

	// ================================================================
	// SECTION D: CAREER INTEREST (Q23–Q30)
	// ML insight: Direct career interest questions are the STRONGEST
	// differentiators. Maximize cross-feature extraction here.
	// ================================================================

	// Q23: Govt job appeal — HIGHEST importance feature (#1)
	23: {
		0: {{FeatGovtInterest, 0.00}, {FeatRiskTolerance, 0.20}}, // no interest → risk tolerant
		1: {{FeatGovtInterest, 0.20}},
		2: {{FeatGovtInterest, 0.55}, {FeatCareerInstability, 0.15}},
		3: {{FeatGovtInterest, 0.80}, {FeatCareerInstability, 0.10}},
		4: {{FeatGovtInterest, 1.00}}, // die-hard govt → pure signal
	},
	// Q24: Entrepreneurship — high interest = TechAffinity + Risk + Instability
	24: {
		0: {{FeatRiskTolerance, 0.05}, {FeatGovtInterest, 0.15}}, // no interest → mild govt
		1: {{FeatRiskTolerance, 0.20}, {FeatLeadershipScore, 0.15}},
		2: {{FeatRiskTolerance, 0.50}, {FeatLeadershipScore, 0.40}, {FeatCareerInstability, 0.40}, {FeatTechAffinity, 0.15}},
		3: {{FeatRiskTolerance, 0.75}, {FeatLeadershipScore, 0.65}, {FeatCareerInstability, 0.60}, {FeatTechAffinity, 0.25}},
		4: {{FeatRiskTolerance, 0.90}, {FeatLeadershipScore, 0.80}, {FeatCareerInstability, 0.75}, {FeatTechAffinity, 0.30}},
	},
	// Q25: MNC / big company — also signals AbroadInterest (MNCs are global)
	25: {
		0: {{FeatTechAffinity, 0.05}},
		1: {{FeatTechAffinity, 0.25}},
		2: {{FeatTechAffinity, 0.55}, {FeatLeadershipScore, 0.20}, {FeatAbroadInterest, 0.20}},
		3: {{FeatTechAffinity, 0.75}, {FeatLeadershipScore, 0.35}, {FeatAbroadInterest, 0.30}}, // MNC lover → abroad
	},
	// Q26: Higher studies in India — academic + mild inverse abroad
	26: {
		0: {{FeatAcademicStrength, 0.10}},
		1: {{FeatAcademicStrength, 0.30}},
		2: {{FeatAcademicStrength, 0.65}, {FeatGovtInterest, 0.10}}, // research in India → mild govt
		3: {{FeatAcademicStrength, 0.90}, {FeatGovtInterest, 0.15}}, // strong India studies → govt/academic path
	},
	// Q27: Abroad dream — #2 MOST IMPORTANT feature
	27: {
		0: {{FeatAbroadInterest, 0.00}, {FeatGovtInterest, 0.30}}, // no abroad → stronger govt signal
		1: {{FeatAbroadInterest, 0.30}},
		2: {{FeatAbroadInterest, 0.75}, {FeatTechAffinity, 0.15}}, // abroad leaning → often tech
		3: {{FeatAbroadInterest, 1.00}, {FeatTechAffinity, 0.20}}, // dream abroad → strong signal
	},
	// Q28: Salary expectations — high salary → tech/abroad, moderate → govt
	28: {
		0: {{FeatIncomeUrgency, 0.60}, {FeatFinancialPressure, 0.30}, {FeatGovtInterest, 0.20}}, // low salary ok → govt
		1: {{FeatIncomeUrgency, 0.35}, {FeatGovtInterest, 0.15}},
		2: {{FeatTechAffinity, 0.30}, {FeatLeadershipScore, 0.20}},
		3: {{FeatLeadershipScore, 0.45}, {FeatAbroadInterest, 0.30}, {FeatTechAffinity, 0.20}},
		4: {{FeatAbroadInterest, 0.55}, {FeatLeadershipScore, 0.40}, {FeatTechAffinity, 0.35}}, // high salary → abroad/tech
	},
	// Q29: Field excitement — STRONG direct career signals
	29: {
		0: {{FeatTechAffinity, 0.95}},                                 // Technology
		1: {{FeatLeadershipScore, 0.65}, {FeatTechAffinity, 0.15}},    // Finance
		2: {{FeatGovtInterest, 0.90}},                                 // Government → near-max
		3: {{FeatAcademicStrength, 0.40}, {FeatAbroadInterest, 0.30}}, // Healthcare → some abroad
		4: {{FeatAcademicStrength, 0.70}, {FeatGovtInterest, 0.15}},   // Education/Research → mild govt
	},
	// Q30: 10-year vision — STRONGEST direct career signals
	30: {
		0: {{FeatTechAffinity, 0.90}, {FeatAcademicStrength, 0.30}},                                // Senior engineer
		1: {{FeatLeadershipScore, 0.85}, {FeatTechAffinity, 0.15}},                                 // Business leader
		2: {{FeatGovtInterest, 1.00}},                                                              // IAS/IPS
		3: {{FeatRiskTolerance, 0.85}, {FeatLeadershipScore, 0.75}, {FeatCareerInstability, 0.55}}, // Own company
		4: {{FeatAcademicStrength, 0.85}, {FeatGovtInterest, 0.15}},                                // Professor → mild govt
		5: {{FeatAbroadInterest, 0.95}, {FeatTechAffinity, 0.30}},                                  // Abroad
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
