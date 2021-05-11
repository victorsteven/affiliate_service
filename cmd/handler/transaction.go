package handler

import (
	"affiliate_service/cmd/domain/transaction"
	"affiliate_service/pkg/storage"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

type transactionHandler struct {
	a storage.AffiliateStorage
	p storage.ProductStorage
	t storage.TransactionStorage
}

func NewTransactionHandler(a storage.AffiliateStorage, p storage.ProductStorage, t storage.TransactionStorage) *transactionHandler {
	return &transactionHandler{a, p, t}
}

func (h *transactionHandler) PurchaseProduct(c *gin.Context) {

	productID := c.Param("productID")
	affiliateCode := c.Param("affiliateCode")

	prod, err := h.p.GetByID(productID)
	if err != nil {
		log.Println("error fetching product by id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product provided"})
		return
	}
	affl, err := h.a.GetByAffiliateCode(affiliateCode)
	if err != nil {
		log.Println("error fetching affiliate by code: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid affiliate provided"})
		return
	}

	// we are assuming that the affiliate has 10% cut
	affiliateCommission := (prod.Amount * 10) / 100

	ts := &transaction.Transaction{
		ID:                  uuid.NewV4().String(),
		AffiliateID:         affl.ID,
		ProductID:           prod.ID,
		ProductAmount:       prod.Amount,
		AffiliateCommission: affiliateCommission,
		CreatedAt:           time.Time{},
	}

	// add to the transaction table and update the affiliate with his commission
	createdTransaction, err := h.t.CreateTransaction(ts)
	if err != nil {
		log.Println("error saving transaction on product purchased: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, try again later"})
		return
	}

	affl.Commission += affiliateCommission

	_, err = h.a.UpdateAffiliate(affl)
	if err != nil {
		log.Println("error updating affiliate record: ", err)
		errDel := h.t.TransactionRollBack(createdTransaction.ID)
		if errDel != nil || err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, try again later"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, try again later"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
	return
}
