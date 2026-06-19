package main

import (
	"fmt"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/darrennnnnn/go-login-api/internal/routes"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/darrennnnnn/go-login-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect()
	db.AutoMigrate(&user.User{})

	repo := user.NewRepository(db)
	service := user.NewService(repo, cfg)
	handler := user.NewHandler(service, cfg)

	router := gin.Default()

	routes.Register(router, handler)

	addr := fmt.Sprintf(":%s", cfg.Server.ServerPort)
	router.Run(addr)
}