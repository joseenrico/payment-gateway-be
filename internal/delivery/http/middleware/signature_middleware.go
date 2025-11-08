package middleware

import (
	"net/http"

	"payment-gateway-manjo/backend/pkg/crypto"
	"payment-gateway-manjo/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SignatureValidator struct {
	secretKey string
}

func NewSignatureValidator(secretKey string) *SignatureValidator {
	return &SignatureValidator{
		secretKey: secretKey,
	}
}

func (sv *SignatureValidator) ValidateQRSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedSignature := c.GetHeader("X-Signature")
		if receivedSignature == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "Unauthorized", "Missing signature")
			c.Abort()
			return
		}

		var requestBody struct {
			MerchantID        string `json:"merchantId"`
			PartnerReferenceNo string `json:"partnerReferenceNo"`
			Amount            struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", err.Error())
			c.Abort()
			return
		}

		signatureString := crypto.GenerateQRSignatureString(
			requestBody.MerchantID,
			requestBody.Amount.Value,
			requestBody.PartnerReferenceNo,
		)

		if !crypto.ValidateSignature(signatureString, receivedSignature, sv.secretKey) {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "Unauthorized", "Invalid signature")
			c.Abort()
			return
		}

		c.Set("requestBody", requestBody)
		c.Next()
	}
}

func (sv *SignatureValidator) ValidatePaymentSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedSignature := c.GetHeader("X-Signature")
		if receivedSignature == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "Unauthorized", "Missing signature")
			c.Abort()
			return
		}

		var requestBody struct {
			OriginalReferenceNo        string `json:"originalReferenceNo"`
			OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo"`
			TransactionStatusDesc      string `json:"transactionStatusDesc"`
			PaidTime                   string `json:"paidTime"`
			Amount                     struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "Bad Request", err.Error())
			c.Abort()
			return
		}

		signatureString := crypto.GeneratePaymentSignatureString(
			requestBody.OriginalReferenceNo,
			requestBody.Amount.Value,
			requestBody.TransactionStatusDesc,
		)

		if !crypto.ValidateSignature(signatureString, receivedSignature, sv.secretKey) {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "Unauthorized", "Invalid signature")
			c.Abort()
			return
		}

		c.Set("requestBody", requestBody)
		c.Next()
	}
}