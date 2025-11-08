package usecase

import (
	"errors"
	"fmt"
	"time"

	"payment-gateway-manjo/backend/internal/domain/entity"
	"payment-gateway-manjo/backend/internal/domain/repository"

	"gorm.io/gorm"
)

type PaymentUsecase interface {
	ProcessPayment(referenceNo string, amount float64, status, paidTime string) (*entity.Transaction, error)
	GetTransactions(merchantID, partnerRefNo, refNo, status string) ([]entity.Transaction, error)
}

type paymentUsecase struct {
	transactionRepo repository.TransactionRepository
}

func NewPaymentUsecase(transactionRepo repository.TransactionRepository) PaymentUsecase {
	return &paymentUsecase{
		transactionRepo: transactionRepo,
	}
}

func (u *paymentUsecase) ProcessPayment(referenceNo string, amount float64, status, paidTime string) (*entity.Transaction, error) {
	transaction, err := u.transactionRepo.FindByReferenceNumber(referenceNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	if transaction.Amount != amount {
		return nil, errors.New("amount mismatch")
	}

	parsedPaidTime, err := time.Parse(time.RFC3339, paidTime)
	if err != nil {
		parsedPaidTime = time.Now()
	}

	transaction.Status = status
	transaction.PaidDate = &parsedPaidTime

	if err := u.transactionRepo.Update(transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}
	return transaction, nil
}

func (u *paymentUsecase) GetTransactions(merchantID, partnerRefNo, refNo, status string) ([]entity.Transaction, error) {
	if merchantID == "" && partnerRefNo == "" && refNo == "" && status == "" {
		return u.transactionRepo.FindAll()
	}
	return u.transactionRepo.FindByFilters(merchantID, partnerRefNo, refNo, status)
}