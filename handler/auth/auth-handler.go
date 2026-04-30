package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/pkg"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/security"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RequestCode(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.SendError(c, http.StatusBadRequest, "Telefone é obrigatório")
			fmt.Printf("Requisição com dados inválidos: %v\n", err)
			fmt.Printf("Phone: %v\n", req.Phone)
			return
		}

		// 1. Pega a lista da variável de ambiente
		allowedPhonesStr := os.Getenv("AUTHORIZED_PHONES")
		allowedPhones := strings.Split(allowedPhonesStr, ",")

		// 2. Verifica se o número está na Whitelist
		isAuthorized := false
		for _, p := range allowedPhones {
			if strings.TrimSpace(p) == req.Phone {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			// Se não estiver na lista, barramos aqui mesmo
			pkg.SendError(c, http.StatusForbidden, "Número não autorizado para acessar métricas")
			fmt.Printf("Tentativa de acesso não autorizada do número: %s\n", req.Phone)
			return
		}

		code := security.GenerateOTP()
		ctx := context.Background()

		err := rdb.Set(ctx, "valid_code:"+code, req.Phone, 20*time.Minute).Err()

		if err != nil {
			pkg.SendError(c, http.StatusInternalServerError, "Erro ao gerar código de autenticação")
			return
		}

		pkg.SendSuccess(c, "code-generated", code)
	}
}

func VerifyCode(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone" binding:"required"`
			Code  string `json:"code" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.SendError(c, http.StatusBadRequest, "Dados invalidos")
			return
		}

		ctx := context.Background()
		storedCode, err := rdb.Get(ctx, "auth"+req.Phone).Result()

		if err == redis.Nil {
			pkg.SendError(c, http.StatusUnauthorized, "Código expirado")
		}

		if storedCode != req.Code {
			pkg.SendError(c, http.StatusUnauthorized, "Código incorreto")
		}

		rdb.Del(ctx, "auth"+req.Phone)

		pkg.SendSuccess(c, "code-verified", "Autenticação bem-sucedida")
	}
}
