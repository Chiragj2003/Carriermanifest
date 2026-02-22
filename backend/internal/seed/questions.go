// Package seed contains the 30 India-focused assessment questions with weighted scoring.
// Each question maps options to scores across 6 career categories + risk factors.
package seed

import (
	"encoding/json"
	"log"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/repository"
)

// QuestionSeed holds a question for seeding.
type QuestionSeed struct {
	Category     string
	QuestionText string
	Options      []dto.QuestionOption
	Weights      []dto.QuestionWeight
	DisplayOrder int
}

// SeedQuestions inserts all 30 questions into the database if none exist.
func SeedQuestions(repo *repository.QuestionRepository) error {
	count, err := repo.CountQuestions()
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Questions already seeded (%d found), skipping", count)
		return nil
	}

	questions := getAllQuestions()

	for _, q := range questions {
		optionsJSON, _ := json.Marshal(q.Options)
		weightsJSON, _ := json.Marshal(q.Weights)
		_, err := repo.Create(q.Category, q.QuestionText, string(optionsJSON), string(weightsJSON), q.DisplayOrder)
		if err != nil {
			log.Printf("Warning: Failed to seed question '%s': %v", q.QuestionText, err)
		}
	}

	log.Printf("✅ Seeded %d assessment questions", len(questions))
	return nil
}

