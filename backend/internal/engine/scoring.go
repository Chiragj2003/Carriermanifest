// Package engine implements the rule-based career scoring system.
// This is the core AI system of CareerManifest - a weighted scoring engine
// that evaluates career paths based on Indian student profiles.
package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/careermanifest/backend/internal/dto"
)

// Career category constants
const (
	CareerIT         = "IT / Software Jobs"
	CareerMBA        = "MBA (India)"
	CareerGovt       = "Government Exams"
	CareerStartup    = "Startup / Entrepreneurship"
	CareerHigherIndia = "Higher Studies (India)"
	CareerMSAbroad   = "MS Abroad"
)

// AllCareers is the list of all career categories.
var AllCareers = []string{
	CareerIT,
	CareerMBA,
	CareerGovt,
	CareerStartup,
	CareerHigherIndia,
	CareerMSAbroad,
}

// ScoringEngine evaluates assessment answers and produces career recommendations.
type ScoringEngine struct{}

// NewScoringEngine creates a new ScoringEngine.
func NewScoringEngine() *ScoringEngine {
	return &ScoringEngine{}
}

// ParsedAnswer holds a single parsed answer with its weight information.
type ParsedAnswer struct {
	QuestionID uint64
	Selected   int
	Weights    dto.QuestionWeight
	Category   string
}

// ComputeResult processes all answers against question weights and produces a full result.
func (e *ScoringEngine) ComputeResult(answers []dto.AnswerItem, questionsJSON []QuestionData) (*dto.AssessmentResult, error) {
	// Step 1: Accumulate scores per career category
	scores := make(map[string]float64)
	maxScores := make(map[string]float64)
	for _, career := range AllCareers {
		scores[career] = 0
		maxScores[career] = 0
	}

	// Risk factor accumulators
	riskFactors := map[string]float64{
		"income_urgency":        0,
		"family_dependency":     0,
		"risk_tolerance":        0,
		"career_instability":    0,
	}
	riskFactorCounts := map[string]int{
		"income_urgency":     0,
		"family_dependency":  0,
		"risk_tolerance":     0,
		"career_instability": 0,
	}

	// Step 2: Process each answer
	for _, answer := range answers {
		// Find the matching question data
		var qData *QuestionData
		for i := range questionsJSON {
			if questionsJSON[i].ID == answer.QuestionID {
				qData = &questionsJSON[i]
				break
			}
		}
		if qData == nil {
			continue // Skip unknown questions
		}

		// Find the weight entry for the selected option
		for _, w := range qData.Weights {
			if w.OptionIndex == answer.Selected {
				// Add career scores
				for career, score := range w.Scores {
					scores[career] += score
				}
				// Add risk factors
				for factor, value := range w.RiskFactors {
					riskFactors[factor] += value
					riskFactorCounts[factor]++
				}
				break
			}
		}

		// Calculate max possible score for each career from this question
		for _, w := range qData.Weights {
			for career, score := range w.Scores {
				if score > maxScores[career] {
					// Track the max among this question's options
				}
				_ = score // We'll calculate max differently
			}
			_ = w
		}
	}

	// Calculate max scores properly: sum of maximum possible score per question per career
	for _, qData := range questionsJSON {
		for _, career := range AllCareers {
			maxForQuestion := 0.0
			for _, w := range qData.Weights {
				if s, ok := w.Scores[career]; ok && s > maxForQuestion {
					maxForQuestion = s
				}
			}
			maxScores[career] += maxForQuestion
		}
	}

	// Step 3: Build career scores with percentages
	var careerScores []dto.CareerScore
	for _, career := range AllCareers {
		maxScore := maxScores[career]
		if maxScore == 0 {
			maxScore = 1 // Prevent division by zero
		}
		percentage := (scores[career] / maxScore) * 100
		careerScores = append(careerScores, dto.CareerScore{
			Category:   career,
			Score:      math.Round(scores[career]*100) / 100,
			MaxScore:   math.Round(maxScore*100) / 100,
			Percentage: math.Round(percentage*100) / 100,
		})
	}

	// Sort by percentage descending
	sort.Slice(careerScores, func(i, j int) bool {
		return careerScores[i].Percentage > careerScores[j].Percentage
	})

	bestCareer := careerScores[0].Category

	// Step 4: Calculate risk using India-realistic formula
	// Normalize risk factors to 0-10 scale
	for factor := range riskFactors {
		if riskFactorCounts[factor] > 0 {
			riskFactors[factor] = riskFactors[factor] / float64(riskFactorCounts[factor])
		}
	}

	// RiskScore = (IncomeUrgency × 0.35) + (FamilyDependency × 0.25) + (RiskTolerance × 0.20) + (CareerInstabilityIndex × 0.20)
	riskScore := (riskFactors["income_urgency"] * 0.35) +
		(riskFactors["family_dependency"] * 0.25) +
		(riskFactors["risk_tolerance"] * 0.20) +
		(riskFactors["career_instability"] * 0.20)

	riskScore = math.Round(riskScore*100) / 100

	riskLevel := "Low"
	if riskScore > 6 {
		riskLevel = "High"
	} else if riskScore > 3 {
		riskLevel = "Medium"
	}

	risk := dto.RiskAssessment{
		Score:   riskScore,
		Level:   riskLevel,
		Factors: riskFactors,
	}

	// Step 5: Generate salary projection for best career
	salaryProjection := getSalaryProjection(bestCareer)

	// Step 6: Generate roadmap for best career
	roadmap := getRoadmap(bestCareer)

	// Step 7: Get required skills
	skills := getRequiredSkills(bestCareer)

	// Step 8: Get suggested exams and colleges
	exams := getSuggestedExams(bestCareer)
	colleges := getSuggestedColleges(bestCareer)

	return &dto.AssessmentResult{
		Scores:            careerScores,
		BestCareerPath:    bestCareer,
		Risk:              risk,
		SalaryProjection:  salaryProjection,
		Roadmap:           roadmap,
		RequiredSkills:    skills,
		SuggestedExams:    exams,
		SuggestedColleges: colleges,
	}, nil
}

