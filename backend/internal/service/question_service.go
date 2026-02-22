package service

import (
	"encoding/json"
	"fmt"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/repository"
)

// QuestionService handles question business logic.
type QuestionService struct {
	questionRepo *repository.QuestionRepository
}

// NewQuestionService creates a new QuestionService.
func NewQuestionService(questionRepo *repository.QuestionRepository) *QuestionService {
	return &QuestionService{questionRepo: questionRepo}
}

// GetActiveQuestions returns all active questions for the assessment form.
func (s *QuestionService) GetActiveQuestions() ([]dto.QuestionDTO, error) {
	questions, err := s.questionRepo.FindAllActive()
	if err != nil {
		return nil, err
	}

	var result []dto.QuestionDTO
	for _, q := range questions {
		var options []dto.QuestionOption
		if err := json.Unmarshal([]byte(q.Options), &options); err != nil {
			continue
		}
		result = append(result, dto.QuestionDTO{
			ID:           q.ID,
			Category:     q.Category,
			QuestionText: q.QuestionText,
			Options:      options,
			DisplayOrder: q.DisplayOrder,
		})
	}

	return result, nil
}

// GetAllQuestions returns all questions with weights (admin).
func (s *QuestionService) GetAllQuestions() ([]dto.QuestionDTO, error) {
	questions, err := s.questionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []dto.QuestionDTO
	for _, q := range questions {
		var options []dto.QuestionOption
		if err := json.Unmarshal([]byte(q.Options), &options); err != nil {
			continue
		}
		var weights []dto.QuestionWeight
		if err := json.Unmarshal([]byte(q.Weights), &weights); err != nil {
			continue
		}
		result = append(result, dto.QuestionDTO{
			ID:           q.ID,
			Category:     q.Category,
			QuestionText: q.QuestionText,
			Options:      options,
			Weights:      weights,
			DisplayOrder: q.DisplayOrder,
			IsActive:     &q.IsActive,
		})
	}

	return result, nil
}

// CreateQuestion creates a new question (admin).
func (s *QuestionService) CreateQuestion(req dto.CreateQuestionRequest) (*dto.QuestionDTO, error) {
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize options: %w", err)
	}

	weightsJSON, err := json.Marshal(req.Weights)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize weights: %w", err)
	}

	q, err := s.questionRepo.Create(req.Category, req.QuestionText, string(optionsJSON), string(weightsJSON), req.DisplayOrder)
	if err != nil {
		return nil, err
	}

	return &dto.QuestionDTO{
		ID:           q.ID,
		Category:     q.Category,
		QuestionText: q.QuestionText,
		Options:      req.Options,
		Weights:      req.Weights,
		DisplayOrder: q.DisplayOrder,
	}, nil
}

// UpdateQuestion updates an existing question (admin).
func (s *QuestionService) UpdateQuestion(id uint64, req dto.UpdateQuestionRequest) error {
	existing, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("question not found")
	}

	category := existing.Category
	if req.Category != "" {
		category = req.Category
	}

	text := existing.QuestionText
	if req.QuestionText != "" {
		text = req.QuestionText
	}

	options := existing.Options
	if req.Options != nil {
		optJSON, _ := json.Marshal(req.Options)
		options = string(optJSON)
	}

	weights := existing.Weights
	if req.Weights != nil {
		wJSON, _ := json.Marshal(req.Weights)
		weights = string(wJSON)
	}

	order := existing.DisplayOrder
	if req.DisplayOrder != 0 {
		order = req.DisplayOrder
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return s.questionRepo.Update(id, category, text, options, weights, order, isActive)
}
