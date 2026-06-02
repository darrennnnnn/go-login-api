package http

import (
	"net/http"

	users "github.com/darrennnnnn/go-login-api/domains"
	"github.com/darrennnnnn/go-login-api/domains/users/models/requests"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-playground/validator/v10"
)

type UserHttp struct {
	uc users.UserUseCase
}

func NewUserHttp(uc users.UserUseCase) *UserHttp {
	return &UserHttp{uc: uc}
}

func (handler *UserHttp) Login(ctx *gin.Context) {
		requestBody := &requests.LoginRequest{}
		
		err := ctx.ShouldBindJSON(&requestBody)
		
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		validate := validator.New()
		err = validate.Struct(requestBody)

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		accessToken, err := handler.uc.Login(requestBody)

		if err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"access_token": accessToken,
		})
}

func (handler *UserHttp) Me(ctx *gin.Context) {
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
		"name":     claims["name"],
		"exp":      claims["exp"],
	})
}