// QuestionData is a simplified question structure for the engine.
type QuestionData struct {
	ID       uint64             `json:"id"`
	Category string             `json:"category"`
	Weights  []dto.QuestionWeight `json:"weights"`
}

// ParseQuestionWeights parses JSON weight data from the database.
func ParseQuestionWeights(weightsJSON string) ([]dto.QuestionWeight, error) {
	var weights []dto.QuestionWeight
	if err := json.Unmarshal([]byte(weightsJSON), &weights); err != nil {
		return nil, fmt.Errorf("failed to parse weights: %w", err)
	}
	return weights, nil
}

// ============================================================
// CAREER-SPECIFIC DATA (India-focused, realistic)
// ============================================================

func getSalaryProjection(career string) dto.SalaryProjection {
	projections := map[string]dto.SalaryProjection{
		CareerIT: {
			Year1: "₹4-8 LPA", Year2: "₹6-12 LPA", Year3: "₹10-18 LPA",
			Year4: "₹14-25 LPA", Year5: "₹18-35 LPA",
		},
		CareerMBA: {
			Year1: "₹8-15 LPA", Year2: "₹10-20 LPA", Year3: "₹14-28 LPA",
			Year4: "₹18-35 LPA", Year5: "₹22-50 LPA",
		},
		CareerGovt: {
			Year1: "₹5-8 LPA", Year2: "₹5.5-9 LPA", Year3: "₹6-10 LPA",
			Year4: "₹7-12 LPA", Year5: "₹8-15 LPA",
		},
		CareerStartup: {
			Year1: "₹0-5 LPA", Year2: "₹0-10 LPA", Year3: "₹5-20 LPA",
			Year4: "₹10-40 LPA", Year5: "₹15-100+ LPA",
		},
		CareerHigherIndia: {
			Year1: "₹0 (Stipend ₹30-50K/mo)", Year2: "₹0 (Stipend ₹35-60K/mo)", Year3: "₹8-15 LPA",
			Year4: "₹10-20 LPA", Year5: "₹14-30 LPA",
		},
		CareerMSAbroad: {
			Year1: "$0 (Studying)", Year2: "$60-90K/year", Year3: "$75-120K/year",
			Year4: "$90-150K/year", Year5: "$100-180K/year",
		},
	}

	if p, ok := projections[career]; ok {
		return p
	}
	return projections[CareerIT]
}

