package user

import (
	"net/http"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
	config *config.Config
}

func NewHandler(service *Service, config *config.Config) *Handler {
	return &Handler{
		service: service,
		config: config,
	}
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	accessToken, err := h.service.Login(req)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"access_token": accessToken,
	})
}

func (h *Handler) Me(ctx *gin.Context) {
	claimsAny, ok := ctx.Get("claims")

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	claims, ok := claimsAny.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": claims["id"],
		"username": claims["username"],
		"name": claims["name"],
		"exp": claims["exp"],
	})
}