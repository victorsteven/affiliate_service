package storage

import (
	"affiliate_service/cmd/domain/product"
	"github.com/jinzhu/gorm"
)

type productStorage struct {
	db *gorm.DB
}

func NewProductStorage(db *gorm.DB) *productStorage {
	return &productStorage{db}
}

type ProductStorage interface {
	GetByID(id string) (*product.Product, error)
}

func (r *productStorage) GetByID(id string) (*product.Product, error) {

	product := &product.Product{}

	err := r.db.Where("id = ?", id).Take(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}
