package transaction

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Transaction struct {
	ID                  string    `json:"id"`
	AffiliateID         string    `json:"affiliate_id"`
	ProductID           string    `json:"product_id"`
	ProductAmount       float64   `json:"product_amount"`
	AffiliateCommission float64   `json:"affiliate_commission"`
	CreatedAt           time.Time `json:"created_at"`
}

func (model *Transaction) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4().String())
}
