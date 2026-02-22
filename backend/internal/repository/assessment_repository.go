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
	res, err := r.db.Exec(
		"INSERT INTO assessments (user_id, answers, result) VALUES (?, ?, ?)",
		userID, answers, result,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}

	id, _ := res.LastInsertId()
	return r.FindByID(uint64(id))
}

// FindByID retrieves an assessment by ID.
func (r *AssessmentRepository) FindByID(id uint64) (*models.Assessment, error) {
	a := &models.Assessment{}
	err := r.db.QueryRow(
		"SELECT id, user_id, answers, result, created_at FROM assessments WHERE id = ?",
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
		"SELECT id, user_id, answers, result, created_at FROM assessments WHERE user_id = ? ORDER BY created_at DESC",
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
		"SELECT JSON_EXTRACT(result, '$.best_career_path') as career, COUNT(*) as cnt FROM assessments GROUP BY career",
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
		// Remove JSON quotes
		if len(career) >= 2 && career[0] == '"' {
			career = career[1 : len(career)-1]
		}
		dist[career] = count
	}
	return dist, nil
}

// GetRiskDistribution returns the count of each risk level.
func (r *AssessmentRepository) GetRiskDistribution() (map[string]int, error) {
	rows, err := r.db.Query(
		"SELECT JSON_EXTRACT(result, '$.risk.level') as risk_level, COUNT(*) as cnt FROM assessments GROUP BY risk_level",
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
		if len(level) >= 2 && level[0] == '"' {
			level = level[1 : len(level)-1]
		}
		dist[level] = count
	}
	return dist, nil
}
