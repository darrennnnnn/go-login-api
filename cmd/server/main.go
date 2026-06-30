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
	cfg := config.Load()

	db := database.Connect(cfg)
	db.AutoMigrate(&user.User{}, &auth.AccessToken{})

	redisClient := redisclient.Connect(cfg)
	defer redisClient.Close()

	validation.Init()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(db, redisClient)
	authService := auth.NewService(authRepo, userRepo, cfg)
	authHandler := auth.NewHandler(authService, userService)

	healthHandler := health.NewHandler(db, redisClient)

	router := gin.Default()

	routes.Register(router, userHandler, authHandler, authService, healthHandler, cfg.JWT.Secret)

	addr := fmt.Sprintf(":%s", cfg.Server.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

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
