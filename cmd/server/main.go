package main

import (
	"fmt"

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

	db := database.Connect()
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

	routes.Register(router, userHandler, authHandler, healthHandler, cfg.JWT.Secret)

	addr := fmt.Sprintf(":%s", cfg.Server.ServerPort)
	router.Run(addr)
}
