package integration_test

import (
	"affiliate_service/cmd/domain/affiliate"
	"affiliate_service/cmd/domain/product"
	"affiliate_service/cmd/handler"
	"affiliate_service/pkg/storage"
	"affiliate_service/pkg/storage/seed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initTest()
	code := m.Run()
	os.Exit(code)
}

var (
	server          *httptest.Server
	client          *http.Client
	seededProduct   *product.Product
	seededAffiliate *affiliate.Affiliate
)

func initTest() {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		panic(err)
	}

	databaseName := os.Getenv("TEST_DB_NAME")
	databaseUser := os.Getenv("TEST_DB_USER")
	databaseHost := os.Getenv("TEST_DB_HOST")
	databasePort := os.Getenv("TEST_DB_PORT")
	databasePassword := os.Getenv("TEST_DB_PASSWORD")

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	o := storage.OpenDB(dbConn)

	affiliateStorage := storage.NewAffiliateDB(o)
	productStorage := storage.NewProductStorage(o)
	transactionStorage := storage.NewTransactionStorage(o)

	//Handlers
	transaction := handler.NewTransactionHandler(affiliateStorage, productStorage, transactionStorage)

	//seed database
	seededProduct, seededAffiliate, err = seed.Load(o)
	if err != nil {
		log.Fatal("Error seeding data: ", err)
	}

	r := gin.Default()

	r.GET("/purchase-product/:productID/:affiliateCode", transaction.PurchaseProduct)

	server = httptest.NewTLSServer(r)

	cc, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = server.Client()
	server.Client().Jar = cc
}
