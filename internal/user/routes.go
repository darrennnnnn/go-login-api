package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/user", h.GetUsers)
	rg.GET("/user/:id", h.GetUserByID)
	rg.DELETE("/user/:id", h.DeleteUser)
}
