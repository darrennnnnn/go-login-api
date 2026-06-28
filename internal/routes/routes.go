package routes

import (
	middlewares "github.com/darrennnnnn/go-login-api/internal/middleware"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, handler *user.Handler, jwtSecret []byte) {
	r.POST("/api/auth/login", handler.Login)
	r.POST("/api/auth/register", handler.CreateUser)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		protected.GET("/auth/me", handler.Me)
		protected.GET("/user", handler.GetUsers)
		protected.GET("/user/:id", handler.GetUserByID)
		protected.DELETE("/user/:id", handler.DeleteUser)
	}
}