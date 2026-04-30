package queryreq

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/config"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/pkg"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/schemas"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/security"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ListQueries(db *gorm.DB, rdb *redis.Client, log *config.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")

		if code == "" {
			pkg.SendError(c, http.StatusBadRequest, "Código de verificação é obrigatório")
			return
		}

		ctx := context.Background()
		redisKey := fmt.Sprintf("valid_code:%s", code)

		// 1. Valida o código no Redis
		_, err := rdb.Get(ctx, redisKey).Result()
		if err == redis.Nil {
			pkg.SendError(c, http.StatusUnauthorized, "Código inválido ou expirado")
			return
		} else if err != nil {
			log.ErrorF("Redis error: %v", err)
			pkg.SendError(c, http.StatusInternalServerError, "Erro interno")
			return
		}

		// 2. Busca os dados no Postgres
		var encryptedQueries []schemas.QueryMetadata
		if err := db.Find(&encryptedQueries).Error; err != nil {
			log.ErrorF("Database error: %v", err)
			pkg.SendError(c, http.StatusInternalServerError, "Erro ao buscar queries")
			return
		}

		// 3. Descriptografia dos SQLs
		// Criamos um slice de resposta para não enviar campos desnecessários do GORM
		type QueryResponse struct {
			ID         string `json:"id"`
			Alias      string `json:"alias"`
			SQL        string `json:"sql"`
			TTLSeconds int    `json:"ttl_seconds"`
		}

		var response []QueryResponse

		for _, q := range encryptedQueries {
			decryptedSQL, err := security.Decrypt(q.SQLEncrypted)
			if err != nil {
				log.ErrorF("Erro ao descriptografar query %d: %v", q.ID, err)
				decryptedSQL = "Erro ao descriptografar este SQL"
			}

			response = append(response, QueryResponse{
				ID:         q.ID.String(),
				Alias:      q.Alias,
				SQL:        decryptedSQL,
				TTLSeconds: q.TTLSeconds,
			})
		}

		// 4. Retorna a lista com o SQL aberto
		pkg.SendSuccess(c, "queries-list", response)
	}
}