func getAllQuestions() []QuestionSeed {
	return []QuestionSeed{
		// ============================================================
		// SECTION A: ACADEMIC BACKGROUND (Questions 1–8)
		// ============================================================
		{
			Category:     "academic",
			QuestionText: "How did you perform in your 10th board exams?",
			Options: []dto.QuestionOption{
				{Label: "Below 60%", Value: 0},
				{Label: "60-75%", Value: 1},
				{Label: "75-85%", Value: 2},
				{Label: "85-95%", Value: 3},
				{Label: "95%+", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 1, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 1, "MS Abroad": 0}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 1}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 5, "MS Abroad": 5}},
			},
			DisplayOrder: 1,
		},
		{
			Category:     "academic",
			QuestionText: "What stream did you pick in 11th/12th?",
			Options: []dto.QuestionOption{
				{Label: "PCM (Physics, Chemistry, Maths)", Value: 0},
				{Label: "PCB (Physics, Chemistry, Biology)", Value: 1},
				{Label: "Commerce (with Maths)", Value: 2},
				{Label: "Commerce (without Maths)", Value: 3},
				{Label: "Arts / Humanities", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 5, "MS Abroad": 5}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 1}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 1}},
			},
			DisplayOrder: 2,
		},
		{
			Category:     "academic",
			QuestionText: "What are you studying (or did you study) in college?",
			Options: []dto.QuestionOption{
				{Label: "Computer Science / IT", Value: 0},
				{Label: "Other Engineering (Mech/Civil/ECE/EEE)", Value: 1},
				{Label: "B.Com / BBA / Economics", Value: 2},
				{Label: "B.Sc (Science)", Value: 3},
				{Label: "BA / Other", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 5}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 5, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 1}},
			},
			DisplayOrder: 3,
		},
		{
			Category:     "academic",
			QuestionText: "How are your college grades looking?",
			Options: []dto.QuestionOption{
				{Label: "Below 6.0 / 55%", Value: 0},
				{Label: "6.0-7.0 / 55-65%", Value: 1},
				{Label: "7.0-8.0 / 65-75%", Value: 2},
				{Label: "8.0-9.0 / 75-85%", Value: 3},
				{Label: "9.0+ / 85%+", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 1, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 1, "MS Abroad": 0}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 2, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 1}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 5, "MS Abroad": 5}},
			},
			DisplayOrder: 4,
		},
		{
			Category:     "academic",
			QuestionText: "Have you taken any competitive exams so far?",
			Options: []dto.QuestionOption{
				{Label: "No, never", Value: 0},
				{Label: "JEE / NEET (engineering/medical entrance)", Value: 1},
				{Label: "CAT / XAT / GMAT (MBA entrance)", Value: 2},
				{Label: "GATE / NET / JAM (higher studies)", Value: 3},
				{Label: "UPSC / SSC / Banking exams", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 2, "Government Exams": 1, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 2, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 2}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 2, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 5, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 1}},
			},
			DisplayOrder: 5,
		},
		{
			Category:     "academic",
			QuestionText: "How comfortable are you with coding and tech?",
			Options: []dto.QuestionOption{
				{Label: "No coding knowledge", Value: 0},
				{Label: "Basic (know one language)", Value: 1},
				{Label: "Intermediate (built small projects)", Value: 2},
				{Label: "Advanced (internship/work experience)", Value: 3},
				{Label: "Expert (open source/competitive programming)", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 0, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 2, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 2, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 2, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 5, "MS Abroad": 5}},
			},
			DisplayOrder: 6,
		},
		{
			Category:     "academic",
			QuestionText: "How confident are you with English?",
			Options: []dto.QuestionOption{
				{Label: "Basic (can read/write simple sentences)", Value: 0},
				{Label: "Intermediate (comfortable in conversation)", Value: 1},
				{Label: "Advanced (can write essays, present confidently)", Value: 2},
				{Label: "Fluent (near-native, IELTS 7+ equivalent)", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 1, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 0}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 5}},
			},
			DisplayOrder: 7,
		},
		{
			Category:     "academic",
			QuestionText: "Have you done any internships or jobs so far?",
			Options: []dto.QuestionOption{
				{Label: "No experience", Value: 0},
				{Label: "1-2 internships", Value: 1},
				{Label: "1-2 years full-time work", Value: 2},
				{Label: "3+ years work experience", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 1, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 2, "MS Abroad": 3}},
			},
			DisplayOrder: 8,
		},

		// ============================================================
		// SECTION B: FINANCIAL SITUATION (Questions 9–14)
		// ============================================================
		{
			Category:     "financial",
			QuestionText: "Do you have any education loan to repay?",
			Options: []dto.QuestionOption{
				{Label: "No loan", Value: 0},
				{Label: "Small loan (under ₹5 lakh)", Value: 1},
				{Label: "Medium loan (₹5-15 lakh)", Value: 2},
				{Label: "Large loan (₹15 lakh+)", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 4, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}, RiskFactors: map[string]float64{"income_urgency": 2, "family_dependency": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 2}, RiskFactors: map[string]float64{"income_urgency": 4, "family_dependency": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 2, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 2}, RiskFactors: map[string]float64{"income_urgency": 6, "family_dependency": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 1, "Government Exams": 2, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 1, "MS Abroad": 1}, RiskFactors: map[string]float64{"income_urgency": 9, "family_dependency": 7}},
			},
			DisplayOrder: 9,
		},
		{
			Category:     "financial",
			QuestionText: "How many people in your family will depend on your income?",
			Options: []dto.QuestionOption{
				{Label: "None (financially independent family)", Value: 0},
				{Label: "1-2 members", Value: 1},
				{Label: "3-4 members", Value: 2},
				{Label: "5+ members (sole breadwinner expected)", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 5, "MS Abroad": 5}, RiskFactors: map[string]float64{"family_dependency": 1, "income_urgency": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"family_dependency": 4, "income_urgency": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 1}, RiskFactors: map[string]float64{"family_dependency": 7, "income_urgency": 6}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 1, "Government Exams": 5, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 1, "MS Abroad": 0}, RiskFactors: map[string]float64{"family_dependency": 9, "income_urgency": 8}},
			},
			DisplayOrder: 10,
		},
		{
			Category:     "financial",
			QuestionText: "How soon do you need to start earning?",
			Options: []dto.QuestionOption{
				{Label: "No urgency, can study/prepare for 2-3 years", Value: 0},
				{Label: "Preferably within 1 year but flexible", Value: 1},
				{Label: "Yes, must earn within 1 year", Value: 2},
				{Label: "Yes, urgently need income within 6 months", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 4, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 5, "MS Abroad": 5}, RiskFactors: map[string]float64{"income_urgency": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"income_urgency": 4}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 1, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 1, "MS Abroad": 1}, RiskFactors: map[string]float64{"income_urgency": 7}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 0, "Government Exams": 2, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 0, "MS Abroad": 0}, RiskFactors: map[string]float64{"income_urgency": 10}},
			},
			DisplayOrder: 11,
		},
		{
			Category:     "financial",
			QuestionText: "Could your family manage MBA-level fees (₹20 lakh+)?",
			Options: []dto.QuestionOption{
				{Label: "Yes, comfortably", Value: 0},
				{Label: "Yes, with some financial strain", Value: 1},
				{Label: "Only with full education loan", Value: 2},
				{Label: "No, cannot afford", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}, RiskFactors: map[string]float64{"income_urgency": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"income_urgency": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 2}, RiskFactors: map[string]float64{"income_urgency": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 1, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 1}, RiskFactors: map[string]float64{"income_urgency": 7}},
			},
			DisplayOrder: 12,
		},
		{
			Category:     "financial",
			QuestionText: "What's your family's approximate annual income?",
			Options: []dto.QuestionOption{
				{Label: "Below ₹3 lakh", Value: 0},
				{Label: "₹3-8 lakh", Value: 1},
				{Label: "₹8-15 lakh", Value: 2},
				{Label: "₹15-30 lakh", Value: 3},
				{Label: "₹30 lakh+", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 1, "Government Exams": 5, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 0}, RiskFactors: map[string]float64{"income_urgency": 8, "family_dependency": 7}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 1}, RiskFactors: map[string]float64{"income_urgency": 6, "family_dependency": 5}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 3}, RiskFactors: map[string]float64{"income_urgency": 3, "family_dependency": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}, RiskFactors: map[string]float64{"income_urgency": 2, "family_dependency": 1}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 5}, RiskFactors: map[string]float64{"income_urgency": 1, "family_dependency": 0}},
			},
			DisplayOrder: 13,
		},
		{
			Category:     "financial",
			QuestionText: "Where did you grow up?",
			Options: []dto.QuestionOption{
				{Label: "Rural / small town", Value: 0},
				{Label: "Tier-2/3 city", Value: 1},
				{Label: "Tier-1 city (Delhi, Mumbai, Bangalore, etc.)", Value: 2},
				{Label: "Metro with strong professional network", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 5}},
			},
			DisplayOrder: 14,
		},

		// ============================================================
		// SECTION C: PERSONALITY & RISK (Questions 15–22)
		// ============================================================
		{
			Category:     "personality",
			QuestionText: "How do you feel about taking risks in your career?",
			Options: []dto.QuestionOption{
				{Label: "Very risk-averse (prefer guaranteed outcomes)", Value: 0},
				{Label: "Low risk (small calculated risks only)", Value: 1},
				{Label: "Moderate (balanced approach)", Value: 2},
				{Label: "High risk (comfortable with uncertainty)", Value: 3},
				{Label: "Very high risk (thrive in chaos)", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 4, "MS Abroad": 2}, RiskFactors: map[string]float64{"risk_tolerance": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 4, "MS Abroad": 3}, RiskFactors: map[string]float64{"risk_tolerance": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 4}, RiskFactors: map[string]float64{"risk_tolerance": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 2, "MS Abroad": 4}, RiskFactors: map[string]float64{"risk_tolerance": 7}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 1, "MS Abroad": 3}, RiskFactors: map[string]float64{"risk_tolerance": 9}},
			},
			DisplayOrder: 15,
		},
		{
			Category:     "personality",
			QuestionText: "Do you see yourself leading and managing teams?",
			Options: []dto.QuestionOption{
				{Label: "No, I prefer individual work", Value: 0},
				{Label: "Somewhat, in small teams", Value: 1},
				{Label: "Yes, I enjoy leading teams", Value: 2},
				{Label: "Absolutely, I want to build organizations", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 1, "Government Exams": 3, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 5, "MS Abroad": 4}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 4, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 2, "MS Abroad": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 5, "Government Exams": 4, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 2, "MS Abroad": 2}},
			},
			DisplayOrder: 16,
		},
		{
			Category:     "personality",
			QuestionText: "What matters more to you — job security or fast growth?",
			Options: []dto.QuestionOption{
				{Label: "Strong stability (pension, job security)", Value: 0},
				{Label: "Stability with moderate growth", Value: 1},
				{Label: "Growth with some stability", Value: 2},
				{Label: "Maximum growth (even if risky)", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 3, "MS Abroad": 2}, RiskFactors: map[string]float64{"career_instability": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"career_instability": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 4}, RiskFactors: map[string]float64{"career_instability": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 0, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 2, "MS Abroad": 4}, RiskFactors: map[string]float64{"career_instability": 8}},
			},
			DisplayOrder: 17,
		},
		{
			Category:     "personality",
			QuestionText: "How do you deal with pressure and tight deadlines?",
			Options: []dto.QuestionOption{
				{Label: "I avoid stressful situations", Value: 0},
				{Label: "I manage but prefer low-stress", Value: 1},
				{Label: "I handle stress well", Value: 2},
				{Label: "I perform best under pressure", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 1, "Government Exams": 3, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 3, "MS Abroad": 4}},
			},
			DisplayOrder: 18,
		},
		{
			Category:     "personality",
			QuestionText: "How important is having free time outside of work?",
			Options: []dto.QuestionOption{
				{Label: "Extremely important (9-to-5 preferred)", Value: 0},
				{Label: "Important but flexible", Value: 1},
				{Label: "Willing to sacrifice for career growth", Value: 2},
				{Label: "Will work 80+ hours if passionate", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 3}},
			},
			DisplayOrder: 19,
		},
		{
			Category:     "personality",
			QuestionText: "Do you prefer solving new problems or following proven methods?",
			Options: []dto.QuestionOption{
				{Label: "Prefer following established processes", Value: 0},
				{Label: "Mix of both", Value: 1},
				{Label: "Enjoy solving new problems", Value: 2},
				{Label: "Love creating new solutions from scratch", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 2, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 5, "MS Abroad": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 3, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 5, "MS Abroad": 5}},
			},
			DisplayOrder: 20,
		},
		{
			Category:     "personality",
			QuestionText: "How do you feel about constantly learning new things?",
			Options: []dto.QuestionOption{
				{Label: "Prefer to learn once and apply", Value: 0},
				{Label: "Open to occasional learning", Value: 1},
				{Label: "Enjoy regular learning and courses", Value: 2},
				{Label: "Passionate about constant learning", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 1, "MS Abroad": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 4, "Government Exams": 2, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 5, "MS Abroad": 5}},
			},
			DisplayOrder: 21,
		},
		{
			Category:     "personality",
			QuestionText: "How do you feel about speaking in front of people?",
			Options: []dto.QuestionOption{
				{Label: "Very uncomfortable", Value: 0},
				{Label: "Nervous but can manage", Value: 1},
				{Label: "Comfortable and experienced", Value: 2},
				{Label: "Excellent — enjoy presenting", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 1, "Government Exams": 2, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 4, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 5, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 4, "MS Abroad": 4}},
			},
			DisplayOrder: 22,
		},

		// ============================================================
		// SECTION D: CAREER INTEREST (Questions 23–30)
		// ============================================================
		{
			Category:     "career_interest",
			QuestionText: "How much does a government job (IAS/IPS/Banking) appeal to you?",
			Options: []dto.QuestionOption{
				{Label: "Not interested at all", Value: 0},
				{Label: "Slightly interested", Value: 1},
				{Label: "Moderately interested", Value: 2},
				{Label: "Very interested", Value: 3},
				{Label: "It's my dream", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 4}, RiskFactors: map[string]float64{"career_instability": 3}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"career_instability": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 2}, RiskFactors: map[string]float64{"career_instability": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 1, "Government Exams": 5, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 1}, RiskFactors: map[string]float64{"career_instability": 5}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 0, "MBA (India)": 0, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 1, "MS Abroad": 0}, RiskFactors: map[string]float64{"career_instability": 6}},
			},
			DisplayOrder: 23,
		},
		{
			Category:     "career_interest",
			QuestionText: "Have you thought about starting your own business?",
			Options: []dto.QuestionOption{
				{Label: "Not interested", Value: 0},
				{Label: "Maybe someday", Value: 1},
				{Label: "Actively thinking about it", Value: 2},
				{Label: "Already working on an idea", Value: 3},
				{Label: "Already have a side project/startup", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 4, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"career_instability": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 3}, RiskFactors: map[string]float64{"career_instability": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 1, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 2, "MS Abroad": 2}, RiskFactors: map[string]float64{"career_instability": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 1, "MS Abroad": 1}, RiskFactors: map[string]float64{"career_instability": 7}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 1, "MS Abroad": 1}, RiskFactors: map[string]float64{"career_instability": 8}},
			},
			DisplayOrder: 24,
		},
		{
			Category:     "career_interest",
			QuestionText: "How excited are you about working at a big company or MNC?",
			Options: []dto.QuestionOption{
				{Label: "Not interested", Value: 0},
				{Label: "As a backup option", Value: 1},
				{Label: "Good option for me", Value: 2},
				{Label: "My primary goal", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 1, "Government Exams": 4, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 4}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 5, "Government Exams": 1, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 2, "MS Abroad": 4}},
			},
			DisplayOrder: 25,
		},
		{
			Category:     "career_interest",
			QuestionText: "Would you like to pursue Masters or PhD in India?",
			Options: []dto.QuestionOption{
				{Label: "No", Value: 0},
				{Label: "Maybe after some work experience", Value: 1},
				{Label: "Yes, planning to apply", Value: 2},
				{Label: "Absolutely, research is my passion", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 4, "Higher Studies (India)": 0, "MS Abroad": 2}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 4, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 2, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 5, "MS Abroad": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 1, "Government Exams": 1, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 5, "MS Abroad": 4}},
			},
			DisplayOrder: 26,
		},
		{
			Category:     "career_interest",
			QuestionText: "Do you dream of studying or working outside India?",
			Options: []dto.QuestionOption{
				{Label: "No, I want to stay in India", Value: 0},
				{Label: "Open to it but not actively pursuing", Value: 1},
				{Label: "Yes, actively preparing (GRE/TOEFL)", Value: 2},
				{Label: "Definitely going abroad", Value: 3},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 4, "Government Exams": 5, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 0}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 1, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 5}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 1, "Government Exams": 0, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 1, "MS Abroad": 5}},
			},
			DisplayOrder: 27,
		},
		{
			Category:     "career_interest",
			QuestionText: "What kind of starting salary would make you happy?",
			Options: []dto.QuestionOption{
				{Label: "₹2-4 LPA (just need a job)", Value: 0},
				{Label: "₹4-8 LPA (decent start)", Value: 1},
				{Label: "₹8-15 LPA (competitive package)", Value: 2},
				{Label: "₹15-25 LPA (premium placement)", Value: 3},
				{Label: "₹25 LPA+ or $60K+ (top-tier)", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 1, "Government Exams": 4, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 1}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 2, "Government Exams": 4, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 2}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 4, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 3, "MS Abroad": 3}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 4, "MBA (India)": 5, "Government Exams": 1, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 0, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 5}},
			},
			DisplayOrder: 28,
		},
		{
			Category:     "career_interest",
			QuestionText: "Which field gets you the most excited?",
			Options: []dto.QuestionOption{
				{Label: "Technology / Software", Value: 0},
				{Label: "Finance / Banking / Consulting", Value: 1},
				{Label: "Government / Public Service", Value: 2},
				{Label: "Healthcare / Pharma", Value: 3},
				{Label: "Education / Research", Value: 4},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 3, "Government Exams": 1, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 3, "MS Abroad": 4}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 3, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 2, "Government Exams": 5, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 3, "MS Abroad": 1}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 4, "MS Abroad": 4}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 2, "Government Exams": 3, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 5, "MS Abroad": 4}},
			},
			DisplayOrder: 29,
		},
		{
			Category:     "career_interest",
			QuestionText: "Imagine yourself 10 years from now — what do you see?",
			Options: []dto.QuestionOption{
				{Label: "Senior engineer at a top tech company", Value: 0},
				{Label: "Business leader / VP at a corporation", Value: 1},
				{Label: "Government officer (IAS/IPS/IFS)", Value: 2},
				{Label: "Running my own successful company", Value: 3},
				{Label: "Professor / Researcher at a top institution", Value: 4},
				{Label: "Living abroad with a high-paying job", Value: 5},
			},
			Weights: []dto.QuestionWeight{
				{OptionIndex: 0, Scores: map[string]float64{"IT / Software Jobs": 5, "MBA (India)": 2, "Government Exams": 0, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 3, "MS Abroad": 4}},
				{OptionIndex: 1, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 5, "Government Exams": 2, "Startup / Entrepreneurship": 3, "Higher Studies (India)": 2, "MS Abroad": 3}},
				{OptionIndex: 2, Scores: map[string]float64{"IT / Software Jobs": 0, "MBA (India)": 1, "Government Exams": 5, "Startup / Entrepreneurship": 0, "Higher Studies (India)": 1, "MS Abroad": 0}},
				{OptionIndex: 3, Scores: map[string]float64{"IT / Software Jobs": 2, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 5, "Higher Studies (India)": 1, "MS Abroad": 2}},
				{OptionIndex: 4, Scores: map[string]float64{"IT / Software Jobs": 1, "MBA (India)": 1, "Government Exams": 1, "Startup / Entrepreneurship": 1, "Higher Studies (India)": 5, "MS Abroad": 3}},
				{OptionIndex: 5, Scores: map[string]float64{"IT / Software Jobs": 3, "MBA (India)": 3, "Government Exams": 0, "Startup / Entrepreneurship": 2, "Higher Studies (India)": 2, "MS Abroad": 5}},
			},
			DisplayOrder: 30,
		},
	}
}
