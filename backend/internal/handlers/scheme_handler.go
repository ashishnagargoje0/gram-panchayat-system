package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)
type SchemeHandler struct {
	schemeService *service.SchemeService
}

func NewSchemeHandler(schemeService *service.SchemeService) *SchemeHandler {
	return &SchemeHandler{schemeService: schemeService}
}

type CreateSchemeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date"`
}

// GetSchemes - List all active schemes
func (h *SchemeHandler) GetSchemes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	isActive := c.DefaultQuery("is_active", "true")

	filters := map[string]interface{}{
		"is_active": isActive == "true",
	}
	if category != "" {
		filters["category"] = category
	}

	schemes, total, err := h.schemeService.GetSchemes(page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch schemes", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Schemes retrieved successfully", schemes, pagination)
}

// GetScheme - Get single scheme
func (h *SchemeHandler) GetScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err.Error())
		return
	}

	scheme, err := h.schemeService.GetScheme(uint(schemeID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Scheme not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Scheme retrieved successfully", scheme)
}

// CreateScheme - Create new scheme (Admin)
func (h *SchemeHandler) CreateScheme(c *gin.Context) {
	var req CreateSchemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	scheme, err := h.schemeService.CreateScheme(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create scheme", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Scheme created successfully", scheme)
}

// UpdateScheme - Update scheme (Admin)
func (h *SchemeHandler) UpdateScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	scheme, err := h.schemeService.UpdateScheme(uint(schemeID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update scheme", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Scheme updated successfully", scheme)
}

// DeleteScheme - Delete scheme (Admin)
func (h *SchemeHandler) DeleteScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err.Error())
		return
	}

	if err := h.schemeService.DeleteScheme(uint(schemeID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete scheme", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Scheme deleted successfully", nil)
}

// ApplyForScheme - Apply for a scheme
func (h *SchemeHandler) ApplyForScheme(c *gin.Context) {
	userID := c.GetUint("userID")
	schemeID, err := strconv.Atoi(c.Param("schemeId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid scheme ID", err.Error())
		return
	}

	var req struct {
		FormData map[string]interface{} `json:"form_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	application, err := h.schemeService.ApplyForScheme(userID, uint(schemeID), req.FormData)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to apply for scheme", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Application submitted successfully", application)
}

// GetMySchemeApplications - Get user's scheme applications
func (h *SchemeHandler) GetMySchemeApplications(c *gin.Context) {
	userID := c.GetUint("userID")

	applications, err := h.schemeService.GetUserApplications(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch applications", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Applications retrieved successfully", applications)
}

// UpdateSchemeApplication - Update application status (Admin)
func (h *SchemeHandler) UpdateSchemeApplication(c *gin.Context) {
	adminID := c.GetUint("userID")
	applicationID, err := strconv.Atoi(c.Param("applicationId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid application ID", err.Error())
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	application, err := h.schemeService.UpdateApplicationStatus(uint(applicationID), adminID, req.Status, req.Remarks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update application", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Application updated successfully", application)
}