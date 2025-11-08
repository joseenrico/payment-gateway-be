package repository

import "payment-gateway-manjo/backend/internal/domain/entity"

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	FindByReferenceNumber(referenceNumber string) (*entity.Transaction, error)
	FindByPartnerReferenceNumber(partnerRefNo string) (*entity.Transaction, error)
	Update(transaction *entity.Transaction) error
	FindAll() ([]entity.Transaction, error)
	FindByFilters(merchantID, partnerRefNo, refNo, status string) ([]entity.Transaction, error)
}