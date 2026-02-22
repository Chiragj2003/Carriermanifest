package handler

import (
	"net/http"

	"github.com/careermanifest/backend/internal/dto"
	"github.com/careermanifest/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin dashboard endpoints.
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetStats handles GET /api/admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
