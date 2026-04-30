package queryreq

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/config"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/pkg"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/schemas"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/security"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RegisterQueryRequest struct {
	Alias      string `json:"alias" binding:"required"`
	SQL        string `json:"sql" binding:"required"`
	TTLSeconds int    `json:"ttl_seconds"`
	Code       int    `json:"code" binding:"required"`
}

func RegisterQuery(db *gorm.DB, log *config.Logger, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterQueryRequest
		var newQuery schemas.QueryMetadata

		if err := c.ShouldBindJSON(&req); err != nil {
			log.ErrorF("Invalid request payload: %v", err)
			pkg.SendError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		ctx := context.Background()

		redisKey := fmt.Sprintf("valid_code:%d", req.Code)

		_, err := rdb.Get(ctx, redisKey).Result()

		if err == redis.Nil {
			pkg.SendError(c, http.StatusUnauthorized, "Código invalido ou expirado")
			return
		} else if err != nil {
			log.ErrorF("Redis error: %v", err)
			pkg.SendError(c, http.StatusInternalServerError, "Internal validation error")
			return
		}

		encrypted, err := security.Encrypt(req.SQL)

		if err != nil {
			log.ErrorF("Encryption error: %v", err)
			pkg.SendError(c, http.StatusInternalServerError, "Failed to encrypt SQL")
			return
		}

		upperAlias := strings.ToUpper(req.Alias)
		lowerAlias := strings.ToLower(req.Alias)
		var existingQuery schemas.QueryMetadata

		err = db.Where("alias = ?", upperAlias).First(&existingQuery).Error
		if err == nil {
			pkg.SendError(c, http.StatusConflict, fmt.Sprintf("O nome '%s' já esta em uso", lowerAlias))
			return
		}

		newQuery = schemas.QueryMetadata{
			Alias:        upperAlias,
			SQLEncrypted: encrypted,
			TTLSeconds:   req.TTLSeconds,
		}

		if err := db.Create(&newQuery).Error; err != nil {
			log.ErrorF("Database error: %v", err)
			pkg.SendError(c, http.StatusInternalServerError, "Failed to create query")
			return
		}

		pkg.SendSuccess(c, "query-insert", gin.H{
			"id": newQuery.ID,
		})
	}
}
