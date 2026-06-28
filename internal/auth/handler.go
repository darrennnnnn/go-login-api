package auth

import (
	"net/http"

	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/darrennnnnn/go-login-api/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	Service     *Service
	userService *user.Service
}

func NewHandler(service *Service, userService *user.Service) *Handler {
	return &Handler{
		Service:     service,
		userService: userService,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"errors": validation.FormatValidationErrors(err),
		})
		return
	}

	createdUser, err := h.userService.CreateUser(user.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err.Error() {
		case "Email already in use", "Username already in use":
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": createdUser})
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"errors": validation.FormatValidationErrors(err),
		})
		return
	}

	accessToken, err := h.Service.Login(req)
	if err != nil {
		switch err.Error() {
		case "User not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "Unauthorized":
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusAccepted, gin.H{"access_token": accessToken})
}

func (h *Handler) Logout(ctx *gin.Context) {
	claimsAny, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := claimsAny.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenID, ok := claims["id"].(string)
	if !ok || tokenID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.Service.Logout(tokenID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) Me(ctx *gin.Context) {
	claimsAny, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := claimsAny.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       claims["id"],
		"username": claims["username"],
		"email":    claims["email"],
		"exp":      claims["exp"],
	})
}
