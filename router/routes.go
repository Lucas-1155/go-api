package router

import (
	"time"

	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/config"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/handler"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/handler/auth"
	queryreq "github.com/Comal-Developers1/oracle-grafana-API.git/go-api/handler/query-req"
	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	handler.InitHandle()
	dbPostgres := config.GetPostgres()
	dbOracle := config.GetOracle()
	logQuery := config.NewLogger("QUERY-REQ")
	rdb := config.GetRedis()

	v1 := router.Group("/api/v1")
	{
		// Rota de Teste **
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "online",
				"message": "API da Comal operando corretamente",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		// Rotas de Query **

		v1.GET("/admin/getqueries/:alias", queryreq.GetMetrics(dbPostgres, rdb, dbOracle))

		v1.GET("/admin/getqueries", queryreq.GetMetrics(dbPostgres, rdb, dbOracle))

		v1.GET("/admin/queries/getSql", queryreq.ListQueries(dbPostgres, rdb, logQuery))

		v1.POST("/admin/queries", queryreq.RegisterQuery(dbPostgres, logQuery, rdb))

		// Rotas de Autenticação **
		v1.POST("/auth/request", auth.RequestCode(rdb))

		v1.POST("/auth/verify", auth.VerifyCode(rdb))

	}
}
