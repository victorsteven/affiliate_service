package seed

import (
	"affiliate_service/cmd/domain/affiliate"
	"affiliate_service/cmd/domain/product"
	"affiliate_service/cmd/domain/transaction"
	"affiliate_service/cmd/domain/user"
	"affiliate_service/pkg/helpers"

	"github.com/jinzhu/gorm"
)

var users = []user.User{
	{
		Name:  "Steven victor",
		Email: "steven@gmail.com",
	},
	{
		Name:  "Martin Luther",
		Email: "luther@gmail.com",
	},
}

var products = []product.Product{
	{
		Name:   "Iphone 12",
		Amount: 284.30,
	},
	{
		Name:   "Tesla",
		Amount: 1050.25,
	},
}

var affiliates = []affiliate.Affiliate{
	{
		Name:          "Emma Jones",
		Email:         "affilate1@gmail.com",
		Commission:    0,
		AffiliateCode: helpers.GenStr(6),
	},
	{
		Name:          "Greg Vivian",
		Email:         "affilate2@gmail.com",
		Commission:    0,
		AffiliateCode: helpers.GenStr(6),
	},
}

func Load(db *gorm.DB) (*product.Product, *affiliate.Affiliate, error) {

	var returnProduct product.Product
	var returnAffiliate affiliate.Affiliate

	err := db.Debug().DropTableIfExists(&transaction.Transaction{}, &affiliate.Affiliate{}, &product.Product{}, &user.User{}).Error
	if err != nil {
		return nil, nil, err
	}
	err = db.Debug().AutoMigrate(&user.User{}, &product.Product{}, &affiliate.Affiliate{}, &transaction.Transaction{}).Error
	if err != nil {
		return nil, nil, err
	}

	for i, _ := range users {
		err = db.Debug().Model(&user.User{}).Create(&users[i]).Error
		if err != nil {
			return nil, nil, err
		}
		products[i].UserID = users[i].ID

		err = db.Debug().Model(&product.Product{}).Create(&products[i]).Error
		if err != nil {
			return nil, nil, err
		}

		returnProduct = products[i]
	}

	for _, v := range affiliates {
		err = db.Debug().Create(&v).Error
		if err != nil {
			return nil, nil, err
		}
		returnAffiliate = v
	}

	return &returnProduct, &returnAffiliate, nil
}
