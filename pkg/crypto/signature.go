package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"	
)

func GenerateSignature(data string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateSignature(data string, signature string, secretKey string) bool {
	expectedSignature := GenerateSignature(data, secretKey)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func CreateSignatureString(parts ...string) string {
	return strings.Join(parts, "|")
}

func GenerateQRSignatureString(merchantID, amount, partnerRefNo string) string {
	return fmt.Sprintf("%s|%s|%s", merchantID, amount, partnerRefNo)
}

func GeneratePaymentSignatureString(referenceNo, amount, status string) string {
	return fmt.Sprintf("%s|%s|%s", referenceNo, amount, status)
}