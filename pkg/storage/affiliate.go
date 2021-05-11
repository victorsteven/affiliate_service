package storage

import (
	"affiliate_service/cmd/domain/affiliate"
	"github.com/jinzhu/gorm"
)

type affiliateStorage struct {
	db *gorm.DB
}

func NewAffiliateDB(db *gorm.DB) *affiliateStorage {
	return &affiliateStorage{db}
}

type AffiliateStorage interface {
	GetByAffiliateCode(code string) (*affiliate.Affiliate, error)
	UpdateAffiliate(affiliate *affiliate.Affiliate) (*affiliate.Affiliate, error)
}

func (r *affiliateStorage) GetByAffiliateCode(code string) (*affiliate.Affiliate, error) {

	affiliate := &affiliate.Affiliate{}

	err := r.db.Where("affiliate_code = ?", code).Take(&affiliate).Error
	if err != nil {
		return nil, err
	}
	return affiliate, nil
}

func (r *affiliateStorage) UpdateAffiliate(affiliate *affiliate.Affiliate) (*affiliate.Affiliate, error) {

	err := r.db.Save(affiliate).Error
	if err != nil {
		return nil, err
	}
	return affiliate, nil
}
