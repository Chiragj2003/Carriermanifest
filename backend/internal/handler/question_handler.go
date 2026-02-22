package handler

import (
	"net/http"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// QuestionHandler handles question endpoints.
type QuestionHandler struct {
	questionService *service.QuestionService
}

// NewQuestionHandler creates a new QuestionHandler.
func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

// GetActiveQuestions handles GET /api/questions (for assessment form)
func (h *QuestionHandler) GetActiveQuestions(c *gin.Context) {
	questions, err := h.questionService.GetActiveQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if questions == nil {
		questions = []dto.QuestionDTO{}
	}

	c.JSON(http.StatusOK, questions)
}

// GetAllQuestions handles GET /api/admin/questions (admin)
func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	questions, err := h.questionService.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if questions == nil {
		questions = []dto.QuestionDTO{}
	}

	c.JSON(http.StatusOK, questions)
}

// CreateQuestion handles POST /api/admin/questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	question, err := h.questionService.CreateQuestion(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// UpdateQuestion handles PUT /api/admin/questions/:id
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id, err := GetParamID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid question ID"})
		return
	}

	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	if err := h.questionService.UpdateQuestion(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Question updated successfully"})
}
