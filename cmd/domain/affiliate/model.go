package affiliate

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Affiliate struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Commission    float64   `json:"commission"`
	AffiliateCode string    `json:"affiliate_code"`
	CreatedAt     time.Time `json:"created_at"`
}

func (model *Affiliate) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4().String())
}
