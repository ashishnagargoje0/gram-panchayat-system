// internal/handlers/application_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService *service.ApplicationService
}

func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

type CreateApplicationRequest struct {
	Type     string                 `json:"type" binding:"required"` // birth, death, income, caste, residence, marriage
	FormData map[string]interface{} `json:"form_data" binding:"required"`
	Priority string                 `json:"priority"`
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	application, err := h.applicationService.CreateApplication(userID, req.Type, req.FormData, req.Priority)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create application", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Application submitted successfully", application)
}

func (h *ApplicationHandler) GetUserApplications(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	applicationType := c.Query("type")

	filters := map[string]interface{}{}
	if status != "" {
		filters["status"] = status
	}
	if applicationType != "" {
		filters["type"] = applicationType
	}

	applications, total, err := h.applicationService.GetUserApplications(userID, page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch applications", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Applications retrieved", applications, pagination)
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	applicationID, _ := strconv.Atoi(c.Param("id"))

	application, err := h.applicationService.GetApplication(uint(applicationID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Application not found", err.Error())
		return
	}

	// Check if user has permission to view this application
	if role != "admin" && application.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Application retrieved", application)
}

func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	userID := c.GetUint("userID")
	
	applications, err := h.applicationService.GetAllUserApplications(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch applications", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Applications retrieved", applications)
}

func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	adminID := c.GetUint("userID")
	applicationID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Status  string `json:"status" binding:"required"` // approved, rejected, under_review
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	application, err := h.applicationService.UpdateApplicationStatus(uint(applicationID), adminID, req.Status, req.Remarks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update application", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Application status updated", application)
}