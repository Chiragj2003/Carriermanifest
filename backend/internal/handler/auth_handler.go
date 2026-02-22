// Package handler implements HTTP handlers for all API routes.
package handler

import (
	"net/http"
	"strconv"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			status = http.StatusConflict
		}
		c.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Profile handles GET /api/auth/profile
func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint64("user_id")
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GoogleLogin handles POST /api/auth/google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	resp, err := h.authService.GoogleLogin(req)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "invalid Google token" || err.Error() == "Google token audience mismatch" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserID extracts user_id from Gin context (set by auth middleware).
func GetUserID(c *gin.Context) uint64 {
	id, _ := c.Get("user_id")
	if uid, ok := id.(uint64); ok {
		return uid
	}
	return 0
}

// GetParamID extracts an ID from URL params.
func GetParamID(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
