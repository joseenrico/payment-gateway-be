package handler

import (
	"net/http"
	"strconv"

	"payment-gateway-manjo/backend/internal/usecase"
	"payment-gateway-manjo/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

type PaymentNotificationRequest struct {
	OriginalReferenceNo        string `json:"originalReferenceNo" binding:"required"`
	OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo" binding:"required"`
	TransactionStatusDesc      string `json:"transactionStatusDesc" binding:"required"`
	PaidTime                   string `json:"paidTime" binding:"required"`
	Amount                     struct {
		Value    string `json:"value" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	} `json:"amount" binding:"required"`
}

type PaymentNotificationResponse struct {
	ResponseCode          string `json:"responseCode"`
	ResponseMessage       string `json:"responseMessage"`
	TransactionStatusDesc string `json:"transactionStatusDesc"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	requestBodyInterface, exists := c.Get("requestBody")
	if !exists {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", "Invalid request body")
		return
	}

	requestBody := requestBodyInterface.(struct {
		OriginalReferenceNo        string `json:"originalReferenceNo"`
		OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo"`
		TransactionStatusDesc      string `json:"transactionStatusDesc"`
		PaidTime                   string `json:"paidTime"`
		Amount                     struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"amount"`
	})

	amount, err := strconv.ParseFloat(requestBody.Amount.Value, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", "Invalid amount format")
		return
	}

	transaction, err := h.paymentUsecase.ProcessPayment(
		requestBody.OriginalReferenceNo,
		amount,
		requestBody.TransactionStatusDesc,
		requestBody.PaidTime,
	)

	if err != nil {
		if err.Error() == "transaction not found" {
			response.Error(c, http.StatusNotFound, response.CodeNotFound, "Transaction Not Found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "Internal Server Error", err.Error())
		return
	}

	paymentResponse := PaymentNotificationResponse{
		ResponseCode:          response.CodePaymentSuccess,
		ResponseMessage:       "Successful",
		TransactionStatusDesc: transaction.Status,
	}

	c.JSON(http.StatusOK, paymentResponse)
}

func (h *PaymentHandler) GetTransactions(c *gin.Context) {
	merchantID := c.Query("merchantId")
	partnerRefNo := c.Query("partnerReferenceNo")
	refNo := c.Query("referenceNo")
	status := c.Query("status")
	transactions, err := h.paymentUsecase.GetTransactions(merchantID, partnerRefNo, refNo, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "Internal Server Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2000000",
		"responseMessage": "Successful",
		"data":            transactions,
	})
}