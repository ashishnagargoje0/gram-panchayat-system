package handlers

import (
	"net/http"
	"strconv"
	"gram-panchayat/internal/service"
	"gram-panchayat/internal/utils"
	
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// InitiatePayment - Initiate payment process
func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		BillID        uint    `json:"bill_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	paymentData, err := h.paymentService.InitiatePayment(userID, req.BillID, req.Amount, req.PaymentMethod)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment initiation failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment initiated successfully", paymentData)
}

// VerifyPayment - Verify payment callback
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		PaymentID     string `json:"payment_id" binding:"required"`
		TransactionID string `json:"transaction_id" binding:"required"`
		Status        string `json:"status" binding:"required"`
		Signature     string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	payment, err := h.paymentService.VerifyPayment(req.PaymentID, req.TransactionID, req.Status, req.Signature)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment verification failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment verified successfully", payment)
}

// GetPaymentStatus - Get payment status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentId")

	payment, err := h.paymentService.GetPaymentByID(paymentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment status retrieved", payment)
}

// RefundPayment - Process refund (Admin)
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	paymentID := c.Param("paymentId")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	refund, err := h.paymentService.ProcessRefund(paymentID, req.Reason)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Refund processing failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Refund processed successfully", refund)
}