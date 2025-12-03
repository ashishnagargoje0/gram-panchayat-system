package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)

type ComplaintHandler struct {
	complaintService *service.ComplaintService
}

func NewComplaintHandler(complaintService *service.ComplaintService) *ComplaintHandler {
	return &ComplaintHandler{complaintService: complaintService}
}

type CreateComplaintRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Location    string   `json:"location"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Priority    string   `json:"priority"`
}

// CreateComplaint - Register a new complaint
func (h *ComplaintHandler) CreateComplaint(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var req CreateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Validate category
	validCategories := []string{"infrastructure", "water", "electricity", "sanitation", "road", "street_light", "garbage", "other"}
	isValidCategory := false
	for _, cat := range validCategories {
		if req.Category == cat {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category", "Please select a valid category")
		return
	}

	complaint, err := h.complaintService.CreateComplaint(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create complaint", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Complaint registered successfully", complaint)
}

// GetComplaints - List all complaints (filtered by user for citizens)
func (h *ComplaintHandler) GetComplaints(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	category := c.Query("category")
	priority := c.Query("priority")

	filters := map[string]interface{}{}
	if status != "" {
		filters["status"] = status
	}
	if category != "" {
		filters["category"] = category
	}
	if priority != "" {
		filters["priority"] = priority
	}

	// If citizen, only show their complaints
	if role != "admin" {
		filters["user_id"] = userID
	}

	complaints, total, err := h.complaintService.GetComplaints(page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch complaints", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Complaints retrieved successfully", complaints, pagination)
}

// GetComplaint - Get single complaint details
func (h *ComplaintHandler) GetComplaint(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	complaintID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid complaint ID", err.Error())
		return
	}

	complaint, err := h.complaintService.GetComplaint(uint(complaintID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Complaint not found", err.Error())
		return
	}

	// Check permission
	if role != "admin" && complaint.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Complaint retrieved successfully", complaint)
}

// UpdateComplaint - Update complaint status and details (Admin)
func (h *ComplaintHandler) UpdateComplaint(c *gin.Context) {
	adminID := c.GetUint("userID")
	complaintID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid complaint ID", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Add admin info
	updates["assigned_to"] = adminID

	complaint, err := h.complaintService.UpdateComplaint(uint(complaintID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update complaint", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Complaint updated successfully", complaint)
}

// AddComment - Add comment to complaint
func (h *ComplaintHandler) AddComment(c *gin.Context) {
	userID := c.GetUint("userID")
	complaintID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid complaint ID", err.Error())
		return
	}

	var req struct {
		Comment    string `json:"comment" binding:"required"`
		IsInternal bool   `json:"is_internal"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	comment, err := h.complaintService.AddComment(uint(complaintID), userID, req.Comment, req.IsInternal)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add comment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Comment added successfully", comment)
}

// AssignComplaint - Assign complaint to staff (Admin)
func (h *ComplaintHandler) AssignComplaint(c *gin.Context) {
	complaintID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		AssignedTo uint `json:"assigned_to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.complaintService.AssignComplaint(uint(complaintID), req.AssignedTo); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to assign complaint", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Complaint assigned successfully", nil)
}

// GetComplaintStats - Get complaint statistics (Admin)
func (h *ComplaintHandler) GetComplaintStats(c *gin.Context) {
	stats, err := h.complaintService.GetComplaintStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statistics", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}