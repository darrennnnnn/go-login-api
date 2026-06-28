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
	service *Service
	config *config.Config
}

func NewHandler(service *Service, config *config.Config) *Handler {
	return &Handler{
		service: service,
		config: config,
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

	user, err := h.service.CreateUser(req)

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

	if err := h.service.DeleteUser(id); err != nil {
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
	users, err := h.service.GetUsers()

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

	user, err := h.service.GetUserByID(id)

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

	accessToken, err := h.service.Login(req)

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