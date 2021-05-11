package storage

import (
	"affiliate_service/cmd/domain/affiliate"
	"affiliate_service/cmd/domain/product"
	"affiliate_service/cmd/domain/transaction"
	"affiliate_service/cmd/domain/user"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

func OpenDB(database string) *gorm.DB {
	db, err := gorm.Open("postgres", database)
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := Automigrate(db); err != nil {
		panic(err)
	}
	return db
}

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(user.User{}, product.Product{}, affiliate.Affiliate{}, transaction.Transaction{}).Error
}
