// internal/handlers/property_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	propertyService *service.PropertyService
}

func NewPropertyHandler(propertyService *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{propertyService: propertyService}
}

// CreatePropertyRequest - Request structure for creating property
type CreatePropertyRequest struct {
	PropertyType    string  `json:"property_type" binding:"required"` // residential, commercial, agricultural
	Address         string  `json:"address" binding:"required"`
	Area            float64 `json:"area" binding:"required"`
	AnnualTaxAmount float64 `json:"annual_tax_amount"`
}

// GetProperties - Get all properties (filtered by ownership for citizens)
func (h *PropertyHandler) GetProperties(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	propertyType := c.Query("property_type")

	filters := map[string]interface{}{}
	
	// If citizen, only show their properties
	if role != "admin" {
		filters["owner_id"] = userID
	}

	// Add property type filter if provided
	if propertyType != "" {
		filters["property_type"] = propertyType
	}

	properties, total, err := h.propertyService.GetProperties(page, limit, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch properties", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Properties retrieved successfully", properties, pagination)
}

// GetProperty - Get single property by ID
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	property, err := h.propertyService.GetProperty(uint(propertyID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Property not found", err.Error())
		return
	}

	// Verify ownership
	if role != "admin" && property.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to view this property")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Property retrieved successfully", property)
}

// CreateProperty - Register a new property
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var req CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Validate property type
	validTypes := map[string]bool{
		"residential":  true,
		"commercial":   true,
		"agricultural": true,
	}
	if !validTypes[req.PropertyType] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property type", "Property type must be residential, commercial, or agricultural")
		return
	}

	// Validate area
	if req.Area <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid area", "Area must be greater than 0")
		return
	}

	property, err := h.propertyService.CreateProperty(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create property", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Property registered successfully", property)
}

// UpdateProperty - Update property details (Admin only)
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	property, err := h.propertyService.UpdateProperty(uint(propertyID), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update property", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Property updated successfully", property)
}

// DeleteProperty - Delete/deactivate property
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	if err := h.propertyService.DeleteProperty(uint(propertyID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete property", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Property deleted successfully", nil)
}

// GetBills - Get all tax bills for a property
func (h *PropertyHandler) GetBills(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	// Verify ownership
	if role != "admin" {
		property, err := h.propertyService.GetProperty(uint(propertyID))
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Property not found", err.Error())
			return
		}
		if property.OwnerID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to view these bills")
			return
		}
	}

	bills, err := h.propertyService.GetPropertyBills(uint(propertyID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch bills", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bills retrieved successfully", bills)
}

// CreateBill - Generate a new tax bill (Admin only)
func (h *PropertyHandler) CreateBill(c *gin.Context) {
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	var req struct {
		FinancialYear string  `json:"financial_year" binding:"required"` // e.g., "2024-25"
		Quarter       string  `json:"quarter" binding:"required"`        // Q1, Q2, Q3, Q4
		TaxAmount     float64 `json:"tax_amount" binding:"required"`
		DueDate       string  `json:"due_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	bill, err := h.propertyService.CreateBill(uint(propertyID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create bill", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Bill created successfully", bill)
}

// MakePayment - Make a payment for a tax bill
func (h *PropertyHandler) MakePayment(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	var req struct {
		BillID        uint    `json:"bill_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"` // online, cash, cheque
		TransactionID string  `json:"transaction_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Validate payment method
	validMethods := map[string]bool{
		"online": true,
		"cash":   true,
		"cheque": true,
		"upi":    true,
	}
	if !validMethods[req.PaymentMethod] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment method", "Payment method must be online, cash, cheque, or upi")
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid amount", "Amount must be greater than 0")
		return
	}

	// Verify ownership
	if role != "admin" {
		property, err := h.propertyService.GetProperty(uint(propertyID))
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Property not found", err.Error())
			return
		}
		if property.OwnerID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to make payment for this property")
			return
		}
	}

	payment, err := h.propertyService.MakePayment(req.BillID, req.Amount, req.PaymentMethod, req.TransactionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payment successful", payment)
}

// GetPaymentHistory - Get payment history
func (h *PropertyHandler) GetPaymentHistory(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	propertyID := c.Query("property_id")
	status := c.Query("status")

	var payments interface{}
	var total int64
	var err error

	if role == "admin" {
		filters := map[string]interface{}{}
		if propertyID != "" {
			filters["property_id"] = propertyID
		}
		if status != "" {
			filters["status"] = status
		}
		payments, total, err = h.propertyService.GetAllPayments(page, limit, filters)
	} else {
		payments, total, err = h.propertyService.GetUserPayments(userID, page, limit)
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payment history", err.Error())
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		TotalItems: total,
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Payment history retrieved successfully", payments, pagination)
}

// GetPayment - Get single payment details
func (h *PropertyHandler) GetPayment(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	paymentID, err := strconv.Atoi(c.Param("paymentId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err.Error())
		return
	}

	payment, err := h.propertyService.GetPayment(uint(paymentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}

	// Verify ownership
	if role != "admin" {
		// Check if payment belongs to user's property
		if payment.Bill.Property.OwnerID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to view this payment")
			return
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment details retrieved successfully", payment)
}

// GetDueBills - Get all due bills for user's properties
func (h *PropertyHandler) GetDueBills(c *gin.Context) {
	userID := c.GetUint("userID")

	bills, err := h.propertyService.GetUserDueBills(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch due bills", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Due bills retrieved successfully", bills)
}

// GetPropertyStats - Get property statistics (Admin only)
func (h *PropertyHandler) GetPropertyStats(c *gin.Context) {
	stats, err := h.propertyService.GetPropertyStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statistics", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}

// GetRevenueReport - Get revenue report (Admin only)
func (h *PropertyHandler) GetRevenueReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	groupBy := c.DefaultQuery("group_by", "month") // month, quarter, year

	report, err := h.propertyService.GetRevenueReport(startDate, endDate, groupBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate report", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revenue report generated successfully", report)
}

// DownloadReceipt - Download payment receipt
func (h *PropertyHandler) DownloadReceipt(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	paymentID, err := strconv.Atoi(c.Param("paymentId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err.Error())
		return
	}

	// Verify ownership
	if role != "admin" {
		payment, err := h.propertyService.GetPayment(uint(paymentID))
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Payment not found", err.Error())
			return
		}
		if payment.Bill.Property.OwnerID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	}

	// Generate receipt PDF (implementation depends on your PDF library)
	receiptURL, err := h.propertyService.GenerateReceipt(uint(paymentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate receipt", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Receipt generated successfully", gin.H{
		"receipt_url": receiptURL,
	})
}

// SendPaymentReminder - Send payment reminder for due bills (Admin only)
func (h *PropertyHandler) SendPaymentReminder(c *gin.Context) {
	propertyID, err := strconv.Atoi(c.Param("propertyId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid property ID", err.Error())
		return
	}

	if err := h.propertyService.SendPaymentReminder(uint(propertyID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to send reminder", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment reminder sent successfully", nil)
}