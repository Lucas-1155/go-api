package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendError envia uma resposta de erro padronizada
func SendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"message": msg,
	})
}

// SendSuccess envia uma resposta de sucesso padronizada
func SendSuccess(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("operation from handle: %s successful", op),
		"data":    data,
	})
}
