package repository

import (
	"database/sql"
	"fmt"

	"github.com/careermanifest/backend/internal/models"
)

// QuestionRepository handles question database operations.
type QuestionRepository struct {
	db *sql.DB
}

// NewQuestionRepository creates a new QuestionRepository.
func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// Create inserts a new question.
func (r *QuestionRepository) Create(category, text, options, weights string, order int) (*models.Question, error) {
	var id uint64
	err := r.db.QueryRow(
		"INSERT INTO questions (category, question_text, options, weights, display_order) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		category, text, options, weights, order,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	return r.FindByID(id)
}

// FindByID retrieves a question by ID.
func (r *QuestionRepository) FindByID(id uint64) (*models.Question, error) {
	q := &models.Question{}
	err := r.db.QueryRow(
		"SELECT id, category, question_text, options, weights, display_order, is_active, created_at, updated_at FROM questions WHERE id = $1",
		id,
	).Scan(&q.ID, &q.Category, &q.QuestionText, &q.Options, &q.Weights, &q.DisplayOrder, &q.IsActive, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find question: %w", err)
	}
	return q, nil
}

// FindAllActive retrieves all active questions ordered by display_order.
func (r *QuestionRepository) FindAllActive() ([]models.Question, error) {
	rows, err := r.db.Query(
		"SELECT id, category, question_text, options, weights, display_order, is_active, created_at, updated_at FROM questions WHERE is_active = TRUE ORDER BY display_order ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.Category, &q.QuestionText, &q.Options, &q.Weights, &q.DisplayOrder, &q.IsActive, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// FindAll retrieves all questions (admin).
func (r *QuestionRepository) FindAll() ([]models.Question, error) {
	rows, err := r.db.Query(
		"SELECT id, category, question_text, options, weights, display_order, is_active, created_at, updated_at FROM questions ORDER BY display_order ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.Category, &q.QuestionText, &q.Options, &q.Weights, &q.DisplayOrder, &q.IsActive, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// Update modifies a question.
func (r *QuestionRepository) Update(id uint64, category, text, options, weights string, order int, isActive bool) error {
	_, err := r.db.Exec(
		"UPDATE questions SET category=$1, question_text=$2, options=$3, weights=$4, display_order=$5, is_active=$6, updated_at=NOW() WHERE id=$7",
		category, text, options, weights, order, isActive, id,
	)
	return err
}

// CountQuestions returns the total question count.
func (r *QuestionRepository) CountQuestions() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM questions WHERE is_active = TRUE").Scan(&count)
	return count, err
}
