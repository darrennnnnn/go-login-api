package routes

import (
	"github.com/darrennnnnn/go-login-api/internal/auth"
	"github.com/darrennnnnn/go-login-api/internal/health"
	"github.com/darrennnnnn/go-login-api/internal/middleware"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	userHandler *user.Handler,
	authHandler *auth.Handler,
	tokenValidator auth.TokenValidator,
	healthHandler *health.Handler,
	jwtSecret []byte,
) {
	r.GET("/health", healthHandler.Check)

	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/register", authHandler.Register)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtSecret, tokenValidator))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/user", userHandler.GetUsers)
		protected.GET("/user/:id", userHandler.GetUserByID)
		protected.DELETE("/user/:id", userHandler.DeleteUser)
	}
}
