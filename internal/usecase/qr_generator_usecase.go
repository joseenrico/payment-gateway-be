package usecase

import (
	"errors"
	"fmt"
	"time"

	"payment-gateway-manjo/backend/internal/domain/entity"
	"payment-gateway-manjo/backend/internal/domain/repository"

	"github.com/google/uuid"
)

type QRGeneratorUsecase interface {
	GenerateQR(merchantID string, amount float64, currency, partnerRefNo string) (*entity.Transaction, error)
}

type qrGeneratorUsecase struct {
	transactionRepo repository.TransactionRepository
}

func NewQRGeneratorUsecase(transactionRepo repository.TransactionRepository) QRGeneratorUsecase {
	return &qrGeneratorUsecase{
		transactionRepo: transactionRepo,
	}
}

func (u *qrGeneratorUsecase) GenerateQR(merchantID string, amount float64, currency, partnerRefNo string) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	referenceNumber := generateReferenceNumber()
	qrContent := generateQRContent(merchantID, referenceNumber, amount)

	transaction := &entity.Transaction{
		MerchantID:             merchantID,
		Amount:                 amount,
		Currency:               currency,
		PartnerReferenceNumber: partnerRefNo,
		ReferenceNumber:        referenceNumber,
		Status:                 entity.StatusPending,
		TransactionDate:        time.Now(),
		QRContent:              qrContent,
	}

	if err := u.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func generateReferenceNumber() string {
	uuid := uuid.New().String()
	shortUUID := uuid[:10]
	return "A" + shortUUID
}

func generateQRContent(merchantID, referenceNumber string, amount float64) string {
    amountStr := fmt.Sprintf("%.2f", amount)
    return fmt.Sprintf("00020101021226620015ID.CO.MANJO.WWW01189360085801751859910210%s0303UMI51530014ID.CO.QRIS.WWW0215ID102106515192304121.0.21.09.255204481653033605502015802ID5904OLDI6013JAKARTA BARAT61051147062454%02d%s62460525%s07031110806ASPI663040FAD",
        merchantID, len(amountStr), amountStr, referenceNumber)
}