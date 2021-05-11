package storage

import (
	"affiliate_service/cmd/domain/transaction"
	"github.com/jinzhu/gorm"
)

type transactionStorage struct {
	db *gorm.DB
}

func NewTransactionStorage(db *gorm.DB) *transactionStorage {
	return &transactionStorage{db}
}

type TransactionStorage interface {
	CreateTransaction(transaction *transaction.Transaction) (*transaction.Transaction, error)
	TransactionRollBack(id string) error
}

func (r *transactionStorage) CreateTransaction(transaction *transaction.Transaction) (*transaction.Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *transactionStorage) TransactionRollBack(id string) error {
	err := r.db.Unscoped().Where("id = ?", id).Delete(&transaction.Transaction{}).Error
	if err != nil {
		return err
	}
	return nil
}
