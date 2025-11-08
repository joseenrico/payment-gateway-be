package handler

import (
	"net/http"
	"strconv"

	"payment-gateway-manjo/backend/internal/usecase"
	"payment-gateway-manjo/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type QRHandler struct {
	qrUsecase usecase.QRGeneratorUsecase
}

func NewQRHandler(qrUsecase usecase.QRGeneratorUsecase) *QRHandler {
	return &QRHandler{
		qrUsecase: qrUsecase,
	}
}

type GenerateQRRequest struct {
	MerchantID         string `json:"merchantId" binding:"required"`
	PartnerReferenceNo string `json:"partnerReferenceNo" binding:"required"`
	Amount             struct {
		Value    string `json:"value" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	} `json:"amount" binding:"required"`
}

type GenerateQRResponse struct {
	ResponseCode       string `json:"responseCode"`
	ResponseMessage    string `json:"responseMessage"`
	ReferenceNo        string `json:"referenceNo"`
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	QRContent          string `json:"qrContent"`
}

func (h *QRHandler) GenerateQR(c *gin.Context) {
	requestBodyInterface, exists := c.Get("requestBody")
	if !exists {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", "Invalid request body")
		return
	}

	requestBody := requestBodyInterface.(struct {
		MerchantID        string `json:"merchantId"`
		PartnerReferenceNo string `json:"partnerReferenceNo"`
		Amount            struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"amount"`
	})

	amount, err := strconv.ParseFloat(requestBody.Amount.Value, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", "Invalid amount format")
		return
	}

	if amount <= 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", "Amount must be greater than 0")
		return
	}

	transaction, err := h.qrUsecase.GenerateQR(
		requestBody.MerchantID,
		amount,
		requestBody.Amount.Currency,
		requestBody.PartnerReferenceNo,
	)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "Internal Server Error", err.Error())
		return
	}

	qrResponse := GenerateQRResponse{
		ResponseCode:       response.CodeSuccess,
		ResponseMessage:    "Successful",
		ReferenceNo:        transaction.ReferenceNumber,
		PartnerReferenceNo: transaction.PartnerReferenceNumber,
		QRContent:          transaction.QRContent,
	}

	c.JSON(http.StatusOK, qrResponse)
}