package engine

// NumFeatures is the total number of profile features.
const NumFeatures = 9

// Feature indices in the UserProfile vector.
const (
	FeatAcademicStrength  = 0
	FeatFinancialPressure = 1
	FeatRiskTolerance     = 2
	FeatLeadershipScore   = 3
	FeatTechAffinity      = 4
	FeatGovtInterest      = 5
	FeatAbroadInterest    = 6
	FeatIncomeUrgency     = 7
	FeatCareerInstability = 8
)

// FeatureNames maps feature indices to human-readable names.
var FeatureNames = [NumFeatures]string{
	"Academic Strength",
	"Financial Pressure",
	"Risk Tolerance",
	"Leadership",
	"Tech Affinity",
	"Govt Interest",
	"Abroad Interest",
	"Income Urgency",
	"Career Instability",
}

// UserProfile holds aggregated feature scores computed from assessment answers.
// Each field is normalized to a 0–1 scale after aggregation.
type UserProfile struct {
	Features [NumFeatures]float64
}

// AcademicStrength returns the academic strength score.
func (p *UserProfile) AcademicStrength() float64 { return p.Features[FeatAcademicStrength] }

// FinancialPressure returns the financial pressure score.
func (p *UserProfile) FinancialPressure() float64 { return p.Features[FeatFinancialPressure] }

// RiskTolerance returns the risk tolerance score.
func (p *UserProfile) RiskTolerance() float64 { return p.Features[FeatRiskTolerance] }

// LeadershipScore returns the leadership score.
func (p *UserProfile) LeadershipScore() float64 { return p.Features[FeatLeadershipScore] }

// TechAffinity returns the tech affinity score.
func (p *UserProfile) TechAffinity() float64 { return p.Features[FeatTechAffinity] }

// GovtInterest returns the government interest score.
func (p *UserProfile) GovtInterest() float64 { return p.Features[FeatGovtInterest] }

// AbroadInterest returns the abroad interest score.
func (p *UserProfile) AbroadInterest() float64 { return p.Features[FeatAbroadInterest] }

// IncomeUrgency returns the income urgency score.
func (p *UserProfile) IncomeUrgency() float64 { return p.Features[FeatIncomeUrgency] }

// CareerInstability returns the career instability score.
func (p *UserProfile) CareerInstability() float64 { return p.Features[FeatCareerInstability] }

// Vector returns the feature vector as a slice for dot product computation.
func (p *UserProfile) Vector() []float64 {
	v := make([]float64, NumFeatures)
	copy(v, p.Features[:])
	return v
}
