package auth

import (
	"errors"
	"net/http"

	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/darrennnnnn/go-login-api/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service     *Service
	userService *user.Service
}

func NewHandler(service *Service, userService *user.Service) *Handler {
	return &Handler{
		service:     service,
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
		if errors.Is(err, user.ErrEmailInUse) || errors.Is(err, user.ErrUsernameInUse) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user.ToUserResponse(createdUser)})
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

	accessToken, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (h *Handler) Logout(ctx *gin.Context) {
	tokenID, ok := tokenIDFromClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.Logout(tokenID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) Me(ctx *gin.Context) {
	claims, ok := claimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       claims["user_id"],
		"username": claims["username"],
		"email":    claims["email"],
		"exp":      claims["exp"],
	})
}

func claimsFromContext(ctx *gin.Context) (jwt.MapClaims, bool) {
	claimsAny, ok := ctx.Get("claims")
	if !ok {
		return nil, false
	}

	claims, ok := claimsAny.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

func tokenIDFromClaims(ctx *gin.Context) (string, bool) {
	claims, ok := claimsFromContext(ctx)
	if !ok {
		return "", false
	}

	tokenID, ok := claims["token_id"].(string)
	if !ok || tokenID == "" {
		return "", false
	}

	return tokenID, true
}
