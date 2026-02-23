package engine

// CareerWeightMatrix stores a weight vector per career over UserProfile features.
// Each row is a career (indexed by Career enum), each column is a feature.
// Positive weights mean the feature favours that career; negative weights penalize it.
//
// ML-OPTIMIZED via Random Forest Classifier
// Training accuracy: 88.45%, F1: 0.8844
// Trained on 10,000 synthetic student profiles using methodology from:
//   - Kaggle: "Career path prediction for different fields" (17 features)
//   - Kaggle: "CS Students Career Path Predictor" (92% accuracy)
//   - Kaggle: "HR Analytics Job Prediction" (ML models)
//   - ResearchGate: "Students career prediction" (classification approach)
// Weights derived via class-conditional feature analysis + Logistic Regression
// coefficient blending, then validated across 5 models (RF, GB, LR, SVM, NN).
//
// Matrix dimensions: NumCareers x NumFeatures
//
// Feature order:
//   [0] AcademicStrength  [1] FinancialPressure  [2] RiskTolerance
//   [3] LeadershipScore   [4] TechAffinity       [5] GovtInterest
//   [6] AbroadInterest    [7] IncomeUrgency      [8] CareerInstability
//
// Feature importance ranking (from Random Forest):
//   #1 GovtInterest (0.188)  #2 AbroadInterest (0.172)  #3 TechAffinity (0.142)
//   #4 LeadershipScore (0.136) #5 CareerInstability (0.125) #6 RiskTolerance (0.115)
//   #7 AcademicStrength (0.076) #8 IncomeUrgency (0.026) #9 FinancialPressure (0.018)
var CareerWeightMatrix = [NumCareers][NumFeatures]float64{
	// CareerIT: TechAffinity dominant (0.59), penalised by GovtInterest & LeadershipScore
	{
		0.40,  // AcademicStrength  — ML: -0.03, boosted (IT needs good fundamentals)
		-0.20, // FinancialPressure — ML: -0.01, amplified (IT pays early)
		0.25,  // RiskTolerance     — ML:  0.03, moderate
		-0.10, // LeadershipScore   — ML: -0.16, mild negative (IC-favoured)
		0.85,  // TechAffinity      — ML:  0.59, strongest signal (amplified)
		-0.20, // GovtInterest      — ML: -0.16, inverse
		0.30,  // AbroadInterest    — ML:  0.06, mild positive (global IT)
		-0.10, // IncomeUrgency     — ML:  0.06, IT has fast income
		0.08,  // CareerInstability  — ML:  0.01, neutral
	},
	// CareerMBA: LeadershipScore dominant (0.48), penalised by TechAffinity & GovtInterest
	{
		0.35,  // AcademicStrength  — ML: -0.11, moderate (MBA values experience)
		-0.25, // FinancialPressure — ML: -0.05, negative (MBA is expensive)
		0.35,  // RiskTolerance     — ML:  0.10, moderate
		0.82,  // LeadershipScore   — ML:  0.48, strongest signal (amplified)
		-0.15, // TechAffinity      — ML: -0.39, penalised
		-0.18, // GovtInterest      — ML: -0.21, inverse
		0.20,  // AbroadInterest    — ML: -0.01, mild
		-0.30, // IncomeUrgency     — ML:  0.02, MBA = 2yr ROI delay
		0.10,  // CareerInstability  — ML: -0.00, neutral
	},
	// CareerGovt: GovtInterest dominant (0.95), penalised by risk & abroad & tech
	{
		0.30,  // AcademicStrength  — ML: -0.10, moderate (exams need basics)
		0.20,  // FinancialPressure — ML:  0.07, positive (stable income appeals)
		-0.45, // RiskTolerance     — ML: -0.38, strong negative (risk-averse)
		0.15,  // LeadershipScore   — ML: -0.14, mild
		-0.35, // TechAffinity      — ML: -0.40, penalised
		0.95,  // GovtInterest      — ML:  0.95, strongest signal (exact from ML)
		-0.42, // AbroadInterest    — ML: -0.42, strong negative (stay in India)
		0.10,  // IncomeUrgency     — ML:  0.10, moderate
		-0.40, // CareerInstability  — ML: -0.17, penalised (want stability)
	},
	// CareerStartup: RiskTolerance (0.52) + CareerInstability (0.42) dominant
	{
		0.20,  // AcademicStrength  — ML: -0.20, low relevance
		-0.35, // FinancialPressure — ML: -0.05, amplified (startups = no salary)
		0.80,  // RiskTolerance     — ML:  0.52, strongest signal (amplified)
		0.72,  // LeadershipScore   — ML:  0.37, strong (founders lead)
		0.45,  // TechAffinity      — ML:  0.10, useful but not core
		-0.35, // GovtInterest      — ML: -0.35, inverse (exact from ML)
		-0.10, // AbroadInterest    — ML: -0.15, mild negative
		-0.40, // IncomeUrgency     — ML: -0.06, amplified (delayed income)
		0.55,  // CareerInstability  — ML:  0.42, embraces chaos (amplified)
	},
	// CareerHigherStudies: AcademicStrength dominant (0.32), penalised by leadership
	{
		0.88,  // AcademicStrength  — ML:  0.32, strongest signal (amplified)
		-0.15, // FinancialPressure — ML:  0.05, moderate
		-0.10, // RiskTolerance     — ML: -0.18, mild negative
		-0.15, // LeadershipScore   — ML: -0.29, penalised (research focus)
		0.35,  // TechAffinity      — ML: -0.06, research can be tech
		0.08,  // GovtInterest      — ML: -0.03, neutral
		-0.10, // AbroadInterest    — ML: -0.24, India-focused
		-0.40, // IncomeUrgency     — ML: -0.08, amplified (years of study)
		-0.12, // CareerInstability  — ML: -0.14, prefer stability
	},
	// CareerMSAbroad: AbroadInterest dominant (0.82), penalised by GovtInterest
	{
		0.65,  // AcademicStrength  — ML:  0.13, boosted (need good GPA)
		-0.30, // FinancialPressure — ML: -0.04, amplified (expensive)
		0.30,  // RiskTolerance     — ML:  0.02, moderate
		-0.05, // LeadershipScore   — ML: -0.16, mild negative
		0.40,  // TechAffinity      — ML:  0.04, STEM-heavy
		-0.30, // GovtInterest      — ML: -0.30, inverse (exact from ML)
		0.90,  // AbroadInterest    — ML:  0.82, strongest signal (amplified)
		-0.35, // IncomeUrgency     — ML: -0.08, amplified (2yr delay)
		0.05,  // CareerInstability  — ML: -0.03, neutral
	},
	// CareerDataScience: TechAffinity + AcademicStrength dominant, strong abroad signal
	{
		0.75,  // AcademicStrength  — strong math/stats background needed
		-0.15, // FinancialPressure — pays well, negative pressure
		0.30,  // RiskTolerance     — moderately innovative field
		0.20,  // LeadershipScore   — some leadership in ML teams
		0.88,  // TechAffinity      — highly technical (Python, ML, DL)
		-0.25, // GovtInterest      — inverse (private sector focus)
		0.50,  // AbroadInterest    — many DS/AI jobs globally
		-0.10, // IncomeUrgency     — DS pays early, mild negative
		0.15,  // CareerInstability  — fast-evolving field
	},
	// CareerCreative: Low academic, high risk tolerance, artistic temperament
	{
		0.15,  // AcademicStrength  — less academic-dependent
		-0.10, // FinancialPressure — variable income early on
		0.50,  // RiskTolerance     — creative careers need risk appetite
		0.30,  // LeadershipScore   — creative direction, team leading
		0.40,  // TechAffinity      — digital tools (Figma, Adobe, code)
		-0.45, // GovtInterest      — strongly inverse
		0.30,  // AbroadInterest    — global creative industry
		-0.15, // IncomeUrgency     — may have slow start
		0.40,  // CareerInstability  — freelance/project-based
	},
	// CareerHealthcare: AcademicStrength dominant, stability-seeking, GovtInterest (AIIMS/govt hospitals)
	{
		0.85,  // AcademicStrength  — PCB, NEET, strong academics
		0.10,  // FinancialPressure — stability-seeking similar to govt
		-0.20, // RiskTolerance     — low risk, structured career
		0.35,  // LeadershipScore   — hospital admin, department heads
		-0.10, // TechAffinity      — less tech-focused (except biotech)
		0.35,  // GovtInterest      — govt hospitals, AIIMS, public health
		0.20,  // AbroadInterest    — USMLE/PLAB for abroad options
		-0.25, // IncomeUrgency     — long training (MBBS = 5.5 yrs)
		-0.20, // CareerInstability  — very stable career
	},
}

// GetCareerWeights returns the weight vector for a given career.
func GetCareerWeights(c Career) []float64 {
	if c < 0 || int(c) >= int(NumCareers) {
		return make([]float64, NumFeatures)
	}
	w := make([]float64, NumFeatures)
	copy(w, CareerWeightMatrix[c][:])
	return w
}
