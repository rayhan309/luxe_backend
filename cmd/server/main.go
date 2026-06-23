// @title           LUXE eCommerce API
// @version         1.0
// @description     Premium luxury eCommerce REST API
// @termsOfService  http://swagger.io/terms/

// @contact.name   LUXE API Support
// @contact.email  api@luxe.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer <token>"

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/luxe/backend/config"
	deliveryHTTP "github.com/luxe/backend/internal/delivery/http"
	"github.com/luxe/backend/internal/pkg/jwt"
	"github.com/luxe/backend/internal/repository"
	"github.com/luxe/backend/internal/usecase"

	_ "github.com/luxe/backend/docs" // Swagger docs
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	// Load config
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	log.Info().Str("service", cfg.App.Name).Msg("Starting LUXE API...")

	// Connect to MongoDB
	db, err := repository.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer db.Disconnect()
	log.Info().Str("database", cfg.MongoDB.Database).Msg("Connected to MongoDB")

	// Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.ExpiryHours,
		cfg.JWT.RefreshExpHours,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	bannerRepo := repository.NewBannerRepository(db)

	// Initialize use cases
	authUC := usecase.NewAuthUseCase(userRepo, jwtManager)
	productUC := usecase.NewProductUseCase(productRepo, categoryRepo)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	orderUC := usecase.NewOrderUseCase(orderRepo, productRepo, cartRepo, couponRepo)
	cartUC := usecase.NewCartUseCase(cartRepo, productRepo)
	couponUC := usecase.NewCouponUseCase(couponRepo)
	reviewUC := usecase.NewReviewUseCase(reviewRepo, productRepo, orderRepo)
	bannerUC := usecase.NewBannerUseCase(bannerRepo)
	userUC := usecase.NewUserUseCase(userRepo)

	// Initialize handlers
	deps := &deliveryHTTP.Dependencies{
		Auth:     deliveryHTTP.NewAuthHandler(authUC),
		Product:  deliveryHTTP.NewProductHandler(productUC),
		Category: deliveryHTTP.NewCategoryHandler(categoryUC),
		Order:    deliveryHTTP.NewOrderHandler(orderUC),
		Cart:     deliveryHTTP.NewCartHandler(cartUC),
		Review:   deliveryHTTP.NewReviewHandler(reviewUC),
		Coupon:   deliveryHTTP.NewCouponHandler(couponUC),
		Banner:   deliveryHTTP.NewBannerHandler(bannerUC),
		User:     deliveryHTTP.NewUserHandler(userUC),
		Admin:    deliveryHTTP.NewAdminHandler(orderUC, userUC),
		Upload:   deliveryHTTP.NewUploadHandler(cfg.App.UploadDir, cfg.App.MaxFileSize),
		JWT:      jwtManager,
	}

	// Setup router
	router := deliveryHTTP.SetupRouter(deps, cfg.App.FrontendURL)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Server listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exited")
}
