package handler

import (
	"net/http"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AssessmentHandler handles assessment endpoints.
type AssessmentHandler struct {
	assessmentService *service.AssessmentService
}

// NewAssessmentHandler creates a new AssessmentHandler.
func NewAssessmentHandler(assessmentService *service.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{assessmentService: assessmentService}
}

// Submit handles POST /api/assessment
func (h *AssessmentHandler) Submit(c *gin.Context) {
	var req dto.SubmitAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	userID := GetUserID(c)
	resp, err := h.assessmentService.SubmitAssessment(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Assessment processing failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetByID handles GET /api/assessment/:id
func (h *AssessmentHandler) GetByID(c *gin.Context) {
	id, err := GetParamID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid assessment ID"})
		return
	}

	userID := GetUserID(c)
	resp, err := h.assessmentService.GetAssessment(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListByUser handles GET /api/assessment
func (h *AssessmentHandler) ListByUser(c *gin.Context) {
	userID := GetUserID(c)
	items, err := h.assessmentService.GetUserAssessments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if items == nil {
		items = []dto.AssessmentListItem{}
	}

	c.JSON(http.StatusOK, items)
}

// Chat handles POST /api/chat — AI chatbot for follow-up career questions.
func (h *AssessmentHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	userID := GetUserID(c)
	reply, err := h.assessmentService.Chat(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Chat failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ChatResponse{Reply: reply})
}
