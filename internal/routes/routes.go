package routes

import (
	"github.com/darrennnnnn/go-login-api/internal/auth"
	"github.com/darrennnnnn/go-login-api/internal/health"
	middlewares "github.com/darrennnnnn/go-login-api/internal/middleware"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, userHandler *user.Handler, authHandler *auth.Handler, healthHandler *health.Handler, jwtSecret []byte) {
	r.GET("/health", healthHandler.Check)

	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/register", authHandler.Register)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(jwtSecret, authHandler.Service))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/logout", authHandler.Logout)
		protected.GET("/user", userHandler.GetUsers)
		protected.GET("/user/:id", userHandler.GetUserByID)
		protected.DELETE("/user/:id", userHandler.DeleteUser)
	}
}