func getRoadmap(career string) []dto.RoadmapStep {
	roadmaps := map[string][]dto.RoadmapStep{
		CareerIT: {
			{Step: 1, Title: "Learn Programming Fundamentals", Description: "Master one language (Python/Java/JavaScript). Complete DSA basics on LeetCode/GeeksForGeeks.", Duration: "3 months"},
			{Step: 2, Title: "Build Projects & Portfolio", Description: "Build 3-5 real projects. Create GitHub portfolio. Learn Git, APIs, databases.", Duration: "3 months"},
			{Step: 3, Title: "Learn Frameworks & Tools", Description: "Pick a stack (MERN/Spring Boot/Django). Learn Docker, cloud basics (AWS/GCP).", Duration: "2 months"},
			{Step: 4, Title: "DSA & Interview Prep", Description: "Solve 200+ LeetCode problems. Practice system design. Mock interviews.", Duration: "3 months"},
			{Step: 5, Title: "Apply & Network", Description: "Apply on LinkedIn, Naukri, AngelList. Attend hackathons. Get referrals.", Duration: "1 month"},
		},
		CareerMBA: {
			{Step: 1, Title: "CAT/XAT/GMAT Preparation", Description: "Join coaching (IMS/TIME/CL) or self-study. Target 95+ percentile in CAT.", Duration: "6-8 months"},
			{Step: 2, Title: "Build Profile", Description: "Gain 2-3 years work experience. Get leadership roles. Volunteer work.", Duration: "Ongoing"},
			{Step: 3, Title: "Application & Essays", Description: "Research IIMs, XLRI, ISB, FMS. Write compelling SOPs and essays.", Duration: "2 months"},
			{Step: 4, Title: "GD/PI Preparation", Description: "Current affairs, case studies, mock GDs and PIs.", Duration: "2 months"},
			{Step: 5, Title: "Specialization Planning", Description: "Research Finance, Marketing, Operations, HR tracks. Network with alumni.", Duration: "1 month"},
		},
		CareerGovt: {
			{Step: 1, Title: "Choose Your Exam", Description: "UPSC CSE, SSC CGL, Banking (IBPS/SBI), State PSC, Railways. Pick based on your eligibility.", Duration: "1 month"},
			{Step: 2, Title: "Foundation Building", Description: "NCERT books (6-12), basic GK, aptitude. Join coaching if needed (Unacademy/BYJU's).", Duration: "3 months"},
			{Step: 3, Title: "Subject Deep Dive", Description: "Cover full syllabus. Make notes. Previous year papers analysis.", Duration: "6 months"},
			{Step: 4, Title: "Test Series & Revision", Description: "Join test series. Weekly full-length mocks. Analyze mistakes.", Duration: "3 months"},
			{Step: 5, Title: "Prelims → Mains → Interview", Description: "Clear each stage. Personality test prep for UPSC. Document verification.", Duration: "6-12 months"},
		},
		CareerStartup: {
			{Step: 1, Title: "Ideation & Validation", Description: "Identify problems worth solving. Talk to 50+ potential customers. Validate demand.", Duration: "2 months"},
			{Step: 2, Title: "MVP Development", Description: "Build minimum viable product. Use no-code tools if needed. Get first 10 users.", Duration: "2 months"},
			{Step: 3, Title: "Early Traction", Description: "Get to 100+ users. Iterate based on feedback. Find product-market fit.", Duration: "3 months"},
			{Step: 4, Title: "Funding & Team", Description: "Apply to incubators (IIT, NSRCEL, T-Hub). Pitch to angels. Build core team.", Duration: "3 months"},
			{Step: 5, Title: "Scale & Growth", Description: "Optimize unit economics. Hiring. Series A preparation. Scale marketing.", Duration: "6 months"},
		},
		CareerHigherIndia: {
			{Step: 1, Title: "Choose Exam & Specialization", Description: "GATE, NET, JAM, or direct admission. Pick M.Tech/M.Sc/PhD path.", Duration: "1 month"},
			{Step: 2, Title: "Exam Preparation", Description: "GATE: Focus on core subjects + aptitude. Target AIR under 500 for IITs.", Duration: "6 months"},
			{Step: 3, Title: "College Selection", Description: "Research IITs, IISc, NITs, IIITs. Check placement records and research labs.", Duration: "1 month"},
			{Step: 4, Title: "Research & Thesis", Description: "Choose research area. Publish papers. Build academic network.", Duration: "12-18 months"},
			{Step: 5, Title: "Placement/PhD Application", Description: "Campus placements or apply for PhD positions. Build research profile.", Duration: "3 months"},
		},
		CareerMSAbroad: {
			{Step: 1, Title: "GRE & TOEFL/IELTS Prep", Description: "Target GRE 320+, TOEFL 100+ or IELTS 7.5+. Use Magoosh/ETS material.", Duration: "3 months"},
			{Step: 2, Title: "University Shortlisting", Description: "Research universities (US/Canada/Germany/UK). Check admit chances on Yocket/Admits.fyi.", Duration: "2 months"},
			{Step: 3, Title: "SOP, LORs & Application", Description: "Write compelling SOPs. Get 3 strong LORs. Apply to 8-12 universities.", Duration: "3 months"},
			{Step: 4, Title: "Funding & Visa", Description: "Apply for scholarships, TA/RA positions. Education loan. F1/student visa.", Duration: "3 months"},
			{Step: 5, Title: "Pre-Departure", Description: "Housing, bank account, health insurance. Connect with seniors at target university.", Duration: "2 months"},
		},
	}

	if r, ok := roadmaps[career]; ok {
		return r
	}
	return roadmaps[CareerIT]
}

