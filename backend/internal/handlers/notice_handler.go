package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)
type NoticeHandler struct {
	noticeService *service.NoticeService
}

func NewNoticeHandler(noticeService *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{noticeService: noticeService}
}

type CreateNoticeRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	IsPublished bool   `json:"is_published"`
}

// GetNotices - List all published notices
func (h *NoticeHandler) GetNotices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	role := c.GetString("role")

	filters := map[string]interface{}{}
	
	// Non-admin users can only see published notices
	if role != "admin" {
		filters["is_published"] = true
	}
	
	if category != "" {
		filters["category"] = category
	}

	notices, total, err := h.noticeService.GetNotices(page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch notices", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Notices retrieved successfully", notices, pagination)
}

// GetNotice - Get single notice
func (h *NoticeHandler) GetNotice(c *gin.Context) {
	noticeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notice ID", err.Error())
		return
	}

	notice, err := h.noticeService.GetNotice(uint(noticeID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Notice not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notice retrieved successfully", notice)
}

// CreateNotice - Create new notice (Admin)
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	adminID := c.GetUint("userID")
	
	var req CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	notice, err := h.noticeService.CreateNotice(adminID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create notice", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Notice created successfully", notice)
}

// UpdateNotice - Update notice (Admin)
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	noticeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notice ID", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	notice, err := h.noticeService.UpdateNotice(uint(noticeID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update notice", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notice updated successfully", notice)
}

// DeleteNotice - Delete notice (Admin)
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	noticeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notice ID", err.Error())
		return
	}

	if err := h.noticeService.DeleteNotice(uint(noticeID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete notice", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notice deleted successfully", nil)
}

// PublishNotice - Publish/unpublish notice (Admin)
func (h *NoticeHandler) PublishNotice(c *gin.Context) {
	noticeID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		IsPublished bool `json:"is_published"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.noticeService.PublishNotice(uint(noticeID), req.IsPublished); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update publish status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notice publish status updated", nil)
}
