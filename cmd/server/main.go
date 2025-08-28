package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"flypro/internal/config"
	"flypro/internal/handlers"
	"flypro/internal/middleware"
	"flypro/internal/repository"
	"flypro/internal/services"
	"flypro/internal/validators"
)

func main() {
	cfg := config.Load()

	// Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// CORS
	c := cors.DefaultConfig()
	if cfg.CORSAllowedOrigins == "*" {
		c.AllowAllOrigins = true
	} else {
		c.AllowOrigins = cfg.GetCORSOrigins()
	}
	c.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
	r.Use(cors.New(c))

	// DB + Redis
	db := repository.MustOpenGorm(cfg)
	redis := services.MustOpenRedis(cfg)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Services
	currencySvc := services.NewCurrencyService(cfg, redis)
	userSvc := services.NewUserService(userRepo, redis)
	expenseSvc := services.NewExpenseService(expenseRepo, currencySvc, redis)
	reportSvc := services.NewReportService(reportRepo, expenseRepo, currencySvc, redis)

	// Validators
	validators.RegisterCustomValidators()

	// Routes
	h := handlers.NewHandler(userSvc, expenseSvc, reportSvc)
	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:           ":" + cfg.AppPort,
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("listening on :%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