func getRequiredSkills(career string) []string {
	skills := map[string][]string{
		CareerIT: {
			"Data Structures & Algorithms", "Programming (Python/Java/JS)",
			"Web Development (React/Node)", "Database Management (SQL/NoSQL)",
			"System Design", "Cloud Computing (AWS/GCP)",
			"Version Control (Git)", "Problem Solving",
		},
		CareerMBA: {
			"Quantitative Aptitude", "Verbal Ability & Reading Comprehension",
			"Logical Reasoning", "Data Interpretation",
			"Leadership & Teamwork", "Communication Skills",
			"Business Acumen", "Current Affairs",
		},
		CareerGovt: {
			"General Knowledge & Current Affairs", "Quantitative Aptitude",
			"English Language", "Logical Reasoning",
			"Indian Polity & Constitution", "Indian Economy",
			"History & Geography", "Essay Writing",
		},
		CareerStartup: {
			"Product Thinking", "Sales & Marketing",
			"Financial Planning", "Leadership & Team Building",
			"Technical Skills (Full-Stack/No-Code)", "Fundraising & Pitching",
			"Customer Development", "Growth Hacking",
		},
		CareerHigherIndia: {
			"Core Subject Expertise", "Research Methodology",
			"Academic Writing", "GATE/NET Exam Skills",
			"Programming (for CS/IT)", "Lab Work & Experimentation",
			"Paper Reading & Review", "Presentation Skills",
		},
		CareerMSAbroad: {
			"GRE Verbal & Quant", "TOEFL/IELTS English Proficiency",
			"Research Experience", "Academic Writing (SOP)",
			"Core Domain Knowledge", "Programming & Tools",
			"Networking & Communication", "Cross-Cultural Adaptability",
		},
	}

	if s, ok := skills[career]; ok {
		return s
	}
	return skills[CareerIT]
}

func getSuggestedExams(career string) []string {
	exams := map[string][]string{
		CareerIT:          {"GATE CS", "Google Kickstart", "CodeChef/Codeforces", "AWS Certification", "Company-specific OAs"},
		CareerMBA:         {"CAT", "XAT", "GMAT", "NMAT", "SNAP", "IIFT"},
		CareerGovt:        {"UPSC CSE", "SSC CGL", "IBPS PO", "SBI PO", "RBI Grade B", "State PSC"},
		CareerStartup:     {"No specific exams - focus on building", "Y Combinator Application", "Shark Tank India (if applicable)"},
		CareerHigherIndia: {"GATE", "UGC NET", "CSIR NET", "IIT JAM", "JEST"},
		CareerMSAbroad:    {"GRE General", "TOEFL iBT", "IELTS Academic", "GRE Subject (optional)"},
	}

	if e, ok := exams[career]; ok {
		return e
	}
	return exams[CareerIT]
}

func getSuggestedColleges(career string) []string {
	colleges := map[string][]string{
		CareerIT: {
			"IIT Bombay/Delhi/Madras (B.Tech/M.Tech)",
			"NIT Trichy/Warangal/Surathkal",
			"IIIT Hyderabad / BITS Pilani",
			"Top product companies (Google, Microsoft, Amazon)",
		},
		CareerMBA: {
			"IIM Ahmedabad / Bangalore / Calcutta",
			"IIM Lucknow / Indore / Kozhikode",
			"XLRI Jamshedpur / FMS Delhi",
			"ISB Hyderabad / IIM Udaipur (1-year)",
		},
		CareerGovt: {
			"LBSNAA (IAS Training)", "SVPNPA (IPS Training)",
			"Reserve Bank of India", "State Administrative Services",
		},
		CareerStartup: {
			"IIT/IIM Incubators", "NSRCEL (IIM Bangalore)",
			"T-Hub Hyderabad", "Startup India Hub",
			"Y Combinator / Techstars (global)",
		},
		CareerHigherIndia: {
			"IISc Bangalore", "IIT Bombay/Delhi/Madras/Kanpur",
			"TIFR / ISI Kolkata", "JNU / Delhi University",
		},
		CareerMSAbroad: {
			"MIT / Stanford / CMU (US)",
			"UC Berkeley / Georgia Tech / UIUC",
			"ETH Zurich / TU Munich (Europe)",
			"University of Toronto / UBC (Canada)",
		},
	}

	if c, ok := colleges[career]; ok {
		return c
	}
	return colleges[CareerIT]
}
