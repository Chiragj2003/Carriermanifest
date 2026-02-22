// Package repository implements database access for all entities.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/careermanifest/backend/internal/models"
)

// UserRepository handles user database operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user and returns the created user.
func (r *UserRepository) Create(name, email, passwordHash string) (*models.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, 'user')",
		name, email, passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, _ := result.LastInsertId()
	return r.FindByID(uint64(id))
}

// FindByEmail retrieves a user by email address.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return user, nil
}

// CountUsers returns the total number of users.
func (r *UserRepository) CountUsers() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// CreateAdmin creates an admin user (used for seeding).
func (r *UserRepository) CreateAdmin(name, email, passwordHash string) (*models.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, 'admin') ON DUPLICATE KEY UPDATE role='admin'",
		name, email, passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	id, _ := result.LastInsertId()
	if id == 0 {
		return r.FindByEmail(email)
	}
	return r.FindByID(uint64(id))
}
