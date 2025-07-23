package repositories

import (
	"gorm.io/gorm"
	"oe02_go_tam/models"
)

type TransactionRepository interface {
	Create(tx *models.PaymentTransaction) error
	FindByTxnRef(txnRef string) (*models.PaymentTransaction, error)
	Update(tx *models.PaymentTransaction) error
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db}
}

func (r *transactionRepositoryImpl) Create(tx *models.PaymentTransaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepositoryImpl) FindByTxnRef(txnRef string) (*models.PaymentTransaction, error) {
	var tx models.PaymentTransaction
	err := r.db.Where("txn_ref = ?", txnRef).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepositoryImpl) Update(tx *models.PaymentTransaction) error {
	return r.db.Save(tx).Error
}
