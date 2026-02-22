// Package service implements business logic for CareerManifest.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/careermanifest/backend/internal/config"
	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

// Register creates a new user account.
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password with bcrypt (cost 12 for production security)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(req.Name, req.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetProfile retrieves the user profile by ID.
func (s *AuthService) GetProfile(userID uint64) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// generateToken creates a signed JWT with user claims.
func (s *AuthService) generateToken(userID uint64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// googleTokenInfo represents the response from Google's tokeninfo endpoint.
type googleTokenInfo struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
}

// GoogleLogin verifies a Google credential token and creates or logs in the user.
func (s *AuthService) GoogleLogin(req dto.GoogleLoginRequest) (*dto.AuthResponse, error) {
	// Verify the Google ID token via Google's tokeninfo endpoint
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.Credential)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Google token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Google response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid Google token")
	}

	var tokenInfo googleTokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse Google token info: %w", err)
	}

	// Verify the token audience matches our client ID
	if s.cfg.GoogleClientID != "" && tokenInfo.Aud != s.cfg.GoogleClientID {
		return nil, errors.New("Google token audience mismatch")
	}

	// Verify email is verified
	if tokenInfo.EmailVerified != "true" {
		return nil, errors.New("Google email not verified")
	}

	// Check if user already exists
	user, err := s.userRepo.FindByEmail(tokenInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if user == nil {
		// Create new user with a random password hash (Google users don't use passwords)
		randomBytes := make([]byte, 32)
		if _, err := rand.Read(randomBytes); err != nil {
			return nil, fmt.Errorf("failed to generate random password: %w", err)
		}
		randomPassword := hex.EncodeToString(randomBytes)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), 12)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		// Use Google name or email prefix as display name
		name := tokenInfo.Name
		if name == "" {
			name = tokenInfo.GivenName
		}
		if name == "" {
			name = tokenInfo.Email
		}

		user, err = s.userRepo.Create(name, tokenInfo.Email, string(hashedPassword))
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// SeedAdmin creates the default admin user if it doesn't exist.
func (s *AuthService) SeedAdmin() error {
	existing, _ := s.userRepo.FindByEmail(s.cfg.AdminEmail)
	if existing != nil {
		return nil // Admin already exists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminPassword), 12)
	if err != nil {
		return err
	}

	_, err = s.userRepo.CreateAdmin("Admin", s.cfg.AdminEmail, string(hashedPassword))
	return err
}
