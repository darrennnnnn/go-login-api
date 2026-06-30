package routes

import (
	"github.com/darrennnnnn/go-login-api/internal/auth"
	"github.com/darrennnnnn/go-login-api/internal/health"
	"github.com/darrennnnnn/go-login-api/internal/middleware"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	User           *user.Handler
	Auth           *auth.Handler
	Health         *health.Handler
	TokenValidator auth.TokenValidator
	JWTSecret      []byte
}

func Register(r *gin.Engine, deps Deps) {
	r.GET("/health", deps.Health.Check)

	auth.RegisterPublicRoutes(r, deps.Auth)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(deps.JWTSecret, deps.TokenValidator))
	auth.RegisterProtectedRoutes(protected, deps.Auth)
	user.RegisterRoutes(protected, deps.User)
}
