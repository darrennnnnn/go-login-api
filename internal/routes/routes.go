package routes

import (
	middlewares "github.com/darrennnnnn/go-login-api/internal/middleware"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, handler *user.Handler) {
	r.POST("/api/auth/login", handler.Login)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/auth/me", handler.Me)
	}
}