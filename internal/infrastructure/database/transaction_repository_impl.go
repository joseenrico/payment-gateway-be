package database

import (
	"payment-gateway-manjo/backend/internal/domain/entity"
	"payment-gateway-manjo/backend/internal/domain/repository"

	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Create(transaction *entity.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepositoryImpl) FindByReferenceNumber(referenceNumber string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Where("reference_number = ?", referenceNumber).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) FindByPartnerReferenceNumber(partnerRefNo string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Where("partner_reference_number = ?", partnerRefNo).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) Update(transaction *entity.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepositoryImpl) FindAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepositoryImpl) FindByFilters(merchantID, partnerRefNo, refNo, status string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := r.db.Model(&entity.Transaction{})

	if merchantID != "" {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if partnerRefNo != "" {
		query = query.Where("partner_reference_number = ?", partnerRefNo)
	}
	if refNo != "" {
		query = query.Where("reference_number = ?", refNo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}