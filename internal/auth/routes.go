package auth

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(r *gin.Engine, h *Handler) {
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/register", h.Register)
}

func RegisterProtectedRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/auth/me", h.Me)
	rg.POST("/auth/logout", h.Logout)
}
