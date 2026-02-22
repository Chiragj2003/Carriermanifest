package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/engine"
	"github.com/careermanifest/backend/internal/repository"
)

// AssessmentService handles assessment business logic.
type AssessmentService struct {
	assessmentRepo *repository.AssessmentRepository
	questionRepo   *repository.QuestionRepository
	scoringEngine  *engine.ScoringEngine
	llmService     *LLMService
}

// NewAssessmentService creates a new AssessmentService.
func NewAssessmentService(
	assessmentRepo *repository.AssessmentRepository,
	questionRepo *repository.QuestionRepository,
	scoringEngine *engine.ScoringEngine,
	llmService *LLMService,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: assessmentRepo,
		questionRepo:   questionRepo,
		scoringEngine:  scoringEngine,
		llmService:     llmService,
	}
}

// SubmitAssessment processes answers, computes scores, and stores the result.
func (s *AssessmentService) SubmitAssessment(userID uint64, req dto.SubmitAssessmentRequest) (*dto.AssessmentResponse, error) {
	// Fetch all active questions for scoring
	questions, err := s.questionRepo.FindAllActive()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	// Convert to engine-compatible format
	var questionData []engine.QuestionData
	for _, q := range questions {
		weights, err := engine.ParseQuestionWeights(q.Weights)
		if err != nil {
			continue
		}
		questionData = append(questionData, engine.QuestionData{
			ID:           q.ID,
			Category:     q.Category,
			Weights:      weights,
			DisplayOrder: q.DisplayOrder,
		})
	}

	// Run the scoring engine
	result, err := s.scoringEngine.ComputeResult(req.Answers, questionData)
	if err != nil {
		return nil, fmt.Errorf("scoring engine error: %w", err)
	}

	// Optional: Enhance with LLM explanation
	if s.llmService != nil && s.llmService.IsEnabled() {
		explanation, err := s.llmService.GenerateExplanation(result)
		if err == nil && explanation != "" {
			result.AIExplanation = explanation
		}
	}

	// Serialize for storage
	answersJSON, err := json.Marshal(req.Answers)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize answers: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize result: %w", err)
	}

	// Store in database
	assessment, err := s.assessmentRepo.Create(userID, string(answersJSON), string(resultJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to store assessment: %w", err)
	}

	return &dto.AssessmentResponse{
		ID:        assessment.ID,
		UserID:    assessment.UserID,
		Result:    *result,
		CreatedAt: assessment.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetAssessment retrieves a single assessment by ID.
func (s *AssessmentService) GetAssessment(id, userID uint64) (*dto.AssessmentResponse, error) {
	assessment, err := s.assessmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, fmt.Errorf("assessment not found")
	}
	if assessment.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to assessment")
	}

	var result dto.AssessmentResult
	if err := json.Unmarshal([]byte(assessment.Result), &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &dto.AssessmentResponse{
		ID:        assessment.ID,
		UserID:    assessment.UserID,
		Result:    result,
		CreatedAt: assessment.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetUserAssessments retrieves all assessments for a user.
func (s *AssessmentService) GetUserAssessments(userID uint64) ([]dto.AssessmentListItem, error) {
	assessments, err := s.assessmentRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var items []dto.AssessmentListItem
	for _, a := range assessments {
		var result dto.AssessmentResult
		if err := json.Unmarshal([]byte(a.Result), &result); err != nil {
			continue
		}
		items = append(items, dto.AssessmentListItem{
			ID:             a.ID,
			BestCareerPath: result.BestCareerPath,
			RiskLevel:      result.Risk.Level,
			CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		})
	}

	return items, nil
}

// Chat handles an AI chatbot question in the context of an assessment.
func (s *AssessmentService) Chat(userID uint64, req dto.ChatRequest) (string, error) {
	if req.AssessmentID == 0 {
		return "", fmt.Errorf("assessment_id is required")
	}

	assessment, err := s.assessmentRepo.FindByID(req.AssessmentID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch assessment: %w", err)
	}
	if assessment == nil {
		return "", fmt.Errorf("assessment not found")
	}
	if assessment.UserID != userID {
		return "", fmt.Errorf("unauthorized")
	}

	var result dto.AssessmentResult
	if err := json.Unmarshal([]byte(assessment.Result), &result); err != nil {
		return "", fmt.Errorf("failed to parse assessment result: %w", err)
	}

	return s.llmService.Chat(req.Message, &result)
}
