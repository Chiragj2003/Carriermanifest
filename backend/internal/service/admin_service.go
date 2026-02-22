package service

import (
	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/repository"
)

// AdminService handles admin dashboard business logic.
type AdminService struct {
	userRepo       *repository.UserRepository
	assessmentRepo *repository.AssessmentRepository
	questionRepo   *repository.QuestionRepository
}

// NewAdminService creates a new AdminService.
func NewAdminService(
	userRepo *repository.UserRepository,
	assessmentRepo *repository.AssessmentRepository,
	questionRepo *repository.QuestionRepository,
) *AdminService {
	return &AdminService{
		userRepo:       userRepo,
		assessmentRepo: assessmentRepo,
		questionRepo:   questionRepo,
	}
}

// GetStats returns aggregate platform statistics.
func (s *AdminService) GetStats() (*dto.AdminStatsResponse, error) {
	totalUsers, err := s.userRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	totalAssessments, err := s.assessmentRepo.CountAssessments()
	if err != nil {
		return nil, err
	}

	careerDist, err := s.assessmentRepo.GetCareerDistribution()
	if err != nil {
		careerDist = make(map[string]int)
	}

	riskDist, err := s.assessmentRepo.GetRiskDistribution()
	if err != nil {
		riskDist = make(map[string]int)
	}

	totalQuestions, err := s.questionRepo.CountQuestions()
	if err != nil {
		totalQuestions = 0
	}

	return &dto.AdminStatsResponse{
		TotalUsers:         totalUsers,
		TotalAssessments:   totalAssessments,
		TotalQuestions:     totalQuestions,
		CareerDistribution: careerDist,
		RiskDistribution:   riskDist,
	}, nil
}
