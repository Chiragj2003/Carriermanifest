// Package models defines the domain entities for CareerManifest.
package models

import "time"

// User represents a registered user.
type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Assessment represents a completed career assessment.
type Assessment struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Answers   string    `json:"answers"`  // JSON string
	Result    string    `json:"result"`   // JSON string
	CreatedAt time.Time `json:"created_at"`
}

// Question represents an assessment question with scoring weights.
type Question struct {
	ID           uint64    `json:"id"`
	Category     string    `json:"category"`
	QuestionText string    `json:"question_text"`
	Options      string    `json:"options"` // JSON array of option strings
	Weights      string    `json:"weights"` // JSON array of weight objects
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
