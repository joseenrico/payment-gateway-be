package main

import (
	"log"

	"payment-gateway-manjo/backend/internal/delivery/http/handler"
	"payment-gateway-manjo/backend/internal/delivery/http/middleware"
	"payment-gateway-manjo/backend/internal/infrastructure/config"
	"payment-gateway-manjo/backend/internal/infrastructure/database"
	"payment-gateway-manjo/backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	transactionRepo := database.NewTransactionRepository(db)

	qrUsecase := usecase.NewQRGeneratorUsecase(transactionRepo)
	paymentUsecase := usecase.NewPaymentUsecase(transactionRepo)

	qrHandler := handler.NewQRHandler(qrUsecase)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	signatureValidator := middleware.NewSignatureValidator(cfg.Security.SecretKey)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Signature"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	v1 := router.Group("/api/v1")
	{
		qr := v1.Group("/qr")
		{
			qr.POST("/generate", signatureValidator.ValidateQRSignature(), qrHandler.GenerateQR)
			qr.POST("/payment", signatureValidator.ValidatePaymentSignature(), paymentHandler.ProcessPayment)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.GET("", paymentHandler.GetTransactions)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}