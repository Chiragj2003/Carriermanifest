// Package dto defines Data Transfer Objects for API request/response.
package dto

// ============================================================
// AUTH DTOs
// ============================================================

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// GoogleLoginRequest is the payload for Google OAuth sign-in.
type GoogleLoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

// AuthResponse is returned after successful login/register.
type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// UserDTO is the safe user representation (no password).
type UserDTO struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// ============================================================
// ASSESSMENT DTOs
// ============================================================

// AnswerItem represents a single question answer.
type AnswerItem struct {
	QuestionID uint64 `json:"question_id" binding:"required"`
	Selected   int    `json:"selected" binding:"min=0"` // index of selected option
}

// SubmitAssessmentRequest is the payload for submitting an assessment.
type SubmitAssessmentRequest struct {
	Answers []AnswerItem `json:"answers" binding:"required,min=1"`
}

// CareerScore holds a score for a specific career category.
type CareerScore struct {
	Category   string  `json:"category"`
	Score      float64 `json:"score"`
	MaxScore   float64 `json:"max_score"`
	Percentage float64 `json:"percentage"`
}

// RiskAssessment holds the risk analysis result.
type RiskAssessment struct {
	Score   float64            `json:"score"`
	Level   string             `json:"level"` // Low, Medium, High
	Factors map[string]float64 `json:"factors"`
}

// SalaryProjection holds 5-year salary growth data.
type SalaryProjection struct {
	Year1 string `json:"year_1"`
	Year2 string `json:"year_2"`
	Year3 string `json:"year_3"`
	Year4 string `json:"year_4"`
	Year5 string `json:"year_5"`
}

// RoadmapStep is a single step in the preparation roadmap.
type RoadmapStep struct {
	Step        int    `json:"step"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
}

// FeatureContributionDTO describes how a feature contributed to a career score.
type FeatureContributionDTO struct {
	Feature      string  `json:"feature"`
	UserValue    float64 `json:"user_value"`
	CareerWeight float64 `json:"career_weight"`
	Contribution float64 `json:"contribution"`
	Percentage   float64 `json:"percentage"`
}

// CareerExplanationDTO holds the deterministic explanation for a career.
type CareerExplanationDTO struct {
	Career     string                   `json:"career"`
	TopFactors []FeatureContributionDTO `json:"top_factors"`
	Summary    string                   `json:"summary"`
	Penalties  []string                 `json:"penalties,omitempty"`
}

// UserProfileDTO represents the aggregated feature profile of the user.
type UserProfileDTO struct {
	AcademicStrength  float64 `json:"academic_strength"`
	FinancialPressure float64 `json:"financial_pressure"`
	RiskTolerance     float64 `json:"risk_tolerance"`
	LeadershipScore   float64 `json:"leadership_score"`
	TechAffinity      float64 `json:"tech_affinity"`
	GovtInterest      float64 `json:"govt_interest"`
	AbroadInterest    float64 `json:"abroad_interest"`
	IncomeUrgency     float64 `json:"income_urgency"`
	CareerInstability float64 `json:"career_instability"`
}

// VersionInfo holds version metadata for reproducibility and ML provenance.
type VersionInfo struct {
	Assessment    string  `json:"assessment"`
	WeightMatrix  string  `json:"weight_matrix"`
	FeatureMap    string  `json:"feature_map"`
	ModelType     string  `json:"model_type"`
	ModelAccuracy float64 `json:"model_accuracy"`
	ModelF1Score  float64 `json:"model_f1_score"`
}

// AssessmentResult is the full computed result.
type AssessmentResult struct {
	Scores            []CareerScore          `json:"scores"`
	BestCareerPath    string                 `json:"best_career_path"`
	Confidence        float64                `json:"confidence"`
	IsMultiFit        bool                   `json:"is_multi_fit"`
	Risk              RiskAssessment         `json:"risk"`
	Profile           UserProfileDTO         `json:"profile"`
	Explanations      []CareerExplanationDTO `json:"explanations"`
	SalaryProjection  SalaryProjection       `json:"salary_projection"`
	Roadmap           []RoadmapStep          `json:"roadmap"`
	RequiredSkills    []string               `json:"required_skills"`
	SuggestedExams    []string               `json:"suggested_exams"`
	SuggestedColleges []string               `json:"suggested_colleges"`
	Version           VersionInfo            `json:"version"`
	AIExplanation     string                 `json:"ai_explanation,omitempty"`
}

// AssessmentResponse is returned after submitting an assessment.
type AssessmentResponse struct {
	ID        uint64           `json:"id"`
	UserID    uint64           `json:"user_id"`
	Result    AssessmentResult `json:"result"`
	CreatedAt string           `json:"created_at"`
}

// AssessmentListItem is a summary for the dashboard list.
type AssessmentListItem struct {
	ID             uint64 `json:"id"`
	BestCareerPath string `json:"best_career_path"`
	RiskLevel      string `json:"risk_level"`
	CreatedAt      string `json:"created_at"`
}

// ============================================================
// QUESTION DTOs
// ============================================================

// QuestionOption represents a single selectable option.
type QuestionOption struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

// QuestionWeight defines how an option scores for each career.
type QuestionWeight struct {
	OptionIndex int                `json:"option_index"`
	Scores      map[string]float64 `json:"scores"` // career_category -> score
	RiskFactors map[string]float64 `json:"risk_factors,omitempty"`
}

// QuestionDTO is the API representation of a question.
type QuestionDTO struct {
	ID           uint64           `json:"id"`
	Category     string           `json:"category"`
	QuestionText string           `json:"question_text"`
	Options      []QuestionOption `json:"options"`
	Weights      []QuestionWeight `json:"weights,omitempty"` // only for admin
	DisplayOrder int              `json:"display_order"`
	IsActive     *bool            `json:"is_active,omitempty"` // only for admin
}

// CreateQuestionRequest is used by admins to create questions.
type CreateQuestionRequest struct {
	Category     string           `json:"category" binding:"required"`
	QuestionText string           `json:"question_text" binding:"required"`
	Options      []QuestionOption `json:"options" binding:"required,min=2"`
	Weights      []QuestionWeight `json:"weights" binding:"required"`
	DisplayOrder int              `json:"display_order"`
}

// UpdateQuestionRequest is used by admins to update questions.
type UpdateQuestionRequest struct {
	Category     string           `json:"category"`
	QuestionText string           `json:"question_text"`
	Options      []QuestionOption `json:"options"`
	Weights      []QuestionWeight `json:"weights"`
	DisplayOrder int              `json:"display_order"`
	IsActive     *bool            `json:"is_active"`
}

// ============================================================
// ADMIN DTOs
// ============================================================

// AdminStatsResponse shows assessment statistics.
type AdminStatsResponse struct {
	TotalUsers         int            `json:"total_users"`
	TotalAssessments   int            `json:"total_assessments"`
	TotalQuestions     int            `json:"total_questions"`
	CareerDistribution map[string]int `json:"career_distribution"`
	RiskDistribution   map[string]int `json:"risk_distribution"`
}

// ============================================================
// CHAT DTOs
// ============================================================

// ChatRequest is the payload for the AI chatbot.
type ChatRequest struct {
	Message      string `json:"message" binding:"required,min=1,max=1000"`
	AssessmentID uint64 `json:"assessment_id"`
}

// ChatResponse is returned from the AI chatbot.
type ChatResponse struct {
	Reply string `json:"reply"`
}

// ErrorResponse is a standard error payload.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse is a standard success payload.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
