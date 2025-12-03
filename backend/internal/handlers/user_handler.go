package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers - List all users (Admin)
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")
	isActive := c.Query("is_active")
	search := c.Query("search")

	filters := map[string]interface{}{}
	if role != "" {
		filters["role"] = role
	}
	if isActive != "" {
		filters["is_active"] = isActive == "true"
	}

	users, total, err := h.userService.GetUsers(page, limit, filters, search)
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

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Users retrieved successfully", users, pagination)
}

// GetUser - Get single user (Admin)
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUser - Update user details (Admin)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Don't allow password updates through this endpoint
	delete(updates, "password")

	user, err := h.userService.UpdateUser(uint(userID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser - Delete user (Admin)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	// Prevent self-deletion
	currentUserID := c.GetUint("userID")
	if uint(userID) == currentUserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot delete your own account", "")
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete user", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ActivateDeactivateUser - Activate or deactivate user account (Admin)
func (h *UserHandler) ActivateDeactivateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.userService.SetUserActiveStatus(uint(userID), req.IsActive); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update user status", err.Error())
		return
	}

	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	utils.SuccessResponse(c, http.StatusOK, "User "+status+" successfully", nil)
}

// ChangeUserRole - Change user role (Admin)
func (h *UserHandler) ChangeUserRole(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"admin":   true,
		"citizen": true,
		"staff":   true,
	}
	if !validRoles[req.Role] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role", "Role must be admin, citizen, or staff")
		return
	}

	if err := h.userService.ChangeUserRole(uint(userID), req.Role); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to change user role", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User role changed successfully", nil)
}

// GetUserStats - Get user statistics (Admin)
func (h *UserHandler) GetUserStats(c *gin.Context) {
	stats, err := h.userService.GetUserStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statistics", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}

// GetUserActivity - Get user activity log (Admin)
func (h *UserHandler) GetUserActivity(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	activities, err := h.userService.GetUserActivity(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch activity", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User activity retrieved successfully", activities)
}

// ResetUserPassword - Reset user password (Admin)
func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.userService.ResetPassword(uint(userID), req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to reset password", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}