package main

import (
	"affiliate_service/cmd/handler"
	"affiliate_service/pkg/storage"
	"affiliate_service/pkg/storage/seed"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	databaseName := os.Getenv("DATABASE_NAME")
	databaseUser := os.Getenv("DATABASE_USER")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	appAddr := ":" + os.Getenv("PORT")

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	// Repositories
	var affiliateStorage storage.AffiliateStorage
	var productStorage storage.ProductStorage
	var transactionStorage storage.TransactionStorage

	o := storage.OpenDB(dbConn)

	affiliateStorage = storage.NewAffiliateDB(o)
	productStorage = storage.NewProductStorage(o)
	transactionStorage = storage.NewTransactionStorage(o)

	//Handlers
	product := handler.NewTransactionHandler(affiliateStorage, productStorage, transactionStorage)

	//seed database
	_, _, err := seed.Load(o)
	if err != nil {
		log.Fatal("Error seeding data: ", err)
	}

	r := gin.Default()

	r.GET("/purchase-product/:productID/:affiliateCode", product.PurchaseProduct)

	srv := &http.Server{
		Addr:    appAddr,
		Handler: r,
	}
	go func() {
		//service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
