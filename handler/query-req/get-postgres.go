package queryreq

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/pkg"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/schemas"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/security"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Estrutura para o cache no Redis que suporta Stale-While-Revalidate
type CacheEnvelope struct {
	Data       interface{} `json:"data"`
	ValidUntil time.Time   `json:"valid_until"`
}

func GetMetrics(db *gorm.DB, rdb *redis.Client, oracleDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")

		// Se o Grafana bater aqui sem alias (ou com um alias de teste)
		if alias == "" {
			c.JSON(200, gin.H{"status": "OK", "message": "Conectado ao backend da Comal"})
			return
		}

		ctx := c.Request.Context()

		aliasUpper := strings.ToUpper(alias)

		// 1. Tenta buscar no Redis
		cacheKey := "metrics:" + aliasUpper
		cachedData, err := rdb.Get(ctx, cacheKey).Result()

		if err == nil {
			var envelope CacheEnvelope
			json.Unmarshal([]byte(cachedData), &envelope)

			// Lógica Stale-While-Revalidate
			if time.Now().After(envelope.ValidUntil) {
				// DADO ESTÁ VELHO (STALE): Dispara atualização em background
				go refreshCache(aliasUpper, db, rdb, oracleDB)
			}

			// Retorna o dado imediatamente (seja ele fresco ou stale)
			c.JSON(http.StatusOK, envelope.Data)
			return
		}

		// 2. MISS: Se não tem no cache, faz o processo síncrono (primeira vez)
		data, err := refreshCache(aliasUpper, db, rdb, oracleDB)
		if err != nil {
			pkg.SendError(c, http.StatusInternalServerError, "Erro ao processar consulta")
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// Função auxiliar que faz o trabalho sujo de atualizar o Oracle -> Redis
func refreshCache(alias string, db *gorm.DB, rdb *redis.Client, oracleDB *gorm.DB) (interface{}, error) {
	var queryMeta schemas.QueryMetadata
	ctx := context.Background()

	// Busca o SQL criptografado no Postgres
	if err := db.Where("alias = ?", alias).First(&queryMeta).Error; err != nil {
		return nil, err
	}

	// Descriptografa o SQL
	sql, _ := security.Decrypt(queryMeta.SQLEncrypted)

	// Executa no Oracle (SELECT puro)
	var result []map[string]interface{}
	if err := oracleDB.Raw(sql).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Prepara o envelope de cache
	ttl := time.Duration(queryMeta.TTLSeconds) * time.Second
	envelope := CacheEnvelope{
		Data:       result,
		ValidUntil: time.Now().Add(ttl),
	}

	// Salva no Redis (com uma expiração real maior para permitir o Stale)
	jsonData, _ := json.Marshal(envelope)
	rdb.Set(ctx, "metrics:"+alias, jsonData, ttl+(60*time.Second))

	return result, nil
}
