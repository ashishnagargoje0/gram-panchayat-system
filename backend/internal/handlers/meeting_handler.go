package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)

type MeetingHandler struct {
	meetingService *service.MeetingService
}

func NewMeetingHandler(meetingService *service.MeetingService) *MeetingHandler {
	return &MeetingHandler{meetingService: meetingService}
}

func (h *MeetingHandler) GetMeetings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	meetings, total, err := h.meetingService.GetMeetings(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch meetings", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Meetings retrieved", meetings, pagination)
}

func (h *MeetingHandler) GetMeeting(c *gin.Context) {
	meetingID, _ := strconv.Atoi(c.Param("id"))

	meeting, err := h.meetingService.GetMeeting(uint(meetingID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Meeting not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Meeting retrieved", meeting)
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		MeetingType string `json:"meeting_type"`
		ScheduledAt string `json:"scheduled_at" binding:"required"`
		Location    string `json:"location"`
		Agenda      string `json:"agenda"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	meeting, err := h.meetingService.CreateMeeting(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create meeting", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Meeting scheduled", meeting)
}

func (h *MeetingHandler) AddMinutes(c *gin.Context) {
	adminID := c.GetUint("userID")
	meetingID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Content   string `json:"content" binding:"required"`
		Attendees string `json:"attendees"`
		Decisions string `json:"decisions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	minutes, err := h.meetingService.AddMinutes(uint(meetingID), adminID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add minutes", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Minutes added", minutes)
}

// internal/handlers/dashboard_handler.go
type DashboardHandler struct {
	userService        *service.UserService
	applicationService *service.ApplicationService
	complaintService   *service.ComplaintService
}

func NewDashboardHandler(userService *service.UserService, applicationService *service.ApplicationService, complaintService *service.ComplaintService) *DashboardHandler {
	return &DashboardHandler{
		userService:        userService,
		applicationService: applicationService,
		complaintService:   complaintService,
	}
}

func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {
	stats, err := h.applicationService.GetAdminStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dashboard stats", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard stats retrieved", stats)
}

func (h *DashboardHandler) GetCitizenDashboard(c *gin.Context) {
	userID := c.GetUint("userID")

	stats, err := h.applicationService.GetCitizenStats(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dashboard stats", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard stats retrieved", stats)
}