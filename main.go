package main

import (
	"github.com/darrennnnnn/go-login-api/domains/users/handlers/http"
	"github.com/darrennnnnn/go-login-api/domains/users/repositories"
	"github.com/darrennnnnn/go-login-api/domains/users/usecases"
	"github.com/darrennnnnn/go-login-api/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.POST("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hwello",
		})
	})

	router.GET("/:id/:whatever", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"id": ctx.Param("id"),
			"whatever": ctx.Param("whatever"),
		})
	})

	userRepo := repositories.NewUserRepository()
	userUc := usecases.NewUserUseCase(userRepo)
	userHttp := http.NewUserHttp(userUc)

	router.POST("/api/auth/login", userHttp.Login)

	protected := router.Group("")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/api/protected/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
	protected.GET("/api/auth/me", userHttp.Me)

	router.Run(":8080")
}