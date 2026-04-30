package handler

import (
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/config"
)

var (
	Logger *config.Logger // Letra Maiúscula para exportar
)

func InitHandle() {
	Logger = config.GetLogger("handler")
}
