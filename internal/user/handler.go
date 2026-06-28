package user

import (
	"net/http"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/darrennnnnn/go-login-api/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Handler struct {
	Service *Service
	Config *config.Config
}

func NewHandler(service *Service, config *config.Config) *Handler {
	return &Handler{
		Service: service,
		Config: config,
	}
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"errors": validation.FormatValidationErrors(err),
		})
		return
	}

	user, err := h.Service.CreateUser(req)

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

	ctx.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if err := h.Service.DeleteUser(id); err != nil {
		if err.Error() == "User not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *Handler) GetUsers(ctx *gin.Context) {
	users, err := h.Service.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (h *Handler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	user, err := h.Service.GetUserByID(id)

	if err != nil {
		if err.Error() == "User not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
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

	ctx.JSON(http.StatusAccepted, gin.H{
		"access_token": accessToken,
	})
}

func (h *Handler) Logout(ctx *gin.Context) {
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

	tokenID, ok := claims["id"].(string)
	if !ok || tokenID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := h.Service.Logout(tokenID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
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
		"email": claims["email"],
		"exp": claims["exp"],
	})
}