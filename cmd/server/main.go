package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/darrennnnnn/go-login-api/internal/auth"
	"github.com/darrennnnnn/go-login-api/internal/health"
	"github.com/darrennnnnn/go-login-api/internal/routes"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/darrennnnnn/go-login-api/internal/validation"
	"github.com/darrennnnnn/go-login-api/pkg/database"
	"github.com/darrennnnnn/go-login-api/pkg/redisclient"
	"github.com/gin-gonic/gin"
)

func main() {
	// environment variables
	cfg := config.Load()

	// database connection
	db := database.Connect(cfg)
	db.AutoMigrate(&user.User{}, &auth.AccessToken{})

	// redis connection
	redisClient := redisclient.Connect(cfg)
	defer redisClient.Close()

	// validation initialization
	validation.Init()

	// user repository
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// auth repository
	authRepo := auth.NewRepository(db, redisClient)
	authService := auth.NewService(authRepo, userRepo, cfg)
	authHandler := auth.NewHandler(authService, userService)

	// health handler
	healthHandler := health.NewHandler(db, redisClient)

	// router initialization
	router := gin.Default()

	// routes registration
	routes.Register(router, routes.Deps{
		User:           userHandler,
		Auth:           authHandler,
		Health:         healthHandler,
		TokenValidator: authService,
		JWTSecret:      cfg.JWT.Secret,
	})

	// server address
	addr := fmt.Sprintf(":%s", cfg.Server.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// server start
	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// server shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}
