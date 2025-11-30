package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)
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

// internal/handlers/user_handler.go
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")

	filters := map[string]interface{}{}
	if role != "" {
		filters["role"] = role
	}

	users, total, err := h.userService.GetUsers(page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Users retrieved", users, pagination)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(userID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete user", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted", nil)
}