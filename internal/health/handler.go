package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHandler(db *gorm.DB, redisClient *redis.Client) *Handler {
	return &Handler{db: db, redis: redisClient}
}

func (h *Handler) Check(ctx *gin.Context) {
	result := gin.H{
		"status": "ok",
	}

	if err := h.db.Exec("SELECT 1").Error; err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "degraded",
			"db":     err.Error(),
		})
		return
	}

	redisCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.redis.Ping(redisCtx).Err(); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "degraded",
			"redis":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
