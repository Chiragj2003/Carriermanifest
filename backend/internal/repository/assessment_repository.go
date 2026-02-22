package repository

import (
	"database/sql"
	"fmt"

	"github.com/careermanifest/backend/internal/models"
)

// AssessmentRepository handles assessment database operations.
type AssessmentRepository struct {
	db *sql.DB
}

// NewAssessmentRepository creates a new AssessmentRepository.
func NewAssessmentRepository(db *sql.DB) *AssessmentRepository {
	return &AssessmentRepository{db: db}
}

// Create stores a new assessment result.
func (r *AssessmentRepository) Create(userID uint64, answers, result string) (*models.Assessment, error) {
	var id uint64
	err := r.db.QueryRow(
		"INSERT INTO assessments (user_id, answers, result) VALUES ($1, $2, $3) RETURNING id",
		userID, answers, result,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}

	return r.FindByID(id)
}

// FindByID retrieves an assessment by ID.
func (r *AssessmentRepository) FindByID(id uint64) (*models.Assessment, error) {
	a := &models.Assessment{}
	err := r.db.QueryRow(
		"SELECT id, user_id, answers, result, created_at FROM assessments WHERE id = $1",
		id,
	).Scan(&a.ID, &a.UserID, &a.Answers, &a.Result, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find assessment: %w", err)
	}
	return a, nil
}

// FindByUserID retrieves all assessments for a user.
func (r *AssessmentRepository) FindByUserID(userID uint64) ([]models.Assessment, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, answers, result, created_at FROM assessments WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query assessments: %w", err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		if err := rows.Scan(&a.ID, &a.UserID, &a.Answers, &a.Result, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan assessment: %w", err)
		}
		assessments = append(assessments, a)
	}
	return assessments, nil
}

// CountAssessments returns the total number of assessments.
func (r *AssessmentRepository) CountAssessments() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM assessments").Scan(&count)
	return count, err
}

// GetCareerDistribution returns the count of each best career path.
func (r *AssessmentRepository) GetCareerDistribution() (map[string]int, error) {
	rows, err := r.db.Query(
		"SELECT result->>'best_career_path' as career, COUNT(*) as cnt FROM assessments GROUP BY career",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dist := make(map[string]int)
	for rows.Next() {
		var career string
		var count int
		if err := rows.Scan(&career, &count); err != nil {
			continue
		}
		dist[career] = count
	}
	return dist, nil
}

// GetRiskDistribution returns the count of each risk level.
func (r *AssessmentRepository) GetRiskDistribution() (map[string]int, error) {
	rows, err := r.db.Query(
		"SELECT result->'risk'->>'level' as risk_level, COUNT(*) as cnt FROM assessments GROUP BY risk_level",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dist := make(map[string]int)
	for rows.Next() {
		var level string
		var count int
		if err := rows.Scan(&level, &count); err != nil {
			continue
		}
		dist[level] = count
	}
	return dist, nil
}
