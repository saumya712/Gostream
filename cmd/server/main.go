package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"gostream/internal/application"
	"gostream/internal/domain"
	"gostream/internal/infrastructure/auth"
	"gostream/internal/infrastructure/config"
	"gostream/internal/infrastructure/processor"
	"gostream/internal/infrastructure/repository/memory"
	"gostream/internal/infrastructure/repository/postgres"
	"gostream/internal/infrastructure/worker"
	"gostream/internal/transport/handlers"
	"gostream/internal/transport/middleware"
)

func main() {
	// Root context with graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer dbpool.Close()

	// ==========================
	// AUTH WIRING
	// ==========================

	userRepo := postgres.NewUserRepository(dbpool)
	hasher := auth.Newbcrypthasher(cfg.BcryptCost)
	jwtManager := auth.Newjwtmanager(cfg.JWTSecret, cfg.JWTDuration)
	streamHandler := handlers.NewStreamHandler(memory.NewMetricsStore())

	authService := application.NewAuthService(
		userRepo,
		hasher,
		jwtManager,
	)

	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// ==========================
	// LOG SYSTEM WIRING
	// ==========================

	metricsStore := memory.NewMetricsStore()
	workerProcessor := processor.NewWorkerProcessor(metricsStore)

	pool := workerpool.NewPool(
		ctx,
		5,   // worker count
		100, // queue buffer size
		workerProcessor,
	)

	logService := application.NewLogService(pool)

	logHandler := handlers.NewLogHandler(logService)
	metricsHandler := handlers.NewMetricsHandler(metricsStore)

	// ==========================
	// ROUTER SETUP
	// ==========================

	router := gin.Default()


	// Public auth routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Service routes (only SERVICE role can send logs)
	serviceGroup := router.Group("/service")
	serviceGroup.Use(authMiddleware.RequireRole(domain.RoleService))
	{
		serviceGroup.POST("/logs", logHandler.Ingest)
	}

	// Admin routes (only ADMIN role can view metrics)
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMiddleware.RequireRole(domain.Roleadmin))
	{
		adminGroup.GET("/metrics", metricsHandler.GetMetrics)
		adminGroup.GET("/stream", streamHandler.Stream)
	}

	// ==========================
	// HTTP SERVER
	// ==========================

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Println("server running on port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error:", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server shutdown failed:", err)
	}

	pool.Shutdown()

	log.Println("server exited properly")
}
