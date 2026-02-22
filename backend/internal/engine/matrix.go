package engine

// CareerWeightMatrix stores a weight vector per career over UserProfile features.
// Each row is a career (indexed by Career enum), each column is a feature.
// Positive weights mean the feature favours that career; negative weights penalize it.
//
// Matrix dimensions: NumCareers x NumFeatures
//
// Feature order:
//   [0] AcademicStrength  [1] FinancialPressure  [2] RiskTolerance
//   [3] LeadershipScore   [4] TechAffinity       [5] GovtInterest
//   [6] AbroadInterest    [7] IncomeUrgency      [8] CareerInstability
var CareerWeightMatrix = [NumCareers][NumFeatures]float64{
	// CareerIT: Favours tech affinity, academic strength; penalised by financial pressure
	{
		0.70,  // AcademicStrength
		-0.25, // FinancialPressure (IT pays well early — slight negative)
		0.30,  // RiskTolerance
		0.20,  // LeadershipScore
		0.90,  // TechAffinity (strongest signal)
		-0.20, // GovtInterest (inverse)
		0.35,  // AbroadInterest (IT jobs are global)
		-0.15, // IncomeUrgency (IT has fast income)
		0.10,  // CareerInstability
	},
	// CareerMBA: Favours leadership, moderate academic; penalised by high income urgency
	{
		0.50,  // AcademicStrength
		-0.30, // FinancialPressure (MBA is expensive)
		0.40,  // RiskTolerance
		0.85,  // LeadershipScore (strongest signal)
		0.15,  // TechAffinity (not core)
		-0.10, // GovtInterest
		0.30,  // AbroadInterest (MBA global too)
		-0.35, // IncomeUrgency (2yr ROI delay)
		0.15,  // CareerInstability
	},
	// CareerGovt: Favours govt interest, low risk, job security; penalised by instability
	{
		0.40,  // AcademicStrength
		0.25,  // FinancialPressure (govt = stable income)
		-0.50, // RiskTolerance (govt seekers are risk-averse)
		0.30,  // LeadershipScore
		-0.20, // TechAffinity (not tech-heavy)
		0.95,  // GovtInterest (strongest signal)
		-0.40, // AbroadInterest (stay in India)
		0.10,  // IncomeUrgency (moderate — exams take time)
		-0.45, // CareerInstability (want stability)
	},
	// CareerStartup: Favours risk tolerance, leadership; penalised by financial pressure
	{
		0.30,  // AcademicStrength
		-0.40, // FinancialPressure (startups = no guaranteed income)
		0.85,  // RiskTolerance (strongest signal)
		0.80,  // LeadershipScore
		0.55,  // TechAffinity (useful but not mandatory)
		-0.30, // GovtInterest (opposite)
		0.20,  // AbroadInterest
		-0.50, // IncomeUrgency (startups = delayed income)
		0.40,  // CareerInstability (embraces chaos)
	},
	// CareerHigherStudies: Favours academic strength, research mindset
	{
		0.90,  // AcademicStrength (strongest signal)
		-0.30, // FinancialPressure (stipend only)
		0.20,  // RiskTolerance
		0.15,  // LeadershipScore
		0.40,  // TechAffinity (research can be tech)
		0.10,  // GovtInterest (slightly — academic institutions)
		0.25,  // AbroadInterest
		-0.45, // IncomeUrgency (years of study)
		-0.10, // CareerInstability
	},
	// CareerMSAbroad: Favours abroad interest, academic strength, English ability
	{
		0.75,  // AcademicStrength
		-0.35, // FinancialPressure (expensive initially)
		0.35,  // RiskTolerance
		0.20,  // LeadershipScore
		0.45,  // TechAffinity (STEM-heavy)
		-0.30, // GovtInterest (leaving India)
		0.90,  // AbroadInterest (strongest signal)
		-0.40, // IncomeUrgency (2yr delay)
		0.10,  // CareerInstability
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
